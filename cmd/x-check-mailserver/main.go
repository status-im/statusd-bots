package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	stdsignal "os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/status-im/go-ethereum/common"
	"github.com/status-im/status-go/logutils"
	"github.com/status-im/status-go/params"
	"github.com/status-im/status-go/signal"
	"github.com/status-im/statusd-bots/protocol"
	"github.com/status-im/whisper/shhclient"
	whisper "github.com/status-im/whisper/whisperv6"
)

func init() {
	if err := logutils.OverrideRootLog(true, *verbosity, "", false); err != nil {
		log.Fatalf("failed to override root log: %v\n", err)
	}
}

func main() {
	// handle OS signals
	signals := make(chan os.Signal, 1)
	stdsignal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		os.Exit(1)
	}()

	// create config
	config, err := newNodeConfig(*fleet, params.MainNetworkID)
	if err != nil {
		log.Fatalf("failed to create a config: %v", err)
	}

	// collect mail servers
	mailserversToCheck := *mailservers
	if len(mailserversToCheck) == 0 {
		// -2 to get at least two mail servers
		min := rand.Intn(len(config.ClusterConfig.TrustedMailServers) - 1)
		max := min + rand.Intn(len(config.ClusterConfig.TrustedMailServers)-min)
		mailserversToCheck = config.ClusterConfig.TrustedMailServers[min:max]
	}

	config.MaxPeers = len(mailserversToCheck)

	// collect mail server request signals
	mailSignalsForwarder := newSignalForwarder()
	defer close(mailSignalsForwarder.in)
	go mailSignalsForwarder.Start()

	// setup signals handler
	signal.SetDefaultNodeNotificationHandler(
		filterMailTypesHandler(printHandler, mailSignalsForwarder.in),
	)

	// setup work
	workConfig := WorkUnitConfig{
		From:     uint32(time.Now().Add(-*duration).Unix()),
		To:       uint32(time.Now().Add(-5 * time.Minute).Unix()), // subtract 5 mins to cater for TTL, time skew on devices etc.
		Channels: *channels,
	}

	var workUnites []*WorkUnit
	var wg sync.WaitGroup

	wg.Add(len(mailserversToCheck))

	for i, enode := range mailserversToCheck {
		config.ListenAddr = "127.0.0.1:" + strconv.Itoa(44300+i)
		config.DataDir, err = ioutil.TempDir("", "")
		if err != nil {
			log.Fatalf("failed to create temp dir: %v", err)
		}

		nodeConfig := *config
		log.Printf("using node config: %v", nodeConfig)

		work := NewWorkUnit(enode, &nodeConfig)
		go func(work *WorkUnit) {
			if err := work.Execute(workConfig, mailSignalsForwarder); err != nil {
				log.Fatalf("failed to execute work: %v", err)
			}
			wg.Done()
		}(work)
		workUnites = append(workUnites, work)
	}

	wg.Wait()

	exitCode := 0

	for i, j := 0, 1; j < len(workUnites); j++ {
		workA := workUnites[i]
		workB := workUnites[j]

		if len(workA.Messages) != len(workB.Messages) {
			exitCode = 1
		}

		log.Printf("%s vs %s: ", workA.MailServerEnode, workB.MailServerEnode)
		log.Printf("    messages: %d vs %d", len(workA.Messages), len(workB.Messages))
	}

	os.Exit(exitCode)
}

func addPublicChatSymKey(c *shhclient.Client, chat string) (string, error) {
	// This operation can be really slow, hence 10 seconds timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return c.GenerateSymmetricKeyFromPassword(ctx, chat)
}

func subscribeMessages(c *shhclient.Client, chat, symKeyID string, messages chan<- *whisper.Message) (ethereum.Subscription, error) {
	topic, err := protocol.PublicChatTopic([]byte(chat))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return c.SubscribeMessages(ctx, whisper.Criteria{
		SymKeyID: symKeyID,
		MinPow:   0,
		Topics:   []whisper.TopicType{topic},
		AllowP2P: true,
	}, messages)
}

func printHandler(event string) {
	log.Printf("received signal: %v\n", event)
}

type signalEnvelope struct {
	Type  string          `json:"type"`
	Event json.RawMessage `json:"event"`
}

type mailTypeEvent struct {
	RequestID        common.Hash `json:"requestID"`
	Hash             common.Hash `json:"hash"`
	LastEnvelopeHash common.Hash `json:"lastEnvelopeHash"`
}

type mailTypeSignal struct {
	Type           string
	RequestID      string
	LastEnvelopeID []byte
}

type signalForwarder struct {
	sync.Mutex

	in  chan mailTypeSignal
	out map[string]chan<- mailTypeSignal
}

func newSignalForwarder() *signalForwarder {
	return &signalForwarder{
		in:  make(chan mailTypeSignal),
		out: make(map[string]chan<- mailTypeSignal),
	}
}

func (s *signalForwarder) Start() {
	for {
		sig, ok := <-s.in
		if !ok {
			return
		}

		s.Lock()
		out, found := s.out[sig.RequestID]
		if found {
			out <- sig
		}
		s.Unlock()
	}
}

func (s *signalForwarder) cancel(reqID []byte) {
	s.Lock()
	delete(s.out, hex.EncodeToString(reqID))
	s.Unlock()
}

func (s *signalForwarder) Filter(reqID []byte) (<-chan mailTypeSignal, func()) {
	c := make(chan mailTypeSignal)
	s.Lock()
	s.out[hex.EncodeToString(reqID)] = c
	s.Unlock()
	return c, func() { s.cancel(reqID); close(c) }
}

func filterMailTypesHandler(fn func(string), in chan<- mailTypeSignal) func(string) {
	return func(event string) {
		fn(event)

		var envelope signalEnvelope
		if err := json.Unmarshal([]byte(event), &envelope); err != nil {
			log.Fatalf("faild to unmarshal signal Envelope: %v", err)
		}

		switch envelope.Type {
		case signal.EventMailServerRequestCompleted:
			var event mailTypeEvent
			if err := json.Unmarshal(envelope.Event, &event); err != nil {
				log.Fatalf("faild to unmarshal signal event: %v", err)
			}
			in <- mailTypeSignal{
				envelope.Type,
				hex.EncodeToString(event.RequestID.Bytes()),
				event.LastEnvelopeHash.Bytes(),
			}
		case signal.EventMailServerRequestExpired:
			var event mailTypeEvent
			if err := json.Unmarshal(envelope.Event, &event); err != nil {
				log.Fatalf("faild to unmarshal signal event: %v", err)
			}
			in <- mailTypeSignal{
				envelope.Type,
				hex.EncodeToString(event.Hash.Bytes()),
				nil,
			}
		}
	}
}

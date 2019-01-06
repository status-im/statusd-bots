package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	stdlog "log"
	"math/rand"
	"os"
	stdsignal "os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/status-im/status-go/logutils"
	"github.com/status-im/status-go/params"
	"github.com/status-im/status-go/signal"
	"github.com/status-im/statusd-bots/protocol"
	"github.com/status-im/whisper/shhclient"
	whisper "github.com/status-im/whisper/whisperv6"
)

func init() {
	if err := logutils.OverrideRootLog(true, *verbosity, "", false); err != nil {
		stdlog.Fatalf("failed to override root log: %v\n", err)
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
		log.Crit("failed to create a config", "err", err)
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

	for _, enode := range mailserversToCheck {
		config.DataDir, err = ioutil.TempDir("", "")
		if err != nil {
			log.Crit("failed to create temp dir", "err", err)
		}

		nodeConfig := *config
		log.Debug("using node config", "config", nodeConfig)

		work := NewWorkUnit(enode, &nodeConfig)
		go func(work *WorkUnit) {
			if err := work.Execute(workConfig, mailSignalsForwarder); err != nil {
				log.Crit("failed to execute work", "err", err, "enode", work.MailServerEnode)
			}
			wg.Done()
		}(work)
		workUnites = append(workUnites, work)
	}

	wg.Wait()

	// Sort results in descending order with regards to the number
	// of returned messages.
	sort.Slice(workUnites, func(i, j int) bool {
		return len(workUnites[i].Messages) > len(workUnites[j].Messages)
	})

	exitCode := 0
	failedMailServers := make([]string, 0)

	for i, j := 0, 1; j < len(workUnites); j++ {
		workA := workUnites[i]
		workB := workUnites[j]
		areEqual := len(workA.Messages) != len(workB.Messages)

		if !areEqual {
			failedMailServers = append(failedMailServers, workB.MailServerEnode)
			exitCode = 1
		}

		log.Info("MailServer A vs MailServer B",
			"A", workA.MailServerEnode,
			"messagesCountA", len(workA.Messages),
			"B", workB.MailServerEnode,
			"messagesCountB", len(workB.Messages))
	}

	if exitCode != 0 {
		log.Error("the following mail servers failed to return all messages", "enodes", failedMailServers)
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
	log.Debug("received signal", "event", event)
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
			log.Crit("faild to unmarshal signal Envelope", "err", err)
		}

		switch envelope.Type {
		case signal.EventMailServerRequestCompleted:
			var event mailTypeEvent
			if err := json.Unmarshal(envelope.Event, &event); err != nil {
				log.Crit("faild to unmarshal signal event", "event", string(envelope.Event), "err", err)
			}
			in <- mailTypeSignal{
				envelope.Type,
				hex.EncodeToString(event.RequestID.Bytes()),
				event.LastEnvelopeHash.Bytes(),
			}
		case signal.EventMailServerRequestExpired:
			var event mailTypeEvent
			if err := json.Unmarshal(envelope.Event, &event); err != nil {
				log.Crit("faild to unmarshal signal event", "err", err)
			}
			in <- mailTypeSignal{
				envelope.Type,
				hex.EncodeToString(event.Hash.Bytes()),
				nil,
			}
		}
	}
}

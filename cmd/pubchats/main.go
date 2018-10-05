package main

import (
	"context"
	"encoding/hex"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/whisper/shhclient"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/status-im/status-go/logutils"
	"github.com/status-im/status-go/node"
	"github.com/status-im/status-go/params"
	statussignal "github.com/status-im/status-go/signal"
	"github.com/status-im/statusd-bots/protocol"
)

func init() {
	if err := logutils.OverrideRootLog(true, *verbosity, "", true); err != nil {
		log.Fatalf("failed to override root log: %v\n", err)
	}

	statussignal.SetDefaultNodeNotificationHandler(func(event string) {
		log.Printf("received signal: %v\n", event)
	})
}

func main() {
	config, err := newNodeConfig(*fleet, params.MainNetworkID)
	if err != nil {
		log.Fatalf("failed to create a config: %v", err)
	}
	log.Printf("using config: %v\n", config)

	n := node.New()
	if err := n.Start(config); err != nil {
		log.Fatalf("failed to start a node: %v", err)
	}

	rpcClient, err := n.GethNode().Attach()
	if err != nil {
		log.Fatalf("failed to get an rpc: %v", err)
	}
	shh := shhclient.NewClient(rpcClient)

	// used to print a channel name from a whisper topic
	log.Printf("tracked channels: %s", *trackedChannels)
	topicsToNamesMap, err := topicsToNames(*trackedChannels)
	if err != nil {
		log.Fatalf("failed to get topics to names mapping: %v", err)
	}

	done := make(chan struct{})
	messages := make(chan *whisper.Message)
	subErr := make(chan error)

	var wg sync.WaitGroup

	for _, name := range *trackedChannels {
		go func(name string) {
			wg.Add(1)
			defer wg.Done()

			symKeyID, err := addPublicChatSymKey(shh, name)
			if err != nil {
				log.Fatalf("failed to add sym key for channel '%s': %v", name, err)
			}

			sub, err := subscribeMessages(shh, name, symKeyID, messages)
			if err != nil {
				log.Fatalf("failed to subscribe to messages for channel '%s': %v", name, err)
			}
			defer sub.Unsubscribe()

			select {
			case err := <-sub.Err():
				if err != nil {
					subErr <- err
				}
			case <-done:
			}
		}(name)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go startMetricsServer(*metricsAddr)

	log.Println("waiting for messages...")

	for {
		select {
		case msg := <-messages:
			chatName := topicsToNamesMap[msg.Topic]
			source := hex.EncodeToString(msg.Sig)
			log.Printf("received a message: topic=%v (%s) data=%s author=%s", msg.Topic, chatName, msg.Payload, source)
			messagesCounter.WithLabelValues(chatName, source).Inc()
		case err := <-subErr:
			log.Fatalf("subscription error: %v", err)
		case <-signals:
			close(done)
			wg.Wait()
			os.Exit(1)
		}
	}
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
		MinPow:   0.001,
		Topics:   []whisper.TopicType{topic},
	}, messages)
}

func topicsToNames(names []string) (map[whisper.TopicType]string, error) {
	m := make(map[whisper.TopicType]string)
	for _, name := range names {
		t, err := protocol.PublicChatTopic([]byte(name))
		if err != nil {
			return nil, err
		}
		m[t] = name
	}
	return m, nil
}

package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/status-im/status-go/logutils"
	"github.com/status-im/status-go/node"
	"github.com/status-im/status-go/params"
	"github.com/status-im/status-go/services/shhext"
	statussignal "github.com/status-im/status-go/signal"
	"github.com/status-im/status-go/t/helpers"
	"github.com/status-im/statusd-bots/protocol"
	"github.com/status-im/whisper/shhclient"
	whisper "github.com/status-im/whisper/whisperv6"
)

func init() {
	if err := logutils.OverrideRootLog(true, *verbosity, "", false); err != nil {
		log.Fatalf("failed to override root log: %v\n", err)
	}

	statussignal.SetDefaultNodeNotificationHandler(func(event string) {
		log.Printf("received signal: %v\n", event)
	})
}

func main() {
	config, err := newNodeConfig(*address, *fleet, params.MainNetworkID)
	if err != nil {
		log.Fatalf("failed to create a config: %v", err)
	}
	log.Printf("using config: %v", config)

	n := node.New()
	if err := n.Start(config); err != nil {
		log.Fatalf("failed to start a node: %v", err)
	}

	rpcClient, err := n.GethNode().Attach()
	if err != nil {
		log.Fatalf("failed to get an rpc: %v", err)
	}
	shh := shhclient.NewClient(rpcClient)

	shhextService, err := n.ShhExtService()
	if err != nil {
		log.Fatalf("failed go get an shhext service: %v", err)
	}
	shhextAPI := shhext.NewPublicAPI(shhextService)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Println("subscribe for messages...")

	symKeyID, err := addPublicChatSymKey(shh, *channel)
	if err != nil {
		log.Fatalf("failed to add sym key for channel '%s': %v", *channel, err)
	}

	// TODO(adam): this does not work and messages are not displayed.
	messages := make(chan *whisper.Message)
	sub, err := subscribeMessages(shh, *channel, symKeyID, messages)
	if err != nil {
		log.Fatalf("failed to subscribe to messages for channel '%s': %v", *channel, err)
	}
	defer sub.Unsubscribe()

	log.Println("adding Mail Server as a peer")

	mailserverEnode := *mailserver
	if mailserverEnode == "" {
		mailserverEnode = config.ClusterConfig.TrustedMailServers[rand.Intn(len(config.ClusterConfig.TrustedMailServers))]
	}

	if err := n.AddPeer(mailserverEnode); err != nil {
		log.Fatalf("failed to add Mail Server as a peer: %v", err)
	}

	errCh := helpers.WaitForPeerAsync(n.Server(), mailserverEnode, p2p.PeerEventTypeAdd, 5*time.Second)
	if err := <-errCh; err != nil {
		log.Fatalf("failed to wait for peer '%s': %v", mailserverEnode, err)
	}

	log.Println("sending requests to Mail Server")

	topic, err := protocol.PublicChatTopic([]byte(*channel))
	if err != nil {
		log.Fatalf("failed to get topic for channel %s: %v", *channel, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	mailServerSymKeyID, err := shh.GenerateSymmetricKeyFromPassword(ctx, protocol.MailServerPassword)
	if err != nil {
		log.Fatalf("failed to generate sym key for mail server: %v", err)
	}

	for i := 0; i < *concurrency; i++ {
		go func() {
			hash, err := shhextAPI.RequestMessages(nil, shhext.MessagesRequest{
				MailServerPeer: mailserverEnode,
				SymKeyID:       mailServerSymKeyID,
				From:           uint32(time.Now().Add(-*duration).Unix()),
				To:             uint32(time.Now().Unix()),
				Limit:          1000,
				Topic:          topic,
				Timeout:        30,
			})
			if err != nil {
				log.Fatalf("failed to request for messages: %v", err)
			}
			log.Printf("requested for messages with a request hash: %s", hash)
		}()
	}

	for {
		select {
		case msg := <-messages:
			source := hex.EncodeToString(msg.Sig)
			log.Printf("received a message: topic=%v data=%s author=%s", msg.Topic, msg.Payload, source)
		case err := <-sub.Err():
			log.Fatalf("subscription error: %v", err)
		case <-signals:
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

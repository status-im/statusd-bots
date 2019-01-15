package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/status-im/status-go/node"
	"github.com/status-im/status-go/params"
	"github.com/status-im/status-go/services/shhext"
	"github.com/status-im/status-go/signal"
	"github.com/status-im/status-go/t/helpers"
	"github.com/status-im/statusd-bots/protocol"
	"github.com/status-im/whisper/shhclient"
	whisper "github.com/status-im/whisper/whisperv6"
)

// WorkUnit represents a single unit of work.
type WorkUnit struct {
	MailServerEnode string
	Messages        []*whisper.Message

	config    *params.NodeConfig
	node      *node.StatusNode
	shh       *shhclient.Client
	shhextAPI *shhext.PublicAPI
}

// NewWorkUnit creates a new WorkUnit instance.
func NewWorkUnit(mailEnode string, config *params.NodeConfig) *WorkUnit {
	return &WorkUnit{
		MailServerEnode: mailEnode,
		config:          config,
	}
}

// WorkUnitConfig configures the execution of the work.
type WorkUnitConfig struct {
	From     uint32
	To       uint32
	Channels []string
}

// Execute runs the work.
func (u *WorkUnit) Execute(config WorkUnitConfig, mailSignals *signalForwarder) error {
	if err := u.startNode(); err != nil {
		return fmt.Errorf("failed to start node: %v", err)
	}

	if err := u.addPeer(); err != nil {
		return fmt.Errorf("failed to add peer: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	mailServerSymKeyID, err := u.shh.GenerateSymmetricKeyFromPassword(
		ctx, protocol.MailServerPassword)
	if err != nil {
		return fmt.Errorf("failed to generate sym key for mail server: %v", err)
	}

	var topics []whisper.TopicType
	for _, ch := range config.Channels {
		topic, err := protocol.PublicChatTopic([]byte(ch))
		if err != nil {
			return fmt.Errorf("failed to create topic: %v", err)
		}
		topics = append(topics, topic)
	}

	var messageSubErrs []<-chan error
	messages := make(chan *whisper.Message)

	for _, ch := range config.Channels {
		symKeyID, err := addPublicChatSymKey(u.shh, ch)
		if err != nil {
			return fmt.Errorf("failed to add sym key for channel '%s': %v", ch, err)
		}

		sub, err := subscribeMessages(u.shh, ch, symKeyID, messages)
		if err != nil {
			return fmt.Errorf("failed to subscribe for messages: %v", err)
		}
		defer sub.Unsubscribe()
		messageSubErrs = append(messageSubErrs, sub.Err())
	}

	// TODO: sshext.MessagesRequest expects time.Duration but multiplies it by time.Second
	reqTimeout := time.Duration(15)
	reqID, err := u.shhextAPI.RequestMessages(nil, shhext.MessagesRequest{
		MailServerPeer: u.MailServerEnode,
		SymKeyID:       mailServerSymKeyID,
		From:           config.From,
		To:             config.To,
		Limit:          1000,
		Topics:         topics,
		Timeout:        reqTimeout,
	})
	if err != nil {
		return fmt.Errorf("failed to request %s for messages: %v", u.MailServerEnode, err)
	}

	// TODO(adam): change to regular fanout. It might happen that a signal
	// is received before the filter is setup.
	signals, cancelSignalsFilter := mailSignals.Filter([]byte(reqID))
	defer cancelSignalsFilter()

	start := time.Now()

	var lastEnvelopeID []byte

	for {
		select {
		case m := <-messages:
			log.Debug("received a message", "hash", hex.EncodeToString(m.Hash))
			u.Messages = append(u.Messages, m)
		case <-time.After(time.Duration(reqTimeout) * time.Second):
			// As we can not predict when messages finish to come in,
			// we timeout after some time.
			// If lastEnvelopeID is found amoung received messages,
			// it's a successful request. Otherwise, an error is returned.
			for i, m := range u.Messages {
				if bytes.Equal(lastEnvelopeID, m.Hash) {
					log.Info("received a message equal to lastEnvelopeID",
						"hash", hex.EncodeToString(lastEnvelopeID),
						"index", i,
						"messagesCount", len(u.Messages))
					return u.stopNode()
				}
			}
			return fmt.Errorf("did not receive a message equal to lastEnvelopeID")
		case err := <-merge(messageSubErrs...):
			return fmt.Errorf("subscription for messages errored: %v", err)
		case s := <-signals:
			switch s.Type {
			case signal.EventMailServerRequestCompleted:
				lastEnvelopeID = s.LastEnvelopeID

				log.Info("received EventMailServerRequestCompleted", "latency", time.Since(start), "enode", u.MailServerEnode, "lastEnvelopeID", lastEnvelopeID)

				if allZeros(lastEnvelopeID) {
					log.Info("lastEnvelopeID is empty so return early")
					return u.stopNode()
				}
			case signal.EventMailServerRequestExpired:
				return fmt.Errorf("request for messages expired")
			}
		}
	}
}

func (u *WorkUnit) startNode() error {
	u.node = node.New()
	if err := u.node.Start(u.config); err != nil {
		return fmt.Errorf("failed to start a node: %v", err)
	}

	rpcClient, err := u.node.GethNode().Attach()
	if err != nil {
		return fmt.Errorf("failed to get an rpc: %v", err)
	}
	u.shh = shhclient.NewClient(rpcClient)

	shhextService, err := u.node.ShhExtService()
	if err != nil {
		return fmt.Errorf("failed go get an shhext service: %v", err)
	}
	u.shhextAPI = shhext.NewPublicAPI(shhextService)

	return nil
}

func (u *WorkUnit) stopNode() error {
	return u.node.Stop()
}

func (u *WorkUnit) addPeer() error {
	if err := u.node.AddPeer(u.MailServerEnode); err != nil {
		return err
	}

	return <-helpers.WaitForPeerAsync(
		u.node.Server(),
		u.MailServerEnode,
		p2p.PeerEventTypeAdd,
		5*time.Second,
	)
}

func allZeros(b []byte) bool {
	zero := byte(0)
	for _, n := range b {
		if n != zero {
			return false
		}
	}
	return true
}

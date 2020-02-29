package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/status-im/status-go/protocol/sqlite"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"

	"github.com/ethereum/go-ethereum/p2p"
	gethbridge "github.com/status-im/status-go/eth-node/bridge/geth"
	"github.com/status-im/status-go/eth-node/types"
	"github.com/status-im/status-go/node"
	"github.com/status-im/status-go/params"
	"github.com/status-im/status-go/protocol"
	"github.com/status-im/status-go/t/helpers"
)

// WorkUnit represents a single unit of work.
// It will make a request for historic messages
// to a mailserver and collect received envelopes.
type WorkUnit struct {
	MailServerEnode string
	MessageHashes   []types.HexBytes // a list of collected messages.

	config    *params.NodeConfig
	node      *node.StatusNode
	key       *ecdsa.PrivateKey
	messenger *protocol.Messenger
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
	From  uint32
	To    uint32
	Chats []string
}

// Execute runs the work.
func (u *WorkUnit) Execute(config WorkUnitConfig) error {
	if err := u.startNode(); err != nil {
		return fmt.Errorf("failed to start node: %v", err)
	}

	if err := u.startMessenger(); err != nil {
		return fmt.Errorf("failed to create messenger: %v", err)
	}

	if err := u.addPeer(u.MailServerEnode); err != nil {
		return fmt.Errorf("failed to add peer: %v", err)
	}

	for _, chatName := range config.Chats {
		chat := protocol.CreatePublicChat(chatName, u.messenger.Timesource())
		if err := u.messenger.Join(chat); err != nil {
			return err
		}
		if err := u.messenger.SaveChat(&chat); err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	msEnode := enode.MustParse(u.MailServerEnode)
	_, err := u.messenger.RequestHistoricMessages(ctx, msEnode.ID().Bytes(), config.From, config.To, nil)
	if err != nil {
		return err
	}

	<-time.After(time.Second) // wait for the Whisper loop to turn (whisper loop turns every 300ms)
	resp, err := u.messenger.RetrieveAll()
	if err != nil {
		return err
	}

	for _, m := range resp.Messages {
		u.MessageHashes = append(u.MessageHashes, types.FromHex(m.ID))
	}

	for _, rawResp := range resp.RawMessages {
		for _, m := range rawResp.Messages {
			u.MessageHashes = append(u.MessageHashes, m.ID)
		}
	}

	return nil
}

func (u *WorkUnit) startNode() error {
	u.node = node.New()
	if err := u.node.Start(u.config, &accounts.Manager{}); err != nil {
		return fmt.Errorf("failed to start a node: %v", err)
	}
	return nil
}

func (u *WorkUnit) stopNode() error {
	return u.node.Stop()
}

func (u *WorkUnit) startMessenger() error {
	key, err := crypto.GenerateKey()
	if err != nil {
		return err
	}
	u.key = key
	node := gethbridge.NewNodeBridge(u.node.GethNode())
	db, err := sqlite.OpenInMemory()
	if err != nil {
		return err
	}
	messenger, err := protocol.NewMessenger(key, node, "instalation-01", protocol.WithDatabase(db))
	if err != nil {
		return err
	}
	if err := messenger.Start(); err != nil {
		return err
	}
	u.messenger = messenger
	return nil
}

func (u *WorkUnit) addPeer(enodeAddr string) error {
	if err := u.node.AddPeer(enodeAddr); err != nil {
		return err
	}

	return <-helpers.WaitForPeerAsync(
		u.node.Server(),
		enodeAddr,
		p2p.PeerEventTypeAdd,
		5*time.Second,
	)
}

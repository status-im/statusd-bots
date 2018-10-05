package main

import (
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/status-im/status-go/params"
)

func newNodeConfig(fleet string, networkID uint64) (*params.NodeConfig, error) {
	c, err := params.NewNodeConfigWithDefaults(
		*datadir, networkID, params.WithFleet(fleet))
	if err != nil {
		return nil, err
	}

	c.ListenAddr = *address
	c.MaxPeers = 10
	c.IPCEnabled = true
	c.HTTPEnabled = false

	c.RequireTopics = map[discv5.Topic]params.Limits{
		discv5.Topic("whisper"): params.Limits{Min: 3, Max: 3},
	}

	return c, nil
}

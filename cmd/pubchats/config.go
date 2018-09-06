package main

import (
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/status-im/status-go/params"
)

func newNodeConfig(fleet string, networkID uint64) (*params.NodeConfig, error) {
	c, err := params.NewNodeConfig(*datadir, "", fleet, networkID)
	if err != nil {
		return nil, err
	}

	c.ListenAddr = *address
	c.MaxPeers = 10
	c.IPCEnabled = false
	c.RPCEnabled = false
	c.HTTPHost = ""

	c.LightEthConfig.Enabled = false
	c.LightEthConfig.Genesis = ""

	c.Rendezvous = true
	c.RequireTopics = map[discv5.Topic]params.Limits{
		discv5.Topic("whisper"): params.Limits{Min: 3, Max: 3},
	}

	c.WhisperConfig.Enabled = true
	c.WhisperConfig.LightClient = false
	c.WhisperConfig.MinimumPoW = 0.001
	c.WhisperConfig.TTL = 120
	c.WhisperConfig.EnableNTPSync = true

	return c, nil
}

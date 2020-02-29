package main

import (
	"github.com/status-im/status-go/params"
)

func newNodeConfig(fleet string, networkID uint64) (*params.NodeConfig, error) {
	c, err := params.NewNodeConfig("", networkID)
	if err != nil {
		return nil, err
	}

	if err := params.WithFleet(fleet)(c); err != nil {
		return nil, err
	}

	c.ListenAddr = ":0"
	c.MaxPeers = 10
	c.IPCEnabled = true
	c.HTTPEnabled = false
	c.NoDiscovery = true
	c.Rendezvous = false

	c.ClusterConfig.Enabled = false
	c.ClusterConfig.StaticNodes = nil

	c.WhisperConfig.Enabled = true

	return c, nil
}

package main

import (
	"github.com/status-im/status-go/params"
)

func newNodeConfig(fleet string, networkID uint64) (*params.NodeConfig, error) {
	c, err := params.NewNodeConfig(*datadir, networkID)
	if err != nil {
		return nil, err
	}

	if err := params.WithFleet(params.FleetBeta)(c); err != nil {
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
	c.WhisperConfig.EnableNTPSync = true

	return c, nil
}

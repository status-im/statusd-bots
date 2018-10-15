package main

import (
	"github.com/status-im/status-go/params"
)

func newNodeConfig(address, fleet string, networkID uint64) (*params.NodeConfig, error) {
	c, err := params.NewNodeConfigWithDefaults(
		*datadir, networkID, params.WithFleet(fleet))
	if err != nil {
		return nil, err
	}

	c.ListenAddr = address
	c.MaxPeers = 10
	c.IPCEnabled = true
	c.HTTPEnabled = false

	c.ClusterConfig.StaticNodes = nil

	return c, nil
}

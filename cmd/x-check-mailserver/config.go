package main

import (
	"errors"
	"strings"

	"github.com/status-im/status-go/params"
)

func newNodeConfig(fleet string, networkID uint64, privKey string) (*params.NodeConfig, error) {
	c, err := params.NewNodeConfig("", networkID)
	if err != nil {
		return nil, err
	}

	if err := params.WithFleet(fleet)(c); err != nil {
		return nil, err
	}

	cleanPrivKey := strings.TrimPrefix(privKey, "0x")
	if len(cleanPrivKey) == 64 {
		c.NodeKey = cleanPrivKey
	} else if len(cleanPrivKey) != 64 && len(cleanPrivKey) != 0 {
		return nil, errors.New("wrong private key length, expected 64 char hexadecimal")
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

package main

import (
	"github.com/spf13/pflag"
	"github.com/status-im/status-go/params"
)

var (
	datadir         = pflag.StringP("datadir", "d", "", "directory for data")
	address         = pflag.StringP("addr", "a", "127.0.0.1:30303", "listener IP address")
	fleet           = pflag.StringP("fleet", "f", params.FleetBeta, "cluster fleet")
	trackedChannels = pflag.StringSliceP("channel", "c", []string{}, "public channels to track")
	verbosity       = pflag.StringP("verbosity", "v", "INFO", "verbosity level of status-go, options: crit, error, warning, info, debug")
	metricsAddr     = pflag.StringP("metrics-addr", "m", ":8080", "metrics server listening address")
)

func init() {
	pflag.Parse()
}

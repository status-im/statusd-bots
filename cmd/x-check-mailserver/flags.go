package main

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/status-im/status-go/params"
)

var (
	datadir     = pflag.StringP("datadir", "d", "", "directory for data")
	fleet       = pflag.StringP("fleet", "f", params.FleetBeta, "cluster fleet")
	mailservers = pflag.StringArrayP("mailservers", "m", nil, "a list of mail servers")
	duration    = pflag.DurationP("duration", "l", time.Hour*24, "length of time span from now")
	channel     = pflag.StringP("channel", "p", "status", "name of the channel")
	verbosity   = pflag.StringP("verbosity", "v", "INFO", "verbosity level of status-go, options: crit, error, warning, info, debug")
)

func init() {
	pflag.Parse()
}

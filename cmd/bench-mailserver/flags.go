package main

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/status-im/status-go/params"
)

var (
	datadir     = pflag.StringP("datadir", "d", "", "directory for data")
	address     = pflag.StringP("addr", "a", "127.0.0.1:30303", "listener IP address")
	fleet       = pflag.StringP("fleet", "f", params.FleetBeta, "cluster fleet")
	mailserver  = pflag.StringP("mailserver", "m", "", "MailServer address (by default a random one from the fleet is selected)")
	concurrency = pflag.IntP("concurrency", "c", 5, "number of concurrent requests")
	duration    = pflag.DurationP("duration", "l", time.Hour*24, "length of time span from now")
	channel     = pflag.StringP("channel", "p", "status", "name of the channel")
	verbosity   = pflag.StringP("verbosity", "v", "INFO", "verbosity level of status-go, options: crit, error, warning, info, debug")
)

func init() {
	pflag.Parse()
}

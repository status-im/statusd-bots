package main

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/status-im/status-go/params"
)

var (
	fleet       = pflag.StringP("fleet", "f", params.FleetProd, "cluster fleet")
	datadir     = pflag.StringP("datadir", "d", "", "home directory for node data")
	privkey     = pflag.StringP("privkey", "p", "", "private key for connecting to nodes, hexadecimal")
	mailservers = pflag.StringArrayP("mailservers", "m", nil, "a list of mail servers")
	duration    = pflag.DurationP("duration", "l", time.Hour*24, "length of time span from now")
	channels    = pflag.StringArrayP("channels", "c", []string{"status"}, "name of one or more channels")
	verbosity   = pflag.StringP("verbosity", "v", "INFO", "verbosity level of status-go, options: crit, error, warn, info, debug")
)

func init() {
	pflag.Parse()
}

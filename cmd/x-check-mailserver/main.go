package main

import (
	"io/ioutil"
	stdlog "log"
	"math/rand"
	"os"
	stdsignal "os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/status-im/status-go/logutils"
	"github.com/status-im/status-go/params"
)

func init() {
	if err := logutils.OverrideRootLog(true, *verbosity, logutils.FileOptions{}, false); err != nil {
		stdlog.Fatalf("failed to override root log: %v\n", err)
	}
}

func main() {
	// handle OS signals
	signals := make(chan os.Signal, 1)
	stdsignal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		os.Exit(1)
	}()

	// create config
	config, err := newNodeConfig(*fleet, params.MainNetworkID, *privkey)
	if err != nil {
		log.Crit("failed to create a config", "err", err)
	}

	// collect mail servers
	mailserversToCheck := *mailservers
	if len(mailserversToCheck) == 0 {
		// -1 to get at least two mail servers
		min := rand.Intn(len(config.ClusterConfig.TrustedMailServers) - 1)
		max := min + rand.Intn(len(config.ClusterConfig.TrustedMailServers)-min)
		mailserversToCheck = config.ClusterConfig.TrustedMailServers[min:max]
	}

	config.MaxPeers = len(mailserversToCheck)

	// setup work
	workConfig := WorkUnitConfig{
		Channels: *channels,
		// starting time for the envelope query
		From: uint32(time.Now().Add(-*duration).Unix()),
		// subtract 5 mins to cater for TTL, time skew on devices etc.
		To: uint32(time.Now().Add(-5 * time.Minute).Unix()),
	}

	var (
		workUnites []*WorkUnit
		wg         sync.WaitGroup
	)

	wg.Add(len(mailserversToCheck))

	for _, msEnode := range mailserversToCheck {
		var nodeId = enode.MustParse(msEnode).ID().String()
		config.DataDir, err = ioutil.TempDir(*datadir, nodeId)
		if err != nil {
			log.Crit("failed to create temp dir", "err", err)
		}

		nodeConfig := *config
		log.Debug("using node config", "config", nodeConfig)

		work := NewWorkUnit(msEnode, &nodeConfig, *privkey)
		go func(work *WorkUnit) {
			if err := work.Execute(workConfig); err != nil {
				log.Crit("failed to execute work", "err", err, "enode", work.MailServerEnode)
			}
			wg.Done()
		}(work)
		workUnites = append(workUnites, work)
	}

	wg.Wait()

	// Sort results in descending order with regards to the number of received hashes.
	sort.Slice(workUnites, func(i, j int) bool {
		return len(workUnites[i].MessageHashes) > len(workUnites[j].MessageHashes)
	})

	exitCode := 0
	failedMailServers := make([]string, 0)

	for i, j := 0, 1; j < len(workUnites); j++ {
		workA := workUnites[i]
		workB := workUnites[j]
		areEqual := len(workA.MessageHashes) == len(workB.MessageHashes)

		if !areEqual {
			failedMailServers = append(failedMailServers, workB.MailServerEnode)
			exitCode = 1
		}

		log.Info("MailServer A vs MailServer B",
			"A", workA.MailServerEnode,
			"messagesCountA", len(workA.MessageHashes),
			"B", workB.MailServerEnode,
			"messagesCountB", len(workB.MessageHashes))
	}

	if exitCode != 0 {
		log.Error("the following mail servers failed to return all messages", "enodes", failedMailServers)
	}

	os.Exit(exitCode)
}

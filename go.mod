module github.com/status-im/statusd-bots

go 1.13

replace github.com/ethereum/go-ethereum v1.9.5 => github.com/status-im/go-ethereum v1.9.5-status.7

replace github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/ethereum/go-ethereum v1.9.5
	github.com/ipfs/go-log v1.0.0 // indirect
	github.com/multiformats/go-multihash v0.0.10 // indirect
	github.com/prometheus/client_golang v1.2.1
	github.com/spf13/pflag v1.0.3
	github.com/status-im/status-go v0.46.0
	github.com/status-im/status-go/whisper/v6 v6.2.4
	github.com/status-im/whisper v1.6.2
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
)

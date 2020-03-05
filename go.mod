module github.com/status-im/statusd-bots

go 1.13

replace github.com/ethereum/go-ethereum v1.9.5 => github.com/status-im/go-ethereum v1.9.5-status.7

replace github.com/Sirupsen/logrus v1.4.2 => github.com/sirupsen/logrus v1.4.2

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20190717161051-705d9623b7c1

require (
	github.com/ethereum/go-ethereum v1.9.5
	github.com/spf13/pflag v1.0.3
	github.com/status-im/status-go v0.48.2
)

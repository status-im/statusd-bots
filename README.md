statusd-bots
============

*This is an alternative project to [status-im/status-go-bots](https://github.com/status-im/status-go-bots). The goal is to explore another possibility of creating bots that are self-contained Whisper nodes and do not require a Whisper node running aside.*

## Setup

Project dependencies are managed with [`dep`](https://github.com/golang/dep). Please install it first.

In order to install all dependencies, execute:

```
$ make dependencies
```

This project uses [`github.com/ethereum/go-ethereum`](https://github.com/ethereum/go-ethereum) as a dependency but its source is changed to Status' fork [status-im/go-ethereum](https://github.com/status-im/go-ethereum) which include patches required by another dependency [`github.com/status-im/status-go`](https://github.com/status-im/status-go).

## Bots

### pubchats

It follows Status (Whisper) public chats and provide logs and Prometheus metrics. Public chats can be read and write by any node and the mechanism to encrypt and find such messages is known.

```
$ ./bin/pubchats -h
Usage of ./bin/pubchats:
  -a, --addr string           listener IP address (default "127.0.0.1:30303")
  -c, --channel strings       public channels to track
  -d, --datadir string        directory for data
  -f, --fleet string          cluster fleet (default "eth.beta")
  -m, --metrics-addr string   metrics server listening address (default ":8080")
  -v, --verbosity string      verbosity level of status-go, options: crit, error, warning, info, debug (default "INFO")
```

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

It can be used to follow Status (Whisper) public chats and provide metrics. Public chats can be read and write by any node and the mechanism to encrypt and find such messages is known.


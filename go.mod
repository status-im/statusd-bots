module github.com/status-im/statusd-bots

go 1.13

require (
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/aristanetworks/goarista v0.0.0-20181109020153-5faa74ffbed7
	github.com/beevik/ntp v0.2.0
	github.com/beorn7/perks v1.0.0
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/coreos/go-semver v0.3.0
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea
	github.com/edsrzf/mmap-go v0.0.0-20170320065105-0bce6a688712
	github.com/elastic/gosigar v0.10.5 // indirect
	github.com/ethereum/go-ethereum v1.9.5
	github.com/fd/go-nat v1.0.0
	github.com/fjl/memsize v0.0.0-20180929194037-2a09253e352a
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08 // indirect
	github.com/go-playground/locales v0.12.1
	github.com/go-playground/universal-translator v0.16.0
	github.com/go-stack/stack v1.8.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang-migrate/migrate v3.5.4+incompatible // indirect
	github.com/golang-migrate/migrate/v4 v4.7.0 // indirect
	github.com/golang/mock v1.2.0
	github.com/golang/protobuf v1.3.1
	github.com/golang/snappy v0.0.1
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/gxed/GoEndian v0.0.0-20160916112711-0f5c6873267e
	github.com/gxed/eventfd v0.0.0-20160916113412-80a92cca79a8
	github.com/gxed/hashland v0.0.0-20180221191214-d9f6b97f8db2
	github.com/hashicorp/golang-lru v0.5.1
	github.com/huin/goupnp v1.0.0
	github.com/ipfs/go-datastore v0.3.1 // indirect
	github.com/ipfs/go-log v1.0.0 // indirect
	github.com/jackpal/gateway v1.0.5
	github.com/jackpal/go-nat-pmp v1.0.1
	github.com/jbenet/go-temp-err-catcher v0.0.0-20150120210811-aac704a3f4f2
	github.com/jbenet/goprocess v0.1.3
	github.com/karalabe/hid v0.0.0-20181128192157-d815e0c1a2e2
	github.com/karalabe/usb v0.0.0-20191104083709-911d15fe12a9 // indirect
	github.com/libp2p/go-conn-security v0.1.0 // indirect
	github.com/libp2p/go-libp2p v0.0.0-20180609053045-a08d9e63dbf0
	github.com/libp2p/go-libp2p-circuit v0.0.0-20180924121056-eca2b86a1bcf
	github.com/libp2p/go-libp2p-host v0.1.0 // indirect
	github.com/libp2p/go-libp2p-interface-connmgr v0.1.0 // indirect
	github.com/libp2p/go-libp2p-interface-pnet v0.1.0 // indirect
	github.com/libp2p/go-libp2p-metrics v0.1.0 // indirect
	github.com/libp2p/go-libp2p-nat v0.0.5 // indirect
	github.com/libp2p/go-libp2p-net v0.1.0 // indirect
	github.com/libp2p/go-libp2p-peer v0.2.0 // indirect
	github.com/libp2p/go-libp2p-peerstore v0.1.3
	github.com/libp2p/go-libp2p-protocol v0.1.0 // indirect
	github.com/libp2p/go-libp2p-secio v0.2.1 // indirect
	github.com/libp2p/go-libp2p-swarm v0.2.2 // indirect
	github.com/libp2p/go-libp2p-transport v0.1.0 // indirect
	github.com/libp2p/go-stream-muxer v0.1.0 // indirect
	github.com/libp2p/go-ws-transport v0.1.2 // indirect
	github.com/mattn/go-colorable v0.1.1
	github.com/mattn/go-isatty v0.0.5
	github.com/matttproud/golang_protobuf_extensions v1.0.1
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/minio/sha256-simd v0.1.1
	github.com/mr-tron/base58 v1.1.2
	github.com/multiformats/go-multiaddr v0.1.1
	github.com/multiformats/go-multiaddr-dns v0.0.1
	github.com/multiformats/go-multiaddr-net v0.1.1 // indirect
	github.com/multiformats/go-multihash v0.0.10 // indirect
	github.com/mutecomm/go-sqlcipher v0.0.0-20170920224653-f799951b4ab2
	github.com/opentracing/opentracing-go v1.1.0
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709
	github.com/prometheus/client_golang v1.0.0
	github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
	github.com/prometheus/common v0.4.1
	github.com/prometheus/procfs v0.0.2
	github.com/prometheus/prometheus v2.1.0+incompatible
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rjeczalik/notify v0.9.2
	github.com/rs/cors v1.6.0
	github.com/spaolacci/murmur3 v1.1.0
	github.com/spf13/pflag v1.0.3
	github.com/status-im/doubleratchet v1.0.0
	github.com/status-im/go-multiaddr-ethv4 v1.0.0
	github.com/status-im/keycard-go v0.0.0-20191119114148-6dd40a46baa0 // indirect
	github.com/status-im/migrate v3.4.0+incompatible // indirect
	github.com/status-im/migrate/v4 v4.6.2-status.2 // indirect
	github.com/status-im/rendezvous v1.0.1-0.20181122054443-e0f2e0d17d81
	github.com/status-im/status-go v0.16.5-0.20181114070358-52a1bdfed669
	github.com/status-im/whisper v1.3.0
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570 // indirect
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	github.com/syndtr/goleveldb v1.0.0
	github.com/tyler-smith/go-bip39 v1.0.2 // indirect
	github.com/whyrusleeping/go-logging v0.0.0-20170515211332-0457bb6b88fc
	github.com/whyrusleeping/go-notifier v0.0.0-20170827234753-097c5d47330f
	github.com/whyrusleeping/go-smux-multiplex v3.0.16+incompatible
	github.com/whyrusleeping/go-smux-multistream v2.0.2+incompatible
	github.com/whyrusleeping/go-smux-yamux v2.0.8+incompatible
	github.com/whyrusleeping/mafmt v1.2.8
	github.com/whyrusleeping/multiaddr-filter v0.0.0-20160516205228-e903e4adabd7
	github.com/whyrusleeping/yamux v1.1.2
	github.com/wsddn/go-ecdh v0.0.0-20161211032359-48726bab9208 // indirect
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190626221950-04f50cda93cb
	golang.org/x/text v0.3.2
	gopkg.in/go-playground/validator.v9 v9.23.0
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20181125150206-ccb656ba24c2
	gopkg.in/urfave/cli.v1 v1.20.0
)

replace github.com/ethereum/go-ethereum v1.9.5 => github.com/status-im/go-ethereum v1.9.5-status.5

package main

import (
	"TestLibp2p/protocol"
	"fmt"
	"github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
)

func main() {
	log.SetAllLoggers(log.LevelWarn)
	log.SetLogLevel("rendezvous", "info")

	cfg := protocol.ParseFlags()

	host, err := libp2p.New(libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%d", cfg.ListenHost, cfg.ListenPort)))
	if err != nil {
		panic(err)
	}
	protocol.StartPeer(host, cfg)
	select {}
}

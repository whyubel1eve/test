package main

import (
	"TestLibp2p/protocol"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"log"
)

func main() {

	sourcePort := flag.Int("p", 0, "source port")
	dest := flag.String("d", "", "destination multiAddress")
	flag.Parse()

	sourceAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort)
	host, err := libp2p.New(
		libp2p.ListenAddrStrings(sourceAddr),
	)
	if err != nil {
		log.Printf("failed to create p2p host: %v", err)
	}

	if *dest == "" {
		protocol.StartPeer(host)
	} else {
		err := protocol.StartPeerAndConnect(host, *dest)
		if err != nil {
			log.Fatalf("Failed to start or connect peer: %v", err)
		}
	}
	select {}
}

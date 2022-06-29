package protocol

import (
	"bufio"
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"log"
	"os"
)

const (
	ProtocolID = "/chat/1.0.0"
)

func StartPeer(h host.Host) {
	h.SetStreamHandler(ProtocolID, ChatHandler)
	addr := h.Addrs()[0]
	id := h.ID()
	peerInfoStr := fmt.Sprintf("%s/p2p/%s", addr, id)
	log.Printf("Peer infomation: %s", peerInfoStr)
}
func StartPeerAndConnect(h host.Host, dest string) error {
	peerInfoStr := fmt.Sprintf("%s/p2p/%s", h.Addrs()[0], h.ID())
	log.Printf("Peer infomation: %s", peerInfoStr)
	destInfo, err := peer.AddrInfoFromString(dest)
	if err != nil {
		return err
	}
	h.Peerstore().AddAddr(destInfo.ID, destInfo.Addrs[0], peerstore.PermanentAddrTTL)

	s, err := h.NewStream(context.Background(), destInfo.ID, ProtocolID)
	if err != nil {
		return err
	}
	log.Println("Established connection to destination")

	// Create a buffered stream so that read and writes are non blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go ReadData(rw)
	go WriteData(rw)
	return nil
}

func ChatHandler(s network.Stream) {
	log.Println("Got a new stream!")
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go ReadData(rw)
	go WriteData(rw)
}
func WriteData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		rw.WriteString(fmt.Sprintf("%s", sendData))
		rw.Flush()
	}
}
func ReadData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')
		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

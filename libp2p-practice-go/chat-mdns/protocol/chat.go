package protocol

import (
	"bufio"
	"context"
	"fmt"
	"github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/protocol"
	"os"
)

var logger = log.Logger("rendezvous")

func StartPeer(h host.Host, cfg *Config) {
	ctx := context.Background()
	h.SetStreamHandler(protocol.ID(cfg.ProtocolID), ChatHandler)
	addr := h.Addrs()[0]
	id := h.ID()
	peerInfoStr := fmt.Sprintf("%s/p2p/%s", addr, id)
	logger.Info("Peer infomation: ", peerInfoStr)

	peerChan := initMDNS(h, cfg.RendezvousString)

	peer := <-peerChan // will block untill we discover a peer
	if peer.ID != h.ID() {
		logger.Info("Found peer:", peer, ", connecting...")

		if err := h.Connect(ctx, peer); err != nil {
			logger.Warn("Connection failed:", err)
		}

		// open a stream, this stream will be handled by handleStream other end
		stream, err := h.NewStream(ctx, peer.ID, protocol.ID(cfg.ProtocolID))

		if err != nil {
			logger.Warn("Stream open failed", err)
		} else {
			rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			logger.Info("Connected to:", peer)
			go WriteData(rw)
			go ReadData(rw)
		}
	}

}

func ChatHandler(s network.Stream) {
	logger.Info("Got a new stream!")
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
			panic(err)
		}
		rw.WriteString(fmt.Sprintf("%s", sendData))
		rw.Flush()
	}
}
func ReadData(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}
	}
}

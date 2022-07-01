package chat

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func SetupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, "pubsub-chatRoom", &discoveryNotifee{h: h})
	return s.Start()
}
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if n.h.ID() != pi.ID {
		fmt.Printf("\x1B[1;34mdiscovered new peer %s\n\x1b[0m> ", pi.ID.Pretty())
		err := n.h.Connect(context.Background(), pi)
		if err != nil {
			fmt.Printf("error connecting to peer %s: %s\n", pi.ID.Pretty(), err)
		}
	}
}

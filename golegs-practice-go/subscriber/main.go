package main

import (
	"Legs"
	"context"
	"flag"
	"fmt"
	"github.com/filecoin-project/go-legs"
	leveldb "github.com/ipfs/go-ds-leveldb"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
)

func main() {

	dest := flag.String("d", "", "destination address")
	flag.Parse()
	// blockstore
	db, err := leveldb.NewDatastore("./leveldb", nil)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	bs := blockstore.NewBlockstore(db)

	// libp2p host
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/6666"),
	)

	// linksystem
	ls := Legs.MkLinkSystem(bs)

	// subscriber
	sub, err := legs.NewSubscriber(host, db, ls, "/legs/hi", nil)
	defer sub.Close()
	if err != nil {
		panic(err)
	}
	destInfo, err := peer.AddrInfoFromString(*dest)
	host.Peerstore().AddAddr(destInfo.ID, destInfo.Addrs[0], peerstore.PermanentAddrTTL)
	err = host.Connect(ctx, *destInfo)
	if err != nil {
		panic(err)
	}

	watcher, cancelWatcher := sub.OnSyncFinished()
	defer cancelWatcher()
	go func() {
		for syncFin := range watcher {
			fmt.Println("Finished sync to", syncFin.Cid, "with peer:", syncFin.PeerID)
		}
	}()

	select {}
}

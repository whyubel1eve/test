package main

import (
	"Legs"
	"context"
	"fmt"
	"github.com/filecoin-project/go-legs/dtsync"
	"github.com/ipfs/go-cid"
	leveldb "github.com/ipfs/go-ds-leveldb"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	"github.com/ipld/go-ipld-prime"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/libp2p/go-libp2p"
	"log"
	"time"
)

func main() {
	// blockstore
	db, err := leveldb.NewDatastore("./leveldb", nil)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	bs := blockstore.NewBlockstore(db)

	// libp2p host
	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/5555"),
	)
	peerInfoStr := fmt.Sprintf("%s/p2p/%s", host.Addrs()[0], host.ID())
	log.Printf("Peer infomation: %s\n", peerInfoStr)

	// linksystem
	ls := Legs.MkLinkSystem(bs)

	pub, err := dtsync.NewPublisher(host, db, ls, "/legs/hi")
	if err != nil {
		panic(err)
	}

	// store a node
	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,    // Usually '1'.
		Codec:    0x71, // 0x71 means "dag-cbor"
		MhType:   0x13, // 0x13 means "sha2-512"
		MhLength: 64,   // sha2-512 hash has a 64-byte sum.
	}}

	np := basicnode.Prototype.Any
	nb := np.NewBuilder()
	ma, _ := nb.BeginMap(2)
	ma.AssembleKey().AssignString("hey")
	ma.AssembleValue().AssignString("it works")
	ma.AssembleKey().AssignString("yes")
	ma.AssembleValue().AssignString("true")
	ma.Finish()
	n := nb.Build()
	link, err := ls.Store(ipld.LinkContext{}, lp, n)
	if err != nil {
		panic(err)
	}
	// publish
	time.Sleep(time.Second * 15)
	err = pub.UpdateRoot(ctx, link.(cidlink.Link).Cid)
	if err != nil {
		panic(err)
	}
	log.Println("Publish: ", link.(cidlink.Link).Cid)
	select {}
}

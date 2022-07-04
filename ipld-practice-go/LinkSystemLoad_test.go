package ipld

import (
	"fmt"
	"os"
	"testing"

	"github.com/ipfs/go-cid"
	_ "github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/storage/memstore"
)

func TestLinkSystemLoad(t *testing.T) {

	var store memstore.Store
	ls := cidlink.DefaultLinkSystem()
	ls.SetWriteStorage(&store)

	lp := cidlink.LinkPrototype{Prefix: cid.Prefix{
		Version:  1,    // Usually '1'.
		Codec:    0x71, // 0x71 means "dag-cbor"
		MhType:   0x13, // 0x13 means "sha2-512"
		MhLength: 64,   // sha2-512 hash has a 64-byte sum.
	}}

	n := fluent.MustBuildMap(basicnode.Prototype.Map, 1, func(na fluent.MapAssembler) {
		na.AssembleEntry("hello").AssignString("world")
	})

	lnk, err := ls.Store(
		linking.LinkContext{}, // The zero value is fine.  Configure it it you want cancellability or other features.
		lp,                    // The LinkPrototype says what codec and hashing to use.
		n,                     // And here's our data.
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("link: %s\n", lnk)
	fmt.Printf("concrete type: `%T`\n", lnk)
	lnkNode := basicnode.NewLink(lnk)
	dagjson.Encode(lnkNode, os.Stdout)
	fmt.Println()

	//========================load========================

	cid, _ := cid.Decode("bafyrgqhai26anf3i7pips7q22coa4sz2fr4gk4q4sqdtymvvjyginfzaqewveaeqdh524nsktaq43j65v22xxrybrtertmcfxufdam3da3hbk")
	lnk = cidlink.Link{Cid: cid}
	lsys := cidlink.DefaultLinkSystem()

	lsys.SetReadStorage(&store)
	np := basicnode.Prototype.Any

	n, err = lsys.Load(
		linking.LinkContext{},
		lnk,
		np,
	)
	if err != nil {
		panic(err)
	}
	dagjson.Encode(n, os.Stdout)
	fmt.Println()

}

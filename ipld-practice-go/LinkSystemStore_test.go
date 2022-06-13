package ipld

import (
	"fmt"
	"testing"

	"github.com/ipfs/go-cid"
	_ "github.com/ipld/go-ipld-prime/codec/dagcbor" // 注意，这里一定要加
	"github.com/ipld/go-ipld-prime/fluent"
	"github.com/ipld/go-ipld-prime/linking"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/storage/memstore"
)

// linksystem link使用
func TestLinkSystemStore(t *testing.T) {
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
}

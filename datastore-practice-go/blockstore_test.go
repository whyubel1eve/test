package TestDatastore

import (
	"context"
	"fmt"
	"testing"

	blocks "github.com/ipfs/go-block-format"
	leveldb "github.com/ipfs/go-ds-leveldb"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
)

func TestLeveldbBlockstore(t *testing.T) {
	db, err := leveldb.NewDatastore("./leveldb", nil)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	bs := blockstore.NewBlockstore(db)
	block := blocks.NewBlock([]byte("Good things never fall apart"))
	c := block.Cid()
	bs.Put(ctx, block)
	has, err := bs.Has(ctx, c)
	if err != nil {
		panic(err)
	}
	if has {
		b, err:= bs.Get(ctx, c)
		if err != nil {
			panic(err)
		}
		size, err := bs.GetSize(ctx, c)
		if err != nil {
			panic(err)
		}
		fmt.Printf("value: %s, size: %d\n", b.RawData(), size)
	}
}
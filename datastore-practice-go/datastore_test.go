package TestDatastore

import (
	"context"
	"fmt"
	"testing"

	"github.com/ipfs/go-datastore"
	leveldb "github.com/ipfs/go-ds-leveldb"
)

func TestLeveldbDatastore(t *testing.T) {
	db, err := leveldb.NewDatastore("./leveldb", nil)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	k := datastore.NewKey("hello")
	db.Put(ctx, k, []byte("great!"))
	has, err := db.Has(ctx, k)
	if err != nil {
		panic(err)
	}
	if has {
		v, err := db.Get(ctx, k)
		if err != nil {
			panic(err)
		}
		size, err := db.GetSize(ctx, k)
		if err != nil {
			panic(err)
		}
		fmt.Printf("value: %s, size: %d\n", v, size)
	}
}

package TestDatastore

import (
	"fmt"
	"testing"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestLevelDB(t *testing.T) {
	db, err := leveldb.OpenFile("./db", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Put([]byte("hello"), []byte("world"), nil)
	data, err := db.Get([]byte("hello"), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}
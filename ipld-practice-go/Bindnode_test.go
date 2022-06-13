package ipld

import (
	"fmt"
	"os"
	"testing"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/bindnode"
)

// 从已有的go结构体创建node，schemaType为nil
func TestBindnode(t *testing.T) {
	type Person struct {
		Name    string
		Age     int64
		Friends []string
	}
	person := &Person{
		Name:    "kenny",
		Friends: []string{"simple", "niko"},
	}
	node := bindnode.Wrap(person, nil)

	fmt.Println(node.Type().TypeKind())
	nodeRepr := node.Representation()
	fmt.Println(nodeRepr.Kind())
	dagjson.Encode(nodeRepr, os.Stdout)

}

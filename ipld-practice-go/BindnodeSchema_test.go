package ipld

import (
	"fmt"
	"os"
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/bindnode"
)

// 从已有的go结构体创建node，有schemaType，并且要和Go类型兼容，否则根据其推断
func TestBindnodeSchema(t *testing.T) {
	ts, err := ipld.LoadSchemaBytes([]byte(`
		type Person struct {
			Name String
			Age optional Int
			Friends [String]
		}
	`))
	if err != nil {
		panic(err)
	}
	schemaType := ts.TypeByName("Person")

	type Person struct {
		Name    string
		Age     *int64
		Friends []string
	}
	person := &Person{
		Name:    "kenny",
		Friends: []string{"alex", "amos"},
	}
	node := bindnode.Wrap(person, schemaType)
	fmt.Println(node.Type().TypeKind())

	nodeRepr := node.Representation()
	dagjson.Encode(nodeRepr, os.Stdout)
}

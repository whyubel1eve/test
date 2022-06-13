package ipld

import (
	"os"
	"testing"

	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

// 使用builder创建node
func TestNodeBuilder(t *testing.T) {

	np := basicnode.Prototype.Any
	nb := np.NewBuilder()
	ma, _ := nb.BeginMap(2)
	ma.AssembleKey().AssignString("hey")
	ma.AssembleValue().AssignString("it works")
	ma.AssembleKey().AssignString("yes")
	ma.AssembleValue().AssignString("true")
	ma.Finish()
	n := nb.Build()

	dagjson.Encode(n, os.Stdout)
}

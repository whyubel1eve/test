package ipld

import (
	"fmt"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"testing"
)

func TestBasic(t *testing.T) {
	/*
		Kind_Int node test
	*/
	n1 := basicnode.NewInt(9)
	fmt.Println(n1.Kind())

	nb := basicnode.Prototype__Int{}.NewBuilder()
	nb.AssignInt(9)

	n2 := nb.Build()
	num, _ := n2.AsInt()
	fmt.Println(num)

}

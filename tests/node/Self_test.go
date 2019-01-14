package node

import (
  "fmt"
  "testing"
  
  "github.com/pokt-network/pocket-core/node"
)

func TestSelf(t *testing.T) {
  fmt.Println(node.GetSelf())
}

package dht

import (
  "fmt"
  "math/big"
)

type Node struct {
  id          []byte
  nodeId      string
  address     string
  port        string
  next        *Node
  finger[160] string
  fNode[160]  *Node
}


func makeDHTNode (id *string, addr string, port string) *Node {

  nod := new (Node)
  if id != nil {
    nod.nodeId = *id
    nod.id     = []byte(nod.nodeId)
  } else {
    nod.nodeId = generateNodeId()
    nod.id     = []byte(nod.nodeId)
  }
  nod.address = addr
  nod.port    = port
  nod.next    = nil

  for i := 0; i < 160; i++ {
    nod.finger[i],_ = calcFinger(nod.id, i, 160);
  }

  return nod
}


func (n *Node) updateFingers() {
  for i := 0; i < 160; i++ {
    n.fNode[i] = n.lookup(n.finger[i])
  }
}


func (n *Node) addToRing (new *Node) {

  if n.next == nil {  // First insert
    n.next   = new
    new.next = n

    return
  }

  for between(n.id, n.next.id, new.id) == false {
    n = n.next
  }
  new.next = n.next
  n.next   = new
}


func (n *Node) printRing() {

  nod := n
  for {
    nBigInt := big.Int{}
    nBigInt.SetString(nod.nodeId, 16)
    fmt.Printf("%s %s\n", nod.nodeId, nBigInt.String())
    nod = nod.next
    if nod == n {
      break
    }
  }
}


func (n *Node) lookup(key string) *Node {

  for between(n.id, n.next.id, []byte(key)) == false {
    n = n.next
  }
  return n
}


func (n *Node) mLookup(key string) (*Node, int) {

  j := 0
  for between(n.id, n.next.id, []byte(key)) == false {
    j++
    n = n.next
  }
  return n, j
}


func (n *Node) fLookup(key string) (*Node, int) {

  var nod *Node
  for i := 0; i < 160; i++ {
    if between([]byte(n.finger[i]), []byte(n.finger[(i + 1) % 160]), []byte(key)) {
//fmt.Printf("\n%s, %s,%i\n", n.finger[i], n.finger[(i + 1) % 160], i)
      nod = n.fNode[i]
      break;
//      return n.fNode[i]
    }
  }
  jumps := 0
  for between(nod.id, nod.next.id, []byte(key)) == false {
    jumps++
    nod = nod.next
  }
  return nod, jumps
}


func (nod *Node) testCalcFingers(n, m int) {

  calcFinger(nod.id, n, m);
}

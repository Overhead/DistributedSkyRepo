package dht

import (
  "fmt"
  "net"
  "bytes"
  "time"
  "encoding/json"
  "math/big"
  "strconv"
  "strings"
  "os"
)

type Node struct {
//  id      []byte
  nodeId       string
  localAddress string
  localPort    string
  localNode    *net.UDPAddr
  nextNode     *net.UDPAddr
  nextKey        string
  prevKey        string
  fingerVal[160] string
  fingerNod[160] *net.UDPAddr
}


type Message struct {
  Idx  int
  Key  string
  Src  *net.UDPAddr
  Dst  *net.UDPAddr
//     …
}


func DHTnode(id *string, addr string, port string) (*Node) {
// Main loop UDP server

  var buf [512]byte

  nod := new (Node)
  if id != nil {
    nod.nodeId = *id
  } else {
    nod.nodeId = generateNodeId()
  }
  nod.localAddress  = addr
  nod.localPort     = port
  nod.nextNode      = nil
  nod.prevKey       = ""
  nod.nextKey       = ""
  for i := 0; i < 160; i++ {
    nod.fingerVal[i],_ = calcFinger([]byte(nod.nodeId), i, 160)
    nod.fingerNod[i]   = nil
  }

fmt.Println("Creating node: ", nod.nodeId)
  service := addr + ":" + port
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err)
  nod.localNode = localAddr

  service = ":" + port
  localAddr, err = net.ResolveUDPAddr("udp", service)
  checkError(err)
  conn, err := net.ListenUDP("udp", localAddr)
  checkError(err)
fmt.Println("Node created: ", nod.nodeId)

  for {
    n, _, err := conn.ReadFromUDP(buf[0:])
    checkError(err)
//fmt.Printf("Received: (%s) %s \n", nod.localPort, buf[0:n])
    go handleRequest(n, buf, nod)
  }
}


func handleRequest(len int, buffer [512]byte, nod *Node) {
  dec := json.NewDecoder(bytes.NewReader([]byte(buffer[0:])))
  var msg Message
  err := dec.Decode(&msg)
  checkError(err)
//fmt.Println("Got: ", msg.Idx)

  switch msg.Idx {
  case 1:   // Join ring from outside, this is the joining node
            // We have only key value of joining node, find hash value for 
            // for reference node in ring and its' successor
    msg.Key = nod.nodeId
    nod1 := getRemoteInfo(nod, msg.Dst)
  // Now we have key value of node in nod1.Key and next node ptr in nod1.Dst
    if nod1.Dst == nil {
      // Now we have only one node in the ring and joining another
//fmt.Println("Single node")
      nod.nextNode = msg.Dst
//fmt.Println("nod: ", msg)
//fmt.Println("other: ", nod1)
      updateNextPointer(msg.Src, msg.Dst)
fmt.Println("Joining ready: ", nod.nodeId)
      return
    }
    nod2 := getRemoteInfo(nod, nod1.Dst)
  // Now we have successors Key value in nod2.Key (and successor in nod2.Dst)
    for between([]byte(nod1.Key), []byte(nod2.Key), []byte(msg.Key)) == false {
//fmt.Println("Selecting nextnode")
      nod1.Key = nod2.Key
      nod1.Dst = nod2.Dst
      nod1.Src = nod2.Src
      nod2 = getRemoteInfo(nod, nod1.Dst)
    }
    nod.nextNode = nod1.Dst
    updateNextPointer(nod.localNode, nod1.Src)
fmt.Println("Joining ready: ", nod.nodeId)
    initCalcFingers(nod.nodeId, nod.nextNode)
  case 2:   // Join ring, reply with nodeId
    conn, err := net.DialUDP("udp", nil, msg.Src)
    checkError(err)
    msg.Src = nod.localNode
    msg.Key = nod.nodeId
    msg.Dst = nod.nextNode
    defer conn.Close()
    buffer, err := json.Marshal(msg)
    checkError(err)
    _, err = conn.Write(buffer)
    checkError(err)
  case 3:   // Join ring, update nextNode pointer
    nod.nextNode = msg.Src
//fmt.Println("remote: ", nod)
  case 4:   // Printring, loop until back to beginning
    if strings.EqualFold(msg.Key, nod.nodeId) != true {
      nBigInt := big.Int{}
      nBigInt.SetString(nod.nodeId, 16)
      fmt.Printf("%s %s\n", nod.nodeId, nBigInt.String())
      if msg.Key == "" {
        msg.Key = nod.nodeId
      }
      conn, err := net.DialUDP("udp", nil, nod.nextNode)
      checkError(err)
      defer conn.Close()
      buffer, err := json.Marshal(msg)
      checkError(err)
      _, err = conn.Write(buffer)
      checkError(err)
    }
  case 5:    // Update fingers
    if msg.Key == "" {
      msg.Key = nod.nodeId
    }
    updateFingers(nod, msg.Key)
  case 6:    // Update prevKey
    nod.prevKey = msg.Key
  case 7:    // Lookup, first node
    if between([]byte(nod.prevKey), []byte(nod.nodeId), []byte(msg.Key)) {
      // I am responsible
fmt.Println("Responsible", nod.nodeId)
      return
    }
    // Not responsible, call next and wait...
    port, err := strconv.Atoi(nod.localPort)
    port += 210
    service := fmt.Sprintf("%s:%d", nod.localAddress, port)
    localAddr, err := net.ResolveUDPAddr("udp", service)
    checkError(err)
    msg.Src = localAddr
    lookupNextNode(nod, &msg)
    answ := waitForResult(port)
fmt.Println("Responsible", answ.Key)
  case 8:    // Lookup, following
    if between([]byte(nod.prevKey),[]byte(nod.nodeId),[]byte(msg.Key))==false {
    // Not responsible, just call next
      lookupNextNode(nod, &msg)
      return
    }
    // I am responsible
fmt.Println("Responsible", nod.nodeId)
    sendRespons(nod, &msg)
  case 9:    // dumpFingers
fmt.Println("dump>Finger")
    for i := 0; i < 160; i++ {
      fmt.Printf("%s  %s\n", nod.fingerVal[i], nod.fingerNod[i])
    }
  }
}


func waitForResult(port int) *Message {

  service := fmt.Sprintf(":%d", port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err)
  conn, err := net.ListenUDP("udp", localAddr)
  checkError(err)
  defer conn.Close()
  var buf [512]byte
  _, err = conn.Read(buf[0:])
  checkError(err)
//fmt.Printf("Received3: (%s) %s \n", nod.localPort, buf[0:n])

  dec  := json.NewDecoder(bytes.NewReader([]byte(buf[0:])))
  answ := new (Message)
  err = dec.Decode(&answ)
  checkError(err)
  return answ
}


func sendRespons(nod *Node, msg *Message) {

  msg.Key = nod.nodeId
  msg.Dst = msg.Src
  msg.Src = nod.localNode
  conn, err := net.DialUDP("udp", nil, msg.Dst)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)
}


func lookupNextNode(nod *Node, msg *Message) {

  nodLow  := []byte(nod.nodeId)
  nodHigh := []byte(nod.fingerVal[159])
  if nodHigh == nil {
    nodHigh = []byte(nod.nodeId)
  }
  i := 158
  for between(nodLow, nodHigh, []byte(msg.Key)) {
    nodHigh = []byte(nod.fingerVal[i])
    i--
    if i == 0 {
      break;
    }
  }
  if nod.fingerNod[i] != nil {
    msg.Dst = nod.fingerNod[i]
  } else {
    msg.Dst = nod.nextNode
  }
  msg.Idx = 8
  conn, err := net.DialUDP("udp", nil, msg.Dst)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)
}


func updateFingers(nod *Node, key string) {

fmt.Println("Updating fingers node: ", nod.nodeId)

//  updatePrevKey(nod.nodeId, nod.nextNode)
  msg        := getRemoteInfo(nod, nod.nextNode)
  nodLow     := []byte(nod.nodeId)
  nodHigh    := []byte(msg.Key)
  nod.nextKey = msg.Key
fmt.Println("msg.Key: ", msg.Key)
fmt.Println("msg.This: ", msg.Src)
fmt.Println("msg.Next: ", msg.Dst)

  mess := new (Message)
  next := msg.Dst
  mess.Src = msg.Src
  for i := 0; i < 160; i++ {
    for between(nodLow, nodHigh, []byte(nod.fingerVal[i])) == false {
      nodLow  = nodHigh
      mess    = getRemoteInfo(nod, next)
      nodHigh = []byte(mess.Key)
      next    = mess.Dst
    }
    nod.fingerNod[i] = mess.Src
  }

  if key != nod.nodeId {
fmt.Println("Key: ", msg.Key)
fmt.Println("Next: ", nod.nextNode)
    initCalcFingers(msg.Key, nod.nextNode)
  }

fmt.Println("Fingers done node: ", nod.nodeId)

}


func initCalcFingers(key string, dst *net.UDPAddr) {

fmt.Println("Init finger: ", key)
  time.Sleep(1000 * time.Millisecond)
  mess := new (Message)
  mess.Idx = 05
  mess.Key = key
  conn, err := net.DialUDP("udp", nil, dst)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(mess)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)
}


func updateNextPointer(src, dst *net.UDPAddr) {

  mess := new (Message)
  mess.Idx = 03
  mess.Key = ""
  mess.Src = src
  mess.Dst = dst
  conn, err := net.DialUDP("udp", nil, dst)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(mess)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)
}


func updatePrevKey(key string, dst *net.UDPAddr) {

  mess := new (Message)
  mess.Idx = 06
  mess.Key = key
  mess.Dst = dst
  conn, err := net.DialUDP("udp", nil, dst)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(mess)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)
}


func getRemoteInfo(nod *Node, dst *net.UDPAddr) *Message {

  mess := new (Message)
  mess.Idx = 02
  mess.Key = ""

  port, err := strconv.Atoi(nod.localPort)
  port += 110
  service := fmt.Sprintf("%s:%d", nod.localAddress, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err)

  mess.Src = localAddr
  mess.Dst = dst
  conn, err := net.DialUDP("udp", nil, dst)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(mess)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)

  service = fmt.Sprintf(":%d", port)
  localAddr, err = net.ResolveUDPAddr("udp", service)
  checkError(err)
  conn2, err := net.ListenUDP("udp", localAddr)
  checkError(err)
  defer conn2.Close()
  var buf [512]byte
  _, err = conn2.Read(buf[0:])
  checkError(err)
//fmt.Printf("Received2: (%s) %s \n", nod.localPort, buf[0:n])

  dec  := json.NewDecoder(bytes.NewReader([]byte(buf[0:])))
  answ := new (Message)
  err = dec.Decode(&answ)
  checkError(err)
  return answ
}


func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
    os.Exit(1)
  }
}

package dht

import (
  "fmt"
  "net"
  "bytes"
  "time"
  "encoding/json"
  "math/big"
  "math/rand"
  "strconv"
  "strings"
  "os"
  "sync"
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
  Idx     int
  Key     string
  Src     *net.UDPAddr
  Dst     *net.UDPAddr
  Status  int
  Version int
  Info    interface{}
  Gen     int
}


type portPool struct {
  mu      sync.Mutex
  baseVal int
  val     int
}


var answPort = new (portPool)


func DHTnode(id *string, addr string, port string) (*Node) {
// Main loop UDP server

  var buf [512]byte
  nod   := new (Node)
  table := createTable()
  if id != nil {
    nod.nodeId = *id
  } else {
    nod.nodeId = generateNodeId()
  }
fmt.Println("Creating node: ", nod.nodeId)
  nod.localAddress  = addr
  nod.localPort     = port
  nod.nextNode      = nil
  nod.prevKey       = ""
  nod.nextKey       = ""
  for i := 0; i < 160; i++ {
    nod.fingerVal[i],_ = calcFinger([]byte(nod.nodeId), i, 160)
    nod.fingerNod[i]   = nil
  }

  service := addr + ":" + port
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err)
  nod.localNode = localAddr

  service = ":" + port
  localAddr, err = net.ResolveUDPAddr("udp", service)
  checkError(err)
  conn, err := net.ListenUDP("udp", localAddr)
  checkError(err)
  portnum,_ := strconv.Atoi(nod.localPort)
  answPort.initPool(portnum)

fmt.Println("Node created: ", nod.nodeId)
  go doReplication(table, nod)
  for {
    n, _, err := conn.ReadFromUDP(buf[0:])
    checkError(err)
//fmt.Printf("Received: (%s) %s \n", node.localPort, buf[0:n])
    go handleRequest(table, n, buf, nod)
  }
}


func handleRequest(t *Table, len int, buffer [512]byte, nod *Node) {
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
    answ := initLookup(nod, &msg)
fmt.Println("Responsible-7 ", answ.Src)
  case 8:    // Lookup, following
    if between([]byte(nod.prevKey),[]byte(nod.nodeId),[]byte(msg.Key))==false {
    // Not responsible, just call next
      lookupNextNode(nod, &msg)
      return
    }
    // I am responsible
//fmt.Println("Responsible-8 ", msg.Src)
    sendRespons(nod, &msg)
  case 9:    //  init Insert
//fmt.Println("Init insert: ", msg.Info)
    answ := initLookup(nod, &msg)
    msg.Dst    = answ.Src
//_=dhtPing(nod, msg.Dst)
    msg.Src    = nod.localNode
    msg.Idx    = 13
    msg.Gen    = 0
    msg.Status = 3
    rnd        := rand.New(rand.NewSource(99))
    rndVal     := rnd.Intn(25000)
    msg.Version = rndVal
//fmt.Println("Resp-9: ", answ.Key)
    port := answPort.getPort() // += 210
    service := fmt.Sprintf("%s:%d", nod.localAddress, port)
    localAddr, err := net.ResolveUDPAddr("udp", service)
    checkError(err)
    msg.Src = localAddr
    doRemote(&msg)
    answ = waitForResult(port)
//fmt.Println("Insert return status", answ.Status)
  case 10:   //  init Get data
//fmt.Println("Init get", msg.Key)
    answ := initLookup(nod, &msg)
    msg.Dst = answ.Src
    msg.Idx = 14
//fmt.Println("Resp-10: ", answ.Key)
    port := answPort.getPort() // += 210
    service := fmt.Sprintf("%s:%d", nod.localAddress, port)
    localAddr, err := net.ResolveUDPAddr("udp", service)
    checkError(err)
    msg.Src = localAddr
    doRemote(&msg)
    answ = waitForResult(port)
//fmt.Println("Get return status", answ.Status)
fmt.Println("Get return data", answ.Info)
  case 11:   //  init Update
    answ := initLookup(nod, &msg)
    msg.Dst = answ.Src
    msg.Src = nod.localNode
    msg.Idx = 15
//fmt.Println("Resp-11: ", answ.Key)
    port := answPort.getPort()
    service := fmt.Sprintf("%s:%d", nod.localAddress, port)
    localAddr, err := net.ResolveUDPAddr("udp", service)
    checkError(err)
    msg.Src = localAddr
    doRemote(&msg)
    answ = waitForResult(port)
//fmt.Println("Get return status", answ.Status)
  case 12:   //  init mark as Deleted
//fmt.Println("Init delete", msg.Key)
    answ := initLookup(nod, &msg)
    msg.Dst = answ.Src
    msg.Idx = 16
//fmt.Println("Resp-12: ", answ.Key)
    port := answPort.getPort() // += 210
    service := fmt.Sprintf("%s:%d", nod.localAddress, port)
    localAddr, err := net.ResolveUDPAddr("udp", service)
    checkError(err)
    msg.Src = localAddr
    doRemote(&msg)
    answ = waitForResult(port)
//fmt.Println("Get return status", answ.Status)
  case 13:   //  do Insert
    stat := t.Insert(msg.Gen, msg.Status, msg.Version, msg.Key, msg.Info)
    msg.Status = stat
    msg.Info   = nil
//fmt.Println("Inserting on node: ", nod.nodeId)
    replyStat(&msg)
  case 14:   //  do Get data
    data, gen := t.getRecord(msg.Key)
//fmt.Println("Returning: ", data)
//fmt.Println("From: ", nod.nodeId)
    if gen < 0 {
      msg.Status = 0
      msg.Info   = nil
    } else {
      msg.Status = 1
      msg.Gen    = gen
      msg.Info   = data
    }
//fmt.Println("Fetching from node: ", nod.nodeId)
    replyStat(&msg)
  case 15:   //  do Update
    stat := t.update(msg.Gen, msg.Status, 0, msg.Key, msg.Info)
    msg.Status = stat
    msg.Info   = nil
//fmt.Println("Inserting on node: ", nod.nodeId)
    replyStat(&msg)
    invalidateCopy(nod, &msg)
  case 16:   // mark as Deleted
//fmt.Println("Removing: ", msg.Key)
    stat := t.Delete(msg.Key)
//fmt.Println("Rturning status ", stat, msg.Gen, nod.nodeId)
    if msg.Gen < 4 {
      markCopyAsDeleted(nod, &msg)
    }
    msg.Status = stat
    replyStat(&msg)
  case 17:   // Cont mark as deleted
//fmt.Println("mark copy as deleted: ", nod.nodeId)
    _ = t.Delete(msg.Key)
    if msg.Gen < 4 {
      markCopyAsDeleted(nod, &msg)
    }
  case 18:   // invalidate copy
//fmt.Println("invalidate copy: ", nod.nodeId)
    doInvalidateCopy(t, msg.Key)
    if msg.Gen < 4 {
      invalidateCopy(nod, &msg)
    }
  case 19:   // Get record status
    stat, vers := getRecordStatus(t, msg.Key)
    msg.Status  = stat
    msg.Version = vers
    replyStat(&msg)
  case 20:   // Update copy
//    if msg.Gen < 3 {
//      _ = updateCopy(nod, msg.Gen, msg.Version, msg.Key, msg.Info)
//    }
    stat := 3
    if msg.Gen >= 3 {
      stat = 2
    }
//fmt.Println("Updating: ", nod.nodeId, msg.Key)
    status := t.update(msg.Gen, stat, msg.Version, msg.Key, msg.Info)
    msg.Status = status
    replyStat(&msg)
  case 21:   // Sending back table at startup
//fmt.Println("Update records: ", nod.nodeId)
    t.updGetRecords(&msg)
  case 22:   // dhtPing
//fmt.Println("Got ping")
    msg.Idx = 23
    replyStat(&msg)
  case 23:   // Ping reply, do nothing
//fmt.Println("Ping reply")
  }
}


func replyStat(msg *Message) {

  conn, err := net.DialUDP("udp", nil, msg.Src)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)
}


func doRemote(msg *Message) {

  conn, err := net.DialUDP("udp", nil, msg.Dst)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)
}


func initLookup(nod *Node, msg *Message) *Message {

  if between([]byte(nod.prevKey), []byte(nod.nodeId), []byte(msg.Key)) {
    // I am responsible
//fmt.Println("Responsible-il ", nod.nodeId)
    msg.Src = nod.localNode
    msg.Gen = 0
    msg.Key = nod.nodeId
    return msg
  }
  // Not responsible, call next and wait...
  port := answPort.getPort()
  service := fmt.Sprintf("%s:%d", nod.localAddress, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err)
  msg.Src = localAddr

  lookupNextNode(nod, msg)
  answ := waitForResult(port)
  return answ
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
  answ.Info = nil
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

  dest := msg.Dst
  i = 0
  for ;dhtPing(nod, dest) == 0; i++ {
    dest = nod.fingerNod[i]
    if (dest == nil) {
      dest = msg.Dst
    }
  }
  msg.Dst = dest
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

//fmt.Println("Updating fingers node: ", nod.nodeId)

//  updatePrevKey(nod.nodeId, nod.nextNode)
  msg        := getRemoteInfo(nod, nod.nextNode)
  nodLow     := []byte(nod.nodeId)
  nodHigh    := []byte(msg.Key)
  nod.nextKey = msg.Key
  mess := new (Message)
  mess.Info = nil
  next := msg.Dst
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
//fmt.Println("Key: ", msg.Key)
//fmt.Println("Next: ", nod.nextNode)
    initCalcFingers(msg.Key, nod.nextNode)
  }

//fmt.Println("Fingers done node: ", nod.nodeId)

}


func initCalcFingers(key string, dst *net.UDPAddr) {

//fmt.Println("Init finger: ", key)
  time.Sleep(1000 * time.Millisecond)
  mess := new (Message)
  mess.Info = nil
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
  mess.Info = nil
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
  mess.Info = nil
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
  mess.Info = nil
  mess.Idx = 02
  mess.Key = ""

  port := answPort.getPort()
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
  answ.Info = nil
  err = dec.Decode(&answ)
  checkError(err)
  return answ
}


func (p *portPool) initPool (num int) {

  p.mu.Lock()
  num -= 1111
  base := num * 250
  p.baseVal = base + 2000
  p.val     = 0
  p.mu.Unlock()
}


func (p *portPool) getPort () int {

  p.mu.Lock()
  retval := p.baseVal + p.val
  p.val++
  p.val %= 200
  p.mu.Unlock()
  return retval
}


func dhtPing(nod *Node,dst *net.UDPAddr) int {

  msg := new (Message)
  msg.Idx  = 22
  msg.Dst  = dst
  msg.Info = nil
  port := answPort.getPort() // += 210
  service := fmt.Sprintf("%s:%d", nod.localAddress, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err)
  msg.Src = localAddr
  doCall(nod, msg)
  conn, err := net.DialUDP("udp", nil, msg.Dst)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)
//fmt.Println("Ping sent", nod.nodeId)

  conn, err = net.ListenUDP("udp", localAddr)
  checkError(err)
  tout := time.Duration(3 * time.Second)
  conn.SetDeadline(time.Now().Add(tout))
  defer conn.Close()
  var buf [512]byte
  _, err = conn.Read(buf[0:])
  if err == nil {
//fmt.Printf("Received ping: (%s)\n", nod.nodeId)
    dec  := json.NewDecoder(bytes.NewReader([]byte(buf[0:])))
    answ := new (Message)
    answ.Info = nil
    err = dec.Decode(&answ)
    checkError(err)
    return 1
  }
  checkError(err)
  return 0
}


func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
    os.Exit(1)
  }
}

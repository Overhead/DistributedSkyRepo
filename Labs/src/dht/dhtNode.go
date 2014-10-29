package dht

import (
  "fmt"
  "net"
  "bytes"
  "time"
  "encoding/json"
  "math/big"
  "math/rand"
//  "strconv"
  "strings"
  "os"
  "os/exec"
  "sync"
)

type Node struct {
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
var fingerStarted bool
var numCopies     int


func main() {
  DHTnode("", "0.0.0.0", "1075")
}


func DHTnode(id string, addr string, port string) (*Node) {
// Main loop UDP server

  var buf [512]byte
  numCopies     = 1
  fingerStarted = false
  nod   := new (Node)
  table := createTable()
  if id != "" {
    nod.nodeId = id
  } else {
    nod.nodeId = generateNodeId()
  }
fmt.Println("Creating node: ", nod.nodeId)
// Get local address
  address, err := exec.Command("hostname", "-I").Output()
  checkError(err, 99)
  addr = string(address)
fmt.Println("My IP: ", addr)
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
  checkError(err, 101)
  nod.localNode = localAddr

  service = ":" + port
  localAddr, err = net.ResolveUDPAddr("udp", service)
  checkError(err, 102)
  answPort.initPool()
  go startWebSocket(nod)
  go doStartupReplication(table, nod)

  conn, err := net.ListenUDP("udp", localAddr)
  checkError(err, 103)

fmt.Println("Node created: ", nod.nodeId)
//  sendToDaddy(nod)
  for {
    n, _, err := conn.ReadFromUDP(buf[0:])
    checkError(err, 104)
fmt.Println("Mess: ", buf)
    go handleRequest(table, n, buf, nod)
  }
}


func handleRequest(t *Table, len int, buffer [512]byte, nod *Node) {
  dec := json.NewDecoder(bytes.NewReader([]byte(buffer[0:])))
  var msg Message
  err := dec.Decode(&msg)
  checkError(err, 105)
fmt.Println("Got: ", msg)

  switch msg.Idx {
  case 1:   // Join ring from outside, this is the joining node
            // We have only key value of joining node, find hash value
            // for reference node in ring and its' successor
    msg.Key = nod.nodeId
    nod1 := getRemoteInfo(nod, msg.Dst)
    nod.nextKey = nod1.Key
  // Now we have key value of node in nod1.Key and next node ptr in nod1.Dst
    if nod1.Dst == nil {
      // Now we have only one node in the ring and joining another
      nod.nextNode = msg.Dst
      updateNextPointer(msg.Src, msg.Dst)
      nod.prevKey = msg.Key
      updatePrevKey(nod.nodeId, msg.Dst)
      return
    }
    nod2 := getRemoteInfo(nod, nod1.Dst)
  // Now we have successors Key value in nod2.Key (and successor in nod2.Dst)
    for between([]byte(nod1.Key), []byte(nod2.Key), []byte(msg.Key)) == false {
      nod1.Key = nod2.Key
      nod1.Dst = nod2.Dst
      nod1.Src = nod2.Src
      nod2 = getRemoteInfo(nod, nod1.Dst)
    }
    // We will add this node between nod1 and nod2
    nod.nextNode = nod1.Dst
    updateNextPointer(msg.Src, nod1.Src)
    nod.prevKey  = nod1.Key
    updatePrevKey(nod.nodeId, nod1.Dst)
  case 2:   // Join ring, reply with nodeId
    conn, err := net.DialUDP("udp", nil, msg.Src)
    checkError(err, 106)
    msg.Src = nod.localNode
    msg.Key = nod.nodeId
    msg.Dst = nod.nextNode
    defer conn.Close()
    buffer, err := json.Marshal(msg)
    checkError(err, 107)
    _, err = conn.Write(buffer)
    checkError(err, 108)
  case 3:   // Join ring, update nextNode pointer
    nod.nextNode = msg.Src
  case 4:   // Printring, loop until back to beginning
    mess := new (Message)
    mess.Idx  = 4
    mess.Info = nil
    mess.Key  = ""
    mess.Dst  = msg.Src
    mess.Src  = nod.localNode
    if strings.EqualFold(msg.Key, nod.nodeId) != true {
      nBigInt := big.Int{}
      nBigInt.SetString(nod.nodeId, 16)
      mess.Key  = nod.nodeId
      doRemote(mess)
      if msg.Key == "" {
        msg.Key = nod.nodeId
      }
      if nod.nextNode != nil {
        time.Sleep(100 * time.Millisecond)
        conn, err := net.DialUDP("udp", nil, nod.nextNode)
        checkError(err, 109)
        defer conn.Close()
        buffer, err := json.Marshal(msg)
        checkError(err, 110)
        _, err = conn.Write(buffer)
        checkError(err, 111)
      } else {
        mess.Key = ""
        doRemote(mess)
      }
      return
    } else {
      doRemote(mess)
    }
  case 5:    // Starting replication and fingering
    if fingerStarted == false {
      fingerStarted = true
      go maintainNode(t, nod)  // Fingers and replication
    }
  case 6:    // Update prevKey
    nod.prevKey = msg.Key
  case 7:    // Lookup, first node
    if fingerStarted == false {
      fingerStarted = true
      go maintainNode(t, nod)  // Fingers and replication
    }
    retAddr := msg.Src
    answ    := initLookup(nod, &msg)
fmt.Println("Responsible-7 ", answ.Src)
    msg.Src = answ.Src
    msg.Key = answ.Key
    msg.Dst = retAddr
    doRemote(&msg)    
  case 8:    // Lookup, following
    if between([]byte(nod.prevKey),[]byte(nod.nodeId),[]byte(msg.Key))==false {
    // Not responsible, just call next
      lookupNextNode(nod, &msg)
      return
    }
    // I am responsible
fmt.Println("Responsible-8 ", msg.Src)
    sendRespons(nod, &msg)
  case 9:    //  init Insert
fmt.Println("Init insert-9: ", msg)
    retAddr := msg.Src
    key     := msg.Key
    answ    := initLookup(nod, &msg)
    dstAddr := answ.Src
    msg.Dst    = answ.Src
    msg.Src    = nod.localNode
    msg.Idx    = 13
    msg.Gen    = 0
    msg.Status = 3
    msg.Key    = key
    rnd        := rand.New(rand.NewSource(99))
    rndVal     := rnd.Intn(25000)
fmt.Println("Resp-9: ", answ)
fmt.Println("Msg-9: ", msg)
    msg.Version = rndVal
    port := answPort.getPort() // += 210
    service := fmt.Sprintf("%s:%d", nod.localAddress, port)
    localAddr, err := net.ResolveUDPAddr("udp", service)
    checkError(err, 112)
    msg.Src = localAddr
    doRemote(&msg)
    answ = waitForResult(port)
fmt.Println("Insert return status-9", answ.Status)
    msg.Src    = dstAddr
    msg.Dst    = retAddr
    msg.Status = answ.Status
    msg.Info   = nil
    doRemote(&msg)
  case 10:   //  init Get data
fmt.Println("Init get-10", msg.Key)
    retAddr := msg.Src
    key     := msg.Key
    answ    := initLookup(nod, &msg)
    dstAddr := answ.Src
    msg.Dst = answ.Src
    msg.Idx = 14
fmt.Println("Resp-10: ", answ.Key)
    port := answPort.getPort() // += 210
    service := fmt.Sprintf("%s:%d", nod.localAddress, port)
    localAddr, err := net.ResolveUDPAddr("udp", service)
    checkError(err, 113)
    msg.Src = localAddr
    msg.Key = key
    doRemote(&msg)
    answ = waitForResult(port)
fmt.Println("Get return status-10", answ.Status)
fmt.Println("Get return data-10", answ.Info)
    msg.Src    = dstAddr
    msg.Dst    = retAddr
    msg.Status = answ.Status
    msg.Info   = answ.Info
    doRemote(&msg)
  case 11:   //  init Update
fmt.Println("Upd-11: ", msg)
    retAddr   := msg.Src
    key       := msg.Key
    answ      := initLookup(nod, &msg)
    dstAddr   := answ.Src
    msg.Dst    = answ.Src
    msg.Src    = nod.localNode
    msg.Idx    = 15
    msg.Gen    = 0
    msg.Status = 2
fmt.Println("Resp-11: ", answ.Key)
    port := answPort.getPort()
    service := fmt.Sprintf("%s:%d", nod.localAddress, port)
    localAddr, err := net.ResolveUDPAddr("udp", service)
    checkError(err, 114)
    msg.Src = localAddr
    msg.Key = key
    doRemote(&msg)
    answ = waitForResult(port)
    msg.Src    = dstAddr
    msg.Dst    = retAddr
    msg.Status = answ.Status
    msg.Info   = nil
    doRemote(&msg)
fmt.Println("Get return status", answ.Status)
  case 12:   //  init mark as Deleted
fmt.Println("Init delete-12", msg.Key)
    retAddr := msg.Src
    answ    := initLookup(nod, &msg)
    dstAddr := answ.Src
    msg.Dst  = answ.Src
    msg.Idx  = 16
    msg.Gen  = 0
fmt.Println("Resp-12: ", answ.Key)
    port := answPort.getPort() // += 210
    service := fmt.Sprintf("%s:%d", nod.localAddress, port)
    localAddr, err := net.ResolveUDPAddr("udp", service)
    checkError(err, 115)
    msg.Src    = localAddr
    msg.Status = answ.Status
    doRemote(&msg)
    answ = waitForResult(port)
    msg.Src    = dstAddr
    msg.Dst    = retAddr
    msg.Status = answ.Status
    msg.Info   = nil
    doRemote(&msg)
fmt.Println("Get return status", answ.Status)
  case 13:   //  do Insert
fmt.Println("Ins-13: ", msg)
    stat := t.Insert(msg.Gen, msg.Status, msg.Version, msg.Key, msg.Info)
    msg.Status = stat
    msg.Info   = nil
fmt.Println("Inserting on node-13: ", nod.nodeId)
    replyStat(&msg)
  case 14:   //  do Get data
fmt.Println("Key-14: ", msg.Key)
    data, gen := t.getRecord(msg.Key)
fmt.Println("Returning-14: ", data, gen)
fmt.Println("From-14: ", nod.nodeId)
    if gen < 0 {
      msg.Status = 0
      msg.Info   = nil
    } else {
      msg.Status = 1
      msg.Gen    = gen
      msg.Info   = data
    }
    msg.Key = nod.nodeId
fmt.Println("Fetching from node-14: ", nod.nodeId)
    replyStat(&msg)
  case 15:   //  do Update
fmt.Println("Upd-15: ", msg)
    key  := msg.Key
    stat := t.update(msg.Gen, msg.Status, 0, msg.Key, msg.Info)
    msg.Status = stat
    msg.Info   = nil
    msg.Key    = nod.nodeId
fmt.Println("Updating on node: ", nod.nodeId, msg)
    msg.Key    = key
    invalidateCopy(nod, &msg)
    replyStat(&msg)
  case 16:   // mark as Deleted
fmt.Println("Removing-16: ", msg.Key)
    stat    := t.Delete(msg.Key)
fmt.Println("Rturning status-16 ", stat, msg.Gen, nod.nodeId)
    if msg.Gen <= numCopies {
      markCopyAsDeleted(nod, &msg)
    }
    msg.Status = stat
    msg.Key    = nod.nodeId
    replyStat(&msg)
  case 17:   // Cont mark as deleted
fmt.Println("mark copy as deleted-17: ", nod.nodeId)
    _ = t.Delete(msg.Key)
    if msg.Gen <= numCopies {
      markCopyAsDeleted(nod, &msg)
    }
  case 18:   // invalidate copy
fmt.Println("invalidate copy: ", nod.nodeId)
    doInvalidateCopy(t, msg.Key)
    if msg.Gen <= numCopies {
      invalidateCopy(nod, &msg)
    }
  case 19:   // Get record status
    stat, vers := getRecordStatus(t, msg.Key)
    msg.Status  = stat
    msg.Version = vers
    replyStat(&msg)
  case 20:   // Update copy
//    if msg.Gen < numCopies {
//      _ = updateCopy(nod, msg.Gen, msg.Version, msg.Key, msg.Info)
//    }
    stat := 3
    if msg.Gen >= numCopies {
      stat = 2
    }
fmt.Println("Updating: ", nod.nodeId, msg.Key)
    status := t.update(msg.Gen, stat, msg.Version, msg.Key, msg.Info)
    msg.Status = status
    replyStat(&msg)
  case 21:   // Sending back table at startup
fmt.Println("Update records: ", nod.nodeId)
    t.updGetRecords(&msg)
  case 22:   // dhtPing
//fmt.Println("Got ping")
    msg.Idx = 23
    replyStat(&msg)
  case 23:   // Ping reply, do nothing
//fmt.Println("Ping reply")
  case 24:   // dump fingers
fmt.Println("dump>Finger")
    mess := new (Message)
    mess.Idx  = 24
    mess.Info = nil
    mess.Dst  = msg.Src
    for i := 0; i < 160; i++ {
fmt.Printf("%s  %s\n", nod.fingerVal[i], nod.fingerNod[i])
      mess.Key = nod.fingerVal[i]
      mess.Src = nod.fingerNod[i]
      mess.Gen = i
      doRemote(mess)
      time.Sleep(100 * time.Millisecond)
    }
    mess.Key = ""
    doRemote(mess)
  }
}


func sendToDaddy(nod *Node) {

  msg := new (Message)
  msg.Key  = nod.nodeId
  msg.Info = nil
  msg.Src  = nod.localNode
  localAddr, err := net.ResolveUDPAddr("udp", "172.17.42.1:12000")
  checkError(err, 199)
  msg.Dst  = localAddr
  doRemote(msg)
}


func replyStat(msg *Message) {

  conn, err := net.DialUDP("udp", nil, msg.Src)
  checkError(err, 116)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err, 117)
  _, err = conn.Write(buffer)
  checkError(err, 118)
}


func doRemote(msg *Message) {

fmt.Println("Sending to: ", msg)
  conn, err := net.DialUDP("udp", nil, msg.Dst)
  checkError(err, 119)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err, 120)
  _, err = conn.Write(buffer)
  checkError(err, 121)
}


func initLookup(nod *Node, msg *Message) *Message {

  if between([]byte(nod.prevKey), []byte(nod.nodeId), []byte(msg.Key)) {
    // I am responsible
fmt.Println("Responsible-il ", nod.nodeId)
    msg.Src = nod.localNode
    msg.Gen = 0
    msg.Key = nod.nodeId
    return msg
  }
  // Not responsible, call next and wait...
  port := answPort.getPort()
  service := fmt.Sprintf("%s:%d", nod.localAddress, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 122)
  msg.Src = localAddr
fmt.Println("Check next node")
  lookupNextNode(nod, msg)
  answ := waitForResult(port)
  return answ
}


func waitForResult(port int) *Message {

  service := fmt.Sprintf(":%d", port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 123)
  conn, err := net.ListenUDP("udp", localAddr)
  checkError(err, 124)
  defer conn.Close()
  var buf [512]byte
  _, err = conn.Read(buf[0:])
  checkError(err, 125)
//fmt.Printf("Received3: (%s) %s \n", nod.localPort, buf[0:n])

  dec  := json.NewDecoder(bytes.NewReader([]byte(buf[0:])))
  answ := new (Message)
  answ.Info = nil
  err = dec.Decode(&answ)
  checkError(err, 126)
  return answ
}


func sendRespons(nod *Node, msg *Message) {

  msg.Key = nod.nodeId
  msg.Dst = msg.Src
  msg.Src = nod.localNode
  conn, err := net.DialUDP("udp", nil, msg.Dst)
  checkError(err, 127)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err, 128)
  _, err = conn.Write(buffer)
  checkError(err, 129)
}


func lookupNextNode(nod *Node, msg *Message) {

  nodLow  := []byte(nod.nodeId)
  nodHigh := []byte(nod.fingerVal[159])
  dest    := nod.nextNode
  i       := 0
  if nodHigh != nil {
    i := 158
    for between(nodLow, nodHigh, []byte(msg.Key)) {
      nodHigh = []byte(nod.fingerVal[i])
      if nodHigh == nil {
        i = 1    
      }
      i--
      if i == 0 {
        break;
      }
    }
    if nod.fingerNod[i] != nil {
      dest = nod.fingerNod[i]
    }
  }

  for dhtPing(nod, dest) == false {
    for j := i;; i++ {
      if nod.fingerNod[j] != dest {
        j++
      }
      if j == 160 {
        dest = nil
        break
      }
    }
    if dest == nil {
      break
    }
  }

fmt.Println("Next node: ", nod.nextNode, msg.Dst)
msg.Dst = nod.nextNode

  if dest == nil {
    msg.Idx  = 0
    msg.Info = nil
    msg.Key  = ""
    msg.Dst  = msg.Src
  } else {
    msg.Dst = dest
  }
  msg.Idx = 8
  conn, err := net.DialUDP("udp", nil, msg.Dst)
  checkError(err, 130)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err, 131)
  _, err = conn.Write(buffer)
  checkError(err, 132)
}


func maintainNode(t *Table, nod *Node) {

  time.Sleep(5 * time.Second)
  for {
    checkFingers(nod)
    time.Sleep(5 * time.Second)
    doReplication(t, nod)
    time.Sleep(55 * time.Second)
  }
}


func checkFingers(nod *Node) {

fmt.Println("Check fingers: ", nod.nodeId)
  // Check if nextnode is unreacheable and find a new one
  if dhtPing(nod, nod.nextNode) == false {
    tmp := nod.nextNode
    var i int
    for i = 0; i < 160; i++ {
      if nod.fingerNod[i] != tmp {
        tmp = nod.fingerNod[i]
        if dhtPing(nod, tmp) == true {
          nod.nextNode = tmp
          break
        }
      }
    }
    if (i == 160) || (tmp == nil) {
      // No next node found alive exit
      os.Exit(1)
    }
  }

fmt.Println("Updating fingers node: ", nod.nodeId)
  updatePrevKey(nod.nodeId, nod.nextNode)
  msg        := getRemoteInfo(nod, nod.nextNode)
  nodLow     := []byte(nod.nodeId)
  nodHigh    := []byte(msg.Key)
  nod.nextKey = msg.Key
  next := msg.Dst
  for i := 0; i < 160; i++ {
    // Fingerval is not between nodes, try next
    for between(nodLow, nodHigh, []byte(nod.fingerVal[i])) == false {
      nodLow  = nodHigh
      msg     = getRemoteInfo(nod, next)
      nodHigh = []byte(msg.Key)
      next    = msg.Dst
    }
    nod.fingerNod[i] = msg.Src
  }
fmt.Println("Fingers done node: ", nod.nodeId)

  var tmp *net.UDPAddr
  tmp = nil
  for i := 0; i < 160; i++ {
    if nod.fingerNod[i] != tmp {
      tmp = nod.fingerNod[i]
      fmt.Println("", tmp, nod.fingerNod[i])
    }
  }
}


func updateNextPointer(src, dst *net.UDPAddr) {

  mess := new (Message)
  mess.Info = nil
  mess.Idx = 03
  mess.Key = ""
  mess.Src = src
  mess.Dst = dst
  conn, err := net.DialUDP("udp", nil, dst)
  checkError(err, 136)
  defer conn.Close()
  buffer, err := json.Marshal(mess)
  checkError(err, 137)
  _, err = conn.Write(buffer)
  checkError(err, 138)
}


func updatePrevKey(key string, dst *net.UDPAddr) {

  mess := new (Message)
  mess.Info = nil
  mess.Idx = 06
  mess.Key = key
  mess.Dst = dst
  conn, err := net.DialUDP("udp", nil, dst)
  checkError(err, 139)
  defer conn.Close()
  buffer, err := json.Marshal(mess)
  checkError(err, 140)
  _, err = conn.Write(buffer)
  checkError(err, 141)
}


func getRemoteInfo(nod *Node, dst *net.UDPAddr) *Message {

  mess := new (Message)
  mess.Info = nil
  mess.Idx = 02
  mess.Key = ""

  port := answPort.getPort()
  service := fmt.Sprintf("%s:%d", nod.localAddress, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
fmt.Println("Addr: ", localAddr, service)
  checkError(err,142)

  mess.Src = localAddr
  mess.Dst = dst
  conn, err := net.DialUDP("udp", nil, dst)
  checkError(err, 143)
  defer conn.Close()
  buffer, err := json.Marshal(mess)
  checkError(err, 144)
  _, err = conn.Write(buffer)
  checkError(err, 145)

  service = fmt.Sprintf(":%d", port)
  localAddr, err = net.ResolveUDPAddr("udp", service)
  checkError(err, 146)
  conn2, err := net.ListenUDP("udp", localAddr)
  checkError(err, 147)
  defer conn2.Close()
  var buf [512]byte
  n, err := conn2.Read(buf[0:])
  checkError(err, 148)
fmt.Printf("Received2: (%s) %s \n", nod.localPort, buf[0:n])

  dec  := json.NewDecoder(bytes.NewReader([]byte(buf[0:])))
  answ := new (Message)
  answ.Info = nil
  err = dec.Decode(&answ)
  checkError(err, 149)
  return answ
}


func (p *portPool) initPool () {

  p.mu.Lock()
  p.baseVal = 8000
  p.val     = 0
  p.mu.Unlock()
}


func (p *portPool) getPort () int {

  p.mu.Lock()
  retval := p.baseVal + p.val
  p.val++
  p.val %= 25
  p.mu.Unlock()
  return retval
}


func dhtPing(nod *Node,dst *net.UDPAddr) bool {

  if dst == nil {
    return true
  }
  msg := new (Message)
  msg.Idx  = 22
  msg.Dst  = dst
  msg.Info = nil
  port := answPort.getPort() // += 210
  service := fmt.Sprintf("%s:%d", nod.localAddress, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 150)
  msg.Src = localAddr
//  doCall(nod, msg)
  conn, err := net.DialUDP("udp", nil, msg.Dst)
  checkError(err, 151)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err, 152)
//fmt.Println("Ping sent", nod.nodeId)
  service = fmt.Sprintf(":%d", port)
  localAddr, err = net.ResolveUDPAddr("udp", service)
  checkError(err, 1520)
  _, err = conn.Write(buffer)
  checkError(err, 153)

  conn, err = net.ListenUDP("udp", localAddr)
  checkError(err, 154)
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
    checkError(err, 155)
    return true
  }
//  checkError(err, 156)
fmt.Println("Ping failed with: ", dst)
  return false
}


func checkError(err error, i int) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error ", i, err.Error())
    os.Exit(1)
  }
}

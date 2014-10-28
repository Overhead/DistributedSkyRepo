package dht

import (
  "code.google.com/p/go.net/websocket"
  "fmt"
  "log"
  "net"
  "bytes"
  "strings"
  "net/http"
  "encoding/json"
)

type Msg struct {
  Action int
  Key    string
  NewKey string
  Date   int64
}


var port      int
var node      *Node

func startWebSocket(nod *Node) {
  node = nod
  port = 1075
  http.Handle("/node", websocket.Handler(nodeHandler))
  http.Handle("/manage", websocket.Handler(manageHandler))
fmt.Println("Starting websocket")
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic("ListenAndServe: " + err.Error())
  }
}


func manageHandler(ws *websocket.Conn) {
  msg := make([]byte, 512)
  n, err := ws.Read(msg)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("Receive: %s\n", msg[:n])
  var res Msg
  json.Unmarshal([]byte(msg[:n]), &res)

  fmt.Println(res)
  fmt.Println(res.Action)
  fmt.Println(res.Key)
  fmt.Println(res.Date)

  m, err := ws.Write(msg[:n])
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("Send: %s\n", msg[:m])
}


func nodeHandler(ws *websocket.Conn) {
  msg := make([]byte, 512)
  n, err := ws.Read(msg)
  if err != nil {
    checkError(err, 500)
//    log.Fatal(err)
  }
fmt.Printf("Receive: %s\n", msg[:n])
  var res Msg
  json.Unmarshal([]byte(msg[:n]), &res)

  fmt.Println(res)
  fmt.Println(res.Action)
  fmt.Println(res.Key)
  fmt.Println(res.Date)


  var mess *Message
  switch res.Action {
  case 1:  // Get  - Search
    mess = doGet(&res)
  case 2:  // Put  - update
    mess = doPut(&res)
  case 3:  // Post - insert
    mess = doPost(&res)
  case 4:  // Del  - Del
    mess = doDel(&res)
  }

  switch mess.Status {
  case 0:
    msg = []byte(fmt.Sprintf("Failure"))
  case 1:
    msg = []byte(fmt.Sprintf("Success"))
  }
  n = 7
  if (res.Action == 1) && (mess.Status == 1) {
    result := fmt.Sprintf("%v", mess.Info)
    msg = []byte(result)
    n   = len(result)
  }
/* Code to do printring
result := doPrintring()
msg = []byte(result)
n   = len(result)
*/
  m, err := ws.Write(msg[:n])
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("Send: %s\n", msg[:m])
}


func doPrintring() string {   // Search
  msg := new (Message)
  msg.Idx  = 04
  msg.Key  = ""
  msg.Info = nil
  msg.Dst  = node.localNode
  remote, listen := doReplyAddresses(node.localAddress)
  msg.Src = remote
  sendMsg(msg)
  conn := setListen(listen)
  var buffer bytes.Buffer
  for {
    answ := getRecord(conn)
    if answ.Key == "" {
      conn.Close()
      break
    }
    buffer.WriteString(fmt.Sprintf("%s ", answ.Key))
  }
  buff := strings.Replace(buffer.String()," ", "#", -1)
  return strings.Trim(buff, "#")
}


func doGet(res *Msg) *Message {   // Search
  hashk := sha1hash(res.Key)
  msg := new (Message)
  msg.Idx  = 10
  msg.Key  = hashk
  msg.Info = nil
  msg.Dst  = node.localNode
  remote, listen := doReplyAddresses(node.localAddress)
  msg.Src = remote
  sendMsg(msg)
  conn := setListen(listen)
  answ := getRecord(conn)
  return answ
}


func doPut(res *Msg) *Message {  // Update
  hashk := sha1hash(res.Key)
  msg := new (Message)
  msg.Idx  = 11
  msg.Key  = hashk
  msg.Info = res.NewKey
  msg.Dst  = node.localNode
  remote, listen := doReplyAddresses(node.localAddress)
  msg.Src = remote
  sendMsg(msg)
  conn := setListen(listen)
  answ := getRecord(conn)
  return answ
}


func doPost(res *Msg) *Message {   // Insert
  hashk := sha1hash(res.Key)
  msg := new (Message)
  msg.Idx  = 9
  msg.Key  = hashk
  msg.Info = res.Key
  msg.Dst  = node.localNode
  remote, listen := doReplyAddresses(node.localAddress)
  msg.Src = remote
  sendMsg(msg)
  conn := setListen(listen)
  answ := getRecord(conn)
  return answ
}


func doDel(res *Msg) *Message {   // Delete
  hashk := sha1hash(res.Key)
  msg := new (Message)
  msg.Idx  = 12
  msg.Key  = hashk
  msg.Info = nil
  msg.Dst  = node.localNode
  remote, listen := doReplyAddresses(node.localAddress)
  msg.Src = remote
  sendMsg(msg)
  conn := setListen(listen)
  answ := getRecord(conn)
  return answ
}


func sendMsg(msg *Message) {
  conn, err := net.DialUDP("udp", nil, msg.Dst)
  checkError(err, 501)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err, 502)
   _, err = conn.Write(buffer)
  checkError(err, 503)
}


func doReplyAddresses(addr string) (*net.UDPAddr, *net.UDPAddr) {
  port    := answPort.getPort()
  service := fmt.Sprintf("%s:%d", addr, port)
  localAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 504)
  service     = fmt.Sprintf(":%d", port)
  local, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 505)
  return localAddr, local
}


func setListen(addr *net.UDPAddr) *net.UDPConn {
  conn, err := net.ListenUDP("udp", addr)
  checkError(err, 506)
  return conn
}


func getRecord(conn *net.UDPConn) *Message {
  var buf [512]byte
  _, err := conn.Read(buf[0:])
  checkError(err, 507)
  dec  := json.NewDecoder(bytes.NewReader([]byte(buf[0:])))
  answ := new (Message)
  err   = dec.Decode(&answ)
  checkError(err, 508)
  return answ
}

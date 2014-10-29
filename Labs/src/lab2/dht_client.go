package dht

import (
  "fmt"
  "net"
  "encoding/json"
  "time"
)

type Host struct {
  id          []byte
  address     string
  port        string
}


func (n *Host) String() string {
  return fmt.Sprintf("%x", n.id)
}


func makeDHTNode(id *string, addr string, port string) *net.UDPAddr {

  service := fmt.Sprintf("%s:%s", addr, port)

  udpAddr, err := net.ResolveUDPAddr("udp", service)
  checkError(err)

  go DHTnode(id, addr, port)

  time.Sleep(500 * time.Millisecond)
  return udpAddr
}


func addToRing(nod, new *net.UDPAddr) {

  conn, err := net.DialUDP("udp", nil, new)
  checkError(err)

  var msg Message
  msg.Idx = 01
  msg.Key = ""
  msg.Src = new
  msg.Dst = nod

  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)

//fmt.Println("Sending: ", buffer[0:])
   _, err = conn.Write(buffer)
  checkError(err)

  time.Sleep(2000 * time.Millisecond)
}


func printRing(nod *net.UDPAddr) {

  conn, err := net.DialUDP("udp", nil, nod)
  checkError(err)

  var msg Message
  msg.Idx = 04
  msg.Key = ""
  msg.Src = nil
  msg.Dst = nod

  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)

//fmt.Println("Sending: ", buffer[0:])
   _, err = conn.Write(buffer)
  checkError(err)
}


func lookup(nod *net.UDPAddr, key string) {

  conn, err := net.DialUDP("udp", nil, nod)
  checkError(err)

  var msg Message
  msg.Idx = 07
  msg.Key = key
  msg.Dst = nod

  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)

   _, err = conn.Write(buffer)
  checkError(err)
}


func testDumpFingers(nod *net.UDPAddr) {

  conn, err := net.DialUDP("udp", nil, nod)
  checkError(err)

  var msg Message
  msg.Idx = 9
  msg.Key = ""
  msg.Src = nil
  msg.Dst = nod

  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)

fmt.Println("Sending: ", buffer[0:])
   _, err = conn.Write(buffer)
  checkError(err)
}


func initCalc(nod *net.UDPAddr) {

  conn, err := net.DialUDP("udp", nil, nod)
  checkError(err)

  var msg Message
  msg.Idx = 05
  msg.Key = ""

  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)

   _, err = conn.Write(buffer)
  checkError(err)
}


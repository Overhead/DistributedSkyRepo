package main

import (
  "fmt"
  "net"
  "encoding/json"
  "os"
)


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


// Call with:  join <node in ring> <node to join the ring> [<portnum>]
func main () {

  if len(os.Args) < 3 {
    fmt.Println("To few args")
    os.Exit(1)
  }
  src  := os.Args[1]
  dest := os.Args[2]
  port := "1075"
  if len(os.Args) == 4 {
    port = os.Args[3]
  }

  service  := src + ":" + port
  old, err := net.ResolveUDPAddr("udp", service)
  checkError(err)

  service   = dest + ":" + port
  new, err := net.ResolveUDPAddr("udp", service)

  var msg Message
  msg.Idx = 01
  msg.Key = ""
  msg.Src = new
  msg.Dst = old

  conn, err := net.DialUDP("udp", nil, new)
  checkError(err)

  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)

//fmt.Println("Sending: ", buffer[0:])
   _, err = conn.Write(buffer)
  checkError(err)
}


func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
    os.Exit(1)
  }
}

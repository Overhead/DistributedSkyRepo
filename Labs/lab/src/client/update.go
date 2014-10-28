package main

import (
  "fmt"
  "net"
  "encoding/json"
  "os"
  "os/exec"
  "strings"
  "bytes"
  "crypto/sha1"
  "github.com/nu7hatch/gouuid"
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


// Call with:  update <node in ring> <insert key> <insert string> [<portnum>]
func main () {

  if len(os.Args) < 4 {
    fmt.Println("To few args")
    os.Exit(1)
  }
  dest := os.Args[1]
  key  := os.Args[2]
  data := os.Args[3]
  port := "1075"
  if len(os.Args) == 5 {
    port = os.Args[4]
  }

  address, err := exec.Command("hostname", "-I").Output()
  checkError(err)
  adres := strings.Fields(string(address))
  addr  := ""
  for i:=0; i < len(adres); i++ {
    if adres[i][0:3] == "172" {
      addr = adres[i]
    }
  }
fmt.Println("My IP: ", addr)

  hashk := sha1hash(key)
fmt.Println("Key", hashk)

  service  := addr + ":12500"
  src, err := net.ResolveUDPAddr("udp", service)
  checkError(err)

  service   = dest + ":" + port
  nod, err := net.ResolveUDPAddr("udp", service)

  var msg Message
  msg.Idx  = 11
  msg.Key  = hashk
  msg.Info = data
  msg.Gen  = 0
  msg.Src  = src
  msg.Dst  = nod

  conn, err := net.DialUDP("udp", nil, nod)

  checkError(err)

  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)

//fmt.Println("Sending: ", buffer[0:])
   _, err = conn.Write(buffer)
  checkError(err)

  localAddr, err := net.ResolveUDPAddr("udp", ":12500")
  checkError(err)
  conn, err = net.ListenUDP("udp", localAddr)
  checkError(err)
  answ := waitForRec(conn)
  conn.Close()

  if answ.Key != "" {
    fmt.Println("Responsible node: ",   answ.Key)
    fmt.Println("Responsible status: ", answ.Status)
    fmt.Println("Responsible IP: ",     answ.Src)
  }
}


func waitForRec(conn *net.UDPConn) *Message {

  var buf [512]byte
  _, err := conn.Read(buf[0:])
  checkError(err)
  dec  := json.NewDecoder(bytes.NewReader([]byte(buf[0:])))
  answ := new (Message)
  answ.Info = nil
  err = dec.Decode(&answ)
  checkError(err)
  return answ
}


func generateNodeId() string {
  u, err := uuid.NewV4()
  if err != nil {
    panic(err)
  }

  return sha1hash(u.String())
}


func sha1hash(str string) string {
  // calculate sha-1 hash
  hasher := sha1.New()
  hasher.Write([]byte(str))

  return fmt.Sprintf("%x", sha1.Sum([]byte(str)))
//  return fmt.Sprintf("%x", hasher.Sum(nil))

}


func checkError(err error) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
    os.Exit(1)
  }
}

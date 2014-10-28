package main

import (
  "fmt"
  "net"
  "encoding/json"
  "os"
  "os/exec"
//  "math/big"
  "bytes"
  "strings"
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


// Call with:  printring <entrypoint node in ring> [<portnum>]
func main () {

  if len(os.Args) < 2 {
    fmt.Println("To few args")
    os.Exit(1)
  }
  rem  := os.Args[1]
  port := "1075"
  if len(os.Args) == 3 {
    port = os.Args[2]
  }

  address, err := exec.Command("hostname", "-I").Output()
  checkError(err, 901)
  adres := strings.Fields(string(address))
  addr  := ""
  for i:=0; i < len(adres); i++ {
    if adres[i][0:3] == "172" {
      addr = adres[i]
    }
  }
fmt.Println("My IP: ", addr)

  service  := addr + ":12500"
  src, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 902)

  service   = rem + ":" + port
  nod, err := net.ResolveUDPAddr("udp", service)
  checkError(err, 903)

  var msg Message
  msg.Idx = 05
  msg.Key = ""
  msg.Src = src
  msg.Dst = nod

  conn, err := net.DialUDP("udp", nil, nod)
  checkError(err, 904)

  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err, 905)

//  fmt.Println("Sending: ", buffer[0:])
   _, err = conn.Write(buffer)
  checkError(err, 906)

  localAddr, err := net.ResolveUDPAddr("udp", ":12500")
  checkError(err, 907)
  conn, err = net.ListenUDP("udp", localAddr)
  checkError(err, 908)
/*
  for {
    answ := waitForRec(conn)
    if answ.Key == "" {
      conn.Close()
      break
    }
    nBigInt := big.Int{}
    nBigInt.SetString(answ.Key, 16)
//    fmt.Printf("%s %s\n", answ.Key, nBigInt.String())
    fmt.Printf("%s ", answ.Key)
    fmt.Println(" ", answ.Src)
  }
  conn.Close()
*/
}


func waitForRec(conn *net.UDPConn) *Message {

  var buf [512]byte
  _, err := conn.Read(buf[0:])
  checkError(err, 909)
  dec  := json.NewDecoder(bytes.NewReader([]byte(buf[0:])))
  answ := new (Message)
  answ.Info = nil
  err = dec.Decode(&answ)
  checkError(err, 910)
  return answ
}


func checkError(err error, i int) {
  if err != nil {
    fmt.Fprintf(os.Stderr, "Fatal error ", i, err.Error())
    os.Exit(1)
  }
}

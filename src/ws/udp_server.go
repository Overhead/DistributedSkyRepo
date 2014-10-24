package main

import (
	"fmt"
	"net"
	"encoding/json"
	"os"
	"bytes"
)

type JoinMsg struct {
	Idx	int
	Key string
	Src *net.UDPAddr
	Dst *net.UDPAddr

}

func handleRequest(len int, buffer [512]byte) {
		dec := json.NewDecoder(bytes.NewReader([]byte(buffer[0:])))
		msg := JoinMsg{}
		err := dec.Decode(&msg)
		checkError(err)
		fmt.Printf("Msg recieved with Idx: %s", msg.Idx)
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Fatal error %s\n", err)
		os.Exit(1)
	}
	
}

func main() {
	var buf [512]byte 
	udpAddr, err := net.ResolveUDPAddr("udp", ":1075")
	checkError(err)
	fmt.Printf("Listen on %s\n", udpAddr)

	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	
	for {
			n,_, err := conn.ReadFromUDP(buf[0:])
			checkError(err)
			go handleRequest(n, buf)
		}

}

package main

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

type Msg struct {
		Action int
		Key string
		NewKey string
		Date int64
}

type DhtNode struct {
	Ip string
	Id string
}

type DhtNodes struct {
	Nodes []*DhtNode
}

func echoHandler(ws *websocket.Conn) {
	msg := make([]byte, 512)
	n, err := ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg[:n])
	var res Msg
	json.Unmarshal([]byte(msg[:n]), &res)

	node1 := DhtNode{Ip: "123", Id: "15123"}
	node2 := DhtNode{Ip: "1512", Id: "qssad"}
	list := []*DhtNode{}

	list  = append(list, &node1)
	list  = append(list, &node2)

	nodes := DhtNodes{Nodes: list }
	response, _ := json.Marshal(nodes)	
	m, err := ws.Write([]byte(string(response)))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", msg[:m])
}

func main() {
	http.Handle("/node", websocket.Handler(echoHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

package main

import (
	"fmt"
	"math/big"
	"net"
	"encoding/json"
  "flag"
	"os"
	"bytes"
	"dht"
)

type Node struct { 
	ID string
  LocalIP string
	LocalPort int
	Successor *RemoteNode
	Fingers map[int] *RemoteNode
}

type RemoteNode struct {
	ID string
  IP string
	Port int
}

type Msg struct {
	Key	string
	Action int
	Src	string
	Dst	string
}

type Transport struct {
	bindAddress string
}

/*
func (transport *Transport) listen() {
	var buf [512]byte 
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress)
	checkError(err)
	fmt.Printf("Start listen\n")
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	for {
		n,_, err := conn.ReadFromUDP(buf[0:])
		fmt.Printf("Retreived msg\n")
		checkError(err)
		go handleRequest(n, buf)
	}
} */

func handleRequest(len int, buffer [512]byte, curNode *Node, transport *Transport) {
		dec := json.NewDecoder(bytes.NewReader([]byte(buffer[0:])))
		msg := Msg{}
		err := dec.Decode(&msg)
		checkError(err)
		udpAddr, err := net.ResolveUDPAddr("udp", msg.Src)
		checkError(err)
		// we got a message
		switch msg.Action {
		case 1: //Join initiate msg
						fmt.Printf("Recieved Join from %s", msg.Src)
						checkError(err)
						rm := RemoteNode{msg.Key, udpAddr.IP.String(),udpAddr.Port}
						curNode.addToRing(&rm)
						reply := Msg{}
						reply.Action = 2
						reply.Key = curNode.ID
						reply.Src = msg.Dst
						reply.Dst = msg.Src
						fmt.Printf("Sending Join-reply to %s", reply.Dst)
						transport.send(&reply)
		case 2: //Join reply
						fmt.Printf("Joined reply from %s with ID-%s\n", msg.Src, msg.Key)
		case 3: //Lookup request
						fmt.Printf("Src %s sent lookup request \n", msg.Src)
						lookupNode := curNode.lookup(msg.Key)
						source := fmt.Sprintf("%s:%d",lookupNode.IP,lookupNode.Port)
						reply := Msg{}
						reply.Key = lookupNode.ID
						reply.Src = source
						reply.Dst = msg.Src
						transport.send(&reply)
		case 4: //Set successor
						fmt.Printf("Src %s sent set successor!\n", msg.Src)
						remoteAddr, err := net.ResolveUDPAddr("udp", msg.Src)
						checkError(err)
						rm := RemoteNode{msg.Key, remoteAddr.IP.String(), remoteAddr.Port}
						curNode.Successor = &rm
						fmt.Printf("This Node-%s has updated it's successor to port-%d", curNode.ID, curNode.Successor.Port)
		}

}

func (transport *Transport) send(msg *Msg) {
	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()
	json, err := json.Marshal(msg)
	checkError(err)
	_, err = conn.Write(json)
}


func makeDHTNode (id *string, localAddress string, localPort int) (*Node) {
	newNode := new(Node)	

	if(id != nil) {
		newNode.ID = *id
	} else {
		newNode.ID = dht.GenerateNodeId()
	}
 
	newNode.Successor = nil
	newNode.Fingers = make(map[int]*RemoteNode)
	newNode.LocalIP = localAddress
	newNode.LocalPort = localPort
	
	
		//fmt.Printf("Node-%s added finger with index %d and value %s\n", newNode.ID, i, hex)
	//fmt.Print("ID: " + newNode.ID + "\n")

	return newNode

}


func (currentNode* Node) addToRing(newNode* RemoteNode) {

	fmt.Printf("Adding Node-%s to ring\n", newNode.ID)
	if currentNode.Successor == nil { //When first node is added to ring
		currentNode.Successor = newNode
		newNode.sendUpdateSuccessor(currentNode.ID,currentNode.LocalIP, currentNode.LocalPort) 
		//newNode.Successor = currentNode
		//currentNode.initFingers()
		//currentNode.updateOthersFinger()
		fmt.Printf("1 Node-%s added successor %s \n", currentNode.ID, newNode.ID)
		return 
	}

	node := currentNode.lookup(newNode.ID)	 //Find node that is responsible for the new node

	if node != nil {
		newNode.sendUpdateSuccessor(node.Successor.ID,node.Successor.IP, node.Successor.Port) //Solve problem with successor
		node.sendUpdateSuccessor(newNode.ID, newNode.IP, newNode.Port)
		//newNode.Successor = node.Successor
		//node.Successor = newNode
		//node.initFingers()
		//node.updateOthersFinger()
		fmt.Printf("2 Node-%s added successor %s \n", node.ID, newNode.ID)	
	} else {
			fmt.Printf("3 Node is nil \n")	
	}
	
}

func (curNode* Node) forwardLookup(id string) *Msg {
 	msg := Msg{}
	msg.Action = 3
	msg.Key = id
  msg.Dst = fmt.Sprintf("%s:%d",curNode.Successor.IP,curNode.Successor.Port)
  msg.Src = fmt.Sprintf("%s:%d",curNode.LocalIP,curNode.LocalPort)

	remoteAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
  checkError(err)
  conn, err := net.DialUDP("udp", nil, remoteAddr)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)


  localAddr, err := net.ResolveUDPAddr("udp", msg.Src)
  checkError(err)
  conn2, err := net.ListenUDP("udp", localAddr)
  checkError(err)
  defer conn2.Close()
  var buf [512]byte
  _, err = conn2.Read(buf[0:])
  checkError(err)
	//fmt.Printf("Received3: (%s) %s \n", nod.localPort, buf[0:n])

  dec  := json.NewDecoder(bytes.NewReader([]byte(buf[0:])))
  answ := new (Msg)
  err = dec.Decode(&answ)
  checkError(err)
  return answ
}

//Sends a msg to the given node to update its successor with the info that are put as parameters
func (node* RemoteNode) sendUpdateSuccessor(key, ip string, port int) {
	msg := Msg{}
	msg.Action = 4
	msg.Key = key
  msg.Dst = fmt.Sprintf("%s:%d",ip,port)
  msg.Src = fmt.Sprintf("%s:%d",node.IP,node.Port)

 	remoteAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	checkError(err)
  conn, err := net.DialUDP("udp", nil, remoteAddr)
  checkError(err)
  defer conn.Close()
  buffer, err := json.Marshal(msg)
  checkError(err)
  _, err = conn.Write(buffer)
  checkError(err)
}

func (curNode* Node) lookup(id string) *RemoteNode {
	if curNode.Successor == nil { //No successor, so this node is responsible
		rm := RemoteNode{curNode.ID, curNode.LocalIP, curNode.LocalPort}
		return &rm
	} else if dht.Between([]byte(curNode.ID), []byte(curNode.Successor.ID), []byte(id)) { //Between this and successor, so this is responsible 
		rm := RemoteNode{curNode.ID, curNode.LocalIP, curNode.LocalPort}
		return &rm
	} else {
		msg := curNode.forwardLookup(id)
		remoteAddr, err := net.ResolveUDPAddr("udp", msg.Src)
		checkError(err)
		rm := RemoteNode{msg.Key, remoteAddr.IP.String(), remoteAddr.Port}
		return  &rm //No finger table, just send request to successor node
	}
/*else if len(curNode.Fingers) != 0 { //The finger table is not empty, so check it
		lastFinger := curNode.Fingers[len(curNode.Fingers)] //Last finger node in the map

		//Node we are looking for is not between current and the last finger, so just send it to last finger directly
		if !dht.Between([]byte(curNode.ID), []byte(lastFinger.ID), []byte(id)) {
			return lastFinger.lookup(id) 
		} else { //Id is between some other finger
			//Loop through all fingers and see if they are between the id we are looking for		
			for _, nextFinger := range curNode.Fingers {
				if dht.Between([]byte(nextFinger.ID), []byte(nextFinger.Successor.ID), []byte(id)) { 					
					return nextFinger.lookup(id) 				
				} else {
					continue
				}
			}			
		}
	}*/
		msg := curNode.forwardLookup(id) //Forward lookup to next node
		remoteAddr, err := net.ResolveUDPAddr("udp", msg.Src) //Get the addr of remote node
		checkError(err)
		rm := RemoteNode{msg.Key, remoteAddr.IP.String(), remoteAddr.Port}
		return  &rm  //Return a remoteNode object
}

func (curNode* Node) printRing(){
	/*
	fmt.Printf("%s \n",curNode.ID) //Print First
	for nextN, thisN := curNode.Successor, curNode ; nextN.ID != thisN.ID; {
		fmt.Printf("%s \n",nextN.ID) //Print second, then loop and print rest until ID is same as first
		nextN = nextN.Successor
	}*/

}

func (curNode* Node) initFingers(){
	/*var nrBits = 3

	for i := 1; i < nrBits+1; i++ {
		hex, _ := dht.CalcFinger([]byte(curNode.ID), i, nrBits)
	
		if hex == "" {
			hex = "00"
		}
	
		curNode.Fingers[i] = curNode.lookup(hex)
		fmt.Printf("Node-%s added finger %d as Node-%s\n",curNode.ID, i, curNode.Fingers[i].ID)
	}*/
}

func (curNode* Node) updateOthersFinger(){
	/*if(curNode.Successor == nil){
		//fmt.Printf("Nothing to update")
		return
	} else {
		for nextN, thisN := curNode.Successor, curNode ; nextN.ID != thisN.ID; { //Loop through all and update their fingers
			nextN.initFingers()
			nextN = nextN.Successor
		}
	}*/
}

func (curNode* Node) testCalcFingers(k int, m int) {
	dht.CalcFinger([]byte(curNode.ID), k, m)
	fmt.Printf("Finger %d for Node-%s is Node-%s\n", k, curNode.ID, curNode.Fingers[k].ID)
}

func (curNode* Node) printFinger(k int, m int) {
	dht.CalcFinger([]byte(curNode.ID), k, m)
	fmt.Printf("Finger %d for Node-%s is Node-%s\n", k, curNode.ID, curNode.Fingers[k].ID)
}


func (curNode* Node) find_distance(b []byte, bits int) *big.Int{

	result := dht.Distance([]byte(curNode.ID), b, bits)
	fmt.Printf("Disance from %s to %s is %d \n", curNode.ID, b, result)
	return result
}

func (curNode* Node) is_between(b, id string) bool{
	return dht.Between([]byte(curNode.ID), []byte(b), []byte(id))
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Fatal error %s\n", err)
		os.Exit(1)
	}
	
}

func main() {
	
 	localAddr := flag.String("localAddress", "localhost", "IP of this node")
	localPort := flag.Int("localPort", 2020, "Port of this node")
	remoteAddr := flag.String("remoteAddress", "localhost", "IP of remote node to join")
	remotePort := flag.Int("remotePort", 0, "Port of remote node")
	flag.Parse()

	portString := fmt.Sprintf("%d", *localPort)
	portString2 := fmt.Sprintf("%d", *remotePort);

	node :=	makeDHTNode(nil, *localAddr, *localPort)

	transport := new(Transport)
	transport.bindAddress = fmt.Sprintf(*localAddr+":"+portString)
  
	var buf [512]byte 
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress)
	checkError(err)
	fmt.Printf("Start listen on %s\n",udpAddr)

	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	
	if *remotePort != 0 { //If node is started and asked to join someone, send a join msg
		msg := Msg{}
		msg.Action = 1
		msg.Key = node.ID
		msg.Src = transport.bindAddress
		msg.Dst = fmt.Sprintf(*remoteAddr+":"+portString2)		
		fmt.Printf("Sending Join msg to %s\n", msg.Dst)
		transport.send(&msg)
	} else {
		fmt.Printf("No node to connect to\n")	
	}

	for {
			n,_, err := conn.ReadFromUDP(buf[0:])
			checkError(err)
			go handleRequest(n, buf, node, transport)
		}


}

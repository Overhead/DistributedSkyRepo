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
	"time"
	"errors"
)

type Node struct { 
	ID string
  LocalIP string
	LocalPort int
	Successor *RemoteNode
	Fingers map[int] *RemoteNode
	//Connection *net.UDPConn
	Transp *Transport
	msgChannel chan *Msg
}

type RemoteNode struct {
	ID string
  IP string
	Port int
	Successor *RemoteNode
}

type Msg struct {
	Key	string
	NextKey string
	Action int
	Src	string
	Dst	string
	NextAddr string
}

type Transport struct {
	bindAddress string
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
						fmt.Printf("Recieved Join from %s with its ID-%s\n", msg.Src,msg.Key)
						checkError(err)
						rm := RemoteNode{msg.Key, udpAddr.IP.String(),udpAddr.Port, nil}
						curNode.addToRing(&rm)
						/*ringMsg := curNode.printRing() //Create a message that will be sent through ring and print all nodes
						transport.send(ringMsg)*/
		case 2: //Join reply
						fmt.Printf("Joined reply from %s with ID-%s\n", msg.Src, msg.Key)
		case 3: //Lookup request
						fmt.Printf("Src %s sent lookup request on Key-%s\n", msg.Src, msg.Key)
						lookupNode := curNode.lookup(msg.Key)
						reply := Msg{}
						reply.Action = 4
						reply.Key = lookupNode.ID
						reply.NextKey = lookupNode.Successor.ID
						reply.Src = fmt.Sprintf("%s:%d",lookupNode.IP,lookupNode.Port) //Set the nodes address to src
						reply.Dst = msg.Src //Send msg back to the one it was retreived from
						reply.NextAddr = fmt.Sprintf("%s:%d",lookupNode.Successor.IP,lookupNode.Successor.Port)
						transport.send(&reply)
		case 4: //Get lookupmsg and put it on channel
						curNode.msgChannel <- &msg
		case 5: //Set successor
						fmt.Printf("Src %s sent set successor!\n", msg.Src)
						remoteAddr, err := net.ResolveUDPAddr("udp", msg.Src)
						checkError(err)
						rm := RemoteNode{msg.Key, remoteAddr.IP.String(), remoteAddr.Port, nil}
						curNode.Successor = &rm
						fmt.Printf("This Node-%d has updated its successor to port-%d\n", curNode.LocalPort, curNode.Successor.Port)
		case 6: //Print node id msg
						fmt.Printf("Ring Node-%s, Src %s and Dst %s\n", curNode.ID, msg.Src, msg.Dst)
						msg.Dst = fmt.Sprintf("%s:%d",curNode.Successor.IP,curNode.Successor.Port)
						if msg.Src != msg.Dst {
							transport.send(&msg)
						}
		case 7: //Got a msg telling me to initiate update-finger table
						fmt.Printf("Node-%s told me to update fingertable", msg.Src)
						curNode.updateOthersFinger()
							
		case 8: //Update own fingertable and forward it around ring
						curNode.initFingers() //Set its own finger-table
						msg.Dst = fmt.Sprintf("%s:%d",curNode.Successor.IP,curNode.Successor.Port)
						if msg.Src != msg.Dst {
							transport.send(&msg) //Tell successor to do the same
						}
		}

}

func makeDHTNode (id string, localAddress string, localPort int) (*Node) {
	newNode := new(Node)	

	if(id != "") {
		newNode.ID = id
	} else {
		newNode.ID = dht.GenerateNodeId()
	}
 
	newNode.Successor = nil
	newNode.Fingers = make(map[int]*RemoteNode)
	newNode.msgChannel = make(chan *Msg)
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
		currentNode.initFingers()
		currentNode.updateOthersFinger()
		fmt.Printf("1 Node-%d added successor %d \n", currentNode.LocalPort, newNode.Port)
		return 
	}

	node := currentNode.lookup(newNode.ID)	 //Find node that is responsible for the new node

	if node != nil {
		newNode.sendUpdateSuccessor(node.Successor.ID,node.Successor.IP, node.Successor.Port) //Solve problem with successor
		node.sendUpdateSuccessor(newNode.ID, newNode.IP, newNode.Port)
		//node.initFingers()
	 //node.updateOthersFinger()
		currentNode.sendUpdateFingerTable(node.IP, node.Port)
		//fmt.Printf("2 Node-%d added successor %d \n", node.Port, newNode.Port)
	} else {
			fmt.Printf("3 Node is nil \n")	
	}
	
}

func (curNode* Node) lookup(id string) *RemoteNode {
	if curNode.Successor == nil { //No successor, so this node is responsible
		suc := RemoteNode{curNode.Successor.ID, curNode.Successor.IP, curNode.Successor.Port, nil}
		rm := RemoteNode{curNode.ID, curNode.LocalIP, curNode.LocalPort, &suc}
		fmt.Printf("Nil suc\n")
		return &rm
	} else if dht.Between([]byte(curNode.ID), []byte(curNode.Successor.ID), []byte(id)) { //Between this and successor, so this is responsible 
		suc := RemoteNode{curNode.Successor.ID, curNode.Successor.IP, curNode.Successor.Port, nil}
		rm := RemoteNode{curNode.ID, curNode.LocalIP, curNode.LocalPort, &suc}
		fmt.Printf("Between this and successor\n")
		return &rm
	} else if len(curNode.Fingers) != 0 { //The finger table is not empty, so check it
		lastFinger := curNode.Fingers[len(curNode.Fingers)] //Last finger node in the map

		//Node we are looking for is not between current and the last finger, so just send it to last finger directly
		if !dht.Between([]byte(curNode.ID), []byte(lastFinger.ID), []byte(id)) {
			fmt.Printf("NOT between this and last finger\n")
			rm, err := curNode.forwardLookup(id, lastFinger.IP, lastFinger.Port)
			if err != nil {
				fmt.Printf("Error on finger-lookup %s", err)
				return nil
			} else {
					return rm //default return
			}
		} else { //Id is between some other finger
			//Loop through all fingers and see if they are between the id we are looking for		
			for _, nextFinger := range curNode.Fingers {
				if dht.Between([]byte(nextFinger.ID), []byte(nextFinger.Successor.ID), []byte(id)) { 		//If node we are looking for are between current and it's successor, send it there			
					fmt.Printf("Send to finger %s\n", nextFinger.ID)					
					rm, err := curNode.forwardLookup(id, nextFinger.IP, nextFinger.Port)
					if err != nil {
						fmt.Printf("Error on finger-lookup %s", err)
						return nil
					} else {
							return rm //Return the lookup node
					}			
				} else {
					continue //Not between, then continue loop
				}
			}			
		}
	} else { //Empty finger-table, so send request to successor
		fmt.Printf("Empty finger table, send to successor\n")
		rm, err := curNode.forwardLookup(id, curNode.Successor.IP, curNode.Successor.Port)
		if err != nil {
			fmt.Printf("Error on lookup %s", err)
		} else {
				return  rm
		}
	}
		fmt.Printf("default lookup, send to successor\n")
		rm, err := curNode.forwardLookup(id, curNode.Successor.IP, curNode.Successor.Port) //Default lookup on successor
		if err != nil {
			fmt.Printf("Error on lookup %s", err)
			return nil
		} else {
				return rm
		}
}

func (curNode* Node) forwardLookup(id, remoteIp string, remotePort int) (*RemoteNode, error) {
 	msg := Msg{}
	msg.Action = 3
	msg.Key = id
  msg.Dst = fmt.Sprintf("%s:%d",remoteIp, remotePort)
  msg.Src = fmt.Sprintf("%s:%d",curNode.LocalIP,curNode.LocalPort)

	curNode.Transp.send(&msg)

	answ, errMsg := curNode.getMsg(5)
	
	if answ != nil {
		remoteAddr, err := net.ResolveUDPAddr("udp", answ.Src) //Get the addr of remote node
		checkError(err)
		sucAddr, err := net.ResolveUDPAddr("udp", answ.NextAddr)
		checkError(err)
		suc := RemoteNode{answ.NextKey, sucAddr.IP.String(), sucAddr.Port, nil}
		rm := RemoteNode{answ.Key, remoteAddr.IP.String(), remoteAddr.Port, &suc}
		return  &rm, nil  
	} else {
		return nil, errMsg
	}
}

func (curNode *Node) sendUpdateFingerTable(ip string, port int) {
	msg := Msg{}
	msg.Action = 7
	msg.Key = curNode.ID
  msg.Dst = fmt.Sprintf("%s:%d", ip, port)
  msg.Src = fmt.Sprintf("%s:%d", curNode.LocalIP, curNode.LocalPort)
	curNode.Transp.send(&msg)
}

//Sends a msg to the given node to update its successor with the info that are put as parameters
func (node* RemoteNode) sendUpdateSuccessor(key, ip string, port int) {
	msg := Msg{}
	msg.Action = 5
	msg.Key = key
  msg.Dst = fmt.Sprintf("%s:%d",node.IP,node.Port)
  msg.Src = fmt.Sprintf("%s:%d",ip,port)

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

func (curNode* Node) printRing() *Msg{
	msg := Msg{}
	msg.Action = 6
	msg.Key = curNode.ID
  msg.Dst = fmt.Sprintf("%s:%d",curNode.LocalIP,curNode.LocalPort)
  msg.Src = fmt.Sprintf("%s:%d",curNode.LocalIP,curNode.LocalPort)
	return &msg
/*
	fmt.Printf("%s \n",curNode.ID) //Print First
	for nextN, thisN := curNode.Successor, curNode ; nextN.ID != thisN.ID; {
		fmt.Printf("%s \n",nextN.ID) //Print second, then loop and print rest until ID is same as first
		nextN = nextN.Successor
	}*/
}

func (curNode* Node) initFingers(){
	var nrBits = 3

	for i := 1; i < nrBits+1; i++ {
		hex, _ := dht.CalcFinger([]byte(curNode.ID), i, nrBits)
	
		if hex == "" {
			hex = "00"
		}
	
		curNode.Fingers[i] = curNode.lookup(hex)
		fmt.Printf("Node-%s added finger %d as Node-%s\n",curNode.ID, i, curNode.Fingers[i].ID)
	}
}

func (curNode* Node) updateOthersFinger(){
	msg := Msg{}
	msg.Key = curNode.ID
	msg.Action = 8
	msg.Src = fmt.Sprintf("%s:%d", curNode.LocalIP, curNode.LocalPort)
	msg.Dst = fmt.Sprintf("%s:%d", curNode.LocalIP, curNode.LocalPort)
	curNode.Transp.send(&msg)
}

func (curNode* Node) getMsg(n int) (*Msg, error) {
	if n == 0 {
		return nil, errors.New("No messages retrieved")
	} else {
		msg := <- curNode.msgChannel		
		if msg == nil {
				time.Sleep(1 * 1e9)
				return curNode.getMsg(n-1)
		}	else {
			return msg, nil
		}			
	}
	return nil, errors.New("No messages retrieved")
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
	
	id := flag.String("ID", "", "ID on node")
 	localAddr := flag.String("localAddress", "127.0.0.1", "IP of this node")
	localPort := flag.Int("localPort", 2020, "Port of this node")
	remoteAddr := flag.String("remoteAddress", "127.0.0.1", "IP of remote node to join")
	remotePort := flag.Int("remotePort", 0, "Port of remote node")
	flag.Parse()

	portString := fmt.Sprintf("%d", *localPort)
	portString2 := fmt.Sprintf("%d", *remotePort);

	node :=	makeDHTNode(*id, *localAddr, *localPort)
	transport := new(Transport)
	transport.bindAddress = fmt.Sprintf(*localAddr+":"+portString)
	node.Transp = transport  

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

	<- node.msgChannel
}

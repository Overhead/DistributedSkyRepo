package lab1

import (
	"fmt"
	"math/big"
)

type Node struct { 
	ID string
  LocalIP string
	LocalPort string
	Successor *Node
	Fingers map[int]string
}

func makeDHTNode (id *string, localAddress string, localPort string) (*Node) {
	var nrBits = 3

	newNode := new(Node)	

	if(id != nil) {
		newNode.ID = *id
	} else {
		newNode.ID = generateNodeId()
	}
 
	newNode.Successor = nil
	newNode.Fingers = make(map[int]string)
	newNode.LocalIP = localAddress
	newNode.LocalPort = localPort
	
	for i := 1; i < nrBits+1; i++ {
		hex, _ := calcFinger([]byte(newNode.ID), i, nrBits)
		
		if hex == "" {
			hex = "00"
		}
		
		newNode.Fingers[i] = hex
		//fmt.Printf("Node-%s added finger with index %d and value %s\n", newNode.ID, i, hex)
	}
	//fmt.Print("ID: " + newNode.ID + "\n")

	return newNode

}


func (currentNode* Node) addToRing(newNode* Node) {

	if currentNode.Successor == nil { //When first node is added to ring
		currentNode.Successor = newNode
		newNode.Successor = currentNode
		//fmt.Printf("1 Node-%s added successor %s \n", currentNode.ID, newNode.ID)
		return 
	}

	node := currentNode.lookup(newNode.ID)	 //Find node that is responsible for the new node

	if node != nil {
		newNode.Successor = node.Successor
		node.Successor = newNode
		//fmt.Printf("2 Node-%s added successor %s \n", node.ID, newNode.ID)	
	} else {
			fmt.Printf("3 Node is nil \n")	
	}
}

func (currentNode* Node) lookup(id string) *Node{

	//Task2 Objective 1
	/*if currentNode.Successor == nil {
		return currentNode	
	} else if between([]byte(currentNode.ID), []byte(currentNode.Successor.ID), []byte(id)) {
		return currentNode
	} else {
		distToID := distance([]byte(curNode.ID), []byte(id), 3)
	 	dist3 := distance([]byte(curNode.ID), []byte(curNode.Fingers[3]), 3)
		if distToID >= dist3 { //Distance is greater or equal than furthest finger node, so just send it directly there
			return currentNode.Fingers[3].lookup(id) 
		} else if between([]byte(currentNode.Fingers[2].ID), []byte(currentNode.Fingers[3].ID), []byte(id)) { //Id is finger2 or between 2 and 3, so send to 2
			return currentNode.Fingers[2].lookup(id) 
		} else {
			return currentNode.Fingers[1].lookup(id)	//Just send to 1 if nothing else goes through
		}
	}*/


	//Task 1 Objective 1
	if(currentNode.Successor == nil){
		//fmt.Printf("Node-%s is responsible for %s \n", currentNode.ID, id)
		return currentNode
	} else if between([]byte(currentNode.ID), []byte(currentNode.Successor.ID), []byte(id)) {
		//fmt.Printf("Node-%s is responsible for %s \n", currentNode.ID, id)
		return currentNode
	} else 
	{
		//fmt.Printf("Node-%s is NOT responsible for %s \n", currentNode.ID, id)
		return currentNode.Successor.lookup(id) //Do same method, just for successor node
	}
	return nil
}

func (curNode* Node) printRing(){
	
	fmt.Printf("%s \n",curNode.ID) //Print First
	for nextN, thisN := curNode.Successor, curNode ; nextN.ID != thisN.ID; {
		fmt.Printf("%s \n",nextN.ID) //Print second, then loop and print rest until ID is same as first
		nextN = nextN.Successor
	}

}

func (curNode* Node) find_distance(b []byte, bits int) *big.Int{

	result := distance([]byte(curNode.ID), b, bits)
	fmt.Printf("Disance from %s to %s is %d \n", curNode.ID, b, result)
	return result
}


func (curNode* Node) testCalcFingers(k int, m int) {
	calcFinger([]byte(curNode.ID), k, m)
}

func (curNode* Node) printFinger(k int, m int) {
	calcFinger([]byte(curNode.ID), k, m)
}

/*
func (curNode* Node) find_successor(id string) *Node{
	node := curNode.find_predecessor(id)
	return node.Successor
}

func (curNode* Node) find_predecessor(id string) *Node{
	node := curNode

	if !between([]byte(node.ID), []byte(node.Successor.ID), []byte(id)) {
			node = node.closest_preceding_finger(id)
	}

	return node
}

func (curNode* Node) closest_preceding_finger(id string) *Node{
	for key, fingerNode := range curNode.Fingers {
		fmt.Println("Key:", key, "Value:", fingerNode)
		if between([]byte(curNode.ID), []byte(id), []byte(fingerNode.ID)) {
			return fingerNode
		}
	}
	return curNode
}*/

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
	Fingers map[int]*Node
}

func makeDHTNode (id *string, localAddress string, localPort string) (*Node) {
	newNode := new(Node)	

	if(id != nil) {
		newNode.ID = *id
	} else {
		newNode.ID = generateNodeId()
	}
 
	newNode.Successor = nil
	newNode.Fingers = make(map[int]*Node)
	newNode.LocalIP = localAddress
	newNode.LocalPort = localPort
	
	
		//fmt.Printf("Node-%s added finger with index %d and value %s\n", newNode.ID, i, hex)
	//fmt.Print("ID: " + newNode.ID + "\n")

	return newNode

}


func (currentNode* Node) addToRing(newNode* Node) {

	if currentNode.Successor == nil { //When first node is added to ring
		currentNode.Successor = newNode
		newNode.Successor = currentNode
		currentNode.initFingers()
		currentNode.updateOthersFinger()
		//fmt.Printf("1 Node-%s added successor %s \n", currentNode.ID, newNode.ID)
		return 
	}

	node := currentNode.lookup(newNode.ID)	 //Find node that is responsible for the new node

	if node != nil {
		newNode.Successor = node.Successor
		node.Successor = newNode
		currentNode.initFingers()
		currentNode.updateOthersFinger()
		//fmt.Printf("2 Node-%s added successor %s \n", node.ID, newNode.ID)	
	} else {
			fmt.Printf("3 Node is nil \n")	
	}
	
}

func (currentNode* Node) lookup(id string) *Node{
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


	//Task2 Objective 1
/*
	if currentNode.Successor == nil {
		return currentNode	
	} else if between([]byte(currentNode.ID), []byte(currentNode.Successor.ID), []byte(id)) {
		return currentNode
	} else {
		distToID := distance([]byte(curNode.ID), []byte(id), 3) //Distance to the ID we are looking for
		lastFingerIndex := curNode.Fingers[len(curNode.Fingers]) //Index of the last finger
	 	dist3 := distance([]byte(curNode.ID), []byte(curNode.Fingers[lastFingerIndex), 3) //Find idstance between currentNode and last finger node
		if distToID >= dist3 { //Distance is greater or equal than furthest finger node, so just send it directly there
			return currentNode.Fingers[].lookup(id) 
		} else { //Id is between some other finger

			for i, nextFinger, := 1, currentNode.Fingers[i] ; i < (len(currentNode.Fingers)+1); i++) //Plan is to loop through all fingers and see if they are between the id we are looking for
				if(!between([]byte(currentNode.ID), []byte(nextFinger.ID), []byte(id)) { //If not, just increment and try nextone
					nextFinger = currentNode.Fingers[i]		
				} else {	
					return currentNode.Fingers[i].lookup(id) 
				}	
		} 
	}*/
}

func (curNode* Node) printRing(){
	
	fmt.Printf("%s \n",curNode.ID) //Print First
	for nextN, thisN := curNode.Successor, curNode ; nextN.ID != thisN.ID; {
		fmt.Printf("%s \n",nextN.ID) //Print second, then loop and print rest until ID is same as first
		nextN = nextN.Successor
	}

}

func (curNode* Node) initFingers(){
	var nrBits = 3

	for i := 1; i < nrBits+1; i++ {
		hex, _ := calcFinger([]byte(curNode.ID), i, nrBits)
	
		if hex == "" {
			hex = "00"
		}
	
		curNode.Fingers[i] = curNode.lookup(hex)
		fmt.Printf("Node-%s added finger %d as Node-%s\n",curNode.ID, i, curNode.Fingers[i].ID)
	}
}

func (curNode* Node) updateOthersFinger(){
	if(curNode.Successor == nil){
		//fmt.Printf("Nothing to update")
		return
	} else {
		for nextN, thisN := curNode.Successor, curNode ; nextN.ID != thisN.ID; { //Loop through all and update their fingers
			nextN.initFingers()
			nextN = nextN.Successor
		}
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
	fmt.Printf("Finger %d for Node-%s is Node-%s\n", k, curNode.ID, curNode.Fingers[k].ID)
	
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

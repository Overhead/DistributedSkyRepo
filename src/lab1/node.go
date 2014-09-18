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

	fmt.Printf("Adding Node-%s to ring\n", newNode.ID)
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
		node.initFingers()
		node.updateOthersFinger()
		//fmt.Printf("2 Node-%s added successor %s \n", node.ID, newNode.ID)	
	} else {
			fmt.Printf("3 Node is nil \n")	
	}
	
}

func (curNode* Node) lookup(id string) *Node{
	//Task 1 Objective 1
	/*if(curNode.Successor == nil){
		//fmt.Printf("Node-%s is responsible for %s \n", curNode.ID, id)
		return curNode
	} else if between([]byte(curNode.ID), []byte(curNode.Successor.ID), []byte(id)) {
		//fmt.Printf("Node-%s is responsible for %s \n", curNode.ID, id)
		return curNode
	} else 
	{
		//fmt.Printf("Node-%s is NOT responsible for %s \n", curNode.ID, id)
		return curNode.Successor.lookup(id) //Do same method, just for successor node
	}
	return nil*/

	fmt.Printf("Current Node-%s\n", curNode.ID)
	//Task2 Objective 1
	if curNode.Successor == nil {
  	fmt.Printf("Successor is nil\n")
		return curNode	
	} else if between([]byte(curNode.ID), []byte(curNode.Successor.ID), []byte(id)) {
  	fmt.Printf("Id is between\n")
		return curNode
	} else if len(curNode.Fingers) != 0 {
		distToID := distance([]byte(curNode.ID), []byte(id), 3) //Distance to the ID we are looking for
		lastFinger := curNode.Fingers[len(curNode.Fingers)] //Last finger node
	 	dist3 := distance([]byte(curNode.ID), []byte(lastFinger.ID), 3) //Find distance between currentNode and last finger node

		//Distance is greater or equal than furthest finger node, so just send it directly there
		if lastFinger.ID != curNode.ID && distToID.Cmp(dist3) == 1 || distToID.Cmp(dist3) == 0 {
			fmt.Printf("Cur-%s and last-%s\n", curNode.ID, lastFinger.ID)
			return lastFinger.lookup(id) 
		} else { //Id is between some other finger
			//Plan is to loop through all fingers and see if they are between the id we are looking for
			for i := 1 ; i < (len(curNode.Fingers)+1); i++ { 
				nextFinger := curNode.Fingers[i]
				//If not, continue
				if !between([]byte(nextFinger.ID), []byte(nextFinger.Successor.ID), []byte(id)) { 					
					continue		
				} else {
					return nextFinger.lookup(id) 		
				}
			}
		}
	} else {
		return curNode.Successor.lookup(id)
	}
	return curNode.Successor.lookup(id)
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

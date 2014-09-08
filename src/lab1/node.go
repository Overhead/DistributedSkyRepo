package lab1

import (
	"fmt"
)

type Node struct { 
	ID string
  LocalIP string
	LocalPort string
	Successor *Node
}

func makeDHTNode (id *string, localAddress string, localPort string) (*Node) {
	newNode := new(Node)	

	if(id != nil) {
		newNode.ID = *id
	} else {
		newNode.ID = generateNodeId()
	}
 
	newNode.LocalIP = localAddress
	newNode.LocalPort = localPort
	//fmt.Print("ID: " + newNode.ID + "\n")

	return newNode

}


func (currentNode* Node) addToRing(newNode* Node) {
	if(currentNode.Successor == nil){
		currentNode.Successor = newNode	
	} else if between([]byte(currentNode.ID), []byte(currentNode.Successor.ID), []byte(newNode.ID)) {
		currentNode.Successor = newNode	
	} else {
		currentNode.Successor.addToRing(newNode)	
	}
}

func (currentNode* Node) lookup(id string) {
	var result = ""	

	if(currentNode.Successor == nil){
		result += currentNode.ID	
	} else if between([]byte(currentNode.ID), []byte(currentNode.Successor.ID), []byte(id)) {
			if(currentNode.ID == id){
				result += currentNode.ID			
			} else if(currentNode.Successor.ID == id){
				result += currentNode.Successor.ID			
			}
	} else 
	{
		currentNode.Successor.lookup(id)
	}

	fmt.Printf("Node-%s is responsible for %s \n", result, id)
}


func (currentNode* Node) find_successor(id string) string{
	return ""
}

func (currentNode* Node) find_predecessor(id string) string{
	return ""
}

func (currentNode* Node) closest_preceding_finger(id string) string{
	return ""
}

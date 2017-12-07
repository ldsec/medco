package loader

import "github.com/dedis/onet/log"

type Concept interface {

}

// Node defines a 'node' in the concept tree
type Node struct {
	Concept string
	Parent *Node
	Cpt  Concept
	Children []*Node
}

// AddNode adds a node to the tree (and every other node in the path that does not exist) and return a pointer to the root
func AddNode(parent *Node, nodeToAdd *Node, path []string){

	// we have reached the end
	if len(path) == 1 {
		// check if the node already exists from a previous execution
		for _,child := range parent.Children {
			if child.Concept == path[0] {
				// if it exist simply update it
				child.Concept = nodeToAdd.Concept
				return
			}
		}
		// if it does not exist create it
		parent.Children = append(parent.Children,nodeToAdd)
		nodeToAdd.Parent = parent
		return
	}

	nextConcept := path[0]
	path = path[1:]

	// it's the root
	if nextConcept == parent.Concept {
		AddNode(parent, nodeToAdd, path)
	} else {
		// check if one of the children
		for _,child := range parent.Children {
			if child.Concept == nextConcept {
				AddNode(child,nodeToAdd,path)
				return
			}
		}
		// None of the children so create a new Node
		newNode := &Node{
			Concept: nextConcept,
			Parent:  parent,
		}
		parent.Children = append(parent.Children, newNode)
		AddNode(newNode, nodeToAdd, path)
	}
}

func PrintTree(root *Node){
	log.LLvl1(root.Concept)
	for _,child := range root.Children{
		PrintTree(child)
	}

}






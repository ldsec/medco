package loader_test

import (
	"testing"
	"github.com/lca1/medco-loader/loader"
	"strings"
	"github.com/dedis/onet/log"
)

func TestAddNode(t *testing.T) {
	root := &loader.Node{
		Concept: "SHRINE",
	}

	nodeToAdd :=  &loader.Node{
		Concept: "Demographics",
	}

	path := "\\SHRINE\\Demographics\\"
	path = path[1:len(path)-1]

	loader.AddNode(root, nodeToAdd, strings.Split(path,"\\"))
	log.LLvl1("###---INSERTED: \\SHRINE\\Demographics\\---#")
	loader.PrintTree(root)

	nodeToAdd =  &loader.Node{
		Concept: "Gender",
	}

	path = "\\SHRINE\\Demographics\\Gender\\"
	path = path[1:len(path)-1]

	loader.AddNode(root, nodeToAdd, strings.Split(path,"\\"))
	log.LLvl1("###---INSERTED: \\SHRINE\\Demographics\\Gender\\---#")
	loader.PrintTree(root)
}

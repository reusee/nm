package nm

import (
	"log"
	"os"
	"strings"
	"testing"
)

var testNode *Node

func TestMain(m *testing.M) {
	f, err := os.Open("qq.html")
	if err != nil {
		log.Fatal(err)
	}
	testNode, err = Parse(f)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestSequence(t *testing.T) {
	program := Program([]Inst{
		{Predict, func(node *Node) bool {
			return node.Tag == "html"
		}, 0, 0},
		{Predict, func(node *Node) bool {
			return node.Tag == "body"
		}, 0, 0},
		{Predict, func(node *Node) bool {
			return node.Tag == "div"
		}, 0, 0},
		{Ok, nil, 0, 0},
	})
	for _, node := range testNode.Children {
		res := Match(node, program)
		for _, n := range res {
			if strings.Join(n.TagPath(), "|") != "html|body|div" {
				t.Fatal("match")
			}
		}
	}
}

func TestOr(t *testing.T) {
	program := Program([]Inst{
		{Predict, func(node *Node) bool {
			return node.Tag == "html"
		}, 0, 0},
		{Split, nil, 2, 4},
		{Predict, func(node *Node) bool {
			return node.Tag == "head"
		}, 0, 0},
		{Jump, nil, 5, 0},
		{Predict, func(node *Node) bool {
			return node.Tag == "body"
		}, 0, 0},
		{Ok, nil, 0, 0},
	})
	for _, node := range testNode.Children {
		res := Match(node, program)
		if len(res) != 2 {
			t.Fatal("match")
		}
	}
}

func TestZeroOrOne(t *testing.T) {
	program := Program([]Inst{
		{Predict, func(node *Node) bool {
			return node.Tag == "html"
		}, 0, 0},
		{Split, nil, 2, 3},
		{Predict, func(node *Node) bool {
			return node.Tag == "head"
		}, 0, 0},
		{Ok, nil, 0, 0},
	})
	for _, node := range testNode.Children {
		res := Match(node, program)
		if len(res) != 2 {
			t.Fatal("match")
		}
	}
}

func TestZeroOrMany(t *testing.T) {
	program := Program([]Inst{
		{Predict, func(node *Node) bool {
			return node.Tag == "html"
		}, 0, 0},
		{Split, nil, 2, 4},
		{Predict, func(node *Node) bool {
			return true
		}, 0, 0},
		{Jump, nil, 1, 0},
		{Ok, nil, 0, 0},
	})
	for _, node := range testNode.Children {
		res := Match(node, program)
		if len(res) != 3 {
			t.Fatal("match")
		}
	}
}

func TestOneOrMany(t *testing.T) {
	program := Program([]Inst{
		{Predict, func(node *Node) bool {
			return node.Tag == "html"
		}, 0, 0},
		{Predict, func(node *Node) bool {
			return node.Tag == "body"
		}, 0, 0},
		{Predict, func(node *Node) bool {
			return node.Tag == "div"
		}, 0, 0},
		{Split, nil, 2, 4},
		{Ok, nil, 0, 0},
	})
	for _, node := range testNode.Children {
		res := Match(node, program)
		for _, r := range res {
			tagPath := r.TagPath()
			if tagPath[0] != "html" || tagPath[1] != "body" || tagPath[2] != "div" {
				t.Fatal("match")
			}
			for _, tag := range tagPath[3:] {
				if tag != "div" {
					t.Fatal("match")
				}
			}
		}
	}
}

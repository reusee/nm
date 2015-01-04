package nm

import (
	"bytes"
	"fmt"
)

type Node struct {
	Parent   *Node
	Children []*Node

	Tag       string
	Text      string
	TextParts []string
	Attr      map[string]string

	Id    string
	Class []string

	Raw    string
	rawBuf *bytes.Buffer
}

func (n *Node) Compare(right *Node) error {
	genErr := func(msg string, args ...interface{}) error {
		return fmt.Errorf("%s\n---left---\n%s\n---right---\n%s\n------\n",
			fmt.Sprintf(msg, args...), n.Raw, right.Raw)
	}
	if len(n.Children) != len(right.Children) {
		return genErr("number of children")
	}
	for i, l := range n.Children {
		r := right.Children[i]
		err := l.Compare(r)
		if err != nil {
			return err
		}
	}
	if n.Tag != right.Tag {
		return genErr("tag <%s> <%s>", n.Tag, right.Tag)
	}
	if n.Text != right.Text {
		return genErr("text")
	}
	if len(n.TextParts) != len(right.TextParts) {
		return genErr("textparts length")
	}
	for i, l := range n.TextParts {
		r := right.TextParts[i]
		if l != r {
			return genErr("textparts")
		}
	}
	if len(n.Attr) != len(right.Attr) {
		return genErr("number of attr")
	}
	for key, value := range n.Attr {
		rvalue := right.Attr[key]
		if value != rvalue {
			return genErr("attr %s: %s <-> %s", key, value, rvalue)
		}
	}
	if n.Raw != right.Raw {
		return genErr("raw")
	}
	return nil
}

func (n *Node) Index() (ret int) {
	ret = -1
	if n.Parent != nil {
		for i, node := range n.Parent.Children {
			if n == node {
				ret = i
				break
			}
		}
	}
	return
}

func (n *Node) TagPath() []string {
	node := n
	var path []string
	for node.Tag != "ROOT" {
		path = append(path, node.Tag)
		node = node.Parent
	}
	for i := 0; i < len(path)/2; i++ {
		j := len(path) - 1 - i
		path[i], path[j] = path[j], path[i]
	}
	return path
}

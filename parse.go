package nm

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"code.google.com/p/go.net/html"
)

func Parse(r io.Reader) ([]*Node, error) {
	root := &Node{
		Tag:    "ROOT",
		rawBuf: new(bytes.Buffer),
	}
	tokenizer := html.NewTokenizer(r)
	currentNode := root
	writeRaw := func() {
		raw := tokenizer.Raw()
		currentNode.rawBuf.Write(raw)
		node := currentNode
		for node.Parent != nil {
			node = node.Parent
			node.rawBuf.Write(raw)
		}
	}
parse:
	for {
		what := tokenizer.Next()
		switch what {
		case html.ErrorToken:
			break parse
		case html.TextToken:
			text := strings.TrimSpace(string(tokenizer.Text()))
			if len(text) > 0 {
				currentNode.Text += text
				currentNode.TextParts = append(currentNode.TextParts, text)
			}
			writeRaw()
		case html.StartTagToken:
			node := &Node{
				Parent: currentNode,
				rawBuf: new(bytes.Buffer),
				Attr:   make(map[string]string),
			}
			currentNode.Children = append(currentNode.Children, node)
			currentNode = node
			writeRaw()
			name, hasAttr := tokenizer.TagName()
			currentNode.Tag = string(name)
			if hasAttr {
				key, val, more := tokenizer.TagAttr()
				currentNode.Attr[string(key)] = string(val)
				for more {
					key, val, more = tokenizer.TagAttr()
					currentNode.Attr[string(key)] = string(val)
				}
			}
			currentNode.collectIdAndClass()
		case html.EndTagToken:
			name, _ := tokenizer.TagName()
			for string(name) != currentNode.Tag { // skip mismatched tag
				currentNode.Raw = string(currentNode.rawBuf.Bytes())
				currentNode = currentNode.Parent
				if currentNode == nil {
					return nil, fmt.Errorf("start tag not found for end tag %s", name)
				}
			}
			writeRaw()
			currentNode.Raw = string(currentNode.rawBuf.Bytes())
			currentNode = currentNode.Parent
		case html.SelfClosingTagToken:
			node := &Node{
				Parent: currentNode,
				Raw:    string(tokenizer.Raw()),
				Attr:   make(map[string]string),
			}
			name, hasAttr := tokenizer.TagName()
			node.Tag = string(name)
			if hasAttr {
				key, val, more := tokenizer.TagAttr()
				node.Attr[string(key)] = string(val)
				for more {
					key, val, more = tokenizer.TagAttr()
					node.Attr[string(key)] = string(val)
				}
			}
			node.collectIdAndClass()
			currentNode.Children = append(currentNode.Children, node)
			writeRaw()
		case html.CommentToken:
			writeRaw()
		}
	}
	root.Raw = string(root.rawBuf.Bytes())

	return root.Children, nil
}

func (n *Node) collectIdAndClass() {
	// id and class
	n.Id = n.Attr["id"]
	for _, class := range strings.Split(n.Attr["class"], " ") {
		class = strings.TrimSpace(class)
		if len(class) > 0 {
			n.Class = append(n.Class, class)
		}
	}
}

func ParseString(s string) ([]*Node, error) {
	return Parse(strings.NewReader(s))
}

func ParseBytes(bs []byte) ([]*Node, error) {
	return Parse(bytes.NewReader(bs))
}

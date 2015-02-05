package nm

import (
	"fmt"
	"os"

	"github.com/reusee/paza"
)

var set *paza.Set

func init() {
	set = paza.NewSet()
	set.Add("expr", set.OrdChoice(
		set.NamedConcat("or-expr", "expr", set.NamedRune("or-op", '|'), "simple-expr"),
		"simple-expr"))
	set.Add("simple-expr", set.OrdChoice(
		set.NamedConcat("concat-expr", "simple-expr", set.Regex(`\s+`), "basic-expr"),
		"basic-expr"))
	set.Add("basic-expr", set.OrdChoice(
		"star-expr",
		"plus-expr",
		"option-expr",
		"elementary-expr"))
	set.Add("star-expr", set.Concat(
		"elementary-expr", set.NamedRune("star-op", '*')))
	set.Add("plus-expr", set.Concat(
		"elementary-expr", set.NamedRune("plus-op", '+')))
	set.Add("option-expr", set.Concat(
		"elementary-expr", set.NamedRune("option-op", '?')))
	set.Add("elementary-expr", set.OrdChoice(
		set.NamedConcat("group-expr", set.NamedRegex("left-paren", `\(`),
			"expr", set.NamedRegex("right-paren", `\)`)),
		"predict",
	))

	set.Add("predict", set.OrdChoice(
		set.NamedConcat("predict-with-attr",
			set.NamedRepeat("basic-predicts", 1, -1, "basic-predict"),
			set.NamedRepeat("option-attr-predict", 0, 1, "attr-predict")),
		"attr-predict"))
	set.Add("basic-predict", set.OrdChoice(
		"id-predict", "class-predict", "tag-predict"))
	set.Add("attr-predict", set.Concat(
		set.Regex(`\[`),
		set.NamedRepeat("option-attr-expr", 0, 1, "attr-expr"),
		set.Regex(`\]`)))
	set.Add("identifier", set.Regex(`[a-zA-Z0-9-_]+`))
	set.Add("id-predict", set.Concat(set.Rune('#'), "identifier"))
	set.Add("class-predict", set.Concat(set.Rune('.'), "identifier"))
	set.Add("tag-predict", set.Concat("identifier"))

	set.Add("attr-expr", set.OrdChoice(
		set.NamedConcat("attr-or-expr", "attr-expr", set.NamedRegex("attr-or-op", "||"), "attr-simple-expr"),
		"attr-simple-expr"))
	set.Add("attr-simple-expr", set.OrdChoice(
		set.NamedConcat("attr-and-expr", "attr-simple-expr", set.NamedRegex("attr-and-op", "&&"), "attr-basic-expr"),
		"attr-basic-expr"))
	set.Add("attr-basic-expr", set.OrdChoice(
		set.NamedConcat("attr-group-expr", set.NamedRegex("attr-left-paren", `\(`),
			"attr-expr", set.NamedRegex("attr-right-paren", `\)`)),
		"attr-elementary-expr"))
	set.Add("attr-elementary-expr", set.Concat(
		"identifier",
		set.NamedRegex("attr-op", `=|!=|~=`),
		set.NamedOrdChoice("value",
			"single-quoted",
			"double-quoted",
			"back-quoted",
			"text")))
}

func Compile(code string) (Program, error) {
	input := paza.NewInput([]byte(code))
	ok, l, node := set.Call("expr", input, 0)
	if !ok {
		return nil, fmt.Errorf("invalid expression")
	}
	if l != len(code) {
		node.Dump(os.Stdout, input)
		return nil, fmt.Errorf("invalid expression")
	}
	node = simplify(node)
	//node.Dump(os.Stdout, input)
	ast := genAst(node, input)
	//ast.dump(0)
	return genProgram(ast)
}

func simplify(node *paza.Node) (ret *paza.Node) {
	for len(node.Subs) == 1 {
		sub := node.Subs[0]
		if node.Start == sub.Start && node.Len == sub.Len {
			node = sub
		}
	}
	for i, sub := range node.Subs {
		node.Subs[i] = simplify(sub)
	}
	return node
}

type Ast struct {
	Op      astOp
	Left    *Ast
	Right   *Ast
	Predict func(*Node) bool
}

type astOp int

const (
	opConcat astOp = iota
	opPredict
	opStar
)

func genAst(node *paza.Node, input *paza.Input) *Ast {
	switch node.Name {
	case "concat-expr":
		return &Ast{
			Op:    opConcat,
			Left:  genAst(node.Subs[0], input),
			Right: genAst(node.Subs[2], input),
		}
	case "predict-with-attr":
		p1 := genPredict(node.Subs[0], input)
		p2 := genPredict(node.Subs[1], input)
		return &Ast{
			Op: opPredict,
			Predict: func(node *Node) bool {
				return p1(node) && p2(node)
			},
		}
	case "star-expr":
		return &Ast{
			Op:   opStar,
			Left: genAst(node.Subs[0], input),
		}
	case "attr-predict":
		p := genPredict(node.Subs[1], input)
		return &Ast{
			Op:      opPredict,
			Predict: p,
		}
	default:
		panic("not handle parse node " + node.Name)
	}
	return nil
}

func truePredict(node *Node) bool {
	return true
}

func genPredict(node *paza.Node, input *paza.Input) func(node *Node) bool {
	switch node.Name {
	case "identifier": // tag
		tag := string(input.Text[node.Start : node.Start+node.Len])
		return func(n *Node) bool {
			return n.Tag == tag
		}
	case "option-attr-predict", "option-attr-expr":
		if len(node.Subs) > 0 {
			return genPredict(node.Subs[0], input)
		} else {
			return truePredict
		}
	case "basic-predicts":
		var predicts []func(*Node) bool
		for _, sub := range node.Subs {
			predicts = append(predicts, genPredict(sub, input))
		}
		return func(n *Node) bool {
			for _, predict := range predicts {
				if !predict(n) {
					return false
				}
			}
			return true
		}
	case "id-predict":
		id := string(input.Text[node.Start+1 : node.Start+node.Len])
		return func(n *Node) bool {
			return n.Id == id
		}
	case "class-predict":
		class := string(input.Text[node.Start+1 : node.Start+node.Len])
		return func(n *Node) bool {
			for _, cls := range n.Class {
				if cls == class {
					return true
				}
			}
			return false
		}
	default:
		panic("not handle predict node " + node.Name)
	}
	return nil
}

func genProgram(ast *Ast) (Program, error) {
	return nil, nil
}

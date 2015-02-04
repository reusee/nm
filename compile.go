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
			set.Repeat(1, -1, "basic-predict"), set.Repeat(0, 1, "attr-predict")),
		"attr-predict"))
	set.Add("basic-predict", set.OrdChoice(
		"id-predict", "class-predict", "tag-predict"))
	set.Add("attr-predict", set.Concat(
		set.Regex(`\[`), set.Repeat(0, 1, "attr-expr"), set.Regex(`\]`)))
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
		pt("%d\n", l)
		node.Dump(os.Stdout, input)
		return nil, fmt.Errorf("invalid expression")
	}
	node.Dump(os.Stdout, input)
	return nil, nil
}

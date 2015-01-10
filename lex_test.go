package nm

import "testing"

var lexTestCases = []struct {
	str   string
	items []item
}{
	{"html",
		[]item{item{what: itemTag, text: []rune("html")}},
	},
	{" html",
		[]item{item{what: itemTag, text: []rune("html")}},
	},
	{" \thtml ",
		[]item{item{what: itemTag, text: []rune("html")}},
	},
	{"html head", []item{
		item{what: itemTag, text: []rune("html")},
		item{what: itemTag, text: []rune("head")},
	}},
	{"#id1 #id2", []item{
		item{what: itemId, text: []rune("id1")},
		item{what: itemId, text: []rune("id2")},
	}},
	{"#id1#id2", []item{
		item{what: itemId, text: []rune("id1")},
		item{what: itemId, text: []rune("id2")},
	}},
	{".class1.class2", []item{
		item{what: itemClass, text: []rune("class1")},
		item{what: itemClass, text: []rune("class2")},
	}},
	{"tag#id.class", []item{
		item{what: itemTag, text: []rune("tag")},
		item{what: itemId, text: []rune("id")},
		item{what: itemClass, text: []rune("class")},
	}},
	{"(html)(head)", []item{
		item{what: itemLeftParen},
		item{what: itemTag, text: []rune("html")},
		item{what: itemRightParen},
		item{what: itemLeftParen},
		item{what: itemTag, text: []rune("head")},
		item{what: itemRightParen},
	}},
	{"[", []item{
		item{what: itemLeftBracket},
	}},
	{"[]", []item{
		item{what: itemLeftBracket},
		item{what: itemRightBracket},
	}},
	{"[id", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
	}},
	{"[ id", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
	}},
	{"[ id=", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("=")},
	}},
	{"[id=id2", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("=")},
		item{what: itemString, text: []rune("id2")},
	}},
	{"[id equals id2", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("equals")},
		item{what: itemString, text: []rune("id2")},
	}},
	{"[id equals \"id2\"", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("equals")},
		item{what: itemString, text: []rune("id2")},
	}},
	{"[id equals 'id2'", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("equals")},
		item{what: itemString, text: []rune("id2")},
	}},
	{"[id equals `id2`", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("equals")},
		item{what: itemString, text: []rune("id2")},
	}},
	{"[id exists", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("exists")},
	}},
	{"[id exists,", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("exists")},
	}},
	{"[ id exists, ", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("exists")},
	}},
	{"[ id exists, class exists", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("exists")},
		item{what: itemAttr, text: []rune("class")},
		item{what: itemOp, text: []rune("exists")},
	}},
	{"[ id exists, class exists]", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("exists")},
		item{what: itemAttr, text: []rune("class")},
		item{what: itemOp, text: []rune("exists")},
		item{what: itemRightBracket},
	}},
	{"[ id exists, class exists ]", []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("id")},
		item{what: itemOp, text: []rune("exists")},
		item{what: itemAttr, text: []rune("class")},
		item{what: itemOp, text: []rune("exists")},
		item{what: itemRightBracket},
	}},
	{`[foo="\t\n \r\b\f\\"]`, []item{
		item{what: itemLeftBracket},
		item{what: itemAttr, text: []rune("foo")},
		item{what: itemOp, text: []rune("=")},
		item{what: itemString, text: []rune("\t\n \r\b\f\\")},
		item{what: itemRightBracket},
	}},
}

var lexFailTestCases = []struct {
	str string
	msg string
}{
	{"%", "invalid char at pos 0"},
	{".foo#", "invalid id at pos 5"},
	{"#id.", "invalid class at pos 4"},
	{"[[", "invalid char at pos 1"},
	{"[id]", "invalid char at pos 3"},
	{`[foo='`, "invalid string at pos 6"},
	{"[foo=`", "invalid string at pos 6"},
}

func TestLexer(t *testing.T) {
	for _, c := range lexTestCases {
		items, err := lex(c.str)
		if err != nil {
			t.Fatal(err)
		}
		if !itemsMatch(items, c.items) {
			t.Fatalf("not match: %s, expecting %v, got %v", c.str, c.items, items)
		}
	}
}

func TestLexerFail(t *testing.T) {
	for _, c := range lexFailTestCases {
		_, err := lex(c.str)
		if err == nil {
			t.Fatalf("should fail: %s", c.str)
		}
		if err.Error() != c.msg {
			t.Fatalf("expecting fail message: %s, got %s", c.msg, err)
		}
	}
}

package nm

import (
	"fmt"
	"unicode"
)

type lexState struct {
	text []rune
	n    int

	items []item
}

const (
	eof rune = -1
)

func (s *lexState) cur() rune {
	if s.n >= len(s.text) {
		return eof
	}
	return s.text[s.n]
}

func (s *lexState) append(item item) {
	s.items = append(s.items, item)
}

func (s *lexState) nextWhile(predict func(rune) bool) {
	for {
		cur := s.cur()
		if cur == eof {
			return
		}
		if predict(cur) {
			s.n++
		} else {
			break
		}
	}
}

type lexer func(*lexState) lexer

type item struct {
	what int
	text []rune
}

const (
	itemLeftParen int = iota
	itemRightParen
	itemId
	itemClass
	itemTag
	itemLeftBracket
	itemRightBracket
	itemAttr
	itemOp
	itemString
	itemExprOp
)

type lexError struct {
	msg string
	pos int
}

func (e lexError) Error() string {
	return fmt.Sprintf("%s at pos %d", e.msg, e.pos)
}

func lex(str string) (items []item, err error) {
	s := &lexState{
		text: []rune(str),
	}
	lexer := lexItem
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = lexError{
					msg: r.(string),
					pos: s.n,
				}
			}
		}()
		for lexer != nil {
			lexer = lexer(s)
		}
	}()
	items = s.items
	return
}

func lexItem(s *lexState) lexer {
	cur := s.cur()
	switch cur {
	case '#':
		return lexId
	case '.':
		return lexClass
	case '(':
		s.append(item{what: itemLeftParen})
		s.n++
		return lexItem
	case ')':
		s.append(item{what: itemRightParen})
		s.n++
		return lexItem
	case '[':
		s.append(item{what: itemLeftBracket})
		s.n++
		return lexConstraint
	case '|', '?', '*', '+':
		s.append(item{what: itemExprOp, text: []rune{cur}})
		s.n++
		return lexItem
	case eof:
		return nil
	}
	if unicode.IsSpace(cur) {
		s.n++
		return lexItem
	} else if isTagChar(cur) {
		return lexTag
	}
	panic("invalid char")
}

func isIdChar(r rune) bool {
	if unicode.IsLetter(r) || unicode.IsDigit(r) ||
		r == '_' || r == '-' {
		return true
	}
	return false
}

func lexId(s *lexState) lexer {
	s.n++ // skip #
	start := s.n
	s.nextWhile(isIdChar)
	if s.n > start {
		s.append(item{what: itemId, text: s.text[start:s.n]})
		return lexItem
	}
	panic("invalid id")
}

var isClassChar = isIdChar

func lexClass(s *lexState) lexer {
	s.n++ // skip .
	start := s.n
	s.nextWhile(isClassChar)
	if s.n > start {
		s.append(item{what: itemClass, text: s.text[start:s.n]})
		return lexItem
	}
	panic("invalid class")
}

func isTagChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func lexTag(s *lexState) lexer {
	start := s.n
	s.nextWhile(isTagChar)
	if s.n > start {
		s.append(item{what: itemTag, text: s.text[start:s.n]})
	}
	return lexItem
}

func lexConstraint(s *lexState) lexer {
	cur := s.cur()
	if cur == eof {
		return nil
	}
	switch cur {
	case ']':
		s.append(item{what: itemRightBracket})
		s.n++
		return lexItem
	}
	if isAttrChar(cur) {
		return lexAttr
	} else if unicode.IsSpace(cur) {
		s.n++
		return lexConstraint
	}
	panic("invalid char")
}

var isAttrChar = isIdChar

func lexAttr(s *lexState) lexer {
	start := s.n
	s.nextWhile(isAttrChar)
	if s.n > start {
		s.append(item{what: itemAttr, text: s.text[start:s.n]})
	}
	return lexOp
}

func isPunctOpChar(r rune) bool {
	switch r {
	case '=', '~', '!', '@', '#', '%', '^', '&', '*', '-', '+', '|', '/', '<', '>', '?':
		return true
	}
	return false
}

func isLetterOpChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func lexOp(s *lexState) lexer {
	cur := s.cur()
	if cur == eof {
		return nil
	} else if unicode.IsSpace(cur) {
		s.n++
		return lexOp
	}
	var predict func(rune) bool
	if isPunctOpChar(cur) {
		predict = isPunctOpChar
	} else if isLetterOpChar(cur) {
		predict = isLetterOpChar
	} else {
		panic("invalid char")
	}
	start := s.n
	s.nextWhile(predict)
	if s.n > start {
		s.append(item{what: itemOp, text: s.text[start:s.n]})
	}
	return lexValue
}

func lexValue(s *lexState) lexer {
	cur := s.cur()
	if unicode.IsSpace(cur) {
		s.n++
		return lexValue
	}
	switch cur {
	case eof:
		return nil
	case '"':
		return lexDoubleQuotedString
	case '\'':
		return lexSingleQuotedString
	case '`':
		return lexBackQuotedString
	case ',':
		s.n++
		return lexConstraint
	case ']':
		s.append(item{what: itemRightBracket})
		s.n++
		return lexItem
		//case '{': TODO
		//	return lexList
	}
	return lexUnquotedString
}

func makeStringLexer(delim rune) lexer {
	return func(s *lexState) lexer {
		var str []rune
	end:
		for {
			s.n++
			cur := s.cur()
			switch cur {
			case '\\': // escape
				s.n++
				c := s.cur()
				switch c {
				case 't':
					str = append(str, '\t')
				case 'n':
					str = append(str, '\n')
				case 'r':
					str = append(str, '\r')
				case 'b':
					str = append(str, '\b')
				case 'f':
					str = append(str, '\f')
				default:
					str = append(str, c)
				}
			case delim: // end
				s.n++
				break end
			case eof:
				panic("invalid string")
			default:
				str = append(str, cur)
			}
		}
		s.append(item{what: itemString, text: str})
		return lexConstraint
	}
}

var lexSingleQuotedString, lexDoubleQuotedString lexer

func init() {
	lexSingleQuotedString = makeStringLexer('\'')
	lexDoubleQuotedString = makeStringLexer('"')
}

func lexBackQuotedString(s *lexState) lexer {
	start := s.n + 1
end:
	for {
		s.n++
		cur := s.cur()
		switch cur {
		case eof:
			panic("invalid string")
		case '`':
			break end
		}
	}
	s.append(item{what: itemString, text: s.text[start:s.n]})
	s.n++
	return lexConstraint
}

func isUnquotedStringChar(r rune) bool {
	return !unicode.IsSpace(r) && r != ',' && r != ']'
}

func lexUnquotedString(s *lexState) lexer {
	start := s.n
	s.nextWhile(isUnquotedStringChar)
	s.append(item{what: itemString, text: []rune(s.text[start:s.n])})
	return lexConstraint
}

/*
func lexList(s *lexState) lexer {
	//TODO
	return nil
}
*/

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
	panic("invalid char")
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

package nm

import "fmt"

var (
	pt = fmt.Printf
)

func runesEq(left, right []rune) bool {
	if len(left) != len(right) {
		return false
	}
	for i, l := range left {
		r := right[i]
		if l != r {
			return false
		}
	}
	return true
}

func itemsMatch(left, right []item) bool {
	if len(left) != len(right) {
		return false
	}
	for i, l := range left {
		r := right[i]
		if l.what != r.what || !runesEq(l.text, r.text) {
			return false
		}
	}
	return true
}

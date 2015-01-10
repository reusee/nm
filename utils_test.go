package nm

import "testing"

func TestRunesEq(t *testing.T) {
	if !runesEq([]rune("foobar"), []rune("foobar")) {
		t.Fatal("eq")
	}
	if runesEq([]rune("foobar"), []rune("foobarbaz")) {
		t.Fatal("eq")
	}
	if runesEq([]rune("foobar"), []rune("foobaR")) {
		t.Fatal("eq")
	}
}

func TestItemsMatch(t *testing.T) {
	if !itemsMatch([]item{}, []item{}) {
		t.Fatal("eq")
	}
	if itemsMatch([]item{}, []item{item{}}) {
		t.Fatal("eq")
	}
	if itemsMatch([]item{item{}}, []item{item{what: itemId}}) {
		t.Fatal("eq")
	}
}

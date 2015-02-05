package nm

import (
	"fmt"
	"strings"
)

var (
	pt = fmt.Printf
)

func (a *Ast) dump(level int) {
	pt("%s%s\n", strings.Repeat("  ", level), a.Op.String())
	if a.Left != nil {
		a.Left.dump(level + 1)
	}
	if a.Right != nil {
		a.Right.dump(level + 1)
	}
	if a.Predict != nil {
		pt("%s%v\n", strings.Repeat("  ", level+1), a.Predict)
	}
}

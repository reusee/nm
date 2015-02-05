package nm

import (
	"testing"
	"time"
)

func TestCompile(t *testing.T) {
	code := `html body div#foo div.bar ul li p a []*`
	t0 := time.Now()
	program, err := Compile(code)
	pt("%v\n", time.Now().Sub(t0))
	if err != nil {
		t.Fatal(err)
	}
	_ = program
}

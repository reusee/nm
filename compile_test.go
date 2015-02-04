package nm

import (
	"testing"
	"time"
)

func TestCompile(t *testing.T) {
	code := `A[] B C`
	t0 := time.Now()
	program, err := Compile(code)
	pt("%v\n", time.Now().Sub(t0))
	if err != nil {
		t.Fatal(err)
	}
	_ = program
}

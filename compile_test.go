package nm

import "testing"

func _TestCompile(t *testing.T) { //TODO
	codes := []string{
		`a b`,
		`a*`,
		`a b*`,
		//`html body div#foo div.bar ul li p a []*`,
	}
	for _, code := range codes {
		program := Compile(code)
		for _, inst := range program {
			pt("%v\n", inst)
		}
		pt("\n")
	}
}

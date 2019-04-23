package syntax

import "testing"

func TestGrammerDescription(t *testing.T) {
	gd := NewGrammerDescription("re_grammar.txt")
	for k, v := range gd.content {
		t.Logf("%d %v\n", k, v)
	}
}

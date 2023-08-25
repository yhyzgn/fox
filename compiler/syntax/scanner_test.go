package syntax

import (
	"fmt"
	"strings"
	"testing"
)

func errh(line, col uint, msg string) {
	panic(fmt.Sprintf("%d:%d: %s", line, col, msg))
}

func TestSmoke(t *testing.T) {
	const src = "if (+foo\t+=..123/***/0.9_0e-0i'a'`raw`\"string\"..f;//$"
	tokens := []token{If, Lparen, Op, Name, AssignOp, Dot, Literal, Literal, Literal, Literal, Literal, Dot, Dot, Name, Semi, EOF}

	var got scanner
	got.init(strings.NewReader(src), errh, 0)
	for _, want := range tokens {
		got.next()
		if got.tok != want {
			t.Errorf("%d:%d: got %s; want %s", got.line, got.col, got.tok, want)
			continue
		} else {
			t.Logf("%d:%d: got %s passed;", got.line, got.col, got.tok)
		}
	}
}

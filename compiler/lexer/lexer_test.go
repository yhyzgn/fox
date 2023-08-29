// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-28 17:08
// version: 1.0.0
// desc   :

package lexer

import (
	"fmt"
	"strings"
	"testing"
)

func errHandler(line, col uint, msg string) {
	panic(fmt.Sprintf("%d:%d: %s", line, col, msg))
}

func TestSmoke(t *testing.T) {
	const src = "if (+foo\t+=..123/*块注释*/...666/**文档注释*/0.9_0e-0i'a'`raw`\"string\"....234..f.;->// 行注释\n"
	tokens := []token{If, Lparen, Add, Identifier, AddAssign, Dot, Literal, Comment, DotDot, Literal, Comment, Literal, Literal, Literal, Literal, DotDotDot, Literal, DotDot, Identifier, Dot, Semi, Arrow, Comment, EOF}

	var got lexer
	got.init(strings.NewReader(src), errHandler)
	for _, want := range tokens {
		got.next()
		if got.tok != want {
			t.Errorf("%d:%d: got %s; want %s", got.line, got.col, got.tok, want)
			continue
		}
		t.Logf("%d:%d: got %s, kind = %s, literal = %s", got.line, got.col, got.tok, got.kind.String(), got.literal)
		if got.tok == EOF {
			break
		}
	}
}

func TestIf(t *testing.T) {
	const src = "if x >= 200 { return true }"
	var got lexer
	got.init(strings.NewReader(src), errHandler)
	for {
		got.next()
		t.Logf("%d:%d: got %s, kind = %s, literal = %s", got.line, got.col, got.tok, got.kind.String(), got.literal)
		if got.tok == EOF {
			break
		}
	}
}

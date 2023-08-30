// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-28 15:32
// version: 1.0.0
// desc   :

package lexer

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

type lexer struct {
	source
	line, col uint
	blank     bool
	tok       token
	kind      LiteralKind
	literal   string
	bad       bool
}

func (l *lexer) init(src io.Reader, errHandler typeErrorHandler) {
	l.source.init(src, errHandler)
}
func (l *lexer) errorf(format string, args ...interface{}) {
	l.error(fmt.Sprintf(format, args...))
}

func (l *lexer) errorAtf(offset int, format string, args ...interface{}) {
	l.errHandler(l.line, l.col+uint(offset), fmt.Sprintf(format, args...))
}

func (l *lexer) next() {
redo:
	l.stop()

	l.kind = None
	l.literal = ""

	startLine, startCol := l.pos()
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.nextCh()
	}

	l.line, l.col = l.pos()
	l.blank = l.line > startLine || startCol == columnBase

	l.start()

	if isLetter(l.ch) || l.ch >= utf8.RuneSelf && l.atIdentChar(true) {
		l.nextCh()
		l.ident()
		return
	}

	switch l.ch {
	case -1:
		l.tok = EOF
		l.literal = "EOF"

	case ';':
		l.nextCh()
		l.literal = ";"
		l.tok = Semi

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		l.number(false)

	case '"':
		l.stdString()

	case '`':
		l.rawString()

	case '\'':
		l.rune()

	case '(':
		l.nextCh()
		l.tok = Lparen
		l.literal = "("

	case '[':
		l.nextCh()
		l.tok = Lbracket
		l.literal = "["

	case '{':
		l.nextCh()
		l.tok = Lbrace
		l.literal = "{"

	case ')':
		l.nextCh()
		l.tok = Rparen
		l.literal = ")"

	case ']':
		l.nextCh()
		l.tok = Rbracket
		l.literal = "]"

	case '}':
		l.nextCh()
		l.tok = Rbrace
		l.literal = "}"

	case ',':
		l.nextCh()
		l.tok = Comma
		l.literal = ","

	case ':':
		l.nextCh()
		if l.ch == '=' {
			l.nextCh()
			l.tok = Define
			l.literal = ":="
			break
		}
		l.tok = Colon
		l.literal = ":"

	case '.':
		l.nextCh()
		if isDecimal(l.ch) {
			l.number(true)
			break
		}
		if l.ch == '.' {
			l.nextCh()
			if isDecimal(l.ch) {
				l.tok = Dot
				l.literal = "."
				break
			}
			if l.ch == '.' {
				l.nextCh()
				if isDecimal(l.ch) {
					l.tok = DotDot
					l.literal = ".."
					break
				}
				l.tok = DotDotDot
				l.literal = "..."
				break
			}
			l.tok = DotDot
			l.literal = ".."
			break
		}
		l.tok = Dot
		l.literal = "."

	case '+':
		l.nextCh()
		if l.ch == '+' {
			l.nextCh()
			l.tok = AddAdd
			l.literal = "++"
			break
		}
		if l.ch == '=' {
			l.nextCh()
			l.tok = AddAssign
			l.literal = "+="
			break
		}
		l.tok = Add
		l.literal = "+"

	case '-':
		l.nextCh()
		if l.ch == '-' {
			l.nextCh()
			l.tok = SubSub
			l.literal = "--"
			break
		}
		if l.ch == '=' {
			l.nextCh()
			l.tok = SubAssign
			l.literal = "-="
			break
		}
		if l.ch == '>' {
			l.nextCh()
			l.tok = Arrow
			l.literal = "->"
			break
		}
		l.tok = Sub
		l.literal = "-"

	case '*':
		l.nextCh()
		if l.ch == '=' {
			l.nextCh()
			l.tok = MulAssign
			l.literal = "*="
			break
		}
		l.tok = Mul
		l.literal = "*"

	case '/':
		l.nextCh()
		if l.ch == '/' {
			l.nextCh()
			l.lineComment()
			break
		}
		if l.ch == '*' {
			l.nextCh()
			if l.ch == '*' {
				l.nextCh()
				l.docComment()
				break
			}
			l.blockComment()
			break
		}
		if l.ch == '=' {
			l.nextCh()
			l.tok = DivAssign
			l.literal = "/="
			break
		}
		l.tok = Div
		l.literal = "/"

	case '%':
		l.nextCh()
		if l.ch == '=' {
			l.nextCh()
			l.tok = ModAssign
			l.literal = "%="
			break
		}
		l.tok = Mod
		l.literal = "%"

	case '&':
		l.nextCh()
		if l.ch == '&' {
			l.nextCh()
			l.tok = AndAnd
			l.literal = "&&"
			break
		}
		if l.ch == '=' {
			l.nextCh()
			l.tok = AndAssign
			l.literal = "&="
			break
		}
		if l.ch == '^' {
			l.nextCh()
			l.tok = AndNot
			l.literal = "&^"
			break
		}
		l.tok = And
		l.literal = "&"

	case '|':
		l.nextCh()
		if l.ch == '|' {
			l.nextCh()
			l.tok = OrOr
			l.literal = "||"
			break
		}
		if l.ch == '=' {
			l.nextCh()
			l.tok = OrAssign
			l.literal = "|="
			break
		}
		l.tok = Or
		l.literal = "|"

	case '^':
		l.nextCh()
		if l.ch == '=' {
			l.nextCh()
			l.tok = XorAssign
			l.literal = "^="
			break
		}
		l.tok = Xor
		l.literal = "^"

	case '<':
		l.nextCh()
		if l.ch == '=' {
			l.nextCh()
			l.tok = Leq
			l.literal = "<="
			break
		}
		if l.ch == '<' {
			l.nextCh()
			if l.ch == '=' {
				l.nextCh()
				l.tok = ShlAssign
				l.literal = "<<="
				break
			}
			l.tok = Shl
			l.literal = "<<"
			break
		}
		if l.ch == '-' {
			l.nextCh()
			l.tok = Receive
			l.literal = "<-"
			break
		}
		l.tok = Lss
		l.literal = "<"

	case '>':
		l.nextCh()
		if l.ch == '=' {
			l.nextCh()
			l.tok = Geq
			l.literal = ">="
			break
		}
		if l.ch == '>' {
			l.nextCh()
			if l.ch == '=' {
				l.nextCh()
				l.tok = ShrAssign
				l.literal = ">>="
				break
			}
			l.tok = Shr
			l.literal = ">>"
			break
		}
		l.tok = Gtr
		l.literal = ">"

	case '=':
		l.nextCh()
		if l.ch == '=' {
			l.nextCh()
			l.tok = Eql
			l.literal = "=="
			break
		}
		l.tok = Assign
		l.literal = "="

	case '!':
		l.nextCh()
		if l.ch == '=' {
			l.nextCh()
			l.tok = Neq
			l.literal = "!="
			break
		}
		l.tok = Not
		l.literal = "!"

	case '~':
		l.nextCh()
		l.tok = Tilde
		l.literal = "~"

	case '?':
		l.nextCh()
		l.tok = Question
		l.literal = "?"

	default:
		l.errorf("Invalid character %#U", l.ch)
		l.nextCh()
		goto redo
	}
}

func (l *lexer) lineComment() {
	for {
		l.nextCh()
		if l.ch == '\n' {
			l.tok = Comment
			l.kind = LineCommentLit
			l.literal = string(l.segment())
			break
		}
	}
}

func (l *lexer) blockComment() {
	if l.skipComment() {
		l.tok = Comment
		l.kind = BlockCommentLit
		l.literal = string(l.segment())
	}
}

func (l *lexer) docComment() {
	if l.skipComment() {
		l.tok = Comment
		l.kind = DocCommentLit
		l.literal = string(l.segment())
	}
}

func (l *lexer) skipComment() bool {
	for l.ch >= 0 {
		for l.ch == '*' {
			l.nextCh()
			if l.ch == '/' {
				l.nextCh()
				return true
			}
		}
		l.nextCh()
	}
	l.errorAtf(0, "comment not terminated")
	return false
}

func (l *lexer) rune() {
	ok := true
	l.nextCh()

	for n := 0; ; n++ {
		if l.ch == '\'' {
			if ok {
				if n == 0 {
					l.errorf("Empty char literal unescaped")
					ok = false
				} else if n != 1 {
					l.errorAtf(0, "More than one character in char literal")
					ok = false
				}
			}
			l.nextCh()
			break
		}

		if l.ch == '\\' {
			l.nextCh()
			if !l.escape('\'') {
				ok = false
			}
			continue
		}

		if l.ch == '\n' {
			if ok {
				l.errorf("Newline in char literal")
				ok = false
			}
			break
		}

		if l.ch < 0 {
			if ok {
				l.errorAtf(0, "Char literal not terminated")
				ok = false
			}
			break
		}
		l.nextCh()
	}

	l.setLiteral(CharLit, ok)
}

func (l *lexer) rawString() {
	ok := true
	l.nextCh()

	for {
		if l.ch == '`' {
			l.nextCh()
			break
		}

		if l.ch < 0 {
			l.errorAtf(0, "String not terminated")
			ok = false
			break
		}
		l.nextCh()
	}
	l.setLiteral(StringLit, ok)
}

func (l *lexer) stdString() {
	ok := true
	l.nextCh()

	for {
		if l.ch == '"' {
			l.nextCh()
			break
		}

		if l.ch == '\\' {
			l.nextCh()
			if !l.escape('"') {
				ok = false
			}
			continue
		}

		if l.ch == '\n' {
			l.errorf("Newline in string")
			ok = false
			break
		}

		if l.ch < 0 {
			l.errorAtf(0, "String not terminated")
			ok = false
			break
		}
		l.nextCh()
	}
	l.setLiteral(StringLit, ok)
}

func (l *lexer) number(seenPoint bool) {
	ok := true
	kind := IntLit
	base := 10
	prefix := rune(0)
	digsep := 0
	invalid := -1

	// integer
	if !seenPoint {
		if l.ch == '0' {
			l.nextCh()
			switch lower(l.ch) {
			case 'x':
				l.nextCh()
				base, prefix = 16, 'x'
			case 'o':
				l.nextCh()
				base, prefix = 8, 'o'
			case 'b':
				l.nextCh()
				base, prefix = 2, 'b'
			default:
				base, prefix = 10, '0'
				digsep = 1
			}
		}
		digsep |= l.digits(base, &invalid)
		if l.ch == '.' {
			if prefix == 'o' || prefix == 'b' {
				l.errorf("Invalid radix point in %s literal", baseName(base))
				ok = false
			}
			l.nextCh()
			seenPoint = true
		}
	}

	if seenPoint {
		kind = FloatLit
		digsep |= l.digits(base, &invalid)
	}

	if digsep&1 == 0 && ok {
		l.errorf("%s literal has no digits", baseName(base))
		ok = false
	}

	if e := lower(l.ch); e == 'e' || e == 'p' {
		if ok {
			switch {
			case e == 'e' && prefix != 0 && prefix != '0':
				l.errorf("%q exponent requires decimal mantissa", l.ch)
				ok = false
			case e == 'p' && prefix != 'x':
				l.errorf("%q exponent requires hexadecimal mantissa", l.ch)
				ok = false
			}
		}
		l.nextCh()
		kind = FloatLit
		if l.ch == '+' || l.ch == '-' {
			l.nextCh()
		}
		digsep = l.digits(10, nil) | digsep&2
		if digsep&1 == 0 && ok {
			l.errorf("Exponent has no digits")
			ok = false
		}
	} else if prefix == 'x' && kind == FloatLit && ok {
		l.errorf("Hexadecimal mantissa requires a 'p' exponent")
		ok = false
	}

	if l.ch == 'i' {
		kind = ImagLit
		l.nextCh()
	}

	l.setLiteral(kind, ok)

	if kind == IntLit && invalid >= 0 && ok {
		l.errorAtf(invalid, "Invalid digit %q in %s literal", l.literal[invalid], baseName(base))
		ok = false
	}

	if digsep&2 != 0 && ok {
		if i := invalidSep(l.literal); i >= 0 {
			l.errorAtf(i, "'_' must separate successive digits")
			ok = false
		}
	}

	l.bad = !ok
}

func (l *lexer) ident() {
	for isLetter(l.ch) || isDecimal(l.ch) {
		l.nextCh()
	}

	if l.ch >= utf8.RuneSelf {
		for l.atIdentChar(false) {
			l.nextCh()
		}
	}

	lit := l.segment()
	if len(lit) >= 2 {
		if tok := keywordReversedMap[string(lit)]; tok != 0 && tokStrFast(tok) == string(lit) {
			l.tok = tok
			l.literal = string(lit)
			return
		}
	}

	l.literal = string(lit)
	l.tok = Identifier
}

func (l *lexer) digits(base int, invalid *int) (digsep int) {
	if base <= 10 {
		max := rune('0' + base)
		for isDecimal(l.ch) || l.ch == '_' {
			ds := 1
			if l.ch == '_' {
				ds = 2
			} else if l.ch >= max && *invalid < 0 {
				_, col := l.pos()
				*invalid = int(col - l.col)
			}
			digsep |= ds
			l.nextCh()
		}
	} else {
		for isHex(l.ch) || l.ch == '_' {
			ds := 1
			if l.ch == '_' {
				ds = 2
			}
			digsep |= ds
			l.nextCh()
		}
	}
	return
}

func (l *lexer) setLiteral(kind LiteralKind, ok bool) {
	l.tok = Literal
	l.kind = kind
	l.literal = string(l.segment())
	l.bad = !ok
}

func (l *lexer) atIdentChar(first bool) bool {
	switch {
	case unicode.IsLetter(l.ch) || l.ch == '_':
	case unicode.IsDigit(l.ch):
		if first {
			l.errorf("Identifier can't begin with digit %#U", l.ch)
		}
	case l.ch >= utf8.RuneSelf:
		l.errorf("Invalid character %#U in identifier", l.ch)
	default:
		return false
	}
	return true
}

func (l *lexer) escape(quote rune) bool {
	var (
		n            int
		base, max, x uint32
	)

	switch l.ch {
	case quote, 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\':
		l.nextCh()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		l.nextCh()
		n, base, max = 2, 16, 255
	case 'u':
		l.nextCh()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		l.nextCh()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		if l.ch < 0 {
			return true
		}
		l.errorf("Unknown escape")
		return false
	}

	for i := n; i > 0; i-- {
		if l.ch < 0 {
			return true
		}

		d := base
		if isDecimal(l.ch) {
			d = uint32(l.ch) - '0'
		} else if 'a' <= lower(l.ch) && lower(l.ch) <= 'f' {
			d = uint32(lower(l.ch)) - 'a' + 10
		}

		if d > base {
			l.errorf("Invalid character %q in %s escape", l.ch, baseName(int(base)))
			return false
		}

		x = x*base + d
		l.nextCh()
	}

	if x > max && base == 8 {
		l.errorf("Octal escape value %d > 255", x)
		return false
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		l.errorf("Escape is invalid Unicode code point %#U", x)
		return false
	}

	return true
}

func invalidSep(x string) int {
	x1 := ' '
	d := '.'
	i := 0

	if len(x) >= 2 && x[0] == '0' {
		x1 = lower(rune(x[1]))
		if x1 == 'x' || x1 == 'o' || x1 == 'b' {
			d = '0'
			i = 2
		}
	}

	for ; i < len(x); i++ {
		p := d
		d = rune(x[i])
		switch {
		case d == '_':
			if p != '0' {
				return i
			}
		case isDecimal(d) || x1 == 'x' && isHex(d):
			d = '0'
		default:
			if p == '_' {
				return i - 1
			}
			d = '.'
		}
	}
	if d == '_' {
		return len(x) - 1
	}
	return -1
}

func tokStrFast(tok token) string {
	return keywordMap[tok]
}

func baseName(base int) string {
	switch base {
	case 2:
		return "binary"
	case 8:
		return "octal"
	case 10:
		return "decimal"
	case 16:
		return "hexadecimal"
	}
	panic("invalid base")
}

func lower(ch rune) rune {
	return ('a' - 'A') | ch
}

func isLetter(ch rune) bool {
	return 'a' <= lower(ch) && lower(ch) <= 'z' || ch == '_'
}

func isDecimal(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isHex(ch rune) bool {
	return '0' <= ch && ch <= '9' || 'a' <= lower(ch) && lower(ch) <= 'f'
}

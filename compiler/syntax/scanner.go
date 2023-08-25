package syntax

import (
	"fmt"
	"io"
	"unicode"
	"unicode/utf8"
)

const (
	comments   uint = 1 << iota // call handler for all comments
	directives                  // call handler for directives only
)

var keywordMap map[string]token

type scanner struct {
	source
	mode      uint
	line, col uint
	blank     bool
	lit       string
	tok       token
	kind      LitKind
	bad       bool
	op        token
	prec      int
}

func init() {
	keywordMap = make(map[string]token, 1<<8)
	for tok := Package; tok <= Run; tok++ {
		if keywordMap[tok.String()] != 0 {
			panic("imperfect token")
		}
		keywordMap[tok.String()] = tok
	}
}

func (s *scanner) init(src io.Reader, errHandler typeErrorHandler, mode uint) {
	s.source.init(src, errHandler)
	s.mode = mode
}

func (s *scanner) errorf(format string, args ...interface{}) {
	s.error(fmt.Sprintf(format, args...))
}

func (s *scanner) errorAtf(offset int, format string, args ...interface{}) {
	s.errHandler(s.line, s.col+uint(offset), fmt.Sprintf(format, args...))
}

func (s *scanner) next() {
redo:
	s.stop()
	startLine, startCol := s.pos()
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.nextCh()
	}

	s.line, s.col = s.pos()
	s.blank = s.line > startLine || startCol == columnBase
	s.start()
	if isLetter(s.ch) || s.ch >= utf8.RuneSelf && s.atIdentChar(true) {
		s.nextCh()
		s.ident()
		return
	}

	switch s.ch {
	case -1:
		s.tok = EOF

	case '\n':
		s.nextCh()
		s.lit = "newline"
		s.tok = Semi

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.number(false)

	case '"':
		s.stdString()

	case '`':
		s.rawString()

	case '\'':
		s.rune()

	case '(':
		s.nextCh()
		s.tok = Lparen

	case '[':
		s.nextCh()
		s.tok = Lbracket

	case '{':
		s.nextCh()
		s.tok = Lbrace

	case ')':
		s.nextCh()
		s.tok = Rparen

	case ']':
		s.nextCh()
		s.tok = Rbracket

	case '}':
		s.nextCh()
		s.tok = Rbrace

	case ',':
		s.nextCh()
		s.tok = Comma

	case ';':
		s.nextCh()
		s.lit = "semicolon"
		s.tok = Semi

	case ':':
		s.nextCh()
		if s.ch == '=' {
			s.nextCh()
			s.tok = Define
			break
		}
		s.tok = Colon

	case '.':
		s.nextCh()
		if isDecimal(s.ch) {
			s.number(true)
			break
		}

		if s.ch == '.' {
			s.nextCh()
			if s.ch == '.' {
				s.nextCh()
				s.tok = DotDotDot
				break
			}
			s.rewind()
			s.nextCh()
		}
		s.tok = Dot

	case '+':
		s.nextCh()
		s.op, s.prec = Add, precAdd
		if s.ch != '+' {
			goto assignop
		}
		s.nextCh()
		s.tok = IncOp

	case '-':
		s.nextCh()
		s.op, s.prec = Sub, precAdd
		if s.ch != '-' {
			goto assignop
		}
		s.nextCh()
		s.tok = IncOp

	case '*':
		s.nextCh()
		s.op, s.prec = Mul, precMul
		if s.ch != '-' {
			goto assignop
		}
		s.nextCh()
		s.tok = Star

	case '/':
		s.nextCh()
		if s.ch == '/' {
			s.nextCh()
			s.lineComment()
			goto redo
		}

		if s.ch == '*' {
			s.nextCh()
			s.fullComment()
			if line, _ := s.pos(); line > s.line {
				s.lit = "newline"
				s.tok = Semi
				break
			}
			goto redo
		}

		s.op, s.prec = Div, precMul
		goto assignop

	case '%':
		s.nextCh()
		s.op, s.prec = Rem, precMul
		goto assignop

	case '&':
		s.nextCh()
		if s.ch == '&' {
			s.nextCh()
			s.op, s.prec = AndAnd, precAndAnd
			s.tok = Op
			break
		}
		s.op, s.prec = And, precMul
		if s.ch == '^' {
			s.nextCh()
			s.op = AndNot
		}
		goto assignop

	case '|':
		s.nextCh()
		if s.ch == '|' {
			s.nextCh()
			s.op, s.prec = OrOr, precOrOr
			s.tok = Op
			break
		}
		s.op, s.prec = Or, precAdd
		goto assignop

	case '^':
		s.nextCh()
		s.op, s.prec = Xor, precAdd
		goto assignop

	case '<':
		s.nextCh()
		if s.ch == '=' {
			s.nextCh()
			s.op, s.prec = Leq, precCmp
			s.tok = Op
			break
		}
		if s.ch == '<' {
			// todo 泛型考虑
			s.nextCh()
			s.op, s.prec = Shl, precMul
			goto assignop
		}
		if s.ch == '-' {
			s.nextCh()
			s.tok = Receive
			break
		}
		s.op, s.prec = Lss, precCmp
		s.tok = Op

	case '>':
		s.nextCh()
		if s.ch == '=' {
			s.nextCh()
			s.op, s.prec = Geq, precCmp
			s.tok = Op
			break
		}
		if s.ch == '>' {
			// todo 泛型考虑
			s.nextCh()
			s.op, s.prec = Shr, precMul
			goto assignop
		}
		s.op, s.prec = Gtr, precCmp
		s.tok = Op

	case '=':
		s.nextCh()
		if s.ch == '=' {
			s.nextCh()
			s.op, s.prec = Eql, precCmp
			s.tok = Op
			break
		}
		s.tok = Assign

	case '!':
		s.nextCh()
		if s.ch == '=' {
			s.nextCh()
			s.op, s.prec = Neq, precCmp
			s.tok = Op
			break
		}
		s.op, s.prec = Not, 0
		s.tok = Op

	case '~':
		s.nextCh()
		s.op, s.prec = Tilde, 0
		s.tok = Op

	default:
		s.errorf("Invalid Character %#U", s.ch)
		s.nextCh()
		goto redo
	}
	return

assignop:
	if s.ch == '=' {
		s.nextCh()
		s.tok = AssignOp
		return
	}
	s.tok = Op
}

func (s *scanner) lineComment() {
	if s.mode&comments != 0 {
		s.skipLine()
		s.comment(string(s.segment()))
		return
	}

	// are we saving directives? or is this definitely not a directive?
	if s.mode&directives == 0 || (s.ch != 'g' && s.ch != 'l') {
		s.stop()
		s.skipLine()
		return
	}

	// recognize go: or line directives
	prefix := "go:"
	if s.ch == 'l' {
		prefix = "line "
	}
	for _, m := range prefix {
		if s.ch != m {
			s.stop()
			s.skipLine()
			return
		}
		s.nextCh()
	}

	// directive text
	s.skipLine()
	s.comment(string(s.segment()))
}

func (s *scanner) fullComment() {
	/* opening has already been consumed */

	if s.mode&comments != 0 {
		if s.skipComment() {
			s.comment(string(s.segment()))
		}
		return
	}

	if s.mode&directives == 0 || s.ch != 'l' {
		s.stop()
		s.skipComment()
		return
	}

	// recognize line directive
	const prefix = "line "
	for _, m := range prefix {
		if s.ch != m {
			s.stop()
			s.skipComment()
			return
		}
		s.nextCh()
	}

	// directive text
	if s.skipComment() {
		s.comment(string(s.segment()))
	}
}

func (s *scanner) comment(text string) {
	s.errorAtf(0, "%s", text)
}

func (s *scanner) skipLine() {
	// don't consume '\n' - needed for nlsemi logic
	for s.ch >= 0 && s.ch != '\n' {
		s.nextCh()
	}
}

func (s *scanner) skipComment() bool {
	for s.ch >= 0 {
		for s.ch == '*' {
			s.nextCh()
			if s.ch == '/' {
				s.nextCh()
				return true
			}
		}
		s.nextCh()
	}
	s.errorAtf(0, "Comment not terminated")
	return false
}

func (s *scanner) rune() {
	ok := true
	s.nextCh()

	n := 0
	for ; ; n++ {
		if s.ch == '\'' {
			if ok {
				if n == 0 {
					s.errorf("Empty char literal or unescaped '")
					ok = false
				} else if n != 1 {
					s.errorAtf(0, "More than one character in char literal")
					ok = false
				}
			}
			s.nextCh()
			break
		}

		if s.ch == '\\' {
			s.nextCh()
			if !s.escape('\'') {
				ok = false
			}
			continue
		}

		if s.ch == '\n' {
			if ok {
				s.errorf("Newline in char literal")
				ok = false
			}
			break
		}

		if s.ch < 0 {
			if ok {
				s.errorAtf(0, "Char literal not terminated")
				ok = false
			}
			break
		}
		s.nextCh()
	}
	s.setLit(CharLit, ok)
}

func (s *scanner) rawString() {
	ok := true
	s.nextCh()

	for {
		if s.ch == '`' {
			s.nextCh()
			break
		}

		if s.ch < 0 {
			s.errorAtf(0, "String not terminated")
			ok = false
			break
		}
		s.nextCh()
	}

	s.setLit(StringLit, ok)
}

func (s *scanner) stdString() {
	ok := true
	s.nextCh()

	for {
		if s.ch == '"' {
			s.nextCh()
			break
		}

		if s.ch == '\\' {
			s.nextCh()
			if !s.escape('"') {
				ok = false
			}
			continue
		}

		if s.ch == '\n' {
			s.errorf("Newline in string")
			ok = false
			break
		}

		if s.ch < 0 {
			s.errorAtf(0, "String not terminated")
			ok = false
			break
		}
		s.nextCh()
	}
	s.setLit(StringLit, ok)
}

func (s *scanner) number(seenPoint bool) {
	ok := true
	kind := IntLit
	base := 10
	prefix := rune(0)
	digsep := 0
	invalid := -1

	// integer
	if !seenPoint {
		if s.ch == '0' {
			s.nextCh()
			switch lower(s.ch) {
			case 'x':
				s.nextCh()
				base, prefix = 16, 'x'
			case 'o':
				s.nextCh()
				base, prefix = 8, 'o'
			case 'b':
				s.nextCh()
				base, prefix = 2, 'b'
			default:
				base, prefix = 10, '0'
				digsep = 1
			}
		}
		digsep |= s.digits(base, &invalid)
		if s.ch == '.' {
			if prefix == 'o' || prefix == 'b' {
				s.errorf("Invalid radix point in %s literal", baseName(base))
				ok = false
			}
			s.nextCh()
			seenPoint = true
		}
	}

	if seenPoint {
		kind = FloatLit
		digsep |= s.digits(base, &invalid)
	}

	if digsep&1 == 0 && ok {
		s.errorf("%s literal has no digits", baseName(base))
		ok = false
	}

	if e := lower(s.ch); e == 'e' || e == 'p' {
		if ok {
			switch {
			case e == 'e' && prefix != 0 && prefix != '0':
				s.errorf("%q exponent requires decimal mantissa", s.ch)
				ok = false
			case e == 'p' && prefix != 'x':
				s.errorf("%q exponent requires hexadecimal mantissa", s.ch)
				ok = false
			}
		}
		s.nextCh()
		kind = FloatLit
		if s.ch == '+' || s.ch == '-' {
			s.nextCh()
		}
		digsep = s.digits(10, nil) | digsep&2
		if digsep&1 == 0 && ok {
			s.errorf("Exponent has no digits")
			ok = false
		}
	} else if prefix == 'x' && kind == FloatLit && ok {
		s.errorf("Hexadecimal mantissa requires a 'p' exponent")
		ok = false
	}

	if s.ch == 'i' {
		kind = ImagLit
		s.nextCh()
	}

	s.setLit(kind, ok)

	if kind == IntLit && invalid >= 0 && ok {
		s.errorAtf(invalid, "Invalid digit %q in %s literal", s.lit[invalid], baseName(base))
		ok = false
	}

	if digsep&2 != 0 && ok {
		if i := invalidSep(s.lit); i >= 0 {
			s.errorAtf(i, "'_' must separate successive digits")
			ok = false
		}
	}

	s.bad = !ok
}

func (s *scanner) digits(base int, invalid *int) (digsep int) {
	if base <= 10 {
		max := rune('0' + base)
		for isDecimal(s.ch) || s.ch == '_' {
			ds := 1
			if s.ch == '_' {
				ds = 2
			} else if s.ch >= max && *invalid < 0 {
				_, col := s.pos()
				*invalid = int(col - s.col)
			}
			digsep |= ds
			s.nextCh()
		}
	} else {
		for isHex(s.ch) || s.ch == '_' {
			ds := 1
			if s.ch == '_' {
				ds = 2
			}
			digsep |= ds
			s.nextCh()
		}
	}
	return
}

func (s *scanner) setLit(kind LitKind, ok bool) {
	s.tok = Literal
	s.lit = string(s.segment())
	s.bad = !ok
	s.kind = kind
}

func (s *scanner) ident() {
	for isLetter(s.ch) || isDecimal(s.ch) {
		s.nextCh()
	}

	if s.ch >= utf8.RuneSelf {
		for s.atIdentChar(false) {
			s.nextCh()
		}
	}

	lit := s.segment()
	if len(lit) >= 2 {
		if tok := keywordMap[string(lit)]; tok != 0 && tokStrFast(tok) == string(lit) {
			s.tok = tok
			return
		}
	}

	s.lit = string(lit)
	s.tok = Name
}

func (s *scanner) escape(quote rune) bool {
	var (
		n            int
		base, max, x uint32
	)

	switch s.ch {
	case quote, 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\':
		s.nextCh()
		return true
	case '0', '1', '2', '3', '4', '5', '6', '7':
		n, base, max = 3, 8, 255
	case 'x':
		s.nextCh()
		n, base, max = 2, 16, 255
	case 'u':
		s.nextCh()
		n, base, max = 4, 16, unicode.MaxRune
	case 'U':
		s.nextCh()
		n, base, max = 8, 16, unicode.MaxRune
	default:
		if s.ch < 0 {
			return true
		}
		s.errorf("Unknown escape")
		return false
	}

	for i := n; i > 0; i-- {
		if s.ch < 0 {
			return true
		}

		d := base
		if isDecimal(s.ch) {
			d = uint32(s.ch) - '0'
		} else if 'a' <= lower(s.ch) && lower(s.ch) <= 'f' {
			d = uint32(lower(s.ch)) - 'a' + 10
		}

		if d > base {
			s.errorf("Invalid character %q in %s escape", s.ch, baseName(int(base)))
			return false
		}

		x = x*base + d
		s.nextCh()
	}

	if x > max && base == 8 {
		s.errorf("Octal escape value %d > 255", x)
		return false
	}

	if x > max || 0xD800 <= x && x < 0xE000 {
		s.errorf("Escape is invalid Unicode code point %#U", x)
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
	return _token_name[_token_index[tok-1]:_token_index[tok]]
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

func (s *scanner) atIdentChar(first bool) bool {
	switch {
	case unicode.IsLetter(s.ch) || s.ch == '_':
	case unicode.IsDigit(s.ch):
		if first {
			s.errorf("Identifier can't begin with digit %#U", s.ch)
		}
	case s.ch >= utf8.RuneSelf:
		s.errorf("Invalid character %#U in identifier", s.ch)
	default:
		return false
	}
	return true
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

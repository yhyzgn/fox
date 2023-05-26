// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-07-14 11:32
// version: 1.0.0
// desc   :

package syntax

import (
	"io"
	"unicode/utf8"
)

type typeErrorHandler func(line, column uint, message string)

// The source buffer is accessed using three indices b (begin),
// r (read), and e (end):
//
// - If b >= 0, it points to the beginning of a segment of most
//   recently read characters (typically a Go literal).
//
// - r points to the byte immediately following the most recently
//   read character ch, which starts at r-chw.
//
// - e points to the byte immediately following the last byte that
//   was read into the buffer.
//
// The buffer content is terminated at buf[e] with the sentinel
// character utf8.RuneSelf. This makes it possible to test for
// the common case of ASCII characters with a single 'if' (see
// nextch method).
//
//                +------ content in use -------+
//                v                             v
// buf [...read...|...segment...|ch|...unread...|s|...free...]
//                ^             ^  ^            ^
//                |             |  |            |
//                b         r-chw  r            e
//
// Invariant: -1 <= b < r <= e < len(buf) && buf[e] == sentinel

type source struct {
	in           io.Reader
	buf          []byte
	line, column uint
	ch           rune
	ioErr        error
	errHandler   typeErrorHandler
	b, r, e      int
	chw          int
}

const (
	sentinel   = utf8.RuneSelf
	lineBase   = 1
	columnBase = 1
	bom        = 0xfeff
)

func (s *source) init(in io.Reader, errHandler typeErrorHandler) {
	s.in = in
	s.errHandler = errHandler
	if nil == s.buf {
		s.buf = make([]byte, nextSize(0))
	}
	s.buf[0] = sentinel
	s.line, s.column = 0, 0
	s.ch = ' '
	s.b, s.r, s.e = -1, 0, 0
	s.chw = 0
	s.ioErr = nil
}

func (s *source) pos() (line, column uint) {
	return lineBase + s.line, columnBase + s.column
}

func (s *source) start() {
	s.b = s.r - s.chw
}

func (s *source) stop() {
	s.b = -1
}

func (s *source) segment() []byte {
	return s.buf[s.b : s.r-s.chw]
}

func (s *source) error(msg string) {
	line, column := s.pos()
	s.errHandler(line, column, msg)
}

func (s *source) rewind() {
	if s.b < 0 {
		panic("No active segment")
	}
	s.column -= uint(s.r - s.b)
	s.r = s.b
	s.nextCh()
}

func (s *source) nextCh() {
redo:
	s.column += uint(s.chw)
	if s.ch == '\n' {
		s.line++
		s.column = 0
	}

	if s.ch = rune(s.buf[s.r]); s.ch < sentinel {
		s.r++
		s.chw = 1
		if s.ch == 0 {
			s.error("Invalid NUL Character")
			goto redo
		}
		return
	}

	if s.e-s.r < utf8.UTFMax && !utf8.FullRune(s.buf[s.r:s.e]) && nil == s.ioErr {
		s.fill()
	}

	// EOF
	if s.r == s.e {
		if s.ioErr != io.EOF {
			s.error("IO error: " + s.ioErr.Error())
			s.ioErr = nil
		}
		s.ch = -1
		s.chw = 0
		return
	}

	s.ch, s.chw = utf8.DecodeRune(s.buf[s.r:s.e])
	s.r += s.chw

	if s.ch == utf8.RuneError && s.chw == 1 {
		s.error("Invalid UTF-8 Encoding")
		goto redo
	}

	// BOM
	if s.ch == bom {
		if s.line > 0 || s.column > 0 {
			s.error("Invalid bom int the middle of the file")
		}
		goto redo
	}
}

func (s *source) fill() {
	b := s.r
	if s.b > 0 {
		b = s.b
		s.b = 0
	}
	content := s.buf[b:s.e]

	if len(content)*2 > len(s.buf) {
		s.buf = make([]byte, nextSize(len(s.buf)))
		copy(s.buf, content)
	} else if b > 0 {
		copy(s.buf, content)
	}

	s.r -= b
	s.e -= b

	for i := 0; i < 10; i++ {
		var n int
		n, s.ioErr = s.in.Read(s.buf[s.e : len(s.buf)-1])
		if n < 0 {
			panic("Negative Read")
		}
		if n > 0 || nil != s.ioErr {
			s.e += n
			s.buf[s.e] = sentinel
			return
		}
	}

	s.buf[s.e] = sentinel
	s.ioErr = io.ErrNoProgress
}

func nextSize(size int) int {
	const min = 4 << 10
	const max = 1 << 20
	if size < min {
		return min
	}
	if size <= max {
		return size << 1
	}
	return size + max
}

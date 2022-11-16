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

type source struct {
	in           io.Reader
	buf          []byte
	line, column uint
	ch           rune
	err          error
	errHandler   typeErrorHandler
}

const (
	sentinel   = utf8.RuneSelf
	lineBase   = 1
	columnBase = 1
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
	s.err = nil
}

func (s *source) pos() (line, column uint) {
	return lineBase + s.line, columnBase + s.column
}

func (s *source) error(msg string) {
	line, column := s.pos()
	s.errHandler(line, column, msg)
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

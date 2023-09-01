// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-31 18:02
// version: 1.0.0
// desc   :

package syntax

import "fmt"

type SyntaxError struct {
	Pos *Pos
	Msg string
}

func NewError(pos *Pos, msg string) SyntaxError {
	return SyntaxError{
		Pos: pos,
		Msg: msg,
	}
}

func (err SyntaxError) Error() string {
	return fmt.Sprintf("%s: %s", err.Pos, err.Msg)
}

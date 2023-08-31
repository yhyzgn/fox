// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-31 17:55
// version: 1.0.0
// desc   :

package loader

type File struct {
	filename string
	pos      *Pos
	trimmed  bool
	err      *Error
}

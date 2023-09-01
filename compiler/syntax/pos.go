// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-31 18:01
// version: 1.0.0
// desc   :

package syntax

import "fmt"

type Pos struct {
	fileObj   *FileObj
	line, col uint
}

func NewPos(fileObj *FileObj, lineBase, colBase uint) *Pos {
	return &Pos{
		fileObj: fileObj,
		line:    lineBase,
		col:     colBase,
	}
}

func (pos Pos) RelFilename() string { return pos.fileObj.filename }

func (pos Pos) String() string {
	//rel := position_{pos.RelFilename(), pos.RelLine(), pos.RelCol()}
	//abs := position_{pos.Base().Pos().RelFilename(), pos.Line(), pos.Col()}
	//s := rel.String()
	//if rel != abs {
	//	s += "[" + abs.String() + "]"
	//}
	//return s
	return fmt.Sprintf("%s: %d:%d", pos.RelFilename(), pos.line, pos.col)
}

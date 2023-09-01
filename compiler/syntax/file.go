// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-31 17:55
// version: 1.0.0
// desc   :

package syntax

const (
	LineBase = 1
	ColBase  = 1
)

type FileObj struct {
	filename string
	pos      *Pos
	trimmed  bool
}

func NewFileObj(filename string, trimmed bool) *FileObj {
	pos := &Pos{
		line: LineBase,
		col:  ColBase,
	}
	fo := &FileObj{
		filename: filename,
		pos:      pos,
		trimmed:  trimmed,
	}
	pos.fileObj = fo
	return fo
}

type File struct {
	PkgName *Name
}

type fileItem struct {
	file *File
	err  chan SyntaxError
}

func newFileItem(file *File) *fileItem {
	return &fileItem{
		file: file,
		err:  make(chan SyntaxError),
	}
}

func (f *fileItem) error(err SyntaxError) {
	f.err <- err
}

func (f *fileItem) errHandler(err error) {
	f.error(err.(SyntaxError))
}

// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-31 17:31
// version: 1.0.0
// desc   :

package syntax

import (
	"compiler/types"
	"io"
	"os"
	"path/filepath"
)

func LoadPkg(filenames []string) {
	if len(filenames) == 0 {
		// 未手动传入待编译的文件，则自动扫描当前目录的文件列表
		dir, err := os.Getwd()
		if nil != err {
			panic(err)
		}
		files, err := filepath.Glob(dir + string(filepath.Separator) + "*.x")
		if nil != err {
			panic(err)
		}
		filenames = files
	}

	fileItemList := make([]*fileItem, len(filenames))
	for i := range fileItemList {
		item := fileItem{err: make(chan SyntaxError)}
		fileItemList[i] = &item
	}

	go func() {
		for i, filename := range filenames {
			fi := fileItemList[i]

			go func(name string, item *fileItem) {
				defer close(item.err)
				fileObj := NewFileObj(name, false)

				f, err := os.Open(name)
				if nil != err {
					item.error(NewError(fileObj.pos, err.Error()))
					return
				}
				defer f.Close()

				item.file, err = parse(fileObj, f, item.errHandler)
				if nil != err {
					item.error(NewError(fileObj.pos, err.Error()))
				}
			}(filename, fi)
		}
	}()

	for _, it := range fileItemList {
		for e := range it.err {
			panic(e)
		}
	}
}

func parse(fileObj *FileObj, src io.Reader, errHandler types.ErrHandler) (_ *File, first error) {
	defer func() {
		if p := recover(); p != nil {
			if err, ok := p.(SyntaxError); ok {
				first = err
				return
			}
			panic(p)
		}
	}()

	var psr parser
	psr.init(fileObj, src, errHandler)
	return psr.fileOrNil(), psr.err
}

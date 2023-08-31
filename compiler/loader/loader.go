// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-31 17:31
// version: 1.0.0
// desc   :

package loader

import (
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

	for _, filename := range filenames {
		println(filename)
	}
}

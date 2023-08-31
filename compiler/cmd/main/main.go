// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-25 15:29
// version: 1.0.0
// desc   :

package main

import (
	"compiler/loader"
	"flag"
)

func main() {
	flag.Parse()

	loader.LoadPkg(flag.Args())
}

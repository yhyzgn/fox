// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-25 16:58
// version: 1.0.0
// desc   :

package lexer

var keywordMap map[token]string

func init() {
	keywordMap = make(map[token]string, 1<<6)
	for i := Pkg; i <= Goto; i++ {
		keywordMap[i] = i.String()
	}
}

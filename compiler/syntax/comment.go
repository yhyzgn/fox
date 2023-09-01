// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-09-01 12:35
// version: 1.0.0
// desc   :

package syntax

import "fmt"

type CommentInfo struct {
	content string
	kind    LiteralKind
}

func NewLineComment(content string) CommentInfo {
	return CommentInfo{
		content: content,
		kind:    LineCommentLit,
	}
}

func NewBlockComment(content string) CommentInfo {
	return CommentInfo{
		content: content,
		kind:    BlockCommentLit,
	}
}

func NewDocComment(content string) CommentInfo {
	return CommentInfo{
		content: content,
		kind:    DocCommentLit,
	}
}

func (i *CommentInfo) String() string {
	return fmt.Sprintf("%s: %s", i.kind.String(), i.content)
}

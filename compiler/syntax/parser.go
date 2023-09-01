// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-09-01 10:42
// version: 1.0.0
// desc   :

package syntax

import (
	"compiler/types"
	"errors"
	"fmt"
	"io"
)

type parser struct {
	fileObj    *FileObj
	errHandler types.ErrHandler
	err        error
	pkgName    *Name
	lexer
}

func (p *parser) init(fileObj *FileObj, src io.Reader, errHandler types.ErrHandler) {
	p.fileObj = fileObj
	p.errHandler = errHandler
	p.lexer.init(src, func(line, column uint, msg string) {
		p.errHandler(errors.New(fmt.Sprintf("%d:%d: %s", line, column, msg)))
	})
}

func (p *parser) fileOrNil() *File {
	p.next()

	f := new(File)

	// 注释
	var docComment string
	for {
		if p.tok != Comment {
			break
		}
		if p.kind != DocCommentLit {
			p.next()
			continue
		}
		docComment = p.literal
		p.next()
	}

	// 包声明
	if !p.got(Pkg) {
		p.syntaxError("pkg statement must be first of file.")
		return nil
	}
	p.pkgName = p.name(docComment)
	p.except(Semi)

	//
	return f
}

func (p *parser) got(tok token) bool {
	if p.tok == tok {
		p.next()
		return true
	}
	return false
}

func (p *parser) except(tok token) {
	if !p.got(tok) {
		p.syntaxError("expected [ " + tok.String() + " ], but got [ " + p.literal + " ].")
	}
}

func (p *parser) name(docComment string) *Name {
	if p.tok == Identifier {
		nm := NewName(p.pos(), p.literal, docComment)
		p.next()
		return nm
	}

	n := NewName(p.pos(), "_", docComment)
	p.syntaxError("expected name")
	return n
}

func (p *parser) pos() Pos {
	return p.posAt(p.line, p.col)
}

func (p *parser) posAt(line uint, col uint) Pos {
	return *NewPos(p.fileObj, line, col)
}

func (p *parser) syntaxError(msg string) {
	panic(msg)
}

type Name struct {
	pos        Pos
	literal    string
	docComment CommentInfo
}

func NewName(pos Pos, literal, docComment string) *Name {
	return &Name{
		pos:        pos,
		literal:    literal,
		docComment: NewDocComment(docComment),
	}
}

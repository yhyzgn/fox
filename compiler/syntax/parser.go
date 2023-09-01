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

// SourceFile = PackageDecl ";" { ImportDecl ";" } { TopLevelDecl ";" }
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

	// PackageDecl
	if !p.got(Pkg) {
		p.syntaxError("pkg statement must be first of file.")
		return nil
	}
	p.pkgName = p.name(docComment)
	p.except(Semi)

	// 跳过注释
	for {
		if p.tok != Comment {
			break
		}
		p.next()
	}

	// ImportDecl = "import" ImportSpec ";"
	prev := Import
	for p.tok != EOF {
		if p.tok == Import && prev != Import {
			// import 语句只能在顶部
			p.syntaxError("Import must appeared before other declarations")
		}
		prev = p.tok

		switch p.tok {
		case Import:
			p.next()
			dl := p.importDecl()
			f.DeclList = append(f.DeclList, dl)

		case Const:
			p.next()

		case Pub:
			p.next()

		case Pri:
			p.next()

		case Class:
			p.next()

		case Interface:
			p.next()

		case Enum:
			p.next()

		case Annotate:
			p.next()

		case Fn:
			p.next()

		default:
			p.next()
		}
	}

	return f
}

// ImportDecl = "import" ImportSpec ";"
// ImportSpec = [ PkgPath | PkgName]
// PkgPath = [PkgName | PkgName "." PkgName]
// PkgName = identifier
func (p *parser) importDecl() Decl {
	nm := new(Name)
	nm.pos = p.pos()
	nm.literal = ""
	for p.tok != Semi {
		switch p.tok {
		case Identifier, Dot:
			nm.literal += p.literal
		default:
			p.syntaxError("Unknown symbol token: " + p.literal)
		}
		p.next()
	}
	p.next()
	d := new(ImportDecl)
	d.Name = nm
	return d
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
	// 这里需要报错，并终止程序
	panic(msg)
}

// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-07-14 14:40
// version: 1.0.0
// desc   :

package syntax

type token uint

//go:generate stringer -type token -linecomment token.go

const (
	_ token = iota //

	Name    // name
	Literal // literal

	Assign  // =
	Define  // :=
	Receive // <-
	Star    // *

	Not   // !
	Tilde // ~

	OrOr   // ||
	AndAnd // &&

	Eql // ==
	Neq // !=
	Lss // <
	Leq // <=
	Gtr // >
	Geq // >=

	Add      // +
	Sub      // -
	Or       // |
	And      // &
	Xor      // ^
	AndNot   // &^
	Mul      // *
	Div      // /
	Rem      // %
	Shl      // <<
	Shr      // >>
	Op       // op
	AssignOp // op=
	IncOp    // opop

	Lparen    // (
	Lbracket  // [
	Lbrace    // {
	Rparen    // )
	Rbracket  // ]
	Rbrace    // }
	Comma     // ,
	Semi      // ;
	Colon     // :
	Dot       // .
	DotDotDot // ...

	Package   // package
	Import    // import
	Sealed    // sealed
	Class     // class
	Interface // interface
	Abstract  // abstract
	Fn        // fn
	Enum      // enum
	Public    // public
	Protect   // protect
	Private   // private
	Defer     // defer
	Const     // const
	GOTO      // goto
	This      // this
	Super     // super
	Return    // return
	Break     // break
	Continue  // continue
	Except    // except
	For       // for
	Do        // do
	While     // while
	If        // if
	Else      // else
	When      // when
	Case      // case
	Default   // default
	Static    // static
	Throws    // throws
	Throw     // throw
	Print     // print
	Printf    // printf
	Println   // println
	Nil       // nil
	Chan      // chan
	Run       // run
	EOF       // EOF
	ERROR     // error
)

type LitKind uint8

const (
	IntLit LitKind = iota
	FloatLit
	ImagLit
	CharLit
	StringLit
)

// Operator precedences
const (
	_ = iota
	precOrOr
	precAndAnd
	precCmp
	precAdd
	precMul
)

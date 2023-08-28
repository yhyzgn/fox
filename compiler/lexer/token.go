// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-08-25 15:45
// version: 1.0.0
// desc   :

package lexer

type token uint32

const (
	_ token = iota

	Identifier // identifier
	Literal    // literal

	Pkg       // pkg
	Import    // import
	Class     // class
	Interface // interface
	Enum      // enum
	Abs       // abs
	Sealed    // sealed
	Pub       // pub
	Pro       // pro
	Pri       // pri
	Fn        // fn
	Defer     // defer
	Const     // const
	This      // this
	Super     // super
	Return    // return
	Break     // break
	Continue  // continue
	Expect    // expect
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
	Nil       // nil
	Chan      // chan
	Run       // run
	Goto      // goto

	Byte   // byte
	Short  // short
	Int    // int
	Long   // long
	UByte  // ubyte
	UShort // ushort
	UInt   // uint
	ULong  // ulong
	Float  // float
	Double // double
	Char   // char
	String // string
	Bool   // bool

	Assign    // =
	Define    // :=
	Receive   // <-
	Arrow     // ->
	Not       // !
	Tilde     // ~
	Add       // +
	Sub       // -
	Or        // |
	And       // &
	Xor       // ^
	AndNot    // &^
	Mul       // *
	Div       // /
	Mod       // %
	Shl       // <<
	Shr       // >>
	OrOr      // ||
	AndAnd    // &&
	Eql       // ==
	Neq       // !=
	Lss       // <
	Leq       // <=
	Gtr       // >
	Geq       // >=
	AddAssign // +=
	SubAssign // -=
	MulAssign // *=
	DivAssign // /=
	ModAssign // %=
	AndAssign // &=
	OrAssign  // |=
	XorAssign // ^=
	ShlAssign // <<=
	ShrAssign // >>=
	AddAdd    // ++
	SubSub    // --
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
	DotDot    // ..
	DotDotDot // ...

	LineComment  // //
	BlockComment // /* comment */
	DocComment   // /** doc */

	EOF   // EOF
	Error // error
)

type LiteralKind uint8

const (
	IntLit LiteralKind = iota
	FloatLit
	ImagLit
	CharLit
	StringLit
)

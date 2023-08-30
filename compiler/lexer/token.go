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
	Annotate  // annotate
	Abs       // abs
	Sealed    // sealed
	Is        // is
	Pub       // pub
	Pri       // pri
	Fn        // fn
	Defer     // defer
	Const     // const
	This      // this
	Super     // super
	Return    // return
	Break     // break
	Continue  // continue
	Try       // try
	Catch     // catch
	Finally   // finally
	For       // for
	If        // if
	Else      // else
	When      // when
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
	Comma     // ,
	Semi      // ;
	Colon     // :
	Define    // :=
	Receive   // <-
	Arrow     // ->
	Not       // !
	Tilde     // ~
	Question  // ?
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
	Dot       // .
	DotDot    // ..
	DotDotDot // ...

	Comment // // /* comment */ /** doc */

	EOF   // EOF
	Error // error
)

type LiteralKind uint8

const (
	None LiteralKind = iota
	IntLit
	FloatLit
	ImagLit
	CharLit
	StringLit
	LineCommentLit
	BlockCommentLit
	DocCommentLit
)

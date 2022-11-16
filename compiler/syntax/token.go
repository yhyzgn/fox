// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2022-07-14 14:40
// version: 1.0.0
// desc   :

package syntax

type token uint

const (
	_          token = iota //
	_EOF                    // EOF
	_Name                   // name
	_Literal                // literal
	_Operator               // op
	_AssignOp               // op=
	_OpOp                   // opop
	_Assign                 // =
	_Define                 // :=
	_Arrow                  // <-
	_Star                   // *
	_Lparen                 // (
	_Lbracket               // [
	_Lbrace                 // {
	_Rparen                 // )
	_Rbracket               // ]
	_Rbrace                 // }
	_Comma                  // ,
	_Semi                   // ;
	_Colon                  // :
	_Dot                    // .
	_DotDotDot              // ...
	_Package                // package
	_Import                 // import
	_Fun                    // fun
	_Class                  // class
	_Interface              // interface
	_Abstract               // abstract
	_Enum                   // enum
	_Public                 // public
	_Protect                // protect
	_Private                // private
	_Defer                  // defer
	_Const                  // const
	_This                   // this
	_Super                  // super
	_Return                 // return
	_Break                  // break
	_Continue               // continue
	_Try                    // try
	_Catch                  // catch
	_Finally                // finally
	_For                    // for
	_Do                     // do
	_While                  // while
	_If                     // if
	_Else                   // else
	_Switch                 // switch
	_Case                   // case
	_Default                // default
	_Static                 // static
	_Final                  // final
	_Throws                 // throws
	_Throw                  // throw
	_Print                  // print
	_Printf                 // printf
	_Println                // println
	_Nil                    // nil
	_Chan                   // chan
	_Run                    // run
)

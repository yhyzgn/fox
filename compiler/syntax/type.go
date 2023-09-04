// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2023-09-01 16:02
// version: 1.0.0
// desc   :

package syntax

import (
	"fmt"
	"go/constant"
)

type (
	Decl interface {
		decl()
	}

	Expr interface {
		typeInfo
		expr()
	}

	ImportDecl struct {
		Name *Name
		declAdapter
	}

	ConstDecl struct {
		NameList []*Name
		Type     Expr
		Values   Expr
		declAdapter
	}
)

type Name struct {
	pos     Pos
	literal string
}

func NewName(pos Pos, literal string) *Name {
	return &Name{
		pos:     pos,
		literal: literal,
	}
}

func (n *Name) String() string {
	return fmt.Sprintf("%s: %s", n.pos.String(), n.literal)
}

type declAdapter struct {
}

func (d *declAdapter) decl() {
}

type typeInfo interface {
	SetTypeInfo(TypeAndValue)
	GetTypeInfo() TypeAndValue
}

type Type interface {
	Underlying() Type
	String() string
}

type TypeAndValue struct {
	Type  Type
	Value constant.Value
	exprFlags
}

type expr struct {
	typeAndValue
}

func (*expr) aExpr() {}

type exprFlags uint8

func (f exprFlags) IsVoid() bool      { return f&1 != 0 }
func (f exprFlags) IsType() bool      { return f&2 != 0 }
func (f exprFlags) IsBuiltin() bool   { return f&4 != 0 }
func (f exprFlags) IsValue() bool     { return f&8 != 0 }
func (f exprFlags) IsNil() bool       { return f&16 != 0 }
func (f exprFlags) Addressable() bool { return f&32 != 0 }
func (f exprFlags) Assignable() bool  { return f&64 != 0 }
func (f exprFlags) HasOk() bool       { return f&128 != 0 }

func (f *exprFlags) SetIsVoid()      { *f |= 1 }
func (f *exprFlags) SetIsType()      { *f |= 2 }
func (f *exprFlags) SetIsBuiltin()   { *f |= 4 }
func (f *exprFlags) SetIsValue()     { *f |= 8 }
func (f *exprFlags) SetIsNil()       { *f |= 16 }
func (f *exprFlags) SetAddressable() { *f |= 32 }
func (f *exprFlags) SetAssignable()  { *f |= 64 }
func (f *exprFlags) SetHasOk()       { *f |= 128 }

type typeAndValue struct {
	tv TypeAndValue
}

func (x *typeAndValue) SetTypeInfo(tv TypeAndValue) {
	x.tv = tv
}
func (x *typeAndValue) GetTypeInfo() TypeAndValue {
	return x.tv
}

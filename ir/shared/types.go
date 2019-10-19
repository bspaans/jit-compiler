package shared

import (
	"strings"

	"github.com/bspaans/jit/lib"
)

//go:generate stringer -type=TypeNr
type TypeNr int

const (
	T_Uint8    TypeNr = iota
	T_Uint64   TypeNr = iota
	T_Float64  TypeNr = iota
	T_Bool     TypeNr = iota
	T_Array    TypeNr = iota
	T_Function TypeNr = iota
	T_Struct   TypeNr = iota
)

type Type interface {
	Type() TypeNr
	String() string
	Width() lib.Size
}

type BaseType struct {
	TypeNr TypeNr
}

func (b *BaseType) Type() TypeNr {
	return b.TypeNr
}

func (b *BaseType) String() string {
	return map[TypeNr]string{
		T_Uint8:    "uint8",
		T_Uint64:   "uint64",
		T_Float64:  "float64",
		T_Bool:     "bool",
		T_Array:    "array",
		T_Function: "func",
	}[b.TypeNr]
}

func (b *BaseType) Width() lib.Size {
	return map[TypeNr]lib.Size{
		T_Uint8:   lib.BYTE,
		T_Uint64:  lib.QUADWORD,
		T_Float64: lib.QUADWORD,
		T_Bool:    lib.BYTE,
	}[b.TypeNr]
}

var (
	TUint8   = &BaseType{T_Uint8}
	TUint64  = &BaseType{T_Uint64}
	TFloat64 = &BaseType{T_Float64}
	TBool    = &BaseType{T_Bool}
)

type TArray struct {
	ItemType Type
	Size     int
}

func (t *TArray) Type() TypeNr {
	return T_Array
}
func (b *TArray) String() string {
	return "[" + b.ItemType.String() + "]"
}
func (b *TArray) Width() lib.Size {
	return lib.QUADWORD
}

type TFunction struct {
	ReturnType Type
	Args       []Type
	ArgNames   []string
}

func (t *TFunction) Type() TypeNr {
	return T_Function
}
func (b *TFunction) String() string {
	args := []string{}
	for _, a := range b.Args {
		args = append(args, a.String())
	}
	return "(" + strings.Join(args, ", ") + ") " + b.ReturnType.String()
}
func (b *TFunction) Width() lib.Size {
	return lib.QUADWORD
}

type TStruct struct {
	FieldTypes []Type
	Fields     []string
}

func (t *TStruct) Type() TypeNr {
	return T_Struct
}
func (b *TStruct) String() string {
	args := []string{}
	for i, a := range b.FieldTypes {
		args = append(args, b.Fields[i]+" "+a.String())
	}
	return "struct {" + strings.Join(args, ", ") + "}"
}
func (b *TStruct) Width() lib.Size {
	return lib.QUADWORD
}

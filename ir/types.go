package ir

type TypeNr int

const (
	T_Uint8   TypeNr = iota
	T_Uint64  TypeNr = iota
	T_Float64 TypeNr = iota
	T_Bool    TypeNr = iota
	T_Array   TypeNr = iota
)

type Type interface {
	Type() TypeNr
}

type BaseType struct {
	TypeNr TypeNr
}

func (b *BaseType) Type() TypeNr {
	return b.TypeNr
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

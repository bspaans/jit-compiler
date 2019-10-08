package shared

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
	String() string
}

type BaseType struct {
	TypeNr TypeNr
}

func (b *BaseType) Type() TypeNr {
	return b.TypeNr
}

func (b *BaseType) String() string {
	return map[TypeNr]string{
		T_Uint8:   "uint8",
		T_Uint64:  "uint64",
		T_Float64: "float64",
		T_Bool:    "bool",
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

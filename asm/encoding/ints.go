package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Value interface {
	Type() Type
	String() string
	Encode() []uint8
}

type Uint8 uint8
type Uint16 uint16
type Uint32 uint32
type Uint64 uint64
type Int32 int32

func (i Uint8) Type() Type {
	return T_Uint8
}
func (i Uint8) String() string {
	return fmt.Sprintf("$%d", i)
}
func (i Uint8) Encode() []uint8 {
	result := make([]byte, 1)
	result[0] = uint8(i)
	return result
}

func (i Uint16) Type() Type {
	return T_Uint16
}
func (i Uint16) String() string {
	return fmt.Sprintf("$%d", i)
}
func (i Uint16) Encode() []uint8 {
	result := make([]byte, 2)
	binary.LittleEndian.PutUint16(result, uint16(i))
	return result
}

func (i Uint32) Type() Type {
	return T_Uint32
}
func (i Uint32) String() string {
	return fmt.Sprintf("$%d", i)
}
func (i Uint32) Encode() []uint8 {
	result := make([]byte, 4)
	binary.LittleEndian.PutUint32(result, uint32(i))
	return result
}

func (i Uint64) Type() Type {
	return T_Uint64
}
func (i Uint64) String() string {
	return fmt.Sprintf("$%d", i)
}
func (i Uint64) Encode() []uint8 {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, uint64(i))
	return result
}

func (i Int32) Type() Type {
	return T_Int32
}
func (i Int32) String() string {
	return fmt.Sprintf("$%d", i)
}
func (i Int32) Encode() []uint8 {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int32(i))
	return buf.Bytes()
}

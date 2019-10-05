package asm

import "errors"

type JNE struct {
	Dest Value
}

func (j *JNE) Encode() (MachineCode, error) {
	if j.Dest == nil {
		return nil, errors.New("Missing destination")
	}
	var result []uint8
	if j.Dest.Type() == T_Uint8 {
		result = []uint8{0x75}
		for _, b := range j.Dest.(Uint8).Encode() {
			result = append(result, b)
		}
	} else {
		return nil, errors.New("Unsupported destination")
	}
	return result, nil

}
func (j *JNE) String() string {
	return "jne " + j.Dest.String()
}

type JMP struct {
	Dest Value
}

func (j *JMP) Encode() (MachineCode, error) {
	if j.Dest == nil {
		return nil, errors.New("Missing destination")
	}
	var result []uint8
	if j.Dest.Type() == T_Uint8 {
		result = []uint8{0xEB}
		for _, b := range j.Dest.(Uint8).Encode() {
			result = append(result, b)
		}
	} else if j.Dest.Type() == T_Uint32 {
		result = []uint8{0xE9}
		for _, b := range j.Dest.(Uint32).Encode() {
			result = append(result, b)
		}
	} else {
		return nil, errors.New("Unsupported destination")
	}
	return result, nil

}
func (j *JMP) String() string {
	return "jmp " + j.Dest.String()
}

package encoding

import (
	"fmt"

	"github.com/bspaans/jit-compiler/lib"
)

type Scale uint8

const (
	Scale1 Scale = 0
	Scale2 Scale = 1
	Scale4 Scale = 2
	Scale8 Scale = 3
)

func ScaleForItemWidth(itemWidth lib.Size) Scale {
	if itemWidth == lib.BYTE {
		return Scale1
	} else if itemWidth == lib.WORD {
		return Scale2
	} else if itemWidth == lib.DOUBLE {
		return Scale4
	} else if itemWidth == lib.QUADWORD {
		return Scale8
	}
	return Scale1
}

func (s Scale) String() string {
	return fmt.Sprintf("%d", 1<<s)
}

type SIB struct {
	Scale Scale // 2 bits
	Index uint8 // 3 bits
	Base  uint8 // 3 bits
}

func NewSIB(scale Scale, index, base uint8) *SIB {
	return &SIB{scale, index, base}
}

func (s *SIB) Encode() uint8 {
	result := s.Base & 7
	result += (s.Index & 7) << 3
	result += (uint8(s.Scale) << 6)
	return result
}

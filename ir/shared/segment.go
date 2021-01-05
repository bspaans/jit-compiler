package shared

import (
	"encoding/hex"
	"strings"
)

type Segment struct {
	Data []uint8
}

func NewSegment() *Segment {
	return &Segment{[]uint8{}}
}

func (s Segment) String() string {
	h := hex.EncodeToString(s.Data)
	result := []rune{' ', ' '}
	for i, c := range h {
		result = append(result, c)
		if i%2 == 1 && i+1 < len(h) {
			result = append(result, ' ')
		}
		if i%16 == 15 && i+1 < len(h) {
			result = append(result, '\n', ' ', ' ')
		}
	}
	return string(result)
}

type SegmentType uint8

const (
	ReadOnly   SegmentType = 1
	ReadWrite  SegmentType = 2
	Executable SegmentType = 3
)

type SegmentPointer struct {
	SegmentType
	Offset uint
}

type Segments struct {
	Segments map[SegmentType]*Segment
}

func NewSegments() *Segments {
	return &Segments{map[SegmentType]*Segment{
		ReadOnly:   NewSegment(),
		ReadWrite:  NewSegment(),
		Executable: NewSegment(),
	}}
}

func (s *Segments) Add(ty SegmentType, data ...uint8) *SegmentPointer {
	offset := len(s.Segments[ty].Data)
	s.Segments[ty].Data = append(s.Segments[ty].Data, data...)
	return &SegmentPointer{
		SegmentType: ty,
		Offset:      uint(offset),
	}
}

func (s *Segments) Encode() []uint8 {
	sub := append(s.Segments[ReadOnly].Data, s.Segments[ReadWrite].Data...)
	return append(sub, s.Segments[Executable].Data...)
}

func (s *Segments) GetAddress(p *SegmentPointer) int {
	if p.SegmentType == ReadOnly {
		return int(p.Offset) + 2
	}
	readOnly := uint(len(s.Segments[ReadOnly].Data))
	if p.SegmentType == ReadWrite {
		return int(readOnly+p.Offset) + 2
	}
	readWrite := uint(len(s.Segments[ReadWrite].Data))
	if p.SegmentType == Executable {
		return int(readOnly+readWrite+p.Offset) + 2
	}
	panic("Unknown segment type")
	return 0
}

func (s *Segments) String() string {
	return strings.Join([]string{
		".rodata",
		s.Segments[ReadOnly].String(),
		".data",
		s.Segments[ReadWrite].String(),
		".text",
		s.Segments[Executable].String(),
	}, "\n")
}

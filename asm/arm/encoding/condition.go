package encoding

type Condition uint8

const (
	// Integer: equal; Float: equal; Condition: Z==1
	EQ Condition = 0b0000
	// Integer: not equal; Float: not equal or unordered; Condition: Z==0
	NE Condition = 0b0001
	// Integer: carry set; Float: greater than, equal or unordered; Condition: C==1
	CS_or_HS Condition = 0b0010
	// Integer: carry clear; Float: less than; Condition: C==0
	CC_or_LO Condition = 0b0011
	// Integer: minus, negative; Float: less than; Condition: N==1
	MI Condition = 0b0100
	// Integer: plus, positive or zero; Float: greater than, equal, unordered; Condition: N==0
	PL Condition = 0b0101
	// Integer: overflow; Float: unordered; Condition: V==1
	VS Condition = 0b0110
	// Integer: no overflow; Float: ordered; Condition: V==0
	VC Condition = 0b0111
	// Integer: unsigned higher; Float: greater than, unordered; Condition: C==1 && Z == 0
	HI Condition = 0b1000
	// Integer: unsigned lower or same; Float: less than or equal; Condition: !(C==1 && Z == 0)
	LS Condition = 0b1001
	// Integer: signed greater than or equal; Float: greater than or equal; Condition: N == V
	GE Condition = 0b1010
	// Integer: signed less than; Float: less than or unordered; Condition: N != V
	LT Condition = 0b1011
	// Integer: signed greater than; Float: greater than; Condition: Z==0 && N==V
	GT Condition = 0b1100
	// Integer: signed less than or equal; Float: less than, equal, unordered; Condition: !(Z==0 && N==V)
	LE Condition = 0b1101
	// Always
	AL Condition = 0b1110
)

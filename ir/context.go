package ir

type IR_Context struct {
	Registers          []bool
	RegistersAllocated uint8
	VariableMap        map[string]uint8
}

func NewIRContext() *IR_Context {
	ctx := &IR_Context{
		Registers:          make([]bool, 8),
		RegistersAllocated: 0,
		VariableMap:        map[string]uint8{},
	}
	// Always allocate rsp
	// Should track usage?
	ctx.Registers[4] = true
	ctx.RegistersAllocated = 1
	return ctx
}

func (i *IR_Context) AllocateRegister() uint8 {
	if i.RegistersAllocated >= 8 {
		panic("Register allocation limit. Needs stack handling")
	}
	for j := 0; j < len(i.Registers); j++ {
		if !i.Registers[j] {
			i.Registers[j] = true
			i.RegistersAllocated += 1
			return uint8(j)
		}
	}
	panic("Register allocation limit reached with incorrect allocation counter. Needs stack handling")
}

func (i *IR_Context) DeallocateRegister(reg uint8) {
	i.Registers[reg] = false
	i.RegistersAllocated -= 1
}

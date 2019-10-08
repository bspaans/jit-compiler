package ir

import (
	"errors"
	"fmt"

	"github.com/bspaans/jit/asm"
)

type IRType int

const (
	Assignment IRType = iota
	If         IRType = iota
	While      IRType = iota
	Return     IRType = iota
	AndThen    IRType = iota
)

type IR_Assignment struct {
	*BaseIR
	Variable string
	Expr     IRExpression
}

func NewIR_Assignment(variable string, expr IRExpression) *IR_Assignment {
	return &IR_Assignment{
		BaseIR:   NewBaseIR(Assignment),
		Variable: variable,
		Expr:     expr,
	}
}

// Allocates a new register and assigns it the value of the expression.
func (i *IR_Assignment) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	r, found := ctx.VariableMap[i.Variable]
	if !found {
		r = ctx.AllocateRegister()
		ctx.VariableMap[i.Variable] = r
		ctx.VariableTypes[i.Variable] = i.Expr.ReturnType(ctx)
	}
	reg := asm.Get64BitRegisterByIndex(r)
	return i.Expr.Encode(ctx, reg)
}

func (i *IR_Assignment) String() string {
	return fmt.Sprintf("%s = %s", i.Variable, i.Expr.String())
}

func (i *IR_Assignment) AddToDataSection(ctx *IR_Context) {
	i.Expr.AddToDataSection(ctx)
}

type IR_If struct {
	*BaseIR
	Condition IRExpression
	Stmt1     IR
	Stmt2     IR
}

func NewIR_If(condition IRExpression, stmt1, stmt2 IR) *IR_If {
	return &IR_If{
		BaseIR:    NewBaseIR(If),
		Condition: condition,
		Stmt1:     stmt1,
		Stmt2:     stmt2,
	}
}

func (i *IR_If) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	// TODO: introduce return type function
	if i.Condition.Type() == Bool || i.Condition.Type() == Equals {
		r := ctx.AllocateRegister()
		defer ctx.DeallocateRegister(r)
		reg := asm.Get64BitRegisterByIndex(r)
		result, err := i.Condition.Encode(ctx, reg)
		if err != nil {
			return nil, err
		}
		stmt1Len, err := IR_Length(i.Stmt1, ctx)
		if err != nil {
			return nil, err
		}
		stmt2Len, err := IR_Length(i.Stmt2, ctx)
		if err != nil {
			return nil, err
		}
		instr := []asm.Instruction{
			&asm.CMP{asm.Uint32(1), reg},
			&asm.JNE{asm.Uint8(stmt1Len)},
		}
		for _, inst := range instr {
			ctx.AddInstruction(inst)
			result = append(result, inst)
		}
		s1, err := i.Stmt1.Encode(ctx)
		if err != nil {
			return nil, err
		}
		for _, instr := range s1 {
			result = append(result, instr)
		}
		jmp := &asm.JMP{asm.Uint8(stmt2Len)}
		ctx.AddInstruction(jmp)
		result = append(result, jmp)

		s2, err := i.Stmt2.Encode(ctx)
		if err != nil {
			return nil, err
		}
		for _, instr := range s2 {
			result = append(result, instr)
		}
		return result, nil
	}
	return nil, errors.New("Unsupported if IR expression")
}

func (i *IR_If) String() string {
	return fmt.Sprintf("if %s; then %s; else %s;", i.Condition.String(), i.Stmt1.String(), i.Stmt2.String())
}

func (i *IR_If) AddToDataSection(ctx *IR_Context) {
	i.Condition.AddToDataSection(ctx)
	i.Stmt1.AddToDataSection(ctx)
	i.Stmt2.AddToDataSection(ctx)
}

type IR_Return struct {
	*BaseIR
	Expr IRExpression
}

func NewIR_Return(expr IRExpression) *IR_Return {
	return &IR_Return{
		BaseIR: NewBaseIR(Return),
		Expr:   expr,
	}
}

func (i *IR_Return) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	result, err := i.Expr.Encode(ctx, asm.Rax)
	if err != nil {
		return nil, err
	}
	instr := []asm.Instruction{
		&asm.MOV{asm.Rax, &asm.DisplacedRegister{asm.Rsp, 8}},
		&asm.RET{},
	}
	for _, inst := range instr {
		ctx.AddInstruction(inst)
		result = append(result, inst)
	}
	return result, nil
}

func (i *IR_Return) String() string {
	return fmt.Sprintf("return %s", i.Expr.String())
}

func (i *IR_Return) AddToDataSection(ctx *IR_Context) {
	i.Expr.AddToDataSection(ctx)
}

type IR_While struct {
	*BaseIR
	Condition IRExpression
	Stmt      IR
}

func NewIR_While(condition IRExpression, stmt IR) *IR_While {
	return &IR_While{
		BaseIR:    NewBaseIR(While),
		Condition: condition,
		Stmt:      stmt,
	}
}

func (i *IR_While) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	r := ctx.AllocateRegister()
	defer ctx.DeallocateRegister(r)
	reg := asm.Get64BitRegisterByIndex(r)
	beginning := ctx.InstructionPointer
	result, err := i.Condition.Encode(ctx, reg)
	if err != nil {
		return nil, err
	}
	stmtLen, err := IR_Length(i.Stmt, ctx)
	if err != nil {
		return nil, err
	}
	instr := []asm.Instruction{
		&asm.CMP{asm.Uint32(1), reg},
		&asm.JNE{asm.Uint8(stmtLen + 2)},
	}
	for _, inst := range instr {
		ctx.AddInstruction(inst)
		result = append(result, inst)
	}
	s1, err := i.Stmt.Encode(ctx)
	if err != nil {
		return nil, err
	}
	for _, instr := range s1 {
		result = append(result, instr)
	}
	jmp := &asm.JMP{asm.Uint8(uint8(0xff - (int(ctx.InstructionPointer+1) - int(beginning))))}
	result = append(result, jmp)
	ctx.AddInstruction(jmp)
	fmt.Println("InstructionPointer", ctx.InstructionPointer, beginning, ctx.InstructionPointer-beginning)
	return result, nil
}

func (i *IR_While) String() string {
	return fmt.Sprintf("while %s { %s }", i.Condition.String(), i.Stmt.String())
}

func (i *IR_While) AddToDataSection(ctx *IR_Context) {
	i.Condition.AddToDataSection(ctx)
	i.Stmt.AddToDataSection(ctx)
}

type IR_AndThen struct {
	*BaseIR
	Stmt1 IR
	Stmt2 IR
}

func NewIR_AndThen(stmt1, stmt2 IR) *IR_AndThen {
	return &IR_AndThen{
		BaseIR: NewBaseIR(AndThen),
		Stmt1:  stmt1,
		Stmt2:  stmt2,
	}
}

func (i *IR_AndThen) Encode(ctx *IR_Context) ([]asm.Instruction, error) {
	result, err := i.Stmt1.Encode(ctx)
	if err != nil {
		return nil, err
	}
	s2, err := i.Stmt2.Encode(ctx)
	if err != nil {
		return nil, err
	}
	for _, instr := range s2 {
		result = append(result, instr)
	}
	return result, nil
}

func (i *IR_AndThen) String() string {
	return fmt.Sprintf("%s ; %s", i.Stmt1.String(), i.Stmt2.String())
}

func (i *IR_AndThen) AddToDataSection(ctx *IR_Context) {
	i.Stmt1.AddToDataSection(ctx)
	i.Stmt2.AddToDataSection(ctx)
}

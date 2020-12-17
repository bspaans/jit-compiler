package shared

import "fmt"

type SSA_Context struct {
	variableCounter int
}

func NewSSA_Context() *SSA_Context {
	return &SSA_Context{}
}

func (s *SSA_Context) GenerateVariable() string {
	s.variableCounter++
	return fmt.Sprintf("__ssa_%d", s.variableCounter)
}

type SSA_Rewrite struct {
	Variable string
	Expr     IRExpression
}

func NewSSA_Rewrite(v string, e IRExpression) SSA_Rewrite {
	return SSA_Rewrite{v, e}
}

type SSA_Rewrites []SSA_Rewrite

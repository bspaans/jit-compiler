package ir

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bspaans/jit-compiler/ir/expr"
	"github.com/bspaans/jit-compiler/ir/shared"
	"github.com/bspaans/jit-compiler/ir/statements"
)

type ParseResult struct {
	Result interface{}
	Error  error
	Rest   string
}

func NilParseResult(str string) *ParseResult {
	return &ParseResult{
		Result: nil,
		Rest:   str,
	}
}
func ParseSuccess(val interface{}, rest string) *ParseResult {
	return &ParseResult{
		Result: val,
		Rest:   rest,
	}
}

func ParseError(err error) *ParseResult {
	return &ParseResult{
		Error: err,
	}
}

func Lazy(p func() Parser) Parser {
	return func(str string) *ParseResult {
		return p()(str)
	}
}

type Parser func(string) *ParseResult

func (p Parser) And(f Parser) Parser {
	return p.AndThen(func(_ *ParseResult) Parser {
		return f
	})
}

func (p Parser) AndThen(f func(*ParseResult) Parser) Parser {
	return func(str string) *ParseResult {
		pResult := p(str)
		if pResult.Result == nil {
			return pResult
		}
		if pResult.Error != nil {
			return pResult
		}
		return f(pResult)(pResult.Rest)
	}
}

func (p Parser) Fmap(f func(*ParseResult) *ParseResult) Parser {
	return func(str string) *ParseResult {
		pResult := p(str)
		if pResult.Result == nil {
			return pResult
		}
		if pResult.Error != nil {
			return pResult
		}
		return f(pResult)
	}
}

func (p Parser) Success(value interface{}) Parser {
	return func(str string) *ParseResult {
		pResult := p(str)
		if pResult.Result == nil {
			return pResult
		}
		if pResult.Error != nil {
			return pResult
		}
		return ParseSuccess(value, pResult.Rest)
	}
}

func (p Parser) Many() Parser {
	return func(str string) *ParseResult {
		result := []interface{}{}
		for {
			subResult := p(str)
			if subResult.Result == nil {
				break
			}
			if subResult.Error != nil {
				return subResult
			}
			str = subResult.Rest
			result = append(result, subResult.Result)
		}
		return ParseSuccess(result, str)
	}
}

func (p Parser) Many1() Parser {
	return func(str string) *ParseResult {
		return p.Many().Fmap(func(subResult *ParseResult) *ParseResult {
			if subResult.Error != nil {
				return subResult
			}
			if len(subResult.Result.([]interface{})) == 0 {
				return NilParseResult(str)
			}
			return subResult
		})(str)
	}
}

func ParseEnclosed(open, p, closed Parser) Parser {
	return open.And(p).AndThen(func(r *ParseResult) Parser {
		return closed.Fmap(func(p *ParseResult) *ParseResult {
			return ParseSuccess(r.Result, p.Rest)
		})
	})
}
func ParseList(p Parser) Parser {
	return ParseListWithSeparator(p, ParseByte(','))
}

func ParseListWithSeparator(p, separator Parser) Parser {
	return func(str string) *ParseResult {
		result := []interface{}{}
		for {
			sub := p(str)
			if sub.Result == nil && len(result) == 0 {
				return ParseSuccess(result, sub.Rest)
			} else if sub.Result == nil {
				return ParseSuccess(result, sub.Rest)
			} else if sub.Error != nil {
				return sub
			}
			result = append(result, sub.Result)
			str = sub.Rest

			comma := ParseSpace().And(separator).And(ParseWhiteSpace())(str)
			if comma.Result == nil || comma.Error != nil {
				return ParseSuccess(result, sub.Rest)
			}
			str = comma.Rest
		}
	}
}

func OneOf(ps []Parser) Parser {
	return func(str string) *ParseResult {
		for _, p := range ps {
			subResult := p(str)
			if subResult.Result != nil {
				return subResult
			}
		}
		return NilParseResult(str)
	}
}

func ParseByte(char byte) Parser {
	return func(str string) *ParseResult {
		if len(str) == 0 {
			return NilParseResult(str)
		}
		if str[0] != char {
			return NilParseResult(str)
		}
		return ParseSuccess(char, str[1:])
	}
}

func ParseString(s string) Parser {
	return func(str string) *ParseResult {
		if !strings.HasPrefix(str, s) {
			return NilParseResult(str)
		}
		return ParseSuccess(s, str[len(s):])
	}
}

func ParseTypeUint8() Parser {
	return ParseString("uint8").Success(shared.TUint8)
}
func ParseTypeUint64() Parser {
	return ParseString("uint64").Success(shared.TUint64)
}
func ParseTypeFloat64() Parser {
	return ParseString("float64").Success(shared.TFloat64)
}
func ParseSimpleType() Parser {
	return OneOf([]Parser{
		ParseTypeUint8(),
		ParseTypeUint64(),
		ParseTypeFloat64(),
	})
}
func ParseTypeArray() Parser {
	return ParseString("[]").And(Lazy(ParseType)).Fmap(func(p *ParseResult) *ParseResult {
		return ParseSuccess(&shared.TArray{p.Result.(shared.Type), 0}, p.Rest)
	})
}

func ParseType() Parser {
	return OneOf([]Parser{
		ParseSimpleType(),
		ParseTypeArray(),
	})
}

func ParseArrayItems() Parser {
	itemParser := OneOf([]Parser{
		ParseFloat64(),
		ParseUint64(),
		ParseBool(),
	})
	return ParseList(itemParser).Fmap(func(items *ParseResult) *ParseResult {
		return ParseSuccess(InterfaceArrayToIRExpressionArray(items.Result), items.Rest)
	})
}

func ParseArray() Parser {
	return ParseString("[]").And(ParseType()).AndThen(func(elemType *ParseResult) Parser {
		return ParseByte('{').And(ParseSpace()).And(ParseArrayItems()).AndThen(func(elems *ParseResult) Parser {
			return ParseSpace().And(ParseByte('}')).Fmap(func(b *ParseResult) *ParseResult {
				return ParseSuccess(expr.NewIR_StaticArray(elemType.Result.(shared.Type), elems.Result.([]shared.IRExpression)), b.Rest)
			})
		})
	})
}

func ParseBool() Parser {
	return OneOf([]Parser{
		ParseString("true"),
		ParseString("false"),
	}).Fmap(func(b *ParseResult) *ParseResult {
		if b.Result.(string) == "true" {
			return ParseSuccess(expr.NewIR_Bool(true), b.Rest)
		}
		return ParseSuccess(expr.NewIR_Bool(false), b.Rest)
	})
}

func ParseByteRange(start, end byte) Parser {
	return func(str string) *ParseResult {
		if len(str) == 0 {
			return NilParseResult(str)
		}
		if !(str[0] >= start && str[0] <= end) {
			return NilParseResult(str)
		}
		return ParseSuccess(str[0], str[1:])
	}
}

func InterfaceArrayToByteArray(values interface{}) []byte {
	valueArray := values.([]interface{})
	chars := []byte{}
	for _, v := range valueArray {
		chars = append(chars, v.(byte))
	}
	return chars
}

func InterfaceArrayToIRExpressionArray(values interface{}) []shared.IRExpression {
	valueArray := values.([]interface{})
	chars := []shared.IRExpression{}
	for _, v := range valueArray {
		chars = append(chars, v.(shared.IRExpression))
	}
	return chars
}

func ParseUint64() Parser {
	return ParseByteRange('0', '9').Many1().Fmap(func(sub *ParseResult) *ParseResult {
		chars := InterfaceArrayToByteArray(sub.Result)
		int_, err := strconv.Atoi(string(chars))
		return &ParseResult{
			Result: expr.NewIR_Uint64(uint64(int_)),
			Rest:   sub.Rest,
			Error:  err,
		}
	})
}

func ParseFloat64() Parser {
	return ParseByteRange('0', '9').Many1().AndThen(func(s *ParseResult) Parser {
		return ParseByte('.').And(ParseByteRange('0', '9').Many1()).Fmap(func(r *ParseResult) *ParseResult {
			chars := InterfaceArrayToByteArray(s.Result)
			remainder := InterfaceArrayToByteArray(r.Result)
			f := string(chars) + "." + string(remainder)
			float, err := strconv.ParseFloat(f, 64)
			return &ParseResult{
				Result: expr.NewIR_Float64(float),
				Rest:   r.Rest,
				Error:  err,
			}
		})
	})
}

func ParseIdent() Parser {
	return OneOf([]Parser{
		ParseByteRange('a', 'z'),
		ParseByteRange('A', 'Z'),
	}).AndThen(func(sub1 *ParseResult) Parser {
		return OneOf([]Parser{
			ParseByteRange('a', 'z'),
			ParseByteRange('A', 'Z'),
			ParseByteRange('0', '9'),
		}).Many().Fmap(func(sub2 *ParseResult) *ParseResult {
			chars := InterfaceArrayToByteArray(sub2.Result)
			result := string(sub1.Result.(byte)) + string(chars)
			return &ParseResult{
				Result: result,
				Rest:   sub2.Rest,
				Error:  nil,
			}
		})
	})
}

func ParseVariable() Parser {
	return ParseIdent().Fmap(func(ident *ParseResult) *ParseResult {
		ReservedWords := map[string]bool{
			"if":      true,
			"while":   true,
			"uint64":  true,
			"float64": true,
		}
		result := ident.Result.(string)
		if reserved := ReservedWords[result]; reserved {
			return ParseError(fmt.Errorf("%v is a reserved word", reserved))
		}
		return &ParseResult{
			Result: expr.NewIR_Variable(result),
			Rest:   ident.Rest,
			Error:  nil,
		}
	})
}

func ParseOperator() Parser {
	return ParseSingleExpression().AndThen(func(op1 *ParseResult) Parser {
		return ParseSpace().And(OneOf([]Parser{
			ParseString("+"),
			ParseString("-"),
			ParseString("*"),
			ParseString("/"),
			ParseString("=="),
			ParseString("!="),
		})).AndThen(func(op *ParseResult) Parser {
			return ParseSpace().And(ParseExpression()).Fmap(func(op2 *ParseResult) *ParseResult {
				if op.Result.(string) == "+" {
					return ParseSuccess(expr.NewIR_Add(op1.Result.(shared.IRExpression), op2.Result.(shared.IRExpression)), op2.Rest)
				} else if op.Result.(string) == "-" {
					return ParseSuccess(expr.NewIR_Sub(op1.Result.(shared.IRExpression), op2.Result.(shared.IRExpression)), op2.Rest)
				} else if op.Result.(string) == "*" {
					return ParseSuccess(expr.NewIR_Mul(op1.Result.(shared.IRExpression), op2.Result.(shared.IRExpression)), op2.Rest)
				} else if op.Result.(string) == "/" {
					return ParseSuccess(expr.NewIR_Div(op1.Result.(shared.IRExpression), op2.Result.(shared.IRExpression)), op2.Rest)
				} else if op.Result.(string) == "==" {
					return ParseSuccess(expr.NewIR_Equals(op1.Result.(shared.IRExpression), op2.Result.(shared.IRExpression)), op2.Rest)
				} else if op.Result.(string) == "!=" {
					return ParseSuccess(expr.NewIR_Not(expr.NewIR_Equals(op1.Result.(shared.IRExpression), op2.Result.(shared.IRExpression))), op2.Rest)
				}
				return ParseError(errors.New("Unknown operator"))
			})
		})
	})
}

func ParseSpace() Parser {
	return OneOf([]Parser{ParseByte(' '), ParseByte('\t')}).Many()
}
func ParseWhiteSpace() Parser {
	return OneOf([]Parser{ParseByte(' '), ParseByte('\t'), ParseByte('\n')}).Many()
}
func ParseSpace1() Parser {
	return ParseByte(' ').Many1()
}

func ParseSingleExpression() Parser {
	return OneOf([]Parser{
		ParseStructField(),
		ParseStruct(),
		ParseArrayIndex(),
		ParseBool(),
		ParseFloat64(),
		ParseUint64(),
		ParseFunctionCall(),
		ParseFunction(),
		ParseVariable(),
		ParseArray(),
		ParseEnclosedExpression(),
	})
}

func ParseEnclosedExpression() Parser {
	return ParseEnclosed(ParseSpace().And(ParseByte('(')), Lazy(ParseExpression), ParseSpace().And(ParseByte(')')))
}

func ParseExpression() Parser {
	return OneOf([]Parser{
		ParseOperator(),
		ParseSingleExpression(),
	})
}

func ParseSingleStatement() Parser {
	return ParseSpace().And(OneOf([]Parser{
		ParseIf(),
		ParseAssignment(),
		ParseArrayAssignment(),
		ParseReturn(),
		ParseWhile(),
		ParseFunctionDef(),
	}))
}

func ParseStatement() Parser {
	return ParseEnclosed(ParseWhiteSpace(), OneOf([]Parser{
		ParseAndThen(),
		ParseSingleStatement(),
	}), ParseWhiteSpace())
}

func ParseAssignment() Parser {
	return ParseVariable().AndThen(func(variable *ParseResult) Parser {
		return ParseSpace().And(ParseByte('=')).And(ParseSpace()).And(ParseExpression()).Fmap(func(value *ParseResult) *ParseResult {
			v := variable.Result.(*expr.IR_Variable).Value
			return ParseSuccess(statements.NewIR_Assignment(v, value.Result.(shared.IRExpression)), value.Rest)
		})
	})
}

func ParseArrayAssignment() Parser {
	return ParseVariable().AndThen(func(variable *ParseResult) Parser {
		return ParseSpace().And(ParseByte('[')).And(ParseSpace()).And(ParseExpression()).AndThen(func(index *ParseResult) Parser {
			return ParseSpace().And(ParseByte(']')).And(ParseSpace()).And(ParseByte('=')).And(ParseSpace()).And(ParseExpression()).Fmap(func(value *ParseResult) *ParseResult {
				array := variable.Result.(*expr.IR_Variable).Value

				return ParseSuccess(statements.NewIR_ArrayAssignment(array, index.Result.(shared.IRExpression), value.Result.(shared.IRExpression)), value.Rest)
			})
		})
	})
}

func ParseFunctionDefArgs() Parser {
	itemParser := ParseVariable().AndThen(func(variable *ParseResult) Parser {
		return ParseSpace1().And(ParseType()).Fmap(func(typ *ParseResult) *ParseResult {
			v := []interface{}{variable.Result.(*expr.IR_Variable).Value, typ.Result.(shared.Type)}
			return ParseSuccess(v, typ.Rest)
		})
	})
	return ParseList(itemParser)
}

func ParseFunction() Parser {
	return ParseString("func").And(ParseSpace()).And(ParseByte('(')).And(ParseFunctionDefArgs()).AndThen(func(args *ParseResult) Parser {
		return ParseByte(')').And(ParseSpace()).And(ParseType()).AndThen(func(returns *ParseResult) Parser {
			return ParseBlock().Fmap(func(body *ParseResult) *ParseResult {
				argNames := []string{}
				argTypes := []shared.Type{}
				for _, pair := range args.Result.([]interface{}) {
					lst := pair.([]interface{})
					argNames = append(argNames, lst[0].(string))
					argTypes = append(argTypes, lst[1].(shared.Type))
				}
				signature := &shared.TFunction{
					ReturnType: returns.Result.(shared.Type),
					Args:       argTypes,
					ArgNames:   argNames,
				}
				return ParseSuccess(expr.NewIR_Function(signature, body.Result.(shared.IR)), body.Rest)
			})
		})
	})
}

func ParseFunctionDef() Parser {
	return ParseString("func").And(ParseSpace1()).And(ParseVariable()).AndThen(func(name *ParseResult) Parser {
		return ParseSpace().And(ParseByte('(')).And(ParseFunctionDefArgs()).AndThen(func(args *ParseResult) Parser {
			return ParseByte(')').And(ParseSpace()).And(ParseType()).AndThen(func(returns *ParseResult) Parser {
				return ParseBlock().Fmap(func(body *ParseResult) *ParseResult {
					argNames := []string{}
					argTypes := []shared.Type{}
					for _, pair := range args.Result.([]interface{}) {
						lst := pair.([]interface{})
						argNames = append(argNames, lst[0].(string))
						argTypes = append(argTypes, lst[1].(shared.Type))
					}
					signature := &shared.TFunction{
						ReturnType: returns.Result.(shared.Type),
						Args:       argTypes,
						ArgNames:   argNames,
					}
					f := expr.NewIR_Function(signature, body.Result.(shared.IR))
					return ParseSuccess(statements.NewIR_FunctionDef(name.Result.(*expr.IR_Variable).Value, f), body.Rest)
				})
			})
		})
	})
}

func ParseFunctionArgs() Parser {
	return ParseList(ParseExpression())
}

func ParseFunctionCall() Parser {
	return ParseIdent().AndThen(func(v *ParseResult) Parser {
		return ParseByte('(').And(ParseSpace()).And(ParseFunctionArgs()).AndThen(func(args *ParseResult) Parser {
			return ParseSpace().And(ParseByte(')')).Fmap(func(r *ParseResult) *ParseResult {
				args := InterfaceArrayToIRExpressionArray(args.Result)
				function := v.Result.(string)
				var result shared.IRExpression
				if function == "syscall" {
					result = expr.NewIR_Syscall(args[0], args[1:])
				} else {
					result = expr.NewIR_Call(function, args)
				}

				if ty := ParseType()(function); ty.Result != nil && ty.Error == nil {
					if len(args) == 1 {
						result = expr.NewIR_Cast(args[0], ty.Result.(shared.Type))
					} else {
						return ParseError(fmt.Errorf("Too many parameters for call to %v", function))
					}
				}
				return ParseSuccess(result, r.Rest)
			})
		})
	})
}

func ParseAndThen() Parser {
	return ParseSingleStatement().AndThen(func(a *ParseResult) Parser {
		return ParseSpace().And(OneOf([]Parser{
			ParseByte(';'),
			ParseByte('\n'),
		})).And(ParseSpace()).And(ParseStatement()).Fmap(func(a2 *ParseResult) *ParseResult {
			return ParseSuccess(statements.NewIR_AndThen(a.Result.(shared.IR), a2.Result.(shared.IR)), a2.Rest)
		})
	})
}

func ParseReturn() Parser {
	return ParseString("return").And(ParseSpace1()).And(ParseExpression()).Fmap(func(r *ParseResult) *ParseResult {
		return ParseSuccess(statements.NewIR_Return(r.Result.(shared.IRExpression)), r.Rest)
	})
}

func ParseBlock() Parser {
	return ParseSpace().And(ParseByte('{')).And(ParseSpace()).And(ParseStatement()).AndThen(func(stmt *ParseResult) Parser {
		return ParseSpace().And(ParseByte('}')).And(ParseSpace()).Fmap(func(b *ParseResult) *ParseResult {
			return ParseSuccess(stmt.Result, b.Rest)
		})
	})

}

func ParseIf() Parser {
	return ParseString("if").And(ParseSpace1()).And(ParseExpression()).AndThen(func(cond *ParseResult) Parser {
		return ParseBlock().AndThen(func(stmt1 *ParseResult) Parser {
			return ParseString("else").And(ParseBlock()).Fmap(func(stmt2 *ParseResult) *ParseResult {
				return ParseSuccess(statements.NewIR_If(cond.Result.(shared.IRExpression), stmt1.Result.(shared.IR), stmt2.Result.(shared.IR)), stmt2.Rest)
			})
		})
	})
}

func ParseWhile() Parser {
	return ParseString("while").And(ParseSpace1()).And(ParseExpression()).AndThen(func(cond *ParseResult) Parser {
		return ParseBlock().Fmap(func(stmt1 *ParseResult) *ParseResult {
			return ParseSuccess(statements.NewIR_While(cond.Result.(shared.IRExpression), stmt1.Result.(shared.IR)), stmt1.Rest)
		})
	})
}

func ParseStructType() Parser {
	typ := ParseVariable().AndThen(func(field *ParseResult) Parser {
		return ParseSpace1().And(ParseType()).Fmap(func(ty *ParseResult) *ParseResult {
			return ParseSuccess([]interface{}{field.Result, ty.Result}, ty.Rest)
		})
	})
	return ParseEnclosed(ParseByte('{').And(ParseWhiteSpace()), ParseListWithSeparator(typ, ParseByte('\n')).Fmap(func(fields *ParseResult) *ParseResult {
		result := &shared.TStruct{
			FieldTypes: []shared.Type{},
			Fields:     []string{},
		}
		for _, f := range fields.Result.([]interface{}) {
			field := f.([]interface{})[0].(*expr.IR_Variable).Value
			typ := f.([]interface{})[1].(shared.Type)
			result.Fields = append(result.Fields, field)
			result.FieldTypes = append(result.FieldTypes, typ)
		}
		return ParseSuccess(result, fields.Rest)
	}), ParseWhiteSpace().And(ParseByte('}')))
}

func ParseStruct() Parser {
	return ParseString("struct").And(ParseSpace()).And(ParseStructType()).AndThen(func(fields *ParseResult) Parser {
		return ParseEnclosed(ParseByte('{').And(ParseWhiteSpace()), ParseArrayItems().Fmap(func(items *ParseResult) *ParseResult {
			return ParseSuccess(expr.NewIR_Struct(fields.Result.(*shared.TStruct), items.Result.([]shared.IRExpression)), items.Rest)
		}), ParseWhiteSpace().And(ParseByte('}')))
	})
}

func ParseStructField() Parser {
	return ParseVariable().AndThen(func(v *ParseResult) Parser {
		return ParseByte('.').And(ParseVariable()).Fmap(func(field *ParseResult) *ParseResult {
			return ParseSuccess(expr.NewIR_StructField(v.Result.(shared.IRExpression), field.Result.(*expr.IR_Variable).Value), field.Rest)
		})
	})
}

func ParseArrayIndex() Parser {
	return OneOf([]Parser{
		ParseVariable(),
		ParseArray(),
		ParseEnclosedExpression(),
	}).AndThen(func(e *ParseResult) Parser {
		return ParseByte('[').And(ParseSpace()).And(ParseExpression()).AndThen(func(e2 *ParseResult) Parser {
			return ParseSpace().And(ParseByte(']')).Fmap(func(l *ParseResult) *ParseResult {
				return ParseSuccess(expr.NewIR_ArrayIndex(e.Result.(shared.IRExpression), e2.Result.(shared.IRExpression)), l.Rest)
			})
		})
	})
}

func ParseIR(str string) (shared.IR, error) {
	result := ParseStatement()(str)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.Rest != "" {
		return nil, fmt.Errorf("Failed to parse: %s at %d", str, len(str)-len(result.Rest))
	}
	if result.Result == nil {
		return nil, fmt.Errorf("Nil parse result %s at %d", str, len(str)-len(result.Rest))
	}
	return result.Result.(shared.IR), nil
}

func MustParseIR(str string) shared.IR {
	result, err := ParseIR(str)
	if err != nil {
		panic(err)
	}
	return result
}

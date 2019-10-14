package ir

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bspaans/jit/ir/expr"
	"github.com/bspaans/jit/ir/shared"
	"github.com/bspaans/jit/ir/statements"
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

func ParseVariable() Parser {
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
			return &ParseResult{
				Result: expr.NewIR_Variable(string(sub1.Result.(byte)) + string(chars)),
				Rest:   sub2.Rest,
				Error:  nil,
			}
		})
	})
}

func ParseOperator() Parser {
	return OneOf([]Parser{
		ParseUint64(),
		ParseVariable(),
	}).AndThen(func(op1 *ParseResult) Parser {
		return ParseSpace().And(OneOf([]Parser{
			ParseByte('+'),
		})).AndThen(func(op *ParseResult) Parser {
			return ParseSpace().And(ParseExpression()).Fmap(func(op2 *ParseResult) *ParseResult {
				if op.Result.(byte) == '+' {
					return ParseSuccess(expr.NewIR_Add(op1.Result.(shared.IRExpression), op2.Result.(shared.IRExpression)), op2.Rest)
				}
				return ParseError(errors.New("Unknown operator"))
			})
		})
	})
}

func ParseSpace() Parser {
	return ParseByte(' ').Many()
}

func ParseExpression() Parser {
	return OneOf([]Parser{
		ParseOperator(),
		ParseUint64(),
		ParseVariable(),
	})
}

func ParseStatement() Parser {
	return ParseSpace().And(OneOf([]Parser{
		ParseAndThen(),
		ParseAssigment(),
	}))
}

func ParseAssigment() Parser {
	return ParseVariable().AndThen(func(variable *ParseResult) Parser {
		return ParseSpace().And(ParseByte('=')).And(ParseSpace()).And(ParseExpression()).Fmap(func(value *ParseResult) *ParseResult {
			v := variable.Result.(*expr.IR_Variable).Value
			return ParseSuccess(statements.NewIR_Assignment(v, value.Result.(shared.IRExpression)), value.Rest)
		})
	})
}

func ParseAndThen() Parser {
	return OneOf([]Parser{
		ParseAssigment(),
	}).AndThen(func(a *ParseResult) Parser {
		return ParseSpace().And(OneOf([]Parser{
			ParseByte(';'),
			ParseByte('\n'),
		})).And(ParseSpace()).And(ParseStatement()).Fmap(func(a2 *ParseResult) *ParseResult {
			return ParseSuccess(statements.NewIR_AndThen(a.Result.(shared.IR), a2.Result.(shared.IR)), a2.Rest)
		})
	})
}

func init() {
	fmt.Println(ParseUint64()("123123"))
	fmt.Println(ParseVariable()("a123"))
	fmt.Println(ParseAssigment()("a123 = 1234"))
	fmt.Println(ParseAssigment()("a123 = 1234 + 123"))
	fmt.Println(ParseStatement()("a123 = 1234 + z + b"))
	fmt.Println(ParseStatement()("a123 = 1234 + z + b; z = a + 2"))
	fmt.Println(ParseStatement()("a123 = 1234 + z + b; z = a + 2"))
	fmt.Println(ParseStatement()("a123 = 1234 + z + b; z = a + 2; b = 3"))
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
		return nil, fmt.Errorf("Nil parse result", str, len(str)-len(result.Rest))
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

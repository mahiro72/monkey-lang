package parser_test

import (
	"fmt"
	"testing"

	"github.com/mahiro72/monkey-lang/ast"
	"github.com/mahiro72/monkey-lang/lexer"
	"github.com/mahiro72/monkey-lang/parser"
	testingHelper "github.com/mahiro72/monkey-lang/testing"
	"github.com/mahiro72/monkey-lang/token"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedStatements []ast.Statement
		expectedErrors     []string
	}{
		{
			name: "success: 数値の変数宣言",
			input: `
				let x = 5;
			`,
			expectedStatements: []ast.Statement{
				&ast.LetStatement{
					Token: token.Token{Type: token.LET, Literal: "let"},
					Name: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Value: "x",
					},
					Value: nil,
				},
			},
		},
		{
			name: "failure: let式に数値が入っている",
			input: `
				let 1 = x;
			`,
			expectedStatements: nil,
			expectedErrors:     []string{"expected next token to be IDENT, got INT instead", "no prefix parse function for = found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) == 0 {
				testingHelper.AssertEqual(t, tt.expectedStatements, program.Statements)
			} else {
				testingHelper.AssertEqual(t, tt.expectedErrors, p.Errors())
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedStatements []ast.Statement
		expectedErrors     []string
	}{
		{
			name: "success: 数値のreturn宣言",
			input: `
				return 10;
			`,
			expectedStatements: []ast.Statement{
				&ast.ReturnStatement{
					Token:       token.Token{Type: token.RETURN, Literal: "return"},
					ReturnValue: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) == 0 {
				testingHelper.AssertEqual(t, tt.expectedStatements, program.Statements)
			} else {
				testingHelper.AssertEqual(t, tt.expectedErrors, p.Errors())
			}
		})
	}
}

func TestIdentifierExpression(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedStatements []ast.Statement
		expectedErrors     []string
	}{
		{
			name: "success: 識別子の式",
			input: `
				hoge;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "hoge"},
					Expression: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "hoge"},
						Value: "hoge",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) == 0 {
				testingHelper.AssertEqual(t, tt.expectedStatements, program.Statements)
			} else {
				testingHelper.AssertEqual(t, tt.expectedErrors, p.Errors())
			}
		})
	}
}

func TestParsingPrefixExpression(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedStatements []ast.Statement
		expectedErrors     []string
	}{
		{
			name: "success: 前置演算式(否定)の式",
			input: `
				!hoge;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.BANG, Literal: "!"},
					Expression: &ast.PrefixExpression{
						Token:    token.Token{Type: token.BANG, Literal: "!"},
						Operator: "!",
						Right: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "hoge"},
							Value: "hoge",
						},
					},
				},
			},
		},
		{
			name: "success: 前置演算式(マイナス)の式",
			input: `
				-x;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.MINUS, Literal: "-"},
					Expression: &ast.PrefixExpression{
						Token:    token.Token{Type: token.MINUS, Literal: "-"},
						Operator: "-",
						Right: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "x"},
							Value: "x",
						},
					},
				},
			},
		},
		{
			name: "success: 識別子の式",
			input: `
				hoge;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "hoge"},
					Expression: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "hoge"},
						Value: "hoge",
					},
				},
			},
		},
		{
			name: "success: 数字の式",
			input: `
				5;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) == 0 {
				testingHelper.AssertEqual(t, tt.expectedStatements, program.Statements)
			} else {
				testingHelper.AssertEqual(t, tt.expectedErrors, p.Errors())
			}
		})
	}
}

func TestParsingInfixExpression(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedStatements []ast.Statement
		expectedString     string
		expectedErrors     []string
	}{
		{
			name: "success: 演算式1",
			input: `
				5 + 5;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token: token.Token{Type: token.PLUS, Literal: "+"},
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Operator: "+",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
			expectedString: "(5 + 5)",
		},
		{
			name: "success: 演算式2",
			input: `
				5 + 5 - 2;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token: token.Token{Type: token.MINUS, Literal: "-"},
						Left: &ast.InfixExpression{
							Token: token.Token{Type: token.PLUS, Literal: "+"},
							Left: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
							Operator: "+",
							Right: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
						Operator: "-",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "2"},
							Value: 2,
						},
					},
				},
			},
			expectedString: "((5 + 5) - 2)",
		},
		{
			name: "success: 演算式3",
			input: `
				5 + 5 - 2 + 1;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token: token.Token{Type: token.PLUS, Literal: "+"},
						Left: &ast.InfixExpression{
							Token: token.Token{Type: token.MINUS, Literal: "-"},
							Left: &ast.InfixExpression{
								Token: token.Token{Type: token.PLUS, Literal: "+"},
								Left: &ast.IntegerLiteral{
									Token: token.Token{Type: token.INT, Literal: "5"},
									Value: 5,
								},
								Operator: "+",
								Right: &ast.IntegerLiteral{
									Token: token.Token{Type: token.INT, Literal: "5"},
									Value: 5,
								},
							},
							Operator: "-",
							Right: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "2"},
								Value: 2,
							},
						},
						Operator: "+",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "1"},
							Value: 1,
						},
					},
				},
			},
			expectedString: "(((5 + 5) - 2) + 1)",
		},
		{
			name: "success: 演算式4",
			input: `
				4 * x + 1;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "4"},
					Expression: &ast.InfixExpression{
						Token: token.Token{Type: token.PLUS, Literal: "+"},
						Left: &ast.InfixExpression{
							Token: token.Token{Type: token.ASTERISK, Literal: "*"},
							Left: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "4"},
								Value: 4,
							},
							Operator: "*",
							Right: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
						},
						Operator: "+",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "1"},
							Value: 1,
						},
					},
				},
			},
			expectedString: "((4 * x) + 1)",
		},
		{
			name: "success: 演算式5",
			input: `
				1 + 4 * x;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "1"},
					Expression: &ast.InfixExpression{
						Token: token.Token{Type: token.PLUS, Literal: "+"},
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "1"},
							Value: 1,
						},
						Operator: "+",
						Right: &ast.InfixExpression{
							Token: token.Token{Type: token.ASTERISK, Literal: "*"},
							Left: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "4"},
								Value: 4,
							},
							Operator: "*",
							Right: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
						},
					},
				},
			},
			expectedString: "(1 + (4 * x))",
		},
		{
			name: "success: 演算式6",
			input: `
				x * 2 > y + 3;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "x"},
					Expression: &ast.InfixExpression{
						Token: token.Token{Type: token.GT, Literal: ">"},
						Left: &ast.InfixExpression{
							Token: token.Token{Type: token.ASTERISK, Literal: "*"},
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Operator: "*",
							Right: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "2"},
								Value: 2,
							},
						},
						Operator: ">",
						Right: &ast.InfixExpression{
							Token: token.Token{Type: token.PLUS, Literal: "+"},
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "y"},
								Value: "y",
							},
							Operator: "+",
							Right: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "3"},
								Value: 3,
							},
						},
					},
				},
			},
			expectedString: "((x * 2) > (y + 3))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) == 0 {
				testingHelper.AssertEqual(t, tt.expectedStatements, program.Statements)
				testingHelper.AssertEqual(t, tt.expectedString, fmt.Sprintf("%s", program.Statements[0]))
			} else {
				testingHelper.AssertEqual(t, tt.expectedErrors, p.Errors())
			}
		})
	}
}

func TestParsingBoolean(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedStatements []ast.Statement
		expectedErrors     []string
	}{
		{
			name: "success: bool値",
			input: `
				true;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.TRUE, Literal: "true"},
					Expression: &ast.Boolean{
						Token: token.Token{Type: token.TRUE, Literal: "true"},
						Value: true,
					},
				},
			},
		},
		{
			name: "success: 否定のbool値",
			input: `
				!true;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.BANG, Literal: "!"},
					Expression: &ast.PrefixExpression{
						Token:    token.Token{Type: token.BANG, Literal: "!"},
						Operator: "!",
						Right: &ast.Boolean{
							Token: token.Token{Type: token.TRUE, Literal: "true"},
							Value: true,
						},
					},
				},
			},
		},
		{
			name: "success: boolを用いた式",
			input: `
				isAdmin == false;
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "isAdmin"},
					Expression: &ast.InfixExpression{
						Token: token.Token{Type: token.EQ, Literal: "=="},
						Left: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "isAdmin"},
							Value: "isAdmin",
						},
						Operator: "==",
						Right: &ast.Boolean{
							Token: token.Token{Type: token.FALSE, Literal: "false"},
							Value: false,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) == 0 {
				testingHelper.AssertEqual(t, tt.expectedStatements, program.Statements)
			} else {
				testingHelper.AssertEqual(t, tt.expectedErrors, p.Errors())
			}
		})
	}
}

func TestParseGroupedExpression(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedStatements []ast.Statement
		expectedErrors     []string
	}{
		{
			name: "success: 演算式 (後方優先)",
			input: `
				2 + (1 + 1);
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "2"},
					Expression: &ast.InfixExpression{
						Token: token.Token{Type: token.PLUS, Literal: "+"},
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "2"},
							Value: 2,
						},
						Operator: "+",
						Right: &ast.InfixExpression{
							Token: token.Token{Type: token.PLUS, Literal: "+"},
							Left: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "1"},
								Value: 1,
							},
							Operator: "+",
							Right: &ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "1"},
								Value: 1,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) == 0 {
				testingHelper.AssertEqual(t, tt.expectedStatements, program.Statements)
			} else {
				testingHelper.AssertEqual(t, tt.expectedErrors, p.Errors())
			}
		})
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedStatements []ast.Statement
		expectedErrors     []string
	}{
		{
			name: "success: if式",
			input: `
				if (x < y) { x }
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IF, Literal: "if"},
					Expression: &ast.IfExpression{
						Token: token.Token{Type: token.IF, Literal: "if"},
						Condition: &ast.InfixExpression{
							Token: token.Token{Type: token.LT, Literal: "<"},
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Operator: "<",
							Right: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "y"},
								Value: "y",
							},
						},
						Consequence: &ast.BlockStatement{
							Token: token.Token{Type: token.LBRACE, Literal: "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Expression: &ast.Identifier{
										Token: token.Token{Type: token.IDENT, Literal: "x"},
										Value: "x",
									},
								},
							},
						},
						Alternative: nil,
					},
				},
			},
		},
		{
			name: "success: if-else式",
			input: `
				if (x < y) { x } else { y }
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IF, Literal: "if"},
					Expression: &ast.IfExpression{
						Token: token.Token{Type: token.IF, Literal: "if"},
						Condition: &ast.InfixExpression{
							Token: token.Token{Type: token.LT, Literal: "<"},
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Operator: "<",
							Right: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "y"},
								Value: "y",
							},
						},
						Consequence: &ast.BlockStatement{
							Token: token.Token{Type: token.LBRACE, Literal: "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Expression: &ast.Identifier{
										Token: token.Token{Type: token.IDENT, Literal: "x"},
										Value: "x",
									},
								},
							},
						},
						Alternative: &ast.BlockStatement{
							Token: token.Token{Type: token.LBRACE, Literal: "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.IDENT, Literal: "y"},
									Expression: &ast.Identifier{
										Token: token.Token{Type: token.IDENT, Literal: "y"},
										Value: "y",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) == 0 {
				testingHelper.AssertEqual(t, tt.expectedStatements, program.Statements)
			} else {
				testingHelper.AssertEqual(t, tt.expectedErrors, p.Errors())
			}
		})
	}
}

func TestParseFunctionLiteral(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedStatements []ast.Statement
		expectedErrors     []string
	}{
		{
			name: "success: 関数",
			input: `
				fn (x, y) { x + y; }
			`,
			expectedStatements: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
					Expression: &ast.FunctionLiteral{
						Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
						Parameters: []*ast.Identifier{
							{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							{
								Token: token.Token{Type: token.IDENT, Literal: "y"},
								Value: "y",
							},
						},
						Body: &ast.BlockStatement{
							Token: token.Token{Type: token.LBRACE, Literal: "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Expression: &ast.InfixExpression{
										Token: token.Token{Type: token.PLUS, Literal: "+"},
										Left: &ast.Identifier{
											Token: token.Token{Type: token.IDENT, Literal: "x"},
											Value: "x",
										},
										Operator: "+",
										Right: &ast.Identifier{
											Token: token.Token{Type: token.IDENT, Literal: "y"},
											Value: "y",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "success: 関数2",
			input: `
				let x = add(2, 3 * 4);
			`,
			expectedStatements: []ast.Statement{
				&ast.LetStatement{
					Token: token.Token{Type: token.LET, Literal: "let"},
					Name: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Value: "x",
					},
					Value: &ast.CallExpression{
						Token: token.Token{Type: token.LPAREN, Literal: "("},
						Function: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "add"},
							Value: "add",
						},
						Arguments: []ast.Expression{
							&ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "2"},
								Value: 2,
							},
							&ast.InfixExpression{
								Token: token.Token{Type: token.ASTERISK, Literal: "*"},
								Left: &ast.IntegerLiteral{
									Token: token.Token{Type: token.INT, Literal: "3"},
									Value: 3,
								},
								Operator: "*",
								Right: &ast.IntegerLiteral{
									Token: token.Token{Type: token.INT, Literal: "4"},
									Value: 4,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "success: 関数3",
			input: `
				let f = fn (x) { x * 2; };
				let ret = f(5);
			`,
			expectedStatements: []ast.Statement{
				&ast.LetStatement{
					Token: token.Token{Type: token.LET, Literal: "let"},
					Name: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "f"},
						Value: "f",
					},
					Value: &ast.FunctionLiteral{
						Token: token.Token{Type: token.FUNCTION, Literal: "fn"},
						Parameters: []*ast.Identifier{
							{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
						},
						Body: &ast.BlockStatement{
							Token: token.Token{Type: token.LBRACE, Literal: "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{Type: token.IDENT, Literal: "x"},
									Expression: &ast.InfixExpression{
										Token: token.Token{Type: token.ASTERISK, Literal: "*"},
										Left: &ast.Identifier{
											Token: token.Token{Type: token.IDENT, Literal: "x"},
											Value: "x",
										},
										Operator: "*",
										Right: &ast.IntegerLiteral{
											Token: token.Token{Type: token.INT, Literal: "2"},
											Value: 2,
										},
									},
								},
							},
						},
					},
				},
				&ast.LetStatement{
					Token: token.Token{Type: token.LET, Literal: "let"},
					Name: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "ret"},
						Value: "ret",
					},
					Value: &ast.CallExpression{
						Token: token.Token{Type: token.LPAREN, Literal: "("},
						Function: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "f"},
							Value: "f",
						},
						Arguments: []ast.Expression{
							&ast.IntegerLiteral{
								Token: token.Token{Type: token.INT, Literal: "5"},
								Value: 5,
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) == 0 {
				testingHelper.AssertEqual(t, tt.expectedStatements, program.Statements)
			} else {
				testingHelper.AssertEqual(t, tt.expectedErrors, p.Errors())
			}
		})
	}
}

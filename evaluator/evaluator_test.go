package evaluator_test

import (
	"testing"

	"github.com/mahiro72/monkey-lang/evaluator"
	"github.com/mahiro72/monkey-lang/lexer"
	"github.com/mahiro72/monkey-lang/object"
	"github.com/mahiro72/monkey-lang/parser"
	testingHelper "github.com/mahiro72/monkey-lang/testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedObj object.Object
	}{
		{
			name:        "success: 5",
			input:       `5`,
			expectedObj: &object.Integer{Value: 5},
		},
		{
			name:        "success: -5",
			input:       `-5`,
			expectedObj: &object.Integer{Value: -5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()
			obj := evaluator.Eval(program, env)

			testingHelper.AssertEqual(t, tt.expectedObj, obj)
		})
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedObj object.Object
	}{
		{
			name:        "success: true",
			input:       `true`,
			expectedObj: &object.Boolean{Value: true},
		},
		{
			name:        "success: !false",
			input:       `!false`,
			expectedObj: &object.Boolean{Value: true},
		},
		{
			name:        "success: !!false",
			input:       `!!false`,
			expectedObj: &object.Boolean{Value: false},
		},
		{
			name:        "success: !5",
			input:       `!5`,
			expectedObj: &object.Boolean{Value: false},
		},
		{
			name:        "success: true == true",
			input:       `true == true`,
			expectedObj: &object.Boolean{Value: true},
		},
		{
			name:        "success: true == false",
			input:       `true == false`,
			expectedObj: &object.Boolean{Value: false},
		},
		{
			name:        "success: (3 > 1) == true",
			input:       `(3 > 1) == true`,
			expectedObj: &object.Boolean{Value: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()
			obj := evaluator.Eval(program, env)

			testingHelper.AssertEqual(t, tt.expectedObj, obj)
		})
	}
}

func TestEvalIntegerInfixExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedObj object.Object
	}{
		{
			name:        "success: 5 + 5",
			input:       `5 + 5;`,
			expectedObj: &object.Integer{Value: 10},
		},
		{
			name:        "success: 5 - 5",
			input:       `5 - 5;`,
			expectedObj: &object.Integer{Value: 0},
		},
		{
			name:        "success: 5 * 5",
			input:       `5 * 5;`,
			expectedObj: &object.Integer{Value: 25},
		},
		{
			name:        "success: 5 / 5",
			input:       `5 / 5;`,
			expectedObj: &object.Integer{Value: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()
			obj := evaluator.Eval(program, env)

			testingHelper.AssertEqual(t, tt.expectedObj, obj)
		})
	}
}

func TestEvalIfExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedObj object.Object
	}{
		{
			name:        "success: if (true) { 10 }",
			input:       `if (true) { 10 };`,
			expectedObj: &object.Integer{Value: 10},
		},
		{
			name:        "success: if (false) { 10 }",
			input:       ` if (false) { 10 };`,
			expectedObj: &object.Null{},
		},
		{
			name:        "success: if (false) { 10 } else { 20 }",
			input:       `if (false) { 10 } else { 20 };`,
			expectedObj: &object.Integer{Value: 20},
		},
		{
			name:        "success: if (5 < 10) { 10 } else { 20 }",
			input:       `if (5 < 10) { 10 } else { 20 };`,
			expectedObj: &object.Integer{Value: 10},
		},
		{
			name: "success: ネストされたif",
			input: `
				if (true) {
					if (true) {
						return 10;
					}
					return 5;
				};
			`,
			expectedObj: &object.Integer{Value: 10},
		},
		{
			name: "success: ネストされたif2",
			input: `
				if (true) {
					if (false) {
						return 10;
					}
					return 5;
				};
			`,
			expectedObj: &object.Integer{Value: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()
			obj := evaluator.Eval(program, env)

			testingHelper.AssertEqual(t, tt.expectedObj, obj)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name                string
		input               string
		expectedErrorString string
	}{
		{
			name:                "success: 5 + true;",
			input:               `5 + true;`,
			expectedErrorString: "Error: type mismatch: INTEGER + BOOLEAN",
		},
		{
			name:                "success: 5 + true; 5;",
			input:               `5 + true; 5;`,
			expectedErrorString: "Error: type mismatch: INTEGER + BOOLEAN",
		},
		{
			name:                "success: -true",
			input:               `-true`,
			expectedErrorString: "Error: unknown operator: -BOOLEAN",
		},
		{
			name:                "success: true + false;",
			input:               `true + false;`,
			expectedErrorString: "Error: unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			name:                "success: 5; true + false; 5;",
			input:               `5; true + false; 5;`,
			expectedErrorString: "Error: unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			name:                "success: if ( 10 > 1 ) { true + false; };",
			input:               `if ( 10 > 1 ) { true + false; };`,
			expectedErrorString: "Error: unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			name:                "success: let x = 1; x + y;",
			input:               `let x = 1; x + y;`,
			expectedErrorString: "Error: identifier not found: y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()
			obj := evaluator.Eval(program, env)

			testingHelper.AssertEqual(t, tt.expectedErrorString, obj.Inspect())
		})
	}
}

func TestEvalLetStatement(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedObj object.Object
	}{
		{
			name:        "success: let a = 5; a;",
			input:       `let a = 5; a;`,
			expectedObj: &object.Integer{Value: 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			env := object.NewEnvironment()
			obj := evaluator.Eval(program, env)

			testingHelper.AssertEqual(t, tt.expectedObj, obj)
		})
	}
}

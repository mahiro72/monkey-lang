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
		name 	  string
		input    string
		expectedObj object.Object
	}{
		{
			name: "success: 5",
			input: `5`,
			expectedObj: &object.Integer{Value: 5},
		},
		{
			name: "success: -5",
			input: `-5`,
			expectedObj: &object.Integer{Value: -5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			obj := evaluator.Eval(program)

			testingHelper.AssertEqual(t, tt.expectedObj, obj)
		})
	}
}


func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		name 	  string
		input    string
		expectedObj object.Object
	}{
		{
			name: "success: true",
			input: `true`,
			expectedObj: &object.Boolean{Value: true},
		},
		{
			name: "success: !false",
			input: `!false`,
			expectedObj: &object.Boolean{Value: true},
		},
		{
			name: "success: !!false",
			input: `!!false`,
			expectedObj: &object.Boolean{Value: false},
		},
		{
			name: "success: !5",
			input: `!5`,
			expectedObj: &object.Boolean{Value: false},
		},
		{
			name: "success: true == true",
			input: `true == true`,
			expectedObj: &object.Boolean{Value: true},
		},
		{
			name: "success: true == false",
			input: `true == false`,
			expectedObj: &object.Boolean{Value: false},
		},
		{
			name: "success: (3 > 1) == true",
			input: `(3 > 1) == true`,
			expectedObj: &object.Boolean{Value: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			obj := evaluator.Eval(program)

			testingHelper.AssertEqual(t, tt.expectedObj, obj)
		})
	}
}

func TestEvalIntegerInfixExpression(t *testing.T) {
	tests := []struct {
		name 	  string
		input    string
		expectedObj object.Object
	}{
		{
			name: "success: 5 + 5",
			input: `5 + 5;`,
			expectedObj: &object.Integer{Value: 10},
		},
		{
			name: "success: 5 - 5",
			input: `5 - 5;`,
			expectedObj: &object.Integer{Value: 0},
		},
		{
			name: "success: 5 * 5",
			input: `5 * 5;`,
			expectedObj: &object.Integer{Value: 25},
		},
		{
			name: "success: 5 / 5",
			input: `5 / 5;`,
			expectedObj: &object.Integer{Value: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.ParseProgram()
			obj := evaluator.Eval(program)

			testingHelper.AssertEqual(t, tt.expectedObj, obj)
		})
	}
}



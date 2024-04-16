package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsNumOrDec(t *testing.T) {
	tests := []struct {
		name     string
		char     rune
		expected bool
	}{
		{name: "Zero", char: '0', expected: true},
		{name: "Decimal Point", char: '.', expected: true},
		{name: "Number", char: '5', expected: true},
		{name: "Wrong Symbol", char: 'a', expected: false},
		{name: "Plus Sign", char: '+', expected: false},
	}

	t.Parallel()
	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := isNumOrDec(tc.char)

			assert.Equal(t, tc.expected, actual, "Error in test: %s func: isNumOrDec. Char %q, got: %v, expected: %v", tc.name, tc.char, actual, tc.expected)
		})
	}
}

func TestProcessCurrentNum(t *testing.T) {
	tests := []struct {
		name        string
		currentNum  string
		expected    float64
		expectError bool
	}{
		{name: "Valid positive integer", currentNum: "10", expected: 10, expectError: false},
		{name: "Valid negative integer", currentNum: "-10", expected: -10, expectError: false},
		{name: "Valid float", currentNum: "10.5", expected: 10.5, expectError: false},
		{name: "Empty string", currentNum: "", expected: defaultValue, expectError: true},
		{name: "Invalid format", currentNum: "abbb", expected: defaultValue, expectError: true},
	}

	for _, test := range tests {
		tc := test
		t.Parallel()
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actualResult, err := processCurrentNum(tc.currentNum)

			if tc.expectError {
				assert.Error(t, err, "Expected an error for currentNum: %q", tc.currentNum)
			} else {
				assert.NoError(t, err, "Did not expect an error for currentNum: %q", tc.currentNum)
				assert.Equal(t, tc.expected, actualResult, "Unexpected result for currentNum: %q", tc.currentNum)
			}
		})
	}
}

func TestEvaluateExpressionsBasic(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		result     float64
	}{
		{name: "Addition 1", expression: "1+1", result: 2},
		{name: "Subtraction", expression: "10-4", result: 6},
		{name: "Multiplication", expression: "4*5", result: 20},
		{name: "Division", expression: "20/5", result: 4},
		{name: "Grouped operations", expression: "20/(2+3)", result: 4},
		{name: "Negative numbers", expression: "-1-1", result: -2},
	}

	for _, test := range tests {
		tc := test
		t.Parallel()
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actual, err := evaluateExpression(tc.expression)
			assert.NoError(t, err, "Error evaluating expression: %q", tc.expression)
			assert.Equal(t, tc.result, actual, "Unexpected result for expression: %q", tc.expression)
		})
	}
}

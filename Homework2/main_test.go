package main

import "testing"

func TestIsNumOrDec(t *testing.T) {
	tests := []struct {
		char     rune
		expected bool
	}{
		{'0', true},
		{'.', true},
		{'5', true},
		{'a', false},
		{'+', false},
	}

	for _, testCase := range tests {
		actual := isNumOrDec(testCase.char)
		if actual != testCase.expected {
			t.Errorf("Error in isNumOrDec. Char %q, got: %v, expected: %v", testCase.char, actual, testCase.expected)
		}
	}
}

func TestProcessCurrentNum(t *testing.T) {
	tests := []struct {
		currentNum  string
		expected    float64
		expectError bool
	}{
		{"10", 10, false},
		{"-10", -10, false},
		{"10.5", 10.5, false},
		{"", defaultValue, true},
		{"abbb", defaultValue, true},
	}

	for _, testCase := range tests {
		actualResult, err := processCurrentNum(testCase.currentNum)
		if testCase.expectError {
			if err == nil {
				t.Errorf("processCurrentNum expected error for currentNum: %q but there isn't any.", testCase.currentNum)
			}
			continue
		}
		if err != nil {
			t.Errorf("processCurrentNum error for currentNum: %q but there isn't any.", testCase.currentNum)
		}
		if actualResult != testCase.expected {
			t.Errorf("processCurrentNum for num: %q returned: %v, expected: %v.", testCase.currentNum, actualResult, testCase.expected)
		}
	}
}

func TestEvaluateExpressionsBasic(t *testing.T) {
	tests := []struct {
		expression string
		result     float64
	}{
		{"1+1", 2},
		{"2+2", 4},
		{"10-4", 6},
		{"4*5", 20},
		{"20/5", 4},
		{"20/(2+3)", 4},
		{"-1-1", -2},
	}

	for _, testCase := range tests {
		actual, err := evaluateExpression(testCase.expression)
		if err != nil {
			t.Errorf("Error in epression: %q, error: %v", testCase.expression, err)
		}
		if actual != testCase.result {
			t.Errorf("Error in calculation epression: %q, got: %v, expected: %v", testCase.expression, actual, testCase.result)
		}
	}
}

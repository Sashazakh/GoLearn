package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestErrorOpeningNonExistFile(t *testing.T) {
	_, err := openInputFile("non_exists.txt")
	if err == nil {
		t.Errorf("Expected an error for non-exist file")
	}
}

func TestOpeningExistFile(t *testing.T) {
	testFileName := "test_filename.txt"
	testFile, err := os.Create(testFileName)
	if err != nil {
		t.Errorf("Failed to create file")
	}
	testFile.Close()
	defer os.Remove(testFileName)

	_, err = openInputFile(testFileName)
	if err != nil {
		t.Errorf("Failed opening existing file")
	}
}

func TestProcessLine(t *testing.T) {
	testCases := []struct {
		line     string
		options  Options
		expected string
	}{
		{"Ignore Cases", Options{ignoreCase: true}, "ignore cases"},
		{"First Second", Options{ignoreFields: 1}, "Second"},
		{"123456", Options{ignoreChars: 2}, "3456"},
		{"IGNORE first", Options{ignoreFields: 1, ignoreCase: true}, "first"},
		{"Ignore first two chars", Options{ignoreFields: 1, ignoreChars: 2}, "rst two chars"},
	}

	for _, testCase := range testCases {
		result := processLine(testCase.line, testCase.options)

		if result != testCase.expected {
			t.Errorf("Process line: %q, options: %+v, expected: %q got: %q", testCase.line, testCase.options, testCase.expected, result)
		}
	}
}

func TestProcessFile(t *testing.T) {
	tests := []struct {
		input    string
		options  Options
		expected string
	}{
		{
			input:    "Line1\nLine2\nLine2\nLine3",
			options:  Options{uniqueFlag: true},
			expected: "Line1\nLine3\n",
		},
		{
			input:    "Line1\nLine2\nLine2\nLine3",
			options:  Options{duplicateFlag: true},
			expected: "Line2\n",
		},
		{
			input:    "Line1\nLine2\nLine2\nLine3",
			options:  Options{countFlag: true},
			expected: "1 Line1\n2 Line2\n1 Line3\n",
		},
	}

	for _, test := range tests {
		outputBuffer := new(bytes.Buffer)

		processFile(strings.NewReader(test.input), outputBuffer, test.options)
		actual := outputBuffer.String()
		if actual != test.expected {
			t.Errorf("Expected output %q, got %q", test.expected, actual)
		}
	}
}

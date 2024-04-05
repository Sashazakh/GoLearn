package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func TestErrorOpeningNonExistFile(t *testing.T) {
	_, err := openInputFile("non_exists.txt")
	assert.Error(t, err, "Expected an error when opening a non-existent file")
}

func TestOpeningExistFile(t *testing.T) {
	testFileName := "test_filename.txt"
	testFile, err := os.Create(testFileName)
	assert.NoError(t, err, "Failed to create file")
	testFile.Close()
	defer os.Remove(testFileName)

	file, err := openInputFile(testFileName)
	assert.NoError(t, err, "Failed opening existing file")
	assert.NotNil(t, file, "File should not be nil")
}

func TestProcessLine(t *testing.T) {
	testCases := []struct {
		name     string
		line     string
		options  Options
		expected string
	}{
		{"Ignore Cases", "Ignore Cases", Options{ignoreCase: true}, "ignore cases"},
		{"First Second", "First Second", Options{ignoreFields: 1}, "Second"},
		{"Ignore chars", "123456", Options{ignoreChars: 2}, "3456"},
		{"Ignore first", "IGNORE first", Options{ignoreFields: 1, ignoreCase: true}, "first"},
		{"Ignore first two chars", "Ignore first two chars", Options{ignoreFields: 1, ignoreChars: 2}, "rst two chars"},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Parallel()
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := processLine(tc.line, tc.options)
			assert.Equal(t, tc.expected, result, "Process line: %q, options: %+v", testCase.line, testCase.options)
		})
	}
}

func TestProcessFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  Options
		expected string
	}{
		{
			name:     "Unique flag enabled",
			input:    "Line1\nLine2\nLine2\nLine3",
			options:  Options{isUniqueFlagEnabled: true},
			expected: "Line1\nLine3\n",
		},
		{
			name:     "Duplicate flag enabled",
			input:    "Line1\nLine2\nLine2\nLine3",
			options:  Options{isDuplicateFlagEnabled: true},
			expected: "Line2\n",
		},
		{
			name:     "Count flag enabled",
			input:    "Line1\nLine2\nLine2\nLine3",
			options:  Options{countFlag: true},
			expected: "1 Line1\n2 Line2\n1 Line3\n",
		},
	}

	for _, test := range tests {
		tc := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			outputBuffer := new(bytes.Buffer)
			processFile(strings.NewReader(tc.input), outputBuffer, tc.options)
			actual := outputBuffer.String()

			assert.Equal(t, tc.expected, actual, "Expected output %q", tc.expected)
		})
	}
}

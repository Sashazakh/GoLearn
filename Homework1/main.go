package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type Options struct {
	countFlag              bool
	inputFile              string
	outputFile             string
	ignoreCase             bool
	isDuplicateFlagEnabled bool
	isUniqueFlagEnabled    bool
	ignoreFields           int
	ignoreChars            int
}

func (o *Options) parseArgs() {
	flag.BoolVar(&o.countFlag, "c", false, "Count occurrences")
	flag.StringVar(&o.inputFile, "input", "", "Input")
	flag.StringVar(&o.outputFile, "output", "", "Output")
	flag.BoolVar(&o.ignoreCase, "i", false, "Ignore case")
	flag.BoolVar(&o.isDuplicateFlagEnabled, "d", false, "Print duplicates")
	flag.BoolVar(&o.isUniqueFlagEnabled, "u", false, "Print unique")
	flag.IntVar(&o.ignoreFields, "f", 0, "Ignore fields")
	flag.IntVar(&o.ignoreChars, "s", 0, "Ignore chars")
	flag.Parse()
}

func openInputFile(path string) (*os.File, error) {
	if path != "" {
		return os.Open(path)
	}

	return os.Stdin, nil
}

func openOutputFile(path string) (*os.File, error) {
	if path != "" {
		return os.Create(path)
	}

	return os.Stdout, nil
}

func processFile(input io.Reader, output io.Writer, options Options) {
	scanner := bufio.NewScanner(input)
	var prevLine, originalPrevLine string
	lineCount := 0

	for scanner.Scan() {
		notProcessedLine := scanner.Text()
		line := processLine(scanner.Text(), options)

		if line == prevLine {
			lineCount++
		} else {
			printLine(output, originalPrevLine, lineCount, options)
			lineCount = 1
			originalPrevLine = notProcessedLine
		}

		prevLine = line
	}

	printLine(output, originalPrevLine, lineCount, options)
}

func processLine(line string, options Options) string {
	if options.ignoreCase {
		line = strings.ToLower(line)
	}

	if options.ignoreFields > 0 || options.ignoreChars > 0 {
		fields := strings.Fields(line)
		if len(fields) > options.ignoreFields {
			line = strings.Join(fields[options.ignoreFields:], " ")
		} else {
			line = ""
		}

		if len(line) > options.ignoreChars {
			line = line[options.ignoreChars:]
		}
	}

	return line
}

func printLine(output io.Writer, line string, count int, options Options) {
	if count == 1 && options.isUniqueFlagEnabled {
		fmt.Fprintln(output, line)
	} else if count > 0 && (!options.isUniqueFlagEnabled && !options.isDuplicateFlagEnabled) || (options.isDuplicateFlagEnabled && count > 1) {
		if options.countFlag {
			fmt.Fprintln(output, count, line)
		} else {
			fmt.Fprintln(output, line)
		}
	}
}

func main() {
	var options Options
	options.parseArgs()

	input, err := openInputFile(options.inputFile)
	if err != nil {
		fmt.Println("Error opening file!")
	}
	defer input.Close()

	output, err := openOutputFile(options.outputFile)
	if err != nil {
		fmt.Println("Error creating output file!")
	}
	defer output.Close()

	processFile(input, output, options)
}

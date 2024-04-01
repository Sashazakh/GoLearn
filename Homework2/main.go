package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const defaultValue = float64(0)

func evaluateExpression(expression string) (float64, error) {
	numsStack := Stack[float64]{}
	operationsStack := Stack[rune]{}
	var currentNum string
	var err error

	for _, char := range expression {
		if isNumOrDec(char) || (char == '-' && currentNum == "") {
			currentNum += string(char)
		} else {
			numsStack, err = processAndAppend(numsStack, currentNum)
			if err != nil {
				return defaultValue, err
			}

			numsStack, operationsStack, err = processOperator(char, numsStack, operationsStack)
			if err != nil {
				return defaultValue, err
			}

			currentNum = ""
		}
	}

	numsStack, err = processAndAppend(numsStack, currentNum)
	if err != nil {
		return defaultValue, err
	}

	for !operationsStack.isEmpty() {
		numsStack, operationsStack, err = applyOperation(numsStack, operationsStack)
		if err != nil {
			return defaultValue, err
		}
	}

	result := numsStack.pop()

	if !numsStack.isEmpty() {
		return defaultValue, fmt.Errorf("invalid expression")
	}

	return result, nil
}

func isNumOrDec(c rune) bool {
	return ('0' <= c && c <= '9') || c == '.'
}

func processAndAppend(numsStack Stack[float64], currentNum string) (Stack[float64], error) {
	if currentNum == "" {
		return numsStack, nil
	}
	num, err := processCurrentNum(currentNum)
	if err != nil {
		return numsStack, err
	}

	numsStack.push(num)
	return numsStack, nil
}

func processCurrentNum(currentNum string) (float64, error) {
	if currentNum == "" {
		return defaultValue, fmt.Errorf("no number to process")
	}
	num, err := strconv.ParseFloat(currentNum, 64)
	if err != nil {
		return defaultValue, err
	}
	return num, err
}

func processOperator(operator rune, numsStack Stack[float64], operationsStack Stack[rune]) (Stack[float64], Stack[rune], error) {
	var err error

	switch operator {
	case '(':
		operationsStack.push(operator)
	case ')':
		for operationsStack.len() > 0 {
			topOperation := operationsStack.peek()
			if topOperation == '(' {
				break
			}

			numsStack, operationsStack, err = applyOperation(numsStack, operationsStack)
			if err != nil {
				return EmptyStack[float64](), EmptyStack[rune](), err
			}
		}
		operationsStack.pop()
	default:
		for operationsStack.len() > 0 {
			topOperation := operationsStack.peek()
			if precedence(topOperation) < precedence(operator) {
				break
			}
			numsStack, operationsStack, err = applyOperation(numsStack, operationsStack)
			if err != nil {
				return EmptyStack[float64](), EmptyStack[rune](), err
			}
		}
		operationsStack.push(operator)
	}
	return numsStack, operationsStack, nil
}

func applyOperation(numsStack Stack[float64], operationsStack Stack[rune]) (Stack[float64], Stack[rune], error) {
	if numsStack.len() < 2 {
		return EmptyStack[float64](), EmptyStack[rune](), fmt.Errorf("not enought operands")
	}

	second := numsStack.pop()
	first := numsStack.pop()

	operation := operationsStack.pop()

	var result float64
	switch operation {
	case '+':
		result = first + second
	case '-':
		result = first - second
	case '*':
		result = first * second
	case '/':
		if second == 0 {
			return EmptyStack[float64](), EmptyStack[rune](), fmt.Errorf("division by zero")
		}
		result = first / second
	}

	numsStack.push(result)
	return numsStack, operationsStack, nil
}

func precedence(operation rune) int {
	switch operation {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	}
	return 0
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Input (or type 'end' to finish program):")

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error while reading!")
		}
		text = strings.ReplaceAll(text, " ", "")
		text = strings.Trim(text, "\n")

		if text == "end" {
			break
		}
		result, err := evaluateExpression(text)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Result:", result)
	}
}

func EmptyStack[T any]() Stack[T] {
	return Stack[T]{}
}

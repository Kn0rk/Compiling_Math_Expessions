package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

var global_bin []byte

func contains[T comparable](array []T, val T) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}

type Tokenizer struct {
	line        int
	charOffset  int
	fileContent []byte
	fileOffset  int
}

const (
	NumToken = iota
	IdToken
	PlusToken
	MinusToken
	MultiplicationToken
	DivisionToken
	UnknownToken
)

type Token struct {
	tokenType int
	text      string
	value     int
}

func (t *Tokenizer) advance() error {
	t.fileOffset += 1
	t.charOffset += 1
	var whitespaces = []byte(" \t")
	for t.fileOffset < len(t.fileContent) && contains(whitespaces, t.fileContent[t.fileOffset]) {
		t.charOffset += 1
		t.fileOffset += 1
	}
	var newline = []byte("\n")
	for t.fileOffset < len(t.fileContent) && contains(newline, t.fileContent[t.fileOffset]) {
		t.charOffset = 0
		t.fileOffset += 1
		t.line += 1
	}

	if t.fileOffset >= len(t.fileContent) {
		return io.EOF
	}

	return nil
}

func (t *Tokenizer) currentToken() Token {

	var resultToken = Token{}
	resultToken.text = string(t.fileContent[t.fileOffset])

	if t.fileContent[t.fileOffset] == '+' {
		resultToken.tokenType = PlusToken
		return resultToken
	}

	if t.fileContent[t.fileOffset] == '-' {
		resultToken.tokenType = MinusToken
		return resultToken
	}

	if t.fileContent[t.fileOffset] == '*' {
		resultToken.tokenType = MultiplicationToken
		return resultToken
	}

	if t.fileContent[t.fileOffset] == '/' {
		resultToken.tokenType = DivisionToken
		return resultToken
	}

	var number = []byte("0123456789")
	var numOffset = len(t.fileContent)
	if contains(number, t.fileContent[t.fileOffset]) {
		numOffset = t.fileOffset

		for t.fileOffset+1 < len(t.fileContent) && contains(number, t.fileContent[t.fileOffset+1]) {
			t.charOffset += 1
			t.fileOffset += 1
		}

		if numOffset <= t.fileOffset {
			resultToken.tokenType = NumToken
			value, err := strconv.Atoi(string(t.fileContent[numOffset : t.fileOffset+1]))
			if err != nil {
				panic(1)
			}
			resultToken.value = value
			resultToken.text = string(t.fileContent[numOffset : t.fileOffset+1])
			return resultToken
		}
	}

	resultToken.tokenType = UnknownToken
	return resultToken

}

func ParseNum(t *Tokenizer) error {
	token := t.currentToken()

	if token.tokenType != NumToken {
		return &SyntaxError{
			line:    t.line,
			offset:  t.charOffset,
			message: fmt.Sprintf("expected a number but got '%s'", token.text)}
	}
	// fmt.Printf("push %d \n", token.value)
	global_bin = append(global_bin, translateTerm(token.value)...)
	return nil
}

func ParseTerm(t *Tokenizer) error {
	parse_err := ParseNum(t)

	if parse_err != nil {
		return parse_err
	}
	err := t.advance()
	if err != nil {
		return err
	}

	for {
		token := t.currentToken()
		operator := -1
		switch token.tokenType {
		case MultiplicationToken:
			operator = MultiplicationToken
		case DivisionToken:
			operator = DivisionToken
		default:
			if parse_err == nil {
				return nil
			} else {
				return &SyntaxError{
					line:    t.line,
					offset:  t.charOffset,
					message: fmt.Sprintf("expected either a * or / sign but got %s", token.text),
				}
			}
		}

		err = t.advance()
		if err != nil {
			return &SyntaxError{
				line:    t.line,
				offset:  t.charOffset,
				message: "File ended before Term was completed.",
			}
		}
		err = ParseNum(t)
		if err != nil {
			return err
		}
		// fmt.Printf("* or /")
		global_bin = append(global_bin, translateOperation(operator)...)
		err = t.advance()
		if err != nil {
			return err // eof
		}
	}
}

func ParseExpr(t *Tokenizer) error {
	parse_err := ParseTerm(t)

	if parse_err != nil {
		return parse_err
	}

	for {
		token := t.currentToken()
		operator := -1
		switch token.tokenType {
		case PlusToken:
			operator = PlusToken
		case MinusToken:
			operator = MinusToken
		default:
			if parse_err == nil {
				return nil
			} else {
				return &SyntaxError{
					line:    t.line,
					offset:  t.charOffset,
					message: fmt.Sprintf("expected either a plus or minus sign but got %s", token.text),
				}
			}
		}

		err := t.advance()
		if err != nil {
			return &SyntaxError{
				line:    t.line,
				offset:  t.charOffset,
				message: "File ended before Expr was completed.",
			}
		}
		err = ParseTerm(t)
		// fmt.Printf("+/-\n")
		global_bin = append(global_bin, translateOperation(operator)...)
		if err != nil {
			return err // eof or err
		}
	}

}

func ParseStatement(t *Tokenizer) error {

	for {

		err := ParseExpr(t)
		// call print

		global_bin = append(global_bin, call_print()...)
		global_bin = append(global_bin, newLine()...)
		if err != nil {
			return err
		}
		// fmt.Println("")
		if err == io.EOF {
			return nil
		}
	}
}

func panic_on_err(e error) {
	if errors.Is(e, io.EOF) {
		// ignore
		return
	}
	if e != nil {
		panic(e)
	}
}

func parseFile(filePath string) {
	global_bin = make([]byte, 0)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No file named %v\n", filePath)
		panic(1)
	}
	tok := Tokenizer{fileContent: data}
	err = ParseStatement(&tok)
	panic_on_err(err)

	file, err := os.Create("go_elf")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	defer file.Close()

	// Write ELF header

	// global_bin = append(global_bin, addPrintResult()...)
	var bytes = createBinary(
		global_bin,
		make([]byte, 0),
	)
	// Write program header

	// Write machine code
	file.Write(bytes)

}

func main() {
	parseFile("../../inputs/input1.txt")
	fmt.Println("ELF file 'minimal_elf' created.")

}

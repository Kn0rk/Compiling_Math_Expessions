package reverse_polish

import (
	"io"
	"strconv"
)

func contains[T comparable](array []T, val T) bool {
	for _, item := range array {
		if item == val {
			return true
		}
	}
	return false
}

type lexer struct {
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

type token struct {
	name   int
	lexeme string
	value  int
}

func (t *lexer) advanceToNextLine() error {
	var newline = []byte("\n")
	for t.fileOffset < len(t.fileContent) && !contains(newline, t.fileContent[t.fileOffset]) {
		t.fileOffset += 1
	}
	t.fileOffset += 1
	t.charOffset = 0
	t.line += 1
	if t.fileOffset >= len(t.fileContent) {
		return io.EOF
	}

	return nil

}

func (t *lexer) advance() error {
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

func (t *lexer) currentToken() token {

	var resultToken = token{}
	resultToken.lexeme = string(t.fileContent[t.fileOffset])

	if t.fileContent[t.fileOffset] == '+' {
		resultToken.name = PlusToken
		return resultToken
	}

	if t.fileContent[t.fileOffset] == '-' {
		resultToken.name = MinusToken
		return resultToken
	}

	if t.fileContent[t.fileOffset] == '*' {
		resultToken.name = MultiplicationToken
		return resultToken
	}

	if t.fileContent[t.fileOffset] == '/' {
		resultToken.name = DivisionToken
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
			resultToken.name = NumToken
			value, err := strconv.Atoi(string(t.fileContent[numOffset : t.fileOffset+1]))
			if err != nil {
				panic(1)
			}
			resultToken.value = value
			resultToken.lexeme = string(t.fileContent[numOffset : t.fileOffset+1])
			return resultToken
		}
	}

	resultToken.name = UnknownToken
	return resultToken

}

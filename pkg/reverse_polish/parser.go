package reverse_polish

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var global_bin []byte

func parseNum(t *Tokenizer) error {
	token := t.currentToken()

	if token.name != NumToken {
		return &SyntaxError{
			line:    t.line,
			offset:  t.charOffset,
			message: fmt.Sprintf("expected a number but got '%s'", token.lexeme)}
	}
	// fmt.Printf("push %d \n", token.value)
	global_bin = append(global_bin, translateTerm(token.value).assembly...)
	return nil
}

func parseTerm(t *Tokenizer) error {
	parse_err := parseNum(t)

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
		switch token.name {
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
					message: fmt.Sprintf("expected either a * or / sign but got %s", token.lexeme),
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
		err = parseNum(t)
		if err != nil {
			return err
		}
		global_bin = append(global_bin, translateOperation(operator)...)
		err = t.advance()
		if err != nil {
			return err // eof
		}
	}
}

func parseExpr(t *Tokenizer) error {
	parse_err := parseTerm(t)

	if parse_err != nil {
		return parse_err
	}

	for {
		token := t.currentToken()
		operator := -1
		switch token.name {
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
					message: fmt.Sprintf("expected either a plus or minus sign but got %s", token.lexeme),
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
		err = parseTerm(t)
		global_bin = append(global_bin, translateOperation(operator)...)
		if err != nil {
			return err // eof or err
		}

	}

}

func parseStatement(t *Tokenizer) []SyntaxError {

	var syntaxErrors []SyntaxError = make([]SyntaxError, 0)

	for {
		var buff = global_bin
		global_bin = make([]byte, 0)
		err := parseExpr(t)
		// call print

		global_bin = append(global_bin, call_print()...)
		global_bin = append(global_bin, newLine()...)
		var e *SyntaxError
		if err == nil {
			global_bin = append(buff, global_bin...)
		}
		if errors.As(err, &e) {
			syntaxErrors = append(syntaxErrors, *e)
			err = t.advanceToNextLine()
		}
		if err == io.EOF {
			break
		}
	}
	return syntaxErrors
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

func ParseFile(filePath string) {
	global_bin = make([]byte, 0)
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No file named %v\n", filePath)
		panic(1)
	}
	tok := Tokenizer{fileContent: data}
	syntaxErrors := parseStatement(&tok)
	fmt.Print(syntaxErrors)
	panic_on_err(err)

	file, err := os.Create("go_elf")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	defer file.Close()

	var bytes = createBinary(
		global_bin,
		make([]byte, 0),
	)

	file.Write(bytes)

}

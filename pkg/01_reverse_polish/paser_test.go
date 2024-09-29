package reverse_polish

import (
	"testing"
)

func TestAddition(t *testing.T) {
	str := "daha"
	tok := lexer{fileContent: []byte(str)}
	err := parseStatement(&tok)
	if len(err) == 0 {
		t.Fail()
	}

	tok = lexer{fileContent: []byte("9+6")}
	err = parseStatement(&tok)
	if len(err) != 0 {
		t.Fail()
	}

	tok = lexer{fileContent: []byte("9+	6 -7")}
	err = parseStatement(&tok)
	if len(err) != 0 {
		t.Fail()
	}

	str = "9+ a	6 -7"
	tok = lexer{fileContent: []byte(str)}
	err = parseStatement(&tok)
	if len(err) == 0 {
		t.Fatalf("didnt fail: %s", str)
	}

	tok = lexer{fileContent: []byte("a9+	6 -7")}
	err = parseStatement(&tok)
	if len(err) == 0 {
		t.Fail()
	}

}

// func TestMultiplication(t *testing.T) {
// 	str := "daha"
// 	tok := Tokenizer{fileContent: []byte(str)}
// 	err := parseStatement(&tok)
// 	if _, ok := err.(*SyntaxError); !ok {
// 		t.Fail()
// 	}

// 	tok = Tokenizer{fileContent: []byte("9+6*6")}
// 	err = parseStatement(&tok)
// 	if _, ok := err.(*SyntaxError); ok {
// 		t.Fail()
// 	}

// 	tok = Tokenizer{fileContent: []byte("9*6 -7")}
// 	err = parseStatement(&tok)
// 	if _, ok := err.(*SyntaxError); ok {
// 		t.Fail()
// 	}

// 	str = "9+ a	6 *7"
// 	tok = Tokenizer{fileContent: []byte(str)}
// 	err = parseStatement(&tok)
// 	if _, ok := err.(*SyntaxError); !ok {
// 		t.Fatalf("didnt fail: %s", str)
// 	}

// 	tok = Tokenizer{fileContent: []byte("a9+	6 -7")}
// 	err = parseStatement(&tok)
// 	if _, ok := err.(*SyntaxError); !ok {
// 		t.Fail()
// 	}

// }

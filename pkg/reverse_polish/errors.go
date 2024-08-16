package reverse_polish

import "fmt"

type SyntaxError struct {
	line    int
	offset  int
	message string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("Line %d:%d - SyntaxError - %s", e.line, e.offset, e.message)
}

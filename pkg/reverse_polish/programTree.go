package reverse_polish

type ProgramNode struct {
	children []ProgramNode
	value    string
	line     int32
	offset   int32
	assembly []byte
}

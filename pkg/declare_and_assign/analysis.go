package declareandassign

import (
	"fmt"
)

func (statement *Statement) analyze() error {
	if statement.Declaration != nil {
		err := statement.Declaration.analyze()
		if err != nil {
			return err
		}
	}

	if statement.Assignment != nil {
		// statement.Assignment.
	}
	return nil

}

type DataType int

const (
	Int8 DataType = iota
	Int16
	Int32
	Int64
)

type Variable struct {
	name     string
	dataType DataType
}

var identifierMap = map[string]Variable{}

func (decl *Declaration) analyze() error {

	_, alreadyExists := identifierMap[*decl.Identifier]
	if alreadyExists {
		return fmt.Errorf("%v Line %v:%v Variable '%v' has already been declared.", decl.Pos.Filename, decl.Pos.Line, decl.Pos.Column, *decl.Identifier)
	}

	// get dataType from analyzing expr
	identifierMap[*decl.Identifier] = Variable{
		name:     *decl.Identifier,
		dataType: Int32,
	}

	return nil

}

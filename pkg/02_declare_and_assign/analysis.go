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
		return statement.Assignment.analyze()
	}
	return nil

}

type DataType int

var TypeNames = []string{"Int64", "Float64", "String"}

const (
	Int64 DataType = iota
	Float64
	String
)

type Variable struct {
	name     string
	dataType DataType
}

var identifierMap = map[string]Variable{}

func (decl *Declaration) analyze() error {

	_, alreadyExists := identifierMap[*decl.Identifier]
	if alreadyExists {
		return fmt.Errorf("%v Line %v:%v Variable '%v' has already been declared", decl.Pos.Filename, decl.Pos.Line, decl.Pos.Column, *decl.Identifier)
	}

	// get dataType from analyzing expr
	identifierMap[*decl.Identifier] = Variable{
		name:     *decl.Identifier,
		dataType: Int64,
	}
	err := decl.Expression.analyze()
	if err != nil {
		return err
	}
	return nil

}

func (assignment *Assignment) analyze() error {
	_, alreadyExists := identifierMap[*assignment.Identifier]
	if !alreadyExists {
		return fmt.Errorf("%v Line %v:%v Variable '%v' has not yet been declared", assignment.Pos.Filename, assignment.Pos.Line, assignment.Pos.Column, *assignment.Identifier)
	}
	err := assignment.Expression.analyze()
	if err != nil {
		return err
	}
	return nil
}

func (expr *Expression) analyze() error {
	err := expr.Left.Left.Base.analyze()
	if err != nil {
		return err
	}
	var typeOfExpression = expr.Left.Left.Base.Type
	for _, opTerm := range expr.Right {
		opTerm.Term.Left.Base.analyze()
		if opTerm.Term.Left.Base.Type != typeOfExpression {
			return fmt.Errorf("%v Line %v:%v Type mismatch between %v and %v",
				expr.Pos.Filename, expr.Pos.Line, expr.Pos.Column, TypeNames[typeOfExpression], TypeNames[*&opTerm.Term.Left.Base.Type])
		}
		for _, OpFactor := range opTerm.Term.Right {
			OpFactor.Factor.Base.analyze()
			if OpFactor.Factor.Base.Type != typeOfExpression {
				return fmt.Errorf("%v Line %v:%v Type mismatch between %v and %v",
					expr.Pos.Filename, expr.Pos.Line, expr.Pos.Column, TypeNames[typeOfExpression], TypeNames[*&opTerm.Term.Left.Base.Type])
			}
		}
	}
	return nil
}

func (value *Value) analyze() error {
	if value.Float != nil {
		value.Type = Float64
	}
	if value.Value != nil {
		value.Type = Int64
	}

	if value.String != nil {
		value.Type = String
	}
	return nil

}

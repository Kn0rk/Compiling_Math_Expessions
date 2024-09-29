package declareandassign

import (
	"fmt"
	// "strconv"

	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// var reserved =  lexer.SimpleRule{Name: "Res"}
var lexRules = []lexer.SimpleRule{
	{`Reserved`, `(printInt|printFloat|return|for|if|var)`},
	{`Operator`, `(\+|-|\*|/)`},
	{`Ident`, `[a-zA-Z][a-zA-Z_\d]*`},
	{`AssignmentOperator`, `(=|\+=)`},
	{"EOL", `;`},
	{`String`, `"(?:\\.|[^"])*"`},
	{`Float`, `\d+(?:\.\d+)`},
	{`Int`, `\d+\d*`},
	{`Char`, `'.'`},
	{`Punct`, `[][=]`},
	{"comment", `[#;][^\n]*`},
	{"whitespace", `\s+`},
}
var (
	iniLexer = lexer.MustSimple(lexRules)
	parser   = participle.MustBuild[KnorkLang](
		participle.Lexer(iniLexer),
	// participle.Unquote("String"),
	// participle.Union[Value](String{}, Number{}),
	)
)

type Operator int

const (
	OpMul Operator = iota
	OpDiv
	OpAdd
	OpSub
)

var operatorMap = map[string]Operator{"+": OpAdd, "-": OpSub, "*": OpMul, "/": OpDiv}

func (o *Operator) Capture(s []string) error {
	*o = operatorMap[s[0]]
	return nil
}

type KnorkLang struct {
	Statements []*Statement `@@*`
}

type Value struct {
	Value  *int64   `@(Int)`
	Float  *float64 `|@(Float)`
	String *string  `|@(String)`
	Type   DataType
	// Variable      *string     `| @Ident`
	// Subexpression *Expression `| "(" @@ ")"`
}

type Factor struct {
	Base *Value `@@`
	// Exponent *Value `( "^" @@ )?`
}

type OpFactor struct {
	Operator Operator `@("*" | "/")`
	Factor   *Factor  `@@`
}

type Term struct {
	Left  *Factor     `@@`
	Right []*OpFactor `@@*`
}

type OpTerm struct {
	Operator Operator `@("+" | "-")`
	Term     *Term    `@@`
}

type Expression struct {
	Pos   lexer.Position
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
	EOL   *string   `@(";")`
}

type Declaration struct {
	Pos        lexer.Position
	Prefix     *string             `@("var")`
	Identifier *string             `@Ident`
	Assignment *AssignmentOperator `@("=")`
	Expression *Expression         `@@`
}

type Assignment struct {
	Pos        lexer.Position
	Identifier *string             `@Ident`
	Assignment *AssignmentOperator `@("=")`
	Expression *Expression         `@@`
}

type AssignmentOperator int

const (
	Equals AssignmentOperator = iota
	PlusEquals
)

var assignmentMap = map[string]AssignmentOperator{"=": AssignmentOperator(Equals), "+=": AssignmentOperator(PlusEquals)}

func (o *AssignmentOperator) Capture(s []string) error {
	*o = assignmentMap[s[0]]
	return nil
}

type Statement struct {
	Declaration *Declaration `@@`
	Assignment  *Assignment  `| @@`
	Expression  *Expression  `| @@`
}

func DMain() {
	ini, err := parser.ParseString("", `var a = 4.2 +9;a=2;`)

	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		uerr, ok := err.(*participle.UnexpectedTokenError)
		if ok && uerr.Unexpected.Type == -2 {
			fmt.Printf("Line %v:%v Used reserved keyword: %v\n",
				uerr.Unexpected.Pos.Line,
				uerr.Unexpected.Pos.Column,
				uerr.Unexpected.Value)
		} else {
			panic(err)
		}
	}

	for _, statement := range ini.Statements {
		err = statement.analyze()
		if err != nil {
			panic(err)
		}
	}
}

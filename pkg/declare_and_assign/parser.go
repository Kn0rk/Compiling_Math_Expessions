package declareandassign

import (
	"fmt"

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
	Value         *int64   `@(Int)`
	Float         *float64 `|@(Float)`
	Type          int
	Variable      *string     `| @Ident`
	Subexpression *Expression `| "(" @@ ")"`
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
	Left  *Term     `@@`
	Right []*OpTerm `@@*`
	EOL   *string   `@(";")`
}

type Declaration struct {
	Prefix     *string             `@("var")`
	Identifier *string             `@Ident`
	Assignment *AssignmentOperator `@("=")`
	Expression *Expression         `@@`
}

type Assignment struct {
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

// var parser = participle.MustBuild[KnorkLang]()

func DMain() {
	// var reader = strings.NewReader(`"hi"+7+2`)
	ini, err := parser.ParseString("", `var far =  6.2 + 2;
	id = id +9;
	`)

	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		uerr, ok := err.(*participle.UnexpectedTokenError)
		// fmt.Print(uerr)
		if ok {
			// uerr.Unexpected.Type
			if uerr.Unexpected.Type == -2 {
				fmt.Printf("Line %v:%v Used reserved keyword: %v\n",
					uerr.Unexpected.Pos.Line,
					uerr.Unexpected.Pos.Column,
					uerr.Unexpected.Value)
			} else {
				panic(err)
			}
			panic(0)
		}

		panic(0)
	}
}

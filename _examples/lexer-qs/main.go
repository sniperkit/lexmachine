package main

import (
	"fmt"
	"log"
	"strings"

	cregex "github.com/mingrammer/commonregex"
	// pp "github.com/sniperkit/colly/plugins/app/debug/pp"

	lexmachine "github.com/sniperkit/lexmachine/pkg"
	machines "github.com/sniperkit/lexmachine/pkg/machines"
)

var Literals []string // The tokens representing literal strings
// var Keywords []string       // The keyword tokens
// var Tokens []string         // All of the tokens (including literals and keywords)
var TokenIds map[string]int // A map from the token names to their int ids
var Lexer *lexmachine.Lexer // The lexer object. Use this to construct a Scanner

// Called at package initialization. Creates the lexer and populates token lists.
func init() {
	initTokens()
	var err error
	Lexer, err = initLexer()
	if err != nil {
		panic(err)
	}
}

var (
	crGitRepo = cregex.GitRepoPattern

	// queries
	queryString = []string{
		`col=[1], row=[1]`,
		`cols=(:5), rows=(1:7)`,
		`cols[0:5], rows[1:7]`,
		`rows=(1,10), cols=("id", "name", "full_name", "description", "language", "stargazers_count", "forks_count")`,
		`cells=(A1,B2), row=(1,2), cols=(A,B), rows=[1:7], cols=["id", "name", "full_name", "description", "language", "stargazers_count", "forks_count"]`,
		//
		`block=(A1,B2), row=(1,2), cols=(A,B), prefix_path="./shared/dump.txt" to_file=true format="json" rows=[1:7], cols=["id", "name", "full_name", "description", "language", "stargazers_count", "forks_count"]`,
	}
)

func main() {

	qs := strings.ToLower(queryString[5])
	qs = strings.Replace(qs, " ", "", -1)

	s, err := Lexer.Scanner([]byte(qs))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Type    | Lexeme     | Position")
	fmt.Println("--------+------------+------------")

	for tok, err, eof := s.Next(); !eof; tok, err, eof = s.Next() {
		if ui, is := err.(*machines.UnconsumedInput); is {
			// to skip bad token do:
			s.TC = ui.FailTC
			log.Fatalln("error occured as unconsumed input, msg=", err) // however, we will just fail the program
		} else if err != nil {
			log.Fatalln("error occured as bad token, msg=", err) // however, we will just fail the program
			// log.Fatal(err)
		}
		token := tok.(*lexmachine.Token)
		fmt.Printf("%-7v | %-10v | %v:%v-%v:%v\n",
			Tokens[token.Type],
			string(token.Lexeme),
			token.StartLine,
			token.StartColumn,
			token.EndLine,
			token.EndColumn,
		)
	}

}

func initTokens() {
	Literals = []string{
		"[",
		"]",
		"(",
		")",
		"{",
		"}",
		"=",
		",",
		"...",
		";",
		":",
		// "'",
		// "\"",
		// "->",
		// "--",
	}

	Tokens = append(Tokens, Keywords...)
	Tokens = append(Tokens, Literals...)
	TokenIds = make(map[string]int)
	for i, tok := range Tokens {
		TokenIds[tok] = i
	}
}

var (
	// Rules

	// `SELECT_BLOCK` speficies regex to extract selection a block of cells. eg. `A2:B4`
	SELECT_BLOCK = `([A-Za-z]+[0-9]+)\:([A-Za-z]+[0-9]+)`

	// `SELECT_CELL` speficies regex to extract selection a specific cell. eg. `A2`
	SELECT_CELL = `([A-Za-z]+)([0-9]+)`

	// `SELECT_COLS_RANGE` speficies regex to extract selection of columns. eg `A:B`
	SELECT_COLS_RANGE = `([A-Za-z]+)\:([A-Za-z]+)`

	// SELECT_COLS_LIST
	SELECT_COLS_LIST = `(([a-zA-Z0-9-_\"\']+(,[a-zA-Z0-9-_\"\']+)*))`

	// `SELECT_ROWS` speficies regex to extract selection of rows. eg `2:5`
	SELECT_ROWS = `([0-9]+)\:([0-9]+)`

	// `SELECT_COL` speficies regex to extract selection of specific column. eg `A`
	SELECT_COL = `([A-Za-z]+)`

	// `SELECT_ROW` speficies regex to extract selection of specific row. eg `5`
	SELECT_ROW = `([0-9]+)`

	// `SELECT_NUMERIC` speficies regex to extract selection of columns and rows matched by their key index.
	// eg. `cols[0:5],rows[1:7]`, `cols[:],rows[:]`, `cols[1,2],rows[:]`
	SELECT_NUMERIC = `((col|cols|rows|row))\[((:\d+)|(\d+\:)|(\d+\:\d+)|(\:)|(\d+(,\d+)*))(\])`

	// `SELECT_ALPHA_NUMERIC` speficies regex to extract selection of columns and rows matched by their key index or column name.
	// Examples:
	// - `col["name", "full_name"], rows[1,5,7,8]`,
	// - `rows[1,10], col["id", "name", "full_name", "description", "language", "stargazers_count", "forks_count"]`
	SELECT_ALPHA_NUMERIC = `((col|cols|rows|row))\s\[((:[a-zA-Z0-9-_\"\']+)|([a-zA-Z0-9-_\"\']+\:)|([a-zA-Z0-9-_\"\']+\:[a-zA-Z0-9-_\"\']+)|([a-zA-Z0-9-_\"\']+(,[a-zA-Z0-9-_\"\']+)*)|(\:)|([a-zA-Z0-9-_\"\']+(,[a-zA-Z0-9-_\"\']+)*))(\])`

	// `SELECT_FUNCTION` speficies regex to extract selection of columns and rows matched by arguments.
	// Examples:
	// - `SELECT(col=1, row=2)`
	// - `SELECT(cols=[:], rows=[::5])`
	// - `SELECT(cols=["name","full_name"], rows=[::100])`
	// - `SELECT(cols=[1,2,7,10], row=2)`
	// - `SELECT(cols=[1,2,7,10], row=2)`
	// - `SELECT(cols=[1,2,7,10], rows=[2,4,7])`
	// - `SELECT(cols=[:], rows=[2,4,7])`
	// - `SELECT(cols=["name","full_name"], rows=[2,4,7])`
	// - `SELECT(cols=[0:5], rows=[1:7])`,
	// - `SELECT(cols=[:],rows=[:])`,
	// - `SELECT(cols=[1,2],rows=[:])`
	// - `SELECT(cols=["name","full_name"], rows=[1,5,7,8])`,
	// - `SELECT(rows=[1,10], cols=["id", "name", "full_name", "description", "language", "stargazers_count", "forks_count"])`
	SELECT_FUNCTION = `((col|cols|rows|row))\s\=\s\(((:\d+)|(\d+\:)|(\d+\:\d+)|(\:)|(\d+(,\d+)*))(\s\))`

	Keywords = []string{
		"prefix_path",
		"to_file",
		"databook",
		"dataset",
		"datasets",
		"col",
		"cols",
		"row",
		"rows",
		"slice",
		"sql",
		"format",
		"formats",
		"cells",
	}

	Tokens = []string{
		"ID",
		"DATASET",
		"DATASET_DIR",
		"DATASET_FILEPATH",
		"DATASET_INTO_FILE",
		"BLOCK_ALPHA",
		"BLOCK_FUNC",
		"BLOCK_NUM",
		"EXCEL_BLOCK",
		"EXCEL_CELL",
		"EXCEL_COLS",
		"EXCEL_ROWS",
		"CELLS_LIST",
		"COLS_LIST",
		"COLS_RANGE",
		"ROW",
		"LINK",
		"GIT_REPO",
		"COL",
		"COMMENT",
		"LOWER",
		"UPPER",
		"CAP",
		"ACTION",
		"AXIS",
		"AXIS_X",
		"AXIS_Y",
		"RANGE",
		"ARRAY",
		"LIST",
		"FORMAT",
	}
)

// Creates the lexer object and compiles the NFA.
func initLexer() (*lexmachine.Lexer, error) {
	lexer := lexmachine.NewLexer()

	for _, lit := range Literals {
		r := "\\" + strings.Join(strings.Split(lit, ""), "\\")
		lexer.Add([]byte(r), token(lit))
	}
	for _, name := range Keywords {
		lexer.Add([]byte(strings.ToLower(name)), token(name))
	}

	// lexer.Add([]byte("( |\t|\n|\r)+"), skip)
	// lexer.Add([]byte(`(\(|\)|\[|\])+`), skip)
	// lexer.Add([]byte(`=`), skip)
	// lexer.Add([]byte("\""), skip)
	// lexer.Add([]byte("\""), skip)

	lexer.Add([]byte(`[0-9]*\:[0-9]+`), token("RANGE"))

	lexer.Add([]byte(`((format|formats))\[((:[a-zA-Z0-9-_\"\']+)|([a-zA-Z0-9-_\"\']+(,[a-zA-Z0-9-_\"\']+)*))(\])+`), token("FORMAT"))
	lexer.Add([]byte(`((format|formats))\(((:[a-zA-Z0-9-_\"\']+)|([a-zA-Z0-9-_\"\']+(,[a-zA-Z0-9-_\"\']+)*))(\))+`), token("FORMAT"))

	lexer.Add([]byte(`([a-zA-Z0-9-_\.\:]*)`), token("DATASET_DIR"))

	lexer.Add([]byte(`(col|cols)`), token("AXIS_X"))
	lexer.Add([]byte(`(rows|row)`), token("AXIS_Y"))

	// lexer.Add([]byte(SELECT_BLOCK), token("EXCEL_BLOCK"))
	lexer.Add([]byte(SELECT_CELL), token("EXCEL_CELL"))

	// lexer.Add([]byte("^(?:[\\w]\\:|\\)(\\[a-z_\\-\\s0-9\\.]+)+\\.(txt|gif|pdf|doc|docx|xls|xlsx)$"), token("DATASET_FILEPATH"))
	// lexer.Add([]byte(cregex.LinkPattern), token("LINK"))
	// lexer.Add([]byte(cregex.GitRepoPattern), token("GIT_REPO"))

	lexer.Add([]byte(SELECT_COLS_LIST), token("COLS_LIST"))
	lexer.Add([]byte(SELECT_COLS_RANGE), token("COLS_RANGE"))
	lexer.Add([]byte(SELECT_COL), token("COL"))
	lexer.Add([]byte(SELECT_ROW), token("RANGE"))

	lexer.Add([]byte(SELECT_NUMERIC), token("BLOCK_NUM"))
	lexer.Add([]byte(SELECT_ALPHA_NUMERIC), token("BLOCK_ALPHA"))
	lexer.Add([]byte(SELECT_FUNCTION), token("BLOCK_FUNC"))

	lexer.Add([]byte(`((col|cols|rows|row))\=\(((:\d+)|(\d+\:)|(\d+\:\d+)|(\:)|(\d+(,\d+)*))(\))`), token("DATASET"))
	lexer.Add([]byte(`((col|cols|rows|row))\=\(((:\d+)|(\d+\:)|(\d+\:\d+)|(\:)|(\d+(,\d+)*))(\))`), token("DATASET_INTO_FILE"))

	err := lexer.Compile()
	if err != nil {
		return nil, err
	}
	return lexer, nil
}

// a lexmachine.Action function which skips the match.
func skip(*lexmachine.Scanner, *machines.Match) (interface{}, error) {
	return nil, nil
}

// a lexmachine.Action function with constructs a Token of the given token type by
// the token type's name.
func token(name string) lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(TokenIds[name], string(m.Bytes), m), nil
	}
}

package main

import (
	"fmt"

	pp "github.com/sniperkit/colly/plugins/app/debug/pp"

	lexmachine "github.com/sniperkit/lexmachine/pkg"
	machines "github.com/sniperkit/lexmachine/pkg/machines"
)

var isDebug bool = true

func main() {
	fmt.Println("start parsing query strings with lexmachine framework...")

	TokenIds = make(map[string]int, len(tokens))

	for i, tok := range tokens {
		TokenIds[tok] = i
	}

	lexer := lexmachine.NewLexer()

	lexer.Add([]byte(` `), token("SPACE"))
	lexer.Add([]byte(`!`), token("BANG"))

	// skipping patterns
	lexer.Add(
		[]byte("( |\t|\n)"),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			// skip white space
			return nil, nil
		},
	)
	lexer.Add(
		[]byte("//[^\n]*\n"),
		func(scan *lexmachine.Scanner, match *machines.Match) (interface{}, error) {
			// skip white space
			return nil, nil
		},
	)

	// add custom rules
	lexer.Add(
		[]byte(SELECT_BLOCK),
		func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return 0, nil
		},
	)

	lexer.Add(
		[]byte(SELECT_CELL),
		func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return 0, nil
		},
	)

	lexer.Add(
		[]byte(SELECT_COLS),
		func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return 0, nil
		},
	)

	lexer.Add(
		[]byte(SELECT_ROWS),
		func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return 0, nil
		},
	)

	lexer.Add(
		[]byte(SELECT_COL),
		func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return 0, nil
		},
	)

	lexer.Add(
		[]byte(SELECT_ROW),
		func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return 0, nil
		},
	)

	lexer.Add(
		[]byte(SELECT_NUMERIC),
		func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return 0, nil
		},
	)

	lexer.Add(
		[]byte(SELECT_ALPHA_NUMERIC),
		func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return 0, nil
		},
	)

	lexer.Add(
		[]byte(SELECT_FUNCTION),
		func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
			return 0, nil
		},
	)

	err := lexer.Compile()
	if err != nil {
		pp.Println("error while trying to compile lexer's rules, msg=", err)
	}

	for _, query := range queryStrings {

		// tokenize (lex) a string construct a Scanner object using the lexer.
		scanner, err := lexer.Scanner([]byte(query))
		if err != nil {
			pp.Println("error while creating new scanner instance msg=", err)
		}

		// the scanner object is an iterator which yields the next token (or error) by calling the Next() method
		for tok, err, eos := scanner.Next(); !eos; tok, err, eos = scanner.Next() {
			if _, is := err.(*machines.UnconsumedInput); is {
				// pp.Println("ui=", ui, "is=", is)
			} else if err != nil {
				pp.Println("erro while trying to scan the next token msg=", err)
			}
			if tok != nil {
				pp.Println("tok=", tok, "tok=", tok)
			}

			// pp.Println("eos=", eos)
		}

	}

}

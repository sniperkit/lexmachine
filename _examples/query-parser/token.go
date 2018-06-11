package main

import (
	lexmachine "github.com/sniperkit/lexmachine/pkg"
	machines "github.com/sniperkit/lexmachine/pkg/machines"
)

var (
	TokenIds map[string]int

	// tokens represents...
	tokens = []string{
		"AT", "PLUS", "STAR", "DASH", "SLASH", "BACKSLASH", "CARROT", "BACKTICK", "COMMA", "LPAREN", "RPAREN",
		"BUS", "COMPUTE", "CHIP", "IGNORE", "LABEL", "SET", "NUMBER", "NAME",
		"COMMENT", "SPACE",
	}

	//-- End
)

type Token struct {
	TokenType int
	Lexeme    string
	Match     *machines.Match
}

func NewToken(tokenType string, m *machines.Match) *Token {
	return &Token{
		TokenType: TokenIds[tokenType], // defined above
		Lexeme:    string(m.Bytes),
		Match:     m,
	}
}

func token(tokenType string) func(*lexmachine.Scanner, *machines.Match) (interface{}, error) {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return NewToken(tokenType, m), nil
	}
}

func tokenAction(name string, tokenIds map[string]int) lexmachine.Action {
	return func(s *lexmachine.Scanner, m *machines.Match) (interface{}, error) {
		return s.Token(tokenIds[name], string(m.Bytes), m), nil
	}
}

package scanner

import "github.com/jasonfriedland/crafting-interpreters/pkg/token"

type Scanner struct {
	start, current, line int
	source               string
}

func New(source string) *Scanner {
	return &Scanner{
		source: source,
		line:   1,
	}
}

func (s *Scanner) Scan() ([]*token.Token, error) {
	var tokens []*token.Token
	for s.current >= len(s.source) {
		c := s.advance()
		switch string(c) {
		case "(":
			tokens = append(tokens, &token.Token{
				Type: token.LEFT_PAREN,
				Line: s.line,
			})
		}
	}
	// Lastly append an EOF
	tokens = append(tokens,
		&token.Token{
			Type:    token.EOF,
			Lexeme:  "",
			Literal: nil,
			Line:    s.line,
		},
	)
	return tokens, nil
}

func (s *Scanner) scanToken() string {
	return ""
}

func (s *Scanner) advance() byte {
	s.current += 1
	return s.source[s.current]
}

package scanner

import (
	"io"

	"github.com/jasonfriedland/crafting-interpreters/pkg/token"
)

type Scanner struct {
	start, current, line int
	source               []byte
}

func New(r io.Reader) (*Scanner, error) {
	s := &Scanner{
		line: 1,
	}
	_, err := r.Read(s.source)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Scanner) Scan() ([]*token.Token, error) {
	var tokens []*token.Token
	for i := 0; i < len(s.source); i++ {
		c := s.source[i]
		switch string(c) {
		case "(":
			tokens = append(tokens, &token.Token{
				Type: token.LEFT_PAREN,
				Line: s.line,
			})
		case ")":
			tokens = append(tokens, &token.Token{
				Type: token.RIGHT_PAREN,
				Line: s.line,
			})
		case "\n":
			s.line += 1
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

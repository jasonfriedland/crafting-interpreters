package scanner

import (
	"fmt"
	"io"

	"github.com/jasonfriedland/crafting-interpreters/pkg/token"
)

type Scanner struct {
	line   int
	source []byte
	tokens []*token.Token
}

func New(r io.Reader) (*Scanner, error) {
	if r == nil {
		return nil, fmt.Errorf("nil Reader")
	}
	s := &Scanner{
		line: 1,
	}
	var err error
	s.source, err = io.ReadAll(r) // read into s.source
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Scanner) Scan() error {
	if s == nil || len(s.source) == 0 {
		return fmt.Errorf("invalid scanner")
	}
	for i := 0; i < len(s.source); i++ {
		c := s.source[i]
		switch string(c) {
		case "(":
			s.addToken(token.LEFT_PAREN)
		case ")":
			s.addToken(token.RIGHT_PAREN)
		case "{":
			s.addToken(token.LEFT_BRACE)
		case "}":
			s.addToken(token.RIGHT_BRACE)
		case ",":
			s.addToken(token.COMMA)
		case ".":
			s.addToken(token.DOT)
		case "-":
			s.addToken(token.MINUS)
		case "+":
			s.addToken(token.PLUS)
		case ";":
			s.addToken(token.SEMICOLON)
		case "*":
			s.addToken(token.STAR)
		case "\n":
			s.line += 1
		}

	}
	// Lastly append an EOF
	s.addToken(token.EOF)
	return nil
}

func (s *Scanner) addToken(tokenType token.TokenType) {
	s.tokens = append(s.tokens, &token.Token{
		Type: tokenType,
		Line: s.line,
	})
}

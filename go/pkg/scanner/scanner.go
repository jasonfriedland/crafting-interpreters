package scanner

import (
	"fmt"
	"io"

	"github.com/jasonfriedland/crafting-interpreters/pkg/token"
)

// Scanner is a Scanner type.
type Scanner struct {
	current int
	line    int
	source  []byte
	tokens  []*token.Token
}

// New builds a new Scanner from the input io.Reader.
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

// Scan is the main routine of a Scanner which iterates over the source, and
// generates a slice of tokens.
func (s *Scanner) Scan() error {
	if s == nil || len(s.source) == 0 {
		return fmt.Errorf("invalid scanner")
	}
	for s.current < len(s.source) {
		c := s.next()
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
		case "!":
			if s.match("=") {
				s.addToken(token.BANG_EQUAL)
			} else {
				s.addToken(token.BANG)
			}
		case "=":
			if s.match("=") {
				s.addToken(token.EQUAL_EQUAL)
			} else {
				s.addToken(token.EQUAL)
			}
		case "<":
			if s.match("=") {
				s.addToken(token.LESS_EQUAL)
			} else {
				s.addToken(token.LESS)
			}
		case ">":
			if s.match("=") {
				s.addToken(token.GREATER_EQUAL)
			} else {
				s.addToken(token.GREATER)
			}
		case "/":
			if s.match("/") {
				// A comment goes until the end of the line.
				for string(s.peek()) != "\n" && !s.eof() {
					s.next()
				}
			} else {
				s.addToken(token.SLASH)
			}
		// Ignore the following
		case " ":
		case "\r":
		case "\t":
		// New line
		case "\n":
			s.line += 1
		default:
			// TODO: return an error
		}

	}
	// Lastly append an EOF
	s.addToken(token.EOF)
	return nil
}

// next returns the current byte and advances the current position.
func (s *Scanner) next() byte {
	c := s.source[s.current]
	s.current++
	return c
}

// peek returns the next byte character but doesn't increment the counter.
func (s *Scanner) peek() byte {
	if s.eof() {
		return 0
	}
	return s.source[s.current]
}

// match returns whether the byte at the current position matches the passed in
// byte; used for matching double charater lexemes e.g. ==, != etc.
func (s *Scanner) match(expected string) bool {
	if s.eof() {
		return false
	}
	if string(s.source[s.current]) != expected {
		return false
	}
	s.current++
	return true
}

// eof returns whether we're at the end of the input source.
func (s *Scanner) eof() bool {
	return s.current >= len(s.source)
}

// addToken appends a new Toekn type to the Scanner's internal storage.
func (s *Scanner) addToken(tokenType token.TokenType) {
	s.tokens = append(s.tokens, &token.Token{
		Type: tokenType,
		Line: s.line,
	})
}

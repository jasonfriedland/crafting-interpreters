package scanner

import (
	"fmt"
	"io"
	"strconv"

	"github.com/jasonfriedland/crafting-interpreters/pkg/token"
)

// Scanner is a Scanner type.
type Scanner struct {
	current int // position of scanner
	start   int // start position of current lexeme
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

// Tokens returns the scanner's parsed tokens.
func (s *Scanner) Tokens() []*token.Token {
	return s.tokens
}

// Line returns the scanner's current line number.
func (s *Scanner) Line() int {
	return s.line
}

// Scan is the main routine of a Scanner which iterates over the source, and
// generates a slice of tokens.
func (s *Scanner) Scan() error {
	if s == nil || len(s.source) == 0 {
		return fmt.Errorf("invalid scanner")
	}
	for s.current < len(s.source) {
		c := s.next()
		s.start = s.current
		switch c {
		case '(':
			s.addToken(token.LEFT_PAREN, string(c))
		case ')':
			s.addToken(token.RIGHT_PAREN, string(c))
		case '{':
			s.addToken(token.LEFT_BRACE, string(c))
		case '}':
			s.addToken(token.RIGHT_BRACE, string(c))
		case ',':
			s.addToken(token.COMMA, string(c))
		case '.':
			s.addToken(token.DOT, string(c))
		case '-':
			s.addToken(token.MINUS, string(c))
		case '+':
			s.addToken(token.PLUS, string(c))
		case ';':
			s.addToken(token.SEMICOLON, string(c))
		case '*':
			s.addToken(token.STAR, string(c))
		case '!':
			if s.match("=") {
				s.addToken(token.BANG_EQUAL, "!=")
			} else {
				s.addToken(token.BANG, string(c))
			}
		case '=':
			if s.match("=") {
				s.addToken(token.EQUAL_EQUAL, "==")
			} else {
				s.addToken(token.EQUAL, string(c))
			}
		case '<':
			if s.match("=") {
				s.addToken(token.LESS_EQUAL, "<=")
			} else {
				s.addToken(token.LESS, string(c))
			}
		case '>':
			if s.match("=") {
				s.addToken(token.GREATER_EQUAL, ">=")
			} else {
				s.addToken(token.GREATER, string(c))
			}
		case '/':
			if s.match("/") {
				// A comment goes until the end of the line.
				for string(s.peek()) != "\n" && !s.eof() {
					s.next()
				}
			} else {
				s.addToken(token.SLASH, string(c))
			}
		case '"':
			err := s.parseString()
			if err != nil {
				return err
			}
		// Ignore the following
		case ' ':
		case '\r':
		case '\t':
		// New line
		case '\n':
			s.line += 1
		default:
			if isDigit(c) {
				err := s.parseNumber()
				if err != nil {
					return err
				}
				continue
			}
			if isAlpha(c) {
				s.parseIdent()
				continue
			}
			return fmt.Errorf("illegal character: %v", string(c))
		}

	}
	// Lastly append an EOF
	s.addToken(token.EOF, "")
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

// peekNext returns the next + 1 byte character but doesn't increment the
// counter.
func (s *Scanner) peekNext() byte {
	if s.eof() {
		return 0
	}
	return s.source[s.current+1]
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

// parseString parses a string literal.
func (s *Scanner) parseString() error {
	for string(s.peek()) != `"` && !s.eof() {
		s.next()
	}
	if s.eof() {
		return fmt.Errorf("un-terminated string")
	}
	s.next() // consume terminating "
	v := string(s.source[s.start : s.current-1])
	s.addTokenLiteral(token.STRING, v, fmt.Sprintf(`"%s"`, v)) // quoted literal as lexeme
	return nil
}

// parseNumber parses a number literal.
func (s *Scanner) parseNumber() error {
	for isDigit(s.peek()) {
		s.next()
	}
	// Look for a fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.next()
	}
	for isDigit(s.peek()) {
		s.next()
	}
	lexeme := string(s.source[s.start-1 : s.current])
	v, err := strconv.ParseFloat(lexeme, 64)
	if err != nil {
		return err
	}
	s.addTokenLiteral(token.NUMBER, v, lexeme)
	return nil
}

// parseIdent parses identifier and keyword literals.
func (s *Scanner) parseIdent() {
	for isAlphaNumeric(s.peek()) {
		s.next()
	}
	v := string(s.source[s.start-1 : s.current])
	if tokenType, found := token.Keywords[v]; found {
		s.addToken(tokenType, v)
	} else {
		s.addTokenLiteral(token.IDENTIFIER, v, v)
	}
}

// eof returns whether we're at the end of the input source.
func (s *Scanner) eof() bool {
	return s.current >= len(s.source)
}

// addToken adds a token that doesn't require a literal value.
func (s *Scanner) addToken(tokenType token.TokenType, lexeme string) {
	s.addTokenLiteral(tokenType, nil, lexeme)
}

// addTokenLiteral appends a new Token type with literal value to the Scanner's
// internal storage.
func (s *Scanner) addTokenLiteral(tokenType token.TokenType, literal any, lexeme string) {
	s.tokens = append(s.tokens, &token.Token{
		Type:    tokenType,
		Line:    s.line,
		Literal: literal,
		Lexeme:  lexeme,
	})
}

// isDigit returns whether the arg is a digit.
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isAlpha returns whether the arg is alpha.
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

// isAlphanumeric returns whether the arg is alphanumeric.
func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

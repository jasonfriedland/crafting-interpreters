package scanner

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/jasonfriedland/crafting-interpreters/pkg/token"
)

func TestScanner_Scan(t *testing.T) {
	tests := []struct {
		name    string
		r       func() io.Reader
		want    []*token.Token
		wantErr bool
	}{
		{
			"Empty case",
			func() io.Reader { return nil },
			nil,
			true,
		},
		{
			"EOL match case",
			func() io.Reader { return strings.NewReader("/") },
			[]*token.Token{
				{
					Type:   token.SLASH,
					Lexeme: "/",
					Line:   1,
				},
				{
					Type: token.EOF,
					Line: 1,
				},
			},
			false,
		},
		{
			"Illegal character case",
			func() io.Reader { return strings.NewReader("$ hello") },
			nil,
			true,
		},
		{
			"Multi tokens case",
			func() io.Reader {
				return strings.NewReader("// this is a comment, with + - = != things!\n! = \n( \t\r <= >= \n) < > \n { }\n, .  + \n -  * \n;!= /  == !\n// the end + !=")
			},
			[]*token.Token{
				{
					Type:   token.BANG,
					Lexeme: "!",
					Line:   2,
				},
				{
					Type:   token.EQUAL,
					Lexeme: "=",
					Line:   2,
				},
				{
					Type:   token.LEFT_PAREN,
					Lexeme: "(",
					Line:   3,
				},
				{
					Type:   token.LESS_EQUAL,
					Lexeme: "<=",
					Line:   3,
				},
				{
					Type:   token.GREATER_EQUAL,
					Lexeme: ">=",
					Line:   3,
				},
				{
					Type:   token.RIGHT_PAREN,
					Lexeme: ")",
					Line:   4,
				},
				{
					Type:   token.LESS,
					Lexeme: "<",
					Line:   4,
				},
				{
					Type:   token.GREATER,
					Lexeme: ">",
					Line:   4,
				},
				{
					Type:   token.LEFT_BRACE,
					Lexeme: "{",
					Line:   5,
				},
				{
					Type:   token.RIGHT_BRACE,
					Lexeme: "}",
					Line:   5,
				},
				{
					Type:   token.COMMA,
					Lexeme: ",",
					Line:   6,
				},
				{
					Type:   token.DOT,
					Lexeme: ".",
					Line:   6,
				},
				{
					Type:   token.PLUS,
					Lexeme: "+",
					Line:   6,
				},
				{
					Type:   token.MINUS,
					Lexeme: "-",
					Line:   7,
				},
				{
					Type:   token.STAR,
					Lexeme: "*",
					Line:   7,
				},
				{
					Type:   token.SEMICOLON,
					Lexeme: ";",
					Line:   8,
				},
				{
					Type:   token.BANG_EQUAL,
					Lexeme: "!=",
					Line:   8,
				},
				{
					Type:   token.SLASH,
					Lexeme: "/",
					Line:   8,
				},
				{
					Type:   token.EQUAL_EQUAL,
					Lexeme: "==",
					Line:   8,
				},
				{
					Type:   token.BANG,
					Lexeme: "!",
					Line:   8,
				},
				{
					Type:   token.EOF,
					Lexeme: "",
					Line:   9,
				},
			},
			false,
		},
		{
			"String case",
			func() io.Reader { return strings.NewReader("// this is a comment\n\"bar\"\n\"hello\"\n") },
			[]*token.Token{
				{
					Type:    token.STRING,
					Literal: "bar",
					Lexeme:  "bar",
					Line:    2,
				},
				{
					Type:    token.STRING,
					Literal: "hello",
					Lexeme:  "hello",
					Line:    3,
				},
				{
					Type: token.EOF,
					Line: 4,
				},
			},
			false,
		},
		{
			"String un-terminated case",
			func() io.Reader { return strings.NewReader("\"foo") },
			nil,
			true,
		},
		{
			"Number case",
			func() io.Reader { return strings.NewReader("// this is a comment\n34.202\n86723\n") },
			[]*token.Token{
				{
					Type:    token.NUMBER,
					Literal: 34.202,
					Lexeme:  "34.202",
					Line:    2,
				},
				{
					Type:    token.NUMBER,
					Literal: 86723.0,
					Lexeme:  "86723",
					Line:    3,
				},
				{
					Type: token.EOF,
					Line: 4,
				},
			},
			false,
		},
		{
			"Identifiers case",
			func() io.Reader { return strings.NewReader("// this is a comment\nif true { foo = \"bar\" }\n") },
			[]*token.Token{
				{
					Type:   token.IF,
					Lexeme: "if",
					Line:   2,
				},
				{
					Type:   token.TRUE,
					Lexeme: "true",
					Line:   2,
				},
				{
					Type:   token.LEFT_BRACE,
					Lexeme: "{",
					Line:   2,
				},
				{
					Type:    token.IDENTIFIER,
					Literal: "foo",
					Lexeme:  "foo",
					Line:    2,
				},
				{
					Type:   token.EQUAL,
					Lexeme: "=",
					Line:   2,
				},
				{
					Type:    token.STRING,
					Literal: "bar",
					Lexeme:  "bar",
					Line:    2,
				},
				{
					Type:   token.RIGHT_BRACE,
					Lexeme: "}",
					Line:   2,
				},
				{
					Type: token.EOF,
					Line: 3,
				},
			},
			false,
		},
		{
			"Read source-1.txt case, valid source",
			func() io.Reader {
				f, _ := os.Open("testdata/source-1.txt")
				return f
			},
			[]*token.Token{
				{
					Type:   token.CLASS,
					Lexeme: "class",
					Line:   1,
				},
				{
					Type:    token.IDENTIFIER,
					Literal: "foo",
					Lexeme:  "foo",
					Line:    1,
				},
				{
					Type:   token.LEFT_BRACE,
					Lexeme: "{",
					Line:   1,
				},
				{
					Type:   token.FUN,
					Lexeme: "fun",
					Line:   2,
				},
				{
					Type:    token.IDENTIFIER,
					Literal: "bar",
					Lexeme:  "bar",
					Line:    2,
				},
				{
					Type:   token.LEFT_PAREN,
					Lexeme: "(",
					Line:   2,
				},
				{
					Type:   token.RIGHT_PAREN,
					Lexeme: ")",
					Line:   2,
				},
				{
					Type:   token.LEFT_BRACE,
					Lexeme: "{",
					Line:   2,
				},
				{
					Type:   token.IF,
					Lexeme: "if",
					Line:   3,
				},
				{
					Type:    token.IDENTIFIER,
					Literal: "x",
					Lexeme:  "x",
					Line:    3,
				},
				{
					Type:   token.EQUAL_EQUAL,
					Lexeme: "==",
					Line:   3,
				},
				{
					Type:   token.TRUE,
					Lexeme: "true",
					Line:   3,
				},
				{
					Type:   token.LEFT_BRACE,
					Lexeme: "{",
					Line:   3,
				},
				{
					Type:   token.RETURN,
					Lexeme: "return",
					Line:   4,
				},
				{
					Type:    token.NUMBER,
					Literal: 45.0,
					Lexeme:  "45",
					Line:    4,
				},
				{
					Type:   token.RIGHT_BRACE,
					Lexeme: "}",
					Line:   5,
				},
				{
					Type:   token.VAR,
					Lexeme: "var",
					Line:   6,
				},
				{
					Type:    token.IDENTIFIER,
					Literal: "baz",
					Lexeme:  "baz",
					Line:    6,
				},
				{
					Type:   token.EQUAL,
					Lexeme: "=",
					Line:   6,
				},
				{
					Type:    token.STRING,
					Literal: "testing",
					Lexeme:  "testing",
					Line:    6,
				},
				{
					Type:   token.RIGHT_BRACE,
					Lexeme: "}",
					Line:   7,
				},
				{
					Type:   token.RIGHT_BRACE,
					Lexeme: "}",
					Line:   8,
				},
				{
					Type: token.EOF,
					Line: 9,
				},
			},
			false,
		},
		{
			"Read source-1.txt case, error source",
			func() io.Reader {
				f, _ := os.Open("testdata/source-2.txt")
				return f
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(tt.r())
			err := s.Scan()
			if (err != nil) != tt.wantErr {
				t.Errorf("Scanner.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(s.tokens, tt.want) {
				t.Errorf("Scanner.Scan() = %+v, want %+v", s.tokens, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Scanner
		wantErr bool
	}{
		{
			"Empty case, error",
			args{},
			nil,
			true,
		},
		{
			"Simple case",
			args{
				strings.NewReader("hello testing"),
			},
			&Scanner{
				line:   1,
				source: []byte("hello testing"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isDigit(t *testing.T) {
	type args struct {
		c []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Simple case",
			args{
				[]byte("2"),
			},
			true,
		},
		{
			"False case, string",
			args{
				[]byte("S"),
			},
			false,
		},
		{
			"False case, symbol",
			args{
				[]byte("|"),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDigit(tt.args.c[0]); got != tt.want {
				t.Errorf("isDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}

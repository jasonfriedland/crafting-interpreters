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
					Type: token.SLASH,
					Line: 1,
				},
				{
					Type: token.EOF,
					Line: 1,
				},
			},
			false,
		},
		{
			"Multi tokens case",
			func() io.Reader {
				return strings.NewReader("// this is a comment, with + - = != things!\n! = \n( \t\r <= >= \n) < > ?\n { }\n, .  + \n -  * \n;!= /  == !\n// the end + !=")
			},
			[]*token.Token{
				{
					Type: token.BANG,
					Line: 2,
				},
				{
					Type: token.EQUAL,
					Line: 2,
				},
				{
					Type: token.LEFT_PAREN,
					Line: 3,
				},
				{
					Type: token.LESS_EQUAL,
					Line: 3,
				},
				{
					Type: token.GREATER_EQUAL,
					Line: 3,
				},
				{
					Type: token.RIGHT_PAREN,
					Line: 4,
				},
				{
					Type: token.LESS,
					Line: 4,
				},
				{
					Type: token.GREATER,
					Line: 4,
				},
				{
					Type: token.LEFT_BRACE,
					Line: 5,
				},
				{
					Type: token.RIGHT_BRACE,
					Line: 5,
				},
				{
					Type: token.COMMA,
					Line: 6,
				},
				{
					Type: token.DOT,
					Line: 6,
				},
				{
					Type: token.PLUS,
					Line: 6,
				},
				{
					Type: token.MINUS,
					Line: 7,
				},
				{
					Type: token.STAR,
					Line: 7,
				},
				{
					Type: token.SEMICOLON,
					Line: 8,
				},
				{
					Type: token.BANG_EQUAL,
					Line: 8,
				},
				{
					Type: token.SLASH,
					Line: 8,
				},
				{
					Type: token.EQUAL_EQUAL,
					Line: 8,
				},
				{
					Type: token.BANG,
					Line: 8,
				},
				{
					Type: token.EOF,
					Line: 9,
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
					Line:    2,
				},
				{
					Type:    token.STRING,
					Literal: "hello",
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
					Line:    2,
				},
				{
					Type:    token.NUMBER,
					Literal: 86723.0,
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
					Type: token.IF,
					Line: 2,
				},
				{
					Type: token.TRUE,
					Line: 2,
				},
				{
					Type: token.LEFT_BRACE,
					Line: 2,
				},
				{
					Type:    token.IDENTIFIER,
					Literal: "foo",
					Line:    2,
				},
				{
					Type: token.EQUAL,
					Line: 2,
				},
				{
					Type:    token.STRING,
					Literal: "bar",
					Line:    2,
				},
				{
					Type: token.RIGHT_BRACE,
					Line: 2,
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
					Type: token.CLASS,
					Line: 1,
				},
				{
					Type:    token.IDENTIFIER,
					Literal: "foo",
					Line:    1,
				},
				{
					Type: token.LEFT_BRACE,
					Line: 1,
				},
				{
					Type: token.FUN,
					Line: 2,
				},
				{
					Type:    token.IDENTIFIER,
					Literal: "bar",
					Line:    2,
				},
				{
					Type: token.LEFT_PAREN,
					Line: 2,
				},
				{
					Type: token.RIGHT_PAREN,
					Line: 2,
				},
				{
					Type: token.LEFT_BRACE,
					Line: 2,
				},
				{
					Type: token.IF,
					Line: 3,
				},
				{
					Type:    token.IDENTIFIER,
					Literal: "x",
					Line:    3,
				},
				{
					Type: token.EQUAL_EQUAL,
					Line: 3,
				},
				{
					Type: token.TRUE,
					Line: 3,
				},
				{
					Type: token.LEFT_BRACE,
					Line: 3,
				},
				{
					Type: token.RETURN,
					Line: 4,
				},
				{
					Type:    token.NUMBER,
					Literal: 45.0,
					Line:    4,
				},
				{
					Type: token.RIGHT_BRACE,
					Line: 5,
				},
				{
					Type: token.VAR,
					Line: 6,
				},
				{
					Type:    token.IDENTIFIER,
					Literal: "baz",
					Line:    6,
				},
				{
					Type: token.EQUAL,
					Line: 6,
				},
				{
					Type:    token.STRING,
					Literal: "testing",
					Line:    6,
				},
				{
					Type: token.RIGHT_BRACE,
					Line: 7,
				},
				{
					Type: token.RIGHT_BRACE,
					Line: 8,
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

package scanner

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/jasonfriedland/crafting-interpreters/pkg/token"
)

func TestScanner_Scan(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		want    []*token.Token
		wantErr bool
	}{
		{
			"Empty case",
			nil,
			nil,
			true,
		},
		{
			"EOL match case",
			strings.NewReader("/"),
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
			"Multi case",
			strings.NewReader("// this is a comment, with + - = != things!\n!he = llo\n( how\t\r <= >= \n) < > are you?\ni {am}\nwell, well. lets + this\nand - the * as\nwell;!= / and == !\n// the end + !="),
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
			strings.NewReader("// this is a comment\n\"bar\"\n\"hello\"\n"),
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
			strings.NewReader("\"foo"),
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, _ := New(tt.r)
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

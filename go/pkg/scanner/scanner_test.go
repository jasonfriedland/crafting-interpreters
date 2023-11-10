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
			"Multi case",
			strings.NewReader("!he = llo\n( how <= >= \n) < > are you?\ni {am}\nwell, well. lets + this\nand - the * as\nwell;!= and == !"),
			[]*token.Token{
				{
					Type: token.BANG,
					Line: 1,
				},
				{
					Type: token.EQUAL,
					Line: 1,
				},
				{
					Type: token.LEFT_PAREN,
					Line: 2,
				},
				{
					Type: token.LESS_EQUAL,
					Line: 2,
				},
				{
					Type: token.GREATER_EQUAL,
					Line: 2,
				},
				{
					Type: token.RIGHT_PAREN,
					Line: 3,
				},
				{
					Type: token.LESS,
					Line: 3,
				},
				{
					Type: token.GREATER,
					Line: 3,
				},
				{
					Type: token.LEFT_BRACE,
					Line: 4,
				},
				{
					Type: token.RIGHT_BRACE,
					Line: 4,
				},
				{
					Type: token.COMMA,
					Line: 5,
				},
				{
					Type: token.DOT,
					Line: 5,
				},
				{
					Type: token.PLUS,
					Line: 5,
				},
				{
					Type: token.MINUS,
					Line: 6,
				},
				{
					Type: token.STAR,
					Line: 6,
				},
				{
					Type: token.SEMICOLON,
					Line: 7,
				},
				{
					Type: token.BANG_EQUAL,
					Line: 7,
				},
				{
					Type: token.EQUAL_EQUAL,
					Line: 7,
				},
				{
					Type: token.BANG,
					Line: 7,
				},
				{
					Type: token.EOF,
					Line: 7,
				},
			},
			false,
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

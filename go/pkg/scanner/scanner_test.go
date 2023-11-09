package scanner

import (
	"reflect"
	"testing"

	"github.com/jasonfriedland/crafting-interpreters/pkg/token"
)

func TestScanner_Scan(t *testing.T) {
	type fields struct {
		start   int
		current int
		line    int
		source  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*token.Token
		wantErr bool
	}{
		{
			"Empty case",
			fields{},
			[]*token.Token{
				{
					Type: token.EOF,
				},
			},
			false,
		},
		{
			"Parens case",
			fields{
				line:   1,
				source: []byte("hello ( how ) are you?"),
			},
			[]*token.Token{
				{
					Type: token.LEFT_PAREN,
					Line: 1,
				},
				{
					Type: token.RIGHT_PAREN,
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
			"Parens case with newlines",
			fields{
				line:   1,
				source: []byte("hello\n( how\n) are you?\n"),
			},
			[]*token.Token{
				{
					Type: token.LEFT_PAREN,
					Line: 2,
				},
				{
					Type: token.RIGHT_PAREN,
					Line: 3,
				},
				{
					Type: token.EOF,
					Line: 4,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Scanner{
				start:   tt.fields.start,
				current: tt.fields.current,
				line:    tt.fields.line,
				source:  tt.fields.source,
			}
			got, err := s.Scan()
			if (err != nil) != tt.wantErr {
				t.Errorf("Scanner.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scanner.Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}

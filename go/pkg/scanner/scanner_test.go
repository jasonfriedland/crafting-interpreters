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
		source  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*token.Token
		wantErr bool
	}{
		{
			"",
			fields{
				source: "hello (\n how are you\n((\n",
			},
			nil,
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

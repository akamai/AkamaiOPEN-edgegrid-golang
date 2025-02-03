package log

import (
	"reflect"
	"testing"
)

func TestFields_Get(t *testing.T) {
	tests := []struct {
		name string
		f    Fields
		want []any
	}{
		{
			name: "happy path",
			f:    Fields{"String": "string", "int": 0, "map": map[string]int{"int": 1}},
			want: []any{"String", "string", "int", 0, "map", map[string]int{"int": 1}},
		},
		{
			name: "empty fields struct",
			f:    Fields{},
			want: []any{},
		},
		{
			name: "nil value",
			f:    map[string]interface{}{"test": nil},
			want: []any{"test", nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fields.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

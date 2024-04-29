package log

import (
	"reflect"
	"testing"
)

func Test_logArgs(t *testing.T) {
	type args struct {
		args []any
	}
	tests := []struct {
		name string
		args args
		want []any
	}{
		{
			name: "only []any... args",
			args: args{
				args: []any{"int", 1, "string", "string", "struct", struct{ a int }{1}},
			},
			want: []any{"int", 1, "string", "string", "struct", struct{ a int }{1}},
		},
		{
			name: "only fields args",
			args: args{
				args: []any{Fields{"int": 1, "string": "string"}},
			},
			want: []any{"int", 1, "string", "string"},
		},
		{
			name: "mixed args",
			args: args{
				args: []any{Fields{"intField": 2, "stringField": "stringField"}, "int", 1, "string", "string"},
			},
			want: []any{"intField", 2, "stringField", "stringField", "int", 1, "string", "string"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logArgs(tt.args.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("logArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

package log

import (
	"log/slog"
	"reflect"
	"testing"
)

func Test_parseDefault(t *testing.T) {
	type args struct {
		buf   []byte
		key   string
		value any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "primitive type",
			args: args{
				buf:   []byte{},
				key:   "key",
				value: 1,
			},
			want: "key=1",
		},
		{
			name: "array/slice",
			args: args{
				buf:   []byte{},
				key:   "key",
				value: [3]int{1, 2, 3},
			},
			want: "key=[1 2 3]",
		},
		{
			name: "struct",
			args: args{
				buf: []byte{},
				key: "key",
				value: struct {
					A string
					b int
				}{"testString", 0},
			},
			want: "key={testString 0}",
		},
		{
			name: "map",
			args: args{
				buf:   []byte{},
				key:   "key",
				value: map[string]any{"a": false, "b": 12},
			},
			want: "key=map[a:false b:12]",
		},
		{
			name: "Fields, with key",
			args: args{
				buf:   []byte{},
				key:   "key",
				value: Fields{"a": false, "b": 12},
			},
			want: "key={\"a\":false,\"b\":12}",
		},
		{
			name: "Fields, without key",
			args: args{
				buf:   []byte{},
				key:   "",
				value: Fields{"a": false, "b": 12},
			},
			want: "a=false b=12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseDefault(tt.args.buf, tt.args.key, tt.args.value); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("parseDefault() = %s, want %s", string(got), tt.want)
			}
		})
	}
}

func Test_normalizeGroupsAndAttributes(t *testing.T) {
	type args struct {
		groupOfAttrs []groupOrAttrs
		numAttrs     int
	}
	tests := []struct {
		name string
		args args
		want []groupOrAttrs
	}{
		{
			name: "attrs with group",
			args: args{
				groupOfAttrs: []groupOrAttrs{{
					group: "123",
					attrs: []slog.Attr{{
						Key:   "key",
						Value: slog.IntValue(1),
					}},
				}},
				numAttrs: 1,
			},
			want: []groupOrAttrs{
				{
					group: "123",
					attrs: []slog.Attr{{
						Key:   "key",
						Value: slog.IntValue(1),
					}},
				},
			},
		},
		{
			name: "no attrs with group",
			args: args{
				groupOfAttrs: []groupOrAttrs{{
					group: "123",
					attrs: []slog.Attr{},
				}},
				numAttrs: 0,
			},
			want: []groupOrAttrs{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeGroupsAndAttributes(tt.args.groupOfAttrs, tt.args.numAttrs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalizeGroupsAndAttributes() = %v, want %v", got, tt.want)
			}
		})
	}
}

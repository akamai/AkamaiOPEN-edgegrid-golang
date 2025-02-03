package log

import (
	"context"
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	type args struct {
		ctx context.Context
		v   Interface
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "happy path",
			args: args{
				ctx: context.TODO(),
				v:   NOPLogger(),
			},
			want: context.WithValue(context.TODO(), logKey{}, NOPLogger()),
		},
		{
			name: "logger is nil",
			args: args{
				ctx: context.TODO(),
				v:   nil,
			},
			want: context.WithValue(context.TODO(), logKey{}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContext(tt.args.ctx, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want Interface
	}{
		{
			name: "happy path",
			args: args{
				ctx: context.WithValue(context.TODO(), logKey{}, NOPLogger()),
			},
			want: NOPLogger(),
		},
		{
			name: "logger not in context",
			args: args{
				ctx: context.TODO(),
			},
			want: Default(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromContext(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

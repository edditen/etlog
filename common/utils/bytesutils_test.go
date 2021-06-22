package utils

import (
	"reflect"
	"testing"
)

func Test_isWhiteByte(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "when a then return false",
			args: args{
				b: 'a',
			},
			want: false,
		},
		{
			name: "when space then return true",
			args: args{
				b: ' ',
			},
			want: true,
		},
		{
			name: "when tab then return true",
			args: args{
				b: '\t',
			},
			want: true,
		},
		{
			name: "when new line then return true",
			args: args{
				b: '\n',
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isWhiteByte(tt.args.b); got != tt.want {
				t.Errorf("isWhiteByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrimBytes(t *testing.T) {
	type args struct {
		bs []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "when a,b,c then return a,b,c",
			args: args{
				bs: []byte{'a', 'b', 'c'},
			},
			want: []byte{'a', 'b', 'c'},
		},
		{
			name: "when space a,b,c then return a,b,c",
			args: args{
				bs: []byte{' ', '\t', 'a', 'b', 'c', ' ', '\n'},
			},
			want: []byte{'a', 'b', 'c'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TrimBytes(tt.args.bs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrimBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

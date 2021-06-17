package utils

import "testing"

func TestLastSubstring(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "when multiple string then return last one",
			args: args{
				s:   "com.github.etlog",
				sep: ".",
			},
			want: "etlog",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LastSubstring(tt.args.s, tt.args.sep); got != tt.want {
				t.Errorf("LastSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFirstSubstring(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "when path deep 1 then return the path",
			args: args{
				s:   "log",
				sep: "/",
			},
			want: "log",
		},
		{
			name: "when path deep 2 then return the path",
			args: args{
				s:   "log/info.log",
				sep: "/",
			},
			want: "log",
		},
		{
			name: "when path deep 3 then return the path",
			args: args{
				s:   "data/log/info.log",
				sep: "/",
			},
			want: "data/log",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstSubstring(tt.args.s, tt.args.sep); got != tt.want {
				t.Errorf("FirstSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}

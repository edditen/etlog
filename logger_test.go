package etlog

import (
	"reflect"
	"testing"
)

func TestNewStdLogger(t *testing.T) {
	tests := []struct {
		name string
		want *StdLogger
	}{
		{
			name: "when new then return new instance",
			want: &StdLogger{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStdLogger(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStdLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStdLogger_Debug(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "when debug then console message",
			args: args{
				msg: "hello",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := NewStdLogger()
			sl.Debug(tt.args.msg)
		})
	}
}

func TestStdLogger_Error(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "when error then console message",
			args: args{
				msg: "hello",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := NewStdLogger()
			sl.Error(tt.args.msg)
		})
	}
}

func TestStdLogger_Fatal(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "when fatal then console message",
			args: args{
				msg: "hello",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := NewStdLogger()
			sl.Fatal(tt.args.msg)
		})
	}
}

func TestStdLogger_Info(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "when info then console message",
			args: args{
				msg: "hello",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := NewStdLogger()
			sl.Info(tt.args.msg)
		})
	}
}

func TestStdLogger_Warn(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "when warn then console message",
			args: args{
				msg: "hello",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := NewStdLogger()
			sl.Warn(tt.args.msg)
		})
	}
}

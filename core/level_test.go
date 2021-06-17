package core

import (
	"fmt"
	"testing"
)

func TestNewLevel(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name string
		args args
		want Level
	}{
		{
			name: "when debug then return DEBUG",
			args: args{
				level: "debug",
			},
			want: DEBUG,
		},
		{
			name: "when info then return INFO",
			args: args{
				level: "info",
			},
			want: INFO,
		},
		{
			name: "when data then return DATA",
			args: args{
				level: "data",
			},
			want: DATA,
		},
		{
			name: "when warn then return WARN",
			args: args{
				level: "warn",
			},
			want: WARN,
		},
		{
			name: "when error then return ERROR",
			args: args{
				level: "error",
			},
			want: ERROR,
		},
		{
			name: "when fatal then return FATAL",
			args: args{
				level: "fatal",
			},
			want: FATAL,
		},
		{
			name: "when FATAL then return FATAL",
			args: args{
				level: "FATAL",
			},
			want: FATAL,
		},
		{
			name: "when empty then return default",
			args: args{
				level: "",
			},
			want: defaultLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLevel(tt.args.level); got != tt.want {
				t.Errorf("NewLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		name string
		l    Level
		want string
	}{
		{
			name: "when DEBUG then return DEBUG",
			l:    DEBUG,
			want: "DEBUG",
		},
		{
			name: "when INFO then return INFO",
			l:    INFO,
			want: "INFO",
		},
		{
			name: "when DATA then return DATA",
			l:    DATA,
			want: "DATA",
		},
		{
			name: "when WARN then return WARN",
			l:    WARN,
			want: "WARN",
		},
		{
			name: "when ERROR then return ERROR",
			l:    ERROR,
			want: "ERROR",
		},
		{
			name: "when FATAL then return FATAL",
			l:    FATAL,
			want: "FATAL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fmt.Sprintf("%s", tt.l); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLarger(t *testing.T) {
	t.Run("when INFO and DEBUG then INFO larger then DEBUG", func(t *testing.T) {
		t.Log("default:", int(defaultLevel))
		t.Log("DEBUG:", int(DEBUG))
		t.Log("INFO:", int(INFO))
		t.Log("DATA:", int(DATA))
		t.Log("WARN:", int(WARN))
		t.Log("ERROR:", int(ERROR))
		t.Log("FATAL:", int(FATAL))
		if INFO > DEBUG {
			t.Log("INFO larger than DEBUG")
		}
	})
}

package core

import (
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestSimpleFormatter_Format(t *testing.T) {
	t.Run("when short format then simple output", func(t *testing.T) {
		formatter := NewSimpleFormatter()
		meta := &LogEntry{
			Time:   time.Date(2021, 6, 15, 12, 20, 45, 152*1e6, time.UTC),
			Level:  INFO,
			Msg:    "hello world",
			Err:    errors.New("oops"),
			Fields: map[string]interface{}{"rootCause": "nil pointer"},
		}
		got := formatter.Format(meta)
		t.Log(got)
		expect := "2021-06-15 12:20:45.152000 [INFO]\t|msg:=hello world|err:=oops|fields:={\"rootCause\":\"nil pointer\"}\n"
		if got.String() != expect {
			t.Errorf("got: %s, expected: %s", got, expect)
		}
		got.Free()

	})
}

func TestFullFormatter_Format(t *testing.T) {

	type args struct {
		meta *LogEntry
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "when miss src then output default",
			args: args{
				meta: &LogEntry{
					Time:   time.Date(2021, 6, 15, 12, 20, 45, 152*1e6, time.UTC),
					Level:  INFO,
					Msg:    "hello world",
					UseLoc: false,
				},
			},
			want: "2021-06-15 12:20:45.152000|INFO|-|-|hello world|-|\n",
		},
		{
			name: "when src then output full",
			args: args{
				meta: &LogEntry{
					Time:     time.Date(2021, 6, 15, 12, 20, 45, 152*1e6, time.UTC),
					Level:    INFO,
					Msg:      "hello world",
					UseLoc:   true,
					Line:     "hello.go:123",
					FuncName: "TestFormatter.func1",
				},
			},
			want: "2021-06-15 12:20:45.152000|INFO|hello.go:123|TestFormatter.func1|hello world|-|\n",
		},
		{
			name: "when fields then output full",
			args: args{
				meta: &LogEntry{
					Time:     time.Date(2021, 6, 15, 12, 20, 45, 152*1e6, time.UTC),
					Level:    INFO,
					Msg:      "hello world",
					UseLoc:   true,
					Line:     "hello.go:123",
					FuncName: "TestFormatter.func1",
					Err:      errors.New("oops"),
					Fields:   map[string]interface{}{"Hello": "world", "abc": 123},
				},
			},
			want: "2021-06-15 12:20:45.152000|INFO|hello.go:123|TestFormatter.func1|hello world|oops|{\"Hello\":\"world\",\"abc\":123}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewFullFormatter()
			got := s.Format(tt.args.meta)
			if got.String() != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
			got.Free()
		})
	}
}

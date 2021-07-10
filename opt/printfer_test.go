package opt

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestInternalLog_Printf(t *testing.T) {
	type fields struct {
		p Printfer
	}
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "when stdout then console log",
			fields: fields{
				p: log.New(os.Stdout, "", log.LstdFlags),
			},
			args: args{
				format: "hello %v",
				v:      []interface{}{"etlog"},
			},
		},
		{
			name: "when stderr then err log",
			fields: fields{
				p: log.New(os.Stderr, "", log.LstdFlags),
			},
			args: args{
				format: "hello %v",
				v:      []interface{}{"etlog"},
			},
		},
		{
			name: "when writer then log",
			fields: fields{
				p: log.New(&bytes.Buffer{}, "", log.LstdFlags),
			},
			args: args{
				format: "hello %v",
				v:      []interface{}{"etlog"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			il := &internalLog{
				p: tt.fields.p,
			}
			il.Printf(tt.args.format, tt.args.v)
		})
	}
}

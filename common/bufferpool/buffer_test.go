package bufferpool

import (
	"bytes"
	"testing"
)

func TestBuffer_AppendValue(t *testing.T) {
	type fields struct {
		buf  *bytes.Buffer
		pool Pool
	}
	type args struct {
		val interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "when append nil then get nil",
			fields: fields{
				buf: new(bytes.Buffer),
			},
			args: args{
				val: nil,
			},
			want: "<nil>",
		},
		{
			name: "when append int then get int",
			fields: fields{
				buf: new(bytes.Buffer),
			},
			args: args{
				val: 123,
			},
			want: "123",
		},
		{
			name: "when append abc then get abc",
			fields: fields{
				buf: new(bytes.Buffer),
			},
			args: args{
				val: "abc",
			},
			want: "abc",
		},
		{
			name: "when append true then get true",
			fields: fields{
				buf: new(bytes.Buffer),
			},
			args: args{
				val: true,
			},
			want: "true",
		},
		{
			name: "when append float then get float",
			fields: fields{
				buf: new(bytes.Buffer),
			},
			args: args{
				val: 1.23,
			},
			want: "1.23",
		},
		{
			name: "when append byte then get byte",
			fields: fields{
				buf: new(bytes.Buffer),
			},
			args: args{
				val: byte('a'),
			},
			want: "a",
		},
		{
			name: "when append bytes then get bytes",
			fields: fields{
				buf: new(bytes.Buffer),
			},
			args: args{
				val: []byte{'a', 'b', 'c'},
			},
			want: "abc",
		},
		{
			name: "when append struct then get struct string",
			fields: fields{
				buf: new(bytes.Buffer),
			},
			args: args{
				val: map[string]interface{}{"hello": "world"},
			},
			want: "map[hello:world]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				buf:  tt.fields.buf,
				pool: tt.fields.pool,
			}
			b.AppendValue(tt.args.val)
			if got := b.String(); got != tt.want {
				t.Errorf("got: %s, want: %s", got, tt.want)
			}
		})
	}
}

package bufferpool

import (
	"bytes"
	"fmt"
	"strconv"
)

// Buffer is a thin wrapper around a byte slice. It's intended to be pooled, so
// the only way to construct one is via a Pool.
type Buffer struct {
	buf  *bytes.Buffer
	pool Pool
}

func newBuffer() *Buffer {
	return &Buffer{
		buf: &bytes.Buffer{},
	}
}

// AppendByte writes a single byte to the Buffer.
func (b *Buffer) AppendByte(v byte) {
	b.buf.WriteByte(v)
}

// AppendString writes a string to the Buffer.
func (b *Buffer) AppendString(s string) {
	b.buf.WriteString(s)
}

// AppendInt appends an integer to the underlying buffer (assuming base 10).
func (b *Buffer) AppendInt(i int64) {
	b.buf.WriteString(strconv.FormatInt(i, 10))
}

// AppendUint appends an unsigned integer to the underlying buffer (assuming
// base 10).
func (b *Buffer) AppendUint(i uint64) {
	b.buf.WriteString(strconv.FormatInt(int64(i), 10))
}

// AppendBool appends a bool to the underlying buffer.
func (b *Buffer) AppendBool(v bool) {
	b.buf.WriteString(strconv.FormatBool(v))
}

// AppendFloat appends a float to the underlying buffer. It doesn't quote NaN
// or +/- Inf.
func (b *Buffer) AppendFloat(f float64) {
	b.buf.WriteString(strconv.FormatFloat(f, 'f', -1, 32))
}

func (b *Buffer) AppendBytes(bs []byte) {
	b.buf.Write(bs)
}

// Len returns the length of the underlying byte slice.
func (b *Buffer) Len() int {
	return b.buf.Len()
}

// Cap returns the capacity of the underlying byte slice.
func (b *Buffer) Cap() int {
	return b.buf.Cap()
}

// Bytes returns a mutable reference to the underlying byte slice.
func (b *Buffer) Bytes() []byte {
	return b.buf.Bytes()
}

// String returns a string copy of the underlying byte slice.
func (b *Buffer) String() string {
	return b.buf.String()
}

// Reset resets the underlying byte slice. Subsequent writes re-use the slice's
// backing array.
func (b *Buffer) Reset() {
	b.buf.Reset()
}

// Write implements io.Writer.
func (b *Buffer) Write(bs []byte) (int, error) {
	b.buf.Write(bs)
	return b.buf.Len(), nil
}

// Free returns the Buffer to its Pool.
//
// Callers must not retain references to the Buffer after calling Free.
func (b *Buffer) Free() {
	if b.pool != nil {
		b.pool.Return(b)
	}
}

// AppendValue append interface value
func (b *Buffer) AppendValue(val interface{}) {
	if val == nil {
		b.AppendString("<nil>")
		return
	}

	if s, ok := val.(string); ok {
		b.AppendString(s)
		return
	}

	if i, ok := val.(int64); ok {
		b.AppendInt(i)
		return
	}

	if f, ok := val.(float64); ok {
		b.AppendFloat(f)
		return
	}

	if bVal, ok := val.(bool); ok {
		b.AppendBool(bVal)
		return
	}

	if bVal, ok := val.(byte); ok {
		b.AppendByte(bVal)
		return
	}

	if bs, ok := val.([]byte); ok {
		b.AppendBytes(bs)
		return
	}

	b.AppendString(fmt.Sprint(val))
}

// AppendKeyValue append string key and interface value
func (b *Buffer) AppendKeyValue(key string, value interface{}) {
	b.AppendString(key)
	b.AppendByte('=')
	b.AppendValue(value)
	b.AppendByte(',')
}

// AppendNewLine append string key and interface value
func (b *Buffer) AppendNewLine() {
	b.buf.WriteByte('\n')
}

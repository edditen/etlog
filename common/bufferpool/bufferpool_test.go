package bufferpool

import (
	"testing"
)

func BenchmarkBorrow(b *testing.B) {
	b.Run("when no buffer", func(b *testing.B) {
		p := NewChannelPool(1024)
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1e5; j++ {
				buf := newBuffer(p)
				buf.AppendBool(true)
			}
		}
	})

	b.Run("when used channel pool", func(b *testing.B) {
		p := NewChannelPool(1024)
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1e5; j++ {
				buf := p.Borrow()
				buf.AppendBool(true)
				buf.Reset()
				buf.Free()
			}
		}
	})

	b.Run("when used sync pool", func(b *testing.B) {
		p := NewSyncPool(1024)
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1e5; j++ {
				buf := p.Borrow()
				buf.AppendBool(true)
				buf.Reset()
				buf.Free()
			}
		}
	})

}

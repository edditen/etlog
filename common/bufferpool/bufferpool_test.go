package bufferpool

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkBorrow(b *testing.B) {
	b.Run("when no pool", func(b *testing.B) {
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
				buf.Free()
			}
		}
	})

}

func BenchmarkBorrow2(b *testing.B) {
	b.Run("count used sync pool time cost", func(b *testing.B) {
		var count int64
		total := int64(1e7)
		p := NewSyncPool(1024)
		wg := new(sync.WaitGroup)

		startTime := time.Now()
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					buf := p.Borrow()
					buf.AppendBool(true)
					buf.Free()
					if atomic.AddInt64(&count, 1) >= total {
						break
					}
				}
			}()
		}
		wg.Wait()
		duration := time.Now().Sub(startTime)
		b.Logf("duration: %v", duration)
	})

	b.Run("count used channel pool time cost", func(b *testing.B) {
		var count int64
		total := int64(1e7)
		p := NewChannelPool(1024)
		wg := new(sync.WaitGroup)

		startTime := time.Now()
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					buf := p.Borrow()
					buf.AppendBool(true)
					buf.Free()
					if atomic.AddInt64(&count, 1) >= total {
						break
					}
				}
			}()
		}
		wg.Wait()
		duration := time.Now().Sub(startTime)
		b.Logf("duration: %v", duration)
	})

	b.Run("count not used pool time cost", func(b *testing.B) {
		var count int64
		total := int64(1e7)
		p := NewChannelPool(1024)
		wg := new(sync.WaitGroup)

		startTime := time.Now()
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					buf := newBuffer(p)
					buf.AppendBool(true)
					if atomic.AddInt64(&count, 1) >= total {
						break
					}
				}
			}()
		}
		wg.Wait()
		duration := time.Now().Sub(startTime)
		b.Logf("duration: %v", duration)
	})
}

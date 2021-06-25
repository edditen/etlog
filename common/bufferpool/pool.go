package bufferpool

import (
	"sync"
)

type Pool interface {
	// Borrow retrieves a Buffer from the pool, creating one if necessary.
	Borrow() *Buffer
	// Return returns a Buffer to the pool.
	Return(buf *Buffer)
}

type ChanPool struct {
	pool chan *Buffer
}

// NewChanPool using chan constructs a new pool of Buffer.
func NewChanPool(max int) *ChanPool {
	return &ChanPool{
		pool: make(chan *Buffer, max),
	}
}

func (p *ChanPool) Borrow() *Buffer {
	var buf *Buffer
	select {
	case buf = <-p.pool:
		buf.Reset()
	default:
		buf = newBuffer()
	}
	buf.pool = p

	return buf
}

func (p *ChanPool) Return(buf *Buffer) {
	select {
	case p.pool <- buf:
	default:
		// let it go, let it go...
	}

}

type SyncPool struct {
	pool *sync.Pool
}

// NewSyncPool using sync.Pool constructs a new Pool.
func NewSyncPool() *SyncPool {
	return &SyncPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return newBuffer()
			},
		}}
}

func (p *SyncPool) Borrow() *Buffer {
	buf := p.pool.Get().(*Buffer)
	buf.Reset()
	buf.pool = p
	return buf
}

func (p *SyncPool) Return(buf *Buffer) {
	p.pool.Put(buf)
}

package bufferpool

import "sync"

type Pool interface {
	Borrow() *Buffer
	Return(buf *Buffer)
}

type ChanPool struct {
	pool chan *Buffer
}

// NewChanPool creates a new pool of Buffer.
func NewChanPool(max int) *ChanPool {
	return &ChanPool{
		pool: make(chan *Buffer, max),
	}
}

// Borrow a Buffer from the pool.
func (p *ChanPool) Borrow() *Buffer {
	var buf *Buffer
	select {
	case buf = <-p.pool:
		buf.pool = p
		buf.Reset()
	default:
		buf = newBuffer(p)
	}

	return buf
}

// Return returns a Buffer to the pool.
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

// NewSyncPool constructs a new Pool.
func NewSyncPool(max int) *SyncPool {
	return &SyncPool{pool: &sync.Pool{
		New: func() interface{} {
			return &Buffer{bs: make([]byte, 0, max)}
		},
	}}
}

// Borrow retrieves a Buffer from the pool, creating one if necessary.
func (p *SyncPool) Borrow() *Buffer {
	buf := p.pool.Get().(*Buffer)
	buf.Reset()
	buf.pool = p
	return buf
}

func (p *SyncPool) Return(buf *Buffer) {
	p.pool.Put(buf)
}

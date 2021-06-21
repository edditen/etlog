package bufferpool

import "sync"

type Pool interface {
	Borrow() *Buffer
	Return(buf *Buffer)
}

type ChannelPool struct {
	pool chan *Buffer
}

// NewChannelPool creates a new pool of Buffer.
func NewChannelPool(max int) *ChannelPool {
	return &ChannelPool{
		pool: make(chan *Buffer, max),
	}
}

// Borrow a Buffer from the pool.
func (p *ChannelPool) Borrow() *Buffer {
	var buf *Buffer
	select {
	case buf = <-p.pool:
		buf.Reset()
	default:
		buf = newBuffer(p)
	}

	return buf
}

// Return returns a Buffer to the pool.
func (p *ChannelPool) Return(buf *Buffer) {
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

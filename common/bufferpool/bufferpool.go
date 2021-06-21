package bufferpool

const _size = 1024 // by default, create 1 KiB buffers

var (
	pool = NewSyncPool(_size)
)

// Borrow retrieves a buffer from the pool, creating one if necessary.
func Borrow() *Buffer {
	return pool.Borrow()
}

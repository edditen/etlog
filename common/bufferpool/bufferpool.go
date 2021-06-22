package bufferpool

var (
	pool = NewSyncPool()
)

// Borrow retrieves a buffer from the pool, creating one if necessary.
func Borrow() *Buffer {
	return pool.Borrow()
}

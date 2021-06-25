package queue

import (
	"sync"
	"sync/atomic"
	"time"
)

// The end of PriorityQueue implementation.

// DelayQueue is an unbounded blocking queue of *Delayed* elements, in which
// an element can only be taken when its delay has expired. The head of the
// queue is the *Delayed* element whose delay expired furthest in the past.
type DelayQueue struct {
	C chan interface{}

	mu sync.Mutex
	pq *PriorityQueue

	// Similar to the sleeping state of runtime.timers.
	sleeping int32
	wakeupC  chan struct{}
}

// NewDelayQueue creates an instance of delayQueue with the specified size.
func NewDelayQueue(size int) *DelayQueue {
	dq := &DelayQueue{
		C:       make(chan interface{}, size),
		pq:      NewPriorityQueue(size),
		wakeupC: make(chan struct{}),
	}
	go dq.init()
	return dq
}

// Run starts an infinite loop, in which it continually waits for an element
// to expire and then send the expired
func (dq *DelayQueue) init() {
	defer close(dq.C)
	for {

		dq.mu.Lock()
		item, delta := dq.pq.PeekAndShift(time.Now().UnixNano())

		if item == nil {
			// No items left or at least one Item is pending.

			// We must ensure the atomicity of the whole operation, which is
			// composed of the above PeekAndShift and the following StoreInt32,
			// to avoid possible race conditions between Offer and Poll.
			atomic.StoreInt32(&dq.sleeping, 1)
		}
		dq.mu.Unlock()

		if item == nil {
			if delta == 0 {
				// No items left.
				select {
				case <-dq.wakeupC:
					// Wait until a new Item is added.
					continue
				}
			} else if delta > 0 {
				// At least one Item is pending.
				select {
				case <-dq.wakeupC:
					// A new Item with an "earlier" expiration than the current "earliest" one is added.
					continue
				case <-time.After(time.Duration(delta) * time.Nanosecond):
					// The current "earliest" Item expires.

					// Reset the sleeping state since there's no need to receive from wakeupC.
					if atomic.SwapInt32(&dq.sleeping, 0) == 0 {
						// A caller of Offer() is being blocked on sending to wakeupC,
						// drain wakeupC to unblock the caller.
						<-dq.wakeupC
					}
					continue
				}
			}
		} else {
			dq.C <- item.value
		}

	}

}

// Offer inserts the element into the current queue.
func (dq *DelayQueue) Offer(elem interface{}, expireAt time.Time) error {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	item := &Item{value: elem, priority: expireAt.UnixNano()}
	if err := dq.pq.Push(item); err != nil {
		return err
	}

	index := item.index

	if index == 0 {
		// A new Item with the earliest expiration is added.
		if atomic.CompareAndSwapInt32(&dq.sleeping, 1, 0) {
			dq.wakeupC <- struct{}{}
		}
	}
	return nil
}

// Take inserts the element into the current queue.
func (dq *DelayQueue) Take(timeout time.Duration) interface{} {
	select {
	case val := <-dq.C:
		return val
	case <-time.After(timeout):
		return nil
	}
}

func (dq *DelayQueue) Cap() int {
	return dq.pq.Cap()
}

func (dq *DelayQueue) Len() int {
	return dq.pq.Len()
}

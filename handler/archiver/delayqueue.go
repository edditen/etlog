package archiver

import (
	"container/heap"
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrQueueIsFull = errors.New("queue is full")
)

// The start of PriorityQueue implementation.
// Borrowed from https://github.com/nsqio/nsq/blob/master/internal/pqueue/pqueue.go

// Item is a priority queue as implemented by a min heap
// ie. the 0th element is the *lowest* value
type Item struct {
	priority int64
	value    interface{}
	index    int
}

func (i *Item) String() string {
	return fmt.Sprintf("(priority:%d, value:%s, index:%d)",
		i.priority, i.value, i.index)
}

func (i *Item) Less(other interface{}) bool {
	return i.priority < other.(*Item).priority
}

type sorter struct {
	items []*Item
}

func newSorter(size int) *sorter {
	return &sorter{items: make([]*Item, 0, size)}
}

// Push heap.Interface: Push, Pop, Len, Less, Swap
func (s *sorter) Push(elem interface{}) {
	s.items = append(s.items, elem.(*Item))
}

func (s *sorter) Pop() interface{} {
	n := s.Len()
	if n > 0 {
		elem := s.items[n-1]
		s.items = s.items[0 : n-1]
		return elem
	}
	return nil
}

func (s *sorter) Len() int {
	return len(s.items)
}

func (s *sorter) Cap() int {
	return cap(s.items)
}

func (s *sorter) Less(i, j int) bool {
	return s.items[i].Less(s.items[j])
}

func (s *sorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
	s.items[i].index = i
	s.items[j].index = j
}

// PriorityQueue priority queue struct
type PriorityQueue struct {
	s  *sorter
	mu *sync.RWMutex
}

func NewPriorityQueue(size int) *PriorityQueue {
	q := &PriorityQueue{
		s:  newSorter(size),
		mu: new(sync.RWMutex),
	}
	heap.Init(q.s)
	return q
}

func (pq *PriorityQueue) Push(elem *Item) error {
	if pq.Len() == pq.Cap() {
		return ErrQueueIsFull
	}

	pq.mu.Lock()
	defer pq.mu.Unlock()

	heap.Push(pq.s, elem)
	return nil
}

func (pq *PriorityQueue) Pop() *Item {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	return heap.Pop(pq.s).(*Item)
}

func (pq *PriorityQueue) Top() *Item {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	if len(pq.s.items) > 0 {
		return pq.s.items[0]
	}
	return nil
}

func (pq *PriorityQueue) Fix(elem *Item, i int) {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	pq.s.items[i] = elem
	heap.Fix(pq.s, i)
}

func (pq *PriorityQueue) Remove(i int) *Item {
	pq.mu.Lock()
	defer pq.mu.Unlock()

	return heap.Remove(pq.s, i).(*Item)
}

func (pq *PriorityQueue) Len() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return pq.s.Len()
}

func (pq *PriorityQueue) Cap() int {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return pq.s.Cap()
}

func (pq *PriorityQueue) PeekAndShift(max int64) (*Item, int64) {
	item := pq.Top()
	if item == nil {
		return nil, 0
	}

	if item.priority > max {
		return nil, item.priority - max
	}
	pq.Remove(0)

	return item, 0
}

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

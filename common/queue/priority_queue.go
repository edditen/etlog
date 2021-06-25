package queue

import (
	"container/heap"
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

var (
	errQueueIsFull = errors.New("queue is full")
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
	if pq.Full() {
		return errQueueIsFull
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

func (pq *PriorityQueue) Full() bool {
	pq.mu.RLock()
	defer pq.mu.RUnlock()

	return pq.s.Len() == pq.s.Cap()
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

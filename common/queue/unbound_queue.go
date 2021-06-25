package queue

import (
	"sync"
	"time"
)

const (
	defaultInitCap     = 100
	defaultScaleFactor = 0.8
	defaultDownFactor  = 0.2
)

type UnboundQueue struct {
	blockingC  chan interface{}
	upFactor   float32
	downFactor float32
	rwMu       *sync.RWMutex
	mu         *sync.Mutex
	ticker     *time.Ticker
}

func NewUnboundQueue() *UnboundQueue {
	uq := &UnboundQueue{
		blockingC:  make(chan interface{}, defaultInitCap),
		upFactor:   defaultScaleFactor,
		downFactor: defaultDownFactor,
		rwMu:       new(sync.RWMutex),
		mu:         new(sync.Mutex),
		ticker:     time.NewTicker(5 * time.Second),
	}

	go uq.run()

	return uq
}

func (uq *UnboundQueue) run() {
	for {
		select {
		case <-uq.ticker.C:
			if uq.scaleNeeded() {
				uq.shrink()
			}
		}
	}
}

func (uq *UnboundQueue) Take(timeout time.Duration) (interface{}, error) {
	if uq.blockingC == nil {
		return nil, ErrNotInit
	}

	uq.rwMu.RLock()
	defer uq.rwMu.RUnlock()

	select {
	case val := <-uq.blockingC:
		return val, nil
	case <-time.After(timeout):
		return nil, ErrTakeTimeout
	}
}

func (uq *UnboundQueue) Offer(val interface{}) error {
	if uq.blockingC == nil {
		return ErrNotInit
	}

	if uq.scaleNeeded() {
		uq.scale()
	}

	uq.rwMu.RLock()
	defer uq.rwMu.RUnlock()

	select {
	case uq.blockingC <- val:
		return nil
	default:
		// should not happened ever
		return ErrQueueFull
	}
}

func (uq *UnboundQueue) Len() int {
	uq.rwMu.RLock()
	uq.rwMu.RUnlock()
	return len(uq.blockingC)
}

func (uq *UnboundQueue) Cap() int {
	uq.rwMu.RLock()
	uq.rwMu.RUnlock()
	return cap(uq.blockingC)
}

func (uq *UnboundQueue) scaleNeeded() bool {
	uq.rwMu.RLock()
	uq.rwMu.RUnlock()
	threshold := int(uq.upFactor * float32(cap(uq.blockingC)))
	if len(uq.blockingC) > threshold {
		return true
	}
	return false
}

func (uq *UnboundQueue) scale() {
	uq.mu.Lock()
	defer uq.mu.Unlock()

	// double check
	if !uq.scaleNeeded() {
		return
	}

	uq.rwMu.Lock()
	defer uq.rwMu.Unlock()
	newCap := cap(uq.blockingC) * 2
	uq.blockingC = uq.newAndCopy(newCap)
}

func (uq *UnboundQueue) shrink() {
	uq.mu.Lock()
	defer uq.mu.Unlock()

	// double check
	if !uq.shrinkNeeded() {
		return
	}

	uq.rwMu.Lock()
	defer uq.rwMu.Unlock()

	newCap := cap(uq.blockingC) / 2
	uq.blockingC = uq.newAndCopy(newCap)
}

func (uq *UnboundQueue) shrinkNeeded() bool {
	uq.rwMu.RLock()
	uq.rwMu.RUnlock()
	threshold := int(uq.downFactor * float32(cap(uq.blockingC)))
	if len(uq.blockingC) < threshold && len(uq.blockingC)/2 >= defaultInitCap {
		return true
	}
	return false
}

func (uq *UnboundQueue) newAndCopy(newCap int) (newC chan interface{}) {
	newC = make(chan interface{}, newCap)
	chanCopy(uq.blockingC, newC)
	return
}

func chanCopy(src, dst chan interface{}) {
	for len(src) > 0 {
		select {
		case it := <-src:
			select {
			case dst <- it:
			default:
				return
			}
		default:
			return
		}
	}
}

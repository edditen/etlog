package queue

import "time"

type BlockingQueue struct {
	blockingC chan interface{}
}

func NewBlockingQueue(cap int) *BlockingQueue {
	return &BlockingQueue{
		blockingC: make(chan interface{}, cap),
	}
}

func (bq *BlockingQueue) Take(timeout time.Duration) (interface{}, error) {
	if bq.blockingC == nil {
		return nil, ErrNotInit
	}

	select {
	case val := <-bq.blockingC:
		return val, nil
	case <-time.After(timeout):
		return nil, ErrTakeTimeout
	}
}

func (bq *BlockingQueue) Offer(val interface{}) error {
	if bq.blockingC == nil {
		return ErrNotInit
	}

	select {
	case bq.blockingC <- val:
		return nil
	default:
		return ErrQueueFull
	}
}

func (bq *BlockingQueue) Len() int {
	return len(bq.blockingC)
}

func (bq *BlockingQueue) Cap() int {
	return cap(bq.blockingC)
}

package queue

import (
	"github.com/pkg/errors"
	"time"
)

var (
	ErrNotInit     = errors.New("queue is not init")
	ErrQueueFull   = errors.New("queue is full")
	ErrTakeTimeout = errors.New("take from queue timeout")
)

type Queue interface {
	Take(timeout time.Duration) (interface{}, error)
	Offer(val interface{}) error
	Len() int
	Cap() int
}

package archiver

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestPriorityQueue(t *testing.T) {

	t.Run("test priority queue", func(t *testing.T) {
		q := NewPriorityQueue(100)

		q.Push(&Item{priority: 8, value: "8"})
		q.Push(&Item{priority: 7, value: "7"})
		q.Push(&Item{priority: 2, value: "6"})
		q.Push(&Item{priority: 1, value: "5"})
		q.Push(&Item{priority: 5, value: "5"})
		q.Push(&Item{priority: 4, value: "4"})

		x := q.Pop()
		fmt.Println(x.priority, x.value)

		for q.Len() > 0 {
			x = q.Pop()
			fmt.Println(x.priority, x.value)
		}

	})

	t.Run("test delay queue", func(t *testing.T) {
		q := NewDelayQueue(100)
		now := time.Now()
		log.Println("start time", now)
		q.Offer("500", now.Add(500*time.Millisecond))
		q.Offer("200", now.Add(200*time.Millisecond))
		q.Offer("10", now.Add(10*time.Millisecond))
		q.Offer("40", now.Add(40*time.Millisecond))
		q.Offer("400", now.Add(400*time.Millisecond))
		for i := 0; i < 5; i++ {
			elen := q.Take(5 * time.Second)
			log.Println(elen)
		}
		elen := q.Take(1 * time.Second)
		log.Println("done", elen)

	})
}

package queue

import (
	"fmt"
	"testing"
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
}

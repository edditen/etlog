package archiver

import (
	"fmt"
	"log"
	"sync"
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
		q := NewDelayQueue(5)
		beginTime := time.Now()

		log.Println("start time", beginTime)
		if err := q.Offer("500", beginTime.Add(500*time.Millisecond)); err != nil {
			t.Errorf("dont want error %s", err)
		}

		if err := q.Offer("200", beginTime.Add(200*time.Millisecond)); err != nil {
			t.Errorf("dont want error %s", err)
		}

		if err := q.Offer("10", beginTime.Add(10*time.Millisecond)); err != nil {
			t.Errorf("dont want error %s", err)
		}

		if err := q.Offer("40", beginTime.Add(40*time.Millisecond)); err != nil {
			t.Errorf("dont want error %s", err)
		}

		if err := q.Offer("400", beginTime.Add(400*time.Millisecond)); err != nil {
			t.Errorf("dont want error %s", err)
		}

		if err := q.Offer("99", beginTime.Add(99*time.Millisecond)); err == nil {
			t.Errorf("want error %s", err)
		}

		/////////
		if elen := q.Take(1 * time.Second); elen != "10" {
			t.Errorf("want 10, got: %s", elen)
		}
		if delay := time.Now().Sub(beginTime); delay < 10*time.Millisecond ||
			delay > 15*time.Millisecond {
			t.Errorf("want 10ms, got delay: %s", delay)
		}

		if elen := q.Take(1 * time.Second); elen != "40" {
			t.Errorf("want 40, got: %s", elen)
		}
		if delay := time.Now().Sub(beginTime); delay < 40*time.Millisecond ||
			delay > 50*time.Millisecond {
			t.Errorf("want 40ms, got delay: %s", delay)
		}

		if elen := q.Take(1 * time.Second); elen != "200" {
			t.Errorf("want 200, got: %s", elen)
		}
		if delay := time.Now().Sub(beginTime); delay < 200*time.Millisecond ||
			delay > 205*time.Millisecond {
			t.Errorf("want 200ms, got delay: %s", delay)
		}

		if elen := q.Take(1 * time.Second); elen != "400" {
			t.Errorf("want 400, got: %s", elen)
		}
		if delay := time.Now().Sub(beginTime); delay < 400*time.Millisecond ||
			delay > 405*time.Millisecond {
			t.Errorf("want 400ms, got delay: %s", delay)
		}

		if elen := q.Take(1 * time.Second); elen != "500" {
			t.Errorf("want 500, got: %s", elen)
		}
		if delay := time.Now().Sub(beginTime); delay < 500*time.Millisecond ||
			delay > 505*time.Millisecond {
			t.Errorf("want 500ms, got delay: %s", delay)
		}

		newBegin := time.Now()
		if elen := q.Take(1000 * time.Millisecond); elen != nil {
			t.Errorf("want nil, got: %s", elen)
		}
		if delay := time.Now().Sub(newBegin); delay < 1000*time.Millisecond ||
			delay > 1005*time.Millisecond {
			t.Errorf("want 1000ms, got delay: %s", delay)
		}

		log.Println("done")

	})
}

func TestChan(t *testing.T) {
	t.Run("test close chan", func(t *testing.T) {
		exitC := make(chan interface{})

		wg := new(sync.WaitGroup)
		wg.Add(2)
		go func() {
			defer wg.Done()

			select {
			case <-exitC:
				log.Println("chan shutdown")
			}

			select {
			case <-exitC:
				log.Println("chan shutdown")
			default:
				log.Println("default")
			}

		}()

		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			close(exitC)

		}()

		wg.Wait()
		log.Println("done")

	})
}

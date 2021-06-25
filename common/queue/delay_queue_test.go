package queue

import (
	"log"
	"sync"
	"testing"
	"time"
)

const (
	diff = 10
)

func TestDelayQueue(t *testing.T) {

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
			delay > (10+diff)*time.Millisecond {
			t.Errorf("want 10ms, got delay: %s", delay)
		}

		if elen := q.Take(1 * time.Second); elen != "40" {
			t.Errorf("want 40, got: %s", elen)
		}
		if delay := time.Now().Sub(beginTime); delay < 40*time.Millisecond ||
			delay > (40+diff)*time.Millisecond {
			t.Errorf("want 40ms, got delay: %s", delay)
		}

		if elen := q.Take(1 * time.Second); elen != "200" {
			t.Errorf("want 200, got: %s", elen)
		}
		if delay := time.Now().Sub(beginTime); delay < 200*time.Millisecond ||
			delay > (200+diff)*time.Millisecond {
			t.Errorf("want 200ms, got delay: %s", delay)
		}

		if elen := q.Take(1 * time.Second); elen != "400" {
			t.Errorf("want 400, got: %s", elen)
		}
		if delay := time.Now().Sub(beginTime); delay < 400*time.Millisecond ||
			delay > (400+diff)*time.Millisecond {
			t.Errorf("want 400ms, got delay: %s", delay)
		}

		if elen := q.Take(1 * time.Second); elen != "500" {
			t.Errorf("want 500, got: %s", elen)
		}
		if delay := time.Now().Sub(beginTime); delay < 500*time.Millisecond ||
			delay > (500+diff)*time.Millisecond {
			t.Errorf("want 500ms, got delay: %s", delay)
		}

		newBegin := time.Now()
		if elen := q.Take(1000 * time.Millisecond); elen != nil {
			t.Errorf("want nil, got: %s", elen)
		}
		if delay := time.Now().Sub(newBegin); delay < 1000*time.Millisecond ||
			delay > (1000+diff)*time.Millisecond {
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

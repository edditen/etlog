package queue

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestUnboundQueue_Offer(t *testing.T) {
	t.Run("offer more than 100 then ok", func(t *testing.T) {
		q := NewUnboundQueue()
		for i := 0; i < defaultInitCap; i++ {
			if err := q.Offer(i); err != nil {
				t.Errorf("not want err here")
				return
			} else {
				t.Log(fmt.Sprintf("i=%d,\tlen=%d,\tcap=%d", i, q.Len(), q.Cap()))
			}
		}

		if err := q.Offer(12); err != nil {
			t.Errorf("not want err here")
			return
		}

		t.Log(fmt.Sprintf("len=%d,\tcap=%d", q.Len(), q.Cap()))

		if c := q.Cap(); c != 2*defaultInitCap {
			t.Errorf("want:%d, got:%d", 2*defaultInitCap, c)
			return
		}

		if err := q.Offer(45); err != nil {
			t.Errorf("not want err here")
			return
		}

		t.Log(fmt.Sprintf("len=%d,\tcap=%d", q.Len(), q.Cap()))

	})

}

func TestUnboundQueue_Take(t *testing.T) {
	t.Run("offer and take 100 then ok", func(t *testing.T) {
		q := NewUnboundQueue()
		for i := 0; i < defaultInitCap; i++ {
			if err := q.Offer(i); err != nil {
				t.Errorf("not want err here")
				return
			} else {
				t.Log(fmt.Sprintf("i=%d,\tlen=%d,\tcap=%d", i, q.Len(), q.Cap()))
			}
		}

		t.Log("=============")

		for i := 0; i < defaultInitCap; i++ {
			if v, err := q.Take(10 * time.Millisecond); err != nil {
				t.Errorf("not want err here")
				return
			} else {
				t.Log(fmt.Sprintf("i=%d,\tlen=%d,\tcap=%d", i, q.Len(), q.Cap()))
				if !reflect.DeepEqual(v, i) {
					t.Errorf("want:%v, got:%v", i, v)
					return
				}
			}
		}

	})
}

func BenchmarkUnboundQueue(b *testing.B) {
	b.Run("bench offer and take concurrently", func(b *testing.B) {
		q := NewUnboundQueue()
		wg := new(sync.WaitGroup)
		iters := 100
		count := 100000
		var total int64
		total = int64(iters) * int64(count)
		totalNum := total

		b.Logf("begin len:%d, cap:%d", q.Len(), q.Cap())
		begin := time.Now()
		for i := 0; i < iters; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < count; j++ {
					q.Offer(j)
				}

			}()
		}

		wg2 := new(sync.WaitGroup)
		b.Logf("begin2 len:%d, cap:%d", q.Len(), q.Cap())
		begin2 := time.Now()
		for i := 0; i < iters; i++ {
			wg2.Add(1)
			go func() {
				defer wg2.Done()
				for atomic.AddInt64(&total, -1) > 0 {

					if _, err := q.Take(1000 * time.Millisecond); err != nil {
						b.Log(err)
					}

				}

			}()
		}

		wg.Wait()
		offerCost := time.Now().Sub(begin)
		b.Logf("end len:%d, cap:%d", q.Len(), q.Cap())

		wg2.Wait()
		takeCost := time.Now().Sub(begin2)
		b.Logf("end2 len:%d, cap:%d", q.Len(), q.Cap())

		b.Logf("total:%d,\tofferCost:%v,\teach:%dns", totalNum, offerCost, offerCost.Nanoseconds()/totalNum)
		b.Logf("total:%d,\ttakeCost:%v,\teach:%dns", totalNum, takeCost, offerCost.Nanoseconds()/totalNum)
		time.Sleep(2 * time.Second)
		b.Logf("final len:%d, cap:%d", q.Len(), q.Cap())
		time.Sleep(2 * time.Second)
		b.Logf("final2 len:%d, cap:%d", q.Len(), q.Cap())
		time.Sleep(2 * time.Second)
		b.Logf("final3 len:%d, cap:%d", q.Len(), q.Cap())
		time.Sleep(2 * time.Second)
		b.Logf("final4 len:%d, cap:%d", q.Len(), q.Cap())
		time.Sleep(2 * time.Second)
		b.Logf("final5 len:%d, cap:%d", q.Len(), q.Cap())
		time.Sleep(2 * time.Second)
		b.Logf("final6 len:%d, cap:%d", q.Len(), q.Cap())
		time.Sleep(2 * time.Second)
		b.Logf("final7 len:%d, cap:%d", q.Len(), q.Cap())
		time.Sleep(2 * time.Second)
		b.Logf("final8 len:%d, cap:%d", q.Len(), q.Cap())
	})
}

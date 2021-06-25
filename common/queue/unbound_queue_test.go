package queue

import (
	"fmt"
	"reflect"
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
		for i := 0; i < b.N; i++ {
			q.Offer(i)
		}
	})
}

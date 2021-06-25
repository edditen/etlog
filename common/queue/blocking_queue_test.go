package queue

import (
	"reflect"
	"testing"
	"time"
)

func TestBlockingQueue_Take(t *testing.T) {
	type fields struct {
		blockingC chan interface{}
	}
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:   "when queue not init then get error",
			fields: fields{},
			args: args{
				timeout: 10 * time.Millisecond,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when queue init then timeout",
			fields: fields{
				blockingC: make(chan interface{}),
			},
			args: args{
				timeout: 10 * time.Millisecond,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bq := &BlockingQueue{
				blockingC: tt.fields.blockingC,
			}
			got, err := bq.Take(tt.args.timeout)
			t.Log(got, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("Take() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Take() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOfferAndTake(t *testing.T) {
	t.Run("when offer and take then ok", func(t *testing.T) {
		bq := NewBlockingQueue(10)

		want := "hello"
		if err := bq.Offer(want); err != nil {
			t.Errorf("not want get error: %+v", err)
			return
		}

		got, err := bq.Take(10 * time.Millisecond)
		if err != nil {
			t.Errorf("not want get error: %+v", err)
			return
		}

		if !reflect.DeepEqual(got, "hello") {
			t.Errorf("Take() got = %v, want %v", got, want)
		}

	})

	t.Run("when queue is full then get full error", func(t *testing.T) {
		bq := NewBlockingQueue(1)

		want := "hello"
		if err := bq.Offer(want); err != nil {
			t.Errorf("not want get error: %+v", err)
			return
		}

		err := bq.Offer(want)
		t.Log(err)
		if err == nil {
			t.Errorf("want get error, but is nil")
			return
		}

	})
}

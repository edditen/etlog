package utils

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	type args struct {
		rate           int64
		intervalMillis int64
	}
	tests := []struct {
		name string
		args args
		want *RateLimiter
	}{

		{
			name: "when rate and interval then new",
			args: args{
				rate:           100,
				intervalMillis: 50,
			},
			want: &RateLimiter{
				rate:     100,
				interval: time.Duration(50) * time.Millisecond,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRateLimiter(tt.args.rate, tt.args.intervalMillis)
			if !reflect.DeepEqual(got.rate, tt.want.rate) || !reflect.DeepEqual(got.interval, tt.want.interval) {
				t.Errorf("NewRateLimiter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRateLimiter_Allowable(t *testing.T) {
	type fields struct {
		lastResetTime time.Time
		rate          int64
		interval      time.Duration
		token         *int64
		mu            *sync.RWMutex
	}
	type testCase = struct {
		name       string
		fields     fields
		allowTimes int64
		want       bool
	}

	t.Run("when rate and interval then new", func(t *testing.T) {
		tests := &testCase{
			fields: fields{
				lastResetTime: time.Now(),
				rate:          100,
				interval:      time.Duration(1000) * time.Millisecond,
				token:         buildAddr(100),
				mu:            new(sync.RWMutex),
			},
			allowTimes: 100,
		}
		r := &RateLimiter{
			lastResetTime: tests.fields.lastResetTime,
			rate:          tests.fields.rate,
			interval:      tests.fields.interval,
			token:         tests.fields.token,
			mu:            tests.fields.mu,
		}

		var i int64
		for i = 0; i < tests.allowTimes; i++ {
			if allowable := r.Allowable(); !allowable {
				t.Errorf("Allowable() = %v, want %v", allowable, true)
			}
		}
		if allowable := r.Allowable(); allowable {
			t.Errorf("Allowable() = %v, want %v", allowable, false)
		}
	})

	t.Run("when reset token then ok", func(t *testing.T) {
		tests := &testCase{
			fields: fields{
				lastResetTime: time.Now(),
				rate:          100,
				interval:      time.Duration(500) * time.Millisecond,
				token:         buildAddr(100),
				mu:            new(sync.RWMutex),
			},
			allowTimes: 100,
		}
		r := &RateLimiter{
			lastResetTime: tests.fields.lastResetTime,
			rate:          tests.fields.rate,
			interval:      tests.fields.interval,
			token:         tests.fields.token,
			mu:            tests.fields.mu,
		}

		var i int64
		for i = 0; i < tests.allowTimes; i++ {
			if allowable := r.Allowable(); !allowable {
				t.Errorf("Allowable() = %v, want %v", allowable, true)
			}
		}

		time.Sleep(501 * time.Millisecond)

		for i = 0; i < tests.allowTimes; i++ {
			if allowable := r.Allowable(); !allowable {
				t.Errorf("Allowable() = %v, want %v", allowable, true)
			}
		}

		if allowable := r.Allowable(); allowable {
			t.Errorf("Allowable() = %v, want %v", allowable, false)
		}

	})
}

func Test_abs(t *testing.T) {
	type args struct {
		val int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{

		{
			name: "when positive then the same",
			args: args{
				val: 12,
			},
			want: 12,
		},
		{
			name: "when negative then the opposite",
			args: args{
				val: -12,
			},
			want: 12,
		},
		{
			name: "when 0 then 0",
			args: args{
				val: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRateLimiter(1, 1).abs(tt.args.val); got != tt.want {
				t.Errorf("abs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildAdder(t *testing.T) {
	type args struct {
		rate int64
	}
	tests := []struct {
		name string
		args args
		want *int64
	}{

		{
			name: "when int64 then *int64",
			args: args{
				rate: 12,
			},
			want: buildAddr(12),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildAddr(tt.args.rate); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decrement(t *testing.T) {
	t.Run("atomic count", func(t *testing.T) {
		var ops int64
		var want int64
		want = 6
		ops = 50*1000 + want

		var wg sync.WaitGroup
		for i := 0; i < 50; i++ {
			wg.Add(1)

			go func() {
				for c := 0; c < 1000; c++ {
					NewRateLimiter(1, 1).decrement(&ops)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		if !reflect.DeepEqual(ops, want) {
			t.Errorf("got: %d, want: %d\n", ops, want)
		}
	})
}

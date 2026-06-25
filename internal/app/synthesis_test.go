package app

import (
	"testing"
	"time"
)

func TestSynthesisTimeoutFloorsAndScales(t *testing.T) {
	cases := []struct {
		fetch time.Duration
		want  time.Duration
	}{
		{0, 120 * time.Second},
		{15 * time.Second, 120 * time.Second},
		{30 * time.Second, 120 * time.Second},
		{60 * time.Second, 240 * time.Second},
		{90 * time.Second, 360 * time.Second},
	}
	for _, c := range cases {
		if got := synthesisTimeout(c.fetch); got != c.want {
			t.Errorf("synthesisTimeout(%v) = %v, want %v", c.fetch, got, c.want)
		}
	}
}

package stopwatch

import (
	"testing"
	"time"
)

func TestGoodFlow(t *testing.T) {
	s := NewStopwatch()
	if _, err := s.Start(); err != nil {
		t.Fatalf("stopwatch could not be started: %v", err)
	}

	if _, err := s.Flash(); err != nil {
		t.Fatalf("stopwatch could not be flashed: %v", err)
	}

	if _, err := s.Stop(); err != nil {
		t.Fatalf("stopwatch could not be stopped: %v", err)
	}

	if s.Total() != s.Duration() {
		t.Fatalf("total not equal to duration for a single run; total=%s; duration=%s", s.Total().String(), s.Duration().String())
	}
}

func TestErrorFlow(t *testing.T) {
	s := NewStopwatch()
	if _, err := s.Stop(); err == nil {
		t.Fatalf("expected an error but didn't receive one when calling stop on an un-started stopwatch")
	}
	if _, err := s.Flash(); err == nil {
		t.Fatalf("expected an error but didn't receive one when calling flash on an un-started stopwatch")
	}
	if s.Total() != time.Duration(0) {
		t.Fatalf("got a non-zero total for an un-started stopwatch (%s)", s.Total().String())
	}
	if s.Duration() != time.Duration(0) {
		t.Fatalf("got a non-zero duration for an un-started stopwatch (%s)", s.Duration().String())
	}
	if s.Flashes() != nil {
		t.Fatalf("got an array of flashes but didn't expect one")
	}
	if s.Average() != 0.0 {
		t.Fatalf("received a non-zero average time: %02f", s.Average())
	}
}

func TestPerf(t *testing.T) {
	var (
		dur time.Duration
		err error
	)

	bench := testing.Benchmark(func(b *testing.B) {
		s := NewStopwatch()
		b.StartTimer()
		s.Start()

		time.Sleep(5 * time.Second)
		dur, err = s.Stop()
		b.StopTimer()

		if err != nil {
			b.Fatalf("received an error stopping the stopwatch: %v", err)
		}
	})

	t.Logf("memallocs=%d", bench.MemAllocs)
	t.Logf("membytes=%d", bench.MemBytes)
	t.Logf("bytes=%d", bench.Bytes)
	t.Logf("duration=%s", bench.T.String())
	t.Logf("stopwatch_duration=%s", dur.String())
}

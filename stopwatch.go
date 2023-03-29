package stopwatch

import (
	"sync"
	"time"
)

type instance struct {
	startedAt time.Time
	stoppedAt time.Time
	flashes   []time.Time
}

func newInstance(startedAt time.Time) instance {
	return instance{
		startedAt: startedAt,
		flashes:   make([]time.Time, 0),
	}
}

type stopwatch struct {
	sync.Mutex
	isRunning bool
	current   int
	instances []instance
}

func (s *stopwatch) latest() instance {
	return s.instances[s.current]
}

type Stopwatch interface {
	Start() (time.Time, error)
	Stop() (time.Duration, error)
	Flash() (time.Duration, error)

	Duration() time.Duration
	Total() time.Duration
	Flashes() []time.Time

	Average() float64
}

func NewStopwatch() Stopwatch {
	return &stopwatch{
		isRunning: false,
		current:   -1,
		instances: make([]instance, 0),
	}
}

func (s *stopwatch) Start() (time.Time, error) {
	at := time.Now()
	s.Lock()
	defer s.Unlock()

	if s.isRunning {
		return time.Time{}, ErrAlreadyRunning
	}

	s.isRunning = true
	s.current++
	ins := newInstance(at)
	s.instances = append(s.instances, ins)

	return at, nil
}

func (s *stopwatch) Stop() (time.Duration, error) {
	at := time.Now()
	s.Lock()
	defer s.Unlock()

	if !s.isRunning {
		return time.Duration(0), ErrNotRunning
	}

	s.isRunning = false
	s.instances[s.current].stoppedAt = at

	return s.latest().stoppedAt.Sub(s.latest().startedAt), nil
}

func (s *stopwatch) Flash() (time.Duration, error) {
	at := time.Now()
	s.Lock()
	defer s.Unlock()

	if !s.isRunning {
		return time.Duration(0), ErrNotRunning
	}

	s.instances[s.current].flashes = append(s.instances[s.current].flashes, at)
	return at.Sub(s.instances[s.current].startedAt), nil
}

func (s *stopwatch) Duration() time.Duration {
	at := time.Now()
	s.Lock()
	defer s.Unlock()

	if s.current < 0 {
		return time.Duration(0)
	}

	if s.isRunning {
		return at.Sub(s.latest().startedAt)
	}
	return s.latest().stoppedAt.Sub(s.latest().startedAt)
}

func (s *stopwatch) Flashes() []time.Time {
	s.Lock()
	defer s.Unlock()

	if s.current < 0 {
		return nil
	}
	return s.latest().flashes
}

func (s *stopwatch) aggDuration(at time.Time) time.Duration {
	var totalDuration time.Duration

	for i := 0; i < len(s.instances)-1; i++ {
		totalDuration += s.instances[i].stoppedAt.Sub(s.instances[i].startedAt)
	}

	if len(s.instances) > 0 {
		if s.isRunning {
			totalDuration += at.Sub(s.latest().startedAt)
		} else {
			totalDuration += s.latest().stoppedAt.Sub(s.latest().startedAt)
		}
	}

	return totalDuration
}

func (s *stopwatch) Total() time.Duration {
	at := time.Now()
	s.Lock()
	defer s.Unlock()

	return s.aggDuration(at)
}

func (s *stopwatch) Average() float64 {
	at := time.Now()
	s.Lock()
	defer s.Unlock()

	totalDuration := s.aggDuration(at)
	totalInstances := len(s.instances)
	if totalInstances == 0 {
		totalInstances = 1
	}

	return float64(totalDuration) / float64(totalInstances)
}

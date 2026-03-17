package timer

import (
	"sync"
	"time"
)

type State int

const (
	Stopped State = iota
	Running
	Paused
)

type Timer struct {
	mu        sync.Mutex
	state     State
	interval  time.Duration
	remaining time.Duration
	ticker    *time.Ticker
	done      chan struct{}

	OnTick    func(remaining time.Duration)
	OnExpired func()
}

func New(intervalMinutes int) *Timer {
	d := time.Duration(intervalMinutes) * time.Minute
	return &Timer{
		interval:  d,
		remaining: d,
		state:     Stopped,
	}
}

func (t *Timer) State() State {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.state
}

func (t *Timer) Remaining() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.remaining
}

func (t *Timer) Start() {
	t.mu.Lock()
	if t.state == Running {
		t.mu.Unlock()
		return
	}
	t.state = Running
	t.remaining = t.interval
	t.done = make(chan struct{})
	t.ticker = time.NewTicker(1 * time.Second)
	done := t.done
	ticker := t.ticker
	t.mu.Unlock()

	go t.run(done, ticker)
}

func (t *Timer) run(done chan struct{}, ticker *time.Ticker) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			t.mu.Lock()
			if t.state != Running {
				t.mu.Unlock()
				continue
			}
			t.remaining -= time.Second
			remaining := t.remaining
			expired := t.remaining <= 0
			onTick := t.OnTick
			onExpired := t.OnExpired
			if expired {
				t.state = Stopped
				ticker.Stop()
			}
			t.mu.Unlock()

			if onTick != nil && !expired {
				onTick(remaining)
			}

			if expired {
				if onExpired != nil {
					onExpired()
				}
				return
			}
		}
	}
}

func (t *Timer) Pause() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state == Running {
		t.state = Paused
	}
}

func (t *Timer) Resume() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.state == Paused {
		t.state = Running
	}
}

func (t *Timer) Reset() {
	t.Stop()
	t.Start()
}

func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ticker != nil {
		t.ticker.Stop()
	}
	if t.done != nil {
		select {
		case <-t.done:
			// already closed
		default:
			close(t.done)
		}
	}
	t.state = Stopped
	t.remaining = t.interval
}

func (t *Timer) SetInterval(minutes int) {
	t.mu.Lock()
	t.interval = time.Duration(minutes) * time.Minute
	t.mu.Unlock()
}

package timer

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	tm := New(5)
	if tm.State() != Stopped {
		t.Error("expected initial state to be Stopped")
	}
	if tm.Remaining() != 5*time.Minute {
		t.Errorf("expected remaining=5m, got %v", tm.Remaining())
	}
}

func TestStartAndTick(t *testing.T) {
	tm := New(1) // 1 minute
	var tickCount atomic.Int32
	tm.OnTick = func(remaining time.Duration) {
		tickCount.Add(1)
	}

	tm.Start()
	defer tm.Stop()

	time.Sleep(2500 * time.Millisecond)
	if tickCount.Load() < 1 {
		t.Error("expected at least 1 tick")
	}
	if tm.State() != Running {
		t.Error("expected Running state")
	}
}

func TestPauseResume(t *testing.T) {
	tm := New(1)
	tm.Start()
	defer tm.Stop()

	time.Sleep(100 * time.Millisecond)
	tm.Pause()
	if tm.State() != Paused {
		t.Error("expected Paused state")
	}

	remaining1 := tm.Remaining()
	time.Sleep(1500 * time.Millisecond)
	remaining2 := tm.Remaining()

	if remaining2 != remaining1 {
		t.Error("expected remaining to not change while paused")
	}

	tm.Resume()
	if tm.State() != Running {
		t.Error("expected Running state after resume")
	}
}

func TestSetInterval(t *testing.T) {
	tm := New(5)
	tm.SetInterval(10)
	tm.Start()
	defer tm.Stop()

	if tm.Remaining() != 10*time.Minute {
		t.Errorf("expected 10m remaining after SetInterval, got %v", tm.Remaining())
	}
}

func TestExpired(t *testing.T) {
	tm := New(1)
	tm.SetInterval(0) // We need to test expiry quickly

	// Create a timer with very short duration manually
	tm2 := &Timer{
		interval:  2 * time.Second,
		remaining: 2 * time.Second,
		state:     Stopped,
	}

	var expired atomic.Bool
	tm2.OnExpired = func() {
		expired.Store(true)
	}
	tm2.OnTick = func(remaining time.Duration) {}

	tm2.Start()
	time.Sleep(3500 * time.Millisecond)

	if !expired.Load() {
		t.Error("expected timer to have expired")
	}
}

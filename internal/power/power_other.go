//go:build !darwin || !cgo

package power

type Event int

const (
	ScreenSleep Event = iota
	ScreenWake
)

type Watcher struct {
	Events chan Event
}

func NewWatcher() *Watcher {
	ch := make(chan Event)
	close(ch)
	return &Watcher{Events: ch}
}

func (w *Watcher) Stop() {}

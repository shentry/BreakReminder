//go:build darwin && cgo

package power

/*
#cgo darwin CFLAGS: -x objective-c -fobjc-arc
#cgo darwin LDFLAGS: -framework Cocoa

#import <Cocoa/Cocoa.h>

extern void powerHandleEvent(uintptr_t handle, int event);

enum {
	BREventScreenSleep = 0,
	BREventScreenWake  = 1,
};

typedef struct {
	uintptr_t handle;
	// NSWorkspace observers
	void *screenSleepObs;
	void *screenWakeObs;
	void *sessionResignObs;
	void *sessionBecomeObs;
	// NSDistributedNotificationCenter observers (screen lock/unlock)
	void *screenLockedObs;
	void *screenUnlockedObs;
} BRPowerWatcher;

static BRPowerWatcher *BRPowerWatcherStart(uintptr_t handle) {
	BRPowerWatcher *w = calloc(1, sizeof(BRPowerWatcher));
	if (w == NULL) return NULL;
	w->handle = handle;

	NSOperationQueue *mainQ = [NSOperationQueue mainQueue];

	// 1) NSWorkspace notifications for display sleep/wake
	NSNotificationCenter *wsCenter = [[NSWorkspace sharedWorkspace] notificationCenter];

	id obs1 = [wsCenter addObserverForName:NSWorkspaceScreensDidSleepNotification
	                                object:nil queue:mainQ
	                            usingBlock:^(__unused NSNotification *n) {
		powerHandleEvent(handle, BREventScreenSleep);
	}];
	id obs2 = [wsCenter addObserverForName:NSWorkspaceScreensDidWakeNotification
	                                object:nil queue:mainQ
	                            usingBlock:^(__unused NSNotification *n) {
		powerHandleEvent(handle, BREventScreenWake);
	}];
	id obs3 = [wsCenter addObserverForName:NSWorkspaceSessionDidResignActiveNotification
	                                object:nil queue:mainQ
	                            usingBlock:^(__unused NSNotification *n) {
		powerHandleEvent(handle, BREventScreenSleep);
	}];
	id obs4 = [wsCenter addObserverForName:NSWorkspaceSessionDidBecomeActiveNotification
	                                object:nil queue:mainQ
	                            usingBlock:^(__unused NSNotification *n) {
		powerHandleEvent(handle, BREventScreenWake);
	}];

	// 2) NSDistributedNotificationCenter for screen lock/unlock
	//    Most reliable way to detect Ctrl+Cmd+Q lock screen
	NSDistributedNotificationCenter *distCenter = [NSDistributedNotificationCenter defaultCenter];

	id obs5 = [distCenter addObserverForName:@"com.apple.screenIsLocked"
	                                  object:nil queue:mainQ
	                              usingBlock:^(__unused NSNotification *n) {
		powerHandleEvent(handle, BREventScreenSleep);
	}];
	id obs6 = [distCenter addObserverForName:@"com.apple.screenIsUnlocked"
	                                  object:nil queue:mainQ
	                              usingBlock:^(__unused NSNotification *n) {
		powerHandleEvent(handle, BREventScreenWake);
	}];

	w->screenSleepObs   = (__bridge_retained void *)obs1;
	w->screenWakeObs    = (__bridge_retained void *)obs2;
	w->sessionResignObs = (__bridge_retained void *)obs3;
	w->sessionBecomeObs = (__bridge_retained void *)obs4;
	w->screenLockedObs  = (__bridge_retained void *)obs5;
	w->screenUnlockedObs = (__bridge_retained void *)obs6;

	return w;
}

static void BRRemoveWsObs(NSNotificationCenter *center, void **ref) {
	if (ref == NULL || *ref == NULL) return;
	id obs = (__bridge_transfer id)(*ref);
	[center removeObserver:obs];
	*ref = NULL;
}

static void BRRemoveDistObs(NSDistributedNotificationCenter *center, void **ref) {
	if (ref == NULL || *ref == NULL) return;
	id obs = (__bridge_transfer id)(*ref);
	[center removeObserver:obs];
	*ref = NULL;
}

static void BRPowerWatcherStop(BRPowerWatcher *w) {
	if (w == NULL) return;
	NSNotificationCenter *wsCenter = [[NSWorkspace sharedWorkspace] notificationCenter];
	BRRemoveWsObs(wsCenter, &w->screenSleepObs);
	BRRemoveWsObs(wsCenter, &w->screenWakeObs);
	BRRemoveWsObs(wsCenter, &w->sessionResignObs);
	BRRemoveWsObs(wsCenter, &w->sessionBecomeObs);

	NSDistributedNotificationCenter *distCenter = [NSDistributedNotificationCenter defaultCenter];
	BRRemoveDistObs(distCenter, &w->screenLockedObs);
	BRRemoveDistObs(distCenter, &w->screenUnlockedObs);

	free(w);
}
*/
import "C"

import (
	"runtime/cgo"
	"sync"
)

type Event int

const (
	ScreenSleep Event = iota
	ScreenWake
)

type Watcher struct {
	Events chan Event

	mu      sync.Mutex
	stopped bool
	last    Event
	hasLast bool
	handle  cgo.Handle
	native  *C.BRPowerWatcher
}

func NewWatcher() *Watcher {
	w := &Watcher{
		Events: make(chan Event, 4),
	}
	w.handle = cgo.NewHandle(w)
	w.native = C.BRPowerWatcherStart(C.uintptr_t(w.handle))
	if w.native == nil {
		w.handle.Delete()
		w.stopped = true
		close(w.Events)
	}
	return w
}

func (w *Watcher) Stop() {
	if w == nil {
		return
	}
	w.mu.Lock()
	if w.stopped {
		w.mu.Unlock()
		return
	}
	w.stopped = true
	native := w.native
	w.native = nil
	handle := w.handle
	ch := w.Events
	w.Events = nil
	w.mu.Unlock()

	if native != nil {
		C.BRPowerWatcherStop(native)
	}
	if handle != 0 {
		handle.Delete()
	}
	if ch != nil {
		close(ch)
	}
}

func (w *Watcher) dispatch(event Event) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.stopped || (w.hasLast && w.last == event) || w.Events == nil {
		return
	}
	w.last = event
	w.hasLast = true
	select {
	case w.Events <- event:
	default:
	}
}

//export powerHandleEvent
func powerHandleEvent(handle C.uintptr_t, event C.int) {
	defer func() { _ = recover() }()

	watcher, _ := cgo.Handle(handle).Value().(*Watcher)
	if watcher == nil {
		return
	}
	switch {
	case event == C.BREventScreenSleep:
		watcher.dispatch(ScreenSleep)
	case event == C.BREventScreenWake:
		watcher.dispatch(ScreenWake)
	}
}

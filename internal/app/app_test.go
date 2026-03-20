//go:build darwin

package app

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	fynetest "fyne.io/fyne/v2/test"

	"github.com/zhangxinyu/breakreminder/internal/config"
	"github.com/zhangxinyu/breakreminder/internal/notification"
	"github.com/zhangxinyu/breakreminder/internal/timer"
	"github.com/zhangxinyu/breakreminder/internal/tray"
)

func TestSendSystemNotification(t *testing.T) {
	testApp := fynetest.NewApp()
	a := &App{fyneApp: testApp}

	notification := fyne.NewNotification("该休息一下了！", "喝水：补充水分")

	fynetest.AssertNotificationSent(t, notification, func() {
		a.sendSystemNotification(notification.Title, notification.Content)
	})
}

func TestOnSkipSystemNotificationWaitsForBreakDuration(t *testing.T) {
	testApp := fynetest.NewApp()
	cfg := config.Config{
		IntervalMinutes:   1,
		NotificationStyle: config.NotifySystem,
		BreakDurationSec:  1,
	}

	a := &App{
		fyneApp: testApp,
		cfg:     cfg,
		timer:   timer.New(cfg.IntervalMinutes),
		tray:    tray.New(testApp),
	}
	a.tray.SetCallbacks(nil, nil, nil, nil)
	a.notifier = notification.New(cfg.NotificationStyle, a.sendSystemNotification, nil)

	// Prevent systray access when the delayed finish updates pause state.
	a.tray.SyncPauseState(false, false)

	a.timer.Start()
	a.onSkip()
	fyne.DoAndWait(func() {})

	if a.timer.State() != timer.Stopped {
		t.Fatalf("expected timer to stay stopped during break, got %v", a.timer.State())
	}

	time.Sleep(1200 * time.Millisecond)
	fyne.DoAndWait(func() {})

	if a.timer.State() != timer.Running {
		t.Fatalf("expected timer to restart after break duration, got %v", a.timer.State())
	}
}

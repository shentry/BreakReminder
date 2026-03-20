package app

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"

	"github.com/zhangxinyu/breakreminder/internal/activities"
	"github.com/zhangxinyu/breakreminder/internal/config"
	"github.com/zhangxinyu/breakreminder/internal/notification"
	"github.com/zhangxinyu/breakreminder/internal/power"
	"github.com/zhangxinyu/breakreminder/internal/timer"
	"github.com/zhangxinyu/breakreminder/internal/tray"
	"github.com/zhangxinyu/breakreminder/internal/ui"
)

type PauseReason uint8

const (
	PauseManual PauseReason = 1 << iota
	PauseScreen
)

type App struct {
	fyneApp      fyne.App
	cfg          config.Config
	timer        *timer.Timer
	tray         *tray.Tray
	notifier     *notification.Notifier
	breakWin     *ui.BreakWindow
	settingsWin  *ui.SettingsWindow
	pauseReasons PauseReason
	powerWatcher *power.Watcher
	breakTimerMu sync.Mutex
	breakTimer   *time.Timer
}

func New(fyneApp fyne.App) *App {
	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}

	a := &App{
		fyneApp: fyneApp,
		cfg:     cfg,
	}

	a.timer = timer.New(cfg.IntervalMinutes)
	a.tray = tray.New(fyneApp)
	a.notifier = notification.New(cfg.NotificationStyle, a.sendSystemNotification, a.showBreakPopup)

	a.breakWin = ui.NewBreakWindow(fyneApp, a.onBreakFinished)
	a.settingsWin = ui.NewSettingsWindow(fyneApp, cfg, a.onSettingsSaved)

	// Wire timer callbacks
	a.timer.OnTick = func(remaining time.Duration) {
		fyne.Do(func() {
			a.tray.UpdateCountdown(remaining)
		})
	}
	a.timer.OnExpired = func() {
		fyne.Do(func() {
			a.startBreak(activities.Random())
		})
	}

	// Wire tray callbacks
	a.tray.SetCallbacks(
		a.onPauseResume,
		a.onSkip,
		a.onSettingsOpen,
		a.onQuit,
	)

	return a
}

func (a *App) Run() {
	a.tray.Setup()
	hideDock()

	a.powerWatcher = power.NewWatcher()
	go func() {
		for event := range a.powerWatcher.Events {
			switch event {
			case power.ScreenSleep:
				fyne.Do(func() {
					a.setPauseReason(PauseScreen, true)
					a.breakWin.PauseCountdown()
				})
			case power.ScreenWake:
				fyne.Do(func() {
					a.setPauseReason(PauseScreen, false)
					a.breakWin.ResumeCountdown()
				})
			}
		}
	}()

	a.timer.Start()
	a.fyneApp.Run()
}

func (a *App) setPauseReason(reason PauseReason, active bool) {
	before := a.pauseReasons
	if active {
		a.pauseReasons |= reason
	} else {
		a.pauseReasons &^= reason
	}

	switch {
	case before == 0 && a.pauseReasons != 0:
		a.timer.Pause()
	case before != 0 && a.pauseReasons == 0:
		a.timer.Resume()
	}

	a.tray.SyncPauseState(
		a.pauseReasons&PauseManual != 0,
		a.pauseReasons&PauseScreen != 0,
	)
	if a.pauseReasons == 0 {
		a.tray.UpdateCountdown(a.timer.Remaining())
	}
}

func (a *App) sendSystemNotification(title, body string) {
	if a.fyneApp != nil {
		a.fyneApp.SendNotification(fyne.NewNotification(title, body))
	}
}

func (a *App) showBreakPopup(activity activities.Activity) {
	activateApp()
	a.breakWin.Show(activity, a.cfg.BreakDurationSec)
}

func (a *App) startBreak(activity activities.Activity) {
	a.notifier.Notify(activity)
	if a.cfg.NotificationStyle == config.NotifySystem {
		a.scheduleSystemBreakEnd()
	}
}

func (a *App) scheduleSystemBreakEnd() {
	duration := time.Duration(a.cfg.BreakDurationSec) * time.Second
	if duration <= 0 {
		a.onBreakFinished()
		return
	}

	a.breakTimerMu.Lock()
	if a.breakTimer != nil {
		a.breakTimer.Stop()
	}
	a.breakTimer = time.AfterFunc(duration, func() {
		fyne.Do(func() {
			a.onBreakFinished()
		})
	})
	a.breakTimerMu.Unlock()
}

func (a *App) clearBreakTimer() {
	a.breakTimerMu.Lock()
	defer a.breakTimerMu.Unlock()
	if a.breakTimer != nil {
		a.breakTimer.Stop()
		a.breakTimer = nil
	}
}

func (a *App) onBreakFinished() {
	a.clearBreakTimer()
	a.timer.Reset()
	if a.pauseReasons != 0 {
		a.timer.Pause()
	}
	a.tray.SyncPauseState(a.pauseReasons&PauseManual != 0, a.pauseReasons&PauseScreen != 0)
}

func (a *App) onPauseResume() {
	if a.pauseReasons&PauseManual != 0 {
		a.setPauseReason(PauseManual, false)
	} else {
		a.setPauseReason(PauseManual, true)
	}
}

func (a *App) onSkip() {
	a.timer.Stop()
	fyne.Do(func() {
		a.startBreak(activities.Random())
	})
}

func (a *App) onSettingsOpen() {
	a.settingsWin.Show()
}

func (a *App) onSettingsSaved(cfg config.Config) {
	a.cfg = cfg
	_ = config.Save(cfg)
	a.notifier.SetStyle(cfg.NotificationStyle)
	a.timer.SetInterval(cfg.IntervalMinutes)
	a.timer.Reset()
	if a.pauseReasons != 0 {
		a.timer.Pause()
	}
	a.tray.SyncPauseState(a.pauseReasons&PauseManual != 0, a.pauseReasons&PauseScreen != 0)
	if a.pauseReasons == 0 {
		a.tray.UpdateCountdown(a.timer.Remaining())
	}
}

func (a *App) onQuit() {
	a.clearBreakTimer()
	if a.powerWatcher != nil {
		a.powerWatcher.Stop()
	}
	a.timer.Stop()
	a.fyneApp.Quit()
}

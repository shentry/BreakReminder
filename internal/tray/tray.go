package tray

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/systray"

	"github.com/zhangxinyu/breakreminder/assets"
)

type Tray struct {
	fyneApp       fyne.App
	menu          *fyne.Menu
	countdownItem *fyne.MenuItem
	pauseItem     *fyne.MenuItem
	isPaused      bool
	isAutoPaused  bool

	onPauseResume func()
	onSkip        func()
	onSettings    func()
	onQuit        func()
}

func New(fyneApp fyne.App) *Tray {
	return &Tray{
		fyneApp: fyneApp,
	}
}

func (t *Tray) SetCallbacks(onPauseResume, onSkip, onSettings, onQuit func()) {
	t.onPauseResume = onPauseResume
	t.onSkip = onSkip
	t.onSettings = onSettings
	t.onQuit = onQuit
}

func (t *Tray) Setup() {
	deskApp, ok := t.fyneApp.(desktop.App)
	if !ok {
		return
	}

	t.countdownItem = fyne.NewMenuItem("距离下次休息：--:--", nil)
	t.countdownItem.Disabled = true

	t.pauseItem = fyne.NewMenuItem("暂停", func() {
		if t.onPauseResume != nil {
			t.onPauseResume()
		}
	})

	skipItem := fyne.NewMenuItem("立即休息", func() {
		if t.onSkip != nil {
			t.onSkip()
		}
	})

	settingsItem := fyne.NewMenuItem("设置...", func() {
		if t.onSettings != nil {
			t.onSettings()
		}
	})

	quitItem := fyne.NewMenuItem("退出", func() {
		if t.onQuit != nil {
			t.onQuit()
		}
	})

	t.menu = fyne.NewMenu("BreakReminder",
		t.countdownItem,
		fyne.NewMenuItemSeparator(),
		t.pauseItem,
		skipItem,
		fyne.NewMenuItemSeparator(),
		settingsItem,
		fyne.NewMenuItemSeparator(),
		quitItem,
	)

	deskApp.SetSystemTrayMenu(t.menu)
	deskApp.SetSystemTrayIcon(assets.AppIcon)
}

func (t *Tray) UpdateCountdown(remaining time.Duration) {
	m := int(remaining.Minutes())
	s := int(remaining.Seconds()) % 60
	countdown := fmt.Sprintf("%02d:%02d", m, s)

	systray.SetTitle(countdown)

	t.countdownItem.Label = fmt.Sprintf("距离下次休息：%s", countdown)
	t.menu.Refresh()
}

func (t *Tray) SetPaused(paused bool) {
	t.SyncPauseState(paused, t.isAutoPaused)
}

func (t *Tray) SyncPauseState(manualPaused, autoPaused bool) {
	t.isPaused = manualPaused
	t.isAutoPaused = autoPaused

	if t.pauseItem == nil || t.menu == nil {
		return
	}

	if t.isPaused || t.isAutoPaused {
		t.pauseItem.Label = "继续"
		if t.isAutoPaused {
			systray.SetTitle("已暂停 (息屏)")
		} else {
			systray.SetTitle("已暂停")
		}
	} else {
		t.pauseItem.Label = "暂停"
	}

	t.menu.Refresh()
}

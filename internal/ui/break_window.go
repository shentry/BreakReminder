package ui

import (
	"fmt"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/zhangxinyu/breakreminder/internal/activities"
)

type BreakWindow struct {
	mu       sync.Mutex
	window   fyne.Window
	timer    *time.Ticker
	done     chan struct{}
	onFinish func()
	paused   bool

	remaining int
	total     int

	countdownLabel *widget.Label
	progressBar    *widget.ProgressBar
}

func NewBreakWindow(app fyne.App, onFinish func()) *BreakWindow {
	w := app.NewWindow("休息时间")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(500, 400))
	w.CenterOnScreen()

	return &BreakWindow{
		window:   w,
		onFinish: onFinish,
	}
}

func (bw *BreakWindow) Show(activity activities.Activity, durationSec int) {
	bw.mu.Lock()
	bw.stopTickerLocked()
	bw.paused = false
	bw.remaining = durationSec
	bw.total = durationSec
	bw.mu.Unlock()

	// Header
	header := MakeGradientHeader("休息一下吧", "站起来活动活动身体", 80)

	// Activity card
	badge := MakeBadge(string(activity.Category), ColorPrimaryLight, ColorPrimaryDark)

	titleText := widget.NewLabelWithStyle(activity.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	descText := widget.NewLabel(activity.Description)
	descText.Alignment = fyne.TextAlignCenter
	descText.Wrapping = fyne.TextWrapWord

	activityCard := MakeCard(container.NewVBox(
		container.NewCenter(badge),
		container.NewCenter(titleText),
		descText,
	))

	// Countdown
	bw.countdownLabel = widget.NewLabelWithStyle(
		formatDuration(time.Duration(durationSec)*time.Second),
		fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})

	// Progress bar
	bw.progressBar = widget.NewProgressBar()
	bw.progressBar.Max = 1.0
	bw.progressBar.SetValue(1.0)

	// Buttons
	doneBtn := widget.NewButton("  完成休息  ", func() {
		bw.close()
	})
	doneBtn.Importance = widget.HighImportance

	skipBtn := widget.NewButton("  跳过  ", func() {
		bw.close()
	})
	skipBtn.Importance = widget.MediumImportance

	buttons := container.NewHBox(layout.NewSpacer(), doneBtn, skipBtn, layout.NewSpacer())

	// Layout
	body := container.NewVBox(
		activityCard,
		layout.NewSpacer(),
		container.NewCenter(bw.countdownLabel),
		container.NewPadded(bw.progressBar),
		buttons,
	)

	content := container.NewBorder(header, nil, nil, nil,
		container.NewPadded(body),
	)

	bw.window.SetContent(content)
	bw.window.SetCloseIntercept(func() {
		bw.close()
	})
	bw.window.Show()
	bw.window.RequestFocus()

	bw.startCountdown()
}

func (bw *BreakWindow) PauseCountdown() {
	bw.mu.Lock()
	defer bw.mu.Unlock()
	if bw.paused || bw.timer == nil || bw.remaining <= 0 {
		return
	}
	bw.paused = true
	bw.stopTickerLocked()
}

func (bw *BreakWindow) ResumeCountdown() {
	bw.mu.Lock()
	if !bw.paused || bw.remaining <= 0 {
		bw.mu.Unlock()
		return
	}
	bw.paused = false
	bw.done = make(chan struct{})
	bw.timer = time.NewTicker(1 * time.Second)
	done := bw.done
	ticker := bw.timer
	bw.mu.Unlock()

	go bw.runCountdown(done, ticker)
}

func (bw *BreakWindow) startCountdown() {
	bw.mu.Lock()
	if bw.remaining <= 0 {
		bw.mu.Unlock()
		return
	}
	bw.done = make(chan struct{})
	bw.timer = time.NewTicker(1 * time.Second)
	done := bw.done
	ticker := bw.timer
	bw.mu.Unlock()

	go bw.runCountdown(done, ticker)
}

func (bw *BreakWindow) runCountdown(done chan struct{}, ticker *time.Ticker) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			bw.mu.Lock()
			if bw.timer != ticker || bw.done != done || bw.paused {
				bw.mu.Unlock()
				continue
			}
			bw.remaining--
			remaining := bw.remaining
			total := bw.total
			bw.mu.Unlock()

			if remaining <= 0 {
				fyne.Do(func() {
					bw.close()
				})
				return
			}
			fyne.Do(func() {
				bw.countdownLabel.SetText(formatDuration(time.Duration(remaining) * time.Second))
				bw.progressBar.SetValue(float64(remaining) / float64(total))
			})
		}
	}
}

func (bw *BreakWindow) stopTickerLocked() {
	if bw.timer != nil {
		bw.timer.Stop()
		bw.timer = nil
	}
	if bw.done != nil {
		select {
		case <-bw.done:
		default:
			close(bw.done)
		}
		bw.done = nil
	}
}

func (bw *BreakWindow) close() {
	bw.mu.Lock()
	bw.stopTickerLocked()
	bw.paused = false
	bw.remaining = 0
	bw.total = 0
	bw.mu.Unlock()

	bw.window.Hide()
	if bw.onFinish != nil {
		bw.onFinish()
	}
}

func formatDuration(d time.Duration) string {
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

package ui

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/zhangxinyu/breakreminder/internal/autostart"
	"github.com/zhangxinyu/breakreminder/internal/config"
)

type SettingsWindow struct {
	window fyne.Window
	onSave func(cfg config.Config)
}

func NewSettingsWindow(app fyne.App, cfg config.Config, onSave func(config.Config)) *SettingsWindow {
	w := app.NewWindow("偏好设置")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(380, 420))
	w.CenterOnScreen()

	sw := &SettingsWindow{
		window: w,
		onSave: onSave,
	}

	// --- Interval ---
	intervalOptions := []string{"10", "15", "20", "25", "30", "45", "60", "90", "120"}
	intervalSelect := widget.NewSelect(intervalOptions, nil)
	intervalSelect.PlaceHolder = "选择"
	current := strconv.Itoa(cfg.IntervalMinutes)
	for _, opt := range intervalOptions {
		if opt == current {
			intervalSelect.SetSelected(opt)
			break
		}
	}
	if intervalSelect.Selected == "" {
		intervalSelect.SetSelected(current)
	}

	// --- Break Duration ---
	breakMinutes := cfg.BreakDurationSec / 60
	if breakMinutes <= 0 {
		breakMinutes = 5
	}
	breakEntry := widget.NewEntry()
	breakEntry.SetText(strconv.Itoa(breakMinutes))
	breakEntry.Validator = func(s string) error {
		if s == "" {
			return nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if n < 1 || n > 60 {
			return strconv.ErrRange
		}
		return nil
	}

	// --- Notification Style ---
	notifyRadio := widget.NewRadioGroup(
		[]string{"系统通知", "弹窗提醒", "两者都用"},
		nil,
	)
	switch cfg.NotificationStyle {
	case config.NotifySystem:
		notifyRadio.SetSelected("系统通知")
	case config.NotifyPopup:
		notifyRadio.SetSelected("弹窗提醒")
	default:
		notifyRadio.SetSelected("两者都用")
	}
	notifyRadio.Horizontal = true

	// --- Checkboxes ---
	launchCheck := widget.NewCheck("开机自动启动", nil)
	launchCheck.SetChecked(cfg.LaunchAtLogin)

	soundCheck := widget.NewCheck("通知时播放提示音", nil)
	soundCheck.SetChecked(cfg.SoundEnabled)

	// --- Buttons ---
	saveBtn := widget.NewButton("保存", func() {
		breakMin, err := strconv.Atoi(breakEntry.Text)
		if err != nil || breakMin < 1 {
			breakMin = 5
		}
		if breakMin > 60 {
			breakMin = 60
		}

		intervalMin, err := strconv.Atoi(intervalSelect.Selected)
		if err != nil || intervalMin < 1 {
			intervalMin = 30
		}

		newCfg := config.Config{
			IntervalMinutes:   intervalMin,
			BreakDurationSec:  breakMin * 60,
			NotificationStyle: radioToStyle(notifyRadio.Selected),
			LaunchAtLogin:     launchCheck.Checked,
			SoundEnabled:      soundCheck.Checked,
		}

		_ = autostart.SetEnabled(newCfg.LaunchAtLogin)

		if sw.onSave != nil {
			sw.onSave(newCfg)
		}
		w.Hide()
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton("取消", func() {
		w.Hide()
	})

	// --- Form using Fyne's native Form widget ---
	form := widget.NewForm(
		widget.NewFormItem("提醒间隔", container.NewBorder(nil, nil, nil,
			widget.NewLabel("分钟"), intervalSelect)),
		widget.NewFormItem("休息时长", container.NewBorder(nil, nil, nil,
			widget.NewLabel("分钟"), breakEntry)),
		widget.NewFormItem("通知方式", notifyRadio),
	)

	sep := widget.NewSeparator()

	content := container.NewVBox(
		container.NewPadded(form),
		sep,
		container.NewPadded(container.NewVBox(
			container.NewHBox(widget.NewIcon(theme.SettingsIcon()), widget.NewLabel("其他")),
			launchCheck,
			soundCheck,
		)),
		layout.NewSpacer(),
		widget.NewSeparator(),
		container.NewPadded(container.NewHBox(
			layout.NewSpacer(), saveBtn, cancelBtn,
		)),
	)

	w.SetContent(content)
	w.SetCloseIntercept(func() {
		w.Hide()
	})

	return sw
}

func (sw *SettingsWindow) Show() {
	sw.window.Show()
	sw.window.RequestFocus()
}

func radioToStyle(sel string) config.NotificationStyle {
	switch sel {
	case "系统通知":
		return config.NotifySystem
	case "弹窗提醒":
		return config.NotifyPopup
	default:
		return config.NotifyBoth
	}
}

package notification

import (
	"fyne.io/fyne/v2"
	"github.com/zhangxinyu/breakreminder/internal/activities"
	"github.com/zhangxinyu/breakreminder/internal/config"
)

type Notifier struct {
	fyneApp   fyne.App
	style     config.NotificationStyle
	showPopup func(activity activities.Activity)
}

func New(fyneApp fyne.App, style config.NotificationStyle, showPopup func(activities.Activity)) *Notifier {
	return &Notifier{
		fyneApp:   fyneApp,
		style:     style,
		showPopup: showPopup,
	}
}

func (n *Notifier) SetStyle(style config.NotificationStyle) {
	n.style = style
}

func (n *Notifier) Notify(activity activities.Activity) {
	switch n.style {
	case config.NotifySystem:
		n.sendSystem(activity)
	case config.NotifyPopup:
		n.sendPopup(activity)
	case config.NotifyBoth:
		n.sendSystem(activity)
		n.sendPopup(activity)
	}
}

func (n *Notifier) sendSystem(activity activities.Activity) {
	notif := fyne.NewNotification("该休息一下了！", activity.Name+"："+activity.Description)
	n.fyneApp.SendNotification(notif)
}

func (n *Notifier) sendPopup(activity activities.Activity) {
	if n.showPopup != nil {
		n.showPopup(activity)
	}
}

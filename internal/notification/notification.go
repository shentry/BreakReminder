package notification

import (
	"github.com/zhangxinyu/breakreminder/internal/activities"
	"github.com/zhangxinyu/breakreminder/internal/config"
)

type Notifier struct {
	style                  config.NotificationStyle
	sendSystemNotification func(title, body string)
	showPopup              func(activity activities.Activity)
}

func New(style config.NotificationStyle, sendSystemNotification func(string, string), showPopup func(activities.Activity)) *Notifier {
	return &Notifier{
		style:                  style,
		sendSystemNotification: sendSystemNotification,
		showPopup:              showPopup,
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
	if n.sendSystemNotification != nil {
		n.sendSystemNotification("该休息一下了！", activity.Name+"："+activity.Description)
	}
}

func (n *Notifier) sendPopup(activity activities.Activity) {
	if n.showPopup != nil {
		n.showPopup(activity)
	}
}

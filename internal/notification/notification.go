package notification

import (
	"os/exec"
	"runtime"

	"github.com/zhangxinyu/breakreminder/internal/activities"
	"github.com/zhangxinyu/breakreminder/internal/config"
)

type Notifier struct {
	style     config.NotificationStyle
	showPopup func(activity activities.Activity)
}

func New(style config.NotificationStyle, showPopup func(activities.Activity)) *Notifier {
	return &Notifier{
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
	title := "该休息一下了！"
	body := activity.Name + "：" + activity.Description
	if runtime.GOOS == "darwin" {
		script := `display notification "` + escapeAppleScript(body) + `" with title "` + escapeAppleScript(title) + `" sound name "Glass"`
		_ = exec.Command("osascript", "-e", script).Start()
	}
}

// escapeAppleScript escapes double quotes and backslashes for AppleScript strings.
func escapeAppleScript(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		if s[i] == '"' || s[i] == '\\' {
			out = append(out, '\\')
		}
		out = append(out, s[i])
	}
	return string(out)
}

func (n *Notifier) sendPopup(activity activities.Activity) {
	if n.showPopup != nil {
		n.showPopup(activity)
	}
}

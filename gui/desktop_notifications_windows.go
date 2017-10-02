package gui

import (
	"os/exec"
	"syscall"

	"github.com/coyim/coyim/ui"
)

const notificationFeaturesSupported = notificationStyles

type desktopNotifications struct {
	notificationStyle   string
	notificationUrgent  bool
	notificationExpires bool
}

func newDesktopNotifications() *desktopNotifications {
	return createDesktopNotifications()
}

func (dn *desktopNotifications) show(jid, from, message string) error {
	from = ui.EscapeAllHTMLTags(string(ui.StripSomeHTML([]byte(from))))
	summary, body := dn.format(from, message, false)

	notification := Notification{
		Title:   "CoyIM",
		Message: summary + body,
		Icon:    coyimIcon.getPath(),
	}
	return notification.Popup()
}

type Notification struct {
	Title   string
	Message string
	Icon    string
}

func (n *Notification) Popup() error {
	cmd := exec.Command("toast.exe", "-t", n.Title, "-m", n.Message, "-p", n.Icon)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}

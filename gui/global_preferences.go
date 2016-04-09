package gui

import (
	"strings"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

// Debugging:
// - RawLogFile

type settingsPanel struct {
	b                   *builder
	dialog              gtki.Dialog
	notebook            gtki.Notebook
	singleWindow        gtki.CheckButton
	renderSlashMe       gtki.CheckButton
	notificationsType   gtki.ComboBox
	urgentNotifications gtki.CheckButton
	expireNotifications gtki.CheckButton
	notificationCommand gtki.Entry
	notificationTimeout gtki.SpinButton
	rawLogFile          gtki.Entry
}

func createSettingsPanel() *settingsPanel {
	p := &settingsPanel{b: newBuilder("GlobalPreferences")}
	p.b.getItems(
		"GlobalPreferences", &p.dialog,
		"notebook1", &p.notebook,
		"singleWindow", &p.singleWindow,
		"slashMe", &p.renderSlashMe,
		"notificationsType", &p.notificationsType,
		"notificationUrgent", &p.urgentNotifications,
		"notificationExpires", &p.expireNotifications,
		"notificationCommand", &p.notificationCommand,
		"notificationTimeout", &p.notificationTimeout,
		"rawLogFile", &p.rawLogFile,
	)

	return p
}

var notificationsTypes = []string{
	"off",
	"only-presence-of-new-information",
	"with-author-but-no-content",
	"with-content",
}

func indexOfNotificationType(s string) int {
	for ix, v := range notificationsTypes {
		if s == v {
			return ix
		}
	}
	return -1
}

func valOr(v, def int) int {
	if v == 0 {
		return def
	}
	return v
}

func notifyCommand(cmd []string) string {
	return strings.Join(cmd, " ")
}

func (u *gtkUI) showGlobalPreferences() {
	settings := u.settings
	config := u.config

	panel := createSettingsPanel()

	panel.singleWindow.SetActive(settings.GetSingleWindow())
	panel.renderSlashMe.SetActive(settings.GetSlashMe())
	panel.urgentNotifications.SetActive(settings.GetNotificationUrgency())
	panel.expireNotifications.SetActive(settings.GetNotificationExpires())
	panel.notificationsType.SetActive(indexOfNotificationType(settings.GetNotificationStyle()))

	if config != nil {
		panel.notificationCommand.SetText(notifyCommand(config.NotifyCommand))
		panel.notificationTimeout.SetValue(float64(valOr(config.IdleSecondsBeforeNotification, 60)))
		panel.rawLogFile.SetText(config.RawLogFile)
	} else {
		// turn them off or something
	}

	panel.b.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
		},
		"on_cancel_signal": func() {
			panel.dialog.Destroy()
		},
	})

	panel.dialog.SetTransientFor(u.window)
	panel.dialog.ShowAll()
	panel.notebook.SetCurrentPage(0)
}

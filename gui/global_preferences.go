package gui

import (
	"strings"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type settingsPanel struct {
	b                        *builder
	dialog                   gtki.Dialog
	notebook                 gtki.Notebook
	singleWindow             gtki.CheckButton
	renderSlashMe            gtki.CheckButton
	sendWithShiftEnter       gtki.CheckButton
	notificationsType        gtki.ComboBox
	urgentNotifications      gtki.CheckButton
	expireNotifications      gtki.CheckButton
	notificationCommand      gtki.Entry
	notificationTimeout      gtki.SpinButton
	rawLogFile               gtki.Entry
	notificationCommandLabel gtki.Label
	notificationTimeoutLabel gtki.Label
	rawLogFileLabel          gtki.Label
}

func createSettingsPanel() *settingsPanel {
	p := &settingsPanel{b: newBuilder("GlobalPreferences")}
	p.b.getItems(
		"GlobalPreferences", &p.dialog,
		"notebook1", &p.notebook,
		"singleWindow", &p.singleWindow,
		"slashMe", &p.renderSlashMe,
		"sendWithShiftEnter", &p.sendWithShiftEnter,
		"notificationsType", &p.notificationsType,
		"notificationUrgent", &p.urgentNotifications,
		"notificationExpires", &p.expireNotifications,
		"notificationCommand", &p.notificationCommand,
		"notificationTimeout", &p.notificationTimeout,
		"rawLogFile", &p.rawLogFile,
		"notificationCommandLabel", &p.notificationCommandLabel,
		"notificationTimeoutLabel", &p.notificationTimeoutLabel,
		"rawLogFileLabel", &p.rawLogFileLabel,
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

	orgSingleWindow := settings.GetSingleWindow()
	panel.singleWindow.SetActive(orgSingleWindow)

	orgSlashMe := settings.GetSlashMe()
	panel.renderSlashMe.SetActive(orgSlashMe)

	orgShiftEnter := settings.GetShiftEnterForSend()
	panel.sendWithShiftEnter.SetActive(orgShiftEnter)

	orgUrgentNot := settings.GetNotificationUrgency()
	panel.urgentNotifications.SetActive(orgUrgentNot)

	orgExpireNot := settings.GetNotificationExpires()
	panel.expireNotifications.SetActive(orgExpireNot)

	orgExpireType := settings.GetNotificationStyle()
	panel.notificationsType.SetActive(indexOfNotificationType(orgExpireType))

	if config != nil {
		panel.notificationCommand.SetText(notifyCommand(config.NotifyCommand))
		panel.notificationTimeout.SetValue(float64(valOr(config.IdleSecondsBeforeNotification, 60)))
		panel.rawLogFile.SetText(config.RawLogFile)
	} else {
		panel.notificationCommand.SetSensitive(false)
		panel.notificationCommandLabel.SetSensitive(false)
		panel.notificationTimeout.SetSensitive(false)
		panel.notificationTimeoutLabel.SetSensitive(false)
		panel.rawLogFile.SetSensitive(false)
		panel.rawLogFileLabel.SetSensitive(false)
	}

	panel.b.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			if newSingleWindow := panel.singleWindow.GetActive(); newSingleWindow != orgSingleWindow {
				settings.SetSingleWindow(newSingleWindow)
			}

			if newSlashMe := panel.renderSlashMe.GetActive(); newSlashMe != orgSlashMe {
				settings.SetSlashMe(newSlashMe)
			}

			if newShiftEnter := panel.sendWithShiftEnter.GetActive(); newShiftEnter != orgShiftEnter {
				settings.SetShiftEnterForSend(newShiftEnter)
			}

			if newUrgentNot := panel.urgentNotifications.GetActive(); newUrgentNot != orgUrgentNot {
				settings.SetNotificationUrgency(newUrgentNot)
			}

			if newExpireNot := panel.expireNotifications.GetActive(); newExpireNot != orgExpireNot {
				settings.SetNotificationExpires(newExpireNot)
			}

			newExpireType := panel.notificationsType.GetActive()
			if newExpireType >= 0 && newExpireType < len(notificationsTypes) && notificationsTypes[newExpireType] != orgExpireType {
				settings.SetNotificationStyle(notificationsTypes[newExpireType])
			}

			if config != nil {
				tx, _ := panel.notificationCommand.GetText()
				if strings.TrimSpace(tx) != "" {
					config.NotifyCommand = strings.Split(tx, " ")
				} else {
					config.NotifyCommand = nil
				}

				val := panel.notificationTimeout.GetValueAsInt()
				if val == 60 {
					val = 0
				}
				config.IdleSecondsBeforeNotification = val
				tx, _ = panel.rawLogFile.GetText()
				tx = strings.TrimSpace(tx)
				config.RawLogFile = tx
				u.saveConfigOnly()
				panel.dialog.Destroy()
			}
		},
		"on_cancel_signal": func() {
			panel.dialog.Destroy()
		},
	})

	panel.dialog.SetTransientFor(u.window)
	panel.dialog.ShowAll()
	panel.notebook.SetCurrentPage(0)
}

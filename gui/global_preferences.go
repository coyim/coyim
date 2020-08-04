package gui

import (
	"strings"

	"github.com/coyim/gotk3adapter/gtki"
)

type settingsPanel struct {
	b                        *builder
	dialog                   gtki.Dialog      `gtk-widget:"GlobalPreferences"`
	notebook                 gtki.Notebook    `gtk-widget:"notebook1"`
	singleWindow             gtki.CheckButton `gtk-widget:"singleWindow"`
	showEmptyGroups          gtki.CheckButton `gtk-widget:"showEmptyGroups"`
	sendWithShiftEnter       gtki.CheckButton `gtk-widget:"sendWithShiftEnter"`
	emacsKeyboard            gtki.CheckButton `gtk-widget:"emacsKeyboard"`
	notificationsType        gtki.ComboBox    `gtk-widget:"notificationsType"`
	urgentNotifications      gtki.CheckButton `gtk-widget:"notificationUrgent"`
	expireNotifications      gtki.CheckButton `gtk-widget:"notificationExpires"`
	notificationCommand      gtki.Entry       `gtk-widget:"notificationCommand"`
	notificationTimeout      gtki.SpinButton  `gtk-widget:"notificationTimeout"`
	rawLogFile               gtki.Entry       `gtk-widget:"rawLogFile"`
	notificationCommandLabel gtki.Label       `gtk-widget:"notificationCommandLabel"`
	notificationTimeoutLabel gtki.Label       `gtk-widget:"notificationTimeoutLabel"`
	rawLogFileLabel          gtki.Label       `gtk-widget:"rawLogFileLabel"`
}

func createSettingsPanel() *settingsPanel {
	p := &settingsPanel{b: newBuilder("GlobalPreferences")}
	panicOnDevError(p.b.bindObjects(p))
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

	orgShiftEnter := settings.GetShiftEnterForSend()
	panel.sendWithShiftEnter.SetActive(orgShiftEnter)

	emacsKeyBindings := settings.GetEmacsKeyBindings()
	panel.emacsKeyboard.SetActive(emacsKeyBindings)

	orgShowEmptyGroups := settings.GetShowEmptyGroups()
	panel.showEmptyGroups.SetActive(orgShowEmptyGroups)

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
		"on_save": func() {
			if newSingleWindow := panel.singleWindow.GetActive(); newSingleWindow != orgSingleWindow {
				settings.SetSingleWindow(newSingleWindow)
			}

			if newShiftEnter := panel.sendWithShiftEnter.GetActive(); newShiftEnter != orgShiftEnter {
				settings.SetShiftEnterForSend(newShiftEnter)
			}

			if newEmacsKeyBindings := panel.emacsKeyboard.GetActive(); newEmacsKeyBindings != emacsKeyBindings {
				settings.SetEmacsKeyBindings(newEmacsKeyBindings)
				u.keyboardSettings.emacs = settings.GetEmacsKeyBindings()
				u.keyboardSettings.update()
			}

			if newShowEmptyGroups := panel.showEmptyGroups.GetActive(); newShowEmptyGroups != orgShowEmptyGroups {
				settings.SetShowEmptyGroups(newShowEmptyGroups)
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
			}
			panel.dialog.Destroy()
			u.roster.redraw()
			u.deNotify.updateWith(settings)
			u.updateUnifiedOrNot()
		},
		"on_cancel": func() {
			panel.dialog.Destroy()
		},
	})

	panel.dialog.SetTransientFor(u.window)
	panel.dialog.ShowAll()

	if notificationFeaturesSupported&notificationUrgency == 0 {
		panel.b.get("notificationUrgencyInstructions").(gtki.Widget).SetVisible(false)
		panel.b.get("notificationUrgentLabel").(gtki.Widget).SetVisible(false)
		panel.urgentNotifications.SetVisible(false)
	}
	if notificationFeaturesSupported&notificationExpiry == 0 {
		panel.b.get("notificationExpiryInstructions").(gtki.Widget).SetVisible(false)
		panel.b.get("notificationExpiresLabel").(gtki.Widget).SetVisible(false)
		panel.expireNotifications.SetVisible(false)
	}
	if notificationFeaturesSupported&notificationStyles == 0 {
		panel.b.get("notificationTypeLabel").(gtki.Widget).SetVisible(false)
		panel.b.get("notificationTypeInstructions").(gtki.Widget).SetVisible(false)
		panel.notificationsType.SetVisible(false)
	}

	panel.notebook.SetCurrentPage(0)
}

package gui

import (
	"errors"
	"strings"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type changePasswordData struct {
	builder               *builder
	dialog                gtki.Dialog      `gtk-widget:"ChangePassword"`
	formBox               gtki.Box         `gtk-widget:"form-box"`
	messagesBox           gtki.Box         `gtk-widget:"messages-box"`
	passwordEntry         gtki.Entry       `gtk-widget:"new-password-entry"`
	repeatPasswordEntry   gtki.Entry       `gtk-widget:"repeat-password-entry"`
	formBoxLabel          gtki.Label       `gtk-widget:"form-box-label"`
	callbackLabel         gtki.Label       `gtk-widget:"callback-label"`
	formImage             gtki.Image       `gtk-widget:"form-image"`
	callbackImage         gtki.Image       `gtk-widget:"callback-image"`
	changePasswordSpinner gtki.Spinner     `gtk-widget:"change-password-spinner"`
	callbackGrid          gtki.Grid        `gtk-widget:"callback-grid"`
	buttonChange          gtki.Button      `gtk-widget:"button-change"`
	buttonCancel          gtki.Button      `gtk-widget:"button-cancel"`
	buttonOk              gtki.Button      `gtk-widget:"button-ok"`
	checkboxSavePassword  gtki.CheckButton `gtk-widget:"save-new-password-checkbox"`
}

func getBuilderAndChangePasswordData() *changePasswordData {
	data := &changePasswordData{}

	dialogID := "ChangePassword"
	data.builder = newBuilder(dialogID)
	panicOnDevError(data.builder.bindObjects(data))

	return data
}

func validateNewPassword(newPassword, repeatedPassword string) error {
	if newPassword != repeatedPassword {
		return errors.New(i18n.Local("The passwords do not match"))
	} else if newPassword == "" {
		return errors.New(i18n.Local("The password can't be empty"))
	}

	return nil
}

func changePassword(account *account, newPassword string, u *gtkUI, data *changePasswordData) {
	accountInfo := account.Account()
	accountInfoParts := strings.SplitN(accountInfo, "@", 2)
	username := accountInfoParts[0]
	server := accountInfoParts[1]

	if err := account.session.Conn().ChangePassword(username, server, newPassword); err == nil {
		data.changePasswordSpinner.Stop()

		config := account.session.GetConfig()
		saveNewPassword := data.checkboxSavePassword.GetActive()
		if saveNewPassword {
			config.Password = newPassword
			u.SaveConfig()
		}

		data.formBox.Hide()
		data.callbackGrid.Show()
		data.callbackGrid.SetMarginTop(35)
		data.callbackLabel.SetText(i18n.Localf("Password changed successfully for %s.", config.Account))
		setImageFromFile(data.callbackImage, "success.svg")
		data.buttonOk.Show()
	} else {
		data.formBox.Hide()
		data.callbackGrid.Show()
		data.callbackLabel.SetText(i18n.Localf("Password change failed.\n Error: %s", err.Error()))
		setImageFromFile(data.callbackImage, "failure.svg")
	}
}

func (u *gtkUI) buildChangePasswordDialog(account *account) {
	assertInUIThread()

	data := getBuilderAndChangePasswordData()

	data.builder.ConnectSignals(map[string]interface{}{
		"on_ok":     data.dialog.Destroy,
		"on_cancel": data.dialog.Destroy,
		"on_change": func() {
			newPassword, _ := data.passwordEntry.GetText()
			repeatedPassword, _ := data.repeatPasswordEntry.GetText()

			data.formBoxLabel.SetText(i18n.Local(""))

			if err := validateNewPassword(newPassword, repeatedPassword); err != nil {
				data.formBoxLabel.Show()
				data.formBoxLabel.SetText(i18n.Localf("Error: %s", err.Error()))
				setImageFromFile(data.formImage, "failure.svg")
			} else {
				data.formImage.Hide()
				data.changePasswordSpinner.Start()
				data.formBoxLabel.Show()
				data.messagesBox.SetMarginTop(35)
				data.formBoxLabel.SetText(i18n.Local("Attempting to change password..."))
				data.buttonChange.Hide()
				data.buttonCancel.Hide()

				data.passwordEntry.SetCanFocus(false)
				data.repeatPasswordEntry.SetCanFocus(false)
				go changePassword(account, newPassword, u, data)
			}
		},
	})

	data.dialog.SetTransientFor(u.window)
	data.dialog.ShowAll()
	data.callbackGrid.Hide()
	data.formBoxLabel.Hide()
	data.buttonOk.Hide()
}

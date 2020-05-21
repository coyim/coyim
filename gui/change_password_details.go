package gui

import (
	"errors"
	"strings"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type changePasswordData struct {
	builder               *builder
	dialog                gtki.Dialog
	formBox               gtki.Box
	messagesBox           gtki.Box
	passwordEntry         gtki.Entry
	repeatPasswordEntry   gtki.Entry
	formBoxLabel          gtki.Label
	callbackLabel         gtki.Label
	formImage             gtki.Image
	callbackImage         gtki.Image
	changePasswordSpinner gtki.Spinner
	callbackGrid          gtki.Grid
	buttonChange          gtki.Button
	buttonCancel          gtki.Button
	buttonOk              gtki.Button
	checkboxSavePassword  gtki.CheckButton
}

func getBuilderAndChangePasswordData() *changePasswordData {
	data := &changePasswordData{}

	dialogID := "ChangePassword"
	data.builder = newBuilder(dialogID)

	data.builder.getItems(
		dialogID, &data.dialog,
		"form-box", &data.formBox,
		"messages-box", &data.messagesBox,
		"new-password-entry", &data.passwordEntry,
		"repeat-password-entry", &data.repeatPasswordEntry,
		"form-box-label", &data.formBoxLabel,
		"callback-label", &data.callbackLabel,
		"form-image", &data.formImage,
		"callback-image", &data.callbackImage,
		"change-password-spinner", &data.changePasswordSpinner,
		"callback-grid", &data.callbackGrid,
		"button-change", &data.buttonChange,
		"button-cancel", &data.buttonCancel,
		"button-ok", &data.buttonOk,
		"save-new-password-checkbox", &data.checkboxSavePassword,
	)

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
	accountInfo := account.session.GetConfig().Account
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
		"on_ok_signal":     data.dialog.Destroy,
		"on_cancel_signal": data.dialog.Destroy,
		"on_change_signal": func() {
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

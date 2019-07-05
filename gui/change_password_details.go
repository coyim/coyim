package gui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type changePasswordData struct {
	builder               *builder
	dialog                gtki.Dialog
	formBox               gtki.Box
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
}

func getBuilderAndChangePasswordData() *changePasswordData {
	data := &changePasswordData{}

	dialogID := "ChangePassword"
	data.builder = newBuilder(dialogID)

	data.builder.getItems(
		dialogID, &data.dialog,
		"form-box", &data.formBox,
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
	)

	return data
}

func validateNewPassword(newPassword, repeatedPassword string) error {
	var err error

	if newPassword == "" || repeatedPassword == "" {
		err = errors.New("The field can't be empty")
	} else {
		if newPassword != repeatedPassword {
			err = errors.New("The passwords do not match")
		}
	}

	return err
}

// TODO: refactor and change me
func changePassword(account *account, newPassword string, u *gtkUI, data *changePasswordData) {

	//NOTE: This block is commented if we're using just one stream for change the password.
	// Prefer using cached password if present
	// if account.cachedPassword != "" {
	// 	oldPassword = account.cachedPassword
	// } else {
	// 	oldPassword = account.session.GetConfig().Password
	// }

	accountInfo := account.session.GetConfig().Account
	accountInfoParts := strings.SplitN(accountInfo, "@", 2) // Get the username and server domain

	if err := account.session.Conn().ChangePassword2(accountInfoParts[0], accountInfoParts[1], newPassword); err == nil {
		data.changePasswordSpinner.Stop()
		// Clear old password and cached password on successful change.
		// We only save new password, if the user wishes to save it at the re-login.
		account.session.GetConfig().Password = ""
		u.SaveConfig()
		data.formBox.Hide()
		data.callbackGrid.Show()
		data.callbackLabel.SetText("Password changed successfully")
		setImageFromFile(data.callbackImage, "success.svg")
		data.buttonOk.Show()
	} else {
		data.formBox.Hide()
		data.callbackGrid.Show()
		data.callbackLabel.SetText(fmt.Sprintf("Password change failed.\n Error: %s", err.Error()))
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
				data.formBoxLabel.SetText(i18n.Local(err.Error()))
				setImageFromFile(data.formImage, "failure.svg")
			} else {
				data.changePasswordSpinner.Start()
				data.formBoxLabel.Show()
				data.formBoxLabel.SetText(i18n.Local("Attempting to password change"))
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

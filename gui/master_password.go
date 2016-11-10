package gui

import (
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
	"github.com/twstrike/gotk3adapter/gtki"
)

func (u *gtkUI) captureInitialMasterPassword(k func(), onCancel func()) {
	dialogID := "CaptureInitialMasterPassword"
	builder := newBuilder(dialogID)
	dialogOb := builder.getObj(dialogID)
	pwdDialog := dialogOb.(gtki.Dialog)

	passObj := builder.getObj("password")
	password := passObj.(gtki.Entry)

	pass2Obj := builder.getObj("password2")
	password2 := pass2Obj.(gtki.Entry)

	msgObj := builder.getObj("passMessage")
	messageObj := msgObj.(gtki.Label)
	messageObj.SetSelectable(true)

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			passText1, _ := password.GetText()
			passText2, _ := password2.GetText()
			if len(passText1) == 0 {
				messageObj.SetLabel(i18n.Local("Password can not be empty - please try again"))
				password.GrabFocus()
			} else if passText1 != passText2 {
				messageObj.SetLabel(i18n.Local("Passwords have to be the same - please try again"))
				password.GrabFocus()
			} else {
				u.keySupplier = &onetimeSavedPassword{
					savedPassword: passText1,
					realF:         u.keySupplier,
				}
				pwdDialog.Destroy()
				k()
			}
		},
		"on_cancel_signal": func() {
			pwdDialog.Destroy()
			onCancel()
		},
	})

	doInUIThread(func() {
		pwdDialog.SetTransientFor(u.window)
		pwdDialog.ShowAll()
	})
}

func (u *gtkUI) wouldYouLikeToEncryptYourFile(k func(bool)) {
	dialogID := "AskToEncrypt"
	builder := newBuilder(dialogID)

	dialogOb := builder.getObj(dialogID)
	encryptDialog := dialogOb.(gtki.MessageDialog)
	encryptDialog.SetDefaultResponse(gtki.RESPONSE_YES)
	encryptDialog.SetTransientFor(u.window)

	responseType := gtki.ResponseType(encryptDialog.Run())
	result := responseType == gtki.RESPONSE_YES
	encryptDialog.Destroy()
	k(result)
}

type onetimeSavedPassword struct {
	savedPassword string
	realF         config.KeySupplier
}

func (o *onetimeSavedPassword) Invalidate() {
	o.realF.Invalidate()
}

func (o *onetimeSavedPassword) LastAttemptFailed() {
	o.realF.LastAttemptFailed()
}

func (o *onetimeSavedPassword) GenerateKey(params config.EncryptionParameters) ([]byte, []byte, bool) {
	if o.savedPassword != "" {
		ourPwd := o.savedPassword
		o.savedPassword = ""

		l, r := config.GenerateKeys(ourPwd, params)
		return l, r, true
	}
	return o.realF.GenerateKey(params)
}

func (u *gtkUI) getMasterPassword(params config.EncryptionParameters, lastAttemptFailed bool) ([]byte, []byte, bool) {
	dialogID := "MasterPassword"
	pwdResultChan := make(chan string)
	var cleanup func()

	doInUIThread(func() {
		builder := newBuilder(dialogID)
		dialogOb := builder.getObj(dialogID)
		dialog := dialogOb.(gtki.Dialog)

		cleanup = dialog.Destroy

		passObj := builder.getObj("password")
		password := passObj.(gtki.Entry)

		msgObj := builder.getObj("passMessage")
		messageObj := msgObj.(gtki.Label)
		messageObj.SetSelectable(true)

		if lastAttemptFailed {
			messageObj.SetLabel(i18n.Local("Incorrect password entered, please try again."))
		}

		hadSubmission := false

		builder.ConnectSignals(map[string]interface{}{
			"on_save_signal": func() {
				if !hadSubmission {
					passText, _ := password.GetText()
					if len(passText) > 0 {
						hadSubmission = true
						messageObj.SetLabel(i18n.Local("Checking password..."))
						pwdResultChan <- passText
						close(pwdResultChan)
					}
				}
			},
			"on_cancel_signal": func() {
				if !hadSubmission {
					hadSubmission = true
					close(pwdResultChan)
					u.quit()
				}
			},
		})

		dialog.SetTransientFor(u.window)
		dialog.ShowAll()
	})

	pwd, ok := <-pwdResultChan

	if !ok {
		doInUIThread(cleanup)
		return nil, nil, false
	}

	l, r := config.GenerateKeys(pwd, params)
	doInUIThread(cleanup)
	return l, r, true
}

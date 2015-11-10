package gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
)

func (u *gtkUI) wouldYouLikeToEncryptYourFile(k func(bool)) {
	encryptDialog := gtk.MessageDialogNew(
		nil,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_QUESTION,
		gtk.BUTTONS_YES_NO,
		i18n.Local("Would you like to save your configuration file in an encrypted format? This can be significantly more secure, but you will not be able to access your configuration if you lose the password. You will only be asked for your password when CoyIM starts."),
	)
	encryptDialog.SetTitle(i18n.Local("Encrypt configuration file"))
	encryptDialog.SetDefaultResponse(gtk.RESPONSE_YES)

	responseType := gtk.ResponseType(encryptDialog.Run())
	switch responseType {
	case gtk.RESPONSE_YES:
		k(true)
	case gtk.RESPONSE_NO:
		k(false)
	default:
		k(false)
	}
	encryptDialog.Destroy()
}

func (u *gtkUI) getMasterPassword(params config.EncryptionParameters) ([]byte, []byte, bool) {
	vars := make(map[string]string)
	vars["$title"] = i18n.Local("Enter master password")
	vars["$passwordMessage"] = i18n.Local("Please enter the master password for the configuration file. You will not be asked for this password again until you restart CoyIM.")
	vars["$saveLabel"] = i18n.Local("OK")
	vars["$cancelLabel"] = i18n.Local("Cancel")

	builder, _ := loadBuilderWith("MasterPasswordDefinition", vars)
	dialogOb, _ := builder.GetObject("MasterPassword")
	dialog := dialogOb.(*gtk.Dialog)
	passObj, _ := builder.GetObject("password")
	password := passObj.(*gtk.Entry)
	pwdResultChan := make(chan string)

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			passText, _ := password.GetText()
			pwdResultChan <- passText
			close(pwdResultChan)
			dialog.Destroy()
		},
		"on_cancel_signal": func() {
			close(pwdResultChan)
			dialog.Destroy()
		},
	})

	glib.IdleAdd(func() {
		dialog.SetTransientFor(u.window)
		dialog.ShowAll()
	})

	pwd, ok := <-pwdResultChan

	if !ok {
		return nil, nil, false
	}

	l, r := config.GenerateKeys(pwd, params)
	return l, r, true
}

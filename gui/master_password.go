package gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
)

func (u *gtkUI) wouldYouLikeToEncryptYourFile(k func(bool)) {
	dialogId := "AskToEncrypt"
	builder, _ := loadBuilderWith(dialogId, nil)

	dialogOb, _ := builder.GetObject(dialogId)
	encryptDialog := dialogOb.(*gtk.MessageDialog)
	encryptDialog.SetDefaultResponse(gtk.RESPONSE_YES)
	encryptDialog.SetTransientFor(u.window)

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

	builder, _ := loadBuilderWith("MasterPassword", vars)
	dialogOb, _ := builder.GetObject("MasterPassword")
	dialog := dialogOb.(*gtk.Dialog)
	passObj, _ := builder.GetObject("password")
	password := passObj.(*gtk.Entry)
	pwdResultChan := make(chan string)

	abort := func() {
		close(pwdResultChan)
		dialog.Destroy()
		u.quit()
	}
	dialog.Connect("close", abort)

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			passText, _ := password.GetText()
			pwdResultChan <- passText
			close(pwdResultChan)
			dialog.Destroy()
		},
		"on_cancel_signal": abort,
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

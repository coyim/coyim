package gui

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
)

func (u *gtkUI) wouldYouLikeToEncryptYourFile(k func(bool)) {
	dialogID := "AskToEncrypt"
	builder := builderForDefinition(dialogID)

	dialogOb, _ := builder.GetObject(dialogID)
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
	dialogID := "MasterPassword"
	builder := builderForDefinition(dialogID)
	dialogOb, _ := builder.GetObject(dialogID)
	dialog := dialogOb.(*gtk.Dialog)

	passObj, _ := builder.GetObject("password")
	password := passObj.(*gtk.Entry)
	pwdResultChan := make(chan string)

	msgObj, _ := builder.GetObject("passMessage")
	messageObj := msgObj.(*gtk.Label)

	builder.ConnectSignals(map[string]interface{}{
		"on_save_signal": func() {
			passText, _ := password.GetText()
			if len(passText) > 0 {
				messageObj.SetLabel(i18n.Local("Checking password..."))
				pwdResultChan <- passText
				close(pwdResultChan)
			}
		},
		"on_cancel_signal": func() {
			close(pwdResultChan)
			u.quit()
		},
	})

	glib.IdleAdd(func() {
		dialog.SetTransientFor(u.window)
		dialog.ShowAll()
	})

	pwd, ok := <-pwdResultChan

	if !ok {
		glib.IdleAdd(func() {
			dialog.Destroy()
		})
		return nil, nil, false
	}

	l, r := config.GenerateKeys(pwd, params)
	glib.IdleAdd(func() {
		dialog.Destroy()
	})
	return l, r, true
}

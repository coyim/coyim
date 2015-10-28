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
	reg := createWidgetRegistry()

	pwdResultChan := make(chan string)

	dialog := dialog{
		title:       i18n.Local("Enter master password"),
		position:    gtk.WIN_POS_CENTER,
		id:          "masterPasswordDialog",
		defaultSize: []int{300, 0},
		content: []creatable{
			label{
				text:      i18n.Local("Please enter the master password for the configuration file. You will not be asked for this password again until you restart CoyIM."),
				wrapLines: true,
			},
			entry{
				editable:   true,
				visibility: false,
				id:         "password",
				focused:    true,
				onActivate: func() {
					pwdResultChan <- reg.getText("password")
					close(pwdResultChan)
					reg.dialogDestroy("masterPasswordDialog")
				},
			},
			hbox{
				fromRight: true,
				content: []creatable{
					button{
						text: i18n.Local("OK"),
						onClicked: func() {
							pwdResultChan <- reg.getText("password")
							close(pwdResultChan)
							reg.dialogDestroy("masterPasswordDialog")
						},
					},
					button{
						text: i18n.Local("Cancel"),
						onClicked: func() {
							close(pwdResultChan)
							reg.dialogDestroy("masterPasswordDialog")
						},
					},
				},
			},
		},
	}

	glib.IdleAdd(func() {
		deg, _ := dialog.create(reg)
		deg.(*gtk.Dialog).SetTransientFor(u.window)
		reg.dialogShowAll("masterPasswordDialog")
	})

	pwd, ok := <-pwdResultChan

	if !ok {
		return nil, nil, false
	}

	l, r := config.GenerateKeys(pwd, params)
	return l, r, true
}

package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoom struct {
	errorBox *errorNotification
	builder  *builder

	gtki.Dialog `gtk-widget:"create-chat-dialog"`

	notification         gtki.InfoBar
	notificationArea     gtki.Box          `gtk-widget:"notification-area"`
	form                 gtki.Grid         `gtk-widget:"form"`
	account              gtki.ComboBox     `gtk-widget:"accounts"`
	chatServices         gtki.ComboBoxText `gtk-widget:"chatServices"`
	chatServiceEntry     gtki.Entry        `gtk-widget:"chatServiceEntry"`
	room                 gtki.Entry        `gtk-widget:"room"`
	cancelButton         gtki.Button       `gtk-widget:"button-cancel"`
	createButton         gtki.Button       `gtk-widget:"button-ok"`
	createButtonPrevText string

	cancel chan bool

	u *gtkUI
}

func (v *createMUCRoom) notifyOnError(err string) {
	if v.notification != nil {
		v.notificationArea.Remove(v.notification)
	}

	v.errorBox.ShowMessage(err)
}

func (v *createMUCRoom) clearErrors() {
	v.errorBox.Hide()
}

func (u *gtkUI) onNoAccountsConnected(v *createMUCRoom) {
	v.chatServices.RemoveAll()
}

func (u *gtkUI) updateServicesBasedOnAccount(v *createMUCRoom, acc *account) {
	doInUIThread(v.clearErrors)
	go v.updateChatServices(acc)
}

func (u *gtkUI) newMUCRoomView() *createMUCRoom {
	view := &createMUCRoom{u: u}

	view.builder = newBuilder("MUCCreateRoom")
	panicOnDevError(view.builder.bindObjects(view))
	view.errorBox = newErrorNotification(view.notificationArea)

	accountsInput := view.builder.get("accounts").(gtki.ComboBox)
	ac := u.createConnectedAccountsComponent(accountsInput, view,
		func(acc *account) {
			u.updateServicesBasedOnAccount(view, acc)
		},
		func() {
			u.onNoAccountsConnected(view)
		},
	)

	view.builder.ConnectSignals(map[string]interface{}{
		"on_create_room": func() {
			go view.createRoomHandler(ac.currentAccount())
		},
		"on_cancel":       view.Destroy,
		"on_close_window": ac.onDestroy,
		"on_room_changed": func() {
			view.disableCreationIfFieldsAreEmpty(ac.currentAccount())
		},
		"on_chatServiceEntry_changed": func() {
			view.disableCreationIfFieldsAreEmpty(ac.currentAccount())
		},
	})

	return view
}

func (v *createMUCRoom) updateChatServices(ac *account) {
	enteredService, _ := v.chatServiceEntry.GetText()
	v.chatServices.RemoveAll()

	items, err := ac.session.GetChatServices(jid.Parse(ac.Account()).Host())

	if err != nil {
		v.u.log.WithError(err).Debug("something went wrong trying to get chat services")
		return
	}

	for _, i := range items {
		v.chatServices.AppendText(i.Jid)
	}

	if len(items) > 0 && enteredService == "" {
		v.chatServices.SetActive(0)
	}
}

func (v *createMUCRoom) updateFields(f bool) {
	v.cancelButton.SetSensitive(f)
	v.createButton.SetSensitive(f)
	v.account.SetSensitive(f)
	v.room.SetSensitive(f)
	v.chatServices.SetSensitive(f)
}

func (v *createMUCRoom) createRoomHandler(ac *account) {
	if ac == nil {
		v.errorBox.ShowMessage(i18n.Local("No account selected, please select one account from the list or connect to one."))
		return
	}

	roomName, _ := v.room.GetText()
	service := v.chatServices.GetActiveText()

	if roomName == "" || service == "" {
		v.errorBox.ShowMessage(i18n.Local("Please fill the required fields to create the room."))
		return
	}

	doInUIThread(func() {
		v.updateFields(false)
		v.createButtonPrevText, _ = v.createButton.GetLabel()
		_ = v.createButton.SetProperty("label", i18n.Local("Creating room..."))
	})

	v.cancel = make(chan bool, 1)

	ec := ac.session.CreateRoom(jid.Parse(fmt.Sprintf("%s@%s", roomName, service)).(jid.Bare))

	go func() {
		shouldUpdateUI := false
		isRoomCreated := false

		defer func() {
			if shouldUpdateUI {
				doInUIThread(func() {
					if isRoomCreated {
						v.errorBox.ShowMessage(i18n.Local("Room created with success"))
					} else {
						v.errorBox.ShowMessage(i18n.Local("Could not create the new room"))
					}
					v.updateFields(true)
					_ = v.createButton.SetProperty("label", v.createButtonPrevText)
				})
			}
		}()

		select {
		case err, ok := <-ec:
			if !ok {
				return
			}

			if err != nil {
				v.u.log.WithError(err).Debug("something went wrong trying to create the room")
			} else {
				isRoomCreated = true
			}
			shouldUpdateUI = true
			return
		case <-v.cancel:
			return
		}
	}()
}

func (v *createMUCRoom) allFieldsHaveContent(ac *account) bool {
	accountVal := ac.Account()
	serviceVal := v.chatServices.GetActiveText()
	roomVal, _ := v.room.GetText()

	return accountVal != "" && serviceVal != "" && roomVal != ""
}

func setEnabled(w gtki.Widget, enable bool) {
	w.SetSensitive(enable)
}

func (v *createMUCRoom) disableCreationIfFieldsAreEmpty(ac *account) {
	setEnabled(v.createButton, v.allFieldsHaveContent(ac))
}

func (u *gtkUI) mucCreateChatRoom() {
	view := u.newMUCRoomView()
	view.SetTransientFor(u.window)
	view.Show()
}

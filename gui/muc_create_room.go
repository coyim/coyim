package gui

import (
	"strings"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoom struct {
	errorBox *errorNotification
	builder  *builder

	dialog gtki.Dialog `gtk-widget:"create-chat-dialog"`

	notification          gtki.InfoBar
	notificationArea      gtki.Box          `gtk-widget:"notification-area"`
	account               gtki.ComboBox     `gtk-widget:"accounts"`
	chatServices          gtki.ComboBoxText `gtk-widget:"chatServices"`
	chatServiceEntry      gtki.Entry        `gtk-widget:"chatServiceEntry"`
	room                  gtki.Entry        `gtk-widget:"room"`
	cancelButton          gtki.Button       `gtk-widget:"button-cancel"`
	createButton          gtki.Button       `gtk-widget:"button-ok"`
	createButtonPrevText  string
	previousUpdateChannel chan bool

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
			view.clearErrors()
			go view.createRoomHandler(ac.currentAccount())
		},
		"on_cancel":       view.dialog.Destroy,
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
	if v.previousUpdateChannel != nil {
		v.previousUpdateChannel <- true
	}

	v.previousUpdateChannel = make(chan bool)

	enteredService, _ := v.chatServiceEntry.GetText()
	v.chatServices.RemoveAll()

	csc, ec, endEarly := ac.session.GetChatServices(jid.ParseDomain(ac.Account()))
	go func() {
		hadAny := false
		defer func() {
			if hadAny && enteredService == "" {
				doInUIThread(func() {
					v.chatServices.SetActive(0)
				})
			}
			v.previousUpdateChannel = nil
		}()
		for {
			select {
			case <-v.previousUpdateChannel:
				doInUIThread(func() {
					v.chatServices.RemoveAll()
				})
				endEarly()
				return
			case err, _ := <-ec:
				if err != nil {
					ac.log.WithError(err).Error("something went wrong trying to get chat services")
				}
				return
			case cs, ok := <-csc:
				if !ok {
					return
				}

				hadAny = true
				doInUIThread(func() {
					v.chatServices.AppendText(cs.String())
				})
			}
		}
	}()
}

func (v *createMUCRoom) updateFields(f bool) {
	v.cancelButton.SetSensitive(f)
	v.createButton.SetSensitive(f)
	v.account.SetSensitive(f)
	v.room.SetSensitive(f)
	v.chatServices.SetSensitive(f)
}

func (v *createMUCRoom) getRoomID() jid.Bare {
	roomName, err := v.room.GetText()
	if err != nil {
		v.u.log.WithError(err).Error("something went wrong trying to create the room")
		doInUIThread(func() {
			v.errorBox.ShowMessage(i18n.Local("Could not get the room name, please try again."))
		})
		return nil
	}
	service := v.chatServices.GetActiveText()
	if !jid.ValidLocal(strings.TrimSpace(roomName)) || !jid.ValidDomain(strings.TrimSpace(roomName)) {
		doInUIThread(func() {
			v.errorBox.ShowMessage(i18n.Localf("%s@%s is not a valid room name.", roomName, service))
		})
		return nil
	}

	return jid.NewBare(jid.NewLocal(strings.TrimSpace(roomName)), jid.NewDomain(strings.TrimSpace(service)))
}

func reasonBasedOnError(err interface{}) string {
	switch e := err.(type) {
	case *session.ErrInvalidInformationQueryRequest:
		return i18n.Local("Could not send the information query to the server, please try again.")
	case *session.ErrUnexpectedResponse:
		return i18n.Local("The connection to the server can not be stablished.")
	case *session.ErrInformationQueryResponse:
		return e.Error()
	}
	return ""
}

func (v *createMUCRoom) createRoomHandler(ac *account) {
	if ac == nil {
		doInUIThread(func() {
			v.errorBox.ShowMessage(i18n.Local("No account selected, please select one account from the list or connect to one."))
		})
		return
	}

	roomIdentity := v.getRoomID()

	if roomIdentity != nil {
		doInUIThread(func() {
			v.updateFields(false)
			v.createButtonPrevText, _ = v.createButton.GetLabel()
			_ = v.createButton.SetProperty("label", i18n.Local("Creating room..."))
		})

		ec := ac.session.CreateRoom(roomIdentity)

		go func() {
			isRoomCreated := false
			defer func() {
				doInUIThread(func() {
					v.updateFields(true)
					_ = v.createButton.SetProperty("label", v.createButtonPrevText)
				})

				if isRoomCreated {
					doInUIThread(func() {
						v.u.mucShowRoom(ac, roomIdentity)
						v.dialog.Destroy()
					})
				}
			}()

			err, ok := <-ec
			if !ok {
				isRoomCreated = true
				return
			}

			if err != nil {
				ac.log.WithError(err.(error)).Error("something went wrong trying to create the room")
				var finalErr string
				if r := reasonBasedOnError(err); r != "" {
					finalErr = i18n.Localf("Could not create the new room, because the following reason:\n %s", r)
				} else {
					finalErr = i18n.Local("Could not create the new room")
				}
				doInUIThread(func() {
					v.errorBox.ShowMessage(finalErr)
				})
			}
		}()
	}
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
	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}

package gui

import (
	"fmt"
	"strings"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoom struct {
	errorBox *errorNotification
	builder  *builder

	// TODO[OB]-MUC: This should be assigned to a named field, as mentioned in earlier comments

	gtki.Dialog `gtk-widget:"create-chat-dialog"`

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
			// TODO[OB]-MUC: Signals are already executing in the UI thread, so not necessary to doInUIThread here
			doInUIThread(view.clearErrors)
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
	if v.previousUpdateChannel != nil {
		v.previousUpdateChannel <- true
	}

	v.previousUpdateChannel = make(chan bool)

	enteredService, _ := v.chatServiceEntry.GetText()
	v.chatServices.RemoveAll()

	// TODO[OB]-MUC: you should use jid.ParseDomain() here instead
	csc, ec, endEarly := ac.session.GetChatServices(jid.Parse(ac.Account()).Host())
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
		// TODO[OB]-MUC: It miht be good to log the error here
		doInUIThread(func() {
			v.errorBox.ShowMessage(i18n.Local("Could not get the room name, please try again."))
		})
		return nil
	}
	service := v.chatServices.GetActiveText()
	if strings.TrimSpace(roomName) == "" || strings.TrimSpace(service) == "" {
		// TODO[OB]-MUC: is this all validation that is necessary?
		doInUIThread(func() {
			v.errorBox.ShowMessage(i18n.Local("Please fill the required fields to create the room."))
		})
		return nil
	}

	// TODO[OB]-MUC: Use jid.NewBare instead, here
	ri, ok := jid.Parse(fmt.Sprintf("%s@%s", strings.TrimSpace(roomName), strings.TrimSpace(service))).(jid.Bare)
	if !ok {
		// TODO[OB]-MUC: This shouldn't really be possible. Can you give an example of how it can happen?
		doInUIThread(func() {
			v.errorBox.ShowMessage(i18n.Local("Room name not allowed."))
		})
		return nil
	}

	return ri
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

				if !isRoomCreated {
					doInUIThread(func() {
						v.errorBox.ShowMessage(i18n.Local("Could not create the new room."))
					})
					return
				}

				doInUIThread(func() {
					v.u.mucShowRoom(ac, roomIdentity)
					v.Destroy()
				})
			}()

			err, ok := <-ec
			if !ok {
				return
			}

			if err != nil {
				// TODO[OB]-MUC: Why is this at the Debug level?
				v.u.log.WithError(err).Debug("something went wrong trying to create the room")
				return
			}

			isRoomCreated = true
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
	view.SetTransientFor(u.window)
	view.Show()
}

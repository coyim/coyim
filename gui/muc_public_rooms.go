package gui

import (
	"errors"
	"sync"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	mucListRoomsIndexJid         = 0
	mucListRoomsIndexName        = 1
	mucListRoomsIndexService     = 2
	mucListRoomsIndexDescription = 3
	mucListRoomsIndexOccupants   = 4
)

type roomListingUpdateData struct {
	iter       gtki.TreeIter
	view       *mucPublicRoomsView
	generation int
}

func (u *gtkUI) updatedRoomListing(rl *muc.RoomListing, data interface{}) {
	d := data.(*roomListingUpdateData)

	// If we get an old update, we don't want to do anything at all
	if d.view.generation == d.generation {
		doInUIThread(func() {
			_ = d.view.roomsModel.SetValue(d.iter, mucListRoomsIndexDescription, rl.Description)
			_ = d.view.roomsModel.SetValue(d.iter, mucListRoomsIndexOccupants, rl.Occupants)
		})
	}
}

type mucPublicRoomsView struct {
	u              *gtkUI
	builder        *builder
	ac             *connectedAccountsComponent
	currentAccount *account

	generation    int
	updateLock    sync.RWMutex
	serviceGroups map[string]gtki.TreeIter
	cancel        chan bool

	dialog              gtki.Dialog         `gtk-widget:"publicRooms"`
	roomsModel          gtki.TreeStore      `gtk-widget:"roomsModel"`
	roomsTree           gtki.TreeView       `gtk-widget:"roomsTree"`
	rooms               gtki.ScrolledWindow `gtk-widget:"rooms"`
	spinner             gtki.Spinner        `gtk-widget:"spinner"`
	customService       gtki.Entry          `gtk-widget:"customServiceEntry"`
	notificationArea    gtki.Box            `gtk-widget:"notificationArea"`
	joinButton          gtki.Button         `gtk-widget:"buttonJoin"`
	refreshButton       gtki.Button         `gtk-widget:"buttonRefresh"`
	customServiceButton gtki.Button         `gtk-widget:"buttonCustomService"`

	notification gtki.InfoBar
	errorNotif   *errorNotification
}

func newMUCPublicRoomsView(u *gtkUI) *mucPublicRoomsView {
	view := &mucPublicRoomsView{u: u}
	view.init()
	return view
}

func (prv *mucPublicRoomsView) init() {
	prv.initBuilder()
	prv.initConnectedAccountsComponent()
	prv.initEvents()
	prv.initCommons()
}

func (prv *mucPublicRoomsView) initBuilder() {
	prv.builder = newBuilder("MUCPublicRoomsDialog")
	panicOnDevError(prv.builder.bindObjects(prv))
	prv.errorNotif = newErrorNotification(prv.notificationArea)
}

func (prv *mucPublicRoomsView) initConnectedAccountsComponent() {
	accountsInput := prv.builder.get("accounts").(gtki.ComboBox)
	ac := prv.u.createConnectedAccountsComponent(accountsInput, prv,
		func(acc *account) {
			// This is safe to do because we really have a selected account here
			prv.currentAccount = acc

			go prv.mucUpdatePublicRoomsOn(acc)
		},
		func() {
			prv.rooms.SetVisible(false)
			prv.spinner.Stop()
			prv.spinner.SetVisible(false)
			prv.roomsModel.Clear()
			prv.refreshButton.SetSensitive(false)
			prv.customServiceButton.SetSensitive(false)

			// We don't have a selected account anymore, we should
			// remove the existing reference to a no-selected account
			prv.currentAccount = nil
		},
	)
	prv.ac = ac
}

func (prv *mucPublicRoomsView) initEvents() {
	prv.builder.ConnectSignals(map[string]interface{}{
		"on_cancel":            prv.dialog.Destroy,
		"on_close_window":      prv.ac.onDestroy,
		"on_join":              prv.onJoinRoom,
		"on_activate_room_row": prv.onActivateRoomRow,
		"on_selection_changed": prv.onSelectionChanged,
		"on_custom_service":    prv.onUpdatePublicRooms,
		"on_refresh":           prv.onUpdatePublicRooms,
	})
}

func (prv *mucPublicRoomsView) initCommons() {
	prv.serviceGroups = make(map[string]gtki.TreeIter)
	prv.joinButton.SetSensitive(false)
}

var (
	errNoPossibleSelection = errors.New("problem getting selection")
	errNoSelection         = errors.New("nothing is selected")
	errNoRoomSelected      = errors.New("a service is selected, not a room, so we can't activate it")
	errNoService           = errors.New("no service is available")
)

func (prv *mucPublicRoomsView) getRoomBareFromIter(iter gtki.TreeIter) (jid.Bare, error) {
	roomJidValue, _ := prv.roomsModel.GetValue(iter, mucListRoomsIndexJid)
	ident, _ := roomJidValue.GetString()

	_, ok := prv.serviceGroups[ident]
	if ok {
		return nil, errNoRoomSelected
	}

	serviceValue, _ := prv.roomsModel.GetValue(iter, mucListRoomsIndexService)
	service, _ := serviceValue.GetString()
	_, ok = prv.serviceGroups[service]
	if !ok {
		return nil, errNoService
	}

	return jid.NewBare(jid.NewLocal(ident), jid.NewDomain(service)), nil
}

func (prv *mucPublicRoomsView) getSelectedRoomBare() (jid.Bare, error) {
	selection, err := prv.roomsTree.GetSelection()
	if err != nil {
		return nil, errNoPossibleSelection
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return nil, errNoSelection
	}

	return prv.getRoomBareFromIter(iter)
}

func (prv *mucPublicRoomsView) onJoinRoom() {
	ident, err := prv.getSelectedRoomBare()
	if err != nil {
		prv.currentAccount.log.WithError(err).Error("An error occurred when trying to join the room")
		prv.showUserMessageForError(err)
		return
	}

	go prv.joinRoom(ident)
}

func (prv *mucPublicRoomsView) onActivateRoomRow(_ gtki.TreeView, path gtki.TreePath) {
	iter, err := prv.roomsModel.GetIter(path)
	if err != nil {
		prv.currentAccount.log.WithError(err).Error("Couldn't activate the selected item")
		return
	}

	ident, err := prv.getRoomBareFromIter(iter)
	if err != nil {
		prv.currentAccount.log.WithError(err).Error("Couldn't join to the room based on the current selection")
		prv.showUserMessageForError(err)
		return
	}

	go prv.joinRoom(ident)
}

func (prv *mucPublicRoomsView) onSelectionChanged() {
	_, err := prv.getSelectedRoomBare()
	if err != nil {
		prv.joinButton.SetSensitive(false)
	} else {
		prv.joinButton.SetSensitive(true)
	}
}

func (prv *mucPublicRoomsView) onUpdatePublicRooms() {
	go prv.mucUpdatePublicRoomsOn(prv.ac.currentAccount())
}

// mucUpdatePublicRoomsOn should NOT be called from the UI thread
func (prv *mucPublicRoomsView) mucUpdatePublicRoomsOn(a *account) {
	if prv.cancel != nil {
		prv.cancel <- true
	}

	prv.updateLock.Lock()

	doInUIThread(prv.clearErrors)

	customService, _ := prv.customService.GetText()

	prv.cancel = make(chan bool, 1)

	doInUIThread(func() {
		prv.rooms.SetVisible(false)
		prv.spinner.Start()
		prv.spinner.SetVisible(true)
		prv.roomsModel.Clear()
		prv.refreshButton.SetSensitive(true)
		prv.customServiceButton.SetSensitive(true)
	})
	prv.generation++
	prv.serviceGroups = make(map[string]gtki.TreeIter)

	// We save the generation value here, in case it gets modified inside the view later
	gen := prv.generation

	res, resServices, ec := a.session.GetRooms(jid.Parse(a.Account()).Host(), customService)
	go func() {
		hasSomething := false

		defer func() {
			if !hasSomething {
				doInUIThread(func() {
					prv.spinner.Stop()
					prv.spinner.SetVisible(false)
					prv.rooms.SetVisible(true)
					if customService != "" {
						prv.notifyOnError(i18n.Local("That service doesn't seem to exist"))
					} else {
						prv.notifyOnError(i18n.Local("Your XMPP server doesn't seem to have any chat room services"))
					}
				})
			}

			prv.updateLock.Unlock()
		}()
		for {
			select {
			case sl, ok := <-resServices:
				if !ok {
					return
				}
				if !hasSomething {
					hasSomething = true
					doInUIThread(func() {
						prv.spinner.Stop()
						prv.spinner.SetVisible(false)
						prv.rooms.SetVisible(true)
					})
				}

				serv, ok := prv.serviceGroups[sl.Jid.String()]
				if !ok {
					doInUIThread(func() {
						serv = prv.roomsModel.Append(nil)
						prv.serviceGroups[sl.Jid.String()] = serv
						_ = prv.roomsModel.SetValue(serv, mucListRoomsIndexJid, sl.Jid.String())
						_ = prv.roomsModel.SetValue(serv, mucListRoomsIndexName, sl.Name)
					})
				}
			case rl, ok := <-res:
				if !ok || rl == nil {
					return
				}

				if !hasSomething {
					hasSomething = true
					doInUIThread(func() {
						prv.spinner.Stop()
						prv.spinner.SetVisible(false)
						prv.rooms.SetVisible(true)
					})
				}

				serv, ok := prv.serviceGroups[rl.Service.String()]
				doInUIThread(func() {
					if !ok {
						serv = prv.roomsModel.Append(nil)
						prv.serviceGroups[rl.Service.String()] = serv
						_ = prv.roomsModel.SetValue(serv, mucListRoomsIndexJid, rl.Service.String())
						_ = prv.roomsModel.SetValue(serv, mucListRoomsIndexName, rl.ServiceName)
					}

					iter := prv.roomsModel.Append(serv)
					_ = prv.roomsModel.SetValue(iter, mucListRoomsIndexJid, rl.Jid.Local().String())
					_ = prv.roomsModel.SetValue(iter, mucListRoomsIndexName, rl.Name)
					_ = prv.roomsModel.SetValue(iter, mucListRoomsIndexService, rl.Service.String())
					rl.OnUpdate(prv.u.updatedRoomListing, &roomListingUpdateData{iter, prv, gen})

					prv.roomsTree.ExpandAll()
				})
			case e, ok := <-ec:
				if !ok {
					return
				}
				if e != nil {
					doInUIThread(func() {
						prv.notifyOnError(i18n.Local("Something went wrong when trying to get chat rooms"))
					})
					prv.currentAccount.log.WithError(e).Error("Something went wrong when trying to get chat rooms")
				}
				return
			case _, _ = <-prv.cancel:
				return
			}
		}
	}()
}

func (prv *mucPublicRoomsView) clearErrors() {
	prv.errorNotif.Hide()
}

func (prv *mucPublicRoomsView) notifyOnError(err string) {
	if prv.notification != nil {
		prv.notificationArea.Remove(prv.notification)
	}

	prv.errorNotif.ShowMessage(err)
}

func (prv *mucPublicRoomsView) showUserMessageForError(err error) {
	var userMessage string

	switch err {
	case errNoPossibleSelection:
		userMessage = i18n.Local("We cant't determinate what has been selected, please try again.")
	case errNoRoomSelected:
		userMessage = i18n.Local("The selected item is not a room, select one room from the list to join to.")
	case errNoSelection:
		userMessage = i18n.Local("Please, select one room from the list to join to.")
	case errNoService:
		userMessage = i18n.Local("We cant't determinate which service has been selected, please try again.")
	default:
		userMessage = i18n.Local("An unknow error ocurred, please try again.")
	}

	prv.notifyOnError(userMessage)
}

// mucShowPublicRooms should be called from the UI thread
func (u *gtkUI) mucShowPublicRooms() {
	view := newMUCPublicRoomsView(u)

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}

// joinRoom should not be called from the UI thread
func (prv *mucPublicRoomsView) joinRoom(roomJid jid.Bare) {
	a := prv.ac.currentAccount()
	if a == nil {
		prv.currentAccount.log.WithField("room", roomJid).Debug("joinRoom(): no account is selected")
		prv.notifyOnError(i18n.Local("No account was selected, please select one account from the list."))
		return
	}

	a.log.WithField("room", roomJid).Debug("joinRoom()")
	doInUIThread(func() {
		prv.dialog.Destroy()
		prv.u.mucShowRoom(a, roomJid)
	})
}

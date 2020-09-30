package gui

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	mucListRoomsIndexJid = iota
	mucListRoomsIndexName
	mucListRoomsIndexService
	mucListRoomsIndexDescription
	mucListRoomsIndexOccupants
	mucListRoomsIndexRoomInfo
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
	roomInfos     map[int]*muc.RoomListing
	cancel        chan bool

	dialog           gtki.Dialog         `gtk-widget:"publicRooms"`
	roomsModel       gtki.TreeStore      `gtk-widget:"roomsModel"`
	roomsTree        gtki.TreeView       `gtk-widget:"roomsTree"`
	rooms            gtki.ScrolledWindow `gtk-widget:"rooms"`
	spinner          gtki.Spinner        `gtk-widget:"spinner"`
	customService    gtki.Entry          `gtk-widget:"customServiceEntry"`
	notificationArea gtki.Box            `gtk-widget:"notificationArea"`

	joinButton          gtki.Button `gtk-widget:"buttonJoin"`
	refreshButton       gtki.Button `gtk-widget:"buttonRefresh"`
	customServiceButton gtki.Button `gtk-widget:"buttonCustomService"`

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
	prv.roomInfos = make(map[int]*muc.RoomListing)
	prv.joinButton.SetSensitive(false)
}

func (prv *mucPublicRoomsView) log() coylog.Logger {
	l := prv.u.log
	if prv.currentAccount != nil {
		l = prv.currentAccount.log
	}

	l.WithField("who", "mucPublilcRoomsView")

	return l
}

var (
	errNoPossibleSelection = errors.New("problem getting selection")
	errNoSelection         = errors.New("nothing is selected")
	errNoRoomSelected      = errors.New("a service is selected, not a room, so we can't activate it")
	errNoService           = errors.New("no service is available")
)

func (prv *mucPublicRoomsView) getRoomListingFromIter(iter gtki.TreeIter) (*muc.RoomListing, error) {
	roomInfoRealVal, err := prv.getRoomInfoFromIter(iter)
	if err != nil {
		return nil, err
	}

	rl, ok := prv.roomInfos[roomInfoRealVal]
	if !ok || rl == nil {
		return nil, errNoPossibleSelection
	}
	return rl, nil
}

func (prv *mucPublicRoomsView) getRoomInfoFromIter(iter gtki.TreeIter) (int, error) {
	roomInfoValue, e1 := prv.roomsModel.GetValue(iter, mucListRoomsIndexRoomInfo)
	roomInfoRef, e2 := roomInfoValue.GoValue()
	if e1 != nil || e2 != nil {
		return 0, errNoRoomSelected
	}

	roomInfoRealVal := roomInfoRef.(int)
	return roomInfoRealVal, nil
}

func (prv *mucPublicRoomsView) getRoomIDFromIter(iter gtki.TreeIter) (jid.Bare, error) {
	roomJidValue, _ := prv.roomsModel.GetValue(iter, mucListRoomsIndexJid)
	roomName, _ := roomJidValue.GetString()

	_, ok := prv.serviceGroups[roomName]
	if ok {
		return nil, errNoRoomSelected
	}

	serviceValue, _ := prv.roomsModel.GetValue(iter, mucListRoomsIndexService)
	service, _ := serviceValue.GetString()
	_, ok = prv.serviceGroups[service]
	if !ok {
		return nil, errNoService
	}

	return jid.NewBare(jid.NewLocal(roomName), jid.NewDomain(service)), nil
}

func (prv *mucPublicRoomsView) getRoomFromIter(iter gtki.TreeIter) (jid.Bare, *muc.RoomListing, error) {
	rl, err := prv.getRoomListingFromIter(iter)
	if err != nil {
		return nil, nil, err
	}

	roomID, err := prv.getRoomIDFromIter(iter)
	if err != nil {
		return nil, nil, err
	}

	return roomID, rl, nil
}

func (prv *mucPublicRoomsView) getSelectedRoom() (jid.Bare, *muc.RoomListing, error) {
	selection, err := prv.roomsTree.GetSelection()
	if err != nil {
		return nil, nil, errNoPossibleSelection
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		return nil, nil, errNoSelection
	}

	return prv.getRoomFromIter(iter)
}

func (prv *mucPublicRoomsView) onJoinRoom() {
	ident, rl, err := prv.getSelectedRoom()
	if err != nil {
		prv.log().WithError(err).Error("An error occurred when trying to join the room")
		prv.showUserMessageForError(err)
		return
	}

	go prv.joinRoom(ident, rl)
}

func (prv *mucPublicRoomsView) onActivateRoomRow(_ gtki.TreeView, path gtki.TreePath) {
	iter, err := prv.roomsModel.GetIter(path)
	if err != nil {
		prv.log().WithError(err).Error("Couldn't activate the selected item")
		return
	}

	ident, rl, err := prv.getRoomFromIter(iter)
	if err != nil {
		prv.log().WithError(err).Error("Couldn't join to the room based on the current selection")
		prv.showUserMessageForError(err)
		return
	}

	go prv.joinRoom(ident, rl)
}

func (prv *mucPublicRoomsView) onSelectionChanged() {
	_, _, err := prv.getSelectedRoom()
	prv.joinButton.SetSensitive(err == nil)
}

func (prv *mucPublicRoomsView) onUpdatePublicRooms() {
	go prv.mucUpdatePublicRoomsOn(prv.currentAccount)
}

func (prv *mucPublicRoomsView) getFreshRoomInfoIdentifierAndSet(rl *muc.RoomListing) int {
	for {
		// We are getting a 31 bit integer here to avoid negative numbers
		// We also want to have it be smaller than normal size Golang numbers because
		// this number will be sent out into C, with Glib. For some reason, it did not
		// like negative numbers, and it doesn't work well will 64 bit numbers either
		v := int(rand.Int31())
		_, ok := prv.roomInfos[v]
		if !ok {
			prv.roomInfos[v] = rl
			return v
		}
	}
}

// TODO: This method is a bit of a beast. We should probably refactor and clean up a bit
// This is big enough that using a helper context object might be necessary.

// mucUpdatePublicRoomsOn MUST NOT be called from the UI thread
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
	prv.roomInfos = make(map[int]*muc.RoomListing)

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

					// This will block while finding an unused identifier. However, since
					// we don't expect to get millions of room listings, it's not likely this will ever be a problem.
					roomInfoRef := prv.getFreshRoomInfoIdentifierAndSet(rl)
					_ = prv.roomsModel.SetValue(iter, mucListRoomsIndexRoomInfo, roomInfoRef)

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
					prv.log().WithError(e).Error("Something went wrong when trying to get chat rooms")
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

// mucShowPublicRooms MUST be called from the UI thread
func (u *gtkUI) mucShowPublicRooms() {
	view := newMUCPublicRoomsView(u)

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}

// joinRoom should not be called from the UI thread
func (prv *mucPublicRoomsView) joinRoom(roomJid jid.Bare, roomInfo *muc.RoomListing) {
	// TODO: should we use the current account field here?
	if prv.currentAccount == nil {
		prv.log().WithField("room", roomJid).Debug("joinRoom(): no account is selected")
		prv.notifyOnError(i18n.Local("No account was selected, please select one account from the list."))
		return
	}

	prv.log().WithField("room", roomJid).Debug("joinRoom()")
	doInUIThread(func() {
		prv.dialog.Destroy()
		prv.u.joinRoom(prv.currentAccount, roomJid, nil)
	})
}

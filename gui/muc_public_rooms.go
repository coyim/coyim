package gui

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glibi"
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

	roomsModel gtki.TreeStore

	dialog              gtki.Dialog         `gtk-widget:"public-rooms-dialog"`
	roomsTree           gtki.TreeView       `gtk-widget:"public-rooms-tree"`
	rooms               gtki.ScrolledWindow `gtk-widget:"public-rooms-view"`
	customService       gtki.Entry          `gtk-widget:"custom-service-entry"`
	joinButton          gtki.Button         `gtk-widget:"join-room-button"`
	refreshButton       gtki.Button         `gtk-widget:"refresh-button"`
	customServiceButton gtki.Button         `gtk-widget:"list-rooms-button"`

	notificationArea gtki.Box `gtk-widget:"notifications-area"`
	notifications    *notifications

	spinnerOverlay gtki.Overlay `gtk-widget:"spinner-overlay"`
	spinnerBox     gtki.Box     `gtk-widget:"spinner-box"`
	spinner        *spinner
}

func newMUCPublicRoomsView(u *gtkUI) *mucPublicRoomsView {
	view := &mucPublicRoomsView{u: u}

	view.initBuilder()
	view.initModel()
	view.initNotificationsAndSpinner(u)
	view.initConnectedAccountsComponent()
	view.initDefaults()

	return view
}

func (prv *mucPublicRoomsView) initBuilder() {
	prv.builder = newBuilder("MUCPublicRoomsDialog")
	panicOnDevError(prv.builder.bindObjects(prv))

	prv.builder.ConnectSignals(map[string]interface{}{
		"on_cancel":            prv.onCancel,
		"on_window_closed":     prv.onWindowClose,
		"on_join":              prv.onJoinRoom,
		"on_activate_room_row": prv.onActivateRoomRow,
		"on_selection_changed": prv.onSelectionChanged,
		"on_custom_service":    prv.onUpdatePublicRooms,
		"on_refresh":           prv.onUpdatePublicRooms,
	})
}

func (prv *mucPublicRoomsView) initModel() {
	roomsModel, _ := g.gtk.TreeStoreNew(
		// jid
		glibi.TYPE_STRING,
		// name
		glibi.TYPE_STRING,
		// service
		glibi.TYPE_STRING,
		// description
		glibi.TYPE_STRING,
		// occupants
		glibi.TYPE_INT,
		// room info reference
		glibi.TYPE_INT,
	)

	prv.roomsModel = roomsModel
	prv.roomsTree.SetModel(prv.roomsModel)
}

func (prv *mucPublicRoomsView) initNotificationsAndSpinner(u *gtkUI) {
	prv.notifications = u.newNotifications(prv.notificationArea)

	prv.spinner = newSpinner()
	s := prv.spinner.getWidget()

	// This is a GTK trick to set the size of the spinner,
	// so if the parent has a size of 40x40 for example, with
	// the bellow properties the spinner will be 40x40 too
	s.SetProperty("hexpand", true)
	s.SetProperty("vexpand", true)

	prv.spinnerBox.Add(s)
}

func (prv *mucPublicRoomsView) initConnectedAccountsComponent() {
	accountsInput := prv.builder.get("accounts").(gtki.ComboBox)
	ac := prv.u.createConnectedAccountsComponent(accountsInput, prv.notifications, prv.onAccountsUpdated, prv.onNoAccounts)
	prv.ac = ac
}

func (prv *mucPublicRoomsView) initDefaults() {
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

func (prv *mucPublicRoomsView) onAccountsUpdated(ca *account) {
	prv.currentAccount = ca
	prv.updatePublicRoomsForCurrentAccount()
}

func (prv *mucPublicRoomsView) onNoAccounts() {
	prv.currentAccount = nil

	prv.disableRoomsView()
	prv.hideSpinner()

	prv.roomsModel.Clear()
	prv.refreshButton.SetSensitive(false)
	prv.customServiceButton.SetSensitive(false)
}

func (prv *mucPublicRoomsView) onCancel() {
	prv.dialog.Destroy()
}

func (prv *mucPublicRoomsView) onWindowClose() {
	prv.cancelActiveUpdate()
	prv.ac.onDestroy()
}

func (prv *mucPublicRoomsView) cancelActiveUpdate() {
	if prv.cancel == nil {
		return
	}

	prv.cancel <- true
	prv.cancel = nil
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
	prv.updatePublicRoomsForCurrentAccount()
}

func (prv *mucPublicRoomsView) updatePublicRoomsForCurrentAccount() {
	if prv.currentAccount != nil {
		go prv.mucUpdatePublicRoomsOn(prv.currentAccount)
	}
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

func (prv *mucPublicRoomsView) beforeUpdatingPublicRooms() {
	prv.notifications.clearErrors()
	prv.disableRoomsViewAndShowSpinner()

	prv.roomsModel.Clear()
	prv.refreshButton.SetSensitive(true)
	prv.customServiceButton.SetSensitive(true)
}

func (prv *mucPublicRoomsView) onUpdatePublicRoomsNoResults(customService string) {
	prv.enableRoomsViewAndHideSpinner()
	if customService != "" {
		prv.notifications.error(i18n.Local("That service doesn't seem to exist"))
	} else {
		prv.notifications.error(i18n.Local("Your XMPP server doesn't seem to have any chat room services"))
	}
}

func (prv *mucPublicRoomsView) showSpinner() {
	prv.spinnerOverlay.Show()
	prv.spinner.show()
}

func (prv *mucPublicRoomsView) hideSpinner() {
	prv.spinner.hide()
	prv.spinnerOverlay.Hide()
}

func (prv *mucPublicRoomsView) enableRoomsViewAndHideSpinner() {
	prv.rooms.SetSensitive(true)
	prv.hideSpinner()
}

func (prv *mucPublicRoomsView) disableRoomsViewAndShowSpinner() {
	prv.disableRoomsView()
	prv.showSpinner()
}

func (prv *mucPublicRoomsView) disableRoomsView() {
	prv.rooms.SetSensitive(false)
}

func (prv *mucPublicRoomsView) addNewServiceToModel(roomName, serviceName string) gtki.TreeIter {
	serv := prv.roomsModel.Append(nil)

	prv.serviceGroups[roomName] = serv
	_ = prv.roomsModel.SetValue(serv, mucListRoomsIndexJid, roomName)
	_ = prv.roomsModel.SetValue(serv, mucListRoomsIndexName, serviceName)

	return serv
}

func (prv *mucPublicRoomsView) addNewRoomToModel(parentIter gtki.TreeIter, rl *muc.RoomListing, gen int) {
	iter := prv.roomsModel.Append(parentIter)

	_ = prv.roomsModel.SetValue(iter, mucListRoomsIndexJid, rl.Jid.Local().String())
	_ = prv.roomsModel.SetValue(iter, mucListRoomsIndexName, rl.Name)
	_ = prv.roomsModel.SetValue(iter, mucListRoomsIndexService, rl.Service.String())

	// This will block while finding an unused identifier. However, since
	// we don't expect to get millions of room listings, it's not likely this will ever be a problem.
	roomInfoRef := prv.getFreshRoomInfoIdentifierAndSet(rl)
	_ = prv.roomsModel.SetValue(iter, mucListRoomsIndexRoomInfo, roomInfoRef)

	rl.OnUpdate(prv.u.updatedRoomListing, &roomListingUpdateData{iter, prv, gen})

	prv.roomsTree.ExpandAll()
}

func (prv *mucPublicRoomsView) handleReceivedServiceListing(sl *muc.ServiceListing) {
	_, ok := prv.serviceGroups[sl.Jid.String()]
	if !ok {
		doInUIThread(func() {
			prv.addNewServiceToModel(sl.Jid.String(), sl.Name)
		})
	}
}

func (prv *mucPublicRoomsView) handleReceivedRoomListing(rl *muc.RoomListing, gen int) {
	serv, ok := prv.serviceGroups[rl.Service.String()]
	doInUIThread(func() {
		if !ok {
			serv = prv.addNewServiceToModel(rl.Service.String(), rl.ServiceName)
		}

		prv.addNewRoomToModel(serv, rl, gen)
	})
}

func (prv *mucPublicRoomsView) handleReceivedError(err error) {
	prv.log().WithError(err).Error("Something went wrong when trying to get chat rooms")
	doInUIThread(func() {
		prv.notifications.error(i18n.Local("Something went wrong when trying to get chat rooms"))
	})
}

func (prv *mucPublicRoomsView) listenPublicRoomsResponse(gen int, res <-chan *muc.RoomListing, resServices <-chan *muc.ServiceListing, ec <-chan error) bool {
	select {
	case sl, ok := <-resServices:
		if !ok {
			return false
		}

		prv.handleReceivedServiceListing(sl)
	case rl, ok := <-res:
		if !ok || rl == nil {
			return false
		}

		prv.handleReceivedRoomListing(rl, gen)
	case err, ok := <-ec:
		if !ok {
			return false
		}
		if err != nil {
			prv.handleReceivedError(err)
		}
		return false
	case <-prv.cancel:
		return false
	}
	return true
}

func (prv *mucPublicRoomsView) listenPublicRoomsUpdate(customService string, gen int, res <-chan *muc.RoomListing, resServices <-chan *muc.ServiceListing, ec <-chan error) {
	hasSomething := false

	defer func() {
		if !hasSomething {
			doInUIThread(func() {
				prv.onUpdatePublicRoomsNoResults(customService)
			})
		}

		prv.updateLock.Unlock()
	}()

	for prv.listenPublicRoomsResponse(gen, res, resServices, ec) {
		if !hasSomething {
			hasSomething = true
			doInUIThread(prv.enableRoomsViewAndHideSpinner)
		}
	}
}

// mucUpdatePublicRoomsOn MUST NOT be called from the UI thread
func (prv *mucPublicRoomsView) mucUpdatePublicRoomsOn(a *account) {
	prv.cancelActiveUpdate()

	prv.updateLock.Lock()
	prv.cancel = make(chan bool, 1)

	doInUIThread(prv.beforeUpdatingPublicRooms)

	prv.generation++
	prv.serviceGroups = make(map[string]gtki.TreeIter)
	prv.roomInfos = make(map[int]*muc.RoomListing)

	// We save the generation value here, in case it gets modified inside the view later
	gen := prv.generation

	customService, _ := prv.customService.GetText()

	res, resServices, ec := a.session.GetRooms(jid.Parse(a.Account()).Host(), customService)
	go prv.listenPublicRoomsUpdate(customService, gen, res, resServices, ec)
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

	prv.notifications.error(userMessage)
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
	if prv.currentAccount == nil {
		prv.log().WithField("room", roomJid).Debug("joinRoom(): no account is selected")
		prv.notifications.error(i18n.Local("No account was selected, please select one account from the list."))
		return
	}

	prv.log().WithField("room", roomJid).Debug("joinRoom()")
	doInUIThread(func() {
		prv.dialog.Destroy()
		prv.u.joinRoom(prv.currentAccount, roomJid, nil)
	})
}

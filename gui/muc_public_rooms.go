package gui

import (
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
	u       *gtkUI
	builder *builder
	ac      *connectedAccountsComponent

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
			go prv.mucUpdatePublicRoomsOn(acc)
		},
		func() {
			prv.rooms.SetVisible(false)
			prv.spinner.Stop()
			prv.spinner.SetVisible(false)
			prv.roomsModel.Clear()
			prv.refreshButton.SetSensitive(false)
			prv.customServiceButton.SetSensitive(false)
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

func (prv *mucPublicRoomsView) onJoinRoom() {
	selection, err := prv.roomsTree.GetSelection()
	if err != nil {
		prv.u.log.WithError(err).Debug("couldn't join")
		prv.notifyOnError(i18n.Local("Please, select one room from the list to join to."))
		return
	}

	_, iter, ok := selection.GetSelected()
	if !ok {
		prv.u.log.Debug("nothing is selected")
		prv.notifyOnError(i18n.Local("Please, select one room from the list to join to."))
		return
	}

	val, _ := prv.roomsModel.GetValue(iter, mucListRoomsIndexJid)
	v, _ := val.GetString()

	_, ok = prv.serviceGroups[v]
	if ok {
		prv.u.log.Debug("a service is selected, not a room, so we can't activate it")
		prv.notifyOnError(i18n.Local("The selected item is not a room, select one room from the list to join to."))
		return
	}

	go prv.joinRoom(v)
}

func (prv *mucPublicRoomsView) onActivateRoomRow(_ gtki.TreeView, path gtki.TreePath) {
	iter, err := prv.roomsModel.GetIter(path)
	if err != nil {
		prv.u.log.WithError(err).Debug("couldn't activate")
		return
	}
	val, _ := prv.roomsModel.GetValue(iter, mucListRoomsIndexJid)
	v, _ := val.GetString()
	_, ok := prv.serviceGroups[v]
	if ok {
		prv.u.log.Debug("a service is selected, not a room, so we can't activate it")
		prv.notifyOnError(i18n.Local("The selected item is not a room, select one room from the list to join to."))
		return
	}
	go prv.joinRoom(v)
}

func (prv *mucPublicRoomsView) onSelectionChanged() {
	selection, err := prv.roomsTree.GetSelection()
	if err != nil {
		prv.u.log.WithError(err).Debug("problem getting selection")
		return
	}
	_, iter, ok := selection.GetSelected()
	if !ok {
		prv.joinButton.SetSensitive(false)
		return
	}

	val, _ := prv.roomsModel.GetValue(iter, mucListRoomsIndexJid)
	v, _ := val.GetString()
	_, ok = prv.serviceGroups[v]
	if ok {
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
					prv.u.log.WithError(e).Debug("something went wrong trying to get chat rooms")
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

// mucShowPublicRooms should be called from the UI thread
func (u *gtkUI) mucShowPublicRooms() {
	view := newMUCPublicRoomsView(u)

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}

// joinRoom should not be called from the UI thread
func (prv *mucPublicRoomsView) joinRoom(roomJid string) {
	a := prv.ac.currentAccount()
	if a == nil {
		prv.u.log.WithField("room", roomJid).Debug("joinRoom(): no account is selected")
		prv.notifyOnError(i18n.Local("No account was selected, please select one account from the list."))
		return
	}

	a.log.WithField("room", roomJid).Debug("joinRoom()")
	prv.u.mucShowRoom(a, jid.Parse(roomJid).(jid.Bare))
}

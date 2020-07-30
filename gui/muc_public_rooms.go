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
	builder *builder

	generation    int
	updateLock    sync.RWMutex
	serviceGroups map[string]gtki.TreeIter
	cancel        chan bool

	dialog           gtki.Dialog         `gtk-widget:"public-rooms"`
	roomsModel       gtki.TreeStore      `gtk-widget:"rooms-model"`
	roomsTree        gtki.TreeView       `gtk-widget:"rooms-tree"`
	rooms            gtki.ScrolledWindow `gtk-widget:"rooms"`
	spinner          gtki.Spinner        `gtk-widget:"spinner"`
	customService    gtki.Entry          `gtk-widget:"customServiceEntry"`
	notificationArea gtki.Box            `gtk-widget:"notification-area"`
	joinButton       gtki.Button         `gtk-widget:"button_join"`
	notification     gtki.InfoBar
	errorNotif       *errorNotification
}

func (prv *mucPublicRoomsView) clearErrors() {
	prv.errorNotif.Hide()
}

func (prv *mucPublicRoomsView) notifyOnError(err string) {
	doInUIThread(func() {
		if prv.notification != nil {
			prv.notificationArea.Remove(prv.notification)
		}

		prv.errorNotif.ShowMessage(err)
	})
}

func (prv *mucPublicRoomsView) init() {
	prv.builder = newBuilder("MUCPublicRoomsDialog")
	panicOnDevError(prv.builder.bindObjects(prv))
	prv.serviceGroups = make(map[string]gtki.TreeIter)
	prv.errorNotif = newErrorNotification(prv.notificationArea)
}

// mucUpdatePublicRoomsOn should NOT be called from the UI thread
func (u *gtkUI) mucUpdatePublicRoomsOn(view *mucPublicRoomsView, a *account) {
	if view.cancel != nil {
		view.cancel <- true
	}

	view.updateLock.Lock()

	doInUIThread(view.clearErrors)

	customService, _ := view.customService.GetText()

	view.cancel = make(chan bool, 1)

	doInUIThread(func() {
		view.rooms.SetVisible(false)
		view.spinner.Start()
		view.spinner.SetVisible(true)
		view.roomsModel.Clear()
	})
	view.generation++
	view.serviceGroups = make(map[string]gtki.TreeIter)

	// We save the generation value here, in case it gets modified inside the view later
	gen := view.generation

	res, resServices, ec := a.session.GetRooms(jid.Parse(a.Account()).Host(), customService)
	go func() {
		hasSomething := false

		defer func() {
			if !hasSomething {
				doInUIThread(func() {
					view.spinner.Stop()
					view.spinner.SetVisible(false)
					view.rooms.SetVisible(true)
					if customService != "" {
						view.notifyOnError(i18n.Local("That service doesn't seem to exist"))
					} else {
						view.notifyOnError(i18n.Local("Your XMPP server doesn't seem to have any chat room services"))
					}
				})
			}

			view.updateLock.Unlock()
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
						view.spinner.Stop()
						view.spinner.SetVisible(false)
						view.rooms.SetVisible(true)
					})
				}

				serv, ok := view.serviceGroups[sl.Jid.String()]
				if !ok {
					doInUIThread(func() {
						serv = view.roomsModel.Append(nil)
						view.serviceGroups[sl.Jid.String()] = serv
						_ = view.roomsModel.SetValue(serv, mucListRoomsIndexJid, sl.Jid.String())
						_ = view.roomsModel.SetValue(serv, mucListRoomsIndexName, sl.Name)
					})
				}
			case rl, ok := <-res:
				if !ok || rl == nil {
					return
				}

				if !hasSomething {
					hasSomething = true
					doInUIThread(func() {
						view.spinner.Stop()
						view.spinner.SetVisible(false)
						view.rooms.SetVisible(true)
					})
				}

				serv, ok := view.serviceGroups[rl.Service.String()]
				doInUIThread(func() {
					if !ok {
						serv = view.roomsModel.Append(nil)
						view.serviceGroups[rl.Service.String()] = serv
						_ = view.roomsModel.SetValue(serv, mucListRoomsIndexJid, rl.Service.String())
						_ = view.roomsModel.SetValue(serv, mucListRoomsIndexName, rl.ServiceName)
					}

					iter := view.roomsModel.Append(serv)
					_ = view.roomsModel.SetValue(iter, mucListRoomsIndexJid, string(rl.Jid.Local()))
					_ = view.roomsModel.SetValue(iter, mucListRoomsIndexName, rl.Name)
					_ = view.roomsModel.SetValue(iter, mucListRoomsIndexService, rl.Service.String())
					rl.OnUpdate(u.updatedRoomListing, &roomListingUpdateData{iter, view, gen})

					view.roomsTree.ExpandAll()
				})
			case e, ok := <-ec:
				if !ok {
					return
				}
				if e != nil {
					view.notifyOnError(i18n.Local("Something went wrong when trying to get chat rooms"))
					u.log.WithError(e).Debug("something went wrong trying to get chat rooms")
				}
				return
			case _, _ = <-view.cancel:
				return
			}
		}
	}()
}

// joinRoom should not be called from the UI thread
func (u *gtkUI) joinRoom(roomJid string, a *account) {
	a.log.WithField("room", roomJid).Debug("joinRoom()")
	// TODO: implement
}

// mucShowPublicRooms should be called from the UI thread
func (u *gtkUI) mucShowPublicRooms() {
	view := &mucPublicRoomsView{}
	view.init()

	accountsInput := view.builder.get("accounts").(gtki.ComboBox)
	ac := u.createConnectedAccountsComponent(accountsInput, view,
		func(acc *account) {
			go u.mucUpdatePublicRoomsOn(view, acc)
		},
		func() {
			view.rooms.SetVisible(false)
			view.spinner.Stop()
			view.spinner.SetVisible(false)
			view.roomsModel.Clear()
		},
	)

	view.joinButton.SetSensitive(false)

	view.builder.ConnectSignals(map[string]interface{}{
		"on_cancel_signal":       view.dialog.Destroy,
		"on_close_window_signal": ac.onDestroy,
		"on_join_signal": func() {
			selection, err := view.roomsTree.GetSelection()
			if err != nil {
				u.log.WithError(err).Debug("couldn't join")
				return
			}
			_, iter, ok := selection.GetSelected()
			if !ok {
				u.log.Debug("nothing is selected")
				return
			}

			val, _ := view.roomsModel.GetValue(iter, mucListRoomsIndexJid)
			v, _ := val.GetString()
			_, ok = view.serviceGroups[v]

			if ok {
				u.log.Debug("a service is selected, not a room, so we can't activate it")
				return
			}
			go u.joinRoom(v, ac.currentAccount())
		},
		"on_activate_room_row": func(_ gtki.TreeView, path gtki.TreePath) {
			iter, err := view.roomsModel.GetIter(path)
			if err != nil {
				u.log.WithError(err).Debug("couldn't activate")
				return
			}
			val, _ := view.roomsModel.GetValue(iter, mucListRoomsIndexJid)
			v, _ := val.GetString()
			_, ok := view.serviceGroups[v]
			if ok {
				u.log.Debug("a service is selected, not a room, so we can't activate it")
				return
			}
			go u.joinRoom(v, ac.currentAccount())
		},
		"on_selection_changed": func() {
			selection, err := view.roomsTree.GetSelection()
			if err != nil {
				u.log.WithError(err).Debug("problem getting selection")
				return
			}
			_, iter, ok := selection.GetSelected()
			if !ok {
				view.joinButton.SetSensitive(false)
				return
			}

			val, _ := view.roomsModel.GetValue(iter, mucListRoomsIndexJid)
			v, _ := val.GetString()
			_, ok = view.serviceGroups[v]
			if ok {
				view.joinButton.SetSensitive(false)
			} else {
				view.joinButton.SetSensitive(true)
			}
		},
		"on_custom_service": func() {
			go u.mucUpdatePublicRoomsOn(view, ac.currentAccount())
		},
		"on_refresh": func() {
			go u.mucUpdatePublicRoomsOn(view, ac.currentAccount())
		},
	})

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}

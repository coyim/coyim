package gui

import (
	"fmt"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomListingUpdateData struct {
	iter gtki.TreeIter
	view *mucPublicRoomsView
}

func (u *gtkUI) updatedRoomListing(rl *muc.RoomListing, data interface{}) {
	d := data.(*roomListingUpdateData)
	_ = d.view.roomsModel.SetValue(d.iter, 3, rl.Description)
	_ = d.view.roomsModel.SetValue(d.iter, 4, rl.Occupants)
}

type mucPublicRoomsView struct {
	builder *builder

	dialog       gtki.Dialog    `gtk-widget:"MUCPublicRooms"`
	model        gtki.ListStore `gtk-widget:"accounts-model"`
	accountInput gtki.ComboBox  `gtk-widget:"accounts"`
	roomsModel   gtki.TreeStore `gtk-widget:"rooms-model"`
	roomsTree    gtki.TreeView  `gtk-widget:"rooms-tree"`

	accounts map[string]*account
}

func (prv *mucPublicRoomsView) init() {
	prv.builder = newBuilder("MUCPublicRoomsDialog")
	panicOnDevError(prv.builder.bindObjects(prv))
}

func (prv *mucPublicRoomsView) initAccounts(accounts []*account) {
	prv.accounts = make(map[string]*account)
	for _, acc := range accounts {
		iter := prv.model.Append()
		_ = prv.model.SetValue(iter, 0, acc.session.GetConfig().Account)
		_ = prv.model.SetValue(iter, 1, acc.session.GetConfig().ID())
		prv.accounts[acc.session.GetConfig().ID()] = acc
	}

	if len(accounts) > 0 {
		prv.accountInput.SetActive(0)
	}
}

// TODO: it should be possible to put in your own custom chat service as well

func (u *gtkUI) mucUpdatePublicRoomsOn(view *mucPublicRoomsView, a *account) {
	view.roomsModel.Clear()
	// TODO: it should show a spinner or something
	res, ec := a.session.GetRooms(jid.Parse(a.session.GetConfig().Account).Host())
	go func() {
		for {
			select {
			case rl, ok := <-res:
				if !ok {
					return
				}

				iter := view.roomsModel.Append(nil)
				_ = view.roomsModel.SetValue(iter, 0, string(rl.Jid.Local()))
				_ = view.roomsModel.SetValue(iter, 1, rl.Name)
				_ = view.roomsModel.SetValue(iter, 2, rl.Service.String())

				rl.OnUpdate(u.updatedRoomListing, &roomListingUpdateData{iter, view})
			case e, ok := <-ec:
				if !ok {
					return
				}
				if e != nil {
					fmt.Printf("Had an error: %v\n", e)
				}
				return
			}
		}
	}()
}

func (u *gtkUI) mucShowPublicRooms() {
	view := &mucPublicRoomsView{}
	view.init()

	accounts := u.getAllConnectedAccounts()
	view.initAccounts(accounts)

	view.builder.ConnectSignals(map[string]interface{}{
		"on_cancel_signal": view.dialog.Destroy,
		"on_join_signal":   func() {},
		"on_accounts_changed": func() {
			act := view.accountInput.GetActive()
			if act >= 0 && act < len(accounts) {
				u.mucUpdatePublicRoomsOn(view, accounts[act])
			}
		},
	})

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()

	if len(accounts) > 0 {
		u.mucUpdatePublicRoomsOn(view, accounts[0])
	}
}

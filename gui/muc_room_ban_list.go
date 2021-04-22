package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

func (v *roomView) onModifyBanList() {
	bl := v.newRoomBanListView()
	bl.show()
}

const (
	roomBanListAccountIndex int = iota
	roomBanListAffiliationIndex
	roomBanListReasonIndex
)

type roomBanListView struct {
	roomView *roomView

	dialog             gtki.Window   `gtk-widget:"ban-list-window"`
	addEntryButton     gtki.Button   `gtk-widget:"ban-list-add-entry-button"`
	removeEntryButton  gtki.Button   `gtk-widget:"ban-list-remove-entry-button"`
	list               gtki.TreeView `gtk-widget:"ban-list-treeview"`
	listContent        gtki.Overlay  `gtk-widget:"ban-list-content-overlay"`
	listView           gtki.Overlay  `gtk-widget:"ban-list-overlay-view"`
	listLoadingView    gtki.Box      `gtk-widget:"ban-list-loading-view"`
	noEntriesView      gtki.Box      `gtk-widget:"ban-list-no-entries-view"`
	noEntriesErrorView gtki.Box      `gtk-widget:"ban-list-error-view"`
	applyButton        gtki.Button   `gtk-widget:"ban-list-apply-changes-button"`

	listModel     gtki.ListStore
	cancelChannel chan bool

	log coylog.Logger
}

func (v *roomView) newRoomBanListView() *roomBanListView {
	bl := &roomBanListView{
		roomView: v,
		log: v.log.WithField("where", "roomBanListView"),
	}

	bl.initBuilder()
	bl.initDefaults()
	bl.initBanListModel()

	return bl
}

func (bl *roomBanListView) initBuilder() {
	builder := newBuilder("MUCRoomBannedUsersDialog")
	panicOnDevError(builder.bindObjects(bl))

	builder.ConnectSignals(map[string]interface{}{
		"on_no_items_try_again_clicked": bl.requestBanListAgain,
		"on_error_try_again_clicked":    bl.requestBanListAgain,
		"on_cancel_clicked":             bl.onCancel,
	})
}

func (bl *roomBanListView) initDefaults() {
	bl.dialog.SetTransientFor(bl.roomView.mainWindow())

	bl.addEntryButton.SetSensitive(false)
	bl.removeEntryButton.SetSensitive(false)
	bl.applyButton.SetSensitive(false)
}

func (bl *roomBanListView) initBanListModel() {
	model, _ := g.gtk.ListStoreNew(
		// the user's jid
		glibi.TYPE_STRING,
		// the user's affiliation
		glibi.TYPE_STRING,
		// the reason
		glibi.TYPE_STRING,
	)

	bl.listModel = model
	bl.list.SetModel(bl.listModel)
}

// addListItem MUST be called from the UI thread
func (bl *roomBanListView) addListItem(itm *muc.RoomBanListItem) {
	iter := bl.listModel.Append()

	bl.listModel.SetValue(iter, roomBanListAccountIndex, itm.Jid.String())
	bl.listModel.SetValue(iter, roomBanListAffiliationIndex, affiliationDisplayName(itm.Affiliation))
	bl.listModel.SetValue(iter, roomBanListReasonIndex, itm.Reason)
}

// show MUST be called from the UI thread
func (bl *roomBanListView) show() {
	bl.refreshBanList()
	bl.dialog.Show()
}

// refreshBanList MUST be called from the UI thread
func (bl *roomBanListView) refreshBanList() {
	bl.listModel.Clear()
	bl.showLoadingView()

	go bl.requestBanList()
}

// requestBanListAgain MUST be called from the UI thread
func (bl *roomBanListView) requestBanListAgain() {
	bl.listView.Show()
	bl.noEntriesView.Hide()
	bl.noEntriesErrorView.Hide()

	bl.refreshBanList()
}

// showLoadingView MUST be called from the UI thread
func (bl *roomBanListView) showLoadingView() {
	bl.listLoadingView.Show()
}

// hideLoadingView MUST be called from the UI thread
func (bl *roomBanListView) hideLoadingView() {
	bl.listLoadingView.Hide()
}

// requestBanList MUST NOT be called from the UI thread
func (bl *roomBanListView) requestBanList() {
	bl.cancelChannel = make(chan bool)

	blc, ec := bl.roomView.account.session.GetRoomBanList(bl.roomView.roomID())

	go func() {
		for {
			select {
			case itm, ok := <-blc:
				if !ok {
					doInUIThread(func() {
						bl.hideLoadingView()

						if _, ok := bl.listModel.GetIterFirst(); !ok {
							bl.listView.Hide()
							bl.noEntriesView.Show()
						}
					})
					return
				}

				doInUIThread(func() {
					bl.addListItem(itm)
				})

			case err := <-ec:
				bl.roomView.log.WithError(err).Error("Something happened when requesting the banned users list")
				doInUIThread(func() {
					bl.hideLoadingView()
					bl.listView.Hide()
					bl.noEntriesErrorView.Show()
				})
				return

			case <-bl.cancelChannel:
				return
			}
		}
	}()
}

// onCancel MUST be called from the UI thread
func (bl *roomBanListView) onCancel() {
	go func() {
		if bl.cancelChannel != nil {
			bl.cancelChannel <- true
		}
	}()

	bl.dialog.Destroy()
}

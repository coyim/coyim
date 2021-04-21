package gui

import "github.com/coyim/gotk3adapter/gtki"

func (v *roomView) onModifyBanList() {
	bl := v.newRoomBanListView()
	bl.show()
}

type roomBanListView struct {
	roomView *roomView

	dialog                 gtki.Window         `gtk-widget:"ban-list-window"`
	addEntryButton         gtki.Button         `gtk-widget:"ban-list-add-entry-button"`
	removeEntryButton      gtki.Button         `gtk-widget:"ban-list-remove-entry-button"`
	removeEntryButtonLabel gtki.Label          `gtk-widget:"ban-list-remove-entry-label"`
	listScrolledWindow     gtki.ScrolledWindow `gtk-widget:"ban-list-scrolled-window"`
	listLoadingView        gtki.Box            `gtk-widget:"ban-list-loading-view"`
	noEntriesView          gtki.Box            `gtk-widget:"ban-list-no-entries-view"`
	applyButton            gtki.Button         `gtk-widget:"ban-list-apply-changes-button"`
}

func (v *roomView) newRoomBanListView() *roomBanListView {
	bl := &roomBanListView{
		roomView: v,
	}

	bl.initBuilder()
	bl.initDefaults()

	return bl
}

func (bl *roomBanListView) initBuilder() {
	builder := newBuilder("MUCRoomBannedUsersDialog")
	panicOnDevError(builder.bindObjects(bl))
}

func (bl *roomBanListView) initDefaults() {
	bl.dialog.SetTransientFor(bl.roomView.mainWindow())

	bl.addEntryButton.SetSensitive(false)
	bl.removeEntryButton.SetSensitive(false)
	bl.applyButton.SetSensitive(false)
}

// show MUST be called from the UI thread
func (bl *roomBanListView) show() {
	bl.dialog.Show()
}

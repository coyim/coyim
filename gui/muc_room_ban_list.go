package gui

import (
	"fmt"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

func (v *roomView) onModifyBanList() {
	bl := v.newRoomBanListView()
	bl.show()
}

const (
	roomBanListAccountIndex int = iota
	roomBanListAffiliationIndex
	roomBanListReasonIndex
	roomBanListAffiliationNameIndex
)

type roomBanListView struct {
	roomView *roomView

	dialog             gtki.Window        `gtk-widget:"ban-list-window"`
	addEntryButton     gtki.Button        `gtk-widget:"ban-list-add-entry-button"`
	removeEntryButton  gtki.Button        `gtk-widget:"ban-list-remove-entry-button"`
	list               gtki.TreeView      `gtk-widget:"ban-list-treeview"`
	listSelection      gtki.TreeSelection `gtk-widget:"ban-list-treeview-selection"`
	listContent        gtki.Overlay       `gtk-widget:"ban-list-content-overlay"`
	listView           gtki.Overlay       `gtk-widget:"ban-list-overlay-view"`
	listLoadingView    gtki.Box           `gtk-widget:"ban-list-loading-view"`
	noEntriesView      gtki.Box           `gtk-widget:"ban-list-no-entries-view"`
	noEntriesErrorView gtki.Box           `gtk-widget:"ban-list-error-view"`
	applyButton        gtki.Button        `gtk-widget:"ban-list-apply-changes-button"`
	notificationsBox   gtki.Box           `gtk-widget:"notifications-box"`
	spinnerBox         gtki.Box           `gtk-widget:"spinner-box"`

	listModel       gtki.ListStore
	originalBanList []*muc.RoomBanListItem
	notifications   *notificationsComponent
	spinner         *spinner
	cancelChannel   chan bool

	log coylog.Logger
}

func (v *roomView) newRoomBanListView() *roomBanListView {
	bl := &roomBanListView{
		roomView: v,
		log:      v.log.WithField("where", "roomBanListView"),
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
		"on_jid_edited":                 bl.onUserJidEdited,
		"on_reason_edited":              bl.onReasonEdited,
		"on_add_item":                   bl.onAddNewItem,
		"on_selection_changed":          bl.onSelectionChanged,
		"on_remove_item":                bl.onRemoveItem,
		"on_apply_clicked":              bl.onApplyChanges,
		"on_cancel_clicked":             bl.onCancel,
	})
}

func (bl *roomBanListView) initDefaults() {
	bl.dialog.SetTransientFor(bl.roomView.mainWindow())

	bl.notifications = bl.roomView.u.newNotificationsComponent()
	bl.notificationsBox.Add(bl.notifications.box)

	bl.spinner = bl.roomView.u.newSpinnerComponent()
	bl.spinnerBox.Add(bl.spinner.spinner())

	bl.disableButtonsAndInteractions()
}

func (bl *roomBanListView) initBanListModel() {
	model, _ := g.gtk.ListStoreNew(
		// the user's jid
		glibi.TYPE_STRING,
		// the user's affiliation
		glibi.TYPE_STRING,
		// the reason
		glibi.TYPE_STRING,
		// the affiliation name
		glibi.TYPE_STRING,
	)

	bl.listModel = model
	bl.list.SetModel(bl.listModel)
}

// addListItem MUST be called from the UI thread
func (bl *roomBanListView) addListItem(itm *muc.RoomBanListItem) gtki.TreeIter {
	iter := bl.listModel.Append()

	jid := ""
	if itm.Jid != nil {
		jid = fmt.Sprintf("%s", itm.Jid)
	}

	affiliation := ""
	affiliationName := data.AffiliationOutcast
	if itm.Affiliation != nil {
		affiliation = affiliationDisplayName(itm.Affiliation)
		affiliationName = itm.Affiliation.Name()
	}

	bl.listModel.SetValue(iter, roomBanListAccountIndex, jid)
	bl.listModel.SetValue(iter, roomBanListAffiliationIndex, affiliation)
	bl.listModel.SetValue(iter, roomBanListReasonIndex, itm.Reason)
	bl.listModel.SetValue(iter, roomBanListAffiliationNameIndex, affiliationName)

	return iter
}

// refreshBanList MUST be called from the UI thread
func (bl *roomBanListView) refreshBanList() {
	bl.listModel.Clear()
	bl.disableButtonsAndInteractions()
	bl.showLoadingAndListViews()

	go bl.requestBanList()
}

// requestBanListAgain MUST be called from the UI thread
func (bl *roomBanListView) requestBanListAgain() {
	bl.showLoadingAndListViews()
	bl.refreshBanList()
}

// showLoadingAndListViews MUST be called from the UI thread
func (bl *roomBanListView) showLoadingAndListViews() {
	bl.noEntriesView.Hide()
	bl.noEntriesErrorView.Hide()

	bl.listLoadingView.Show()
	bl.listView.Show()
}

// hideLoadingAndListViews MUST be called from the UI thread
func (bl *roomBanListView) hideLoadingAndListViews() {
	bl.hideLoading()
	bl.listView.Hide()
}

// hideLoadingAndShowListView MUST be called from the UI thread
func (bl *roomBanListView) hideLoadingAndShowListView() {
	bl.hideLoading()
	bl.listView.Show()
}

// hideLoading MUST be called from the UI thread
func (bl *roomBanListView) hideLoading() {
	bl.listLoadingView.Hide()
}

// hasItems MUST be called from the UI thread
func (bl *roomBanListView) hasItems() bool {
	_, ok := bl.listModel.GetIterFirst()
	return ok
}

// requestBanList MUST NOT be called from the UI thread
func (bl *roomBanListView) requestBanList() {
	bl.cancelChannel = make(chan bool)

	blc, ec := bl.roomView.account.session.GetRoomBanList(bl.roomView.roomID())

	go func() {
		defer func() {
			bl.cancelChannel = nil
		}()

		select {
		case items := <-blc:
			bl.onRequestFinish(items)
		case err := <-ec:
			bl.onRequestError(err)
		case <-bl.cancelChannel:
		}
	}()
}

// onRequestFinish MUST NOT be called from the UI thread
func (bl *roomBanListView) onRequestFinish(items []*muc.RoomBanListItem) {
	if len(items) > 0 {
		doInUIThread(func() {
			for _, itm := range items {
				_ = bl.addListItem(itm)
			}
		})
	} else {
		doInUIThread(bl.noEntriesView.Show)
	}

	bl.originalBanList = items

	bl.hideLoadingAndShowListView()
	bl.enableButtonsAndInteractions()
}

// onRequestError MUST NOT be called from the UI thread
func (bl *roomBanListView) onRequestError(err error) {
	bl.roomView.log.WithError(err).Error("Something happened when requesting the banned users list")

	doInUIThread(func() {
		bl.hideLoadingAndListViews()
		bl.noEntriesErrorView.Show()
	})
}

// onUserJidEdited MUST be called from the UI thread
func (bl *roomBanListView) onUserJidEdited(_ gtki.CellRendererText, path string, newValue string) {
	iter, err := bl.listModel.GetIterFromString(path)
	if err != nil {
		bl.log.WithFields(log.Fields{
			"path":     path,
			"newValue": newValue,
		}).WithError(err).Error("Can't get the iter to update the jid of the banned user")
		return
	}

	if newValue != "" && !jid.Parse(newValue).Valid() {
		bl.log.WithFields(log.Fields{
			"path":        path,
			"newJidValue": newValue,
		}).Error("Can't update the jid of the banned user to an invalid value")
		return
	}

	if err = bl.listModel.SetValue(iter, roomBanListAccountIndex, newValue); err != nil {
		bl.log.WithFields(log.Fields{
			"path":        path,
			"newJidValue": newValue,
		}).WithError(err).Error("Can't set the new value for the jid of the banned user")
	}

	bl.enableApplyIfConditionsAreMet()
}

// onReasonEdited MUST be called from the UI thread
func (bl *roomBanListView) onReasonEdited(_ gtki.CellRendererText, path string, newValue string) {
	iter, err := bl.listModel.GetIterFromString(path)
	if err != nil {
		bl.log.WithFields(log.Fields{
			"path":           path,
			"newReasonValue": newValue,
		}).WithError(err).Error("Can't get the iter to update the reason for the banned user")
		return
	}

	if err = bl.listModel.SetValue(iter, roomBanListReasonIndex, newValue); err != nil {
		bl.log.WithFields(log.Fields{
			"path":           path,
			"newReasonValue": newValue,
		}).WithError(err).Error("Can't set the new value for the reason of the banned user")
	}

	bl.enableApplyIfConditionsAreMet()
}

// onAddNewItem MUST be called from the UI thread
func (bl *roomBanListView) onAddNewItem() {
	bl.listView.Show()
	bl.noEntriesView.Hide()

	bl.unselectSelectedRows()
	bl.listSelection.SelectIter(bl.getEmptyRowIter())
	bl.enableApplyIfConditionsAreMet()
}

// getEmptyRowIter MUST be called from the UI thread
func (bl *roomBanListView) getEmptyRowIter() gtki.TreeIter {
	iter, ok := bl.listModel.GetIterFirst()
	for ok {
		account := bl.columnStringValueFromListModelIter(iter, roomBanListAccountIndex)
		if account == "" {
			return iter
		}

		ok = bl.listModel.IterNext(iter)
	}

	return bl.addListItem(&muc.RoomBanListItem{
		Affiliation: affiliationFromKnowString(data.AffiliationOutcast),
	})
}

// onSelectionChanged MUST be called from the UI thread
func (bl *roomBanListView) onSelectionChanged() {
	totalSelected := len(bl.getSeledtedRows())
	bl.removeEntryButton.SetSensitive(totalSelected > 0)
	bl.removeEntryButton.SetTooltipText(i18n.Local("Remove selected item"))
	if totalSelected > 1 {
		bl.removeEntryButton.SetTooltipText(i18n.Local("Remove selected items"))
	}
}

// onRemoveItem MUST be called from the UI thread
func (bl *roomBanListView) onRemoveItem() {
	for _, path := range bl.getSeledtedRows() {
		iter, _ := bl.listModel.GetIter(path)
		bl.listModel.Remove(iter)
	}

	bl.enableApplyIfConditionsAreMet()
}

// unselectSelectedRows MUST be called from the UI thread
func (bl *roomBanListView) unselectSelectedRows() {
	for _, path := range bl.getSeledtedRows() {
		bl.listSelection.UnselectPath(path)
	}
}

// getSeledtedRows MUST be called from the UI thread
func (bl *roomBanListView) getSeledtedRows() []gtki.TreePath {
	return bl.listSelection.GetSelectedRows(bl.listModel)
}

// onApplyChanges MUST be called from the UI thread
func (bl *roomBanListView) onApplyChanges() {
	if !bl.isTheListUpdated() || !bl.isTheListValid() {
		return
	}

	go bl.modifyBanList(bl.currentListFromModel())
}

// modifyBanList MUST NOT be called from the UI thread
func (bl *roomBanListView) modifyBanList(changedItems []*muc.RoomBanListItem) {
	bl.cancelChannel = make(chan bool)

	doInUIThread(func() {
		bl.disableButtonsAndInteractions()
		bl.unselectSelectedRows()
		bl.spinner.show()
	})

	rc, ec := bl.roomView.account.session.ModifyRoomBanList(bl.roomView.roomID(), changedItems)

	select {
	case <-rc:
		doInUIThread(bl.close)
	case err := <-ec:
		bl.log.WithError(err).Error("Something happened when saving the room's ban list")
		doInUIThread(func() {
			bl.notifications.notifyOnError(i18n.Local("The ban list can't be updated. Please, try again."))
			bl.enableButtonsAndInteractions()
			bl.spinner.hide()
		})
	case <-bl.cancelChannel:
	}
}

// onCancel MUST be called from the UI thread
func (bl *roomBanListView) onCancel() {
	go bl.cancelActiveRequestListening()
	bl.close()
}

// disableButtonsAndInteractions MUST be called from the UI thread
func (bl *roomBanListView) disableButtonsAndInteractions() {
	bl.listView.SetSensitive(false)

	bl.addEntryButton.SetSensitive(false)
	bl.removeEntryButton.SetSensitive(false)
	bl.applyButton.SetSensitive(false)
}

// enableButtonsAndInteractions MUST be called from the UI thread
func (bl *roomBanListView) enableButtonsAndInteractions() {
	bl.listView.SetSensitive(true)
	bl.addEntryButton.SetSensitive(true)

	bl.enableApplyIfConditionsAreMet()
}

// cancelActiveRequestListening MUST NOT be called from the UI thread
func (bl *roomBanListView) cancelActiveRequestListening() {
	if bl.cancelChannel != nil {
		bl.cancelChannel <- true
	}
}

// enableApplyIfConditionsAreMet MUST be called from the UI thread
func (bl *roomBanListView) enableApplyIfConditionsAreMet() {
	bl.applyButton.SetSensitive(bl.isTheListUpdated() && bl.isTheListValid())
}

// isTheListUpdated MUST be called from the UI thread
func (bl *roomBanListView) isTheListUpdated() bool {
	currentList := bl.currentListFromModel()

	if len(currentList) != len(bl.originalBanList) {
		return true
	}

	for idx, itm := range bl.originalBanList {
		if banListItemsAreDifferent(currentList[idx], itm) {
			return true
		}
	}

	return false
}

// isTheListValid MUST be called from the UI thread
func (bl *roomBanListView) isTheListValid() bool {
	for _, itm := range bl.currentListFromModel() {
		if itm.Jid.String() == "" {
			return false
		}
	}
	return true
}

// currentListFromModel MUST be called from the UI thread
func (bl *roomBanListView) currentListFromModel() []*muc.RoomBanListItem {
	list := []*muc.RoomBanListItem{}

	iter, ok := bl.listModel.GetIterFirst()
	for ok {
		account := bl.columnStringValueFromListModelIter(iter, roomBanListAccountIndex)
		affiliation := bl.columnStringValueFromListModelIter(iter, roomBanListAffiliationNameIndex)
		reason := bl.columnStringValueFromListModelIter(iter, roomBanListReasonIndex)

		list = append(list, &muc.RoomBanListItem{
			Jid:         jid.Parse(account),
			Affiliation: affiliationFromKnowString(affiliation),
			Reason:      reason,
		})

		ok = bl.listModel.IterNext(iter)
	}

	return list
}

// columnValueFromListModelIter MUST be called from the UI thread
func (bl *roomBanListView) columnStringValueFromListModelIter(iter gtki.TreeIter, column int) string {
	v, _ := bl.listModel.GetValue(iter, column)
	s, _ := v.GetString()

	return s
}

// show MUST be called from the UI thread
func (bl *roomBanListView) show() {
	bl.refreshBanList()
	bl.dialog.Show()
}

// close MUST be called from the UI thread
func (bl *roomBanListView) close() {
	bl.dialog.Destroy()
}

func affiliationFromKnowString(a string) data.Affiliation {
	affiliation, _ := data.AffiliationFromString(a)
	return affiliation
}

func banListItemsAreDifferent(itm1, itm2 *muc.RoomBanListItem) bool {
	return itm1.Jid.String() != itm2.Jid.String() ||
		itm1.Affiliation.IsDifferentFrom(itm2.Affiliation) ||
		itm1.Reason != itm2.Reason
}

package gui

import (
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomPositionList struct {
	list muc.RoomOccupantItemList
	sync.RWMutex
}

func newRoomPositionList() *roomPositionList {
	return &roomPositionList{}
}

func (pl *roomPositionList) set(list muc.RoomOccupantItemList) {
	pl.Lock()
	defer pl.Unlock()

	pl.list = list
}

func (pl *roomPositionList) merge(list muc.RoomOccupantItemList) {
	pl.set(append(pl.positions(), list...))
}

func (pl *roomPositionList) positions() muc.RoomOccupantItemList {
	pl.RLock()
	defer pl.RUnlock()

	ret := muc.RoomOccupantItemList{}
	for _, itm := range pl.list {
		ret = append(ret, itm)
	}

	return ret
}

type roomPositions struct {
	banned *roomPositionList
	none   *roomPositionList
}

func newRoomPositions() *roomPositions {
	return &roomPositions{
		banned: newRoomPositionList(),
		none:   newRoomPositionList(),
	}
}

func (rp *roomPositions) bannedList() muc.RoomOccupantItemList {
	return rp.banned.positions()
}

func (rp *roomPositions) setBanList(list muc.RoomOccupantItemList) {
	rp.banned.set(list)
}

func (rp *roomPositions) updateRemovedOccupantList(list muc.RoomOccupantItemList) {
	rp.none.merge(list)
}

func (rp *roomPositions) positionsToUpdate() muc.RoomOccupantItemList {
	return append(rp.banned.positions(), rp.none.positions()...)
}

func (v *roomView) onRoomPositionsView() {
	rp := v.newRoomPositionsView()
	rp.show()
}

const (
	roomBanListAccountIndex int = iota
	roomBanListAffiliationIndex
	roomBanListReasonIndex
	roomBanListAffiliationNameIndex
)

type roomPositionsView struct {
	roomView              *roomView
	banned                muc.RoomOccupantItemList
	none                  muc.RoomOccupantItemList
	onUpdateOccupantLists *callbacksSet

	dialog  gtki.Window `gtk-widget:"positions-window"`
	content gtki.Box    `gtk-widget:"content"`

	log coylog.Logger
}

func (v *roomView) newRoomPositionsView() *roomPositionsView {
	rpv := &roomPositionsView{
		roomView:              v,
		onUpdateOccupantLists: newCallbacksSet(),
	}

	rpv.log = v.log.WithFields(log.Fields{
		"room":  v.roomID(),
		"where": "roomPositionsView",
	})

	rpv.initBuilder()
	rpv.initDefaults()

	return rpv
}

func (rpv *roomPositionsView) initBuilder() {
	builder := newBuilder("MUCRoomPositionsDialog")
	panicOnDevError(builder.bindObjects(rpv))

	builder.ConnectSignals(map[string]interface{}{
		"on_apply":  rpv.onApply,
		"on_cancel": rpv.dialog.Destroy,
	})
}

func (rpv *roomPositionsView) initDefaults() {
	rpv.dialog.SetTransientFor(rpv.roomView.mainWindow())
	mucStyles.setRoomConfigPageStyle(rpv.content)
}

// onApply MUST be called from the UI thread
func (rpv *roomPositionsView) onApply() {
	rpv.onUpdateOccupantLists.invokeAll()

	rpv.dialog.Destroy()
	rpv.roomView.loadingViewOverlay.onRoomPositionsUpdate()

	rc, ec := rpv.roomView.account.session.UpdateOccupantAffiliations(rpv.roomView.roomID(), rpv.positionsToUpdate())
	go func() {
		select {
		case <-rc:
			doInUIThread(func() {
				rpv.roomView.notifications.info(roomNotificationOptions{
					message:   i18n.Local("The positions were updated."),
					closeable: true,
				})
				rpv.roomView.loadingViewOverlay.hide()
			})
		case <-ec:
			doInUIThread(func() {
				rpv.roomView.notifications.error(roomNotificationOptions{
					message:   i18n.Local("Unable to update positions."),
					closeable: true,
				})
				rpv.roomView.loadingViewOverlay.hide()
			})
		}
	}()
}

// requestOccupantsByAffiliation MUST NOT be called from the UI thread
func (rp *roomPositionsView) requestOccupantsByAffiliation(a data.Affiliation, onSuccess func(muc.RoomOccupantItemList), onError func()) {
	rc, ec := rp.roomView.account.session.GetRoomOccupantsByAffiliation(rp.roomView.roomID(), a)

	select {
	case ol := <-rc:
		onSuccess(ol)
	case <-ec:
		onError()
	}
}

// requestRoomPositions MUST NOT be called from the UI thread
func (rp *roomPositionsView) requestRoomPositions(onSuccess func(), onError func()) {
	rp.requestOccupantsByAffiliation(&data.OutcastAffiliation{},
		func(items muc.RoomOccupantItemList) {
			rp.roomPositions.setBanList(items)

			rp.addPositionComponent(newRoomConfigPositions(&data.OutcastAffiliation{}, rp.roomPositions.bannedList(),
				rp.roomPositions.setBanList, rp.roomPositions.updateRemovedOccupantList, func() {}))
		},
		onError)

	onSuccess()
}

// addPositionComponent MUST be called from the UI thread
func (rpv *roomPositionsView) addPositionComponent(positionComponent hasRoomConfigFormField) {
	rpv.content.Add(positionComponent.fieldWidget())
	positionComponent.refreshContent()
	rpv.onUpdateOccupantLists.add(positionComponent.updateFieldValue)
}

// show MUST be called from the UI thread
func (rpv *roomPositionsView) show() {
	rpv.roomView.loadingViewOverlay.onRoomPositionsRequest()

	go func() {
		rpv.requestRoomPositions(
			func() {
				doInUIThread(func() {
					rpv.roomView.loadingViewOverlay.hide()
					rpv.dialog.Show()
				})
			},
			func() {
				doInUIThread(func() {
					rpv.roomView.loadingViewOverlay.hide()
					rpv.roomView.notifications.error(roomNotificationOptions{
						message:   i18n.Local("We couldn't get the occupants by affiliation"),
						closeable: true,
					})
				})
			},
		)
	}()
}

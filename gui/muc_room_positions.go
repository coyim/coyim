package gui

import (
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomPositions struct {
	banned muc.RoomOccupantItemList
	none   muc.RoomOccupantItemList

	sync.Mutex
}

func newRoomPositions() *roomPositions {
	return &roomPositions{}
}

func (rp *roomPositions) bannedList() muc.RoomOccupantItemList {
	return rp.banned
}

func (rp *roomPositions) setBanList(banned muc.RoomOccupantItemList) {
	rp.Lock()
	defer rp.Unlock()

	rp.banned = banned
}

func (rp *roomPositions) updateRemovedOccupantList(occupantsRemoved muc.RoomOccupantItemList) {
	rp.Lock()
	defer rp.Unlock()

	rp.none = append(rp.none, occupantsRemoved...)
}

func (rp *roomPositions) positionsToUpdate() muc.RoomOccupantItemList {
	rp.Lock()
	defer rp.Unlock()

	positionsToUpdate := muc.RoomOccupantItemList{}
	positionsToUpdate = append(positionsToUpdate, rp.banned...)
	positionsToUpdate = append(positionsToUpdate, rp.none...)

	return positionsToUpdate

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
	roomPositions         *roomPositions
	onUpdateOccupantLists *callbacksSet

	dialog  gtki.Window `gtk-widget:"positions-window"`
	content gtki.Box    `gtk-widget:"content"`

	log coylog.Logger
}

func (v *roomView) newRoomPositionsView() *roomPositionsView {
	rp := &roomPositionsView{
		roomView:              v,
		roomPositions:         newRoomPositions(),
		onUpdateOccupantLists: newCallbacksSet(),
		log:                   v.log.WithField("where", "roomPositionsView"),
	}

	rp.initBuilder()
	rp.initDefaults()

	return rp
}

func (rp *roomPositionsView) initBuilder() {
	builder := newBuilder("MUCRoomPositionsDialog")
	panicOnDevError(builder.bindObjects(rp))

	builder.ConnectSignals(map[string]interface{}{
		"on_apply":  rp.onApply,
		"on_cancel": rp.onCancel,
	})
}

func (rp *roomPositionsView) initDefaults() {
	rp.dialog.SetTransientFor(rp.roomView.mainWindow())
	mucStyles.setRoomConfigPageStyle(rp.content)
}

// onCancel MUST be called from the UI thread
func (rp *roomPositionsView) onCancel() {
	rp.dialog.Destroy()
}

// onCancel MUST be called from the UI thread
func (rp *roomPositionsView) onApply() {
	rp.onUpdateOccupantLists.invokeAll()

	doInUIThread(func() {
		rp.dialog.Destroy()
		rp.roomView.loadingViewOverlay.onRoomPositionsUpdate()
	})

	rc, ec := rp.roomView.account.session.UpdateOccupantAffiliations(rp.roomView.roomID(), rp.roomPositions.positionsToUpdate())
	go func() {
		select {
		case <-rc:
			doInUIThread(func() {
				rp.roomView.notifications.info(roomNotificationOptions{
					message:   i18n.Local("The positions list has been updated."),
					closeable: true,
				})
				rp.roomView.loadingViewOverlay.hide()
			})
		case <-ec:
			doInUIThread(func() {
				rp.roomView.notifications.error(roomNotificationOptions{
					message:   i18n.Local("The positions list couldn't be updated."),
					closeable: true,
				})
				rp.roomView.loadingViewOverlay.hide()
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

func (rp *roomPositionsView) addPositionComponent(positionComponent hasRoomConfigFormField) {
	rp.content.Add(positionComponent.fieldWidget())
	positionComponent.refreshContent()
	rp.onUpdateOccupantLists.add(positionComponent.updateFieldValue)
}

// show MUST be called from the UI thread
func (rp *roomPositionsView) show() {
	rp.roomView.loadingViewOverlay.onRoomPositionsRequest()

	go func() {
		rp.requestRoomPositions(
			func() {
				doInUIThread(func() {
					rp.roomView.loadingViewOverlay.hide()
					rp.dialog.Show()
				})
			},
			func() {
				doInUIThread(func() {
					rp.roomView.loadingViewOverlay.hide()
					rp.roomView.notifications.error(roomNotificationOptions{
						message:   i18n.Local("We couldn't get the occupants by affiliation"),
						closeable: true,
					})
				})
			},
		)
	}()
}

package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

// onRoomPositionsView MUST be called from the UI thread
func (v *roomView) onRoomPositionsView() {
	rpv := v.newRoomPositionsView()
	rpv.show()
}

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

// setBanList MUST be called from the UI thread
func (rpv *roomPositionsView) setBanList(list muc.RoomOccupantItemList) {
	rpv.banned = append(rpv.banned, list...)
}

// updateRemovedOccupantList MUST be called from the UI thread
func (rpv *roomPositionsView) updateRemovedOccupantList(list muc.RoomOccupantItemList) {
	rpv.none = append(rpv.none, list...)
}

func (rpv *roomPositionsView) positionsToUpdate() muc.RoomOccupantItemList {
	return append(rpv.banned, rpv.none...)
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

// requestRoomPositions MUST NOT be called from the UI thread
func (rpv *roomPositionsView) requestRoomPositions(onSuccess func(), onError func()) {
	rc, ec := rpv.roomView.account.session.GetRoomOccupantsByAffiliation(rpv.roomView.roomID(), &data.OutcastAffiliation{})

	select {
	case ol := <-rc:
		doInUIThread(
			func() {
				rpv.setBanList(ol)
				rpv.addPositionComponent(newRoomConfigPositions(&data.OutcastAffiliation{}, rpv.banned,
					rpv.setBanList, rpv.updateRemovedOccupantList, func() {}))
			})
	case <-ec:
		onError()
	}

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

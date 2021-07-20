package gui

import (
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomPositions struct {
	owners muc.RoomOccupantItemList
	admins muc.RoomOccupantItemList
	banned muc.RoomOccupantItemList
	none   muc.RoomOccupantItemList

	sync.Mutex
}

func newRoomPositions() *roomPositions {
	return &roomPositions{}
}

func (rp *roomPositions) ownersList() muc.RoomOccupantItemList {
	return rp.owners
}

func (rp *roomPositions) setOwnerList(owners muc.RoomOccupantItemList) {
	rp.Lock()
	defer rp.Unlock()

	rp.owners = owners
}

func (rp *roomPositions) adminsList() muc.RoomOccupantItemList {
	return rp.admins
}

func (rp *roomPositions) setAdminList(admins muc.RoomOccupantItemList) {
	rp.Lock()
	defer rp.Unlock()

	rp.admins = admins
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
	roomView      *roomView
	roomPositions *roomPositions

	dialog  gtki.Window `gtk-widget:"positions-window"`
	content gtki.Box    `gtk-widget:"content"`

	log coylog.Logger
}

func (v *roomView) newRoomPositionsView() *roomPositionsView {
	rp := &roomPositionsView{
		roomView:      v,
		roomPositions: newRoomPositions(),
		log:           v.log.WithField("where", "roomPositionsView"),
	}

	rp.initBuilder()
	rp.initDefaults()

	return rp
}

func (rp *roomPositionsView) initBuilder() {
	builder := newBuilder("MUCRoomPositionsDialog")
	panicOnDevError(builder.bindObjects(rp))
}

func (rp *roomPositionsView) initDefaults() {
	rp.dialog.SetTransientFor(rp.roomView.mainWindow())
	mucStyles.setRoomConfigPageStyle(rp.content)
}

// show MUST be called from the UI thread
func (rp *roomPositionsView) show() {
	rp.dialog.Show()
}

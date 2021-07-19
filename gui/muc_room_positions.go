package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/gotk3adapter/gtki"
)

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
	roomView *roomView

	dialog  gtki.Window `gtk-widget:"positions-window"`
	content gtki.Box    `gtk-widget:"content"`

	log coylog.Logger
}

func (v *roomView) newRoomPositionsView() *roomPositionsView {
	rp := &roomPositionsView{
		roomView: v,
		log:      v.log.WithField("where", "roomPositionsView"),
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

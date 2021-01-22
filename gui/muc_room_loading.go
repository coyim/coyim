package gui

import (
	"github.com/coyim/coyim/i18n"
)

type roomViewLoadingOverlay struct {
	*loadingOverlayComponent
}

func (v *roomView) newRoomViewLoadingOverlay() *roomViewLoadingOverlay {
	o := &roomViewLoadingOverlay{
		v.u.newLoadingOverlayComponent(),
	}

	v.overlay.AddOverlay(o.overlay)

	return o
}

// onRoomDiscoInfoLoad MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onRoomDiscoInfoLoad() {
	lo.setTitle(i18n.Local("Loading room information"))
	lo.setDescription(i18n.Local("This will only take a few moments."))
	lo.setSolid()
	lo.show()
}

// onJoinRoom MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onJoinRoom() {
	lo.setTitle(i18n.Local("Joining room..."))
	lo.setSolid()
	lo.show()
}

// onRoomDestroy MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onRoomDestroy() {
	lo.setTitle(i18n.Local("Destroying room..."))
	lo.setTransparent()
	lo.show()
}

// onRoomAffiliationConfirmation MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onOccupantAffiliationUpdate() {
	lo.setTitle(i18n.Local("Updating occupant affiliation..."))
	lo.setTransparent()
	lo.show()
}

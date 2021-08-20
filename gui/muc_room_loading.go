package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewLoadingOverlay struct {
	*loadingOverlayComponent
}

func (v *roomView) newRoomViewLoadingOverlay() *roomViewLoadingOverlay {
	return &roomViewLoadingOverlay{
		v.u.newLoadingOverlayComponent(),
	}
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

// onRoomReconnect MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onRoomReconnect() {
	lo.setTitle(i18n.Local("Connecting to the room..."))
	lo.setSolid()
	lo.show()
}

// onRoomReconnectTryAgain MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onRoomReconnectTryAgain() {
	lo.setTitle(i18n.Local("Trying to connect the room..."))
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
	lo.setTitle(i18n.Local("Updating position..."))
	lo.setTransparent()
	lo.show()
}

// onOccupantRoleUpdate MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onOccupantRoleUpdate(role data.Role) {
	m := i18n.Local("Updating role...")
	if role.IsNone() {
		m = i18n.Local("Expelling person from the room...")
	}
	lo.setTitle(m)
	lo.setTransparent()
	lo.show()
}

// onRoomConfigurationRequest MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onRoomConfigurationRequest() {
	lo.setTitle(i18n.Local("Loading room configuration..."))
	lo.setTransparent()
	lo.show()
}

// onRoomPositionsRequest MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onRoomPositionsRequest() {
	lo.setTitle(i18n.Local("Loading room positions..."))
	lo.setTransparent()
	lo.show()
}

// onRoomPositionsUpdate MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onRoomPositionsUpdate() {
	lo.setTitle(i18n.Local("Saving room positions..."))
	lo.setTransparent()
	lo.show()
}

func (lo *roomViewLoadingOverlay) view() gtki.Overlay {
	return lo.overlay
}

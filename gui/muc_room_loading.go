package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewLoadingOverlay struct {
	overlay     gtki.Overlay
	box         gtki.Box
	title       gtki.Label
	description gtki.Label
}

func (v *roomView) newRoomViewLoadingOverlay(o gtki.Overlay, b gtki.Box, t gtki.Label, d gtki.Label) *roomViewLoadingOverlay {
	lo := &roomViewLoadingOverlay{o, b, t, d}
	lo.initDefaults()
	return lo
}

func (lo *roomViewLoadingOverlay) initDefaults() {
	mucStyles.setRoomLoadingViewOverlayTitleStyle(lo.title)
}

func (lo *roomViewLoadingOverlay) setTransparent() {
	mucStyles.setRoomLoadingViewOverlayTransparentStyle(lo.box)
}

func (lo *roomViewLoadingOverlay) setSolid() {
	mucStyles.setRoomLoadingViewOverlaySolidStyle(lo.box)
}

// onRoomDiscoInfoLoad MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onRoomDiscoInfoLoad() {
	lo.setTitle(i18n.Local("Loading room information"))
	lo.setDescription(i18n.Local("This will only take a few moments."))
	lo.setSolid()
	lo.show()
}

// onRoomDestroy MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onRoomDestroy() {
	lo.setTitle(i18n.Local("Destroying room..."))
	lo.setTransparent()
	lo.show()
}

// show MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) show() {
	lo.overlay.Show()
}

// showWithMessage MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) showWithMessage(m string) {
	lo.setTitle(m)
	lo.show()
}

func (lo *roomViewLoadingOverlay) setTitle(t string) {
	lo.title.SetLabel(t)
	lo.title.Show()
}

func (lo *roomViewLoadingOverlay) setDescription(d string) {
	lo.description.SetLabel(d)
	lo.description.Show()
}

// hide MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) hide() {
	lo.setTitle("")
	lo.setDescription("")

	lo.title.Hide()
	lo.description.Hide()

	lo.overlay.Hide()
}

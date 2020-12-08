package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewLoadingOverlay struct {
	overlay gtki.Overlay
	box     gtki.Box
	label   gtki.Label
}

func (v *roomView) newRoomViewLoadingOverlay(o gtki.Overlay, b gtki.Box, l gtki.Label) *roomViewLoadingOverlay {
	lo := &roomViewLoadingOverlay{o, b, l}
	lo.initDefaults()

	return lo
}

func (lo *roomViewLoadingOverlay) initDefaults() {
	mucStyles.setRoomLoadingViewOverlayBoxStyle(lo.box)
}

// onRoomDiscoInfoLoad MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) onRoomDiscoInfoLoad() {
	lo.showWithMessage(i18n.Local("Loading room info...."))
}

// show MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) show() {
	lo.overlay.Show()
}

// showWithMessage MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) showWithMessage(m string) {
	lo.label.SetLabel(m)
	lo.label.Show()
	lo.show()
}

// hide MUST be called from the UI thread
func (lo *roomViewLoadingOverlay) hide() {
	lo.label.SetLabel("")
	lo.label.Hide()
	lo.overlay.Hide()
}

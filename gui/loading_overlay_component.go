package gui

import "github.com/coyim/gotk3adapter/gtki"

type loadingOverlayComponent struct {
	overlay     gtki.Overlay `gtk-widget:"loading-overlay"`
	title       gtki.Label   `gtk-widget:"loading-overlay-title"`
	description gtki.Label   `gtk-widget:"loading-overlay-description"`
	content     gtki.Box     `gtk-widget:"loading-overlay-content"`
	box         gtki.Box     `gtk-widget:"loading-overlay-box"`
}

func (u *gtkUI) newLoadingOverlayComponent() *loadingOverlayComponent {
	lo := &loadingOverlayComponent{}

	builder := newBuilder("LoadingOverlay")
	panicOnDevError(builder.bindObjects(lo))

	mucStyles.setLabelBoldStyle(lo.title)

	return lo
}

func (lo *loadingOverlayComponent) getOverlay() gtki.Overlay {
	return lo.overlay
}

// setTransparent MUST be called from the UI thread
func (lo *loadingOverlayComponent) setTransparent() {
	mucStyles.setRoomLoadingViewOverlayTransparentStyle(lo.box)
	mucStyles.setRoomLoadingViewOverlayContentTransparentStyle(lo.content)
}

// setSolid MUST be called from the UI thread
func (lo *loadingOverlayComponent) setSolid() {
	mucStyles.setRoomLoadingViewOverlaySolidStyle(lo.box)
	mucStyles.setRoomLoadingViewOverlayContentSolidStyle(lo.content)
}

// show MUST be called from the UI thread
func (lo *loadingOverlayComponent) show() {
	lo.overlay.Show()
}

// showWithMessage MUST be called from the UI thread
func (lo *loadingOverlayComponent) showWithMessage(m string) {
	lo.setTitle(m)
	lo.show()
}

// setTitle MUST be called from the UI thread
func (lo *loadingOverlayComponent) setTitle(t string) {
	lo.title.SetLabel(t)
	lo.title.Show()
}

// setDescription MUST be called from the UI thread
func (lo *loadingOverlayComponent) setDescription(d string) {
	lo.description.SetLabel(d)
	lo.description.Show()
}

// hide MUST be called from the UI thread
func (lo *loadingOverlayComponent) hide() {
	lo.setTitle("")
	lo.setDescription("")

	lo.title.Hide()
	lo.description.Hide()

	lo.overlay.Hide()
}

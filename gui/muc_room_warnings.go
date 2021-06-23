package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
)

func (v *roomView) showRoomWarnings(info data.RoomDiscoInfo) {
	v.warnings.add(
		i18n.Local("Please be aware that communication in chat rooms " +
			"is not encrypted - anyone that can intercept communication between you and " +
			"the server - and the server itself - will be able to see what you are saying " +
			"in this chat room. Only join this room and communicate here if you trust the " +
			"server to not be hostile."),
	)

	switch info.AnonymityLevel {
	case "semi":
		v.warnings.add(
			i18n.Local("This room is partially anonymous. This means that " +
				"only moderators can connect your nickname with your real username (your JID)."),
		)
	case "no":
		v.warnings.add(
			i18n.Local("This room is not anonymous. This means that any person " +
				"in this room can connect your nickname with your real username (your JID)."),
		)
	default:
		v.log.WithField("anonymityLevel", info.AnonymityLevel).Warn("Unknown anonymity " +
			"setting for room")
	}

	if info.Logged {
		v.warnings.add(
			i18n.Local("This room is publicly logged, meaning that everything " +
				"you and the others in the room say or do can be made public on a website."),
		)
	}
}

type roomViewWarning struct {
	text string

	bar     gtki.Box   `gtk-widget:"warning-infobar"`
	message gtki.Label `gtk-widget:"message"`
}

func newRoomViewWarning(text string) *roomViewWarning {
	w := &roomViewWarning{
		text: text,
	}

	builder := newBuilder("MUCRoomWarning")
	panicOnDevError(builder.bindObjects(w))

	w.message.SetText(w.text)

	return w
}

type roomViewWarningsInfoBar struct {
	*notificationBar
}

func (v *roomView) newRoomViewWarningsInfoBar() *roomViewWarningsInfoBar {
	ib := &roomViewWarningsInfoBar{
		v.u.newNotificationBar(i18n.Local("Check out the security properties of this room!"), gtki.MESSAGE_WARNING),
	}

	showWarningsButton, _ := g.gtk.ButtonNewWithLabel(i18n.Local("Details"))
	showWarningsButton.Connect("clicked", v.showWarnings)

	ib.addActionWidget(showWarningsButton, gtki.RESPONSE_NONE)

	return ib
}

type roomViewWarningsDirection string

const (
	roomViewWarningsPreviousDirection roomViewWarningsDirection = "previous"
	roomViewWarningsNextDirection     roomViewWarningsDirection = "next"
)

type roomViewWarnings struct {
	roomView            *roomView
	warnings            []*roomViewWarning
	currentWarningIndex int

	dialog      gtki.Window `gtk-widget:"room-warnings-dialog"`
	currentText gtki.Label  `gtk-widget:"room-warnings-current-title"`
	currentInfo gtki.Label  `gtk-widget:"room-warnings-current-info"`
}

func (v *roomView) newRoomViewWarnings() *roomViewWarnings {
	vw := &roomViewWarnings{
		roomView: v,
	}

	vw.initBuilder()
	vw.initDefaults()

	return vw
}

func (vw *roomViewWarnings) initBuilder() {
	builder := newBuilder("MUCRoomWarningsDialog")
	panicOnDevError(builder.bindObjects(vw))

	builder.ConnectSignals(map[string]interface{}{
		"on_warning_go_previous_clicked": func() {
			vw.move(roomViewWarningsPreviousDirection)
		},
		"on_warning_go_next_clicked": func() {
			vw.move(roomViewWarningsNextDirection)
		},
	})
}

// move MUST be called from the UI thread
func (vw *roomViewWarnings) move(direction roomViewWarningsDirection) {
	firstWarningIndex := 0
	lastWarningIndex := vw.total() - 1
	newCurrentWarningIndex := vw.currentWarningIndex

	switch direction {
	case roomViewWarningsPreviousDirection:
		newCurrentWarningIndex = vw.currentWarningIndex - 1
		if newCurrentWarningIndex < firstWarningIndex {
			newCurrentWarningIndex = lastWarningIndex
		}

	case roomViewWarningsNextDirection:
		newCurrentWarningIndex = vw.currentWarningIndex + 1
		if newCurrentWarningIndex > lastWarningIndex {
			newCurrentWarningIndex = firstWarningIndex
		}
	}

	vw.currentWarningIndex = newCurrentWarningIndex

	vw.refresh()
}

func (vw *roomViewWarnings) initDefaults() {
	vw.dialog.SetTransientFor(vw.roomView.window)
	mucStyles.setRoomWarningsStyles(vw.dialog)
}

// add MUST be called from the UI thread
func (vw *roomViewWarnings) add(text string) {
	w := newRoomViewWarning(text)
	vw.warnings = append(vw.warnings, w)
	vw.refresh()
}

// refresh MUST be called from the UI thread
func (vw *roomViewWarnings) refresh() {
	warningText := ""
	warningInfo := ""

	if warning, ok := vw.warningByIndex(vw.currentWarningIndex); ok {
		warningText = warning.text
		warningInfo = i18n.Localf("Warning %d of %d", vw.currentWarningIndex+1, vw.total())
	}

	vw.currentText.SetText(warningText)
	vw.currentInfo.SetText(warningInfo)
}

func (wi *roomViewWarningsInfoBar) hide() {
	wi.infoBar.SetVisible(false)
}

func (vw *roomViewWarnings) warningByIndex(idx int) (*roomViewWarning, bool) {
	if vw.hasWarningIndex(idx) {
		return vw.warnings[idx], true
	}
	return nil, false
}

func (vw *roomViewWarnings) hasWarningIndex(idx int) bool {
	return idx >= 0 && idx < vw.total()
}

func (vw *roomViewWarnings) total() int {
	return len(vw.warnings)
}

// clear MUST be called from the UI thread
func (vw *roomViewWarnings) clear() {
	vw.warnings = nil
	vw.currentWarningIndex = 0
	vw.refresh()
}

// show MUST be called from the UI thread
func (vw *roomViewWarnings) show() {
	vw.dialog.Show()
}

package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewWarningIconType string

const (
	roomViewWarningIconDefault            roomViewWarningIconType = "warning_default"
	roomViewWarningIconNotEncrypted       roomViewWarningIconType = "not_encrypted"
	roomViewWarningIconPartiallyAnonymous roomViewWarningIconType = "partially_anonymous"
	roomViewWarningIconNotAnonymous       roomViewWarningIconType = "not_anonymous"
	roomViewWarningIconPubliclyLogged     roomViewWarningIconType = "publicly_logged"
)

func (icon roomViewWarningIconType) name() string {
	return string(icon)
}

// showWarnings MUST be called from the UI thread
func (v *roomView) showWarnings() {
	v.warnings.show()
}

func (v *roomView) addRoomWarningsBasedOnInfo(info data.RoomDiscoInfo) {
	v.warnings.add(
		roomViewWarningIconNotEncrypted,
		i18n.Local("Communication in this room is not encrypted"),
		i18n.Local("Please be aware that communication in chat rooms is not encrypted - "+
			"anyone that can intercept communication between you and the server - "+
			"and the server itself - will be able to see what you are saying in this "+
			"chat room. Only join this room and communicate here if you trust the server "+
			"to not be hostile."),
	)

	switch info.AnonymityLevel {
	case "semi":
		v.warnings.add(
			roomViewWarningIconPartiallyAnonymous,
			i18n.Local("Partially anonymous room"),
			i18n.Local("This room is partially anonymous. This means that only moderators "+
				"can connect your nickname with your real username (your JID)."),
		)
	case "no":
		v.warnings.add(
			roomViewWarningIconNotAnonymous,
			i18n.Local("Non-anonymous room"),
			i18n.Local("This room is not anonymous. This means that any person in this room "+
				"can connect your nickname with your real username (your JID)."),
		)
	default:
		v.log.WithField("anonymityLevel", info.AnonymityLevel).Warn("Unknown anonymity " +
			"setting for room")
	}

	if info.Logged {
		v.warnings.add(
			roomViewWarningIconPubliclyLogged,
			i18n.Local("Publicly logged room"),
			i18n.Local("This room is publicly logged, meaning that everything you and the "+
				"others in the room say or do can be made public on a website."),
		)
	}
}

type roomViewWarning struct {
	icon        roomViewWarningIconType
	title       string
	description string
}

func newRoomViewWarning(icon roomViewWarningIconType, title, description string) *roomViewWarning {
	return &roomViewWarning{
		icon,
		title,
		description,
	}
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

func (wi *roomViewWarningsInfoBar) hide() {
	wi.infoBar.SetVisible(false)
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

	dialog             gtki.Window `gtk-widget:"room-warnings-dialog"`
	currentIcon        gtki.Image  `gtk-widget:"room-warnings-current-icon"`
	currentTitle       gtki.Label  `gtk-widget:"room-warnings-current-title"`
	currentDescription gtki.Label  `gtk-widget:"room-warnings-current-description"`
	currentInfo        gtki.Label  `gtk-widget:"room-warnings-current-info"`
	movePreviousButton gtki.Button `gtk-widget:"room-warnings-move-previous-button"`
	moveNextButton     gtki.Button `gtk-widget:"room-warnings-move-next-button"`
}

func (v *roomView) newRoomViewWarnings() *roomViewWarnings {
	vw := &roomViewWarnings{
		roomView: v,
	}

	vw.initBuilder()
	vw.initShortcuts()
	vw.initDefaults()

	return vw
}

func (vw *roomViewWarnings) initBuilder() {
	builder := newBuilder("MUCRoomWarningsDialog")
	panicOnDevError(builder.bindObjects(vw))

	builder.ConnectSignals(map[string]interface{}{
		"on_warning_go_previous_clicked": vw.moveLeft,
		"on_warning_go_next_clicked":     vw.moveRight,
		"on_dialog_close":                vw.close,
	})
}

func (vw *roomViewWarnings) initShortcuts() {
	callers := map[string]func(){
		"Escape":         vw.dialog.Hide,
		"<Primary>Left":  vw.moveLeft,
		"<Primary>Right": vw.moveRight,
	}

	for sh, c := range callers {
		connectShortcut(sh, vw.dialog, simpleWindowShortcutCall(c))
	}
}

func (vw *roomViewWarnings) initDefaults() {
	vw.dialog.SetTransientFor(vw.roomView.window)
	mucStyles.setRoomWarningsStyles(vw.dialog)
}

// move MUST be called from the UI thread
func (vw *roomViewWarnings) move(direction roomViewWarningsDirection) {
	firstWarningIndex := 0
	lastWarningIndex := vw.total() - 1
	newCurrentWarningIndex := vw.currentWarningIndex

	switch direction {
	case roomViewWarningsPreviousDirection:
		newCurrentWarningIndex = vw.currentWarningIndex - 1
	case roomViewWarningsNextDirection:
		newCurrentWarningIndex = vw.currentWarningIndex + 1
	}

	if newCurrentWarningIndex >= firstWarningIndex &&
		newCurrentWarningIndex <= lastWarningIndex {
		vw.currentWarningIndex = newCurrentWarningIndex
		vw.refresh()
	}
}

func (vw *roomViewWarnings) moveLeft() {
	vw.move(roomViewWarningsPreviousDirection)
}

func (vw *roomViewWarnings) moveRight() {
	vw.move(roomViewWarningsNextDirection)
}

// add MUST be called from the UI thread
func (vw *roomViewWarnings) add(icon roomViewWarningIconType, title, description string) {
	w := newRoomViewWarning(icon, title, description)
	vw.warnings = append(vw.warnings, w)
	vw.refresh()
}

// refresh MUST be called from the UI thread
func (vw *roomViewWarnings) refresh() {
	vw.refreshCurrentWarningContent()
	vw.refreshMoveButtons()
}

// refreshCurrentWarningContent MUST be called from the UI thread
func (vw *roomViewWarnings) refreshCurrentWarningContent() {
	warningIcon := roomViewWarningIconDefault
	warningTitle := ""
	warningDescription := ""
	warningInfo := ""

	if warning, ok := vw.warningByIndex(vw.currentWarningIndex); ok {
		warningIcon = warning.icon
		warningTitle = warning.title
		warningDescription = warning.description
		warningInfo = i18n.Localf("Warning %[1]d of %[2]d", vw.currentWarningIndex+1, vw.total())
	}

	vw.currentTitle.SetText(warningTitle)
	vw.currentDescription.SetText(warningDescription)
	vw.currentInfo.SetText(warningInfo)
	vw.currentIcon.SetFromPixbuf(getMUCIconPixbuf(warningIcon.name()))
}

// refreshMoveButtons MUST be called from the UI thread
func (vw *roomViewWarnings) refreshMoveButtons() {
	totalWarnings := vw.total()
	firstWarningIndex := 0
	lastWarningIndex := totalWarnings - 1

	vw.movePreviousButton.SetSensitive(vw.currentWarningIndex > firstWarningIndex)
	vw.moveNextButton.SetSensitive(vw.currentWarningIndex < lastWarningIndex)
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
	vw.refresh()
}

// close MUST be called from the UI thread
func (vw *roomViewWarnings) close() {
	vw.dialog.Hide()
}

func simpleWindowShortcutCall(fn func()) func(gtki.Window) {
	return func(gtki.Window) {
		fn()
	}
}

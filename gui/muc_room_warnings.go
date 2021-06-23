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

type roomViewWarningsOverlay struct {
	warnings []*roomViewWarning
	onClose  func()

	box      gtki.Box      `gtk-widget:"warningsBox"`
	revealer gtki.Revealer `gtk-widget:"revealer"`
}

func (v *roomView) newRoomViewWarningsOverlay() *roomViewWarningsOverlay {
	o := &roomViewWarningsOverlay{
		onClose: v.closeNotificationsOverlay,
	}

	builder := newBuilder("MUCRoomWarningsOverlay")
	panicOnDevError(builder.bindObjects(o))

	builder.ConnectSignals(map[string]interface{}{
		"on_close": o.close,
	})

	mucStyles.setRoomWarningsBoxStyle(o.box)

	v.messagesBox.Add(o.revealer)

	return o
}

func (o *roomViewWarningsOverlay) add(text string) {
	w := newRoomViewWarning(text)
	o.warnings = append(o.warnings, w)

	mucStyles.setRoomWarningsMessageBoxStyle(w.bar)

	o.box.PackStart(w.bar, false, false, 5)

	o.box.ShowAll()
}

func (o *roomViewWarningsOverlay) show() {
	o.revealer.SetRevealChild(true)
}

func (wi *roomViewWarningsInfoBar) hide() {
	wi.infoBar.SetVisible(false)
}

func (o *roomViewWarningsOverlay) hide() {
	o.revealer.SetRevealChild(false)
}

func (o *roomViewWarningsOverlay) close() {
	o.hide()
	o.onClose()
}

func (o *roomViewWarningsOverlay) clear() {
	warnings := o.warnings
	for _, w := range warnings {
		o.box.Remove(w.bar)
	}
	o.warnings = nil
}

package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type errorNotification struct {
	gtki.Box //the container

	area  gtki.Box
	label gtki.Label
}

func newErrorNotification(info gtki.Box) *errorNotification {
	errorNotif := &errorNotification{Box: info}

	b := newBuilder("ErrorNotification")
	b.getItems(
		"infobar", &errorNotif.area,
		"message", &errorNotif.label,
	)

	info.Add(errorNotif.area)
	return errorNotif
}

func (n *errorNotification) renderAccountError(label string) {
	prov := providerWithCSS("box { background-color: #4a8fd9;  color: #ffffff; border-radius: 2px; }")
	updateWithStyle(n.area, prov)

	n.label.SetMarginTop(10)
	n.label.SetMarginBottom(10)
	n.label.SetText(i18n.Local(label))

	n.Box.ShowAll()
}

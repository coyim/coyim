package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type errorNotification struct {
	area  gtki.Box   `gtk-widget:"infobar"`
	label gtki.Label `gtk-widget:"message"`
}

func newErrorNotification(info gtki.Box) *errorNotification {
	view := &errorNotification{}

	b := newBuilder("ErrorNotification")
	panicOnDevError(b.bindObjects(view))

	info.Add(view.area)
	return view
}

func (n *errorNotification) ShowMessage(label string) {
	prov := providerWithCSS("box { background-color: #4a8fd9;  color: #ffffff; border-radius: 2px; }")
	updateWithStyle(n.area, prov)

	n.label.SetMarginTop(10)
	n.label.SetMarginBottom(10)
	n.label.SetText(i18n.Local(label))

	parent, _ := n.area.GetParent()
	parent.ShowAll()
}

func (n *errorNotification) Hide() {
	parent, _ := n.area.GetParent()
	parent.Hide()
}

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
	prov := providerWithStyle("box", style{
		"background-color": "#e53e3e",
		"border-radius":    "2px",
		"color":            "#ffffff",
	})

	updateWithStyle(n.area, prov)

	n.label.SetMarginTop(10)
	n.label.SetMarginBottom(10)
	n.label.SetText(i18n.Local(label))

	parent, _ := n.area.GetParent()
	if parent != nil {
		parent.ShowAll()
	}
}

func (n *errorNotification) Hide() {
	parent, _ := n.area.GetParent()
	if parent != nil {
		parent.Hide()
	}
}

func (n *errorNotification) IsVisible() bool {
	parent, _ := n.area.GetParent()

	return parent != nil && parent.IsVisible()
}

type canNotifyErrors interface {
	// clearErrors should ONLY be called from the UI thread
	clearErrors()
	// notifyOnError should ONLY be called from the UI thread
	notifyOnError(err string)
}

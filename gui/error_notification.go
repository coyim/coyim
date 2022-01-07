package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type errorNotification struct {
	area  gtki.Box   `gtk-widget:"infobar"`
	label gtki.Label `gtk-widget:"message"`
}

func newErrorNotification(wl withLog, info gtki.Box) *errorNotification {
	view := &errorNotification{}

	b := newBuilder("ErrorNotification")
	panicOnDevError(b.bindObjects(view))

	info.Add(view.area)

	prov := providerWithStyle(wl, "error notification", "box", style{
		"background-color": "#e53e3e",
		"border-radius":    "2px",
		"color":            "#ffffff",
	})

	updateWithStyle(view.area, prov)

	return view
}

// ShowMessage will display the already internationalized label
func (n *errorNotification) ShowMessage(label string) {
	n.label.SetMarginTop(10)
	n.label.SetMarginBottom(10)
	n.label.SetText(label)

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

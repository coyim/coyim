package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type infoBar struct {
	*box
	internal *gtk.InfoBar
}

func WrapInfoBarSimple(v *gtk.InfoBar) gtki.InfoBar {
	if v == nil {
		return nil
	}
	return &infoBar{WrapBoxSimple(&v.Box).(*box), v}
}

func WrapInfoBar(v *gtk.InfoBar, e error) (gtki.InfoBar, error) {
	return WrapInfoBarSimple(v), e
}

func UnwrapInfoBar(v gtki.InfoBar) *gtk.InfoBar {
	if v == nil {
		return nil
	}
	return v.(*infoBar).internal
}

func (*RealGtk) InfoBarNew() (gtki.InfoBar, error) {
	return WrapInfoBar(gtk.InfoBarNew())
}

func (v *infoBar) AddActionWidget(w gtki.Widget, responseId gtki.ResponseType) {
	v.internal.AddActionWidget(UnwrapWidget(w), gtk.ResponseType(responseId))
}

func (v *infoBar) AddButton(buttonText string, responseId gtki.ResponseType) {
	v.internal.AddButton(buttonText, gtk.ResponseType(responseId))
}

func (v *infoBar) SetDefaultResponse(responseId gtki.ResponseType) {
	v.internal.SetDefaultResponse(gtk.ResponseType(responseId))
}

func (v *infoBar) SetMessageType(messageType gtki.MessageType) {
	v.internal.SetMessageType(gtk.MessageType(messageType))
}

func (v *infoBar) GetMessageType() gtki.MessageType {
	return gtki.MessageType(v.internal.GetMessageType())
}

func (v *infoBar) GetActionArea() (gtki.Widget, error) {
	return WrapWidget(v.internal.GetActionArea())
}

func (v *infoBar) GetContentArea() (gtki.Box, error) {
	return WrapBox(v.internal.GetContentArea())
}

func (v *infoBar) GetShowCloseButton() bool {
	return v.internal.GetShowCloseButton()
}

func (v *infoBar) SetShowCloseButton(setting bool) {
	v.internal.SetShowCloseButton(setting)
}

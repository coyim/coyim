package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewConversation struct {
	view     gtki.Box        `gtk-widget:"roomConversation"`
	messages gtki.TextBuffer `gtk-widget:"messages"`
}

func newRoomViewConversation() *roomViewConversation {
	c := &roomViewConversation{}

	builder := newBuilder("MUCRoomConversation")
	panicOnDevError(builder.bindObjects(c))

	return c
}

func (v *roomViewConversation) showOccupantLeftRoom(nickname jid.Resource) {
	v.addNewMessage(i18n.Localf("%s left the room", nickname))
}

func (v *roomViewConversation) addNewMessage(text string) {
	i := v.messages.GetEndIter()
	v.messages.Insert(i, fmt.Sprintf("%s\n", text))
}

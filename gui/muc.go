package gui

import (
	"log"

	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucMockupView struct {
	gtki.Window
	entry gtki.Entry

	chat interfaces.Chat
}

func (v *mucMockupView) connectOrSendMessage(msg string) {
	log.Printf("--> %q", msg)

	if !v.chat.CheckForSupport(msg) {
		log.Println("No support to MUC")
	} else {
		log.Println("MUC is supported")
	}

	rooms, err := v.chat.QueryRooms(msg)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%s has rooms:", msg)
	for _, i := range rooms {
		log.Printf("- %s\t%s", i.Jid, i.Name)
	}

	response, err := v.chat.QueryRoomInformation(msg)
	if err != nil {
		log.Println(err)
	}

	log.Printf("RoomInfo: %#v", response)

}

func (v *mucMockupView) onSendMessage(_ glibi.Object) {
	//TODO: Why cant I use entry as gtki.Entry?
	//TODO: File a bug againt gotkadapter

	msg, err := v.entry.GetText()
	if err != nil {
		return
	}

	v.entry.SetText("")

	go v.connectOrSendMessage(msg)
}

func (u *gtkUI) addChatRoom() {
	builder := newBuilder("AddChat")
	dialog := builder.get("add-chat-dialog").(gtki.Dialog)

	dialog.SetTransientFor(u.window)
	dialog.Show()
}

func (u *gtkUI) openMUCMockup() {
	accounts := u.getAllConnectedAccounts()
	conn := accounts[0].session.Conn()
	if conn == nil {
		return
	}

	builder := newBuilder("MUCMockup")

	mockup := &mucMockupView{
		chat: conn.GetChatContext(), //TODO: hackish

		Window: builder.get("muc-window").(gtki.Window),
		entry:  builder.get("text-box").(gtki.Entry),
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_send_message": mockup.onSendMessage,
	})

	mockup.SetTransientFor(u.window)
	mockup.Show()
}

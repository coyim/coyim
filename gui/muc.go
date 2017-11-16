package gui

import (
	"log"

	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucMockupView struct {
	gtki.Window `gtk-widget:"muc-window"`
	entry       gtki.Entry `gtk-widget:"text-box"`

	chat     interfaces.Chat
	occupant *xmpp.Occupant
}

func newMockupView(account *account, occupant *xmpp.Occupant) *mucMockupView {
	conn := account.session.Conn()
	if conn == nil {
		return nil
	}

	builder := newBuilder("MUCMockup")
	mockup := &mucMockupView{
		chat:     conn.GetChatContext(),
		occupant: occupant,
	}

	err := builder.bindObjects(mockup)
	if err != nil {
		panic(err)
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_send_message": mockup.onSendMessage,
	})

	mockup.SetTitle(occupant.JID())

	return mockup
}

func (v *mucMockupView) connectOrSendMessage(msg string) {
	log.Printf("--> %q", msg)

	if v.occupant == nil {
		return
	}

	if !v.chat.CheckForSupport(v.occupant.Service) {
		log.Println("No support to MUC")
	} else {
		log.Println("MUC is supported")
	}

	rooms, err := v.chat.QueryRooms(v.occupant.Service)
	if err != nil {
		log.Println(err)
	}

	log.Printf("%s has rooms:", msg)
	for _, i := range rooms {
		log.Printf("- %s\t%s", i.Jid, i.Name)
	}

	response, err := v.chat.QueryRoomInformation(v.occupant.Room.JID())
	if err != nil {
		log.Println("Error to query room information")
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
	mockup := newMockupView(accounts[0], nil)
	mockup.SetTransientFor(u.window)
	mockup.Show()
}

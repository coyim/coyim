package gui

import (
	"errors"
	"log"

	"github.com/coyim/coyim/xmpp"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type addChatView struct {
	accountManager *accountManager

	gtki.Dialog `gtk-widget:"add-chat-dialog"`

	model   gtki.ListStore `gtk-widget:"accounts-model"`
	account gtki.ComboBox  `gtk-widget:"accounts"`
	service gtki.Entry     `gtk-widget:"service"`
	room    gtki.Entry     `gtk-widget:"room"`
	handle  gtki.Entry     `gtk-widget:"handle"`
}

func newChatView(accountManager *accountManager) gtki.Dialog {
	view := &addChatView{
		accountManager: accountManager,
	}

	builder := newBuilder("AddChat")
	err := builder.bindObjects(view)
	if err != nil {
		panic(err)
	}

	builder.ConnectSignals(map[string]interface{}{
		"join_room_handler": view.joinRoomHandler,
		"cancel_handler":    view.Destroy,
	})

	view.populateModel()

	return view
}

func (v *addChatView) populateModel() {
	accs := v.accountManager.getAllConnectedAccounts()
	for _, acc := range accs {
		iter := v.model.Append()
		v.model.SetValue(iter, 0, acc.session.GetConfig().Account)
		v.model.SetValue(iter, 1, acc.session.GetConfig().ID())
	}

	if len(accs) > 0 {
		v.account.SetActive(0)
	}
}

//TODO: This is repeated on AddAccount logic, for example.
func (v *addChatView) getAccount() (*account, bool) {
	iter, err := v.account.GetActiveIter()
	if err != nil {
		return nil, false
	}

	val, err := v.model.GetValue(iter, 1)
	if err != nil {
		return nil, false
	}

	id, err := val.GetString()
	if err != nil {
		return nil, false
	}

	return v.accountManager.getAccountByID(id)
}

func (v *addChatView) validateForm() (*account, *xmpp.Occupant, error) {
	account, ok := v.getAccount()
	if !ok {
		return nil, nil, errors.New("could not find account")
	}

	service, err := v.service.GetText()
	if err != nil {
		return nil, nil, err
	}

	room, err := v.room.GetText()
	if err != nil {
		return nil, nil, err
	}

	handle, err := v.handle.GetText()
	if err != nil {
		return nil, nil, err
	}

	//TODO: VALIDATE!

	occ := &xmpp.Occupant{
		Room: xmpp.Room{
			ID:      room,
			Service: service,
		},
		Nick: handle,
	}

	return account, occ, nil
}

func (v *addChatView) joinRoomHandler() {
	account, occupant, err := v.validateForm()
	if err != nil {
		//TODO: show error
		return
	}

	doInUIThread(func() {
		//TODO: the reference to this object should be kept
		//otherwise it will be garbage collected before we are done
		//with the window. We cant keep this goroutine blocked to avoid
		//the view from leaving the scope, because it would block the
		//glib main thread.
		//It probably works because glib.IdleAdd never releases this fn
		chatRoom := newMockupView(account, occupant)
		if parent, err := v.GetTransientFor(); err != nil {
			chatRoom.SetTransientFor(parent)
		}

		v.Destroy()
		chatRoom.Show()
	})
}

func (u *gtkUI) addChatRoom() {
	view := newChatView(u.accountManager)
	view.SetTransientFor(u.window)
	view.Show()
}

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

func (u *gtkUI) openMUCMockup() {
	accounts := u.getAllConnectedAccounts()
	mockup := newMockupView(accounts[0], nil)
	mockup.SetTransientFor(u.window)
	mockup.Show()
}

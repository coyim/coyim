package gui

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/ui"
	"github.com/coyim/coyim/xmpp"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/utils"
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

	doInUIThread(view.populateModel)

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

func (v *addChatView) validateForm() (*account, *data.Occupant, error) {
	account, ok := v.getAccount()
	if !ok {
		return nil, nil, errors.New("could not find account")
	}

	//TODO: If service is empty, should get it from account's JID
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
	if handle == "" {
		handle = xmpp.ParseJID(account.session.GetConfig().Account).LocalPart
	}

	occ := &data.Occupant{
		Room: data.Room{
			ID:      room,
			Service: service,
		},
		Handle: handle,
	}

	return account, occ, nil
}

func (v *addChatView) joinRoomHandler() {
	account, occupant, err := v.validateForm()
	if err != nil {
		//TODO: show error
		return
	}

	conn := account.session.Conn()
	if conn == nil {
		//TODO: show error
		return
	}

	//TODO: This should notify the user about what is happening (bc it blocks)
	//and also notify when a failure occurs.
	chat := conn.GetChatContext()
	if !chat.CheckForSupport(occupant.Service) {
		//TODO: show error
		log.Println("No support to MUC")
		return
	}

	chatRoom := newChatRoomView(chat, occupant)
	if parent, err := v.GetTransientFor(); err == nil {
		chatRoom.SetTransientFor(parent)
	}
	v.Destroy()

	account.session.Subscribe(chatRoom.eventsChan)
	chatRoom.openWindow()
}

func (u *gtkUI) joinChatRoom() {
	//pass message and presence channels
	view := newChatView(u.accountManager)
	view.SetTransientFor(u.window)
	view.Show()
}

type roomOccupant struct {
	Role        string
	Affiliation string
}

type chatRoomView struct {
	gtki.Window `gtk-widget:"muc-window"`
	subject     gtki.Label `gtk-widget:"subject"`
	entry       gtki.Entry `gtk-widget:"text-box"`

	historyMutex  sync.Mutex
	historyBuffer gtki.TextBuffer     `gtk-widget:"chat-buffer"`
	historyScroll gtki.ScrolledWindow `gtk-widget:"chat-box"`

	occupantsList struct {
		sync.Mutex

		dirty bool
		m     map[string]*roomOccupant
	}
	occupantsView  gtki.TreeView  `gtk-widget:"occupants-view"`
	occupantsModel gtki.ListStore `gtk-widget:"occupants"`

	eventsChan chan interface{}
	chat       interfaces.Chat
	occupant   *data.Occupant

	receivedSelfPresence bool
}

func newChatRoomView(chat interfaces.Chat, occupant *data.Occupant) *chatRoomView {
	builder := newBuilder("ChatRoom")
	v := &chatRoomView{
		chat:     chat,
		occupant: occupant,

		//TODO: This could go somewhere else (account maybe?)
		eventsChan: make(chan interface{}),
	}

	v.occupantsList.m = make(map[string]*roomOccupant, 5)

	err := builder.bindObjects(v)
	if err != nil {
		panic(err)
	}

	builder.ConnectSignals(map[string]interface{}{
		"send_message_handler":             v.onSendMessage,
		"scroll_history_to_bottom_handler": v.scrollHistoryToBottom,

		//TODO: A closed window will leave the room
		//Probably not what we want for the final version
		"leave_room_handler": v.leaveRoom,
	})

	v.SetTitle(occupant.Room.JID())

	return v
}

func (v *chatRoomView) showDebugInfo() {
	//TODO Remove this. It is only for debugging
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

	log.Printf("%s has rooms:", v.occupant.Service)
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

func (v *chatRoomView) openWindow() {
	//TODO: show error
	go v.chat.EnterRoom(v.occupant)

	go v.watchEvents(v.eventsChan)

	v.Show()
}

func (v *chatRoomView) leaveRoom() {
	v.chat.LeaveRoom(v.occupant)
	close(v.eventsChan)
	v.eventsChan = nil
}

func (v *chatRoomView) sameRoom(from string) bool {
	return xmpp.ParseJID(from).Bare() == v.occupant.Room.JID()
}

func (v *chatRoomView) watchEvents(evs <-chan interface{}) {
	for {
		v.redrawOccupantsList()

		ev, ok := <-evs
		if !ok {
			return
		}

		//TODO: Disable controls when the session disconnects

		switch e := ev.(type) {
		case events.ChatPresence:
			if !v.sameRoom(e.ClientPresence.From) {
				log.Println("muc: presence not for this room. %#v", e.ClientPresence)
				continue
			}

			v.updatePresence(e.ClientPresence)
		case events.ChatMessage:
			if !v.sameRoom(e.ClientMessage.From) {
				continue
			}

			//TODO: should check if body is not present, and not if it is empty
			//TODO: check if thread is also not present
			if e.ClientMessage.Subject != nil && e.ClientMessage.Body == "" {
				v.displaySubjectChange(*e.ClientMessage.Subject)
				v.notifySubjectChange(e.ClientMessage.From, *e.ClientMessage.Subject)
				continue
			}

			v.displayReceivedMessage(&e)
		default:
			//Ignore
			log.Printf("chat view got event: %#v", e)
		}
	}
}

func (v *chatRoomView) updatePresence(presence *data.ClientPresence) {
	v.occupantsList.Lock()
	defer v.occupantsList.Unlock()

	v.occupantsList.dirty = true

	if isSelfPresence(presence) {
		v.receivedSelfPresence = true
	}

	if presence.Type == "unavailable" {
		delete(v.occupantsList.m, presence.From)
		v.notifyUserLeftRoom(presence)
	} else {
		v.occupantsList.m[presence.From] = &roomOccupant{
			Role:        presence.Chat.Item.Role,
			Affiliation: presence.Chat.Item.Affiliation,
		}
		v.notifyUserEnteredRoom(presence)
	}
}

func (v *chatRoomView) notifyUserLeftRoom(presence *data.ClientPresence) {
	if !v.receivedSelfPresence {
		return
	}
	message := fmt.Sprintf("%v left the room", utils.ResourceFromJid(presence.From))
	v.notifyStatusChange(message)
}

func (v *chatRoomView) notifyUserEnteredRoom(presence *data.ClientPresence) {
	if !v.receivedSelfPresence {
		return
	}
	message := fmt.Sprintf("%v entered the room", utils.ResourceFromJid(presence.From))
	v.notifyStatusChange(message)
}

func isSelfPresence(presence *data.ClientPresence) bool {
	return presence.Chat.Status.Code == 110
}

func (v *chatRoomView) notifyStatusChange(message string) {
	doInUIThread(func() {
		v.insertNewLine()
		insertTimestamp(v.historyBuffer, time.Now())
		insertAtEnd(v.historyBuffer, message)
	})
}

func (v *chatRoomView) redrawOccupantsList() {
	if !v.occupantsList.dirty {
		return
	}

	v.occupantsList.Lock()
	defer v.occupantsList.Unlock()
	v.occupantsList.dirty = false

	doInUIThread(func() {
		v.occupantsView.SetModel(nil)
		v.occupantsModel.Clear()

		for jid, occupant := range v.occupantsList.m {
			iter := v.occupantsModel.Append()
			v.occupantsModel.SetValue(iter, 0, xmpp.ParseJID(jid).ResourcePart)
			v.occupantsModel.SetValue(iter, 1, occupant.Role)
			v.occupantsModel.SetValue(iter, 2, occupant.Affiliation)
		}

		v.occupantsView.SetModel(v.occupantsModel)
	})
}

func (v *chatRoomView) displaySubjectChange(subject string) {
	v.subject.SetVisible(true)
	v.subject.SetText(subject)
}

func (v *chatRoomView) notifySubjectChange(from, subject string) {
	from = utils.ResourceFromJid(from)
	message := fmt.Sprintf("%s has set the topic to \"%s\"", from, subject)
	v.notifyStatusChange(message)
}

func (v *chatRoomView) displayReceivedMessage(message *events.ChatMessage) {
	//TODO: maybe notify?
	doInUIThread(func() {
		v.appendToHistory(message)
	})
}

func (v *chatRoomView) appendToHistory(message *events.ChatMessage) {
	v.historyMutex.Lock()
	defer v.historyMutex.Unlock()

	v.insertNewLine()

	sent := sentMessage{
		//TODO: Why both?
		message:         message.Body,
		strippedMessage: ui.StripSomeHTML([]byte(message.Body)),

		from:      utils.ResourceFromJid(message.From),
		to:        message.To,
		timestamp: message.When,
	}

	//TODO: use attention?
	entries, _ := sent.Tagged()

	insertTimestamp(v.historyBuffer, message.When)
	for _, e := range entries {
		insertEntry(v.historyBuffer, e)
	}

	v.scrollHistoryToBottom()
}

func (v *chatRoomView) insertNewLine() {
	start := v.historyBuffer.GetCharCount()
	if start != 0 {
		insertAtEnd(v.historyBuffer, "\n")
	}
}

func (v *chatRoomView) scrollHistoryToBottom() {
	scrollToBottom(v.historyScroll)
}

func (v *chatRoomView) onSendMessage(_ glibi.Object) {
	//TODO: Why cant I use entry as gtki.Entry?
	//TODO: File a bug againt gotkadapter

	msg, err := v.entry.GetText()
	if err != nil {
		return
	}

	v.entry.SetText("")

	//TODO: error?
	go v.chat.SendChatMessage(msg, &v.occupant.Room)
}

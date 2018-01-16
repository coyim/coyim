package gui

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/ui"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type addChatView struct {
	accountManager *accountManager
	chatManager    *chatManager
	errorBox       *errorNotification

	gtki.Dialog `gtk-widget:"add-chat-dialog"`

	notification gtki.Box      `gtk-widget:"notification-area"`
	form         gtki.Grid     `gtk-widget:"form"`
	account      gtki.ComboBox `gtk-widget:"accounts"`
	service      gtki.Entry    `gtk-widget:"service"`
	room         gtki.Entry    `gtk-widget:"room"`
	handle       gtki.Entry    `gtk-widget:"handle"`

	model gtki.ListStore `gtk-widget:"accounts-model"`
}

func newChatView(accountManager *accountManager, chatManager *chatManager) *addChatView {
	view := &addChatView{
		accountManager: accountManager,
		chatManager:    chatManager,
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

	view.errorBox = newErrorNotification(view.notification)
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
func (v *addChatView) getAccount() (string, string, error) {
	iter, err := v.account.GetActiveIter()
	if err != nil {
		return "", "", err
	}

	val, err := v.model.GetValue(iter, 0)
	if err != nil {
		return "", "", err
	}

	bareJID, err := val.GetString()
	if err != nil {
		return "", "", err
	}

	val, err = v.model.GetValue(iter, 1)
	if err != nil {
		return "", "", err
	}

	id, err := val.GetString()
	if err != nil {
		return "", "", err
	}

	return id, bareJID, nil
}

func (v *addChatView) setActiveAccount(accIndex int) {
	doInUIThread(func() {
		v.account.SetActive(accIndex)
	})
}

func (v *addChatView) validateForm() (string, *data.Occupant, error) {
	accountID, bareJID, err := v.getAccount()
	if err != nil {
		panic(err)
	}

	service, err := v.service.GetText()
	if err != nil {
		panic(err)
	}

	room, err := v.room.GetText()
	if err != nil {
		panic(err)
	}

	handle, err := v.handle.GetText()
	if err != nil {
		panic(err)
	}

	//TODO: If service is empty, should get it from account's JID?

	//Validate
	if handle == "" {
		j := jid.Parse(bareJID)
		if jj, ok := j.(jid.WithLocal); ok {
			handle = string(jj.Local())
		}
	}

	occ := &data.Occupant{
		Room: data.Room{
			ID:      room,
			Service: service,
		},
		Handle: handle,
	}

	return accountID, occ, nil
}

func (v *addChatView) waitForSelfPresence(chat interfaces.Chat, occupant *data.Occupant) ([]*roomOccupant, error) {
	var ret []*roomOccupant

	//TODO: this should timeout
	//TODO: this should have a cancelation mechanism
	for ev := range chat.Events() {
		switch e := ev.(type) {
		case events.ChatPresence:
			presence := e.ClientPresence
			if jid.NR(presence.From).String() != occupant.Room.JID() {
				continue
			}

			if presence.Type == "error" {
				//TODO: Return error constants?
				return ret, fmt.Errorf("Error %s: %s", presence.Error.Code, presence.Error.Text)
			}

			if presence.Chat == nil {
				continue //TODO: this is a broken server
			}

			ret = append(ret, &roomOccupant{
				OccupantJID: presence.From,
				Role:        presence.Chat.Item.Role,
				Affiliation: presence.Chat.Item.Affiliation,
			})

			if presence.Chat.Status.Code == 110 {
				return ret, nil
			}
		}
	}

	return ret, errors.New("Did not receive a presence response")
}

//TODO: This could all go to the interfaces.Chat
func (v *addChatView) enterRoom(chat interfaces.Chat, occupant *data.Occupant) ([]*roomOccupant, error) {
	if !chat.CheckForSupport(occupant.Service) {
		return nil, errors.New("The service does not support chat.")
	}

	err := chat.EnterRoom(occupant)
	if err != nil {
		return nil, err
	}

	return v.waitForSelfPresence(chat, occupant)
}

//openRoomDialog blocks to do networking and should be called in a goroutine
func (v *addChatView) openRoomDialog(chat interfaces.Chat, occupant *data.Occupant) {
	occupantsInRoom, err := v.enterRoom(chat, occupant)
	if err != nil {
		doInUIThread(func() {
			v.form.Show()
			v.errorBox.ShowMessage(i18n.Local(err.Error()))
		})
		return
	}

	doInUIThread(func() {
		defer v.Destroy()

		chatRoom := newChatRoomView(chat, occupant)
		chatRoom.setOccupantList(occupantsInRoom)
		if parent, err := v.GetTransientFor(); err == nil {
			chatRoom.SetTransientFor(parent)
		}
		chatRoom.openWindow()
	})
}

func (v *addChatView) getChatAndOccupantFromForm() (interfaces.Chat, *data.Occupant, error) {
	accountID, occupant, err := v.validateForm()
	if err != nil {
		return nil, nil, err
	}

	chat, err := v.chatManager.getChatContextForAccount(accountID)
	if err != nil {
		return nil, nil, err
	}

	return chat, occupant, nil
}

func (v *addChatView) joinRoomHandler() {
	v.errorBox.Hide()

	chat, occupant, err := v.getChatAndOccupantFromForm()
	if err != nil {
		v.form.Show()
		v.errorBox.ShowMessage(err.Error())
		return
	}

	v.form.Hide()
	v.errorBox.ShowMessage(i18n.Localf("Joining #%s", occupant.Room.ID))
	go v.openRoomDialog(chat, occupant)
}

type roomConfigView struct {
	dialog gtki.Dialog `gtk-widget:"dialog"`
	grid   gtki.Grid   `gtk-widget:"grid"`

	formFields []formField
	done       chan<- interface{}
}

func newRoomConfigDialog(done chan<- interface{}, fields []formField) *roomConfigView {
	view := &roomConfigView{
		formFields: fields,
		done:       done,
	}

	builder := newBuilder("ConfigureRoom")
	err := builder.bindObjects(view)
	if err != nil {
		panic(err)
	}

	builder.ConnectSignals(map[string]interface{}{
		"on_cancel_signal": view.close,
		"on_save_signal":   view.updateFormWithValuesFromFormFields,
	})

	view.attachFormFields()

	return view
}

func (v *roomConfigView) close() {
	v.dialog.Destroy()
	v.done <- true
}

func (v *roomConfigView) updateFormWithValuesFromFormFields() {
	//Find the fields we need to copy from the form to the account
	for _, field := range v.formFields {
		switch ff := field.field.(type) {
		case *data.TextFormField:
			w := field.widget.(gtki.Entry)
			ff.Result, _ = w.GetText()
		case *data.BooleanFormField:
			w := field.widget.(gtki.CheckButton)
			ff.Result = w.GetActive()
		case *data.SelectionFormField:
			w := field.widget.(gtki.ComboBoxText)
			ff.Result = w.GetActive()
		default:
			log.Printf("We need to implement %#v", ff)
		}
	}

	v.close()
}

func (v *roomConfigView) attachFormFields() {
	for i, field := range v.formFields {
		v.grid.Attach(field.label, 0, i+1, 1, 1)
		v.grid.Attach(field.widget, 1, i+1, 1, 1)
	}
}

//This will be called from a goroutine because otherwise it would block the gtk event loop
//Thats why we need to do everything GTK-related inUIThread
func (v *chatRoomView) renderForm(title, instructions string, fields []interface{}) error {
	done := make(chan interface{})

	doInUIThread(func() {
		formFields := buildWidgetsForFields(fields)
		dialog := newRoomConfigDialog(done, formFields)

		if parent, err := v.GetTransientFor(); err == nil {
			dialog.dialog.SetTransientFor(parent)
		}
		dialog.dialog.ShowAll()
	})

	<-done
	close(done)
	return nil
}

func (v *chatRoomView) showRoomConfigDialog() {
	//Run in a goroutine to not block the GTK event loop
	//TODO: Display error
	go v.chat.RoomConfigForm(&v.occupant.Room, v.renderForm)
}

func (u *gtkUI) joinChatRoom() {
	//pass message and presence channels
	view := newChatView(u.accountManager, u.chatManager)
	view.SetTransientFor(u.window)
	view.Show()
}

type roomOccupant struct {
	OccupantJID string
	Role        string
	Affiliation string
}

type chatRoomView struct {
	gtki.Window `gtk-widget:"muc-window"`
	subject     gtki.Label `gtk-widget:"subject"`
	entry       gtki.Entry `gtk-widget:"text-box"`

	historyMutex  sync.Mutex
	menuBox       gtki.Box            `gtk-widget:"menu-box"`
	historyBuffer gtki.TextBuffer     `gtk-widget:"chat-buffer"`
	historyScroll gtki.ScrolledWindow `gtk-widget:"chat-box"`

	occupantsList struct {
		sync.Mutex

		dirty bool
		m     map[string]*roomOccupant
	}
	occupantsView  gtki.TreeView  `gtk-widget:"occupants-view"`
	occupantsModel gtki.ListStore `gtk-widget:"occupants"`

	chat     interfaces.Chat
	occupant *data.Occupant
}

func newChatRoomView(chat interfaces.Chat, occupant *data.Occupant) *chatRoomView {
	builder := newBuilder("ChatRoom")
	v := &chatRoomView{
		chat:     chat,
		occupant: occupant,
	}

	v.occupantsList.m = make(map[string]*roomOccupant, 5)

	err := builder.bindObjects(v)
	if err != nil {
		panic(err)
	}

	doInUIThread(func() {
		prov := providerWithCSS("box { border-top: 1px solid #d3d3d3; }")
		updateWithStyle(v.menuBox, prov)
	})

	builder.ConnectSignals(map[string]interface{}{
		"send_message_handler":             v.onSendMessage,
		"scroll_history_to_bottom_handler": v.scrollHistoryToBottom,
		"on_change_room_config":            v.showRoomConfigDialog,

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
	go v.watchEvents(v.chat.Events())
	v.Show()
}

func (v *chatRoomView) authenticationError() {
	//TODO: Go to "join chat room" dialog and show error message
	doInUIThread(v.Destroy)
}

func (v *chatRoomView) leaveRoom() {
	v.chat.LeaveRoom(v.occupant)
}

func (v *chatRoomView) sameRoom(from string) bool {
	return jid.NR(from).String() == v.occupant.Room.JID()
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
				log.Printf("muc: presence not for this room. %#v", e.ClientPresence)
				continue
			}

			//See: XEP-0045, section "7.2.6 Password-Protected Rooms"
			if e.ClientPresence.Type == "error" && e.ClientPresence.Chat == nil {
				presenceError := e.ClientPresence.Error
				errorCondition := e.ClientPresence.Error.Condition.XMLName
				if presenceError.Type == "auth" && errorCondition.Local == "not-authorized" && errorCondition.Space == "urn:ietf:params:xml:ns:xmpp-stanzas" {
					v.authenticationError()
					continue
				}
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

	if presence.Type == "unavailable" {
		delete(v.occupantsList.m, presence.From)
		v.notifyUserLeftRoom(presence)
	} else {
		v.occupantsList.m[presence.From] = &roomOccupant{
			OccupantJID: presence.From,
			Role:        presence.Chat.Item.Role,
			Affiliation: presence.Chat.Item.Affiliation,
		}
		v.notifyUserEnteredRoom(presence)
	}
}

func resourceFromJid(j string) string {
	return string(jid.Parse(j).PotentialResource())
}

func (v *chatRoomView) notifyUserLeftRoom(presence *data.ClientPresence) {
	message := fmt.Sprintf("%v left the room", resourceFromJid(presence.From))
	v.notifyStatusChange(message)
}

func (v *chatRoomView) notifyUserEnteredRoom(presence *data.ClientPresence) {
	message := fmt.Sprintf("%v entered the room", resourceFromJid(presence.From))
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

func (v *chatRoomView) setOccupantList(occupants []*roomOccupant) {
	v.occupantsList.Lock()
	defer v.occupantsList.Unlock()
	v.occupantsList.dirty = true

	for _, occ := range occupants {
		v.occupantsList.m[occ.OccupantJID] = occ
	}
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

		for j, occupant := range v.occupantsList.m {
			iter := v.occupantsModel.Append()
			v.occupantsModel.SetValue(iter, 0, string(jid.Parse(j).PotentialResource()))
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
	message := fmt.Sprintf("%s has set the topic to \"%s\"", resourceFromJid(from), subject)
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

		// TODO: this is clearly completely incorrect
		from:      string(message.From.Resource()),
		to:        jid.NR(message.To),
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

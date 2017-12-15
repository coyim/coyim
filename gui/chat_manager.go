package gui

import (
	"errors"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/interfaces"
)

type chatManager struct {
	//TODO: Shall we hide accountManager calls inside chatManager?
	*accountManager
}

func newChatManager(m *accountManager) *chatManager {
	return &chatManager{
		accountManager: m,
	}
}

func (m *chatManager) getChatContextForAccount(accountID string) (interfaces.Chat, error) {
	account, ok := m.accountManager.getAccountByID(accountID)
	if !ok {
		return nil, errors.New(i18n.Local("The selected account could not be found."))
	}

	conn := account.session.Conn()
	if conn == nil {
		return nil, errors.New(i18n.Local("The selected account is not connected."))
	}

	chat := conn.GetChatContext()
	account.session.Subscribe(chat.Events())
	return chat, nil
}

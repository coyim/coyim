package session

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/otrclient"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/otr3"
)

func (s *session) onOtrEventHandlerCreate(peer jid.Any, eh *otrclient.EventHandler, nots chan string, dels chan int) {
	go s.listenToOtrNotifications(nots, peer)
	go s.listenToOtrDelayedMessageDelivery(dels, peer)
}

func (s *session) listenToOtrNotifications(c <-chan string, peer jid.Any) {
	for notification := range c {
		s.notify(peer, notification)
	}
}

func (s *session) listenToOtrDelayedMessageDelivery(c <-chan int, peer jid.Any) {
	for t := range c {
		s.publishEvent(events.DelayedMessageSent{
			Peer:   peer,
			Tracer: t,
		})
	}
}

func (s *session) newOTRKeys(peer jid.WithResource, conversation otrclient.Conversation) {
	s.publishPeerEvent(events.OTRNewKeys, peer)
}

func (s *session) renewedOTRKeys(peer jid.WithResource, conversation otrclient.Conversation) {
	s.publishPeerEvent(events.OTRRenewedKeys, peer)
}

func (s *session) otrEnded(peer jid.WithResource) {
	s.publishPeerEvent(events.OTREnded, peer)
}

// newConversation will create a new OTR conversation with the given peer
// TODO: Why starting a conversation requires being able to translate a message?
// This also assumes it's useful to send friendly message to another person in
// the same language configured on your app.
func (s *session) newConversation(peer jid.Any) *otr3.Conversation {
	conversation := &otr3.Conversation{}
	conversation.SetOurKeys(s.privateKeys)
	conversation.SetFriendlyQueryMessage(i18n.Local("Your peer has requested a private conversation with you, but your client doesn't seem to support the OTR protocol."))

	instanceTag := conversation.InitializeInstanceTag(s.GetConfig().InstanceTag)

	if s.GetConfig().InstanceTag != instanceTag {
		s.cmdManager.ExecuteCmd(otrclient.SaveInstanceTagCmd{
			Account:     s.GetConfig(),
			InstanceTag: instanceTag,
		})
	}

	s.GetConfig().SetOTRPoliciesFor(peer.NoResource().String(), conversation)

	return conversation
}

// ManuallyEndEncryptedChat allows a user to end the encrypted chat from this side
func (s *session) ManuallyEndEncryptedChat(peer jid.Any) error {
	c, ok := s.ConversationManager().GetConversationWith(peer)
	if !ok {
		return fmt.Errorf("couldn't find conversation with %s", peer)
	}

	defer c.EventHandler().ConsumeSecurityChange()
	return c.EndEncryptedChat()
}

func (s *session) terminateConversations() {
	s.convManager.TerminateAll()
}

// StartSMP begins the SMP interactions for a conversation
func (s *session) StartSMP(peer jid.WithResource, question, answer string) {
	conv, ok := s.convManager.GetConversationWith(peer)
	if !ok {
		s.log.Error("tried to start SMP when a conversation does not exist")
		return
	}
	if err := conv.StartAuthenticate(question, []byte(answer)); err != nil {
		s.log.WithError(err).Error("cannot start SMP")
	}
}

// FinishSMP takes a user's SMP answer and finishes the protocol
func (s *session) FinishSMP(peer jid.WithResource, answer string) {
	conv, ok := s.convManager.GetConversationWith(peer)
	if !ok {
		s.log.Error("tried to finish SMP when a conversation does not exist")
		return
	}
	if err := conv.ProvideAuthenticationSecret([]byte(answer)); err != nil {
		s.log.WithError(err).Error("cannot provide an authentication secret for SMP")
	}
}

// AbortSMP will abort the current SMP interaction for a conversation
func (s *session) AbortSMP(peer jid.WithResource) {
	conv, ok := s.convManager.GetConversationWith(peer)
	if !ok {
		s.log.Error("tried to abort SMP when a conversation does not exist")
		return
	}
	if err := conv.AbortAuthentication(); err != nil {
		s.log.WithError(err).Error("cannot abort SMP")
	}
}

func (s *session) CreateSymmetricKeyFor(peer jid.Any) []byte {
	conv, ok := s.convManager.GetConversationWith(peer)
	if !ok {
		return nil
	}

	key, err := conv.CreateExtraSymmetricKey()
	if err != nil {
		return nil
	}

	return key
}

func (s *session) GetAndWipeSymmetricKeyFor(peer jid.Any) []byte {
	conv, ok := s.convManager.GetConversationWith(peer)
	if !ok {
		return nil
	}

	_, _, key := conv.GetAndWipeLastExtraKey()

	return key
}

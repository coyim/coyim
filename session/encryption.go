package session

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/events"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

var knownAlgorithms = map[string]string{
	"jabber:x:encrypted":             "Legacy OpenPGP",
	"urn:xmpp:openpgp:0":             "OpenPGP for XMPP",
	"eu.siacs.conversations.axolotl": "OMEMO",
}

const otrEncryptionNamespace = "urn:xmpp:otr:0"

func (s *session) notify(peer jid.Any, notification string) {
	s.publishEvent(events.Notification{
		Peer:         peer,
		Notification: notification,
	})
}

func (s *session) processEncryption(peer jid.Any, e *data.Encryption) {
	if e.Namespace == otrEncryptionNamespace {
		// We have message marked as OTR, everything good
		s.connectionLogger.Debug("got message marked with OTR encryption tag (XEP-0380)")
		return
	}

	name, ok := knownAlgorithms[e.Namespace]
	if !ok {
		name = e.Name
	}

	s.connectionLogger.WithFields(log.Fields{
		"namespace": e.Namespace,
		"name":      name,
	}).Info("got message marked with unknown encryption tag (XEP-0380)")

	s.notify(peer, i18n.Localf("We received a message encrypted with %s - sadly CoyIM does not support this algorithm. Please let your contact know to encrypt using OTR, nothing else, to communicate with you.", name))
}

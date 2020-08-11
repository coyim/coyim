package session

import (
	"encoding/xml"

	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/xmpp"
	"github.com/coyim/coyim/xmpp/data"
)

func discoIQ(s access.Session, iq *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	s.Log().Info("IQ: http://jabber.org/protocol/disco#info query")

	query := &data.DiscoveryInfoQuery{}
	err := xml.Unmarshal(iq.Query, query)

	if err != nil {
		s.(*session).log.WithError(err).Error("error on parsing disco#info query")
		return nil, "error", false
	}

	return xmpp.DiscoveryReply(s.GetConfig().Account, query.Node), "", false
}

func discoItemsIQ(s access.Session, iq *data.ClientIQ) (ret interface{}, iqtype string, ignore bool) {
	s.Log().Info("IQ: http://jabber.org/protocol/disco#items query")

	query := &data.DiscoveryItemsQuery{}
	err := xml.Unmarshal(iq.Query, query)

	if err != nil {
		s.(*session).log.WithError(err).Error("error on parsing disco#items query")
		return nil, "error", false
	}

	// If someone asks for the rooms, we always return an empty list
	if query.Node == "http://jabber.org/protocol/muc#rooms" {
		return data.DiscoveryItemsQuery{
			Node:           query.Node,
			DiscoveryItems: []data.DiscoveryItem{},
		}, "", false
	}

	return data.ErrorReply{
		Type:  "cancel",
		Error: data.ErrorServiceUnavailable{},
	}, "", false

}

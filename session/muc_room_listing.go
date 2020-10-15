package session

import (
	"sync"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

func hasIdentity(idents []data.DiscoveryIdentity, category, tp string) (name string, ok bool) {
	for _, id := range idents {
		if id.Category == category && id.Type == tp {
			return id.Name, true
		}
	}
	return "", false
}

func stringArrayContains(r []string, a string) bool {
	for _, f := range r {
		if f == a {
			return true
		}
	}

	return false
}

func hasFeatures(features []string, expected ...string) bool {
	for _, exp := range expected {
		if !stringArrayContains(features, exp) {
			return false
		}
	}
	return true
}

func (s *session) getRoomsInService(service jid.Any, name string, results chan<- *muc.RoomListing, resultsServices chan<- *muc.ServiceListing, allRooms *sync.WaitGroup) {
	defer allRooms.Done()

	s.log.WithField("service", service).Debug("getRoomsInService()")
	idents, features, ok := s.Conn().DiscoveryFeaturesAndIdentities(service.String())
	if !ok {
		return
	}

	identName, hasIdent := hasIdentity(idents, "conference", "text")
	if !hasIdent {
		return
	}

	if !hasFeatures(features, "http://jabber.org/protocol/disco#items", "http://jabber.org/protocol/muc") {
		return
	}

	sl := muc.NewServiceListing(service, identName)
	resultsServices <- sl

	items, err := s.Conn().QueryServiceItems(service.String())
	if err != nil {
		s.log.WithError(err).Debug("getRoomsInService() had error")
		return
	}

	s.discoverRoomListingInformation(items.DiscoveryItems, sl, results)
}

func (s *session) discoverRoomListingInformation(items []data.DiscoveryItem, sl *muc.ServiceListing, results chan<- *muc.RoomListing) {
	for _, i := range items {
		rl := muc.NewRoomListing()
		rl.Service = sl.Jid
		rl.ServiceName = sl.Name
		rl.Jid = jid.Parse(i.Jid).(jid.Bare)
		rl.Name = i.Name

		results <- rl

		go s.muc.findOutMoreInformationAboutRoom(rl)
	}
}

func (s *session) getRoomsAsync(server jid.Domain, results chan<- *muc.RoomListing, resultsServices chan<- *muc.ServiceListing, errorResult chan<- error) {
	s.log.WithField("server", server).Debug("getRoomsAsync()")
	ditems, err := s.conn.QueryServiceItems(server.String())
	if err != nil {
		errorResult <- err
		return
	}

	allServices := sync.WaitGroup{}
	allServices.Add(len(ditems.DiscoveryItems))
	for _, di := range ditems.DiscoveryItems {
		go s.getRoomsInService(jid.Parse(di.Jid), di.Name, results, resultsServices, &allServices)
	}
	allServices.Wait()

	// This signals we are done
	results <- nil
}

func (s *session) getRoomsAsyncCustomService(service string, results chan<- *muc.RoomListing, resultsServices chan<- *muc.ServiceListing, errorResult chan<- error) {
	s.log.WithField("service", service).Debug("getRoomsAsyncCustomService()")

	allServices := sync.WaitGroup{}
	allServices.Add(1)
	go s.getRoomsInService(jid.Parse(service), "", results, resultsServices, &allServices)
	allServices.Wait()

	// This signals we are done
	results <- nil
}

func (s *session) GetRooms(server jid.Domain, customService string) (<-chan *muc.RoomListing, <-chan *muc.ServiceListing, <-chan error) {
	s.log.WithField("server", server).Debug("GetRooms()")
	// 20 is a fairly arbitrary work load number that will restrict how many
	// messages we send at the same time, so we don't overwhelm or DoS the servers.
	result := make(chan *muc.RoomListing, 20)
	resultServices := make(chan *muc.ServiceListing, 20)
	errorResult := make(chan error, 1)

	if customService == "" {
		go s.getRoomsAsync(server, result, resultServices, errorResult)
	} else {
		go s.getRoomsAsyncCustomService(customService, result, resultServices, errorResult)
	}

	return result, resultServices, errorResult
}

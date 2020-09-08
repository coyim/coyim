package session

import (
	"fmt"
	"strconv"
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

func extractFormData(fields []data.FormFieldX) map[string][]string {
	result := make(map[string][]string)
	for _, field := range fields {
		result[field.Var] = field.Values
	}

	return result
}

func (s *session) findOutMoreInformationAboutRoom(rl *muc.RoomListing) {
	diq, e := s.Conn().QueryServiceInformation(rl.Jid.String())
	if e != nil {
		s.log.WithError(e).Debug("findOutMoreInformationAboutRoom() had error")
		return
	}

	for _, feat := range diq.Features {
		switch feat.Var {
		case "http://jabber.org/protocol/muc":
			// Supports MUC - probably not useful for us
		case "http://jabber.org/protocol/muc#stable_id":
			// This means the server will use the same id in groupchat messages
		case "http://jabber.org/protocol/muc#self-ping-optimization":
			// This means the chat room supports XEP-0410, that allows
			// users to see if they're still connected to a chat room.
		case "http://jabber.org/protocol/disco#info":
			// Ignore
		case "http://jabber.org/protocol/disco#items":
			// Ignore
		case "urn:xmpp:mam:0":
			// Ignore
		case "urn:xmpp:mam:1":
			// Ignore
		case "urn:xmpp:mam:2":
			// Ignore
		case "urn:xmpp:mam:tmp":
			// Ignore
		case "urn:xmpp:mucsub:0":
			// Ignore
		case "urn:xmpp:sid:0":
			// Ignore
		case "vcard-temp":
			// Ignore
		case "http://jabber.org/protocol/muc#request":
			rl.SupportsVoiceRequests = true
		case "jabber:iq:register":
			rl.AllowsRegistration = true
		case "muc_semianonymous":
			rl.Anonymity = "semi"
		case "muc_nonanonymous":
			rl.Anonymity = "no"
		case "muc_persistent":
			rl.Persistent = true
		case "muc_temporary":
			rl.Persistent = false
		case "muc_unmoderated":
			rl.Moderated = false
		case "muc_moderated":
			rl.Moderated = true
		case "muc_open":
			rl.Open = true
		case "muc_membersonly":
			rl.Open = false
		case "muc_passwordprotected":
			rl.PasswordProtected = true
		case "muc_unsecured":
			rl.PasswordProtected = false
		case "muc_public":
			rl.Public = true
		case "muc_hidden":
			rl.Public = false
		default:
			fmt.Printf("UNKNOWN FEATURE: %s\n", feat.Var)
		}
	}

	for _, form := range diq.Forms {
		formData := extractFormData(form.Fields)

		if form.Type == "result" && len(formData["FORM_TYPE"]) > 0 && formData["FORM_TYPE"][0] == "http://jabber.org/protocol/muc#roominfo" {
			for k, val := range formData {
				switch k {
				case "FORM_TYPE":
					// Ignore, we already checked
				case "muc#roominfo_lang":
					if len(val) > 0 {
						rl.Language = val[0]
					}
				case "muc#roominfo_changesubject":
					if len(val) > 0 {
						rl.OccupantsCanChangeSubject = val[0] == "1"
					}
				case "muc#roomconfig_enablelogging":
					if len(val) > 0 {
						rl.Logged = val[0] == "1"
					}
				case "muc#roomconfig_roomname":
					// Room name - we already have this
				case "muc#roominfo_description":
					if len(val) > 0 {
						rl.Description = val[0]
					}
				case "muc#roominfo_occupants":
					if len(val) > 0 {
						res, e := strconv.Atoi(val[0])
						if e != nil {
							rl.Occupants = res
						}
					}
				case "{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites":
					if len(val) > 0 {
						rl.MembersCanInvite = val[0] == "1"
					}
				case "muc#roomconfig_allowinvites":
					if len(val) > 0 {
						rl.OccupantsCanInvite = val[0] == "1"
					}
				case "muc#roomconfig_allowpm":
					if len(val) > 0 {
						rl.AllowPrivateMessages = val[0]
					}
				case "muc#roominfo_contactjid":
					if len(val) > 0 {
						rl.ContactJid = val[0]
					}
				default:
					fmt.Printf("UNKNOWN FORM VAR: %s\n", k)
				}
			}
		}
	}

	rl.Updated()
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

	sl := muc.NewServiceListing()
	sl.Jid = service
	sl.Name = identName
	resultsServices <- sl

	items, err := s.Conn().QueryServiceItems(service.String())
	if err != nil {
		s.log.WithError(err).Debug("getRoomsInService() had error")
		return
	}

	for _, i := range items.DiscoveryItems {
		rl := muc.NewRoomListing()
		rl.Service = service
		rl.ServiceName = identName
		rl.Jid = jid.Parse(i.Jid).(jid.Bare)
		rl.Name = i.Name

		results <- rl

		go s.findOutMoreInformationAboutRoom(rl)
	}
}

func (s *session) getRoomsAsync(server jid.Domain, results chan<- *muc.RoomListing, resultsServices chan<- *muc.ServiceListing, errorResult chan<- error) {
	s.log.WithField("server", server).Debug("getRoomsAsync()")
	ditems, err := s.conn.QueryServiceItems(server.String())
	if err != nil {
		errorResult <- err
		return
	}

	allRooms := sync.WaitGroup{}
	allRooms.Add(len(ditems.DiscoveryItems))
	for _, di := range ditems.DiscoveryItems {
		go s.getRoomsInService(jid.Parse(di.Jid), di.Name, results, resultsServices, &allRooms)
	}
	allRooms.Wait()

	// This signals we are done
	results <- nil
}

func (s *session) getRoomsAsyncCustomService(service string, results chan<- *muc.RoomListing, resultsServices chan<- *muc.ServiceListing, errorResult chan<- error) {
	s.log.WithField("service", service).Debug("getRoomsAsyncCustomService()")

	allRooms := sync.WaitGroup{}
	allRooms.Add(1)
	go s.getRoomsInService(jid.Parse(service), "", results, resultsServices, &allRooms)
	allRooms.Wait()

	// This signals we are done
	results <- nil
}

func (s *session) GetRooms(server jid.Domain, customService string) (<-chan *muc.RoomListing, <-chan *muc.ServiceListing, <-chan error) {
	s.log.WithField("server", server).Debug("GetRooms()")
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

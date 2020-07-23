package session

import (
	"bytes"
	"encoding/xml"
	"errors"
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
		// TODO: log
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
				default:
					fmt.Printf("UNKNOWN FORM VAR: %s\n", k)
				}
			}
		}
	}

	rl.Updated()
}

func (s *session) getRoomsInService(service jid.Any, name string, results chan<- *muc.RoomListing, allRooms *sync.WaitGroup) {
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

	items, err := s.Conn().QueryServiceItems(service.String())
	if err != nil {
		// TODO: log
		return
	}

	for _, i := range items.DiscoveryItems {
		rl := muc.NewRoomListing()
		rl.Service = service
		rl.ServiceName = identName
		rl.Jid = jid.NR(i.Jid)
		rl.Name = i.Name

		results <- rl

		go s.findOutMoreInformationAboutRoom(rl)
	}

	allRooms.Done()
}

func (s *session) getRoomsAsync(server jid.Domain, results chan<- *muc.RoomListing, errorResult chan<- error) {
	rp, _, err := s.conn.SendIQ(server.String(), "get", &data.DiscoveryItemsQuery{})
	if err != nil {
		errorResult <- err
		return
	}

	r, ok := <-rp
	if !ok {
		errorResult <- errors.New("IQ channel closed")
		return
	}

	var ditems data.DiscoveryItemsQuery
	switch ciq := r.Value.(type) {
	case *data.ClientIQ:
		if ciq.Type != "result" {
			errorResult <- errors.New("got IQ result that is not 'result' type")
			return
		}
		if err := xml.NewDecoder(bytes.NewBuffer(ciq.Query)).Decode(&ditems); err != nil {
			errorResult <- err
			return
		}
	default:
		errorResult <- errors.New("got result to IQ that wasn't an IQ")
		return
	}

	allRooms := sync.WaitGroup{}
	allRooms.Add(len(ditems.DiscoveryItems))
	for _, di := range ditems.DiscoveryItems {
		go s.getRoomsInService(jid.Parse(di.Jid), di.Name, results, &allRooms)
	}
	allRooms.Wait()
	close(results)
	close(errorResult)
}

func (s *session) GetRooms(server jid.Domain) (<-chan *muc.RoomListing, <-chan error) {
	result := make(chan *muc.RoomListing, 20)
	errorResult := make(chan error, 1)

	go s.getRoomsAsync(server, result, errorResult)

	return result, errorResult
}

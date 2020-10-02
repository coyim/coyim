package muc

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

// RoomListing contains the information about a room for listing it
type RoomListing struct {
	Service     jid.Any
	ServiceName string
	Jid         jid.Bare
	Name        string

	SupportsVoiceRequests     bool
	AllowsRegistration        bool
	Anonymity                 string
	Persistent                bool
	Moderated                 bool
	Open                      bool
	PasswordProtected         bool
	Public                    bool
	Language                  string
	OccupantsCanChangeSubject bool
	Description               string
	Occupants                 int
	MembersCanInvite          bool
	OccupantsCanInvite        bool
	AllowPrivateMessages      string // This can be 'anyone', 'participants', 'moderators', 'none'
	ContactJid                string
	Logged                    bool // Notice that this will not always be correct for all servers

	lockUpdates sync.RWMutex
	onUpdates   []func()
}

// NewRoomListing creates and returns a new room listing
func NewRoomListing() *RoomListing {
	return &RoomListing{}
}

// OnUpdate takes a function and some data, and when this room listing is updated, that function
// will be called with the current room listing and the associated data
func (rl *RoomListing) OnUpdate(f func(*RoomListing, interface{}), data interface{}) {
	rl.lockUpdates.Lock()
	defer rl.lockUpdates.Unlock()

	rl.onUpdates = append(rl.onUpdates, func() {
		f(rl, data)
	})
}

// Updated should be called after a room listing has been updated, to notify observers of the update
func (rl *RoomListing) Updated() {
	rl.lockUpdates.RLock()
	defer rl.lockUpdates.RUnlock()

	for _, f := range rl.onUpdates {
		f()
	}
}

// SetFeatures receive a list of features and updates the room listing properties based on each feature
func (rl *RoomListing) SetFeatures(features []data.DiscoveryFeature) {
	rl.lockUpdates.Lock()
	defer rl.lockUpdates.Unlock()

	for _, feat := range features {
		rl.setFeature(feat.Var)
	}
}

// SetFeature updates a room listing propertie based on the given feature
func (rl *RoomListing) setFeature(feature string) {
	switch feature {
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
		fmt.Printf("UNKNOWN FEATURE: %s\n", feature)
	}
}

// SetFormsData extract the forms data and updates the room listing properties based on each data
func (rl *RoomListing) SetFormsData(forms []data.Form) {
	rl.lockUpdates.Lock()
	defer rl.lockUpdates.Unlock()

	for _, form := range forms {
		fields := formFieldsToKeyValue(form.Fields)
		if isValidRoomInfoForm(form, fields) {
			rl.updateWithFormFields(form, fields)
		}
	}
}

func (rl *RoomListing) updateWithFormFields(form data.Form, fields map[string][]string) {
	for field, values := range fields {
		rl.updateWithFormField(field, values)
	}
}

func (rl *RoomListing) updateWithFormField(field string, values []string) {
	switch field {
	case "FORM_TYPE":
		// Ignore, we already checked
	case "muc#roominfo_lang":
		if len(values) > 0 {
			rl.Language = values[0]
		}
	case "muc#roominfo_changesubject":
		if len(values) > 0 {
			rl.OccupantsCanChangeSubject = values[0] == "1"
		}
	case "muc#roomconfig_enablelogging":
		if len(values) > 0 {
			rl.Logged = values[0] == "1"
		}
	case "muc#roomconfig_roomname":
		// Room name - we already have this
	case "muc#roominfo_description":
		if len(values) > 0 {
			rl.Description = values[0]
		}
	case "muc#roominfo_occupants":
		if len(values) > 0 {
			res, e := strconv.Atoi(values[0])
			if e != nil {
				rl.Occupants = res
			}
		}
	case "{http://prosody.im/protocol/muc}roomconfig_allowmemberinvites":
		if len(values) > 0 {
			rl.MembersCanInvite = values[0] == "1"
		}
	case "muc#roomconfig_allowinvites":
		if len(values) > 0 {
			rl.OccupantsCanInvite = values[0] == "1"
		}
	case "muc#roomconfig_allowpm":
		if len(values) > 0 {
			rl.AllowPrivateMessages = values[0]
		}
	case "muc#roominfo_contactjid":
		if len(values) > 0 {
			rl.ContactJid = values[0]
		}
	default:
		fmt.Printf("UNKNOWN FORM VAR: %s\n", field)
	}
}

func formFieldsToKeyValue(fields []data.FormFieldX) map[string][]string {
	result := make(map[string][]string)
	for _, field := range fields {
		result[field.Var] = field.Values
	}

	return result
}

func isValidRoomInfoForm(form data.Form, fields map[string][]string) bool {
	return form.Type == "result" && hasRoomInfoFormType(fields)
}

func hasRoomInfoFormType(fields map[string][]string) bool {
	return len(fields["FORM_TYPE"]) > 0 && fields["FORM_TYPE"][0] == "http://jabber.org/protocol/muc#roominfo"
}

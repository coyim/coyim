package session

import (
	"time"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) loadRoomInfo(roomID jid.Bare) {
	result := make(chan *muc.RoomListing)
	go m.getRoomListing(roomID, result)

	select {
	case rl := <-result:
		m.onRoomConfigurationReceived(roomID, rl)
	case <-time.After(time.Minute * 5):
		m.roomConfigRequestTimeout(roomID)
	}
}

func (m *mucManager) onRoomConfigurationReceived(roomID jid.Bare, rl *muc.RoomListing) {
	m.addRoomInfo(roomID, rl)
	m.roomConfigReceived(roomID, rl.GetConfig())
}

var roomConfigUpdateCallers map[int]func(jid.Bare)

func (m *mucManager) getRoomConfigUpdateCallers() map[int]func(jid.Bare) {
	if len(roomConfigUpdateCallers) == 0 {
		roomConfigUpdateCallers = map[int]func(jid.Bare){
			MUCStatusRoomLoggingEnabled:  m.loggingEnabled,
			MUCStatusRoomLoggingDisabled: m.loggingDisabled,
			MUCStatusRoomNonAnonymous:    m.nonAnonymousRoom,
			MUCStatusRoomSemiAnonymous:   m.semiAnonymousRoom,
			MUCStatusConfigChanged:       m.nonPrivacyConfigChanged,
		}
	}

	return roomConfigUpdateCallers
}

func (m *mucManager) handleRoomConfigUpdate(stanza *xmppData.ClientMessage) {
	roomID := m.retrieveRoomID(stanza.From, "handleRoomConfigUpdate")

	status := mucUserStatuses(stanza.MUCUser.Status)

	for s, f := range m.getRoomConfigUpdateCallers() {
		if status.contains(s) {
			f(roomID)
		}
	}
}

func (m *mucManager) nonPrivacyConfigChanged(roomID jid.Bare) {
	rl, ok := m.getRoomInfo(roomID)
	if !ok {
		rl = m.newRoomListing(roomID)
		m.addRoomInfo(roomID, rl)
	}

	prevConfig := rl.GetConfig()
	m.findOutMoreInformationAboutRoom(rl)
	m.onRoomConfigurationUpdate(roomID, rl.GetConfig(), prevConfig)
}

var roomConfigUpdateCheckers = map[data.RoomConfigType]func(data.RoomConfig, data.RoomConfig) bool{
	data.RoomConfigSupportsVoiceRequests:     roomConfigSupportsVoiceRequestsCheckUpdate,
	data.RoomConfigAllowsRegistration:        roomConfigAllowsRegistrationCheckUpdate,
	data.RoomConfigPersistent:                roomConfigPersistentCheckUpdate,
	data.RoomConfigModerated:                 roomConfigModeratedCheckUpdate,
	data.RoomConfigOpen:                      roomConfigOpenCheckUpdate,
	data.RoomConfigPasswordProtected:         roomConfigPasswordProtectedCheckUpdate,
	data.RoomConfigPublic:                    roomConfigPublicCheckUpdate,
	data.RoomConfigLanguage:                  roomConfigLanguageCheckUpdate,
	data.RoomConfigOccupantsCanChangeSubject: roomConfigOccupantsCanChangeSubjectCheckUpdate,
	data.RoomConfigTitle:                     roomConfigTitleCheckUpdate,
	data.RoomConfigDescription:               roomConfigDescriptionCheckUpdate,
	data.RoomConfigOccupants:                 roomConfigOccupantsCheckUpdate,
	data.RoomConfigOccupantsCanInvite:        roomConfigOccupantsCanInviteCheckUpdate,
	data.RoomConfigAllowPrivateMessages:      roomConfigAllowPrivateMessagesCheckUpdate,
	data.RoomConfigLogged:                    roomConfigLoggedCheckUpdate,
}

func (m *mucManager) onRoomConfigurationUpdate(roomID jid.Bare, currConfig, prevConfig data.RoomConfig) {
	whatChanged := []data.RoomConfigType{}

	for kk, f := range roomConfigUpdateCheckers {
		if f(currConfig, prevConfig) {
			whatChanged = append(whatChanged, kk)
		}
	}

	m.roomConfigChanged(roomID, whatChanged, currConfig)
}

func roomConfigSupportsVoiceRequestsCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.SupportsVoiceRequests != prevConfig.SupportsVoiceRequests
}

func roomConfigAllowsRegistrationCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.AllowsRegistration != prevConfig.AllowsRegistration
}

func roomConfigPersistentCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.Persistent != prevConfig.Persistent
}

func roomConfigModeratedCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.Moderated != prevConfig.Moderated
}

func roomConfigOpenCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.Open != prevConfig.Open
}

func roomConfigPasswordProtectedCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.PasswordProtected != prevConfig.PasswordProtected
}

func roomConfigPublicCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.Public != prevConfig.Public
}

func roomConfigLanguageCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.Language != prevConfig.Language
}

func roomConfigOccupantsCanChangeSubjectCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.OccupantsCanChangeSubject != prevConfig.OccupantsCanChangeSubject
}

func roomConfigTitleCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.Title != prevConfig.Title
}

func roomConfigDescriptionCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.Description != prevConfig.Description
}

func roomConfigOccupantsCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.Occupants != prevConfig.Occupants
}

func roomConfigOccupantsCanInviteCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.OccupantsCanInvite != prevConfig.OccupantsCanInvite
}

func roomConfigAllowPrivateMessagesCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.AllowPrivateMessages != prevConfig.AllowPrivateMessages
}

func roomConfigLoggedCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.Logged != prevConfig.Logged
}

func isRoomConfigUpdate(stanza *xmppData.ClientMessage) bool {
	return hasMUCUserExtension(stanza)
}

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
			MUCStatusRoomLoggingEnabled:  m.handleLoggingEnabled,
			MUCStatusRoomLoggingDisabled: m.handleLoggingDisabled,
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

func (m *mucManager) obtainRoomInfo(roomID jid.Bare) *muc.RoomListing {
	rl, ok := m.getRoomInfo(roomID)
	if !ok {
		rl = m.newRoomListing(roomID)
		m.addRoomInfo(roomID, rl)
	}
	return rl
}

func (m *mucManager) handleLoggingEnabled(roomID jid.Bare) {
	rl := m.obtainRoomInfo(roomID)

	if !rl.Logged {
		m.loggingEnabled(roomID)
		rl.Logged = true
	}
}

func (m *mucManager) handleLoggingDisabled(roomID jid.Bare) {
	rl := m.obtainRoomInfo(roomID)

	if rl.Logged {
		m.loggingDisabled(roomID)
		rl.Logged = false
	}
}

func (m *mucManager) nonPrivacyConfigChanged(roomID jid.Bare) {
	rl := m.obtainRoomInfo(roomID)

	prevConfig := rl.GetConfig()
	m.findOutMoreInformationAboutRoom(rl)
	m.onRoomConfigUpdate(roomID, rl.GetConfig(), prevConfig)
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
	data.RoomConfigMembersCanInvite:          roomConfigMembersCanInviteCheckUpdate,
	data.RoomConfigAllowPrivateMessages:      roomConfigAllowPrivateMessagesCheckUpdate,
	data.RoomConfigLogged:                    roomConfigLoggedCheckUpdate,
	data.RoomConfigMaxHistoryFetch:           roomConfigMaxHistoryFetchCheckUpdate,
}

func (m *mucManager) onRoomConfigUpdate(roomID jid.Bare, currConfig, prevConfig data.RoomConfig) {
	changes := []data.RoomConfigType{}

	for k, f := range roomConfigUpdateCheckers {
		if f(currConfig, prevConfig) {
			changes = append(changes, k)
		}
	}

	// Some XMPP clients sends an update room configuration even when nothing has changed
	// We do this validation to avoid publish and room configuration change event when nothing has changed.
	if len(changes) > 0 {
		m.roomConfigChanged(roomID, changes, currConfig)
	}
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

func roomConfigMembersCanInviteCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.MembersCanInvite != prevConfig.MembersCanInvite
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

func roomConfigMaxHistoryFetchCheckUpdate(currConfig, prevConfig data.RoomConfig) bool {
	return currConfig.MaxHistoryFetch != prevConfig.MaxHistoryFetch
}

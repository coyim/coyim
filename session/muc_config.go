package session

import (
	"github.com/coyim/coyim/session/muc/data"
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
)

func (m *mucManager) handlersForRoomConfigurationChanges() map[int]func(jid.Bare) {
	m.roomConfigChangesHandlersLock.Lock()
	defer m.roomConfigChangesHandlersLock.Unlock()

	if len(m.roomConfigChangesHandlers) == 0 {
		m.roomConfigChangesHandlers = map[int]func(jid.Bare){
			MUCStatusRoomLoggingEnabled:  m.handleLoggingEnabled,
			MUCStatusRoomLoggingDisabled: m.handleLoggingDisabled,
			MUCStatusRoomNonAnonymous:    m.nonAnonymousRoom,
			MUCStatusRoomSemiAnonymous:   m.semiAnonymousRoom,
			MUCStatusConfigChanged:       m.nonPrivacyConfigChanged,
		}
	}

	return m.roomConfigChangesHandlers
}

func (m *mucManager) handleRoomConfigUpdate(stanza *xmppData.ClientMessage) {
	roomID := m.retrieveRoomID(stanza.From, "handleRoomConfigUpdate")

	status := mucUserStatuses(stanza.MUCUser.Status)

	for s, f := range m.handlersForRoomConfigurationChanges() {
		if status.contains(s) {
			f(roomID)
		}
	}
}

func (m *mucManager) handleLoggingEnabled(roomID jid.Bare) {
	rl := m.discoInfoForRoom(roomID)

	if !rl.Logged {
		m.loggingEnabled(roomID)
		rl.Logged = true
	}
}

func (m *mucManager) handleLoggingDisabled(roomID jid.Bare) {
	rl := m.discoInfoForRoom(roomID)

	if rl.Logged {
		m.loggingDisabled(roomID)
		rl.Logged = false
	}
}

func (m *mucManager) nonPrivacyConfigChanged(roomID jid.Bare) {
	rl := m.discoInfoForRoom(roomID)

	prevConfig := rl.GetDiscoInfo()
	m.findOutMoreInformationAboutRoom(rl)
	m.onRoomConfigUpdate(roomID, rl.GetDiscoInfo(), prevConfig)
}

var roomConfigUpdateCheckers = map[data.RoomConfigType]func(data.RoomDiscoInfo, data.RoomDiscoInfo) bool{
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

func (m *mucManager) onRoomConfigUpdate(roomID jid.Bare, currConfig, prevConfig data.RoomDiscoInfo) {
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

func roomConfigSupportsVoiceRequestsCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.SupportsVoiceRequests != prevConfig.SupportsVoiceRequests
}

func roomConfigAllowsRegistrationCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.AllowsRegistration != prevConfig.AllowsRegistration
}

func roomConfigPersistentCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.Persistent != prevConfig.Persistent
}

func roomConfigModeratedCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.Moderated != prevConfig.Moderated
}

func roomConfigOpenCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.Open != prevConfig.Open
}

func roomConfigPasswordProtectedCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.PasswordProtected != prevConfig.PasswordProtected
}

func roomConfigPublicCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.Public != prevConfig.Public
}

func roomConfigLanguageCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.Language != prevConfig.Language
}

func roomConfigOccupantsCanChangeSubjectCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.OccupantsCanChangeSubject != prevConfig.OccupantsCanChangeSubject
}

func roomConfigTitleCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.Title != prevConfig.Title
}

func roomConfigDescriptionCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.Description != prevConfig.Description
}

func roomConfigMembersCanInviteCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.MembersCanInvite != prevConfig.MembersCanInvite
}

func roomConfigAllowPrivateMessagesCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.AllowPrivateMessages != prevConfig.AllowPrivateMessages
}

func roomConfigLoggedCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.Logged != prevConfig.Logged
}

func isRoomConfigUpdate(stanza *xmppData.ClientMessage) bool {
	return hasMUCUserExtension(stanza)
}

func roomConfigMaxHistoryFetchCheckUpdate(currConfig, prevConfig data.RoomDiscoInfo) bool {
	return currConfig.MaxHistoryFetch != prevConfig.MaxHistoryFetch
}

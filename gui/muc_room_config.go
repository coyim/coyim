package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

type hasRoomConfigComponentView interface {
	canConfigureRoom() bool
	launchRoomConfigView()
}

type roomConfigScenario int

const (
	roomConfigScenarioCreate roomConfigScenario = iota
	roomConfigScenarioSubsequent
)

func (u *gtkUI) launchRoomConfigView(scenario roomConfigScenario, data *roomConfigData) {
	data.setScenario(scenario)
	data.ensureRequiredFields()

	var view hasRoomConfigComponentView
	switch scenario {
	case roomConfigScenarioCreate, roomConfigScenarioSubsequent:
		view = u.newRoomConfigAssistant(data)
	default:
		panic(fmt.Sprintf("developer error: trying to launch a not defined "+
			"room config view for scenario \"%v\"", scenario))
	}

	view.launchRoomConfigView()
}

type roomConfigData struct {
	roomConfigScenario              roomConfigScenario
	account                         *account
	roomID                          jid.Bare
	configForm                      *muc.RoomConfigForm
	autoJoinRoomAfterSaved          bool
	doAfterConfigSaved              func(autoJoin bool) // doAfterConfigSaved will be called from the UI thread
	doAfterConfigCanceled           func()              // doAfterConfigCanceled will be called from the UI thread
	doNotAskForConfirmationOnCancel bool
	parentWindow                    gtki.Window
}

func (rcd *roomConfigData) setScenario(scenario roomConfigScenario) {
	rcd.roomConfigScenario = scenario
}

func (rcd *roomConfigData) hasAccount() bool {
	return rcd.account != nil
}

func (rcd *roomConfigData) hasRoomID() bool {
	return rcd.roomID != nil
}

func (rcd *roomConfigData) hasConfigForm() bool {
	return rcd.configForm != nil
}

func (rcd *roomConfigData) hasRequiredFields() bool {
	return rcd.hasAccount() && rcd.hasRoomID() && rcd.hasConfigForm()
}

func (rcd *roomConfigData) ensureRequiredFields() {
	if !rcd.hasRequiredFields() {
		panic("Developer error: account, roomID and configForm should never be nil")
	}
}

// onConfigureRoom MUST be called from the UI thread
func (v *roomView) onConfigureRoom() {
	v.loadingViewOverlay.onRoomConfigurationRequest()

	fc, ec := v.account.session.GetRoomConfigurationForm(v.room.ID)
	go func() {
		select {
		case f := <-fc:
			doInUIThread(func() {
				v.u.launchRoomConfigView(roomConfigScenarioSubsequent, &roomConfigData{
					account:    v.account,
					roomID:     v.room.ID,
					configForm: f,
					doAfterConfigSaved: func(autoJoin bool) {
						v.notifications.info(roomNotificationOptions{
							message:   i18n.Local("The room configuration changed."),
							closeable: true,
						})
					},
					doNotAskForConfirmationOnCancel: true,
					parentWindow:                    v.window,
				})
			})
		case err := <-ec:
			v.log.WithError(err).Error("An error occurred when retrieving the Room Configuration Form")
			doInUIThread(func() {
				v.notifications.error(roomNotificationOptions{
					message:   i18n.Local("Unable to open the room configuration. Please, try again."),
					closeable: true,
				})
			})
		}
		doInUIThread(v.loadingViewOverlay.hide)
	}()
}

type roomConfigChangedTypes []data.RoomConfigType

func (c roomConfigChangedTypes) contains(k data.RoomConfigType) bool {
	for _, kk := range c {
		if kk == k {
			return true
		}
	}
	return false
}

var roomConfigFriendlyMessages map[data.RoomConfigType]func(data.RoomDiscoInfo) string

func initMUCConfigUpdateMessages() {
	roomConfigFriendlyMessages = map[data.RoomConfigType]func(data.RoomDiscoInfo) string{
		data.RoomConfigSupportsVoiceRequests:     roomConfigSupportsVoiceRequests,
		data.RoomConfigAllowsRegistration:        roomConfigAllowsRegistration,
		data.RoomConfigPersistent:                roomConfigPersistent,
		data.RoomConfigModerated:                 roomConfigModerated,
		data.RoomConfigPasswordProtected:         roomConfigPasswordProtected,
		data.RoomConfigPublic:                    roomConfigPublic,
		data.RoomConfigLanguage:                  roomConfigLanguage,
		data.RoomConfigOccupantsCanChangeSubject: roomConfigOccupantsCanChangeSubject,
		data.RoomConfigTitle:                     roomConfigTitle,
		data.RoomConfigDescription:               roomConfigDescription,
		data.RoomConfigMembersCanInvite:          roomConfigMembersCanInvite,
		data.RoomConfigAllowPrivateMessages:      roomConfigAllowPrivateMessages,
		data.RoomConfigMaxHistoryFetch:           roomConfigMaxHistoryFetch,
	}
}

func getRoomConfigUpdatedFriendlyMessages(changes roomConfigChangedTypes, discoInfo data.RoomDiscoInfo) []string {
	messages := []string{}

	for k, f := range roomConfigFriendlyMessages {
		if changes.contains(k) {
			messages = append(messages, f(discoInfo))
		}
	}

	return messages
}

func roomConfigSupportsVoiceRequests(di data.RoomDiscoInfo) string {
	if di.SupportsVoiceRequests {
		return i18n.Local("The room's occupants with no \"voice\" can now " +
			"request permission to speak.")
	}
	return i18n.Local("This room's doesn't support \"voice\" request " +
		"now, which means that visitors can't ask for permission to speak.")
}

func roomConfigAllowsRegistration(di data.RoomDiscoInfo) string {
	if di.AllowsRegistration {
		return i18n.Local("This room supports user registration.")
	}
	return i18n.Local("This room doesn't support user registration.")
}

func roomConfigPersistent(di data.RoomDiscoInfo) string {
	if di.Persistent {
		return i18n.Local("This room is now persistent.")
	}
	return i18n.Local("This room is not persistent anymore.")
}

func roomConfigModerated(di data.RoomDiscoInfo) string {
	if di.Moderated {
		return i18n.Local("Only participants and moderators can now send messages in this room.")
	}
	return i18n.Local("Everyone can now send messages in this room.")
}

func roomConfigPasswordProtected(di data.RoomDiscoInfo) string {
	if di.PasswordProtected {
		return i18n.Local("This room is now protected by a password.")
	}
	return i18n.Local("This room is not protected by a password.")
}

func roomConfigPublic(di data.RoomDiscoInfo) string {
	if di.Public {
		return i18n.Local("This room is publicly listed.")
	}
	return i18n.Local("This room is not publicly listed anymore.")
}

func roomConfigLanguage(di data.RoomDiscoInfo) string {
	return i18n.Localf("The language of this room was changed to %s", getLanguage(di.Language))
}

func getLanguage(languageCode string) string {
	languageTag, _ := language.Parse(languageCode)
	l := display.Self.Name(languageTag)
	if l == "" {
		return languageCode
	}
	return l
}

func roomConfigOccupantsCanChangeSubject(config data.RoomDiscoInfo) string {
	if config.OccupantsCanChangeSubject {
		return i18n.Local("Occupants can now change the subject of this room.")
	}
	return i18n.Local("Occupants cannot change the subject of this room.")
}

func roomConfigTitle(config data.RoomDiscoInfo) string {
	return i18n.Localf("Title was changed to \"%s\".", config.Title)
}

func roomConfigDescription(config data.RoomDiscoInfo) string {
	return i18n.Localf("Description was changed to \"%s\".", config.Description)
}

func roomConfigAllowPrivateMessages(config data.RoomDiscoInfo) string {
	if config.AllowPrivateMessages == data.RoleNone {
		return i18n.Local("No one in this room can send private messages now.")
	}
	return i18n.Localf("Only the \"%s\" can send private messages "+
		"to the room's occupants.", rolePluralName(config.AllowPrivateMessages))
}

func roomConfigMembersCanInvite(config data.RoomDiscoInfo) string {
	if config.MembersCanInvite {
		return i18n.Local("Members can now invite other users to join.")
	}
	return i18n.Local("Members cannot invite other users to join anymore.")
}

func roomConfigMaxHistoryFetch(config data.RoomDiscoInfo) string {
	return i18n.Localf("The room's max history length was changed to %d", config.MaxHistoryFetch)
}

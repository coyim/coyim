package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	log "github.com/sirupsen/logrus"
)

const (
	roomConfigInformationPageIndex = iota
	roomConfigAccessPageIndex
	roomConfigPermissionsPageIndex
	roomConfigOccupantsPageIndex
	roomConfigOthersPageIndex
	roomConfigSummaryPageIndex
)

type mucRoomConfigComponent struct {
	u        *gtkUI
	account  *account
	form     *muc.RoomConfigForm
	roomID   jid.Bare
	autoJoin bool

	infoPage        mucRoomConfigPage
	accessPage      mucRoomConfigPage
	permissionsPage mucRoomConfigPage
	occupantsPage   mucRoomConfigPage
	othersPage      mucRoomConfigPage
	summaryPage     mucRoomConfigPage

	log coylog.Logger
}

func (u *gtkUI) newMUCRoomConfigComponent(account *account, roomID jid.Bare, f *muc.RoomConfigForm, autoJoin bool) *mucRoomConfigComponent {
	c := &mucRoomConfigComponent{
		u:        u,
		account:  account,
		roomID:   roomID,
		form:     f,
		autoJoin: autoJoin,
		log: u.log.WithFields(log.Fields{
			"room":  roomID,
			"where": "roomConfigComponent",
		}),
	}

	c.initConfigPages()

	return c
}

func (c *mucRoomConfigComponent) initConfigPages() {
	c.infoPage = c.newRoomConfigInfoPage()
	c.accessPage = c.newRoomConfigAccessPage()
	c.permissionsPage = c.newRoomConfigPermissionsPage()
	c.occupantsPage = c.newRoomConfigOccupantsPage()
	c.othersPage = c.newRoomConfigOthersPage()
	c.summaryPage = c.newRoomConfigSummaryPage()
}

func (c *mucRoomConfigComponent) updateAutoJoin(v bool) {
	c.autoJoin = v
}

// cancelConfiguration IS SAFE to be called from the UI thread
func (c *mucRoomConfigComponent) cancelConfiguration(onSuccess func(), onError func(error)) {
	ec := c.account.session.CancelRoomConfiguration(c.roomID)

	onSuccessFinal := func() {
		if onSuccess != nil {
			onSuccess()
		}
	}

	onErrorFinal := func(err error) {
		if onError != nil {
			onError(err)
		}
	}

	go func() {
		if err := <-ec; err != nil {
			c.log.WithError(err).Error("An error occurred when trying to cancel the room configuration")
			onErrorFinal(err)
			return
		}

		onSuccessFinal()
	}()
}

// submitConfigurationForm IS SAFE to be called from the UI thread
func (c *mucRoomConfigComponent) submitConfigurationForm(onSuccess func(), onError func(error)) {
	rc, ec := c.account.session.SubmitRoomConfigurationForm(c.roomID, c.form)

	onSuccessFinal := func() {
		if onSuccess != nil {
			onSuccess()
		}
	}

	onErrorFinal := func(err error) {
		if onError != nil {
			onError(err)
		}
	}

	go func() {
		select {
		case <-rc:
			onSuccessFinal()
		case err := <-ec:
			c.log.WithError(err).Error("An error occurred when submitting the configuration form")
			onErrorFinal(err)
		}
	}()
}

func (c *mucRoomConfigComponent) getConfigPage(p int) mucRoomConfigPage {
	switch p {
	case roomConfigInformationPageIndex:
		return c.infoPage
	case roomConfigAccessPageIndex:
		return c.accessPage
	case roomConfigPermissionsPageIndex:
		return c.permissionsPage
	case roomConfigOccupantsPageIndex:
		return c.occupantsPage
	case roomConfigOthersPageIndex:
		return c.othersPage
	case roomConfigSummaryPageIndex:
		return c.summaryPage
	default:
		return nil
	}
}

func configOptionToFriendlyMessage(o string) string {
	switch o {
	case muc.RoomConfigOptionParticipants:
		return i18n.Local("Participants")
	case muc.RoomConfigOptionParticipant:
		return i18n.Local("Participant")
	case muc.RoomConfigOptionModerators:
		return i18n.Local("Moderators")
	case muc.RoomConfigOptionModerator:
		return i18n.Local("Moderator")
	case muc.RoomConfigOptionVisitor:
		return i18n.Local("Visitor")
	case muc.RoomConfigOptionAnyone:
		return i18n.Local("Anyone")
	case muc.RoomConfigOptionNone:
		return i18n.Local("None")
	case muc.RoomConfigOption10:
		return i18n.Local("10")
	case muc.RoomConfigOption20:
		return i18n.Local("20")
	case muc.RoomConfigOption30:
		return i18n.Local("30")
	case muc.RoomConfigOption50:
		return i18n.Local("50")
	case muc.RoomConfigOption100:
		return i18n.Local("100")
	default:
		return i18n.Localf("Unsupported \"%s\" option", o)
	}
}

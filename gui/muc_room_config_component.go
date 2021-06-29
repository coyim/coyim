package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type mucRoomConfigPageID int

const (
	roomConfigInformationPageIndex mucRoomConfigPageID = iota
	roomConfigAccessPageIndex
	roomConfigPermissionsPageIndex
	roomConfigPositionsPageIndex
	roomConfigOthersPageIndex
	roomConfigSummaryPageIndex
)

var roomConfigPages = []mucRoomConfigPageID{
	roomConfigInformationPageIndex,
	roomConfigAccessPageIndex,
	roomConfigPermissionsPageIndex,
	roomConfigPositionsPageIndex,
	roomConfigOthersPageIndex,
	roomConfigSummaryPageIndex,
}

type mucRoomConfigComponent struct {
	u              *gtkUI
	account        *account
	form           *muc.RoomConfigForm
	roomID         jid.Bare
	autoJoin       bool
	setCurrentPage func(indexPage mucRoomConfigPageID)
	pages          []*roomConfigPage
	log            coylog.Logger
}

func (u *gtkUI) newMUCRoomConfigComponent(account *account, roomID jid.Bare, f *muc.RoomConfigForm, autoJoin bool, setCurrentPage func(indexPage mucRoomConfigPageID), parent gtki.Window) *mucRoomConfigComponent {
	c := &mucRoomConfigComponent{
		u:              u,
		account:        account,
		roomID:         roomID,
		form:           f,
		autoJoin:       autoJoin,
		setCurrentPage: setCurrentPage,
		log: u.log.WithFields(log.Fields{
			"room":  roomID,
			"where": "roomConfigComponent",
		}),
	}

	c.initConfigPages(parent)

	return c
}

func (c *mucRoomConfigComponent) initConfigPages(parent gtki.Window) {
	for _, pageID := range roomConfigPages {
		c.pages = append(c.pages, c.newConfigPage(pageID))
	}
}

func (c *mucRoomConfigComponent) updateAutoJoin(v bool) {
	c.autoJoin = v
}

// configureRoom IS SAFE to be called from the UI thread
func (c *mucRoomConfigComponent) configureRoom(onSuccess func(), onError func(error)) {
	rc, ec := c.account.session.UpdateOccupantAffiliations(c.roomID, c.form.GetRoomOccupantsToUpdate())

	go func() {
		select {
		case <-rc:
			go c.submitConfigurationForm(onSuccess, onError)
		case err := <-ec:
			c.log.WithError(err).Error("An error occurred when configurating the occupant affiliations")
			doInUIThread(func() {
				onError(err)
			})
		}
	}()
}

func (c *mucRoomConfigComponent) submitConfigurationForm(onSuccess func(), onError func(error)) {
	rc, ec := c.account.session.SubmitRoomConfigurationForm(c.roomID, c.form)

	go func() {
		select {
		case <-rc:
			doInUIThread(onSuccess)
		case err := <-ec:
			c.log.WithError(err).Error("An error occurred when submitting the configuration form")
			doInUIThread(func() {
				onError(err)
			})
		}
	}()
}

func (c *mucRoomConfigComponent) getConfigPage(pageID mucRoomConfigPageID) (*roomConfigPage, bool) {
	for _, p := range c.pages {
		if p.pageID == pageID {
			return p, true
		}
	}
	return nil, false
}

func (c *mucRoomConfigComponent) friendlyConfigErrorMessage(err error) string {
	switch err {
	case session.ErrRoomConfigSubmit:
		return i18n.Local("We can't apply the given room configuration because an error occurred when trying to send the request for it. Please try again.")
	case session.ErrRoomConfigSubmitResponse:
		return i18n.Local("We can't apply the given room configuration because either you don't have the permissions for doing it or the location is not available right now. Please try again.")
	case session.ErrRoomConfigCancel:
		return i18n.Local("We can't cancel the room configuration process because an error occurred when trying to send the request for it. Please try again.")
	case session.ErrRoomConfigCancelResponse:
		return i18n.Local("We can't cancel the room configuration process because either you don't have the permissions for doing it or the location is not available right now. Please try again.")
	default:
		return i18n.Localf("Unsupported config error: %s", err)
	}
}

func configOptionToFriendlyMessage(o, defaultLabel string) string {
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
	case muc.RoomConfigOptionNobody:
		return i18n.Local("Nobody")
	case muc.RoomConfigOptionNone:
		return i18n.Local("No maximum")
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
	}
	return defaultLabel
}

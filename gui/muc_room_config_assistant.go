package gui

import (
	"fmt"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomConfigAssistant struct {
	u                   *gtkUI
	account             *account
	roomConfigComponent *mucRoomConfigComponent
	roomID              jid.Bare

	assistant          gtki.Assistant `gtk-widget:"room-config-assistant"`
	infoPageBox        gtki.Box       `gtk-widget:"room-config-info-page"`
	accessPageBox      gtki.Box       `gtk-widget:"room-config-access-page"`
	permissionsPageBox gtki.Box       `gtk-widget:"room-config-permissions-page"`
	occupantsPageBox   gtki.Box       `gtk-widget:"room-config-occupants-page"`
	othersPageBox      gtki.Box       `gtk-widget:"room-config-others-page"`
	summaryPageBox     gtki.Box       `gtk-widget:"room-config-summary-page"`

	infoPage        mucRoomConfigPage
	accessPage      mucRoomConfigPage
	permissionsPage mucRoomConfigPage
	occupantsPage   mucRoomConfigPage
	othersPage      mucRoomConfigPage
	summaryPage     mucRoomConfigPage

	onSuccess func(*account, jid.Bare)

	currentPageIndex int

	log coylog.Logger
}

func (u *gtkUI) newRoomConfigAssistant(account *account, roomID jid.Bare, form *muc.RoomConfigForm, onSuccess func(*account, jid.Bare)) *roomConfigAssistant {
	rc := &roomConfigAssistant{
		u:         u,
		account:   account,
		roomID:    roomID,
		onSuccess: onSuccess,
		log: u.log.WithFields(log.Fields{
			"room":  roomID,
			"where": "configureRoomAssistant",
		}),
	}

	rc.initBuilder()
	rc.initRoomConfigComponent(form)
	rc.initRoomConfigPages()
	rc.initDefaults()

	return rc
}

func (rc *roomConfigAssistant) initBuilder() {
	b := newBuilder("MUCRoomConfigAssistant")
	panicOnDevError(b.bindObjects(rc))

	b.ConnectSignals(map[string]interface{}{
		"on_cancel":       rc.onCancel,
		"on_page_changed": rc.onPageChanged,
		"on_apply":        rc.onApply,
	})
}

func (rc *roomConfigAssistant) initRoomConfigComponent(form *muc.RoomConfigForm) {
	rc.roomConfigComponent = rc.u.newMUCRoomConfigComponent(rc.roomID, form)
}

func (rc *roomConfigAssistant) initRoomConfigPages() {
	rc.infoPage = rc.roomConfigComponent.getConfigPage(roomConfigInformationPageIndex)
	rc.accessPage = rc.roomConfigComponent.getConfigPage(roomConfigAccessPageIndex)
	rc.permissionsPage = rc.roomConfigComponent.getConfigPage(roomConfigPermissionsPageIndex)
	rc.occupantsPage = rc.roomConfigComponent.getConfigPage(roomConfigOccupantsPageIndex)
	rc.othersPage = rc.roomConfigComponent.getConfigPage(roomConfigOthersPageIndex)
	rc.summaryPage = rc.roomConfigComponent.getConfigPage(roomConfigSummaryPageIndex)

	rc.infoPageBox.Add(rc.infoPage.getPageView())
	rc.accessPageBox.Add(rc.accessPage.getPageView())
	rc.permissionsPageBox.Add(rc.permissionsPage.getPageView())
	rc.occupantsPageBox.Add(rc.occupantsPage.getPageView())
	rc.othersPageBox.Add(rc.othersPage.getPageView())
	rc.summaryPageBox.Add(rc.summaryPage.getPageView())
}

func (rc *roomConfigAssistant) initDefaults() {
	rc.infoPageBox.SetHExpand(true)
	rc.accessPageBox.SetHExpand(true)
	rc.permissionsPageBox.SetHExpand(true)
	rc.occupantsPageBox.SetHExpand(true)
	rc.othersPageBox.SetHExpand(true)
	rc.summaryPageBox.SetHExpand(true)
}

func (rc *roomConfigAssistant) onCancel() {
	ec := rc.account.session.CancelRoomConfiguration(rc.roomID)
	go func() {
		err := <-ec
		if err != nil {
			// TODO: Show notification related to error produced trying to cancel the room configuration
			rc.log.WithError(err).Error("Error trying to cancel the room configuration")
			return
		}
	}()
	rc.assistant.Destroy()
}

func (rc *roomConfigAssistant) onPageChanged(_ gtki.Assistant, p gtki.Widget) {
	previousPage := rc.pageByIndex(rc.currentPageIndex)
	previousPage.collectData()

	rc.currentPageIndex = rc.assistant.GetCurrentPage()
	currentPage := rc.pageByIndex(rc.currentPageIndex)

	rc.assistant.SetPageComplete(p, true)
	currentPage.refresh()
}

func (rc *roomConfigAssistant) onApply() {
	sc, ec := rc.account.session.SubmitRoomConfigurationForm(rc.roomID, rc.roomConfigComponent.form)
	go func() {
		select {
		case <-sc:
			rc.onSuccess(rc.account, rc.roomID)
			doInUIThread(func() {
				rc.assistant.Destroy()
			})
		case err := <-ec:
			rc.log.WithField("error", err).Error("ERROR RECEIVED")
		}
	}()
}

func (rc *roomConfigAssistant) pageByIndex(p int) mucRoomConfigPage {
	switch p {
	case roomConfigInformationPageIndex:
		return rc.infoPage
	case roomConfigAccessPageIndex:
		return rc.accessPage
	case roomConfigPermissionsPageIndex:
		return rc.permissionsPage
	case roomConfigOccupantsPageIndex:
		return rc.occupantsPage
	case roomConfigOthersPageIndex:
		return rc.othersPage
	case roomConfigSummaryPageIndex:
		return rc.summaryPage
	default:
		panic(fmt.Sprintf("developer error: unsupported room config assistant page \"%d\"", p))
	}
}

func (rc *roomConfigAssistant) show() {
	rc.assistant.ShowAll()
}

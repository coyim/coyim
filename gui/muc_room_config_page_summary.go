package gui

import (
	"log"

	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigSummaryPage struct {
	*roomConfigPageBase

	box gtki.Box `gtk-widget:"room-config-summary-page"`
}

func (c *mucRoomConfigComponent) newRoomConfigSummaryPage() mucRoomConfigPage {
	p := &roomConfigSummaryPage{}

	builder := newBuilder("MUCRoomConfigPageSummary")
	panicOnDevError(builder.bindObjects(p))

	p.roomConfigPageBase = c.newConfigPage(p.box)
	p.onRefresh(p.onSummaryPageRefresh)

	return p
}

func (p *roomConfigSummaryPage) onSummaryPageRefresh() {
	log.Println("MaxHistoryFetch: ", p.form.MaxHistoryFetch)
	log.Println("AllowPrivateMessages: ", p.form.AllowPrivateMessages)
	log.Println("OccupantsCanInvite: ", p.form.OccupantsCanInvite)
	log.Println("OccupantsCanChangeSubject: ", p.form.OccupantsCanChangeSubject)
	log.Println("Logged: ", p.form.Logged)
	log.Println("RetrieveMembersList: ", p.form.RetrieveMembersList)
	log.Println("Language: ", p.form.Language)
	log.Println("AssociatedPublishSubscribeNode: ", p.form.AssociatedPublishSubscribeNode)
	log.Println("MaxOccupantsNumber: ", p.form.MaxOccupantsNumber)
	log.Println("MembersOnly: ", p.form.MembersOnly)
	log.Println("Moderated: ", p.form.Moderated)
	log.Println("PasswordProtected: ", p.form.PasswordProtected)
	log.Println("Persistent: ", p.form.Persistent)
	log.Println("PresenceBroadcast: ", p.form.PresenceBroadcast)
	log.Println("Public: ", p.form.Public)
	log.Println("Admins: ", p.form.Admins)
	log.Println("Description: ", p.form.Description)
	log.Println("Title: ", p.form.Title)
	log.Println("Owners: ", p.form.Owners)
	log.Println("Password: ", p.form.Password)
	log.Println("Whois: ", p.form.Whois)
}

func (p *roomConfigSummaryPage) collectData() {
	// Nothing to do, just implement the interface
}

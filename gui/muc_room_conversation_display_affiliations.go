package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
)

type affiliationUpdateDisplayer interface {
	affiliation() data.Affiliation
	addExtraMessageInfo() string
	updateReason() string
	previousAffiliationIsNone() bool
	displayForAffiliationRemoved() string
	displayForAffiliationOutcast() string
	displayForAffiliationAdded() string
	displayForAffiliationChanged() string
}

func displayAffiliationUpdateMessage(d affiliationUpdateDisplayer) (message string) {
	newAffiliation := d.affiliation()

	switch {
	case newAffiliation.IsNone():
		message = d.displayForAffiliationRemoved()
	case newAffiliation.IsBanned():
		message = d.displayForAffiliationOutcast()
	default:
		if d.previousAffiliationIsNone() {
			message = d.displayForAffiliationAdded()
		} else {
			message = d.displayForAffiliationChanged()
		}
	}

	if append := d.addExtraMessageInfo(); append != "" {
		message = i18n.Localf("%s %s", message, append)
	}

	if reason := d.updateReason(); reason != "" {
		message = i18n.Localf("%s because: %s", message, reason)
	}

	return message
}

func getDisplayForOccupantAffiliationUpdate(affiliationUpdate data.AffiliationUpdate) string {
	d := newAffiliationUpdateDisplayData(affiliationUpdate)
	return displayAffiliationUpdateMessage(d)
}

type affiliationUpdateDisplayData struct {
	nickname            string
	newAffiliation      data.Affiliation
	previousAffiliation data.Affiliation
	actor               string
	actorAffiliation    data.Affiliation
	reason              string
}

func newAffiliationUpdateDisplayData(affiliationUpdate data.AffiliationUpdate) *affiliationUpdateDisplayData {
	d := &affiliationUpdateDisplayData{
		nickname:            affiliationUpdate.Nickname,
		newAffiliation:      affiliationUpdate.New,
		previousAffiliation: affiliationUpdate.Previous,
		reason:              affiliationUpdate.Reason,
	}

	if affiliationUpdate.Actor != nil {
		d.actor = affiliationUpdate.Actor.Nickname
		d.actorAffiliation = affiliationUpdate.Actor.Affiliation
	}

	return d
}

func (d *affiliationUpdateDisplayData) affiliation() data.Affiliation {
	return d.newAffiliation
}

func (d *affiliationUpdateDisplayData) previousAffiliationIsNone() bool {
	return d.previousAffiliation.IsNone()
}

func (d *affiliationUpdateDisplayData) updateReason() string {
	return d.reason
}

func (d *affiliationUpdateDisplayData) displayForAffiliationRemoved() string {
	if d.actor == "" {
		return i18n.Localf("The %s position of %s was removed",
			displayNameForAffiliation(d.previousAffiliation), d.nickname)
	}
	return i18n.Localf("%s removed the %s position from %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForAffiliation(d.previousAffiliation),
		d.nickname,
	)
}

func (d *affiliationUpdateDisplayData) displayForAffiliationOutcast() string {
	if d.actor == "" {
		return i18n.Localf("%s was banned from the room", d.nickname)
	}
	return i18n.Localf("%s banned %s from the room",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		d.nickname,
	)
}

func (d *affiliationUpdateDisplayData) displayForAffiliationAdded() string {
	if d.actor == "" {
		return i18n.Localf("%s is now %s", d.nickname,
			displayNameForAffiliationWithPreposition(d.newAffiliation))
	}
	return i18n.Localf("%s changed the position of %s to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		d.nickname,
		displayNameForAffiliation(d.newAffiliation),
	)
}

func (d *affiliationUpdateDisplayData) displayForAffiliationChanged() string {
	if d.actor == "" {
		return i18n.Localf("The position of %s was changed from %s to %s", d.nickname,
			displayNameForAffiliation(d.previousAffiliation),
			displayNameForAffiliation(d.newAffiliation))
	}
	return i18n.Localf("%s changed the position of %s from %s to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		d.nickname,
		displayNameForAffiliation(d.previousAffiliation),
		displayNameForAffiliation(d.newAffiliation),
	)
}

func (d *affiliationUpdateDisplayData) addExtraMessageInfo() string {
	return ""
}

func getDisplayForSelfOccupantAffiliationUpdate(affiliationUpdate data.AffiliationUpdate) string {
	d := newSelfAffiliationUpdateDisplayData(affiliationUpdate)
	return displayAffiliationUpdateMessage(d)
}

type selfAffiliationUpdateDisplayData struct {
	*affiliationUpdateDisplayData
}

func newSelfAffiliationUpdateDisplayData(affiliationUpdate data.AffiliationUpdate) *selfAffiliationUpdateDisplayData {
	return &selfAffiliationUpdateDisplayData{
		newAffiliationUpdateDisplayData(affiliationUpdate),
	}
}

func (d *selfAffiliationUpdateDisplayData) displayForAffiliationRemoved() string {
	if d.actor == "" {
		return i18n.Localf("Your position of %s was removed", displayNameForAffiliation(d.previousAffiliation))
	}
	return i18n.Localf("%s removed your position as %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForAffiliation(d.previousAffiliation),
	)
}

func (d *selfAffiliationUpdateDisplayData) displayForAffiliationOutcast() string {
	if d.actor == "" {
		return i18n.Local("You were banned from the room")
	}
	return i18n.Localf("%s banned you from the room",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForAffiliation(d.actorAffiliation),
		d.actor,
	)
}

func (d *selfAffiliationUpdateDisplayData) displayForAffiliationAdded() string {
	if d.actor == "" {
		return i18n.Localf("You are now %s", displayNameForAffiliationWithPreposition(d.newAffiliation))
	}
	return i18n.Localf("%s changed your position to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForAffiliation(d.newAffiliation),
	)
}

func (d *selfAffiliationUpdateDisplayData) displayForAffiliationChanged() string {
	if d.actor == "" {
		return i18n.Localf("Your position was changed from %s to %s",
			displayNameForAffiliation(d.previousAffiliation),
			displayNameForAffiliation(d.newAffiliation))
	}
	return i18n.Localf("%s changed your position from %s to %s",
		displayActorWithAffiliation(d.actor, d.actorAffiliation),
		displayNameForAffiliation(d.previousAffiliation),
		displayNameForAffiliation(d.newAffiliation))
}

func (d *selfAffiliationUpdateDisplayData) addExtraMessageInfo() string {
	return ""
}

func displaySelfOccupantAffiliationUpdate(affiliationUpdate data.AffiliationUpdate) string {
	d := newSelfAffiliationUpdateDisplayData(affiliationUpdate)
	return displayAffiliationUpdateMessage(d)
}

func displayActorWithAffiliation(actor string, affiliation data.Affiliation) string {
	if affiliation != nil {
		return i18n.Localf("The %s %s", displayNameForAffiliation(affiliation), actor)
	}
	return actor
}

func displayNameForAffiliation(affiliation data.Affiliation) string {
	switch {
	case affiliation.IsAdmin():
		return i18n.Local("administrator")
	case affiliation.IsOwner():
		return i18n.Local("owner")
	case affiliation.IsBanned():
		return i18n.Local("outcast")
	case affiliation.IsMember():
		return i18n.Local("member")
	default: // Other values get the default treatment
		return ""
	}
}

func displayNameForAffiliationWithPreposition(affiliation data.Affiliation) string {
	switch {
	case affiliation.IsAdmin():
		return i18n.Local("an admininistrator")
	case affiliation.IsOwner():
		return i18n.Local("an owner")
	case affiliation.IsMember():
		return i18n.Local("a member")
	default: // Other values get the default treatment
		return ""
	}
}

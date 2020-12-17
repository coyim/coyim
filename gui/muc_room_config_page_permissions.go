package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigPermissionsPage struct {
	*roomConfigPageBase

	configPermissionsBox gtki.Box      `gtk-widget:"room-config-permissions-page"`
	roomChangeSubject    gtki.Switch   `gtk-widget:"room-changesubject"`
	roomModerated        gtki.Switch   `gtk-widget:"room-moderated"`
	roomWhois            gtki.ComboBox `gtk-widget:"room-whois"`
	roomWhoisModel       gtki.ListStore
	roomWhoisOptions     map[string]int
}

func (c *mucRoomConfigComponent) newRoomConfigPermissionsPage() mucRoomConfigPage {
	p := &roomConfigPermissionsPage{}

	builder := newBuilder("MUCRoomConfigPagePermissions")
	panicOnDevError(builder.bindObjects(p))

	p.roomConfigPageBase = c.newConfigPage(p.configPermissionsBox)
	p.onRefresh(p.refreshWhoisField)

	// These two values are the option name and the friendly label for it
	p.roomWhoisModel, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING, glibi.TYPE_STRING)
	p.roomWhois.SetModel(p.roomWhoisModel)

	p.initDefaultValues()

	return p
}

func (p *roomConfigPermissionsPage) initDefaultValues() {
	setSwitchActive(p.roomChangeSubject, p.form.OccupantsCanChangeSubject)
	setSwitchActive(p.roomModerated, p.form.Moderated)

	p.refreshWhoisField()
}

const (
	configWhoisOptionValueIndex int = iota
	configWhoisOptionLabelIndex
)

func (p *roomConfigPermissionsPage) refreshWhoisField() {
	p.roomWhoisModel.Clear()
	p.roomWhoisOptions = make(map[string]int)

	for index, o := range p.form.Whois.Options() {
		iter := p.roomWhoisModel.Append()
		_ = p.roomWhoisModel.SetValue(iter, configWhoisOptionValueIndex, o)
		_ = p.roomWhoisModel.SetValue(iter, configWhoisOptionLabelIndex, configOptionToFriendlyMessage(o))
		p.roomWhoisOptions[o] = index
	}

	p.activateWhoisOption(p.form.Whois.CurrentValue())
}

func (p *roomConfigPermissionsPage) activateWhoisOption(o string) {
	if index, ok := p.roomWhoisOptions[o]; ok {
		p.roomWhois.SetActive(index)
		return
	}
	p.log.WithField("option", o).Error("Trying to activate an unsupported \"whois\" field option")
}

func (p *roomConfigPermissionsPage) collectData() {
	p.form.OccupantsCanChangeSubject = getSwitchActive(p.roomChangeSubject)
	p.form.Moderated = getSwitchActive(p.roomModerated)

	for o, index := range p.roomWhoisOptions {
		if index == p.roomWhois.GetActive() {
			p.form.Whois.UpdateValue(o)
			return
		}
	}
}

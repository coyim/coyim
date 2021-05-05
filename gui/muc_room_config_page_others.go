package gui

import (
	"strconv"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigOthersPage struct {
	*roomConfigPageBase

	roomMaxHistoryFetchCombo gtki.ComboBoxText `gtk-widget:"room-maxhistoryfetch"`
	roomMaxHistoryFetchEntry gtki.Entry        `gtk-widget:"room-maxhistoryfetch-entry"`
	roomMaxOccupantsCombo    gtki.ComboBoxText `gtk-widget:"room-maxoccupants"`
	roomMaxOccupantsEntry    gtki.Entry        `gtk-widget:"room-maxoccupants-entry"`
	roomEnableLogging        gtki.Switch       `gtk-widget:"room-enablelogging"`
	roomUnknowFieldsBox      gtki.Box          `gtk-widget:"room-config-unknow-fields-box"`

	roomMaxHistoryFetch *roomConfigComboEntry
	roomMaxOccupants    *roomConfigComboEntry

	fields []hasRoomConfigFormField
}

func (c *mucRoomConfigComponent) newRoomConfigOthersPage() mucRoomConfigPage {
	p := &roomConfigOthersPage{}
	p.roomConfigPageBase = c.newConfigPage(pageConfigOthers, "MUCRoomConfigPageOthers", p, nil)

	p.roomMaxHistoryFetch = newRoomConfigCombo(p.roomMaxHistoryFetchCombo, p.roomMaxHistoryFetchEntry)
	p.roomMaxOccupants = newRoomConfigCombo(p.roomMaxOccupantsCombo, p.roomMaxOccupantsEntry)

	p.initDefaultValues()

	return p
}

func (p *roomConfigOthersPage) initDefaultValues() {
	p.roomMaxHistoryFetch.updateCurrentValue(p.form.MaxHistoryFetch.CurrentValue())
	p.roomMaxHistoryFetch.updateOptions(p.form.MaxHistoryFetch.Options())

	p.roomMaxOccupants.updateCurrentValue(p.form.MaxOccupantsNumber.CurrentValue())
	p.roomMaxOccupants.updateOptions(p.form.MaxOccupantsNumber.Options())

	p.roomEnableLogging.SetActive(p.form.Logged)

	for _, f := range p.form.Fields {
		field, err := roomConfigFormFieldFactory(f)
		if err != nil {
			p.log.WithField("field", f.Name).WithError(err).Error("Room configuration form field not supported")
			continue
		}

		p.addField(field)
		p.doAfterRefresh.add(field.refreshContent)
	}
}

// isInvalid MUST be called from the UI thread
func (p *roomConfigOthersPage) isInvalid() bool {
	return p.roomMaxHistoryFetch.isInvalid() || p.roomMaxOccupants.isInvalid()
}

// showValidationErrors MUST be called from the UI thread
func (p *roomConfigOthersPage) showValidationErrors() {
	p.clearErrors()

	if p.roomMaxHistoryFetch.isInvalid() {
		p.roomMaxHistoryFetch.focus()
		return
	}

	if p.roomMaxOccupants.isInvalid() {
		p.roomMaxOccupants.focus()
		return
	}
}

// collectData MUST be called from the UI thread
func (p *roomConfigOthersPage) collectData() {
	p.form.MaxHistoryFetch.UpdateValue(p.roomMaxHistoryFetch.currentValue())
	p.form.MaxOccupantsNumber.UpdateValue(p.roomMaxOccupants.currentValue())
	p.form.Logged = p.roomEnableLogging.GetActive()

	for _, f := range p.fields {
		p.form.UpdateFieldValueByName(f.fieldName(), f.fieldValue())
	}
}

// addField MUST BE called from the UI thread
func (p *roomConfigOthersPage) addField(f hasRoomConfigFormField) {
	p.fields = append(p.fields, f)
	p.roomUnknowFieldsBox.PackStart(f.fieldWidget(), true, false, 0)
}

const (
	roomMaxHistoryFetchValueColumIndex = iota
	roomMaxHistoryFetchLabelColumIndex
)

const (
	roomConfigComboOptionValueIndex = iota
	roomConfigComboOptionLabelIndex
)

type roomConfigComboEntry struct {
	options map[string]string

	model    gtki.ListStore
	comboBox gtki.ComboBoxText
	entry    gtki.Entry
}

func newRoomConfigCombo(cb gtki.ComboBoxText, e gtki.Entry) *roomConfigComboEntry {
	cc := &roomConfigComboEntry{
		comboBox: cb,
		entry:    e,
		options:  make(map[string]string),
	}

	// The following is created with two columns, one is for the "value" and the other for the "label"
	cc.model, _ = g.gtk.ListStoreNew(glibi.TYPE_STRING, glibi.TYPE_STRING)
	cc.comboBox.SetModel(cc.model)
	cc.comboBox.SetIDColumn(roomConfigComboOptionValueIndex)
	cc.comboBox.SetEntryTextColumn(roomConfigComboOptionLabelIndex)

	return cc
}

// updateCurrentValue MUST be called from the UI thread
func (cc *roomConfigComboEntry) updateCurrentValue(v string) {
	cc.entry.SetText(v)
}

// updateOptions MUST be called from the UI thread
func (cc *roomConfigComboEntry) updateOptions(options []string) {
	cc.model.Clear()
	cc.options = make(map[string]string)

	for _, o := range options {
		label := configOptionToFriendlyMessage(o)
		if o == muc.RoomConfigOptionNone {
			label = i18n.Local("No maximum (default)")
		}

		iter := cc.model.Append()
		cc.model.SetValue(iter, roomMaxHistoryFetchValueColumIndex, o)
		cc.model.SetValue(iter, roomMaxHistoryFetchLabelColumIndex, label)

		cc.options[label] = o
	}
}

// isInvalid MUST be called from the UI thread
func (cc *roomConfigComboEntry) isInvalid() bool {
	ct := getEntryText(cc.entry)
	if ct != "" {
		_, err := strconv.Atoi(cc.currentValue())
		if err != nil {
			return true
		}
	}
	return false
}

// currentValue MUST be called from the UI thread
func (cc *roomConfigComboEntry) currentValue() string {
	selected := cc.comboBox.GetActiveID()
	if selected != "" {
		return selected
	}

	ok := false
	entryText := getEntryText(cc.entry)

	selected, ok = cc.options[entryText]
	if ok {
		return selected
	}

	_, err := strconv.Atoi(entryText)
	if err != nil {
		return ""
	}

	return entryText
}

// focus MUST be called from the UI thread
func (cc *roomConfigComboEntry) focus() {
	cc.entry.GrabFocus()
}

package gui

import (
	"strconv"

	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigOthersPage struct {
	*roomConfigPageBase

	configOthersBox          gtki.Box          `gtk-widget:"room-config-others-page"`
	notificationBox          gtki.Box          `gtk-widget:"notification-box"`
	roomMaxHistoryFetchCombo gtki.ComboBoxText `gtk-widget:"room-maxhistoryfetch"`
	roomMaxHistoryFetchEntry gtki.Entry        `gtk-widget:"room-maxhistoryfetch-entry"`
	roomMaxOccupantsCombo    gtki.ComboBoxText `gtk-widget:"room-maxoccupants"`
	roomMaxOccupantsEntry    gtki.Entry        `gtk-widget:"room-maxoccupants-entry"`
	roomEnableLoggin         gtki.Switch       `gtk-widget:"room-enablelogging"`

	roomMaxHistoryFetch *roomConfigComboEntry
	roomMaxOccupants    *roomConfigComboEntry
}

func (c *mucRoomConfigComponent) newRoomConfigOthersPage() mucRoomConfigPage {
	p := &roomConfigOthersPage{}

	builder := newBuilder("MUCRoomConfigPageOthers")
	panicOnDevError(builder.bindObjects(p))

	p.roomConfigPageBase = c.newConfigPage(p.configOthersBox, p.notificationBox)

	p.roomMaxHistoryFetch = newRoomConfigCombo(p.roomMaxHistoryFetchCombo, p.roomMaxHistoryFetchEntry)
	p.roomMaxOccupants = newRoomConfigCombo(p.roomMaxOccupantsCombo, p.roomMaxOccupantsEntry)

	p.initDefaultValues()

	return p
}

func (p *roomConfigOthersPage) initDefaultValues() {
	setSwitchActive(p.roomEnableLoggin, p.form.Logged)

	p.roomMaxHistoryFetch.updateOptions(p.form.MaxHistoryFetch.Options())
	p.roomMaxOccupants.updateOptions(p.form.MaxOccupantsNumber.Options())
}

func (p *roomConfigOthersPage) collectData() {
	p.form.MaxHistoryFetch.UpdateValue(p.roomMaxHistoryFetch.currentValue())
	p.form.MaxOccupantsNumber.UpdateValue(p.roomMaxOccupants.currentValue())
	p.form.Logged = getSwitchActive(p.roomEnableLoggin)
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

func (cc *roomConfigComboEntry) updateOptions(options []string) {
	cc.model.Clear()
	cc.options = make(map[string]string)

	for _, o := range options {
		label := configOptionToFriendlyMessage(o)

		iter := cc.model.Append()
		cc.model.SetValue(iter, roomMaxHistoryFetchValueColumIndex, o)
		cc.model.SetValue(iter, roomMaxHistoryFetchLabelColumIndex, label)

		cc.options[label] = o
	}
}

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

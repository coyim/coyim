package gui

import (
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

const occupantsListJidColumnIndex = 0

type roomConfigOccupants struct {
	affiliation data.Affiliation

	content                  gtki.Box              `gtk-widget:"room-config-occupants-content"`
	header                   gtki.Label            `gtk-widget:"room-config-occupant-header"`
	occupantsListContent     gtki.Box              `gtk-widget:"room-config-occupants-list-content"`
	occupantsList            gtki.TreeView         `gtk-widget:"room-config-occupants-list"`
	occupantsAddButton       gtki.Button           `gtk-widget:"room-config-occupant-add"`
	occupantsRemoveButton    gtki.Button           `gtk-widget:"room-config-occupant-remove"`
	occupantsRemoveLabel     gtki.Label            `gtk-widget:"room-config-occupant-remove-label"`
	occupantsListJidRenderer gtki.CellRendererText `gtk-widget:"room-config-occupant-jid-text-renderer"`

	occupantsListController *mucRoomConfigListController
}

func newRoomConfigOccupants(a data.Affiliation) hasRoomConfigFormField {
	field := &roomConfigOccupants{affiliation: a}

	field.initBuilder()
	field.initDefaults()
	field.initOccupantsLists(nil)
	return field
}

func (p *roomConfigOccupants) initBuilder() {
	builder := newBuilder("MUCRoomConfigAssistantFieldOccupants")
	panicOnDevError(builder.bindObjects(p))
	builder.ConnectSignals(map[string]interface{}{
		"on_occupant_jid_edited": p.onOccupantJidEdited,
	})
}

func (p *roomConfigOccupants) initDefaults() {
	p.initHeaderLabel()
}

func (p *roomConfigOccupants) initHeaderLabel() {
	p.header.SetText(roomConfigOccupantFieldTexts[p.affiliation].headerLabel)
}

func (p *roomConfigOccupants) initOccupantsLists(parent gtki.Window) {
	p.initOwnersListController(parent)
}

func (p *roomConfigOccupants) initOwnersListController(parent gtki.Window) {
	p.occupantsListController = newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:      p.occupantsAddButton,
		removeOccupantButton:   p.occupantsRemoveButton,
		removeOccupantLabel:    p.occupantsRemoveLabel,
		occupantsTreeView:      p.occupantsList,
		parentWindow:           parent,
		addOccupantDialogTitle: roomConfigOccupantFieldTexts[p.affiliation].dialogTitle,
		addOccupantDescription: roomConfigOccupantFieldTexts[p.affiliation].dialogDescription,
		onListUpdated:          p.refreshContentLists,
	})
}

func (p *roomConfigOccupants) refreshContentLists() {
	p.occupantsListContent.SetVisible(p.occupantsListController.hasItems())
}

func (p *roomConfigOccupants) onOccupantJidEdited(_ gtki.CellRendererText, path string, newValue string) {
	p.updateOccupantListCellForString(p.occupantsListController, occupantsListJidColumnIndex, path, newValue)
}

func (p *roomConfigOccupants) updateOccupantListCellForString(controller *mucRoomConfigListController, column int, path string, newValue string) {
	if controller.updateCellForString(column, path, newValue) {
		log.WithFields(log.Fields{
			"path":        path,
			"newText":     newValue,
			"affiliation": p.affiliation.Name(),
		}).Debug("The occupant's jid can't be updated")
	}
}

func (p *roomConfigOccupants) collectFieldValue() {}

func (p *roomConfigOccupants) showValidationErrors() {}

func (p *roomConfigOccupants) fieldWidget() gtki.Widget {
	return p.content
}

func (p *roomConfigOccupants) refreshContent() {
	p.refreshContentLists()
}

func (p *roomConfigOccupants) isValid() bool {
	return true
}

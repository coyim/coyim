package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

const positionsListJidColumnIndex = 0

type roomConfigPositions struct {
	affiliation                     data.Affiliation
	updatePositionListByAffiliation func(affiliation data.Affiliation, occupants []*muc.RoomOccupantItem)

	content                  gtki.Box              `gtk-widget:"room-config-positions-content"`
	header                   gtki.Label            `gtk-widget:"room-config-position-header"`
	description              gtki.Label            `gtk-widget:"room-config-position-description"`
	positionsListContent     gtki.Box              `gtk-widget:"room-config-positions-list-content"`
	positionsList            gtki.TreeView         `gtk-widget:"room-config-positions-list"`
	positionsAddButton       gtki.Button           `gtk-widget:"room-config-position-add"`
	positionsRemoveButton    gtki.Button           `gtk-widget:"room-config-position-remove"`
	positionsRemoveLabel     gtki.Label            `gtk-widget:"room-config-position-remove-label"`
	positionsListJidRenderer gtki.CellRendererText `gtk-widget:"room-config-position-jid-text-renderer"`

	positionsListController *mucRoomConfigListController
}

func newRoomConfigPositions(affiliation data.Affiliation, updatePositionListByAffiliation func(affiliation data.Affiliation, occupants []*muc.RoomOccupantItem)) hasRoomConfigFormField {
	field := &roomConfigPositions{
		affiliation:                     affiliation,
		updatePositionListByAffiliation: updatePositionListByAffiliation,
	}

	field.initBuilder()
	field.initDefaults()
	field.initPositionsLists(nil)
	return field
}

func (p *roomConfigPositions) initBuilder() {
	builder := newBuilder("MUCRoomConfigAssistantFieldPositions")
	panicOnDevError(builder.bindObjects(p))
	builder.ConnectSignals(map[string]interface{}{
		"on_jid_edited": p.onOccupantJidEdited,
	})
}

func (p *roomConfigPositions) initDefaults() {
	p.initPositionLabels()
	mucStyles.setHelpTextStyle(p.content)
}

func (p *roomConfigPositions) initPositionLabels() {
	p.header.SetText(roomConfigOccupantFieldTexts[p.affiliation].headerLabel)
	p.description.SetText(roomConfigOccupantFieldTexts[p.affiliation].descriptionLabel)
}

func (p *roomConfigPositions) initPositionsLists(parent gtki.Window) {
	p.initOwnersListController(parent)
}

func (p *roomConfigPositions) initOwnersListController(parent gtki.Window) {
	p.positionsListController = newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:      p.positionsAddButton,
		removeOccupantButton:   p.positionsRemoveButton,
		removeOccupantLabel:    p.positionsRemoveLabel,
		occupantsTreeView:      p.positionsList,
		parentWindow:           parent,
		addOccupantDialogTitle: roomConfigOccupantFieldTexts[p.affiliation].dialogTitle,
		addOccupantDescription: roomConfigOccupantFieldTexts[p.affiliation].dialogDescription,
		onListUpdated:          p.refreshContentLists,
	})
}

func (p *roomConfigPositions) refreshContentLists() {
	p.positionsListContent.SetVisible(p.positionsListController.hasItems())
}

func (p *roomConfigPositions) onOccupantJidEdited(_ gtki.CellRendererText, path string, newValue string) {
	p.updateOccupantListCellForString(p.positionsListController, positionsListJidColumnIndex, path, newValue)
}

func (p *roomConfigPositions) updateOccupantListCellForString(controller *mucRoomConfigListController, column int, path string, newValue string) {
	if controller.updateCellForString(column, path, newValue) {
		log.WithFields(log.Fields{
			"path":        path,
			"newText":     newValue,
			"affiliation": p.affiliation.Name(),
		}).Debug("The occupant's jid can't be updated")
	}
}

func (p *roomConfigPositions) updateFieldValue() {
	positionsList := []*muc.RoomOccupantItem{}
	for _, item := range p.positionsListController.listItems() {
		positionsList = append(positionsList, &muc.RoomOccupantItem{
			Jid:         jid.Parse(item),
			Affiliation: p.affiliation,
		})
	}
	p.updatePositionListByAffiliation(p.affiliation, positionsList)
}

func (p *roomConfigPositions) showValidationErrors() {}

func (p *roomConfigPositions) fieldWidget() gtki.Widget {
	return p.content
}

func (p *roomConfigPositions) refreshContent() {
	p.refreshContentLists()
}

func (p *roomConfigPositions) isValid() bool {
	return true
}

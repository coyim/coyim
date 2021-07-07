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
	affiliation               data.Affiliation
	originalOccupantsList     muc.RoomOccupantItemList
	setOccupantList           func(occupants muc.RoomOccupantItemList)
	updateRemovedOccupantList func(occupantsToRemove muc.RoomOccupantItemList)

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

func newRoomConfigPositions(affiliation data.Affiliation, occupantsList muc.RoomOccupantItemList,
	setOccupantList func(occupants muc.RoomOccupantItemList),
	updateRemovedOccupantList func(occupantsToRemove muc.RoomOccupantItemList)) hasRoomConfigFormField {
	field := &roomConfigPositions{
		affiliation:               affiliation,
		originalOccupantsList:     occupantsList,
		setOccupantList:           setOccupantList,
		updateRemovedOccupantList: updateRemovedOccupantList,
	}

	field.initBuilder()
	field.initDefaults()
	field.initPositionsLists(nil)
	field.initOccupantList()

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
	p.header.SetText(getFieldTextByAffiliation(p.affiliation).headerLabel)
	p.description.SetText(getFieldTextByAffiliation(p.affiliation).descriptionLabel)
}

func (p *roomConfigPositions) initOccupantList() {
	jids := []string{}
	for _, o := range p.originalOccupantsList {
		jids = append(jids, o.Jid.String())
	}
	p.positionsListController.listComponent.addListItems(jids)
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
		addOccupantDialogTitle: getFieldTextByAffiliation(p.affiliation).dialogTitle,
		addOccupantDescription: getFieldTextByAffiliation(p.affiliation).dialogDescription,
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
	p.refreshOccupantLists(p.currentOccupantList())
}

func (p *roomConfigPositions) refreshOccupantLists(currentList muc.RoomOccupantItemList) {
	occupantsList := muc.RoomOccupantItemList{}
	for _, oi := range currentList {
		oi.MustBeUpdated = p.isNewOccupant(oi)
		occupantsList = append(occupantsList, oi)
	}
	p.setOccupantList(occupantsList)

	deletedOccupantsList := muc.RoomOccupantItemList{}
	for _, oi := range p.originalOccupantsList {
		if !currentList.IncludesJid(oi.Jid) {
			oi.ChangeAffiliationToNone()
			oi.MustBeUpdated = true
			deletedOccupantsList = append(deletedOccupantsList, oi)
		}
	}
	p.updateRemovedOccupantList(deletedOccupantsList)
}

func (p *roomConfigPositions) isNewOccupant(o *muc.RoomOccupantItem) bool {
	return !p.originalOccupantsList.IncludesJid(o.Jid)
}

func (p *roomConfigPositions) currentOccupantList() muc.RoomOccupantItemList {
	positionsList := []*muc.RoomOccupantItem{}
	for _, item := range p.positionsListController.listItems() {
		positionsList = append(positionsList, &muc.RoomOccupantItem{
			Jid:         jid.Parse(item),
			Affiliation: p.affiliation,
		})
	}
	return positionsList
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

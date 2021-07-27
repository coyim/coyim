package gui

import (
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

const positionsListJidColumnIndex = 0

type roomConfigPositionsOptions struct {
	affiliation            data.Affiliation
	occupantList           muc.RoomOccupantItemList
	setOccupantList        func(muc.RoomOccupantItemList) // setOccupantList WILL be called from the UI thread
	setRemovedOccupantList func(muc.RoomOccupantItemList) // setRemovedOccupantList WILL be called from the UI thread
	displayErrors          func()                         // displayErrors WILL be called from the UI thread
	parentWindow           gtki.Window
}

type roomConfigPositions struct {
	affiliation               data.Affiliation
	originalOccupantsList     muc.RoomOccupantItemList
	setOccupantList           func(occupants muc.RoomOccupantItemList)
	updateRemovedOccupantList func(occupantsToRemove muc.RoomOccupantItemList)
	showErrorNotification     func()

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

func newRoomConfigPositions(options roomConfigPositionsOptions) hasRoomConfigFormField {
	field := &roomConfigPositions{
		affiliation:               options.affiliation,
		originalOccupantsList:     options.occupantList,
		setOccupantList:           options.setOccupantList,
		updateRemovedOccupantList: options.setRemovedOccupantList,
		showErrorNotification:     options.displayErrors,
	}

	field.initBuilder()
	field.initDefaults()
	field.initPositionsLists(options.parentWindow)
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

func (p *roomConfigPositions) showValidationErrors() {
	if p.showErrorNotification != nil {
		p.showErrorNotification()
	}
}

func (p *roomConfigPositions) fieldWidget() gtki.Widget {
	return p.content
}

func (p *roomConfigPositions) refreshContent() {
	p.refreshContentLists()
}

func (p *roomConfigPositions) isValid() bool {
	return !(p.affiliation.IsOwner() && len(p.originalOccupantsList) != 0 && len(p.currentOccupantList()) == 0)
}

// fieldKey implements the hasRoomConfigFormField interface
func (p *roomConfigPositions) fieldKey() muc.RoomConfigFieldType {
	return muc.RoomConfigFieldUnexpected
}

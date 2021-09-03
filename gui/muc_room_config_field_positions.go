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

// [ps] Rename this struct to `roomConfigFieldPositions` to fix naming problems in this file
type roomConfigPositions struct {
	builder                   *builder
	affiliation               data.Affiliation
	originalOccupantsList     muc.RoomOccupantItemList
	setOccupantList           func(occupants muc.RoomOccupantItemList)
	updateRemovedOccupantList func(occupantsToRemove muc.RoomOccupantItemList)
	showErrorNotification     func()

	content               gtki.Box      `gtk-widget:"room-config-positions-content"`
	header                gtki.Label    `gtk-widget:"room-config-position-header"`
	description           gtki.Label    `gtk-widget:"room-config-position-description"`
	positionsListContent  gtki.Box      `gtk-widget:"room-config-positions-list-content"`
	positionsList         gtki.TreeView `gtk-widget:"room-config-positions-list"`
	positionsAddButton    gtki.Button   `gtk-widget:"room-config-position-add"`
	positionsRemoveButton gtki.Button   `gtk-widget:"room-config-position-remove"`
	positionsRemoveLabel  gtki.Label    `gtk-widget:"room-config-position-remove-label"`

	positionsListController *mucRoomConfigListController
}

func newRoomConfigPositionsComponent(options roomConfigPositionsOptions) *roomConfigPositions {
	rcp := &roomConfigPositions{
		affiliation:               options.affiliation,
		originalOccupantsList:     options.occupantList,
		setOccupantList:           options.setOccupantList,
		updateRemovedOccupantList: options.setRemovedOccupantList,
		showErrorNotification:     options.displayErrors,
	}

	rcp.loadUIDefinition()
	rcp.initDefaults()

	return rcp
}

func (p *roomConfigPositions) setUIBuilder(b *builder) {
	p.builder = b
}

func (p *roomConfigPositions) loadUIDefinition() {
	buildUserInterface("MUCRoomConfigFieldPositions", p, p.setUIBuilder)
}

type roomConfigPositionsWithApplyButton struct {
	*roomConfigPositions
	applyButton gtki.Button
}

func newRoomConfigPositionsWithApplyButton(applyButton gtki.Button, options roomConfigPositionsOptions) hasRoomConfigFormField {
	rcpb := &roomConfigPositionsWithApplyButton{
		roomConfigPositions: newRoomConfigPositionsComponent(options),
		applyButton:         applyButton,
	}

	rcpb.connectUISignals()
	rcpb.initPositionsLists(options.parentWindow)

	return rcpb
}

func (rcpb *roomConfigPositionsWithApplyButton) connectUISignals() {
	rcpb.builder.ConnectSignals(map[string]interface{}{
		"on_jid_edited": rcpb.onOccupantJidEdited,
	})
}

func (rcpb *roomConfigPositionsWithApplyButton) initPositionsLists(parent gtki.Window) {
	rcpb.positionsListController = newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:      rcpb.positionsAddButton,
		removeOccupantButton:   rcpb.positionsRemoveButton,
		removeOccupantLabel:    rcpb.positionsRemoveLabel,
		occupantsTreeView:      rcpb.positionsList,
		parentWindow:           parent,
		addOccupantDialogTitle: getFieldTextByAffiliation(rcpb.affiliation).dialogTitle,
		addOccupantDescription: getFieldTextByAffiliation(rcpb.affiliation).dialogDescription,
		onListUpdated:          rcpb.refreshContentLists,
	})

	rcpb.addItemsToListController()
}

// refreshContentLists MUST be called from the UI thread
func (rcpb *roomConfigPositionsWithApplyButton) refreshContentLists() {
	rcpb.roomConfigPositions.refreshContentLists()
	rcpb.enableOrDisableApplyButton()
}

// onOccupantJidEdited MUST be called from the UI thread
func (rcpb *roomConfigPositionsWithApplyButton) onOccupantJidEdited(cell gtki.CellRendererText, path string, newValue string) {
	rcpb.roomConfigPositions.onOccupantJidEdited(cell, path, newValue)
	rcpb.enableOrDisableApplyButton()
}

// enableOrDisableApplyButton MUST be called from the UI thread
func (rcpb *roomConfigPositionsWithApplyButton) enableOrDisableApplyButton() {
	rcpb.applyButton.SetSensitive(rcpb.hasListChanged())
}

type roomConfigPositionsField struct {
	*roomConfigPositions
}

func newRoomConfigPositionsField(options roomConfigPositionsOptions) hasRoomConfigFormField {
	rcpf := &roomConfigPositionsField{
		newRoomConfigPositionsComponent(options),
	}

	rcpf.connectUISignals()
	rcpf.initPositionsLists(options.parentWindow)

	return rcpf
}

func (rcpf *roomConfigPositionsField) connectUISignals() {
	rcpf.builder.ConnectSignals(map[string]interface{}{
		"on_jid_edited": rcpf.onOccupantJidEdited,
	})
}

func (rcpf *roomConfigPositionsField) initPositionsLists(parent gtki.Window) {
	rcpf.positionsListController = newMUCRoomConfigListController(&mucRoomConfigListControllerData{
		addOccupantButton:      rcpf.positionsAddButton,
		removeOccupantButton:   rcpf.positionsRemoveButton,
		removeOccupantLabel:    rcpf.positionsRemoveLabel,
		occupantsTreeView:      rcpf.positionsList,
		parentWindow:           parent,
		addOccupantDialogTitle: getFieldTextByAffiliation(rcpf.affiliation).dialogTitle,
		addOccupantDescription: getFieldTextByAffiliation(rcpf.affiliation).dialogDescription,
		onListUpdated:          rcpf.refreshContentLists,
	})

	rcpf.addItemsToListController()
}

func (p *roomConfigPositions) initDefaults() {
	p.initPositionLabels()
	mucStyles.setHelpTextStyle(p.content)
}

func (p *roomConfigPositions) initPositionLabels() {
	p.header.SetText(getFieldTextByAffiliation(p.affiliation).headerLabel)
	p.description.SetText(getFieldTextByAffiliation(p.affiliation).descriptionLabel)
}

// addItemsToListController MUST be called from the UI thread
func (p *roomConfigPositions) addItemsToListController() {
	jids := []string{}
	for _, o := range p.originalOccupantsList {
		jids = append(jids, o.Jid.String())
	}
	p.positionsListController.listComponent.addListItems(jids)
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

func (p *roomConfigPositions) hasListChanged() bool {
	ol := append(muc.RoomOccupantItemList{}, p.originalOccupantsList...)
	cl := append(muc.RoomOccupantItemList{}, p.currentOccupantList()...)

	if len(ol) != len(cl) {
		return true
	}

	for _, i := range cl {
		if !ol.IncludesJid(i.Jid) {
			return true
		}
	}

	return false
}

// fieldKey implements the hasRoomConfigFormField interface
func (p *roomConfigPositions) fieldKey() muc.RoomConfigFieldType {
	return muc.RoomConfigFieldUnexpected
}

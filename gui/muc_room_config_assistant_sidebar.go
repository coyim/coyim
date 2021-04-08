package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

const noRowIndex = -1

type roomConfigAssistantSidebar struct {
	assistant *roomConfigAssistant

	box     gtki.Box     `gtk-widget:"assistant-options-content"`
	listBox gtki.ListBox `gtk-widget:"assistant-options"`
}

func (rc *roomConfigAssistant) newRoomConfigAssistantSidebar() *roomConfigAssistantSidebar {
	sb := &roomConfigAssistantSidebar{
		assistant: rc,
	}

	sb.initBuilder()

	return sb
}

func (sb *roomConfigAssistantSidebar) initBuilder() {
	b := newBuilder("MUCRoomConfigAssistantSidebar")
	panicOnDevError(b.bindObjects(sb))

	b.ConnectSignals(map[string]interface{}{
		"row_selected": sb.onRowSelected,
	})
}

// onRowSelected MUST be called from the UI thread
func (sb *roomConfigAssistantSidebar) onRowSelected(_ gtki.ListBox, r gtki.ListBoxRow) {
	sb.assistant.updateAssistantPage(r.GetIndex())
}

// selectOptionByIndex MUST be called from the UI thread
func (sb *roomConfigAssistantSidebar) selectOptionByIndex(idx int) {
	row := sb.listBox.GetRowAtIndex(idx)
	rowIndex := getListBoxRowIndex(row)
	currentRowIndex := getListBoxRowIndex(sb.listBox.GetSelectedRow())

	if rowIndex != noRowIndex && rowIndex != currentRowIndex {
		sb.listBox.SelectRow(row)
	}
}

func getListBoxRowIndex(r gtki.ListBoxRow) int {
	if r != nil {
		return r.GetIndex()
	}
	return noRowIndex
}

package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigAssistantSidebar struct {
	assistant *roomConfigAssistant

	content gtki.Box     `gtk-widget:"assistant-options-content"`
	options gtki.ListBox `gtk-widget:"assistant-options"`
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

func (sb *roomConfigAssistantSidebar) onRowSelected(_ gtki.ListBox, r gtki.ListBoxRow) {
	sb.assistant.updateAssistantPage(r.GetIndex())
}

func (sb *roomConfigAssistantSidebar) selectOption(indexPage int) {
	sb.options.SelectRow(sb.options.GetRowAtIndex(indexPage))
}

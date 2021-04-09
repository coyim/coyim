package gui

import "github.com/coyim/gotk3adapter/gtki"

// findAssistantHeaderContainer MUST be called from the UI thread
func findAssistantHeaderContainer(a gtki.Assistant) gtki.Container {
	lbl, _ := g.gtk.LabelNew("")
	a.AddActionWidget(lbl)
	parentBox, _ := lbl.GetParentX()
	a.RemoveActionWidget(lbl)
	return parentBox.(gtki.Container)
}

type assistantButtons map[string]gtki.Button

// getButtonsForAssistantHeader MUST be called from the UI thread
func getButtonsForAssistantHeader(a gtki.Assistant) assistantButtons {
	h := findAssistantHeaderContainer(a)
	result := assistantButtons{}

	for _, c := range h.GetChildren() {
		if b, ok := c.(gtki.Button); ok {
			name, _ := g.gtk.GetWidgetBuildableName(b)
			result[name] = b
		}
	}

	return result
}

// updateButtonLabelByName MUST be called from the UI thread
func (list assistantButtons) updateButtonLabelByName(name string, label string) {
	if b, ok := list[name]; ok {
		b.SetLabel(label)
	}
}

// getBottomActionAreaFromAssistant MUST be called from the UI thread
func getBottomActionAreaFromAssistant(a gtki.Assistant) (gtki.Box, bool) {
	return findGtkBoxWithID(a.GetChildren(), "action_area")
}

// getSidebarFromAssistant MUST be called from the UI thread
func getSidebarFromAssistant(a gtki.Assistant) (gtki.Box, bool) {
	return findGtkBoxWithID(a.GetChildren(), "sidebar")
}

// setAssistantSidebar MUST be called from the UI thread
func setAssistantSidebarContent(a gtki.Assistant, content gtki.Widget) {
	if sidebar, ok := getSidebarFromAssistant(a); ok {
		for _, ch := range sidebar.GetChildren() {
			sidebar.Remove(ch)
		}
		sidebar.PackStart(content, false, false, 0)
	}
}

func findGtkBoxWithID(list []gtki.Widget, boxName string) (gtki.Box, bool) {
	for _, widget := range list {
		if box, ok := widget.(gtki.Box); ok {
			if name, _ := g.gtk.GetWidgetBuildableName(box); name == boxName {
				return box, true
			}
			if box, ok = findGtkBoxWithID(box.GetChildren(), boxName); ok {
				return box, true
			}
		}
	}
	return nil, false
}

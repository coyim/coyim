package gui

import "github.com/coyim/gotk3adapter/gtki"

func findAssistantHeaderContainer(a gtki.Assistant) gtki.Container {
	lbl, _ := g.gtk.LabelNew("")
	a.AddActionWidget(lbl)
	parentBox, _ := lbl.GetParentX()
	a.RemoveActionWidget(lbl)
	return parentBox.(gtki.Container)
}

func getButtonsForAssistantHeader(a gtki.Assistant) []gtki.Button {
	h := findAssistantHeaderContainer(a)
	result := []gtki.Button{}

	for _, c := range h.GetChildren() {
		if b, ok := c.(gtki.Button); ok {
			result = append(result, b)
		}
	}

	return result
}

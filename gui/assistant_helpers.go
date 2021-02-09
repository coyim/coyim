package gui

import "github.com/coyim/gotk3adapter/gtki"

func findAssistantHeaderContainer(a gtki.Assistant) gtki.Container {
	lbl, _ := g.gtk.LabelNew("")
	a.AddActionWidget(lbl)
	parentBox, _ := lbl.GetParentX()
	a.RemoveActionWidget(lbl)
	return parentBox.(gtki.Container)
}

type assistantButtons map[string]gtki.Button

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

func (list assistantButtons) updateButtonLabelByName(name string, label string) {
	if b, ok := list[name]; ok {
		b.SetLabel(label)
	}
}

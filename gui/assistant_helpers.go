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

func updateSidebarContent(a gtki.Assistant, content gtki.Box) {
	if s, ok := findBoxWidgetByName(a.GetChildren(), "sidebar"); ok {
		removeAllBoxChildrens(s)
		s.PackStart(content, false, false, 0)
	}
}

func removeAllBoxChildrens(box gtki.Box) {
	for _, ch := range box.GetChildren() {
		box.Remove(ch)
	}
}

func hideActionArea(a gtki.Assistant) {
	if actionArea, ok := findBoxWidgetByName(a.GetChildren(), "action_area"); ok {
		actionArea.SetVisible(false)
	}
}

func findBoxWidgetByName(wl []gtki.Widget, wn string) (gtki.Box, bool) {
	for _, c := range wl {
		if b, ok := c.(gtki.Box); ok {
			if name, _ := g.gtk.GetWidgetBuildableName(b); name == wn {
				return b, true
			}
			if b, ok := findBoxWidgetByName(b.GetChildren(), wn); ok {
				return b, ok
			}
		}
	}
	return nil, false
}

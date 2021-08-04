package gui

import "github.com/coyim/gotk3adapter/gtki"

const (
	assistantButtonBackLastName = "back"
	assistantButtonLastName     = "last"
	assistantButtonForwardName  = "forward"
	assistantButtonApplyName    = "apply"
)

var assistantNavigationButtons = []string{
	assistantButtonBackLastName,
	assistantButtonLastName,
	assistantButtonForwardName,
	assistantButtonApplyName,
}

type assistantButtons map[string]gtki.Button

// getButtonsForAssistantHeader MUST be called from the UI thread
func getButtonsForAssistantHeader(a gtki.Assistant) assistantButtons {
	result := assistantButtons{}

	for _, button := range a.GetButtons() {
		name, _ := g.gtk.GetWidgetBuildableName(button)
		result[name] = button
	}

	return result
}

// updateLastButtonLabel MUST be called from the UI thread
func (list assistantButtons) updateLastButtonLabel(label string) {
	list.updateButtonLabelByName(assistantButtonLastName, label)
}

// updateApplyButtonLabel MUST be called from the UI thread
func (list assistantButtons) updateApplyButtonLabel(label string) {
	list.updateButtonLabelByName(assistantButtonApplyName, label)
}

// updateButtonLabelByName MUST be called from the UI thread
func (list assistantButtons) updateButtonLabelByName(name string, label string) {
	if b, ok := list[name]; ok {
		b.SetLabel(label)
	}
}

// disableNavigationButNotCancel MUST be called from the UI thread
func (list assistantButtons) disableNavigationButNotCancel() {
	for _, buttonName := range assistantNavigationButtons {
		if b, ok := list[buttonName]; ok {
			b.SetSensitive(false)
		}
	}
}

// enableNavigation MUST be called from the UI thread
func (list assistantButtons) enableNavigation() {
	for _, buttonName := range assistantNavigationButtons {
		if b, ok := list[buttonName]; ok {
			b.SetSensitive(true)
		}
	}
}

// removeMarginFromAssistantPages MUST be called from the UI thread
func removeMarginFromAssistantPages(a gtki.Assistant) {
	if notebook, err := a.GetNotebook(); err == nil {
		for _, page := range notebook.GetChildren() {
			page.SetProperty("margin", 0)
		}
	}
}

// setAssistantSidebar MUST be called from the UI thread
func setAssistantSidebarContent(a gtki.Assistant, content gtki.Widget) {
	if sidebar, err := a.GetSidebar(); err == nil {
		for _, ch := range sidebar.GetChildren() {
			sidebar.Remove(ch)
		}
		sidebar.PackStart(content, false, false, 0)
	}
}

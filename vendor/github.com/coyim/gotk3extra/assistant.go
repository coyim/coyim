package gotk3extra

import "github.com/gotk3/gotk3/gtk"

const assistantButtonSizeGroupName = "button_size_group"

func GetAssistantButtonSizeGroup(a *gtk.Assistant) (*gtk.SizeGroup, error) {
	obj, err := GetWidgetTemplateChild(a, assistantButtonSizeGroupName)
	return WrapSizeGroup(obj, err)
}

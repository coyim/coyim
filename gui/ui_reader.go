package gui

import (
	"fmt"
	"github.com/gotk3/gotk3/gtk"
)

type UIDefinition interface {
	getDefinition() string
}

const defsFolder string = "definitions"
const xmlExtension string = ".xml"

func loadBuilderWith(uiName string) (*gtk.Builder, error) {
	//TODO: Add OS-aware path separator
	fileName := defsFolder + "/" + uiName + xmlExtension
	builder, _ := gtk.BuilderNew()

	fileErr := builder.AddFromFile(fileName)
	if fileErr != nil {
		uiDef := getDefinition(uiName)
		if uiDef == nil {
			return nil, fmt.Errorf("There's no definition for %s", uiName)
		}
		addErr := builder.AddFromString(uiDef.getDefinition())
		if addErr != nil {
			return nil, fileErr
		}
	}
	return builder, nil
}

func getDefinition(uiName string) UIDefinition {
	switch uiName {
	default:
		return nil
	case "TestWindow":
		return new(TestWindow)
	}
}

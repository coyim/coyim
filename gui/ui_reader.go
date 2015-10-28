package gui

import (
	"bufio"
	"fmt"
	"github.com/gotk3/gotk3/gtk"
	"os"
	"strings"
)

type UIDefinition interface {
	getDefinition() string
}

const defsFolder string = "gui/definitions"
const xmlExtension string = ".xml"

func loadBuilderWith(uiName string, vars map[string]string) (*gtk.Builder, error) {
	//TODO: Add OS-aware path separator
	fileName := defsFolder + "/" + uiName + xmlExtension
	builder, _ := gtk.BuilderNew()
	var toReplace string
	if doesnotExist(fileName) {
		uiDef := getDefinition(uiName)
		if uiDef == nil {
			return nil, fmt.Errorf("There's no definition for %s", uiName)
		}
		toReplace = uiDef.getDefinition()
	} else {
		toReplace = readFile(fileName)
	}

	replaced := replaceVars(toReplace, vars)

	addErr := builder.AddFromString(replaced)
	if addErr != nil {
		return nil, addErr
	}

	return builder, nil
}

func doesnotExist(fileName string) bool {
	_, fnf := os.Stat(fileName)
	return os.IsNotExist(fnf)
}

func readFile(fileName string) string {
	file, _ := os.Open(fileName)
	reader := bufio.NewScanner(file)
	var content string
	for reader.Scan() {
		content = content + reader.Text()
	}
	file.Close()
	return content
}

func replaceVars(toReplace string, vars map[string]string) string {
	replaced := toReplace
	for k, v := range vars {
		replaced = strings.Replace(replaced, k, v, -1)
	}
	return replaced
}

func getDefinition(uiName string) UIDefinition {
	switch uiName {
	default:
		return nil
	case "TestWindow":
		return new(TestWindow)
	}
}

package gui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type uiDefinition interface {
	getDefinition() string
}

const (
	defsFolder   string = "gui/definitions"
	xmlExtension string = ".xml"
)

//hold a reference to them to prevent garbage collecting
var builders map[string]*gtk.Builder

func loadBuilderWith(uiName string, vars map[string]string) (*gtk.Builder, error) {
	if builders == nil {
		builders = make(map[string]*gtk.Builder)
	}

	builder, ok := builders[uiName]
	if ok {
		return builder, nil
	}

	fileName := filepath.Join(defsFolder, uiName+xmlExtension)
	builder, _ = gtk.BuilderNew()
	var toReplace string
	if doesnotExist(fileName) {
		log.Printf("Loading compiled definition %q")
		uiDef := getDefinition(uiName)
		if uiDef == nil {
			return nil, fmt.Errorf("There's no definition for %s", uiName)
		}
		toReplace = uiDef.getDefinition()
	} else {
		log.Printf("Loading UI definition %q from: %s", uiName, fileName)
		toReplace = readFile(fileName)
	}

	replaced := replaceVars(toReplace, vars)

	addErr := builder.AddFromString(replaced)
	if addErr != nil {
		return nil, addErr
	}

	builders[uiName] = builder
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

func getDefinition(uiName string) uiDefinition {
	switch uiName {
	default:
		return nil
	case "TestWindow":
		return new(testWindow)
	}
}

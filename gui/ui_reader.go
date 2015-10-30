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

func loadBuilderWith(uiName string, vars map[string]string) (*gtk.Builder, error) {
	fileName := filepath.Join(defsFolder, uiName+xmlExtension)
	builder, err := gtk.BuilderNew()
	if err != nil {
		return nil, err
	}

	var toReplace string
	if doesnotExist(fileName) {
		log.Printf("Loading compiled definition %q\n", uiName)
		uiDef := getDefinition(uiName)
		if uiDef == nil {
			return nil, fmt.Errorf("There's no definition for %s", uiName)
		}
		toReplace = uiDef.getDefinition()
	} else {
		log.Printf("Loading UI definition %q from: %s\n", uiName, fileName)
		toReplace = readFile(fileName)
	}

	replaced := replaceVars(toReplace, vars)

	addErr := builder.AddFromString(replaced)
	if addErr != nil {
		log.Printf("Failed to add string %s: %s\n", replaced, addErr.Error())
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

func getDefinition(uiName string) uiDefinition {
	switch uiName {
	default:
		return nil
	case "MainDefinition":
		return new(mainDefinition)
	case "ConversationDefinition":
		return new(conversationDefinition)
	case "TestDefinition":
		return new(testDefinition)
	}
}

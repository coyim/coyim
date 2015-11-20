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
	defsFolder   = "gui/definitions"
	xmlExtension = ".xml"
)

func getDefinitionWithFileFallback(uiName string) string {
	// this makes sure a missing definition wont break only when the app is released
	uiDef := getDefinition(uiName)

	fileName := filepath.Join(defsFolder, uiName+xmlExtension)
	if fileNotFound(fileName) {
		log.Printf("gui: loading compiled definition %q\n", uiName)
		return uiDef.getDefinition()
	}

	return readFile(fileName)
}

func loadBuilderWith(uiName string, vars map[string]string) (*gtk.Builder, error) {
	//TODO: replace this by gettext
	replaced := replaceVars(getDefinitionWithFileFallback(uiName), vars)

	builder, err := gtk.BuilderNew()
	if err != nil {
		return nil, err
	}

	err = builder.AddFromString(replaced)
	if err != nil {
		log.Printf("gui: failed load %s: %s\n", uiName, err.Error())
		return nil, err
	}

	return builder, nil
}

func fileNotFound(fileName string) bool {
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
		panic(fmt.Sprintf("No definition found for %s", uiName))
	case "MainDefinition":
		return new(mainDefinition)
	case "AddContactDefinition":
		return new(accountDetailsDefinition)
	case "AccountDetailsDefinition":
		return new(accountDetailsDefinition)
	case "MasterPasswordDefinition":
		return new(accountDetailsDefinition)
	case "ConversationDefinition":
		return new(conversationDefinition)
	case "ConfigAssistantDefinition":
		return new(configAssistantDefinition)
	case "TorNotRunningDef":
		return new(torNotRunningDef)
	case "ConnectionSettingsDialogDef":
		return new(connectionSettingsDialogDef)
	case "TestDefinition":
		return new(testDefinition)
	}
}

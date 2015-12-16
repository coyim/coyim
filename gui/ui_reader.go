package gui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/gtk"
	"../gui/definitions"
)

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
		return uiDef.String()
	}

	return readFile(fileName)
}

func builderForDefinition(uiName string) *gtk.Builder {
	template := getDefinitionWithFileFallback(uiName)

	builder, err := gtk.BuilderNew()
	if err != nil {
		//We cant recover from this
		panic(err)
	}

	err = builder.AddFromString(template)
	if err != nil {
		//This is a programming error
		panic(fmt.Sprintf("gui: failed load %s: %s\n", uiName, err.Error()))
	}

	return builder
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

func getDefinition(uiName string) fmt.Stringer {
	def, ok := definitions.Get(uiName)
	if !ok {
		panic(fmt.Sprintf("No definition found for %s", uiName))
	}

	return def
}

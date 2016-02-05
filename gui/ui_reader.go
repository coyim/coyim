package gui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"reflect"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/gui/definitions"
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

// This must be called from the UI thread - otherwise bad things will happen sooner or later
func builderForDefinition(uiName string) *gtk.Builder {
	// assertInUIThread()

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

type builder struct {
	*gtk.Builder
}

func newBuilder(uiName string) *builder {
	return &builder{builderForDefinition(uiName)}
}

func (b *builder) getItem(name string, target interface{}) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		panic("builder.getItem() target argument must be a pointer")
	}
	elem := v.Elem()
	elem.Set(reflect.ValueOf(b.get(name)))
}

func (b *builder) getItems(args ...interface{}) {
	for len(args) >= 2 {
		name, ok := args[0].(string)
		if !ok {
			panic("string argument expected in builder.getItems()")
		}
		b.getItem(name, args[1])
		args = args[2:]
	}
}

func (b *builder) get(name string) glib.IObject {
	obj, err := b.GetObject(name)
	if err != nil {
		panic("builder.GetObject() failed: " + err.Error())
	}
	return obj
}

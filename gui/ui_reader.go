//go:generate esc -o definitions.go -modtime 1489449600 -pkg gui -include ".xml$" definitions/

package gui

import (
	"fmt"
	"os"
	"path"
	"strings"

	"reflect"

	"github.com/twstrike/gotk3adapter/glibi"
	"github.com/twstrike/gotk3adapter/gtki"
)

const (
	defsFolder   = "gui/definitions"
	xmlExtension = ".xml"
)

func getActualDefsFolder() string {
	wd, _ := os.Getwd()
	if strings.HasSuffix(wd, "/gui") {
		return "definitions"
	}
	return "gui/definitions"
}

func getDefinitionWithFileFallback(uiName string) string {
	fname := path.Join("/definitions", uiName+xmlExtension)
	embeddedFile, err := FSString(false, fname)
	if err != nil {
		panic(fmt.Sprintf("No definition found for %s", uiName))
	}

	if localFile, err := FSString(true, fname); err == nil {
		return localFile
	}

	return embeddedFile
}

// This must be called from the UI thread - otherwise bad things will happen sooner or later
func builderForDefinition(uiName string) gtki.Builder {
	template := getDefinitionWithFileFallback(uiName)

	builder, err := g.gtk.BuilderNew()
	if err != nil {
		//We cant recover from this
		panic(err)
	}

	//XXX Why are we using AddFromString rather than NewFromString
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

type builder struct {
	gtki.Builder
}

func newBuilder(filename string) *builder {
	return newBuilderFromString(filename)
}

func newBuilderFromString(uiName string) *builder {
	return &builder{builderForDefinition(uiName)}
}

func (b *builder) getObj(name string) glibi.Object {
	obj, _ := b.GetObject(name)
	return obj
}

func (b *builder) getItem(name string, target interface{}) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		panic("builder.getItem() target argument must be a pointer")
	}
	elem := v.Elem()
	elem.Set(reflect.ValueOf(b.get(name)))
}

//TODO: Why not a map[string]interface{}?
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

func (b *builder) get(name string) glibi.Object {
	obj, err := b.GetObject(name)
	if err != nil {
		panic("builder.GetObject() failed: " + err.Error())
	}
	return obj
}

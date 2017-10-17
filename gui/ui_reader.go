//go:generate esc -o definitions.go -modtime 1489449600 -pkg gui -ignore "Makefile" definitions/

package gui

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"reflect"

	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

const (
	xmlExtension = ".xml"
	imagesFolder = "/definitions/images/"
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

	//We dont use NewFromString because it doesnt give us an error message
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

func mustGetImageBytes(filename string) []byte {
	bs, err := FSByte(false, imagesFolder+filename)
	if err != nil {
		panic("Developer error: getting the image " + filename + " but it does not exist")
	}
	return bs
}

func setImageFromFile(i gtki.Image, filename string) {
	pl, err := g.gdk.PixbufLoaderNew()
	if err != nil {
		panic("Developer error: setting the image from " + filename)
	}

	var w sync.WaitGroup
	w.Add(1)
	pl.Connect("area-prepared", w.Done)

	if _, err := pl.Write(mustGetImageBytes(filename)); err != nil {
		log.Println(">> WARN - cannot write to PixbufLoader: " + err.Error())
		return
	}
	if err := pl.Close(); err != nil {
		log.Println(">> WARN - cannot close PixbufLoader: " + err.Error())
		return
	}

	w.Wait() //Waiting for Pixbuf to load before using it

	pb, err := pl.GetPixbuf()
	if err != nil {
		log.Println(">> WARN - cannot write to PixbufLoader: " + err.Error())
		return
	}
	i.SetFromPixbuf(pb)
	return
}

package settings

import (
	"fmt"
	"os"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"
	"github.com/twstrike/coyim/gui/settings/definitions"
)

var g glibi.Glib

// InitSettings should be called before using settings
func InitSettings(gx glibi.Glib) {
	g = gx
}

// TODO: Create a parent with default settings without the default config id to allow setting real defaults for eg SG

var cachedSchema glibi.SettingsSchemaSource

func getSchemaSource() glibi.SettingsSchemaSource {
	if cachedSchema == nil {
		dir := definitions.SchemaInTempDir()
		defer os.Remove(dir)
		fmt.Printf("using directory: %s\n", dir)
		cachedSchema = g.SettingsSchemaSourceNewFromDirectory(dir, nil, true)
	}

	return cachedSchema
}

func getSchema() glibi.SettingsSchema {
	return getSchemaSource().Lookup("im.coy.coyim.MainSettings", false)
}

func getSettingsFor(s string) glibi.Settings {
	return g.SettingsNewFull(getSchema(), nil, fmt.Sprintf("/im/coy/coyim/%s/", s))
}

func getDefaultSettings() glibi.Settings {
	return g.SettingsNewFull(getSchema(), nil, "/im/coy/coyim/")
}

// func RunTest() {
// 	before1 := getSettingsFor("foo1").GetString("hello")
// 	getSettingsFor("foo1").SetString("hello", "goodbye")
// 	after1 := getSettingsFor("foo1").GetString("hello")

// 	before2 := getDefaultSettings().GetString("hello")
// 	getDefaultSettings().SetString("hello", "somewhere")
// 	after2 := getDefaultSettings().GetString("hello")

// 	fmt.Printf("before1: %s\n", before1)
// 	fmt.Printf("after1: %s\n", after1)
// 	fmt.Printf("before2: %s\n", before2)
// 	fmt.Printf("after2: %s\n", after2)
// }

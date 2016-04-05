package settings

import (
	"fmt"
	"os"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/glib"
	"github.com/twstrike/coyim/gui/settings/definitions"
)

// TODO: Create a parent with default settings without the default config id to allow setting real defaults for eg SG

var cachedSchema *glib.SettingsSchemaSource

func getSchemaSource() *glib.SettingsSchemaSource {
	if cachedSchema == nil {
		dir := definitions.SchemaInTempDir()
		defer os.Remove(dir)
		fmt.Printf("using directory: %s\n", dir)
		cachedSchema = glib.SettingsSchemaSourceNewFromDirectory(dir, nil, true)
	}

	return cachedSchema
}

func getSchema() *glib.SettingsSchema {
	return getSchemaSource().Lookup("im.coy.coyim.MainSettings", false)
}

func getSettingsFor(s string) *glib.Settings {
	return glib.SettingsNewFull(getSchema(), nil, fmt.Sprintf("/im/coy/coyim/%s/", s))
}

func getDefaultSettings() *glib.Settings {
	return glib.SettingsNewFull(getSchema(), nil, "/im/coy/coyim/")
}

func RunTest() {
	before1 := getSettingsFor("foo1").GetString("hello")
	getSettingsFor("foo1").SetString("hello", "goodbye")
	after1 := getSettingsFor("foo1").GetString("hello")

	before2 := getDefaultSettings().GetString("hello")
	getDefaultSettings().SetString("hello", "somewhere")
	after2 := getDefaultSettings().GetString("hello")

	fmt.Printf("before1: %s\n", before1)
	fmt.Printf("after1: %s\n", after1)
	fmt.Printf("before2: %s\n", before2)
	fmt.Printf("after2: %s\n", after2)
}

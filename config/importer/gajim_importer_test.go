// +build !windows

package importer

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type gajimImporterPathTest struct {
	origConfigHomePath string
	origDataHomePath   string
	tempPath           string
}

func newGajimImportPathsForTesting() *gajimImporterPathTest {
	dir, _ := ioutil.TempDir("", "")

	gi := &gajimImporterPathTest{
		origConfigHomePath: os.Getenv("XDG_CONFIG_HOME"),
		origDataHomePath:   os.Getenv("XDG_DATA_HOME"),
		tempPath:           dir,
	}

	os.Setenv("XDG_CONFIG_HOME", dir)
	os.Setenv("XDG_DATA_HOME", dir)

	os.MkdirAll(dir, 0755)
	os.MkdirAll(filepath.Join(gi.dir(), "config"), 0755)
	os.MkdirAll(filepath.Join(gi.dir(), "pluginsconfig"), 0755)

	return gi
}

func (gi *gajimImporterPathTest) dir() string {
	return filepath.Join(gi.tempPath, "gajim")
}

// This method should be called at the end of the test execution in
// order to restore the original paths
func (gi *gajimImporterPathTest) restorePaths() {
	os.Setenv("XDG_CONFIG_HOME", gi.origConfigHomePath)
	os.Setenv("XDG_DATA_HOME", gi.origDataHomePath)
	os.RemoveAll(gi.tempPath)
}

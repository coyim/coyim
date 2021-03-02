package importer

import (
	"path/filepath"

	"github.com/coyim/coyim/config"
)

func gajimGetConfigAndDataDirs() (configRoot, dataRoot string) {
	configRoot = filepath.Join(config.SystemConfigDir(), "Gajim")
	dataRoot = configRoot
	return
}

// +build !windows

package importer

import (
	"path/filepath"

	"github.com/coyim/coyim/config"
)

func gajimGetConfigAndDataDirs() (configRoot, dataRoot string) {
	configRoot = filepath.Join(config.SystemConfigDir(), "gajim")
	dataRoot = filepath.Join(config.XdgDataHome(), "gajim")
	return
}

// +build !windows

package importer

import (
	"os"
)

func (s *GajimTryImportSuite) setAppDataHome() {
	s.origEnv1 = os.Getenv("XDG_CONFIG_HOME")
	s.origEnv2 = os.Getenv("XDG_DATA_HOME")

	os.Setenv("XDG_CONFIG_HOME", s.tempPath)
	os.Setenv("XDG_DATA_HOME", s.tempPath)
}

func (s *GajimTryImportSuite) restoreAppDataHome() {
	os.Setenv("XDG_CONFIG_HOME", s.origEnv1)
	os.Setenv("XDG_DATA_HOME", s.origEnv2)
}

func (s *GajimTryImportSuite) appDirName() string {
	return "gajim"
}

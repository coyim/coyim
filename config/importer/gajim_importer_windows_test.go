// +build windows

package importer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	csidlApdata = 0x1a
)

func getWindowsCurrentAppDataPath() string {
	shell32 := syscall.NewLazyDLL("shell32.dll")
	procShGetFolderPath := shell32.NewProc("SHGetFolderPathW")

	b := make([]uint16, syscall.MAX_PATH)
	ret, _, err := syscall.Syscall6(procShGetFolderPath.Addr(), 5, 0, csidlApdata, 0, 0, uintptr(unsafe.Pointer(&b[0])), 0)
	if int(ret) != 0 {
		panic(fmt.Sprintf("SHGetFolderPathW : err %d", int(err)))
	}

	return syscall.UTF16ToString(b)
}

func setWindowsCurrentAppDataPath(path string) {
	shell32 := syscall.NewLazyDLL("shell32.dll")
	procShSetFolderPath := shell32.NewProc("SHSetFolderPathW")

	b, _ := syscall.UTF16PtrFromString(path)
	procShSetFolderPath.Call(uintptr(csidlApdata), 0, 0, uintptr(unsafe.Pointer(b)))
}

type gajimImporterPathTest struct {
	origPath string
	tempPath string
}

func newGajimImportPathsForTesting() *gajimImporterPathTest {
	dir, _ := ioutil.TempDir("", "")

	gi := &gajimImporterPathTest{
		origPath: getWindowsCurrentAppDataPath(),
		tempPath: dir,
	}

	setWindowsCurrentAppDataPath(dir)

	os.MkdirAll(dir, 0755)
	os.MkdirAll(filepath.Join(gi.dir(), "config"), 0755)
	os.MkdirAll(filepath.Join(gi.dir(), "pluginsconfig"), 0755)

	return gi
}

func (gi *gajimImporterPathTest) dir() string {
	return filepath.Join(gi.tempPath, "Gajim")
}

// This method should be called at the end of the test execution in
// order to restore the original paths
func (gi *gajimImporterPathTest) restorePaths() {
	setWindowsCurrentAppDataPath(gi.origPath)
	os.RemoveAll(gi.tempPath)
}

package importer

import (
	"fmt"
	"syscall"
	"unsafe"
)

func (s *GajimSuite) setAppDataDirOSDependent() {
	s.origEnv1 = getWindowsCurrentAppDataPath()
	setWindowsCurrentAppDataPath(s.tempPath)
}

func (s *GajimSuite) restoreAppDataDirOSDependent() {
	setWindowsCurrentAppDataPath(s.origEnv1)
}

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

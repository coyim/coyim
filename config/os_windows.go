package config

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	shell32           = syscall.NewLazyDLL("shell32.dll")
	procGetFolderPath = shell32.NewProc("SHGetFolderPathW")
)

const (
	csidlAppdata = 0x1a
)

func appdataFolderPath() string {
	b := make([]uint16, syscall.MAX_PATH)
	ret, _, err := syscall.Syscall6(procGetFolderPath.Addr(), 5, 0, csidlAppdata, 0, 0, uintptr(unsafe.Pointer(&b[0])), 0)
	if int(ret) != 0 {
		panic(fmt.Sprintf("SHGetFolderPathW : err %d", int(err)))
	}
	return syscall.UTF16ToString(b)
}

// IsWindows returns true if this is running under windows
func IsWindows() bool {
	return true
}

// SystemConfigDir points to the function that gets the configuration directory for this system
var SystemConfigDir = appdataFolderPath

// // SystemConfigDir returns the application data directory, valid on both windows and posix systems
// func SystemConfigDir() string {
// 	return appdataFolderPath()
// }

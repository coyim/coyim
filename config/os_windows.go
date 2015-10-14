package config

import (
	"fmt"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	shell32                 = syscall.NewLazyDLL("shell32.dll")
	proc_sh_get_folder_path = shell32.NewProc("SHGetFolderPathW")
)

const (
	csidl_appdata = 0x1a
)

func appdataFolderPath() string {
	b := make([]uint16, syscall.MAX_PATH)
	ret, _, err := syscall.Syscall6(proc_sh_get_folder_path.Addr(), 5, 0, csidl_appdata, 0, 0, uintptr(unsafe.Pointer(&b[0])), 0)
	if int(ret) != 0 {
		panic(fmt.Sprintf("SHGetFolderPathW : err %d", int(err)))
	}
	return syscall.UTF16ToString(b)
}

func configDir() string {
	return filepath.Join(appdataFolderPath(), "coyim")
}

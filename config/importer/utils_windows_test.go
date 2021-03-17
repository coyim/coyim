// +build windows

package importer

import (
	"fmt"
	"os"
	"os/exec"
)

func makeDirectoryUnnaccesibleOSDependent(dir string) {
	denyWindowsUserDirPermissions(dir)
}

func makeDirectoryAccesibleOSDependent(dir string) {
	grantWindowsUserDirPermissions(dir)
}

func denyWindowsUserDirPermissions(dir string) {
	icaclsExec(dir, "/deny", fmt.Sprintf("%s:(RX,W)", os.Getenv("username")))
}

func grantWindowsUserDirPermissions(dir string) {
	icaclsExec(dir, "/grant", fmt.Sprintf("%s:(RX,W)", os.Getenv("username")))
}

func icaclsExec(dir, action, permissions string) {
	c := exec.Command("icacls", dir, action, permissions)
	c.Run()
}

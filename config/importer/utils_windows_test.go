// +build windows

package importer

import (
	"fmt"
	"os"
	"os/exec"
)

func makeDirectoryUnnaccesibleOSDependent(dir string) {
	done := make(chan bool)
	go denyWindowsUserDirPermissions(dir, done)
	<-done
}

func makeDirectoryAccesibleOSDependent(dir string) {
	done := make(chan bool)
	go grantWindowsUserDirPermissions(dir, done)
	<-done
}

func denyWindowsUserDirPermissions(dir string, done chan bool) {
	icaclsExec(dir, "/deny", fmt.Sprintf("%s:(RX,W)", os.Getenv("username")), done)
}

func grantWindowsUserDirPermissions(dir string, done chan bool) {
	icaclsExec(dir, "/grant", fmt.Sprintf("%s:(RX,W)", os.Getenv("username")), done)
}

func icaclsExec(dir, action, permissions string, done chan bool) {
	c := exec.Command("icacls", dir, action, permissions)
	c.Run()
	done <- true
}

// +build !windows

package importer

import "os"

func makeDirectoryUnnaccesibleOSDependent(dir string) {
	os.Chmod(dir, 0000)
}

func makeDirectoryAccesibleOSDependent(dir string) {
	os.Chmod(dir, 0755)
}

// +build !windows

package importer

import "os"

func makeDirectoryInaccessible(dir string) {
	os.Chmod(dir, 0000)
}

func makeDirectoryAccessible(dir string) {
	os.Chmod(dir, 0755)
}

// +build !windows

package importer

import "os"

func makeDirectoryUnnaccesible(dir string) {
	os.Chmod(dir, 0000)
}

func makeDirectoryAccesible(dir string) {
	os.Chmod(dir, 0755)
}

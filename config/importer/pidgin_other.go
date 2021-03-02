// +build !windows

package importer

func findDirOSDependent() (string, bool) {
	return "", false
}

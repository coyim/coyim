package importer

import (
	"path"
	"runtime"
)

func testResourceFilename(s string) string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Join(path.Dir(filename), s)
}

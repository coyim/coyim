package css

import (
	"embed"
	"io/fs"
	"path"

	"github.com/sirupsen/logrus"
)

//go:embed definitions
var files embed.FS

// Get will return the CSS string corresponding to the name given
func Get(name string) string {
	content, e := fs.ReadFile(files, path.Join("definitions", name))
	if e != nil {
		logrus.WithError(e).WithField("definition", name).Panic("No definition found")
	}
	return string(content)
}

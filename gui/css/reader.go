//go:generate ../../.build-tools/esc -o definitions.go -private -modtime 1489449600 -pkg css -ignore "Makefile" definitions/

package css

import "path"

// Get will return the CSS string corresponding to the name given
func Get(name string) string {
	fname := path.Join("/definitions", name)
	return _escFSMustString(false, fname)
}

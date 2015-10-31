package importer

import (
	"os"

	plist "github.com/DHowett/go-plist"
)

type adiumImporter struct{}

func (ai *adiumImporter) read() {
	f, _ := os.Open("test.txt")
	defer f.Close()

	dec := plist.NewDecoder(f)

	var res map[string]interface{}
	dec.Decode(&res)
}

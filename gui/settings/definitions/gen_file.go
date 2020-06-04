package definitions

import (
	"encoding/hex"
	"io/ioutil"
	"path"
)

func fileContent() []byte {
	decoded, _ := hex.DecodeString(schemaDefinition)
	return decoded
}

func writeSchemaToDir(dir string) {
	ioutil.WriteFile(path.Join(dir, "gschemas.compiled"), fileContent(), 0600)
}

// SchemaInTempDir will create a new temporary directory and put the gsettings schema file in there. It is the callers responsibility to remove the directory
func SchemaInTempDir() string {
	dir, _ := ioutil.TempDir("", "coyim-schema")
	writeSchemaToDir(dir)
	return dir
}

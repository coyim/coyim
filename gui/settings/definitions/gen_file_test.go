package definitions

import (
	"io/ioutil"
	"os"
	"path"

	. "gopkg.in/check.v1"
)

type SettingsDefinitionsSuite struct{}

var _ = Suite(&SettingsDefinitionsSuite{})

func (s *SettingsDefinitionsSuite) Test_fileContent_returnsData(c *C) {
	res := fileContent()
	c.Assert(res[0:10], DeepEquals, []byte{0x47, 0x56, 0x61, 0x72, 0x69, 0x61, 0x6e, 0x74, 0x0, 0x0})
}

func (s *SettingsDefinitionsSuite) Test_SchemaInTempDir(c *C) {
	res := SchemaInTempDir()
	fi, err := os.Stat(res)
	c.Assert(err, IsNil)
	c.Assert(fi.IsDir(), Equals, true)

	data, _ := ioutil.ReadFile(path.Join(res, "gschemas.compiled"))

	c.Assert(data, DeepEquals, fileContent())
}

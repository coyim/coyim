package filetransfer

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"
)

type PackagingSuite struct{}

var _ = Suite(&PackagingSuite{})

func createTemporaryDirectoryStructure(intoDir string) {
	_ = os.MkdirAll(filepath.Join(intoDir, "foo", "bar", "quux"), 0755)
	_ = os.MkdirAll(filepath.Join(intoDir, "foo", "somewhere"), 0755)
	_ = os.MkdirAll(filepath.Join(intoDir, "another place"), 0755)

	_ = ioutil.WriteFile(filepath.Join(intoDir, "a file"), []byte("some content"), 0755)
	_ = ioutil.WriteFile(filepath.Join(intoDir, "another file"), []byte("some more content"), 0755)
	_ = ioutil.WriteFile(filepath.Join(intoDir, "foo", "wut"), []byte("even more content"), 0755)
	_ = ioutil.WriteFile(filepath.Join(intoDir, "another place", "oh well"), []byte("final content"), 0755)
}

func (s *PackagingSuite) Test_pack_and_unpack(c *C) {
	dd := c.MkDir()
	createTemporaryDirectoryStructure(dd)

	f, ex := ioutil.TempFile("", "")
	c.Assert(ex, IsNil)
	defer func() {
		ex2 := os.Remove(f.Name())
		c.Assert(ex2, IsNil)
	}()

	e := pack(dd, f)
	ex = f.Close()
	c.Assert(ex, IsNil)

	c.Assert(e, IsNil)

	dd2 := c.MkDir()
	e2 := unpack(f.Name(), dd2)
	c.Assert(e2, IsNil)

	actualDir := filepath.Join(dd2, filepath.Base(dd))

	st, ee := os.Stat(filepath.Join(actualDir, "foo"))
	c.Assert(ee, IsNil)
	c.Assert(st, Not(IsNil))
	c.Assert(st.IsDir(), Equals, true)

	st, ee = os.Stat(filepath.Join(actualDir, "foo", "bar"))
	c.Assert(ee, IsNil)
	c.Assert(st, Not(IsNil))
	c.Assert(st.IsDir(), Equals, true)

	st, ee = os.Stat(filepath.Join(actualDir, "foo", "bar", "quux"))
	c.Assert(ee, IsNil)
	c.Assert(st, Not(IsNil))
	c.Assert(st.IsDir(), Equals, true)

	st, ee = os.Stat(filepath.Join(actualDir, "foo", "somewhere"))
	c.Assert(ee, IsNil)
	c.Assert(st, Not(IsNil))
	c.Assert(st.IsDir(), Equals, true)

	st, ee = os.Stat(filepath.Join(actualDir, "another place"))
	c.Assert(ee, IsNil)
	c.Assert(st, Not(IsNil))
	c.Assert(st.IsDir(), Equals, true)

	ct, ee2 := ioutil.ReadFile(filepath.Join(actualDir, "a file"))
	c.Assert(ee2, IsNil)
	c.Assert(string(ct), Equals, "some content")

	ct, ee2 = ioutil.ReadFile(filepath.Join(actualDir, "another file"))
	c.Assert(ee2, IsNil)
	c.Assert(string(ct), Equals, "some more content")

	ct, ee2 = ioutil.ReadFile(filepath.Join(actualDir, "foo", "wut"))
	c.Assert(ee2, IsNil)
	c.Assert(string(ct), Equals, "even more content")

	ct, ee2 = ioutil.ReadFile(filepath.Join(actualDir, "another place", "oh well"))
	c.Assert(ee2, IsNil)
	c.Assert(string(ct), Equals, "final content")
}

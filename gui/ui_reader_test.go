package gui

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"time"

	"github.com/coyim/gotk3adapter/gtk_mock"
	"github.com/coyim/gotk3adapter/gtki"

	. "gopkg.in/check.v1"
)

type UIReaderSuite struct{}

var _ = Suite(&UIReaderSuite{})

const testFile string = `
<interface>
  <object class="GtkWindow" id="conversation">
    <property name="default-height">600</property>
    <property name="default-width">500</property>
    <child>
	  <object class="GtkBox" id="vbox">
	    <property name="orientation">GTK_ORIENTATION_VERTICAL</property>  
	  </object>
    </child>
  </object>
</interface>
`

func writeTestFile(name, content string) {
	_ = ioutil.WriteFile(name, []byte(content), 0700)
}

func removeFile(name string) {
	_ = os.Remove(name)
}

type mockBuilder struct {
	gtk_mock.MockBuilder
	stringGiven string
	errorGiven  error
}

func (v *mockBuilder) AddFromString(v1 string) error {
	v.stringGiven = v1
	return v.errorGiven
}

type mockWithBuilder struct {
	gtk_mock.Mock
	errorGiven       error
	secondErrorGiven error
}

func (v *mockWithBuilder) BuilderNew() (gtki.Builder, error) {
	return &mockBuilder{
		errorGiven: v.secondErrorGiven,
	}, v.errorGiven
}

const wrongTemplate string = `
yeah
<interfae>
  <object class="GtkWindow">
    I have a bad format
  </object>
`

func (s *UIReaderSuite) Test_builderForString_panicsIfNotBuilder(c *C) {
	g = Graphics{gtk: &mockWithBuilder{
		errorGiven: errors.New("bla"),
	}}

	c.Assert(func() {
		builderForString(testFile)
	}, PanicMatches, "bla")
}

func (s *UIReaderSuite) Test_builderForString_panicsIfEmptyTemplate(c *C) {
	g = Graphics{gtk: &mockWithBuilder{
		secondErrorGiven: errors.New("foo"),
	}}

	c.Assert(func() {
		builderForString("")
	}, PanicMatches, "gui: wrong template format: foo\n")
}

func (s *UIReaderSuite) Test_builderForString_panicsIfWrongTemplate(c *C) {
	g = Graphics{gtk: &mockWithBuilder{
		secondErrorGiven: errors.New("bla"),
	}}

	c.Assert(func() {
		builderForString(wrongTemplate)
	}, PanicMatches, "gui: wrong template format: bla\n")
}

func (s *UIReaderSuite) Test_builderForString_useTemplateStringIfOk(c *C) {
	g = Graphics{gtk: &mockWithBuilder{}}

	builder := builderForString(testFile)

	str := builder.(*mockBuilder).stringGiven

	c.Assert(str, Equals, testFile)
}

func (s *UIReaderSuite) Test_builderForDefinition_useXMLIfExists(c *C) {
	orgMod := getModTime(getActualDefsFolder() + "/Test.xml")
	defer func() {
		setModTime(getActualDefsFolder()+"/Test.xml", orgMod)
	}()

	g = Graphics{gtk: &mockWithBuilder{}}
	removeFile(getActualDefsFolder() + "/Test.xml")
	writeTestFile(getActualDefsFolder()+"/Test.xml", testFile)
	ui := "Test"

	builder := builderForDefinition(ui)

	str := builder.(*mockBuilder).stringGiven

	c.Assert(str, Equals, testFile)
}

func getModTime(fn string) time.Time {
	file, _ := os.Stat(fn)
	return file.ModTime()
}

func setModTime(fn string, t time.Time) {
	_ = os.Chtimes(fn, t, t)
}

func (s *UIReaderSuite) Test_builderForDefinition_useGoFileIfXMLDoesntExists(c *C) {
	orgMod := getModTime(getActualDefsFolder() + "/Test.xml")
	defer func() {
		writeTestFile(getActualDefsFolder()+"/Test.xml", testFile)
		setModTime(getActualDefsFolder()+"/Test.xml", orgMod)
	}()

	g = Graphics{gtk: &mockWithBuilder{}}
	removeFile(getActualDefsFolder() + "/Test.xml")
	ui := "Test"

	builder := builderForDefinition(ui)

	str := builder.(*mockBuilder).stringGiven

	c.Assert(str, Equals, testFile)
}

func (s *UIReaderSuite) Test_builderForDefinition_shouldReturnErrorWhenDefinitionDoesntExist(c *C) {
	ui := "nonexistent"

	c.Assert(func() {
		builderForDefinition(ui)
	}, Panics, "No definition found for nonexistent")
}

func (s *UIReaderSuite) Test_getImageBytes_forExistingImage(c *C) {
	r := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(
		`
PHN2ZyB3aWR0aD0iMTUiIGhlaWdodD0iMTQiIHZpZXdCb3g9IjAgMCAxNSAxNCIgZmlsbD0ibm9uZSIg
eG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICAgIDxwYXRoIGZpbGwtcnVsZT0iZXZl
bm9kZCIgY2xpcC1ydWxlPSJldmVub2RkIiBkPSJNOC41NzAyMiAxMS40OTY1VjkuOTA3ODdDOC41NzAy
MiA5LjgyOTgzIDguNTQzNzMgOS43NjQzNCA4LjQ5MDc4IDkuNzExMzlDOC40Mzc4MyA5LjY1ODQzIDgu
Mzc1MTIgOS42MzE5NiA4LjMwMjY2IDkuNjMxOTZINi42OTczNEM2LjYyNDg4IDkuNjMxOTYgNi41NjIx
NyA5LjY1ODQzIDYuNTA5MjEgOS43MTEzOUM2LjQ1NjI2IDkuNzY0MzQgNi40Mjk3OCA5LjgyOTgzIDYu
NDI5NzggOS45MDc4N1YxMS40OTY1QzYuNDI5NzggMTEuNTc0NSA2LjQ1NjI2IDExLjY0IDYuNTA5MjEg
MTEuNjkzQzYuNTYyMTcgMTEuNzQ1OSA2LjYyNDg4IDExLjc3MjQgNi42OTczNCAxMS43NzI0SDguMzAy
NjZDOC4zNzUxMiAxMS43NzI0IDguNDM3ODMgMTEuNzQ1OSA4LjQ5MDc4IDExLjY5M0M4LjU0MzczIDEx
LjY0IDguNTcwMjIgMTEuNTc0NSA4LjU3MDIyIDExLjQ5NjVaTTguNTUzNDkgOC4zNjk0M0w4LjcwMzk5
IDQuNTMxN0M4LjcwMzk5IDQuNDY0ODEgOC42NzYxMSA0LjQxMTg2IDguNjIwMzcgNC4zNzI4NEM4LjU0
NzkxIDQuMzExNTMgOC40ODEwNCA0LjI4MDg3IDguNDE5NzIgNC4yODA4N0g2LjU4MDI3QzYuNTE4OTYg
NC4yODA4NyA2LjQ1MjA4IDQuMzExNTMgNi4zNzk2MiA0LjM3Mjg0QzYuMzIzODggNC40MTE4NiA2LjI5
NiA0LjQ3MDM5IDYuMjk2IDQuNTQ4NDJMNi40MzgxNCA4LjM2OTQzQzYuNDM4MTQgOC40MjUxOCA2LjQ2
NjAyIDguNDcxMTYgNi41MjE3NiA4LjUwNzM5QzYuNTc3NSA4LjU0MzYyIDYuNjQ0MzkgOC41NjE3NCA2
LjcyMjQyIDguNTYxNzRIOC4yNjkyMkM4LjM0NzI2IDguNTYxNzQgOC40MTI3NCA4LjU0MzYyIDguNDY1
NyA4LjUwNzM5QzguNTE4NjUgOC40NzExNiA4LjU0NzkyIDguNDI1MTggOC41NTM0OSA4LjM2OTQzWk04
LjQzNjQ0IDAuNTYwMTkyTDE0Ljg1NzcgMTIuMzMyNkMxNS4wNTI4IDEyLjY4MzggMTUuMDQ3MyAxMy4w
MzQ5IDE0Ljg0MSAxMy4zODYxQzE0Ljc0NjMgMTMuNTQ3NyAxNC42MTY3IDEzLjY3NTkgMTQuNDUyMiAx
My43NzA3QzE0LjI4NzggMTMuODY1NCAxNC4xMTA4IDEzLjkxMjggMTMuOTIxMyAxMy45MTI4SDEuMDc4
NjlDMC44ODkxNjggMTMuOTEyOCAwLjcxMjE5MiAxMy44NjU0IDAuNTQ3NzU3IDEzLjc3MDdDMC4zODMz
MjIgMTMuNjc1OSAwLjI1MzczOCAxMy41NDc3IDAuMTU4OTc4IDEzLjM4NjFDLTAuMDQ3MjYyOCAxMy4w
MzQ5IC0wLjA1MjgzODQgMTIuNjgzOCAwLjE0MjI1NSAxMi4zMzI2TDYuNTYzNTUgMC41NjAxOTJDNi42
NTgzMSAwLjM4NzM5NiA2Ljc4OTMgMC4yNTA4MzMgNi45NTY1MiAwLjE1MDQ5OUM3LjEyMzc1IDAuMDUw
MTY1OSA3LjMwNDkgMCA3LjUgMEM3LjY5NTA5IDAgNy44NzYyNSAwLjA1MDE2NTkgOC4wNDM0NyAwLjE1
MDQ5OUM4LjIxMDY5IDAuMjUwODMzIDguMzQxNjkgMC4zODczOTYgOC40MzY0NCAwLjU2MDE5MloiIGZp
bGw9IiNGODlCMUMiLz4KPC9zdmc+Cg==`))
	expectedBytes, _ := ioutil.ReadAll(r)

	c.Assert(expectedBytes, DeepEquals, mustGetImageBytes("alert.svg"))

	r = base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(
		`
PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiIHN0YW5kYWxvbmU9Im5vIj8+Cjxzdmcg
d2lkdGg9IjExcHgiIGhlaWdodD0iMTNweCIgdmlld0JveD0iMCAwIDExIDEzIiB2ZXJzaW9uPSIxLjEi
IHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cu
dzMub3JnLzE5OTkveGxpbmsiPgogICAgPCEtLSBHZW5lcmF0b3I6IFNrZXRjaCA0MS4yICgzNTM5Nykg
LSBodHRwOi8vd3d3LmJvaGVtaWFuY29kaW5nLmNvbS9za2V0Y2ggLS0+CiAgICA8dGl0bGU+R3JvdXA8
L3RpdGxlPgogICAgPGRlc2M+Q3JlYXRlZCB3aXRoIFNrZXRjaC48L2Rlc2M+CiAgICA8ZGVmcz48L2Rl
ZnM+CiAgICA8ZyBpZD0iU01QLWZsb3dzIiBzdHJva2U9Im5vbmUiIHN0cm9rZS13aWR0aD0iMSIgZmls
bD0ibm9uZSIgZmlsbC1ydWxlPSJldmVub2RkIj4KICAgICAgICA8ZyBpZD0iQm9iXzQiIHRyYW5zZm9y
bT0idHJhbnNsYXRlKC04NDMuMDAwMDAwLCAtMTY0LjAwMDAwMCkiPgogICAgICAgICAgICA8ZyBpZD0i
R3JvdXAiIHRyYW5zZm9ybT0idHJhbnNsYXRlKDg0My4wMDAwMDAsIDE2NC4wMDAwMDApIj4KICAgICAg
ICAgICAgICAgIDxwYXRoIGQ9Ik0zLjA1NTU1NTU2LDUuOTA5MDkwOTEgTDcuOTQ0NDQ0NDQsNS45MDkw
OTA5MSBMNy45NDQ0NDQ0NCw0LjEzNjM2MzY0IEM3Ljk0NDQ0NDQ0LDMuNDgzODk4MjUgNy43MDU3MzE1
NSwyLjkyNjg0ODkgNy4yMjgyOTg2MSwyLjQ2NTE5ODg2IEM2Ljc1MDg2NTY3LDIuMDAzNTQ4ODMgNi4x
NzQ3NzE4OSwxLjc3MjcyNzI3IDUuNSwxLjc3MjcyNzI3IEM0LjgyNTIyODExLDEuNzcyNzI3MjcgNC4y
NDkxMzQzMywyLjAwMzU0ODgzIDMuNzcxNzAxMzksMi40NjUxOTg4NiBDMy4yOTQyNjg0NSwyLjkyNjg0
ODkgMy4wNTU1NTU1NiwzLjQ4Mzg5ODI1IDMuMDU1NTU1NTYsNC4xMzYzNjM2NCBMMy4wNTU1NTU1Niw1
LjkwOTA5MDkxIFogTTExLDYuNzk1NDU0NTUgTDExLDEyLjExMzYzNjQgQzExLDEyLjM1OTg0OTcgMTAu
OTEwODgwNSwxMi41NjkxMjc5IDEwLjczMjYzODksMTIuNzQxNDc3MyBDMTAuNTU0Mzk3MywxMi45MTM4
MjY2IDEwLjMzNzk2NDIsMTMgMTAuMDgzMzMzMywxMyBMMC45MTY2NjY2NjcsMTMgQzAuNjYyMDM1NzY0
LDEzIDAuNDQ1NjAyNzQzLDEyLjkxMzgyNjYgMC4yNjczNjExMTEsMTIuNzQxNDc3MyBDMC4wODkxMTk0
NzkyLDEyLjU2OTEyNzkgMCwxMi4zNTk4NDk3IDAsMTIuMTEzNjM2NCBMMCw2Ljc5NTQ1NDU1IEMwLDYu
NTQ5MjQxMTkgMC4wODkxMTk0NzkyLDYuMzM5OTYyOTggMC4yNjczNjExMTEsNi4xNjc2MTM2NCBDMC40
NDU2MDI3NDMsNS45OTUyNjQyOSAwLjY2MjAzNTc2NCw1LjkwOTA5MDkxIDAuOTE2NjY2NjY3LDUuOTA5
MDkwOTEgTDEuMjIyMjIyMjIsNS45MDkwOTA5MSBMMS4yMjIyMjIyMiw0LjEzNjM2MzY0IEMxLjIyMjIy
MjIyLDMuMDAzNzgyMjIgMS42NDIzNTY5MSwyLjAzMTI1NDA2IDIuNDgyNjM4ODksMS4yMTg3NSBDMy4z
MjI5MjA4NywwLjQwNjI0NTkzNyA0LjMyODY5Nzg1LDAgNS41LDAgQzYuNjcxMzAyMTUsMCA3LjY3NzA3
OTEzLDAuNDA2MjQ1OTM3IDguNTE3MzYxMTEsMS4yMTg3NSBDOS4zNTc2NDMwOSwyLjAzMTI1NDA2IDku
Nzc3Nzc3NzgsMy4wMDM3ODIyMiA5Ljc3Nzc3Nzc4LDQuMTM2MzYzNjQgTDkuNzc3Nzc3NzgsNS45MDkw
OTA5MSBMMTAuMDgzMzMzMyw1LjkwOTA5MDkxIEMxMC4zMzc5NjQyLDUuOTA5MDkwOTEgMTAuNTU0Mzk3
Myw1Ljk5NTI2NDI5IDEwLjczMjYzODksNi4xNjc2MTM2NCBDMTAuOTEwODgwNSw2LjMzOTk2Mjk4IDEx
LDYuNTQ5MjQxMTkgMTEsNi43OTU0NTQ1NSBMMTEsNi43OTU0NTQ1NSBaIiBpZD0i74CjLWNvcHktMiIg
ZmlsbD0iIzdFRDMyMSI+PC9wYXRoPgogICAgICAgICAgICAgICAgPHBhdGggZD0iTTguMzIzMjMxNzEs
OC4wMDQ5NDQ4NyBMNS4yMjA1NjI4MywxMS4wMDUwNDYyIEM1LjE1OTU3NjQxLDExLjA2NDAxNjUgNS4w
ODcxNTYxMSwxMS4wOTM1MDEzIDUuMDAzMjk5NzgsMTEuMDkzNTAxMyBDNC45MTk0NDM0NCwxMS4wOTM1
MDEzIDQuODQ3MDIzMTUsMTEuMDY0MDE2NSA0Ljc4NjAzNjcyLDExLjAwNTA0NjIgTDMuMTQ3MDM0NzQs
OS40MjAyMjYwOCBDMy4wODYwNDgzMSw5LjM2MTI1NTczIDMuMDU1NTU1NTYsOS4yOTEyMjk1IDMuMDU1
NTU1NTYsOS4yMTAxNDUyNyBDMy4wNTU1NTU1Niw5LjEyOTA2MTA1IDMuMDg2MDQ4MzEsOS4wNTkwMzQ4
MiAzLjE0NzAzNDc0LDkuMDAwMDY0NDcgTDMuNTY2MzE0MzEsOC41OTQ2NDUzNyBDMy42MjczMDA3NCw4
LjUzNTY3NTAzIDMuNjk5NzIxMDMsOC41MDYxOTAzIDMuNzgzNTc3MzcsOC41MDYxOTAzIEMzLjg2NzQz
MzcsOC41MDYxOTAzIDMuOTM5ODU0LDguNTM1Njc1MDMgNC4wMDA4NDA0Miw4LjU5NDY0NTM3IEw1LjAw
MzI5OTc4LDkuNTYzOTY1NTggTDcuNDY5NDI2MDIsNy4xNzkzNjQxNyBDNy41MzA0MTI0NSw3LjEyMDM5
MzgyIDcuNjAyODMyNzQsNy4wOTA5MDkwOSA3LjY4NjY4OTA4LDcuMDkwOTA5MDkgQzcuNzcwNTQ1NDEs
Ny4wOTA5MDkwOSA3Ljg0Mjk2NTcsNy4xMjAzOTM4MiA3LjkwMzk1MjEzLDcuMTc5MzY0MTcgTDguMzIz
MjMxNzEsNy41ODQ3ODMyNiBDOC4zODQyMTgxMyw3LjY0Mzc1MzYxIDguNDE0NzEwODksNy43MTM3Nzk4
NCA4LjQxNDcxMDg5LDcuNzk0ODY0MDcgQzguNDE0NzEwODksNy44NzU5NDgyOSA4LjM4NDIxODEzLDcu
OTQ1OTc0NTMgOC4zMjMyMzE3MSw4LjAwNDk0NDg3IEw4LjMyMzIzMTcxLDguMDA0OTQ0ODcgWiIgaWQ9
IlBhdGgiIGZpbGw9IiNGRkZGRkYiPjwvcGF0aD4KICAgICAgICAgICAgPC9nPgogICAgICAgIDwvZz4K
ICAgIDwvZz4KPC9zdmc+`))
	expectedBytes, _ = ioutil.ReadAll(r)

	c.Assert(expectedBytes, DeepEquals, mustGetImageBytes("padlock.svg"))
}

func (s *UIReaderSuite) Test_GettingNonExistantImage_Panics(c *C) {
	image := "nonexistent"

	c.Assert(func() {
		mustGetImageBytes(image)
	}, Panics, "Developer error: getting the image "+image+" but it does not exist")
}

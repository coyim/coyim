package gui

import (
	"io/ioutil"
	"os"

	"github.com/twstrike/gotk3adapter/gtk_mock"
	"github.com/twstrike/gotk3adapter/gtki"

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
      <object class="GtkVBox" id="vbox"></object>
    </child>
  </object>
</interface>
`

func writeTestFile(name, content string) {
	ioutil.WriteFile(name, []byte(content), 0700)
}

func removeFile(name string) {
	os.Remove(name)
}

type mockBuilder struct {
	gtk_mock.MockBuilder
	stringGiven string
}

func (v *mockBuilder) AddFromString(v1 string) error {
	v.stringGiven = v1
	return nil
}

type mockWithBuilder struct {
	gtk_mock.Mock
}

func (*mockWithBuilder) BuilderNew() (gtki.Builder, error) {
	return &mockBuilder{}, nil
}

func (s *UIReaderSuite) Test_builderForDefinition_useXMLIfExists(c *C) {
	g = Graphics{gtk: &mockWithBuilder{}}
	removeFile(getActualDefsFolder() + "/Test.xml")
	writeTestFile(getActualDefsFolder()+"/Test.xml", testFile)
	ui := "Test"

	builder := builderForDefinition(ui)

	str := builder.(*mockBuilder).stringGiven

	c.Assert(str, Equals, testFile)
}

func (s *UIReaderSuite) Test_builderForDefinition_useGoFileIfXMLDoesntExists(c *C) {
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

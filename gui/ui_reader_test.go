package gui

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"

	. "gopkg.in/check.v1"
)

type UIReaderSuite struct{}

var _ = Suite(&UIReaderSuite{})

const testFile string = `
<interface>
  <object class="GtkWindow" id="conversation">
    <property name="default-height">$win-height</property>
    <property name="default-width">$win-width</property>

      <child>
	<object class="GtkVBox" id="vbox">
	</object>
      </child>

  </object>
</interface>
`

func writeTestFile(name, content string) {
	desc, _ := os.Create(name)
	desc.WriteString(content)
}

func removeFile(name string) {
	os.Remove(name)
}

func (s *UIReaderSuite) Test_loadBuilderWith_useXMLIfExists(c *C) {
	gtk.Init(nil)
	removeFile("definitions/TestDefinition.xml")
	writeTestFile("definitions/TestDefinition.xml", testFile)
	ui := "TestDefinition"

	vars := make(map[string]string)
	vars["$win-height"] = "500"
	vars["$win-width"] = "400"
	builder, parseErr := loadBuilderWith(ui, vars)
	if parseErr != nil {
		fmt.Errorf("\nFailed!\n%s", parseErr.Error())
		c.Fail()
	}

	win, getErr := builder.GetObject("conversation")
	if getErr != nil {
		fmt.Errorf("\nFailed to get window \n%s", getErr.Error())
		c.Fail()
	}
	w, h := win.(*gtk.Window).GetSize()
	c.Assert(h, Equals, 500)
	c.Assert(w, Equals, 400)
}

func (s *UIReaderSuite) Test_loadBuilderWith_useGoFileIfXMLDoesntExists(c *C) {
	gtk.Init(nil)
	removeFile("definitions/TestDefinition.xml")
	//writeTestFile("definitions/TestDefinition.xml", testFile)
	ui := "TestDefinition"

	vars := make(map[string]string)
	vars["$win-height"] = "500"
	vars["$win-width"] = "400"
	builder, parseErr := loadBuilderWith(ui, vars)
	if parseErr != nil {
		fmt.Errorf("\nFailed!\n%s", parseErr.Error())
		c.Fail()
	}

	win, getErr := builder.GetObject("conversation")
	if getErr != nil {
		fmt.Errorf("\nFailed to get window \n%s", getErr.Error())
		c.Fail()
	}
	w, h := win.(*gtk.Window).GetSize()
	c.Assert(h, Equals, 500)
	c.Assert(w, Equals, 400)
}

func (s *UIReaderSuite) Test_loadBuilderWith_shouldReturnErrorWhenDefinitionDoesntExist(c *C) {
	removeFile("definitions/nonexistent")
	ui := "nonexistent"

	_, parseErr := loadBuilderWith(ui, nil)

	expected := "There's no definition for nonexistent"
	c.Assert(parseErr.Error(), Equals, expected)
}

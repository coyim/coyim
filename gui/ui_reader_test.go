package gui

import (
	"os"

	. "gopkg.in/check.v1"
)

type UIReaderSuite struct{}

var _ = Suite(&UIReaderSuite{})

const testFile string = `
<interface>
  <object class="GtkWindow" id="conversation">
    <property name="default-height">500</property>
    <property name="default-width">400</property>
    <child>
      <object class="GtkVBox" id="vbox"></object>
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

// func (s *UIReaderSuite) Test_builderForDefinition_useXMLIfExists(c *C) {
// 	g.gtk.Init(nil)
// 	removeFile("definitions/Test.xml")
// 	writeTestFile("definitions/Test.xml", testFile)
// 	ui := "Test"

// 	builder := builderForDefinition(ui)

// 	win, getErr := builder.GetObject("conversation")
// 	if getErr != nil {
// 		fmt.Errorf("\nFailed to get window \n%s", getErr.Error())
// 		c.Fail()
// 	}
// 	w, h := win.(gtki.Window).GetSize()
// 	c.Assert(h, Equals, 500)
// 	c.Assert(w, Equals, 400)
// }

// func (s *UIReaderSuite) Test_builderForDefinition_useGoFileIfXMLDoesntExists(c *C) {
// 	g.gtk.Init(nil)
// 	removeFile("definitions/Test.xml")
// 	//writeTestFile("definitions/TestDefinition.xml", testFile)
// 	ui := "Test"

// 	builder := builderForDefinition(ui)

// 	win, getErr := builder.GetObject("conversation")
// 	if getErr != nil {
// 		fmt.Errorf("\nFailed to get window \n%s", getErr.Error())
// 		c.Fail()
// 	}
// 	w, h := win.(gtki.Window).GetSize()
// 	c.Assert(h, Equals, 500)
// 	c.Assert(w, Equals, 400)
// }

// func (s *UIReaderSuite) Test_builderForDefinition_shouldReturnErrorWhenDefinitionDoesntExist(c *C) {
// 	removeFile("definitions/nonexistent")
// 	ui := "nonexistent"

// 	c.Assert(func() {
// 		builderForDefinition(ui)
// 	}, Panics, "No definition found for nonexistent")
// }

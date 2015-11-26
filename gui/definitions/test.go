
package definitions

func init(){
  add(`Test`, &defTest{})
}

type defTest struct{}

func (*defTest) String() string {
	return `
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
}

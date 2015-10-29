
package gui

type testWindow  struct{}

func (w testWindow) getDefinition() string {
	return `

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
}


package definitions

func init(){
  add(`TorNotRunning`, &defTorNotRunning{})
}

type defTorNotRunning struct{}

func (*defTorNotRunning) String() string {
	return `
<interface>
  <object class="GtkDialog" id="TorNotRunningDialog">
    <property name="default-height">150</property>
    <property name="default-width">300</property>
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <property name="title" translatable="yes">Tor is not running</property>
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkLabel" id="label">
            <property name="label" translatable="yes">Tor was not found running on your system. Make sure it is running.</property>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">0</property>
          </packing>
        </child>
      </object>
    </child>
  </object>
</interface>

`
}


package definitions

func init(){
  add(`MasterPasswordDefinition`, &defMasterPasswordDefinition{})
}

type defMasterPasswordDefinition struct{}

func (*defMasterPasswordDefinition) String() string {
	return `
<interface>
  <object class="GtkDialog" id="MasterPassword">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <property name="title">$title</property>
    <property name="default-width">160</property>
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkLabel" id="passMessage" >
            <property name="label" translatable="yes">Please enter the master password for the configuration file. You will not be asked for this password again until you restart CoyIM.</property>
            <property name="wrap">true</property>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">0</property>
          </packing>
        </child>
        <child>
          <object class="GtkEntry" id="password">
            <property name="has-focus">true</property>
            <property name="visibility">false</property>
            <signal name="activate" handler="on_save_signal" />
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">1</property>
          </packing>
        </child>
        <child>
          <object class="GtkBox" id="Hbox">
            <property name="homogeneous">false</property>
            <property name="border-width">5</property>
            <child>
              <object class="GtkButton" id="cancel">
                <property name="label">$cancelLabel</property>
                <signal name="clicked" handler="on_cancel_signal" />
              </object>
              <packing>
                <property name="expand">false</property>
                <property name="fill">true</property>
                <property name="padding">2</property>
                <property name="position">0</property>
              </packing>
            </child>
            <child>
              <object class="GtkButton" id="save">
                <property name="label">$saveLabel</property>
                <signal name="clicked" handler="on_save_signal" />
              </object>
              <packing>
                <property name="expand">false</property>
                <property name="fill">true</property>
                <property name="padding">2</property>
                <property name="position">1</property>
              </packing>
            </child>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">2</property>
          </packing>
        </child>
      </object>
    </child>
  </object>
</interface>

`
}

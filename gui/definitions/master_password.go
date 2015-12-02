package definitions

func init() {
	add(`MasterPassword`, &defMasterPassword{})
}

type defMasterPassword struct{}

func (*defMasterPassword) String() string {
	return `
<interface>
  <object class="GtkDialog" id="MasterPassword">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <property name="title" translatable="yes">Enter master password</property>
    <property name="default-width">300</property>
    <signal name="close" handler="on_cancel_signal" />
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="margin">10</property>
        <property name="spacing">10</property>
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
      </object>
    </child>
    <child internal-child="action_area">
      <object class="GtkButtonBox" id="button_box">
        <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
        <child>
          <object class="GtkButton" id="cancel">
            <property name="label" translatable="yes">Cancel</property>
            <signal name="clicked" handler="on_cancel_signal" />
          </object>
        </child>
        <child>
          <object class="GtkButton" id="save">
            <property name="label" translatable="yes">OK</property>
            <signal name="clicked" handler="on_save_signal" />
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>

`
}

package definitions

func init() {
	add(`RenameContact`, &defRenameContact{})
}

type defRenameContact struct{}

func (*defRenameContact) String() string {
	return `
<interface>
  <object class="GtkDialog" id="RenameContactPopup">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <property name="title" translatable="yes">Add contact</property>
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkBox" id="notification-area">
            <property name="visible">true</property>
            <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">0</property>
          </packing>
        </child>
        <child>
          <object class="GtkLabel" id="accountLabel" >
            <property name="label" translatable="yes">Contact's new name</property>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">1</property>
          </packing>
        </child>
        <child>
          <object class="GtkEntry" id="rename">
            <property name="has-focus">true</property>
            <signal name="activate" handler="on_rename_signal" />
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">2</property>
          </packing>
        </child>
        <child>
          <object class="GtkButton" id="add">
            <property name="label" translatable="yes">Ok</property>
            <signal name="clicked" handler="on_rename_signal" />
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">3</property>
          </packing>
        </child>
      </object>
    </child>
  </object>
</interface>

`
}

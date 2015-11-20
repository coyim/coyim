
package definitions

func init(){
  add(`AddContactDefinition`, &defAddContactDefinition{})
}

type defAddContactDefinition struct{}

func (*defAddContactDefinition) String() string {
	return `
<interface>
  <object class="GtkDialog" id="AddContact">
    <property name="window-position">1</property>
    <property name="title" translatable="yes">Add contact</property>
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkLabel" id="accountsLabel" >
            <property name="label" translatable="yes">Account</property>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">0</property>
          </packing>
        </child>
        <child>
          <object class="GtkComboBox" id="accounts">
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">1</property>
          </packing>
        </child>
        <child>
          <object class="GtkLabel" id="accountLabel" >
            <property name="label" translatable="yes">Contact to add (for example: arnoldsPub@jabber.ccc.de)</property>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">2</property>
          </packing>
        </child>
        <child>
          <object class="GtkEntry" id="address">
            <property name="has-focus">true</property>
            <signal name="activate" handler="on_save_signal" />
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">3</property>
          </packing>
        </child>
        <child>
          <object class="GtkButton" id="add">
            <property name="label" translatable="yes">Add</property>
            <signal name="clicked" handler="on_save_signal" />
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">4</property>
          </packing>
        </child>
      </object>
    </child>
  </object>
</interface>

`
}

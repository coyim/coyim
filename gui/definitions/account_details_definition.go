
package definitions

func init(){
  add(`AccountDetailsDefinition`, &defAccountDetailsDefinition{})
}

type defAccountDetailsDefinition struct{}

func (*defAccountDetailsDefinition) String() string {
	return `
<interface>
  <object class="GtkDialog" id="AccountDetailsDialog">
    <property name="title">$title</property>
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkLabel" id="AccountMessageLabel">
            <property name="label">$accountMessage</property>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">0</property>
          </packing>
        </child>
        <child>
          <object class="GtkEntry" id="account">
            <property name="has-focus">true</property>
            <signal name="activate" handler="on_save_signal" />
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">1</property>
          </packing>
        </child>
        <child>
          <object class="GtkLabel" id="PasswordLabel">
            <property name="label">$pswMessage</property>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">2</property>
          </packing>
        </child>
        <child>
          <object class="GtkEntry" id="password">
            <property name="visibility">false</property>
            <signal name="activate" handler="on_save_signal" />
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">3</property>
          </packing>
        </child>
        <child>
          <object class="GtkButton" id="save">
            <property name="label">$saveLabel</property>
            <signal name="clicked" handler="on_save_signal"/>
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

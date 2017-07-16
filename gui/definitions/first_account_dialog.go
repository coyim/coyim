package definitions

func init() {
	add(`FirstAccountDialog`, &defFirstAccountDialog{})
}

type defFirstAccountDialog struct{}

func (*defFirstAccountDialog) String() string {
	return `<interface>
  <object class="GtkDialog" id="dialog">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <property name="title" translatable="yes">Setup your first account</property>
    <signal name="delete-event" handler="on_cancel_signal" />
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="orientation">vertical</property>
	<property name="margin">20</property>
        <property name="margin_top">30</property>
	<property name="spacing">2</property>
        <child>
          <object class="GtkLabel" id="message">
            <property name="wrap">true</property>
	    <property name="label" translatable="yes">Welcome to CoyIM!</property>
            <property name="margin_bottom">10</property>
            <attributes>
              <attribute name="font-desc" value="&lt;Enter Value&gt; 14"/>
              <attribute name="style" value="normal"/>
              <attribute name="weight" value="semibold"/>
            </attributes>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">0</property>
          </packing>
        </child>
        <child internal-child="action_area">
          <object class="GtkButtonBox" id="button_box">
            <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
            <property name="homogeneous">True</property>
            <property name="margin_top">20</property>
	    <child>
              <object class="GtkButton" id="button_register">
                <property name="label" translatable="yes">Create a new account</property>
                <property name="margin_bottom">5</property>
		<signal name="clicked" handler="on_register_signal" />
              </object>
            </child>
            <child>
              <object class="GtkButton" id="button_existing">
                <property name="label" translatable="yes">Add an existing account</property>
                <property name="margin_bottom">5</property>
		<signal name="clicked" handler="on_existing_signal" />
              </object>
            </child>
            <child>
              <object class="GtkButton" id="button_import">
                <property name="label" translatable="yes">Import account from your computer</property>
                <property name="margin_bottom">5</property>
		<signal name="clicked" handler="on_import_signal" />
              </object>
            </child>
	    <child>
              <object class="GtkButton" id="button_cancel">
                <property name="label">Cancel</property>
                <property name="margin_bottom">5</property>
		<signal name="clicked" handler="on_cancel_signal" />
              </object>
            </child>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`
}

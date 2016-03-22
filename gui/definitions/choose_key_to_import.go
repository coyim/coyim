package definitions

func init() {
	add(`ChooseKeyToImport`, &defChooseKeyToImport{})
}

type defChooseKeyToImport struct{}

func (*defChooseKeyToImport) String() string {
	return `<interface>
  <object class="GtkDialog" id="dialog">
    <property name="title" translatable="yes">Choose a key to import</property>
    <signal name="close" handler="on_cancel_signal" />
    <signal name="delete-event" handler="on_cancel_signal" />

    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="margin">10</property>
        <property name="spacing">10</property>
        <child>
          <object class="GtkLabel" id="message">
            <property name="label" translatable="yes">The file you choose contains more than one private key. Choose from the list below the key you would like to import.</property>
            <property name="wrap">true</property>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
          </packing>
        </child>

        <child>
          <object class="GtkComboBoxText" id="keys">
          </object>
        </child>

        <child internal-child="action_area">
          <object class="GtkButtonBox" id="bbox">
            <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
            <child>
              <object class="GtkButton" id="btn-cancel">
                <property name="label" translatable="yes">Cancel</property>
                <signal name="clicked" handler="on_cancel_signal" />
              </object>
            </child>
            <child>
              <object class="GtkButton" id="btn-import">
                <property name="label" translatable="yes">Import</property>
                <property name="can-default">true</property>
                <signal name="clicked" handler="on_import_signal" />
              </object>
            </child>
          </object>
        </child>
      </object>
    </child>

    <action-widgets>
      <action-widget response="cancel">btn-cancel</action-widget>
      <action-widget response="ok" default="true">btn-import</action-widget>
    </action-widgets>
  </object>
</interface>
`
}

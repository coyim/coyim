package definitions

func init() {
	add(`VerifyFingerprintUnknown`, &defVerifyFingerprintUnknown{})
}

type defVerifyFingerprintUnknown struct{}

func (*defVerifyFingerprintUnknown) String() string {
	return `<interface>
  <object class="GtkDialog" id="dialog">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <child internal-child="vbox">
      <object class="GtkBox" id="box">
        <property name="border-width">10</property>
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkLabel" id="message" />
        </child>
        <child internal-child="action_area">
          <object class="GtkButtonBox" id="button_box">
            <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
            <child>
              <object class="GtkButton" id="button_ok">
                <property name="can-default">true</property>
                <property name="label" translatable="yes">OK</property>
              </object>
            </child>
          </object>
        </child>
      </object>
    </child>
    <action-widgets>
      <action-widget response="ok" default="true">button_ok</action-widget>
    </action-widgets>
  </object>
</interface>
`
}

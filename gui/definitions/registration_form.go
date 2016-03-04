package definitions

func init() {
	add(`RegistrationForm`, &defRegistrationForm{})
}

type defRegistrationForm struct{}

func (*defRegistrationForm) String() string {
	return `<interface>
  <object class="GtkDialog" id="dialog">
    <signal name="response" handler="response-handler" />
    <signal name="close" handler="close-handler" />

    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="visible">true</property>
        <property name="margin-bottom">10</property>
        <property name="margin">10</property>

        <child>
          <object class="GtkGrid" id="grid">
            <property name="visible">true</property>
            <property name="margin-bottom">10</property>
            <property name="row-spacing">12</property>
            <property name="column-spacing">6</property>
            <child>
              <object class="GtkLabel" id="instructions">
                <property name="visible">true</property>
              </object>
              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">0</property>
                <property name="width">2</property>
              </packing>
            </child>
          </object>
        </child>

        <child internal-child="action_area">
          <object class="GtkButtonBox" id="bbox">
            <property name="visible">true</property>
            <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
            <child>
              <object class="GtkButton" id="btn-cancel">
                <property name="visible">true</property>
                <property name="label" translatable="yes">Cancel</property>
              </object>
            </child>
            <child>
              <object class="GtkButton" id="btn-register">
                <property name="visible">true</property>
                <property name="label" translatable="yes">Register</property>
                <property name="can-default">true</property>
              </object>
            </child>
          </object>
        </child>
      </object>
    </child>

    <action-widgets>
      <action-widget response="cancel">btn-cancel</action-widget>
      <action-widget response="apply" default="true">btn-register</action-widget>
    </action-widgets>
  </object>
</interface>
`
}

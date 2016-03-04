package definitions

func init() {
	add(`VerifyIdentityNotification`, &defVerifyIdentityNotification{})
}

type defVerifyIdentityNotification struct{}

func (*defVerifyIdentityNotification) String() string {
	return `<interface>
  <object class="GtkInfoBar" id="infobar">
    <property name="message-type">GTK_MESSAGE_WARNING</property>
    <child internal-child="content_area">
      <object class="GtkBox" id="box">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
        <child>
          <object class="GtkLabel" id="message">
            <property name="wrap">true</property>
          </object>
        </child>
      </object>
    </child>
    <child internal-child="action_area">
      <object class="GtkBox" id="button_box">
        <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
        <child>
          <object class="GtkButton" id="button_verify">
            <property name="label" translatable="yes">Verify</property>
            <property name="can-default">true</property>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`
}

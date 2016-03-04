package definitions

func init() {
	add(`ConnectionFailureNotification`, &defConnectionFailureNotification{})
}

type defConnectionFailureNotification struct{}

func (*defConnectionFailureNotification) String() string {
	return `<interface>
  <object class="GtkInfoBar" id="infobar">
    <property name="message-type">GTK_MESSAGE_ERROR</property>
    <property name="show-close-button">true</property>
    <signal name="response" handler="handleResponse" />
    <child internal-child="content_area">
      <object class="GtkBox" id="box">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
        <child>
          <object class="GtkLabel" id="message">
            <property name="ellipsize">PANGO_ELLIPSIZE_END</property>
            <property name="wrap">true</property>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`
}

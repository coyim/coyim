package definitions

func init() {
	add(`SimpleNotification`, &defSimpleNotification{})
}

type defSimpleNotification struct{}

func (*defSimpleNotification) String() string {
	return `<interface>
  <object class="GtkMessageDialog" id="dialog">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <property name="modal">true</property>
    <property name="message-type">GTK_MESSAGE_INFO</property>
    <property name="buttons">GTK_BUTTONS_OK</property>
  </object>
</interface>
`
}

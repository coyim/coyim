package definitions

func init() {
	add(`Conversation`, &defConversation{})
}

type defConversation struct{}

func (*defConversation) String() string {
	return `<interface>
  <object class="GtkWindow" id="conversation">
    <property name="window-position">GTK_WIN_POS_NONE</property>
    <property name="default-height">500</property>
    <property name="default-width">400</property>
    <property name="destroy-with-parent">true</property>
    <signal name="enable" handler="on_connect" />
    <signal name="disable" handler="on_disconnect" />
    <child>
      <object class="GtkBox" id="box">
        <property name="visible">true</property>
        <property name="homogeneous">false</property>
    	<property name="orientation">GTK_ORIENTATION_VERTICAL</property>
      </object>
    </child>
  </object>
</interface>
`
}

package definitions

func init() {
	add(`XMLConsole`, &defXMLConsole{})
}

type defXMLConsole struct{}

func (*defXMLConsole) String() string {
	return `<?xml version="1.0" encoding="utf-8"?>
<interface>
  <object class="GtkTextBuffer" id="consoleContent">
    <property name="text" translatable="yes">Hello world</property>
  </object>
  <object class="GtkDialog" id="XMLConsole">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <property name="border_width">6</property>
    <property name="title" translatable="yes">XMPP Console: ACCOUNT_NAME</property>
    <property name="resizable">True</property>
    <property name="default-height">400</property>
    <property name="default-width">500</property>
    <property name="destroy-with-parent">true</property>
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="homogeneous">false</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <property name="spacing">6</property>
        <child>
          <object class="GtkScrolledWindow" id="message_scrolledwindow">
            <property name="visible">True</property>
            <property name="can_focus">True</property>
            <property name="no_show_all">True</property>
            <property name="border_width">6</property>
            <property name="shadow_type">etched-in</property>
            <property name="min_content_height">5</property>
            <child>
              <object class="GtkTextView" id="message_textview">
                <property name="visible">true</property>
                <property name="wrap-mode">2</property>
                <property name="editable">false</property>
                <property name="cursor-visible">false</property>
                <property name="pixels-below-lines">5</property>
                <property name="left-margin">5</property>
                <property name="right-margin">5</property>
                <property name="buffer">consoleContent</property>
              </object>
            </child>
          </object>
          <packing>
            <property name="expand">True</property>
            <property name="fill">True</property>
          </packing>
        </child>
        <child internal-child="action_area">
          <object class="GtkButtonBox" id="button_box">
            <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
            <child>
              <object class="GtkButton" id="button_refresh">
                <property name="label">_Refresh</property>
                <property name="use-underline">True</property>
                <signal name="clicked" handler="on_refresh_signal" />
              </object>
            </child>
            <child>
              <object class="GtkButton" id="button_close">
                <property name="label" translatable="yes">_Close</property>
                <property name="use-underline">True</property>
                <property name="can-default">true</property>
                <signal name="clicked" handler="on_close_signal" />
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

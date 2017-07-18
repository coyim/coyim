package definitions

func init() {
	add(`Feedback`, &defFeedback{})
}

type defFeedback struct{}

func (*defFeedback) String() string {
	return `<interface>
  <object class="GtkDialog" id="dialog">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <property name="title" translatable="yes">We would like to receive your feedback</property>
    <property name="border_width">7</property>
    <child internal-child="vbox">
      <object class="GtkBox" id="box">
        <property name="border-width">10</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
	<child>
          <object class="GtkLabel">
            <property name="can_focus">False</property>
            <property name="margin-bottom">3</property>
	    <property name="label" translatable="yes">
Visit the website to get in touch with us:
            </property>
          </object>
          <packing>
            <property name="expand">False</property>
            <property name="fill">True</property>
            <property name="position">0</property>
          </packing>
        </child>
	<child>
          <object class="GtkLabel" id="message">
            <property name="can_focus">False</property>
            <property name="margin-top">10</property>
	    <property name="label" translatable="yes">https://coy.im</property>
            <property name="selectable">True</property>
            <attributes>
              <attribute name="font-desc" value="&lt;Enter Value&gt; 14"/>
	      <attribute name="weight" value="semibold"/>
            </attributes>
          </object>
          <packing>
            <property name="expand">False</property>
            <property name="fill">True</property>
            <property name="position">1</property>
          </packing>
        </child>
        <child>
          <object class="GtkLabel">
            <property name="can_focus">False</property>
            <property name="margin-top">10</property>
            <property name="margin-bottom">10</property>
	    <property name="label" translatable="yes">
Let us know what you think of CoyIM.

This is the only way we can create a better privacy tool.
            </property>
          </object>
          <packing>
            <property name="expand">False</property>
            <property name="fill">True</property>
            <property name="position">2</property>
          </packing>
        </child>
      </object>
    </child>
    <child internal-child="action_area">
      <object class="GtkButtonBox" id="button_box">
        <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
        <child>
          <object class="GtkButton" id="close">
            <property name="label" translatable="yes">Close</property>
            <signal name="clicked" handler="on_close_signal" />
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`
}

package definitions

func init() {
	add(`FeedbackInfo`, &defFeedbackInfo{})
}

type defFeedbackInfo struct{}

func (*defFeedbackInfo) String() string {
	return `
<interface>
  <child>
    <object class="GtkInfoBar" id="feedbackInfo">
      <property name="message-type">GTK_MESSAGE_OTHER</property>
      <child internal-child="content_area">
        <object class="GtkBox" id="box">
          <property name="homogeneous">false</property>
          <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
          <child>
            <object class="GtkLabel" id="feedbackMessage">
              <property name="wrap">true</property>
              <property name="label" translatable="yes">Are you liking it?</property>
            </object>
          </child>
          <child>
            <object class="GtkButton" id="feedbackButton">
              <property name="label" translatable="yes">Feedback</property>
              <signal name="activate" handler="on_click_button" />
              <property name="halign">GTK_ALIGN_END</property>
              <property name="hexpand">true</property>
            </object>
          </child>
        </object>
      </child>
    </object>
  </child>
</interface>

`
}

package definitions

func init() {
	add(`FeedbackInfo`, &defFeedbackInfo{})
}

type defFeedbackInfo struct{}

func (*defFeedbackInfo) String() string {
	return `
<interface>
  <object class="GtkInfoBar" id="feedbackInfo">
    <property name="message-type">GTK_MESSAGE_OTHER</property>
    <property name="show-close-button">true</property>
    <signal name="response" handler="handleResponse" />
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
      </object>
    </child>
    <child internal-child="action_area">
      <object class="GtkBox" id="button_box">
        <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
        <child>
          <object class="GtkButton" id="feedbackButton">
            <property name="label" translatable="yes">Feedback</property>
            <property name="halign">GTK_ALIGN_END</property>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>

`
}

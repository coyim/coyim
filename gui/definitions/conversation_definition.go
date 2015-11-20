
package definitions

func init(){
  add(`ConversationDefinition`, &defConversationDefinition{})
}

type defConversationDefinition struct{}

func (*defConversationDefinition) String() string {
	return `
<interface>
  <object class="GtkWindow" id="conversation">
    <property name="window-position">0</property>
    <property name="default-height">500</property>
    <property name="default-width">400</property>
    <property name="destroy-with-parent">true</property>
    <property name="title">$title</property>
    <child>
      <object class="GtkVBox" id="box">
        <property name="homogeneous">false</property>
        <child>
          <object class="GtkMenuBar" id="menubar">
            <child>
              <object class="GtkMenuItem" id="conversationMenu">
                <property name="label">$DevOptions</property>
                <child type="submenu">
                  <object class="GtkMenu" id="menu">
                    <child>
                      <object class="GtkMenuItem" id="startOTRMenu">
                        <property name="label">$StartOTR</property>
                        <signal name="activate" handler="on_start_otr_signal" />
                      </object>
                    </child>
                    <child>
                      <object class="GtkMenuItem" id="endOTRMenu">
                        <property name="label">$EndOTR</property>
                        <signal name="activate" handler="on_end_otr_signal" />
                      </object>
                    </child>
                    <child>
                      <object class="GtkMenuItem" id="verifyFingerMenu">
                        <property name="label">$VerifyFP</property>
                        <signal name="activate" handler="on_verify_fp_signal" />
                      </object>
                    </child>
                  </object>
                </child>
              </object>
            </child>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">0</property>
          </packing>
        </child>
        <child>
          <object class="GtkScrolledWindow" id="historyScroll">
            <property name="vscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
            <property name="hscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
            <child>
              <object class="GtkTextView" id="history">
                <property name="visible">true</property>
                <property name="wrap-mode">2</property>
                <property name="editable">false</property>
                <property name="cursor-visible">false</property>
              </object>
            </child>
          </object>
          <packing>
            <property name="expand">true</property>
            <property name="fill">true</property>
            <property name="position">1</property>
          </packing>
        </child>
        <child>
          <object class="GtkEntry" id="message">
            <property name="has-focus">true</property>
            <signal name="activate" handler="on_send_message_signal" />
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">true</property>
            <property name="position">2</property>
          </packing>
        </child>
      </object>
    </child>
  </object>
</interface>

`
}

package definitions

func init() {
	add(`ConversationPane`, &defConversationPane{})
}

type defConversationPane struct{}

func (*defConversationPane) String() string {
	return `<interface>
  <object class="GtkBox" id="box">
    <property name="visible">true</property>
    <property name="homogeneous">false</property>
    <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
    <child>
      <object class="GtkMenuBar" id="menubar">
        <property name="visible">true</property>
        <child>
          <object class="GtkMenuItem" id="conversationMenu">
            <property name="visible">true</property>
            <property name="label" translatable="yes">_Encryption</property>
            <property name="use-underline">True</property>
            <child type="submenu">
              <object class="GtkMenu" id="menu">
                <property name="visible">true</property>
                <child>
                  <object class="GtkMenuItem" id="startOTRMenu">
                    <property name="visible">true</property>
                    <property name="label" translatable="yes">Start encrypted chat</property>
                    <signal name="activate" handler="on_start_otr_signal" />
                  </object>
                </child>
                <child>
                  <object class="GtkMenuItem" id="endOTRMenu">
                    <property name="visible">true</property>
                    <property name="label" translatable="yes">End encrypted chat</property>
                    <signal name="activate" handler="on_end_otr_signal" />
                  </object>
                </child>
                <child>
                  <object class="GtkMenuItem" id="verifyFingerMenu">
                    <property name="visible">true</property>
                    <property name="label" translatable="yes">Verify fingerprint</property>
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
        <property name="pack-type">GTK_PACK_START</property>
        <property name="position">0</property>
      </packing>
    </child>
    <child>
      <object class="GtkBox" id="notification-area">
        <property name="visible">true</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkInfoBar" id="security-warning">
            <property name="visible">false</property>
            <property name="message-type">GTK_MESSAGE_WARNING</property>
            <child internal-child="content_area">
              <object class="GtkBox" id="security-warning-box">
                <property name="visible">true</property>
                <property name="homogeneous">false</property>
                <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
                <child>
                  <object class="GtkLabel" id="warning-message">
                    <property name="visible">true</property>
                    <property name="wrap">true</property>
                    <property name="label" translatable="yes">You are talking over an unprotected channel</property>
                  </object>
                </child>
              </object>
            </child>
            <child internal-child="action_area">
              <object class="GtkBox" id="button_box">
                <property name="visible">true</property>
                <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
                <child>
                  <object class="GtkButton" id="button_verify">
                    <property name="visible">true</property>
                    <property name="label" translatable="yes">Secure chat</property>
                    <property name="can-default">true</property>
                    <signal name="clicked" handler="on_start_otr_signal" />
                  </object>
                </child>
              </object>
            </child>
          </object>
          <packing>
            <property name="expand">false</property>
            <property name="fill">false</property>
            <property name="pack-type">GTK_PACK_START</property>
            <property name="position">0</property>
          </packing>
        </child>
      </object>
      <packing>
        <property name="expand">false</property>
        <property name="fill">false</property>
        <property name="pack-type">GTK_PACK_START</property>
        <property name="position">1</property>
      </packing>
    </child>
    <child>
      <object class="GtkScrolledWindow" id="historyScroll">
        <property name="visible">true</property>
        <property name="vscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
        <property name="hscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
        <child>
          <object class="GtkTextView" id="history">
            <property name="visible">true</property>
            <property name="wrap-mode">2</property>
            <property name="editable">false</property>
            <property name="cursor-visible">false</property>
            <property name="pixels-below-lines">5</property>
            <property name="left-margin">5</property>
            <property name="right-margin">5</property>
          </object>
        </child>
      </object>
      <packing>
        <property name="expand">true</property>
        <property name="fill">true</property>
        <property name="pack-type">GTK_PACK_END</property>
        <property name="position">2</property>
      </packing>
    </child>
    
    <child>
      <object class="GtkScrolledWindow" id="pendingScroll">
        <property name="visible">true</property>
        <property name="vscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
        <property name="hscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
        <child>
          <object class="GtkTextView" id="pending">
            <property name="visible">true</property>
            <property name="wrap-mode">2</property>
            <property name="editable">false</property>
            <property name="cursor-visible">false</property>
            <property name="pixels-below-lines">5</property>
            <property name="left-margin">5</property>
            <property name="right-margin">5</property>
          </object>
        </child>
      </object>
      <packing>
        <property name="expand">false</property>
        <property name="fill">true</property>
        <property name="pack-type">GTK_PACK_END</property>
        <property name="position">1</property>
      </packing>
    </child>

    <child>
      <object class="GtkScrolledWindow" id="messageScroll">
        <property name="visible">true</property>
        <property name="vscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
        <property name="hscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
        <property name="shadow-type">in</property>
        <child>
          <object class="GtkTextView" id="message">
            <property name="visible">true</property>
            <property name="has-focus">true</property>
            <property name="wrap-mode">GTK_WRAP_WORD_CHAR</property>
            <property name="editable">true</property>
            <property name="left-margin">3</property>
            <property name="right-margin">3</property>
          </object>
        </child>
      </object>

      <packing>
        <property name="expand">false</property>
        <property name="fill">true</property>
        <property name="pack-type">GTK_PACK_END</property>
        <property name="position">0</property>
      </packing>
    </child>
  </object>
</interface>
`
}

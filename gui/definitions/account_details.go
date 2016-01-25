package definitions

func init() {
	add(`AccountDetails`, &defAccountDetails{})
}

type defAccountDetails struct{}

func (*defAccountDetails) String() string {
	return `
<interface>
  <object class="GtkDialog" id="AccountDetails">
    <property name="title" translatable="yes">Account Details</property>
    <signal name="close" handler="on_cancel_signal" />
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="margin">10</property>

        <child>
          <object class="GtkBox" id="notification-area">
            <property name="visible">true</property>
            <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
          </object>
        </child>

        <child>
          <object class="GtkNotebook" id="notebook1">
            <property name="visible">True</property>
            <property name="show-border">False</property>
            <property name="page">0</property>
            <property name="margin-bottom">10</property>
            <child>
              <object class="GtkGrid" id="grid">
                <property name="margin-top">15</property>
                <property name="margin-bottom">10</property>
                <property name="margin-start">10</property>
                <property name="margin-end">10</property>
                <property name="row-spacing">12</property>
                <property name="column-spacing">6</property>
                <child>
                  <object class="GtkLabel" id="AccountMessageLabel">
                    <property name="label" translatable="yes">Your account&#xA;(example: kim42@dukgo.com)</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">0</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="account">
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">0</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkLabel" id="PasswordLabel">
                    <property name="label" translatable="yes">Password</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="password">
                    <property name="visibility">false</property>
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="showOtherSettings">
                    <property name="label" translatable="yes">Other settings</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">2</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkCheckButton" id="otherSettings">
                    <signal name="toggled" handler="on_toggle_other_settings" />
                  </object>
                  
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">2</property>
                  </packing>
                </child>
              </object>
            </child>

            <child type="tab">
              <object class="GtkLabel" id="label-tab1">
                <property name="label" translatable="yes">Account</property>
                <property name="visible">True</property>
              </object>
              <packing>
                <property name="position">0</property>
                <property name="tab_fill">False</property>
              </packing>
            </child>

            <child>
              <object class="GtkGrid" id="otherOptionsGrid">
                <property name="margin-top">15</property>
                <property name="margin-bottom">10</property>
                <property name="margin-start">10</property>
                <property name="margin-end">10</property>
                <property name="row-spacing">12</property>
                <property name="column-spacing">6</property>
                <child>
                  <object class="GtkLabel" id="serverLabel">
                    <property name="label" translatable="yes">Server (leave empty for default)</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">0</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="server">
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">0</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkLabel" id="portLabel">
                    <property name="label" translatable="yes">Port (leave empty for default)</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="port">
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>

              </object>
            </child>
            <child type="tab">
              <object class="GtkLabel" id="label-tab2">
                <property name="label" translatable="yes">Other options</property>
                <property name="visible">True</property>
              </object>
              <packing>
                <property name="position">1</property>
                <property name="tab_fill">False</property>
              </packing>
            </child>
            <child>
              <object class="GtkLabel" id="label-interesting">
                <property name="label" translatable="yes">Interesting</property>
                <property name="visible">True</property>
              </object>
            </child>
            <child type="tab">
              <object class="GtkLabel" id="label-tab3">
                <property name="label" translatable="yes">Proxies</property>
                <property name="visible">True</property>
              </object>
              <packing>
                <property name="position">2</property>
                <property name="tab_fill">False</property>
              </packing>
            </child>
          </object>
        </child>

        <child internal-child="action_area">
          <object class="GtkButtonBox" id="button_box">
            <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
            <child>
              <object class="GtkButton" id="cancel">
                <property name="label" translatable="yes">Cancel</property>
                <signal name="clicked" handler="on_cancel_signal"/>
              </object>
            </child>
            <child>
              <object class="GtkButton" id="save">
                <property name="label" translatable="yes">Save</property>
                <property name="can-default">true</property>
                <signal name="clicked" handler="on_save_signal"/>
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

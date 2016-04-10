package definitions

func init() {
	add(`GlobalPreferences`, &defGlobalPreferences{})
}

type defGlobalPreferences struct{}

func (*defGlobalPreferences) String() string {
	return `<interface>
  <object class="GtkListStore" id="notification-type-model">
    <columns>
      <!-- notification name -->
      <column type="gchararray"/>
      <!-- notification id -->
      <column type="gchararray"/>
    </columns>
    <data>
      <row>
        <col id="0" translatable="yes">No notifications</col>
        <col id="1">off</col>
      </row>
      <row>
        <col id="0" translatable="yes">Only show that a new message arrived</col>
        <col id="1">only-presence-of-new-information</col>
      </row>
      <row>
        <col id="0" translatable="yes">Show who sent the message</col>
        <col id="1">with-author-but-no-content</col>
      </row>
      <row>
        <col id="0" translatable="yes">Show message</col>
        <col id="1">with-content</col>
      </row>
    </data>
  </object>

  <object class="GtkAdjustment" id="adjustment1">
    <property name="lower">0</property>
    <property name="upper">3600</property>
    <property name="step_increment">60</property>
    <property name="page_increment">60</property>
  </object>

  <object class="GtkDialog" id="GlobalPreferences">
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
              <object class="GtkGrid" id="generalGrid">
                <property name="margin-top">15</property>
                <property name="margin-bottom">10</property>
                <property name="margin-start">10</property>
                <property name="margin-end">10</property>
                <property name="row-spacing">12</property>
                <property name="column-spacing">6</property>
                <child>
                  <object class="GtkLabel" id="singleWindowLabel">
                    <property name="label" translatable="yes">Unify conversations in one window</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">0</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkCheckButton" id="singleWindow">
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">0</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="slashMeLabel">
                    <property name="label" translatable="yes">Render /me commands</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkCheckButton" id="slashMe">
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="sendWithShiftEnterLabel">
                    <property name="label" translatable="yes">Send messages with Shift-Enter</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">2</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkCheckButton" id="sendWithShiftEnter">
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">2</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="showEmptyGroupsLabel">
                    <property name="label" translatable="yes">Display empty groups</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">3</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkCheckButton" id="showEmptyGroups">
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">3</property>
                  </packing>
                </child>
              </object>
            </child>

            <child type="tab">
              <object class="GtkLabel" id="label-general-tab">
                <property name="label" translatable="yes">General</property>
                <property name="visible">True</property>
              </object>
              <packing>
                <property name="position">0</property>
                <property name="tab-fill">False</property>
              </packing>
            </child>

            <child>
              <object class="GtkGrid" id="notificationsGrid">
                <property name="margin-top">15</property>
                <property name="margin-bottom">10</property>
                <property name="margin-start">10</property>
                <property name="margin-end">10</property>
                <property name="row-spacing">12</property>
                <property name="column-spacing">6</property>

                <child>
                  <object class="GtkLabel" id="notificationTypeInstructions">
                    <property name="label" translatable="yes">CoyIM supports notifying you when a new message arrives - these notifications can be turned off, or set to show different amounts of information. It's important to remember that showing the text of the message, or even the author, could potentially be a privacy leak if you're ever in settings where someone can see your screen.</property>
                    <property name="visible">true</property>
                    <property name="wrap">true</property>
                    <property name="max-width-chars">50</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">0</property>
                    <property name="width">2</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="notificationTypeLabel">
                    <property name="label" translatable="yes">Notification type</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkComboBox" id="notificationsType">
                    <property name="model">notification-type-model</property>
                    <property name="has-focus">true</property>
                    <property name="hexpand">True</property>
                    <child>
                      <object class="GtkCellRendererText" id="notification-name-rendered"/>
                      <attributes>
                        <attribute name="text">0</attribute>
                      </attributes>
                    </child>
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="notificationUrgencyInstructions">
                    <property name="label" translatable="yes">Notifications can be set to display urgently - this is useful if you work in fullscreen mode. If the notification is not urgent, it will not display over a fullscreen window.</property>
                    <property name="visible">true</property>
                    <property name="wrap">true</property>
                    <property name="max-width-chars">50</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">2</property>
                    <property name="width">2</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="notificationUrgentLabel">
                    <property name="label" translatable="yes">Should notifications be displayed urgently?</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">3</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkCheckButton" id="notificationUrgent">
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">3</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="notificationExpiryInstructions">
                    <property name="label" translatable="yes">Notifications can stay on the screen until you've gone back to CoyIM, or they can expire after a while. The below setting changes this behavior.</property>
                    <property name="visible">true</property>
                    <property name="wrap">true</property>
                    <property name="max-width-chars">50</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">4</property>
                    <property name="width">2</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="notificationExpiresLabel">
                    <property name="label" translatable="yes">Should notifications expire?</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">5</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkCheckButton" id="notificationExpires">
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">5</property>
                  </packing>
                </child>


                <child>
                  <object class="GtkLabel" id="notificationCommandInstructions">
                    <property name="label" translatable="yes">You can specify a custom command to run whenever a message is received. If specificed, this command will run on every messages except it will wait for a timeout period before running the next time. The command and timeout is specified below. </property>
                    <property name="visible">true</property>
                    <property name="wrap">true</property>
                    <property name="max-width-chars">50</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">6</property>
                    <property name="width">2</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="notificationCommandLabel">
                    <property name="label" translatable="yes">Notification command</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">7</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="notificationCommand">
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">7</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="notificationTimeoutLabel">
                    <property name="label" translatable="yes">Minimum time between notifications in seconds</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">8</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkSpinButton" id="notificationTimeout">
                    <property name="climb-rate">60</property>
                    <property name="adjustment">adjustment1</property>
                    <signal name="activate" handler="on_save_signal" />
                  </object>
                  <packing>
                    <property name="left-attach">1</property>
                    <property name="top-attach">8</property>
                  </packing>
                </child>

              </object>
            </child>
            <child type="tab">
              <object class="GtkLabel" id="label-notifications-tab">
                <property name="label" translatable="yes">Notifications</property>
                <property name="visible">True</property>
              </object>
              <packing>
                <property name="position">1</property>
                <property name="tab-fill">False</property>
              </packing>
            </child>

            <child>
              <object class="GtkGrid" id="debuggingGrid">
                <property name="margin-top">15</property>
                <property name="margin-bottom">10</property>
                <property name="margin-start">10</property>
                <property name="margin-end">10</property>
                <property name="row-spacing">12</property>
                <property name="column-spacing">6</property>
                <child>
                  <object class="GtkLabel" id="rawLogFileInstructions">
                    <property name="label" translatable="yes">If you set this property to a file name, low level information will be logged there. Be very careful - this information is sensitive and could potentially contain very private information. Only turn this setting on if you absolutely need it for debugging. This file will specifically log XMPP traffic information. This setting will only take effect after a restart of CoyIM.</property>
                    <property name="visible">true</property>
                    <property name="wrap">true</property>
                    <property name="max-width-chars">50</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">0</property>
                    <property name="width">2</property>
                  </packing>
                </child>

                <child>
                  <object class="GtkLabel" id="rawLogFileLabel">
                    <property name="label" translatable="yes">Raw log file</property>
                    <property name="justify">GTK_JUSTIFY_RIGHT</property>
                    <property name="halign">GTK_ALIGN_END</property>
                  </object>
                  <packing>
                    <property name="left-attach">0</property>
                    <property name="top-attach">1</property>
                  </packing>
                </child>
                <child>
                  <object class="GtkEntry" id="rawLogFile">
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
              <object class="GtkLabel" id="label-debugging-tab">
                <property name="label" translatable="yes">Debugging</property>
                <property name="visible">True</property>
              </object>
              <packing>
                <property name="position">2</property>
                <property name="tab-fill">False</property>
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

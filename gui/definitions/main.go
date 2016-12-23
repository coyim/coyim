package definitions

func init() {
	add(`Main`, &defMain{})
}

type defMain struct{}

func (*defMain) String() string {
	return `<interface>
  <object class="GtkApplicationWindow" id="mainWindow">
    <property name="can_focus">False</property>
    <property name="title">CoyIM</property>
    <property name="default_width">200</property>
    <property name="default_height">600</property>
    <signal name="destroy" handler="on_close_window_signal" swapped="no"/>
    <child>
      <object class="GtkBox" id="Hbox">
        <property name="visible">True</property>
        <property name="can_focus">False</property>
        <child>
          <object class="GtkBox" id="Vbox">
            <property name="can_focus">False</property>
            <property name="orientation">vertical</property>
            <child>
              <object class="GtkMenuBar" id="menubar">
                <property name="can_focus">False</property>
                <child>
                  <object class="GtkMenuItem" id="ContactsMenu">
                    <property name="can_focus">False</property>
                    <property name="label" translatable="yes">_Contacts</property>
                    <property name="use-underline">True</property>
                    <child type="submenu">
                      <object class="GtkMenu" id="menu">
                        <property name="can_focus">False</property>
                        <child>
                          <object class="GtkMenuItem" id="addMenu">
                            <property name="can_focus">False</property>
                            <property name="label" translatable="yes">Add...</property>
                            <signal name="activate" handler="on_add_contact_window_signal" swapped="no"/>
                          </object>
                        </child>
                        <child>
                          <object class="GtkMenuItem" id="newConvMenu">
                            <property name="can_focus">False</property>
                            <property name="label" translatable="yes">New Conversation...</property>
                            <signal name="activate" handler="on_new_conversation_signal" swapped="no"/>
                          </object>
                        </child>
                      </object>
                    </child>
                  </object>
                </child>
                <child>
                  <object class="GtkMenuItem" id="AccountsMenu">
                    <property name="can_focus">False</property>
                    <property name="label" translatable="yes">_Accounts</property>
                    <property name="use-underline">True</property>
                  </object>
                </child>
                <child>
                  <object class="GtkMenuItem" id="ViewMenu">
                    <property name="can_focus">False</property>
                    <property name="label" translatable="yes">_View</property>
                    <property name="use-underline">True</property>
                    <child type="submenu">
                      <object class="GtkMenu" id="menu2">
                        <property name="can_focus">False</property>
                        <child>
                          <object class="GtkCheckMenuItem" id="CheckItemMerge">
                            <property name="can_focus">False</property>
                            <property name="label" translatable="yes">Merge Accounts</property>
                            <signal name="toggled" handler="on_toggled_check_Item_Merge_signal" swapped="no"/>
                          </object>
                        </child>
                        <child>
                          <object class="GtkCheckMenuItem" id="CheckItemShowOffline">
                            <property name="can_focus">False</property>
                            <property name="label" translatable="yes">Show Offline Contacts</property>
                            <signal name="toggled" handler="on_toggled_check_Item_Show_Offline_signal" swapped="no"/>
                          </object>
                        </child>
                        <child>
                          <object class="GtkCheckMenuItem" id="CheckItemShowWaiting">
                            <property name="can_focus">False</property>
                            <property name="label" translatable="yes">Show Waiting Contacts</property>
                            <signal name="toggled" handler="on_toggled_check_Item_Show_Waiting_signal" swapped="no"/>
                          </object>
                        </child>
                        <child>
                          <object class="GtkCheckMenuItem" id="CheckItemSortStatus">
                            <property name="can_focus">False</property>
                            <property name="label" translatable="yes">Sort By Status</property>
                            <signal name="toggled" handler="on_toggled_check_Item_Sort_By_Status_signal" swapped="no"/>
                          </object>
                        </child>
                      </object>
                    </child>
                  </object>
                </child>
                <child>
                  <object class="GtkMenuItem" id="OptionsMenu">
                    <property name="can-focus">False</property>
                    <property name="label" translatable="yes">_Options</property>
                    <property name="use-underline">True</property>
                    <child type="submenu">
                      <object class="GtkMenu" id="options_submenu">
                        <property name="can_focus">False</property>
                        <child>
                          <object class="GtkCheckMenuItem" id="EncryptConfigurationFileCheckMenuItem">
                            <property name="can_focus">False</property>
                            <property name="label" translatable="yes">Encrypt configuration file</property>
                            <signal name="toggled" handler="on_toggled_encrypt_configuration_file_signal" swapped="no"/>
                          </object>
                        </child>
                        <child>
                          <object class="GtkMenuItem" id="preferencesMenuItem">
                            <property name="can_focus">False</property>
                            <property name="label" translatable="yes">_Preferences...</property>
                            <property name="use-underline">True</property>
                            <signal name="activate" handler="on_preferences_signal" swapped="no"/>
                          </object>
                        </child>
                      </object>
                    </child>
                  </object>
                </child>
                <child>
                  <object class="GtkMenuItem" id="HelpMenu">
                    <property name="can_focus">False</property>
                    <property name="label" translatable="yes">Help</property>
                    <child type="submenu">
                      <object class="GtkMenu" id="menu3">
                        <property name="can_focus">False</property>
                        <child>
                          <object class="GtkMenuItem" id="feedbackMenu">
                            <property name="can_focus">False</property>
                            <property name="label" translatable="yes">Feedback</property>
                            <signal name="activate" handler="on_feedback_dialog_signal" swapped="no"/>
                          </object>
                        </child>
                        <child>
                          <object class="GtkMenuItem" id="aboutMenu">
                            <property name="can_focus">False</property>
                            <property name="label" translatable="yes">About</property>
                            <signal name="activate" handler="on_about_dialog_signal" swapped="no"/>
                          </object>
                        </child>
                      </object>
                    </child>
                  </object>
                </child>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">True</property>
                <property name="position">0</property>
              </packing>
            </child>
            <child>
              <object class="GtkBox" id="notification-area">
                <property name="can_focus">False</property>
                <property name="orientation">vertical</property>
              </object>
              <packing>
                <property name="expand">False</property>
                <property name="fill">False</property>
                <property name="pack_type">end</property>
                <property name="position">0</property>
              </packing>
            </child>
          </object>
          <packing>
            <property name="expand">True</property>
            <property name="fill">True</property>
            <property name="position">0</property>
          </packing>
        </child>
      </object>
    </child>
  </object>
</interface>
`
}

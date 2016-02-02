package definitions

func init() {
	add(`PeerDetails`, &defPeerDetails{})
}

type defPeerDetails struct{}

func (*defPeerDetails) String() string {
	return `
<interface>
  <object class="GtkDialog" id="dialog">
    <signal name="close" handler="on_cancel_signal" />
    <child internal-child="vbox">
      <object class="GtkBox" id="Vbox">
        <property name="margin">10</property>

        <child>
          <object class="GtkGrid" id="grid">
            <property name="margin-top">15</property>
            <property name="margin-bottom">10</property>
            <property name="margin-start">10</property>
            <property name="margin-end">10</property>
            <property name="row-spacing">12</property>
            <property name="column-spacing">6</property>
            <child>
              <object class="GtkLabel" id="nickname-label">
                <property name="label" translatable="yes">Nickname</property>
                <property name="justify">GTK_JUSTIFY_RIGHT</property>
              </object>
              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">0</property>
              </packing>
            </child>
            <child>
              <object class="GtkEntry" id="nickname">
                <signal name="activate" handler="on_save_signal" />
              </object>
              <packing>
                <property name="left-attach">1</property>
                <property name="top-attach">0</property>
              </packing>
            </child>
            <child>
              <object class="GtkLabel" id="groups">
                <property name="label" translatable="yes">Groups</property>
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
                <property name="left-attach">0</property>
                <property name="top-attach">2</property>
              </packing>
            </child>


            <child>
              <object class="GtkTreeView" id="proxies-view">
                <property name="model">proxies-model</property>
                <property name="visible">True</property>
                <property name="can-focus">True</property>
                <property name="headers-visible">False</property>
                <property name="show-expanders">False</property>
                <property name="reorderable">True</property>
                <signal name="row-activated" handler="on_edit_activate_proxy_signal" />
                <child internal-child="selection">
                  <object class="GtkTreeSelection" id="selection">
                    <property name="mode">GTK_SELECTION_SINGLE</property>
                  </object>
                </child>
                <child>
                  <object class="GtkTreeViewColumn" id="proxy-name-column">
                    <property name="title">proxy-name</property>
                    <child>
                      <object class="GtkCellRendererText" id="proxy-name-column-rendered"/>
                      <attributes>
                        <attribute name="text">0</attribute>
                      </attributes>
                    </child>
                  </object>
                </child>
              </object>

              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">3</property>
              </packing>
            </child>

            <child>
              <object class="GtkButtonBox" id="groups-buttons">
                <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
                <child>
                  <object class="GtkMenuButton" id="add-btn">
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

              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">4</property>
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

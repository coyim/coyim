package definitions

func init() {
	add(`PeerDetails`, &defPeerDetails{})
}

type defPeerDetails struct{}

func (*defPeerDetails) String() string {
	return `<interface>
  <object class="GtkListStore" id="current-groups">
    <columns>
      <column type="gchararray"/>
    </columns>
  </object>

  <object class="GtkMenuItem" id="addGroup">
    <property name="visible">True</property>
    <property name="label" translatable="yes">New Group...</property>
    <signal name="activate" handler="on-add-new-group" />
  </object>

  <object class="GtkMenu" id="groups-menu">
    <property name="visible">True</property>
  </object>

  <object class="GtkDialog" id="dialog">
    <property name="title" translatable="yes">Edit contact</property>
    <signal name="close" handler="on-cancel" />
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
              <object class="GtkLabel" id="account-label">
                <property name="label" translatable="yes">Account</property>
                <property name="halign">GTK_ALIGN_END</property>
                <property name="selectable">TRUE</property>
              </object>
              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">0</property>
              </packing>
            </child>
            <child>
              <object class="GtkLabel" id="account-name">
                <property name="halign">GTK_ALIGN_START</property>
                <property name="selectable">TRUE</property>
              </object>
              <packing>
                <property name="left-attach">1</property>
                <property name="top-attach">0</property>
              </packing>
            </child>

            <child>
              <object class="GtkLabel" id="jid-label">
                <property name="label" translatable="yes">Contact</property>
                <property name="halign">GTK_ALIGN_END</property>
                <property name="selectable">TRUE</property>
              </object>
              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">1</property>
              </packing>
            </child>
            <child>
              <object class="GtkLabel" id="jid">
                <property name="halign">GTK_ALIGN_START</property>
                <property name="selectable">TRUE</property>
              </object>
              <packing>
                <property name="left-attach">1</property>
                <property name="top-attach">1</property>
              </packing>
            </child>

            <child>
              <object class="GtkLabel" id="nickname-label">
                <property name="label" translatable="yes">Nickname</property>
                <property name="halign">GTK_ALIGN_END</property>
              </object>
              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">2</property>
              </packing>
            </child>
            <child>
              <object class="GtkEntry" id="nickname">
                <property name="activates-default">True</property>
              </object>
              <packing>
                <property name="left-attach">1</property>
                <property name="top-attach">2</property>
              </packing>
            </child>

            <child>
              <object class="GtkLabel" id="require-encryption-label">
                <property name="label" translatable="yes">Require encryption with this peer</property>
                <property name="halign">GTK_ALIGN_END</property>
                <property name="selectable">TRUE</property>
              </object>
              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">3</property>
              </packing>
            </child>
            <child>
              <object class="GtkCheckButton" id="require-encryption"/>
              <packing>
                <property name="left-attach">1</property>
                <property name="top-attach">3</property>
              </packing>
            </child>


            <child>
              <object class="GtkLabel" id="groups">
                <property name="label" translatable="yes">Groups</property>
                <property name="halign">GTK_ALIGN_START</property>
                <property name="selectable">TRUE</property>
              </object>
              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">4</property>
                <property name="width">2</property>
              </packing>
            </child>
            <!-- TODO: This should be in a scrolled window with min/max size  -->
            <child>
              <object class="GtkTreeView" id="groups-view">
                <property name="model">current-groups</property>
                <property name="visible">True</property>
                <property name="can-focus">True</property>
                <property name="headers-visible">False</property>
                <property name="show-expanders">False</property>
                <property name="reorderable">True</property>
                <child internal-child="selection">
                  <object class="GtkTreeSelection" id="selection">
                    <property name="mode">GTK_SELECTION_SINGLE</property>
                    <signal name="changed" handler="on-group-selection-changed" swapped="no"/>
                  </object>
                </child>
                <child>
                  <object class="GtkTreeViewColumn" id="group-name-column">
                    <property name="title">group-name</property>
                    <child>
                      <object class="GtkCellRendererText" id="group-name-column-rendered"/>
                      <attributes>
                        <attribute name="text">0</attribute>
                      </attributes>
                    </child>
                  </object>
                </child>
              </object>

              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">5</property>
                <property name="width">2</property>
              </packing>
            </child>

            <child>
              <object class="GtkButtonBox" id="groups-buttons">
                <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
                <child>
                  <object class="GtkMenuButton" id="add-btn">
                    <property name="visible">true</property>
                    <property name="popup">groups-menu</property>

                    <property name="direction">GTK_ARROW_DOWN</property>
                    <child>
                      <object class="GtkImage" id="add-btn-img">
                        <property name="icon-name">list-add</property>
                      </object>
                    </child>
                  </object>
                </child>
                <child>
                  <object class="GtkButton" id="remove-btn">
                    <signal name="clicked" handler="on-remove-group"/>
                    <child>
                      <object class="GtkImage" id="remove-btn-img">
                        <property name="icon-name">list-remove</property>
                      </object>
                    </child>
                  </object>
                </child>
              </object>

              <packing>
                <property name="left-attach">0</property>
                <property name="top-attach">6</property>
              </packing>
            </child>
          </object>
        </child>

        <child internal-child="action_area">
          <object class="GtkButtonBox" id="button_box">
            <property name="orientation">GTK_ORIENTATION_HORIZONTAL</property>
            <child>
              <object class="GtkButton" id="btn-cancel">
                <property name="label" translatable="yes">Cancel</property>
                <signal name="clicked" handler="on-cancel"/>
              </object>
            </child>
            <child>
              <object class="GtkButton" id="btn-save">
                <property name="label" translatable="yes">Save</property>
                <property name="can-default">true</property>
                <signal name="clicked" handler="on-save"/>
              </object>
            </child>
          </object>
        </child>

       <!-- TODO: This should be in a scrolled window with min/max size  -->
       <child>
        <object class="GtkBox" id="box">
          <property name="border-width">10</property>
          <property name="homogeneous">false</property>
          <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
          <child>
            <object class="GtkLabel" id="fingerprintsInformation">
              <property name="selectable">true</property>
            </object>
          </child>
          <child>
            <object class="GtkGrid" id="fingerprintsGrid">
              <property name="margin-top">15</property>
              <property name="margin-bottom">10</property>
              <property name="margin-start">10</property>
              <property name="margin-end">10</property>
              <property name="row-spacing">12</property>
              <property name="column-spacing">6</property>
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

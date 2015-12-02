package definitions

func init() {
	add(`ContactPopupMenu`, &defContactPopupMenu{})
}

type defContactPopupMenu struct{}

func (*defContactPopupMenu) String() string {
	return `
<interface>
  <object class="GtkMenu" id="contactMenu">
    <child>
      <object class="GtkMenuItem" id="removeContactMenuItem">
        <property name="label" translatable="yes">Remove</property>
        <signal name="activate" handler="on_remove_contact" />
      </object>
    </child>
    <child>
      <object class="GtkSeparatorMenuItem" id="sep1"/>
    </child>
    <child>
      <object class="GtkMenuItem" id="menuItemForVisibilitySection">
        <child>
          <object class="GtkLabel" id="labelForVisibilitySection">
            <property name="label" translatable="yes">Visibility</property>
            <property name="sensitive">false</property>
          </object>
        </child>
      </object>
    </child>
    <child>
      <object class="GtkMenuItem" id="allowContactToSeeStatusMenuItem">
        <property name="label" translatable="yes">Allow contact to see my status</property>
        <signal name="activate" handler="on_allow_contact_to_see_status" />
      </object>
    </child>
    <child>
      <object class="GtkMenuItem" id="forbidContactToSeeStatusMenuItem">
        <property name="label" translatable="yes">Forbid contact to see my status</property>
        <signal name="activate" handler="on_forbid_contact_to_see_status" />
      </object>
    </child>
    <child>
      <object class="GtkMenuItem" id="askContactToSeeStatusMenuItem">
        <property name="label" translatable="yes">Ask contact to see their status</property>
        <signal name="activate" handler="on_ask_contact_to_see_status" />
      </object>
    </child>
    <child>
      <object class="GtkSeparatorMenuItem" id="sep2"/>
    </child>
    <child>
      <object class="GtkMenuItem" id="dumpInfoMenuItem">
        <property name="label" translatable="yes">Dump info</property>
        <signal name="activate" handler="on_dump_info" />
      </object>
    </child>
  </object>
</interface>

`
}

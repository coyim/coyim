package definitions

func init() {
	add(`Roster`, &defRoster{})
}

type defRoster struct{}

func (*defRoster) String() string {
	return `
<interface>
  <object class="GtkTreeStore" id="roster-model">
    <columns>
      <!-- jid -->
      <column type="gchararray"/>
      <!-- display name -->
      <column type="gchararray"/>
      <!-- account id -->
      <column type="gchararray"/>
      <!-- status color -->
      <column type="gchararray"/>
      <!-- background color -->
      <column type="gchararray"/>
      <!-- weight of font -->
      <column type="gint"/>
      <!-- tooltip -->
      <column type="gchararray"/>
      <!-- icon -->
      <column type="GdkPixbuf"/>
    </columns>
  </object>
  <object class="GtkScrolledWindow" id="roster">
    <property name="hscrollbar-policy">GTK_POLICY_NEVER</property>
    <property name="vscrollbar-policy">GTK_POLICY_AUTOMATIC</property>
    <child>
      <object class="GtkTreeView" id="roster-view">
        <property name="model">roster-model</property>
        <property name="headers-visible">false</property>
        <property name="show-expanders">false</property>
        <property name="level-indentation">3</property>
        <!-- TODO remove magic number -->
        <property name="tooltip-column">6</property>
        <signal name="row-activated" handler="on_activate_buddy" />
        <child internal-child="selection">
          <object class="GtkTreeSelection" id="selection">
            <property name="mode">GTK_SELECTION_SINGLE</property>
          </object>
        </child>
        <child>
          <object class="GtkTreeViewColumn" id="icon-column">
            <child>
              <object class="GtkCellRendererPixbuf" id="icon-column-rendered"/>
              <attributes>
                <attribute name="pixbuf">7</attribute>
                <attribute name="cell-background">4</attribute>
              </attributes>
            </child>
          </object>
        </child>
        <child>
          <object class="GtkTreeViewColumn" id="name-column">
            <property name="title">name</property>
            <child>
              <object class="GtkCellRendererText" id="name-column-rendered"/>
              <attributes>
                <!-- TODO remove magic numbers -->
                <attribute name="text">1</attribute>
                <attribute name="foreground">3</attribute>
                <attribute name="background">4</attribute>
                <attribute name="weight">5</attribute>
              </attributes>
            </child>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>

`
}

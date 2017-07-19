package definitions

func init() {
	add(`TorInstallHelper`, &defTorInstallHelper{})
}

type defTorInstallHelper struct{}

func (*defTorInstallHelper) String() string {
	return `<interface>
  <!-- change this to GtkDialog so text can be selectable-->
  <object class="GtkMessageDialog" id="dialog">
    <property name="window-position">GTK_WIN_POS_CENTER_ALWAYS</property>
    <property name="modal">true</property>
    <property name="border_width">7</property>
    <property name="text" translatable="yes">Install Tor</property>
    <property name="secondary_text" translatable="yes">

    This should be easy:

    1. Go to https://www.torproject.org/

    2. Download and install Tor

    3. Start Tor or open Tor Browser

    4. Close CoyIM

    5. Reopen CoyIM

    6. Enjoy!

    </property>
    <property name="buttons">GTK_BUTTONS_CLOSE</property>
  </object>
</interface>
`
}

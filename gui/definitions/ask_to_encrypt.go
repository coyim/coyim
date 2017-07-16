package definitions

func init() {
	add(`AskToEncrypt`, &defAskToEncrypt{})
}

type defAskToEncrypt struct{}

func (*defAskToEncrypt) String() string {
	return `<interface>
  <object class="GtkMessageDialog" id="AskToEncrypt">
    <property name="window-position">GTK_WIN_POS_CENTER</property>
    <property name="modal">true</property>
    <property name="border_width">7</property>
    <property name="title" translatable="yes">Encrypt your account's information</property>
    <property name="text" translatable="yes">
Would you like to encrypt your account's information?
    </property>
    <property name="secondary_text" translatable="yes">You will not be able to access your account's information file if you lose the

password. You will only be asked for it when CoyIM starts.
    </property>
    <property name="buttons">GTK_BUTTONS_YES_NO</property>
  </object>
</interface>
`
}

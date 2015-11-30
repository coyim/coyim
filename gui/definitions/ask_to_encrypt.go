package definitions

func init() {
	add(`AskToEncrypt`, &defAskToEncrypt{})
}

type defAskToEncrypt struct{}

func (*defAskToEncrypt) String() string {
	return `
<interface>
  <object class="GtkMessageDialog" id="AskToEncrypt">
    <property name="window-position">1</property>
    <property name="modal">true</property>
    <property name="title" translatable="yes">Encrypt configuration file</property>
    <property name="text" translatable="yes">Would you like to save your configuration file in an encrypted format? This can be significantly more secure, but you will not be able to access your configuration if you lose the password. You will only be asked for your password when CoyIM starts.</property>
    <property name="buttons">GTK_BUTTONS_YES_NO</property>
  </object>
</interface>

`
}

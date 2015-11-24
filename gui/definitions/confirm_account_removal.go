
package definitions

func init(){
  add(`ConfirmAccountRemoval`, &defConfirmAccountRemoval{})
}

type defConfirmAccountRemoval struct{}

func (*defConfirmAccountRemoval) String() string {
	return `
<interface>
  <object class="GtkMessageDialog" id="RemoveAccount">
    <property name="window-position">1</property>
    <property name="title">Confirm account removal</property>
    <property name="text">Are you sure?</property>
    <property name="buttons">GTK_BUTTONS_OK_CANCEL</property>
  </object>
</interface>

`
}

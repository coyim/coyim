
package definitions

func init(){
  add(`ConfigAssistantDefinition`, &defConfigAssistantDefinition{})
}

type defConfigAssistantDefinition struct{}

func (*defConfigAssistantDefinition) String() string {
	return `
<?xml version="1.0"?>
<interface>
  <object class="GtkAssistant" id="assistant">
    <property name="visible">true</property>
    <property name="width_request">450</property>
    <property name="height_request">300</property>
    <property name="border_width">10</property>
    <signal name="destroy" handler="close-assistant" />
    <signal name="cancel" handler="close-assistant" />
    <signal name="apply" handler="create-account" />
    <child internal-child="accessible">
      <object class="AtkObject" id="ConfigurationName">
        <property name="accessible-name" translatable="yes">Configuration Assistant</property>
      </object>
    </child>
    <child>
      <object class="GtkLabel" id="welcome">
        <property name="visible">true</property>
        <property name="wrap">true</property>
        <property name="label" translatable="yes">Welcome to CoyIM, the safe and secure xmpp client.</property>
      </object>
      <packing>
        <property name="page_type">GTK_ASSISTANT_PAGE_INTRO</property>
      </packing>
    </child>
    <child>
      <object class="GtkBox" id="detecting-tor">
        <property name="visible">true</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <signal name="realize" handler="detect-tor" />
        <child>
          <object class="GtkLabel" id="TorLabel">
            <property name="visible">true</property>
            <property name="wrap">true</property>
            <property name="label" translatable="yes">Detecting Tor...</property>
          </object>
        </child>
        <child>
          <object class="GtkLabel" id="tor-detected-msg">
            <property name="visible">false</property>
            <property name="wrap">true</property>
            <property name="label" translatable="yes">Tor detected successfully. Continue.</property>
          </object>
        </child>
        <child>
          <object class="GtkLabel" id="tor-not-detected-msg">
            <property name="visible">false</property>
            <property name="wrap">true</property>
            <property name="label" translatable="yes">Failed to detect Tor. Make sure it is running and try again.</property>
          </object>
        </child>
        <!-- TODO: Add a retry button -->
      </object>
      <packing>
        <property name="page_type">GTK_ASSISTANT_PAGE_PROGRESS</property>
      </packing>
    </child>
    <child>
      <object class="GtkBox" id="account-details">
        <property name="visible">true</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <child>
          <object class="GtkLabel" id="AccountLabel">
            <property name="visible">true</property>
            <property name="wrap">true</property>
            <property name="label" translatable="yes">Your account (for example: kim42@dukgo.com)</property>
          </object>
        </child>
        <child>
          <object class="GtkEntry" id="account">
            <property name="visible">true</property>
            <property name="input-purpose">GTK_INPUT_PURPOSE_EMAIL</property>
            <signal name="key-release-event" handler="xmpp-id-changed" />
          </object>
        </child>
        <child>
          <object class="GtkLabel" id="PasswordLabel">
            <property name="visible">true</property>
            <property name="wrap">true</property>
            <property name="label" translatable="yes">Password (optional)</property>
          </object>
        </child>
        <child>
          <object class="GtkEntry" id="password">
            <property name="visible">true</property>
            <property name="visibility">false</property>
            <property name="input-purpose">GTK_INPUT_PURPOSE_PASSWORD</property>
          </object>
        </child>
      </object>
      <packing>
        <property name="page_type">GTK_ASSISTANT_PAGE_CONTENT</property>
      </packing>
    </child>
    <child>
      <object class="GtkBox" id="detecting-xmpp-server">
        <property name="visible">true</property>
        <property name="orientation">GTK_ORIENTATION_VERTICAL</property>
        <signal name="realize" handler="detect-xmpp-server" />
        <child>
          <object class="GtkLabel" id="DetectingLabel">
            <property name="visible">true</property>
            <property name="wrap">true</property>
            <property name="label" translatable="yes">Detecting XMPP server configuration...</property>
          </object>
        </child>
        <child>
          <object class="GtkLabel" id="xmpp-server-msg">
            <property name="visible">false</property>
            <property name="wrap">true</property>
            <property name="label" translatable="yes"></property>
          </object>
        </child>
      </object>
      <packing>
        <property name="page_type">GTK_ASSISTANT_PAGE_PROGRESS</property>
      </packing>
    </child>
    <child>
      <object class="GtkLabel" id="ApplyLabel">
        <property name="visible">True</property>
        <property name="label" translatable="yes">Click "Apply" to accept this configuration.</property>
      </object>
      <packing>
        <property name="page_type">GTK_ASSISTANT_PAGE_CONFIRM</property>
      </packing>
    </child>
  </object>
</interface>


`
}

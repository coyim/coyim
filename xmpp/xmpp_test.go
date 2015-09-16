package xmpp

import (
	"encoding/xml"
	"testing"
)

func TestDiscoReplyVerSimple(t *testing.T) {
	expect := "QgayPKawpkPSDYmwT/WM94uAlu0="
	input := []byte(`
  <query xmlns='http://jabber.org/protocol/disco#info'
         node='http://code.google.com/p/exodus#QgayPKawpkPSDYmwT/WM94uAlu0='>
    <identity category='client' name='Exodus 0.9.1' type='pc'/>
    <feature var='http://jabber.org/protocol/caps'/>
    <feature var='http://jabber.org/protocol/disco#info'/>
    <feature var='http://jabber.org/protocol/disco#items'/>
    <feature var='http://jabber.org/protocol/muc'/>
  </query>
  `)
	var dr DiscoveryReply
	if err := xml.Unmarshal(input, &dr); err != nil {
		t.Fatal(err)
	}
	hash, err := dr.VerificationString()
	if err != nil {
		t.Fatal(err)
	}
	if hash != expect {
		t.Fatalf("Expected: %s Got: %s", expect, hash)
	}
}

func TestDiscoReplyVerComplex(t *testing.T) {
	expect := "q07IKJEyjvHSyhy//CH0CxmKi8w="
	input := []byte(`
  <query xmlns='http://jabber.org/protocol/disco#info'
         node='http://psi-im.org#q07IKJEyjvHSyhy//CH0CxmKi8w='>
    <identity xml:lang='en' category='client' name='Psi 0.11' type='pc'/>
    <identity xml:lang='el' category='client' name='Ψ 0.11' type='pc'/>
    <feature var='http://jabber.org/protocol/caps'/>
    <feature var='http://jabber.org/protocol/disco#info'/>
    <feature var='http://jabber.org/protocol/disco#items'/>
    <feature var='http://jabber.org/protocol/muc'/>
    <x xmlns='jabber:x:data' type='result'>
      <field var='FORM_TYPE' type='hidden'>
        <value>urn:xmpp:dataforms:softwareinfo</value>
      </field>
      <field var='ip_version'>
        <value>ipv4</value>
        <value>ipv6</value>
      </field>
      <field var='os'>
        <value>Mac</value>
      </field>
      <field var='os_version'>
        <value>10.5.1</value>
      </field>
      <field var='software'>
        <value>Psi</value>
      </field>
      <field var='software_version'>
        <value>0.11</value>
      </field>
    </x>
  </query>
`)
	var dr DiscoveryReply
	if err := xml.Unmarshal(input, &dr); err != nil {
		t.Fatal(err)
	}
	hash, err := dr.VerificationString()
	if err != nil {
		t.Fatal(err)
	}
	if hash != expect {
		t.Fatalf("Expected: %s Got: %s", expect, hash)
	}
}

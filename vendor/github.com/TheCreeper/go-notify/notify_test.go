package notify

import "testing"

func TestGetCapabilities(t *testing.T) {
	c, err := GetCapabilities()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Support Action Icons: %v\n", c.ActionIcons)
	t.Logf("Support Actions: %v\n", c.Actions)
	t.Logf("Support Body: %v\n", c.Body)
	t.Logf("Support Body Hyperlinks: %v\n", c.BodyHyperlinks)
	t.Logf("Support Body Images: %v\n", c.BodyImages)
	t.Logf("Support Body Markup: %v\n", c.BodyMarkup)
	t.Logf("Support Icon Multi: %v\n", c.IconMulti)
	t.Logf("Support Icon Static: %v\n", c.IconStatic)
	t.Logf("Support Persistence: %v\n", c.Persistence)
	t.Logf("Support Sound: %v\n", c.Sound)
}

func TestGetServerInformation(t *testing.T) {
	info, err := GetServerInformation()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Server Name: %s\n", info.Name)
	t.Logf("Server Spec Version: %s\n", info.SpecVersion)
	t.Logf("Server Vendor: %s\n", info.Vendor)
	t.Logf("Sserver Version: %s\n", info.Version)
}

func TestNewNotification(t *testing.T) {
	ntf := NewNotification("Notification Test", "Just a test")
	if _, err := ntf.Show(); err != nil {
		t.Fatal(err)
	}
}

func TestCloseNotification(t *testing.T) {
	ntf := NewNotification("Notification Test", "Just a test")
	id, err := ntf.Show()
	if err != nil {
		t.Fatal(err)
	}

	if err = CloseNotification(id); err != nil {
		t.Fatal(err)
	}
}

func TestUrgencyNotification(t *testing.T) {
	ntfLow := NewNotification("Urgency Test", "Testing notification urgency low")
	ntfLow.Hints = make(map[string]interface{})

	ntfLow.Hints[HintUrgency] = UrgencyLow
	_, err := ntfLow.Show()
	if err != nil {
		t.Fatal(err)
	}

	ntfCritical := NewNotification("Urgency Test", "Testing notification urgency critical")
	ntfCritical.Hints = make(map[string]interface{})

	ntfCritical.Hints[HintUrgency] = UrgencyCritical
	_, err = ntfCritical.Show()
	if err != nil {
		t.Fatal(err)
	}
}

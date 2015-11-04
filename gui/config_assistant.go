package gui

import (
	"log"
	"net"
	"strings"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/i18n"
)

type saveAccountFunc func(*config.Account)

func (u *gtkUI) addAndSaveAccountConfig(c *config.Account) {
	u.config.Add(c)
	u.SaveConfig()
}

func (u *gtkUI) showConfigAssistant() error {
	assistant, err := buildConfigAssistant(u.addAndSaveAccountConfig)
	if err != nil {
		return err
	}

	assistant.Show()
	return nil
}

func buildConfigAssistant(saveFn saveAccountFunc) (*gtk.Assistant, error) {
	builder, err := loadBuilderWith("ConfigAssistantDefinition", nil)
	if err != nil {
		return nil, err
	}

	var obj glib.IObject
	obj, err = builder.GetObject("assistant")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	assistant := obj.(*gtk.Assistant)
	intro, _ := assistant.GetNthPage(0)
	assistant.SetPageComplete(intro, true)

	confirm, _ := assistant.GetNthPage(-1)
	assistant.SetPageComplete(confirm, true)

	obj, _ = builder.GetObject("account")
	accountEntry := obj.(*gtk.Entry)

	obj, _ = builder.GetObject("password")
	passwordEntry := obj.(*gtk.Entry)

	builder.ConnectSignals(map[string]interface{}{
		"detect-tor": func(page *gtk.Box) {
			log.Println("Detecting Tor")

			obj, err := builder.GetObject("tor-detected-msg")
			if err != nil {
				return
			}

			detectedMsg := obj.(*gtk.Label)

			obj, err = builder.GetObject("tor-not-detected-msg")
			if err != nil {
				return
			}

			notDetectedMsg := obj.(*gtk.Label)

			//detectedMsg.SetVisible(false)
			//notDetectedMsg.SetVisible(false)

			go func() {
				<-time.After(5 * time.Second) // just to simulate
				_, ok := config.DetectTor()

				glib.IdleAdd(func() {
					detectedMsg.SetVisible(ok)
					notDetectedMsg.SetVisible(!ok)
					assistant.SetPageComplete(page, ok)
				})
			}()
		},

		"xmpp-id-changed": func(entry *gtk.Entry) {
			page, _ := entry.GetParent()
			l := entry.GetTextLength()
			assistant.SetPageComplete(page, l > 0)
		},

		"detect-xmpp-server": func(page *gtk.Box) {
			obj, err := builder.GetObject("xmpp-server-msg")
			if err != nil {
				return
			}
			msgLabel := obj.(*gtk.Label)

			xmppID, err := accountEntry.GetText()
			if err != nil {
				return
			}

			parts := strings.Split(xmppID, "@")
			if len(parts) < 2 {
				//TODO: go back to previous page
				return
			}

			domain := parts[1]
			services, err := config.ResolveXMPPServerOverTor(domain)
			if err != nil {
				//TODO: some network/DNS failure. Should it show a retry option?
				return
			}

			if len(services) > 0 {
				msgLabel.SetVisible(true)
				msgLabel.SetText(i18n.Local("All right with SRV"))
				assistant.SetPageComplete(page, true)
				return
			}

			// Fallback to using the domain at default port
			// TODO: proxy.Dialer does not support DialTimeout
			addr := net.JoinHostPort(domain, "5222")
			torProxy, err := config.NewTorProxy()
			if err != nil {
				//TODO: how to recover from this?
				return
			}

			go func() {
				<-time.After(5 * time.Second) // just to simulate

				conn, err := torProxy.Dial("tcp", addr)
				defer conn.Close()

				glib.IdleAdd(func() {
					if err != nil {
						//TODO: Failed to connect, should ask for XMPP server (and port)
						msgLabel.SetVisible(true)
						msgLabel.SetText(i18n.Localf(
							"Could not detect XMPP server for %s. Please inform the server domain.", domain,
						))
						return
					}

					msgLabel.SetVisible(true)
					msgLabel.SetText(i18n.Local("All right with fallback"))
					assistant.SetPageComplete(page, true)
				})
			}()

		},

		"create-account": func() {
			c, err := config.NewAccount()
			if err != nil {
				return
			}

			c.Account, _ = accountEntry.GetText()
			c.Password, _ = passwordEntry.GetText()

			saveFn(c)
		},
	})

	return assistant, nil
}

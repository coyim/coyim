// +build !nocli

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
	"golang.org/x/crypto/ssh/terminal"
)

type cliUI struct {
	session *Session

	password string
	oldState *terminal.State
	term     *terminal.Terminal
}

func newCLI() *cliUI {
	var err error
	c := &cliUI{}

	if c.oldState, err = terminal.MakeRaw(0); err != nil {
		panic(err.Error())
	}

	c.term = terminal.NewTerminal(os.Stdin, "")
	updateTerminalSize(c.term)
	c.term.SetBracketedPasteMode(true)

	resizeChan := make(chan os.Signal)
	go func() {
		for _ = range resizeChan {
			updateTerminalSize(c.term)
		}
	}()
	signal.Notify(resizeChan, syscall.SIGWINCH)

	return c
}

func (c *cliUI) Loop() {
	quit := make(chan bool)

	//This should be done by any client
	info(c.term, "Fetching roster")
	rosterReply, _, err := c.session.conn.RequestRoster()
	if err != nil {
		c.Alert("Failed to request roster: " + err.Error())
		return
	}

	c.session.conn.SignalPresence("")
	c.session.input = Input{
		term:        c.term,
		uidComplete: new(priorityList),
	}

	ticker := time.NewTicker(1 * time.Second)
	go timeoutLoop(c.session, ticker.C)

	commandChan := make(chan interface{})
	go c.session.input.ProcessCommands(commandChan)

	stanzaChan := make(chan xmpp.Stanza)
	go c.session.readMessages(stanzaChan)

	go commandLoop(c.term, c.session, c.session.config, commandChan, quit)
	go stanzaLoop(c, c.term, c.session, stanzaChan, quit)
	go rosterLoop(c.term, c.session, rosterReply, quit)

	c.term.SetPrompt("> ")
	<-quit // wait
}

func (c *cliUI) Close() {
	if c.oldState != nil {
		terminal.Restore(0, c.oldState)
	}

	if c.term != nil {
		c.term.SetBracketedPasteMode(false)
	}
}

func (c *cliUI) Alert(m string) {
	alert(c.term, m)
}

func (c *cliUI) Enroll(config *Config) bool {
	return enroll(config, c.term)
}

func (c *cliUI) AskForPassword(config *Config) (string, error) {
	var err error
	c.password, err = c.term.ReadPassword(fmt.Sprintf("Password for %s (will not be saved to disk): ", config.Account))

	return c.password, err
}

func (c *cliUI) RegisterCallback() xmpp.FormCallback {
	if *createAccount {
		return func(title, instructions string, fields []interface{}) error {
			//TODO: why does this function needs the
			//TODO: get user from Config
			user := "xxxxxx"
			return promptForForm(c.term, user, c.password, title, instructions, fields)
		}
	}

	return nil
}

func (c *cliUI) ProcessPresence(stanza *xmpp.ClientPresence) {
	s := c.session
	gone := false

	switch stanza.Type {
	case "subscribe":
		// This is a subscription request
		jid := xmpp.RemoveResourceFromJid(stanza.From)
		info(s.term, jid+" wishes to see when you're online. Use '/confirm "+jid+"' to confirm (or likewise with /deny to decline)")
		s.pendingSubscribes[jid] = stanza.Id
		s.input.AddUser(jid)
		return
	case "unavailable":
		gone = true
	case "":
		break
	default:
		return
	}

	from := xmpp.RemoveResourceFromJid(stanza.From)

	if gone {
		if _, ok := s.knownStates[from]; !ok {
			// They've gone, but we never knew they were online.
			return
		}
		delete(s.knownStates, from)
	} else {
		if _, ok := s.knownStates[from]; !ok && isAwayStatus(stanza.Show) {
			// Skip people who are initially away.
			return
		}

		if lastState, ok := s.knownStates[from]; ok && lastState == stanza.Show {
			// No change. Ignore.
			return
		}
		s.knownStates[from] = stanza.Show
	}

	if !s.config.HideStatusUpdates {
		var line []byte
		line = append(line, []byte(fmt.Sprintf("   (%s) ", time.Now().Format(time.Kitchen)))...)
		line = append(line, s.term.Escape.Magenta...)
		line = append(line, []byte(from)...)
		line = append(line, ':')
		line = append(line, s.term.Escape.Reset...)
		line = append(line, ' ')
		if gone {
			line = append(line, []byte("offline")...)
		} else if len(stanza.Show) > 0 {
			line = append(line, []byte(stanza.Show)...)
		} else {
			line = append(line, []byte("online")...)
		}
		line = append(line, ' ')
		line = append(line, []byte(stanza.Status)...)
		line = append(line, '\n')
		s.term.Write(line)
	}
}

func main() {
	flag.Parse()

	ui := newCLI()
	defer ui.Close()

	//Terminal is necessary to print error messages and
	//to ask for password
	config, password, err := loadConfig(ui)
	if err != nil {
		return
	}

	logger := &lineLogger{ui.term, nil}

	// Act on configuration
	conn, err := NewXMPPConn(ui, config, password, ui.RegisterCallback(), logger)
	if err != nil {
		ui.Alert(err.Error())
		return
	}

	//TODO support one session per account
	ui.session = &Session{
		ui: ui,

		account:           config.Account,
		conn:              conn,
		term:              ui.term,
		conversations:     make(map[string]*otr3.Conversation),
		eh:                make(map[string]*eventHandler),
		knownStates:       make(map[string]string),
		privateKey:        new(otr3.PrivateKey),
		config:            config,
		pendingRosterChan: make(chan *rosterEdit),
		pendingSubscribes: make(map[string]string),
		lastActionTime:    time.Now(),
	}

	ui.session.privateKey.Parse(config.PrivateKey)
	ui.session.timeouts = make(map[xmpp.Cookie]time.Time)

	info(ui.term, fmt.Sprintf("Your fingerprint is %x", ui.session.privateKey.DefaultFingerprint()))

	ui.Loop()

	os.Stdout.Write([]byte("\n"))
}

// +build !nocli

package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	coyconf "github.com/twstrike/coyim/config"
	coyui "github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
	"golang.org/x/crypto/ssh/terminal"
)

type cliUI struct {
	session *Session

	password string
	oldState *terminal.State
	term     *terminal.Terminal
	input    *Input

	terminate chan bool
}

func newCLI() *cliUI {
	var err error
	c := &cliUI{
		terminate: make(chan bool),
	}

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

//TODO: This should receive something telling which Session/COnfig should be terminated if we have multiple accounts connected
func (c *cliUI) Disconnected() {
	c.terminate <- true
}

func (c *cliUI) Loop() {
	c.input = &Input{
		term:        c.term,
		uidComplete: new(priorityList),
	}

	go c.session.WatchTimeout()
	go c.session.WatchRosterEvents()

	commandChan := make(chan interface{})
	go c.input.ProcessCommands(commandChan)
	go commandLoop(c, commandChan, c.terminate)

	stanzaChan := make(chan xmpp.Stanza)
	go c.session.readMessages(stanzaChan)
	go stanzaLoop(c, c.session, stanzaChan, c.terminate)

	c.term.SetPrompt("> ")
	<-c.terminate // wait
}

func (c *cliUI) Close() {
	if c.oldState != nil {
		terminal.Restore(0, c.oldState)
	}

	if c.term != nil {
		c.term.SetBracketedPasteMode(false)
	}
}

func (c *cliUI) Info(m string) {
	info(c.term, m)
}

func (c *cliUI) Warn(m string) {
	warn(c.term, m)
}

func (c *cliUI) Alert(m string) {
	alert(c.term, m)
}

func (c *cliUI) Enroll(config *coyconf.Config) bool {
	return enroll(config, c.term)
}

func (c *cliUI) AskForPassword(config *coyconf.Config) (string, error) {
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
		info(c.term, jid+" wishes to see when you're online. Use '/confirm "+jid+"' to confirm (or likewise with /deny to decline)")
		s.pendingSubscribes[jid] = stanza.Id
		c.input.AddUser(jid)
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
		if _, ok := s.knownStates[from]; !ok && coyui.IsAwayStatus(stanza.Show) {
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
		line = append(line, c.term.Escape.Magenta...)
		line = append(line, []byte(from)...)
		line = append(line, ':')
		line = append(line, c.term.Escape.Reset...)
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
		c.term.Write(line)
	}
}

func (c *cliUI) IQReceived(uid string) {
	c.input.AddUser(uid)
}

func (c *cliUI) RosterReceived(roster []xmpp.RosterEntry) {
	for _, entry := range roster {
		c.input.AddUser(entry.Jid)
	}

	info(c.term, "Roster received")
}

func main() {
	flag.Parse()

	ui := newCLI()
	defer ui.Close()

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
		conversations:     make(map[string]*otr3.Conversation),
		eh:                make(map[string]*eventHandler),
		knownStates:       make(map[string]string),
		privateKey:        new(otr3.PrivateKey),
		config:            config,
		pendingRosterChan: make(chan *coyui.RosterEdit),
		pendingSubscribes: make(map[string]string),
		lastActionTime:    time.Now(),
		sessionHandler:    ui,
	}

	ui.session.privateKey.Parse(config.PrivateKey)
	ui.session.timeouts = make(map[xmpp.Cookie]time.Time)

	info(ui.term, fmt.Sprintf("Your fingerprint is %x", ui.session.privateKey.DefaultFingerprint()))

	ui.Loop()
	os.Stdout.Write([]byte("\n"))
}

func (c *cliUI) MessageReceived(from, timestamp string, encrypted bool, message []byte) {

	var line []byte
	if encrypted {
		line = append(line, c.term.Escape.Green...)
	} else {
		line = append(line, c.term.Escape.Red...)
	}

	t := fmt.Sprintf("(%s) %s: ", timestamp, from)
	line = append(line, []byte(t)...)
	line = append(line, c.term.Escape.Reset...)
	line = appendTerminalEscaped(line, coyui.StripHTML(message))
	line = append(line, '\n')
	if c.session.config.Bell {
		line = append(line, '\a')
	}

	c.term.Write(line)
}

func (c *cliUI) NewOTRKeys(uid string, conversation *otr3.Conversation) {
	c.input.SetPromptForTarget(uid, true)
	c.printConversationInfo(uid, conversation)
}

func (c *cliUI) OTREnded(uid string) {
	c.input.SetPromptForTarget(uid, false)
}

func (c *cliUI) printConversationInfo(uid string, conversation *otr3.Conversation) {
	s := c.session
	term := c.term

	fpr := conversation.GetTheirKey().DefaultFingerprint()
	fprUid := s.config.UserIdForFingerprint(fpr)
	info(term, fmt.Sprintf("  Fingerprint  for %s: %x", uid, fpr))
	info(term, fmt.Sprintf("  Session  ID  for %s: %x", uid, conversation.GetSSID()))
	if fprUid == uid {
		info(term, fmt.Sprintf("  Identity key for %s is verified", uid))
	} else if len(fprUid) > 1 {
		alert(term, fmt.Sprintf("  Warning: %s is using an identity key which was verified for %s", uid, fprUid))
	} else if s.config.HasFingerprint(uid) {
		critical(term, fmt.Sprintf("  Identity key for %s is incorrect", uid))
	} else {
		alert(term, fmt.Sprintf("  Identity key for %s is not verified. You should use /otr-auth or /otr-authqa or /otr-authoob to verify their identity", uid))
	}
}

// promptForForm runs an XEP-0004 form and collects responses from the user.
func promptForForm(term *terminal.Terminal, user, password, title, instructions string, fields []interface{}) error {
	info(term, "The server has requested the following information. Text that has come from the server will be shown in red.")

	// formStringForPrinting takes a string form the form and returns an
	// escaped version with codes to make it show as red.
	formStringForPrinting := func(s string) string {
		var line []byte

		line = append(line, term.Escape.Red...)
		line = appendTerminalEscaped(line, []byte(s))
		line = append(line, term.Escape.Reset...)
		return string(line)
	}

	write := func(s string) {
		term.Write([]byte(s))
	}

	var tmpDir string

	showMediaEntries := func(questionNumber int, medias [][]xmpp.Media) {
		if len(medias) == 0 {
			return
		}

		write("The following media blobs have been provided by the server with this question:\n")
		for i, media := range medias {
			for j, rep := range media {
				if j == 0 {
					write(fmt.Sprintf("  %d. ", i+1))
				} else {
					write("     ")
				}
				write(fmt.Sprintf("Data of type %s", formStringForPrinting(rep.MIMEType)))
				if len(rep.URI) > 0 {
					write(fmt.Sprintf(" at %s\n", formStringForPrinting(rep.URI)))
					continue
				}

				var fileExt string
				switch rep.MIMEType {
				case "image/png":
					fileExt = "png"
				case "image/jpeg":
					fileExt = "jpeg"
				}

				if len(tmpDir) == 0 {
					var err error
					if tmpDir, err = ioutil.TempDir("", "xmppclient"); err != nil {
						write(", but failed to create temporary directory in which to save it: " + err.Error() + "\n")
						continue
					}
				}

				filename := filepath.Join(tmpDir, fmt.Sprintf("%d-%d-%d", questionNumber, i, j))
				if len(fileExt) > 0 {
					filename = filename + "." + fileExt
				}
				out, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
				if err != nil {
					write(", but failed to create file in which to save it: " + err.Error() + "\n")
					continue
				}
				out.Write(rep.Data)
				out.Close()

				write(", saved in " + filename + "\n")
			}
		}

		write("\n")
	}

	var err error
	if len(title) > 0 {
		write(fmt.Sprintf("Title: %s\n", formStringForPrinting(title)))
	}
	if len(instructions) > 0 {
		write(fmt.Sprintf("Instructions: %s\n", formStringForPrinting(instructions)))
	}

	questionNumber := 0
	for _, field := range fields {
		questionNumber++
		write("\n")

		switch field := field.(type) {
		case *xmpp.FixedFormField:
			write(formStringForPrinting(field.Text))
			write("\n")
			questionNumber--

		case *xmpp.BooleanFormField:
			write(fmt.Sprintf("%d. %s\n\n", questionNumber, formStringForPrinting(field.Label)))
			showMediaEntries(questionNumber, field.Media)
			term.SetPrompt("Please enter yes, y, no or n: ")

		TryAgain:
			for {
				answer, err := term.ReadLine()
				if err != nil {
					return err
				}
				switch answer {
				case "yes", "y":
					field.Result = true
				case "no", "n":
					field.Result = false
				default:
					continue TryAgain
				}
				break
			}

		case *xmpp.TextFormField:
			switch field.Label {
			case "CAPTCHA web page":
				if strings.HasPrefix(field.Default, "http") {
					// This is a oddity of jabber.ccc.de and maybe
					// others. The URL for the capture is provided
					// as the default answer to a question. Perhaps
					// that was needed with some clients. However,
					// we support embedded media and it's confusing
					// to ask the question, so we just print the
					// URL.
					write(fmt.Sprintf("CAPTCHA web page (only if not provided below): %s\n", formStringForPrinting(field.Default)))
					questionNumber--
					continue
				}

			case "User":
				field.Result = user
				questionNumber--
				continue

			case "Password":
				field.Result = password
				questionNumber--
				continue
			}

			write(fmt.Sprintf("%d. %s\n\n", questionNumber, formStringForPrinting(field.Label)))
			showMediaEntries(questionNumber, field.Media)

			if len(field.Default) > 0 {
				write(fmt.Sprintf("Please enter response or leave blank for the default, which is '%s'\n", formStringForPrinting(field.Default)))
			} else {
				write("Please enter response")
			}
			term.SetPrompt("> ")
			if field.Private {
				field.Result, err = term.ReadPassword("> ")
			} else {
				field.Result, err = term.ReadLine()
			}
			if err != nil {
				return err
			}
			if len(field.Result) == 0 {
				field.Result = field.Default
			}

		case *xmpp.MultiTextFormField:
			write(fmt.Sprintf("%d. %s\n\n", questionNumber, formStringForPrinting(field.Label)))
			showMediaEntries(questionNumber, field.Media)

			write("Please enter one or more responses, terminated by an empty line\n")
			term.SetPrompt("> ")

			for {
				line, err := term.ReadLine()
				if err != nil {
					return err
				}
				if len(line) == 0 {
					break
				}
				field.Results = append(field.Results, line)
			}

		case *xmpp.SelectionFormField:
			write(fmt.Sprintf("%d. %s\n\n", questionNumber, formStringForPrinting(field.Label)))
			showMediaEntries(questionNumber, field.Media)

			for i, opt := range field.Values {
				write(fmt.Sprintf("  %d. %s\n\n", i+1, formStringForPrinting(opt)))
			}
			term.SetPrompt("Please enter the number of your selection: ")

		TryAgain2:
			for {
				answer, err := term.ReadLine()
				if err != nil {
					return err
				}
				answerNum, err := strconv.Atoi(answer)
				answerNum--
				if err != nil || answerNum < 0 || answerNum >= len(field.Values) {
					write("Cannot parse that reply. Try again.")
					continue TryAgain2
				}

				field.Result = answerNum
				break
			}

		case *xmpp.MultiSelectionFormField:
			write(fmt.Sprintf("%d. %s\n\n", questionNumber, formStringForPrinting(field.Label)))
			showMediaEntries(questionNumber, field.Media)

			for i, opt := range field.Values {
				write(fmt.Sprintf("  %d. %s\n\n", i+1, formStringForPrinting(opt)))
			}
			term.SetPrompt("Please enter the numbers of zero or more of the above, separated by spaces: ")

		TryAgain3:
			for {
				answer, err := term.ReadLine()
				if err != nil {
					return err
				}

				var candidateResults []int
				answers := strings.Fields(answer)
				for _, answerStr := range answers {
					answerNum, err := strconv.Atoi(answerStr)
					answerNum--
					if err != nil || answerNum < 0 || answerNum >= len(field.Values) {
						write("Cannot parse that reply. Please try again.")
						continue TryAgain3
					}
					for _, other := range candidateResults {
						if answerNum == other {
							write("Cannot have duplicates. Please try again.")
							continue TryAgain3
						}
					}
					candidateResults = append(candidateResults, answerNum)
				}

				field.Results = candidateResults
				break
			}
		}
	}

	if len(tmpDir) > 0 {
		os.RemoveAll(tmpDir)
	}

	return nil
}

func commandLoop(ui *cliUI, commandChan <-chan interface{}, done chan<- bool) {
	var err error

	term := ui.term
	s := ui.session
	config := ui.session.config

CommandLoop:
	for {
		select {
		case cmd, ok := <-commandChan:
			if !ok {
				warn(term, "Exiting because command channel closed")
				break CommandLoop
			}
			s.lastActionTime = time.Now()
			switch cmd := cmd.(type) {
			case quitCommand:
				for to, conversation := range s.conversations {
					msgs, err := conversation.End()
					if err != nil {
						//TODO: error handle
						panic("this should not happen")
					}
					for _, msg := range msgs {
						s.conn.Send(to, string(msg))
					}
				}
				break CommandLoop
			case versionCommand:
				replyChan, cookie, err := s.conn.SendIQ(cmd.User, "get", xmpp.VersionQuery{})
				if err != nil {

					alert(ui.term, "Error sending version request: "+err.Error())
					continue
				}
				s.timeouts[cookie] = time.Now().Add(5 * time.Second)
				go s.awaitVersionReply(replyChan, cmd.User)
			case rosterCommand:
				info(ui.term, "Current roster:")
				maxLen := 0
				for _, item := range s.roster {
					if maxLen < len(item.Jid) {
						maxLen = len(item.Jid)
					}
				}

				for _, item := range s.roster {
					state, ok := s.knownStates[item.Jid]

					line := ""
					if ok {
						line += "[*] "
					} else if cmd.OnlineOnly {
						continue
					} else {
						line += "[ ] "
					}

					line += item.Jid
					numSpaces := 1 + (maxLen - len(item.Jid))
					for i := 0; i < numSpaces; i++ {
						line += " "
					}
					line += item.Subscription + "\t" + item.Name
					if ok {
						line += "\t" + state
					}
					info(ui.term, line)
				}
			case rosterEditCommand:
				if s.pendingRosterEdit != nil {
					warn(ui.term, "Aborting previous roster edit")
					s.pendingRosterEdit = nil
				}
				rosterCopy := make([]xmpp.RosterEntry, len(s.roster))
				copy(rosterCopy, s.roster)
				go s.editRoster(rosterCopy)
			case rosterEditDoneCommand:
				if s.pendingRosterEdit == nil {
					warn(ui.term, "No roster edit in progress. Use /rosteredit to start one")
					continue
				}
				go s.loadEditedRoster(*s.pendingRosterEdit)
			case toggleStatusUpdatesCommand:
				s.config.HideStatusUpdates = !s.config.HideStatusUpdates
				s.config.Save()
				// Tell the user the current state of the statuses
				if s.config.HideStatusUpdates {
					info(ui.term, "Status updates disabled")
				} else {
					info(ui.term, "Status updates enabled")
				}
			case confirmCommand:
				s.handleConfirmOrDeny(cmd.User, true /* confirm */)
			case denyCommand:
				s.handleConfirmOrDeny(cmd.User, false /* deny */)
			case addCommand:
				s.conn.SendPresence(cmd.User, "subscribe", "" /* generate id */)
			case msgCommand:
				conversation, ok := s.conversations[cmd.to]
				isEncrypted := ok && conversation.IsEncrypted()
				if cmd.setPromptIsEncrypted != nil {
					cmd.setPromptIsEncrypted <- isEncrypted
				}
				if !isEncrypted && config.ShouldEncryptTo(cmd.to) {
					warn(ui.term, fmt.Sprintf("Did not send: no encryption established with %s", cmd.to))
					continue
				}
				var msgs [][]byte
				message := []byte(cmd.msg)
				// Automatically tag all outgoing plaintext
				// messages with a whitespace tag that
				// indicates that we support OTR.
				if config.OTRAutoAppendTag &&
					!bytes.Contains(message, []byte("?OTR")) &&
					(!ok || !conversation.IsEncrypted()) {
					message = append(message, coyui.OTRWhitespaceTag...)
				}
				if ok {
					var err error
					validMsgs, err := conversation.Send(message)
					msgs = otr3.Bytes(validMsgs)
					if err != nil {
						alert(ui.term, err.Error())
						break
					}
				} else {
					msgs = [][]byte{[]byte(message)}
				}
				for _, message := range msgs {
					s.conn.Send(cmd.to, string(message))
				}
			case otrCommand:
				s.conn.Send(string(cmd.User), QueryMessage)
			case otrInfoCommand:
				info(term, fmt.Sprintf("Your OTR fingerprint is %x", s.privateKey.DefaultFingerprint()))
				for to, conversation := range s.conversations {
					if conversation.IsEncrypted() {
						info(ui.term, fmt.Sprintf("Secure session with %s underway:", to))
						ui.printConversationInfo(to, conversation)
					}
				}
			case endOTRCommand:
				to := string(cmd.User)
				conversation, ok := s.conversations[to]
				if !ok {
					alert(ui.term, "No secure session established")
					break
				}
				msgs, err := conversation.End()
				if err != nil {
					//TODO: error handle
					panic("this should not happen")
				}
				for _, msg := range msgs {
					s.conn.Send(to, string(msg))
				}
				ui.input.SetPromptForTarget(cmd.User, false)
				warn(ui.term, "OTR conversation ended with "+cmd.User)
			case authQACommand:
				to := string(cmd.User)
				conversation, ok := s.conversations[to]
				if !ok {
					alert(ui.term, "Can't authenticate without a secure conversation established")
					break
				}
				var ret []otr3.ValidMessage
				if s.eh[to].waitingForSecret {
					s.eh[to].waitingForSecret = false
					ret, err = conversation.ProvideAuthenticationSecret([]byte(cmd.Secret))
				} else {
					ret, err = conversation.StartAuthenticate(cmd.Question, []byte(cmd.Secret))
				}
				msgs := otr3.Bytes(ret)
				if err != nil {
					alert(ui.term, "Error while starting authentication with "+to+": "+err.Error())
				}
				for _, msg := range msgs {
					s.conn.Send(to, string(msg))
				}
			case authOobCommand:
				fpr, err := hex.DecodeString(cmd.Fingerprint)
				if err != nil {
					alert(ui.term, fmt.Sprintf("Invalid fingerprint %s - not authenticated", cmd.Fingerprint))
					break
				}
				existing := s.config.UserIdForFingerprint(fpr)
				if len(existing) != 0 {
					alert(ui.term, fmt.Sprintf("Fingerprint %s already belongs to %s", cmd.Fingerprint, existing))
					break
				}
				s.config.KnownFingerprints = append(s.config.KnownFingerprints, coyconf.KnownFingerprint{Fingerprint: fpr, UserId: cmd.User})
				s.config.Save()
				info(ui.term, fmt.Sprintf("Saved manually verified fingerprint %s for %s", cmd.Fingerprint, cmd.User))
			case awayCommand:
				s.conn.SignalPresence("away")
			case chatCommand:
				s.conn.SignalPresence("chat")
			case dndCommand:
				s.conn.SignalPresence("dnd")
			case xaCommand:
				s.conn.SignalPresence("xa")
			case onlineCommand:
				s.conn.SignalPresence("")
			}
		}
	}

	done <- true
}

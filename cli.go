// +build !nocli

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/event"
	"github.com/twstrike/coyim/session"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
	"github.com/twstrike/otr3"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/proxy"
)

type cliUI struct {
	config  *config.Config
	session *session.Session

	password string
	oldState *terminal.State
	term     *terminal.Terminal
	input    *Input

	terminate chan bool
}

func newCLI() *cliUI {
	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err.Error())
	}

	term := terminal.NewTerminal(os.Stdin, "")
	updateTerminalSize(term)
	term.SetBracketedPasteMode(true)

	resizeChan := make(chan os.Signal)
	go func() {
		for _ = range resizeChan {
			updateTerminalSize(term)
		}
	}()
	signal.Notify(resizeChan, syscall.SIGWINCH)

	return &cliUI{
		term:      term,
		oldState:  oldState,
		terminate: make(chan bool),
		input: &Input{
			term:        term,
			uidComplete: new(priorityList),
		},
	}
}

//TODO: This should receive something telling which Session/COnfig should be terminated if we have multiple accounts connected
func (c *cliUI) Disconnected() {
	c.terminate <- true
}

func (c *cliUI) Loop() {
	go c.session.WatchTimeout()
	go c.session.WatchRosterEvents()
	go c.session.WatchStanzas()
	go c.WatchCommands()

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

func (c *cliUI) RegisterCallback() xmpp.FormCallback {
	if *config.CreateAccount {
		return func(title, instructions string, fields []interface{}) error {
			user := c.config.Account
			return promptForForm(c.term, user, c.password, title, instructions, fields)
		}
	}

	return nil
}

func (c *cliUI) ProcessPresence(stanza *xmpp.ClientPresence, ignore, gone bool) {

	from := xmpp.RemoveResourceFromJid(stanza.From)

	if stanza.Type == "subscribe" {
		info(c.term, from+" wishes to see when you're online. Use '/confirm "+from+"' to confirm (or likewise with /deny to decline)")
		c.input.AddUser(from)
		return
	}

	s := c.session
	if ignore || s.Config.HideStatusUpdates {
		return
	}

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

	u := newCLI()
	defer u.Close()

	var err error
	u.config, err = config.Load(*config.ConfigFile)
	if err != nil {
		filename, e := config.FindConfigFile(os.Getenv("HOME"))
		if e != nil {
			//TODO cant write config file. Should it be a problem?
			return
		}

		u.config = config.NewConfig()
		u.config.Filename = *filename
		u.Alert(err.Error())
		enroll(u.config, u.term)
	}

	//TODO We do not support empty passwords
	var password string
	if len(u.config.Password) == 0 {
		var err error

		password, err = u.term.ReadPassword(
			fmt.Sprintf("Password for %s (will not be saved to disk): ", u.config.Account),
		)
		if err != nil {
			return
		}
	} else {
		password = u.config.Password
	}

	logger := &lineLogger{u.term, nil}

	// Act on configuration
	conn, err := config.NewXMPPConn(u.config, password, u.RegisterCallback(), logger)
	if err != nil {
		u.Alert(err.Error())
		return
	}

	//TODO support one session per account
	u.session = &session.Session{
		//WHY both?
		Account: u.config.Account,
		Config:  u.config,

		Conn:                conn,
		Conversations:       make(map[string]*otr3.Conversation),
		OtrEventHandler:     make(map[string]*event.OtrEventHandler),
		KnownStates:         make(map[string]string),
		PrivateKey:          new(otr3.PrivateKey),
		PendingRosterChan:   make(chan *ui.RosterEdit),
		PendingSubscribes:   make(map[string]string),
		LastActionTime:      time.Now(),
		SessionEventHandler: u,
	}

	u.session.PrivateKey.Parse(u.config.PrivateKey)
	u.session.Timeouts = make(map[xmpp.Cookie]time.Time)

	info(u.term, fmt.Sprintf("Your fingerprint is %x", u.session.PrivateKey.DefaultFingerprint()))

	u.Loop()
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
	line = appendTerminalEscaped(line, ui.StripHTML(message))
	line = append(line, '\n')
	if c.session.Config.Bell {
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
	fprUid := s.Config.UserIdForFingerprint(fpr)
	info(term, fmt.Sprintf("  Fingerprint  for %s: %x", uid, fpr))
	info(term, fmt.Sprintf("  Session  ID  for %s: %x", uid, conversation.GetSSID()))
	if fprUid == uid {
		info(term, fmt.Sprintf("  Identity key for %s is verified", uid))
	} else if len(fprUid) > 1 {
		alert(term, fmt.Sprintf("  Warning: %s is using an identity key which was verified for %s", uid, fprUid))
	} else if s.Config.HasFingerprint(uid) {
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

func (c *cliUI) WatchCommands() {
	defer c.Disconnected()

	commandChan := make(chan interface{})
	go c.input.ProcessCommands(commandChan)
	c.term.SetPrompt("> ")

	term := c.term
	s := c.session
	conf := s.Config

	var err error

CommandLoop:
	for {
		select {
		case cmd, ok := <-commandChan:
			if !ok {
				warn(term, "Exiting because command channel closed")
				break CommandLoop
			}
			s.LastActionTime = time.Now()
			switch cmd := cmd.(type) {
			case quitCommand:
				for to, conversation := range s.Conversations {
					msgs, err := conversation.End()
					if err != nil {
						//TODO: error handle
						panic("this should not happen")
					}
					for _, msg := range msgs {
						s.Conn.Send(to, string(msg))
					}
				}
				break CommandLoop
			case versionCommand:
				replyChan, cookie, err := s.Conn.SendIQ(cmd.User, "get", xmpp.VersionQuery{})
				if err != nil {

					alert(term, "Error sending version request: "+err.Error())
					continue
				}
				s.Timeouts[cookie] = time.Now().Add(5 * time.Second)
				go s.AwaitVersionReply(replyChan, cmd.User)
			case rosterCommand:
				info(term, "Current roster:")
				maxLen := 0
				for _, item := range s.Roster {
					if maxLen < len(item.Jid) {
						maxLen = len(item.Jid)
					}
				}

				for _, item := range s.Roster {
					state, ok := s.KnownStates[item.Jid]

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
					info(term, line)
				}
			case rosterEditCommand:
				if s.PendingRosterEdit != nil {
					warn(term, "Aborting previous roster edit")
					s.PendingRosterEdit = nil
				}
				rosterCopy := make([]xmpp.RosterEntry, len(s.Roster))
				copy(rosterCopy, s.Roster)
				go s.EditRoster(rosterCopy)
			case rosterEditDoneCommand:
				if s.PendingRosterEdit == nil {
					warn(term, "No roster edit in progress. Use /rosteredit to start one")
					continue
				}
				go s.LoadEditedRoster(*s.PendingRosterEdit)
			case toggleStatusUpdatesCommand:
				s.Config.HideStatusUpdates = !s.Config.HideStatusUpdates
				s.Config.Save()
				// Tell the user the current state of the statuses
				if s.Config.HideStatusUpdates {
					info(term, "Status updates disabled")
				} else {
					info(term, "Status updates enabled")
				}
			case confirmCommand:
				s.HandleConfirmOrDeny(cmd.User, true /* confirm */)
			case denyCommand:
				s.HandleConfirmOrDeny(cmd.User, false /* deny */)
			case addCommand:
				s.Conn.SendPresence(cmd.User, "subscribe", "" /* generate id */)
			case msgCommand:
				conversation, ok := s.Conversations[cmd.to]
				isEncrypted := ok && conversation.IsEncrypted()
				if cmd.setPromptIsEncrypted != nil {
					cmd.setPromptIsEncrypted <- isEncrypted
				}
				if !isEncrypted && conf.ShouldEncryptTo(cmd.to) {
					warn(term, fmt.Sprintf("Did not send: no encryption established with %s", cmd.to))
					continue
				}
				var msgs [][]byte
				message := []byte(cmd.msg)
				// Automatically tag all outgoing plaintext
				// messages with a whitespace tag that
				// indicates that we support OTR.
				if conf.OTRAutoAppendTag &&
					!bytes.Contains(message, []byte("?OTR")) &&
					(!ok || !conversation.IsEncrypted()) {
					message = append(message, ui.OTRWhitespaceTag...)
				}
				if ok {
					var err error
					validMsgs, err := conversation.Send(message)
					msgs = otr3.Bytes(validMsgs)
					if err != nil {
						alert(term, err.Error())
						break
					}
				} else {
					msgs = [][]byte{[]byte(message)}
				}
				for _, message := range msgs {
					s.Conn.Send(cmd.to, string(message))
				}
			case otrCommand:
				s.Conn.Send(string(cmd.User), event.QueryMessage)
			case otrInfoCommand:
				info(term, fmt.Sprintf("Your OTR fingerprint is %x", s.PrivateKey.DefaultFingerprint()))
				for to, conversation := range s.Conversations {
					if conversation.IsEncrypted() {
						info(term, fmt.Sprintf("Secure session with %s underway:", to))
						c.printConversationInfo(to, conversation)
					}
				}
			case endOTRCommand:
				to := string(cmd.User)
				conversation, ok := s.Conversations[to]
				if !ok {
					alert(term, "No secure session established")
					break
				}
				msgs, err := conversation.End()
				if err != nil {
					//TODO: error handle
					panic("this should not happen")
				}
				for _, msg := range msgs {
					s.Conn.Send(to, string(msg))
				}
				c.input.SetPromptForTarget(cmd.User, false)
				warn(term, "OTR conversation ended with "+cmd.User)
			case authQACommand:
				to := string(cmd.User)
				conversation, ok := s.Conversations[to]
				if !ok {
					alert(term, "Can't authenticate without a secure conversation established")
					break
				}
				var ret []otr3.ValidMessage
				if s.OtrEventHandler[to].WaitingForSecret {
					s.OtrEventHandler[to].WaitingForSecret = false
					ret, err = conversation.ProvideAuthenticationSecret([]byte(cmd.Secret))
				} else {
					ret, err = conversation.StartAuthenticate(cmd.Question, []byte(cmd.Secret))
				}
				msgs := otr3.Bytes(ret)
				if err != nil {
					alert(term, "Error while starting authentication with "+to+": "+err.Error())
				}
				for _, msg := range msgs {
					s.Conn.Send(to, string(msg))
				}
			case authOobCommand:
				fpr, err := hex.DecodeString(cmd.Fingerprint)
				if err != nil {
					alert(term, fmt.Sprintf("Invalid fingerprint %s - not authenticated", cmd.Fingerprint))
					break
				}
				existing := s.Config.UserIdForFingerprint(fpr)
				if len(existing) != 0 {
					alert(term, fmt.Sprintf("Fingerprint %s already belongs to %s", cmd.Fingerprint, existing))
					break
				}
				s.Config.KnownFingerprints = append(s.Config.KnownFingerprints, config.KnownFingerprint{Fingerprint: fpr, UserId: cmd.User})
				s.Config.Save()
				info(term, fmt.Sprintf("Saved manually verified fingerprint %s for %s", cmd.Fingerprint, cmd.User))
			case awayCommand:
				s.Conn.SignalPresence("away")
			case chatCommand:
				s.Conn.SignalPresence("chat")
			case dndCommand:
				s.Conn.SignalPresence("dnd")
			case xaCommand:
				s.Conn.SignalPresence("xa")
			case onlineCommand:
				s.Conn.SignalPresence("")
			}
		}
	}
}

func enroll(conf *config.Config, term *terminal.Terminal) bool {
	var err error
	warn(term, "Enrolling new config file")

	var domain string
	for {
		term.SetPrompt("Account (i.e. user@example.com, enter to quit): ")
		if conf.Account, err = term.ReadLine(); err != nil || len(conf.Account) == 0 {
			return false
		}

		parts := strings.SplitN(conf.Account, "@", 2)
		if len(parts) != 2 {
			alert(term, "invalid username (want user@domain): "+conf.Account)
			continue
		}
		domain = parts[1]
		break
	}

	term.SetPrompt("Enable debug logging to /tmp/xmpp-client-debug.log? ")
	if debugLog, err := term.ReadLine(); err != nil || !config.ParseYes(debugLog) {
		info(term, "Not enabling debug logging...")
	} else {
		info(term, "Debug logging enabled...")
		conf.RawLogFile = "/tmp/xmpp-client-debug.log"
	}

	term.SetPrompt("Use Tor?: ")
	if useTorQuery, err := term.ReadLine(); err != nil || len(useTorQuery) == 0 || !config.ParseYes(useTorQuery) {
		info(term, "Not using Tor...")
		conf.UseTor = false
	} else {
		info(term, "Using Tor...")
		conf.UseTor = true
	}

	term.SetPrompt("File to import libotr private key from (enter to generate): ")

	var priv otr3.PrivateKey
	for {
		importFile, err := term.ReadLine()
		if err != nil {
			return false
		}
		if len(importFile) > 0 {
			privKeyBytes, err := ioutil.ReadFile(importFile)
			if err != nil {
				alert(term, "Failed to open private key file: "+err.Error())
				continue
			}

			if !priv.Import(privKeyBytes) {
				alert(term, "Failed to parse libotr private key file (the parser is pretty simple I'm afraid)")
				continue
			}
			break
		} else {
			info(term, "Generating private key...")
			priv.Generate(rand.Reader)
			break
		}
	}
	conf.PrivateKey = priv.Serialize()

	conf.OTRAutoAppendTag = true
	conf.OTRAutoStartSession = true
	conf.OTRAutoTearDown = false

	// List well known Tor hidden services.
	knownTorDomain := map[string]string{
		"jabber.ccc.de":             "okj7xc6j2szr2y75.onion",
		"riseup.net":                "4cjw6cwpeaeppfqz.onion",
		"jabber.calyxinstitute.org": "ijeeynrc6x2uy5ob.onion",
		"jabber.otr.im":             "5rgdtlawqkcplz75.onion",
		"wtfismyip.com":             "ofkztxcohimx34la.onion",
	}

	// Autoconfigure well known Tor hidden services.
	if hiddenService, ok := knownTorDomain[domain]; ok && conf.UseTor {
		const torProxyURL = "socks5://127.0.0.1:9050"
		info(term, "It appears that you are using a well known server and we will use its Tor hidden service to connect.")
		conf.Server = hiddenService
		conf.Port = 5222
		conf.Proxies = []string{torProxyURL}
		term.SetPrompt("> ")
		return true
	}

	var proxyStr string
	proxyDefaultPrompt := ", enter for none"
	if conf.UseTor {
		proxyDefaultPrompt = ", which is the default"
	}
	term.SetPrompt("Proxy (i.e socks5://127.0.0.1:9050" + proxyDefaultPrompt + "): ")

	for {
		if proxyStr, err = term.ReadLine(); err != nil {
			return false
		}
		if len(proxyStr) == 0 {
			if !conf.UseTor {
				break
			} else {
				proxyStr = "socks5://127.0.0.1:9050"
			}
		}
		u, err := url.Parse(proxyStr)
		if err != nil {
			alert(term, "Failed to parse "+proxyStr+" as a URL: "+err.Error())
			continue
		}
		if _, err = proxy.FromURL(u, proxy.Direct); err != nil {
			alert(term, "Failed to parse "+proxyStr+" as a proxy: "+err.Error())
			continue
		}
		break
	}

	if len(proxyStr) > 0 {
		conf.Proxies = []string{proxyStr}

		info(term, "Since you selected a proxy, we need to know the server and port to connect to as a SRV lookup would leak information every time.")
		term.SetPrompt("Server (i.e. xmpp.example.com, enter to lookup using unproxied DNS): ")
		if conf.Server, err = term.ReadLine(); err != nil {
			return false
		}
		if len(conf.Server) == 0 {
			var port uint16
			info(term, "Performing SRV lookup")
			if conf.Server, port, err = xmpp.Resolve(domain); err != nil {
				alert(term, "SRV lookup failed: "+err.Error())
				return false
			}
			conf.Port = int(port)
			info(term, "Resolved "+conf.Server+":"+strconv.Itoa(conf.Port))
		} else {
			for {
				term.SetPrompt("Port (enter for 5222): ")
				portStr, err := term.ReadLine()
				if err != nil {
					return false
				}
				if len(portStr) == 0 {
					portStr = "5222"
				}
				if conf.Port, err = strconv.Atoi(portStr); err != nil || conf.Port <= 0 || conf.Port > 65535 {
					info(term, "Port numbers must be 0 < port <= 65535")
					continue
				}
				break
			}
		}
	}

	term.SetPrompt("> ")

	return true
}

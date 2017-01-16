package cli

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/twstrike/coyim/cli/terminal"
	"github.com/twstrike/coyim/client"
	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/roster"
	"github.com/twstrike/coyim/servers"
	sessions "github.com/twstrike/coyim/session/access"
	"github.com/twstrike/coyim/tls"
	"github.com/twstrike/coyim/xmpp/data"
	xi "github.com/twstrike/coyim/xmpp/interfaces"

	"github.com/twstrike/otr3"

	"golang.org/x/net/proxy"
)

type cliUI struct {
	session        sessions.Session
	sessionFactory sessions.Factory
	dialerFactory  func(tls.Verifier) xi.Dialer
	events         chan interface{}
	commands       chan interface{}

	password string

	termControl terminal.Control
	oldState    interface{}
	term        terminal.Terminal

	input *input

	terminate chan bool

	RosterEditor
}

// UI is the user interface functionality exposed to main
type UI interface {
	Loop()
}

// NewCLI creates a new cliUI instance
func NewCLI(version string, cf terminal.ControlFactory, sf sessions.Factory, df func(tls.Verifier) xi.Dialer) UI {
	termControl := cf()
	oldState, err := termControl.MakeRaw(0)
	if err != nil {
		panic(err.Error())
	}

	term := termControl.NewTerminal(os.Stdin, "")
	updateTerminalSize(term, termControl)
	term.SetBracketedPasteMode(true)

	resizeChan := make(chan os.Signal)
	go func() {
		for _ = range resizeChan {
			updateTerminalSize(term, termControl)
		}
	}()
	signal.Notify(resizeChan, syscall.SIGWINCH)

	return &cliUI{
		termControl: termControl,
		term:        term,
		oldState:    oldState,
		terminate:   make(chan bool),
		input: &input{
			tc:          termControl,
			term:        term,
			uidComplete: new(priorityList),
		},
		RosterEditor: RosterEditor{
			PendingRosterChan: make(chan *RosterEdit),
		},
		events:         make(chan interface{}),
		commands:       make(chan interface{}, 5),
		sessionFactory: sf,
		dialerFactory:  df,
	}
}

func (c *cliUI) getMasterPassword(params config.EncryptionParameters, previousFailed bool) ([]byte, []byte, bool) {
	password, err := c.term.ReadPassword("Master password for configuration file: ")
	if err != nil {
		c.alert(err.Error())
		return nil, nil, false
	}

	l, r := config.GenerateKeys(password, params)
	return l, r, true
}

func findAccount(a *string, acs []*config.Account) *config.Account {
	if a != nil && *a != "" {
		for _, ac := range acs {
			if ac.Account == *a {
				return ac
			}
		}
	}
	return acs[0]
}

func (c *cliUI) loadConfig(configFile string) error {
	accounts, ok, err := config.LoadOrCreate(configFile, config.FunctionKeySupplier(c.getMasterPassword))
	if !ok {
		c.alert("Couldn't open encrypted file - did you enter your password correctly?")
		return errors.New("couldn't open encrypted file - did you supply the wrong password?")
	}
	if err != nil {
		c.alert(err.Error())
		acc, e := accounts.AddNewAccount()
		if e != nil {
			return e
		}
		if !c.enroll(accounts, acc) {
			return errors.New("asked to quit")
		}
	}

	account := findAccount(config.AccountFlag, accounts.Accounts)

	var password string
	if len(account.Password) == 0 {
		var err error

		password, err = c.term.ReadPassword(
			fmt.Sprintf("Password for %s (will not be saved to disk): ", account.Account),
		)
		if err != nil {
			c.alert(err.Error())
			return err
		}
	} else {
		password = account.Password
	}

	logger := &lineLogger{c.term, c.termControl, nil}
	if *config.DebugFlag {
		log.SetOutput(logger)
	}

	//TODO: call session.ConnectAndRegister() in this case
	//var registerCallback xmpp.FormCallback
	//if *config.CreateAccount {
	//	registerCallback = c.RegisterCallback
	//}

	c.session = c.sessionFactory(accounts, account, c.dialerFactory)
	c.session.SetSessionEventHandler(c)
	c.session.Subscribe(c.events)
	c.session.SetCommandManager(c)
	c.session.SetConnectionLogger(logger)
	// TODO: this nil is incorrect and will cause failures when trying to use it.
	if err := c.session.Connect(password, c.verifier()); err != nil {
		c.alert(err.Error())
		return err
	}

	return nil
}

func (c *cliUI) quit() {
	c.session.Close()
	c.terminate <- true
}

func (c *cliUI) SaveConf() {
	c.session.Config().Save(config.FunctionKeySupplier(c.getMasterPassword))
}

func (c *cliUI) Loop() {
	defer c.close()

	if err := c.loadConfig(*config.ConfigFile); err != nil {
		return
	}

	go c.watchClientCommands()
	go c.observeSessionEvents()
	go c.watchRosterEdits()
	go c.watchInputCommands()

	<-c.terminate // wait
}

func (c *cliUI) close() {
	if c.oldState != nil {
		c.termControl.Restore(0, c.oldState)
	}

	if c.term != nil {
		c.term.SetBracketedPasteMode(false)
	}
}

func (c *cliUI) info(m string) {
	info(c.term, c.termControl, m)
}

func (c *cliUI) warn(m string) {
	warn(c.term, c.termControl, m)
}

func (c *cliUI) alert(m string) {
	alert(c.term, c.termControl, m)
}

func (c *cliUI) critical(m string) {
	critical(c.term, c.termControl, m)
}

func (c *cliUI) RegisterCallback(title, instructions string, fields []interface{}) error {
	user := c.session.GetConfig().Account
	return c.promptForForm(user, c.password, title, instructions, fields)
}

func (c *cliUI) printConversationInfo(uid string, conversation client.Conversation) {
	s := c.session

	fpr := conversation.TheirFingerprint()
	fprUID := s.GetConfig().UserIDForVerifiedFingerprint(fpr)
	c.info(fmt.Sprintf("  Fingerprint  for %s: %X", uid, fpr))
	c.info(fmt.Sprintf("  Session  ID  for %s: %X", uid, conversation.GetSSID()))
	if fprUID == uid {
		c.info(fmt.Sprintf("  Identity key for %s is verified", uid))
	} else if len(fprUID) > 1 {
		c.alert(fmt.Sprintf("  Warning: %s is using an identity key which was verified for %s", uid, fprUID))
	} else if s.GetConfig().HasFingerprint(uid) {
		c.critical(fmt.Sprintf("  Identity key for %s is incorrect", uid))
	} else {
		c.alert(fmt.Sprintf("  Identity key for %s is not verified. You should use /otr-auth or /otr-authqa or /otr-authoob to verify their identity", uid))
	}
}

// promptForForm runs an XEP-0004 form and collects responses from the user.
func (c *cliUI) promptForForm(user, password, title, instructions string, fields []interface{}) error {
	c.info("The server has requested the following information. Text that has come from the server will be shown in red.")

	// formStringForPrinting takes a string form the form and returns an
	// escaped version with codes to make it show as red.
	formStringForPrinting := func(s string) string {
		var line []byte

		line = append(line, c.termControl.Escape(c.term).Red...)
		line = appendTerminalEscaped(line, []byte(s))
		line = append(line, c.termControl.Escape(c.term).Reset...)
		return string(line)
	}

	write := func(s string) {
		c.term.Write([]byte(s))
	}

	var tmpDir string

	showMediaEntries := func(questionNumber int, medias [][]data.Media) {
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
		case *data.FixedFormField:
			write(formStringForPrinting(field.Text))
			write("\n")
			questionNumber--

		case *data.BooleanFormField:
			write(fmt.Sprintf("%d. %s\n\n", questionNumber, formStringForPrinting(field.Label)))
			showMediaEntries(questionNumber, field.Media)
			c.term.SetPrompt("Please enter yes, y, no or n: ")

		TryAgain:
			for {
				answer, err := c.term.ReadLine()
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

		case *data.TextFormField:
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
			c.term.SetPrompt("> ")
			if field.Private {
				field.Result, err = c.term.ReadPassword("> ")
			} else {
				field.Result, err = c.term.ReadLine()
			}
			if err != nil {
				return err
			}
			if len(field.Result) == 0 {
				field.Result = field.Default
			}

		case *data.MultiTextFormField:
			write(fmt.Sprintf("%d. %s\n\n", questionNumber, formStringForPrinting(field.Label)))
			showMediaEntries(questionNumber, field.Media)

			write("Please enter one or more responses, terminated by an empty line\n")
			c.term.SetPrompt("> ")

			for {
				line, err := c.term.ReadLine()
				if err != nil {
					return err
				}
				if len(line) == 0 {
					break
				}
				field.Results = append(field.Results, line)
			}

		case *data.SelectionFormField:
			write(fmt.Sprintf("%d. %s\n\n", questionNumber, formStringForPrinting(field.Label)))
			showMediaEntries(questionNumber, field.Media)

			for i, opt := range field.Values {
				write(fmt.Sprintf("  %d. %s\n\n", i+1, formStringForPrinting(opt)))
			}
			c.term.SetPrompt("Please enter the number of your selection: ")

		TryAgain2:
			for {
				answer, err := c.term.ReadLine()
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

		case *data.MultiSelectionFormField:
			write(fmt.Sprintf("%d. %s\n\n", questionNumber, formStringForPrinting(field.Label)))
			showMediaEntries(questionNumber, field.Media)

			for i, opt := range field.Values {
				write(fmt.Sprintf("  %d. %s\n\n", i+1, formStringForPrinting(opt)))
			}
			c.term.SetPrompt("Please enter the numbers of zero or more of the above, separated by spaces: ")

		TryAgain3:
			for {
				answer, err := c.term.ReadLine()
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

func (c *cliUI) watchRosterEdits() {
	for edit := range c.PendingRosterChan {
		if !edit.IsComplete {
			c.info("Please edit " + edit.FileName + " and run /rostereditdone when complete")
			c.PendingRosterEdit = edit
			continue
		}

		parsedRoster, err := parseEditedRoster(edit.Contents)
		if err != nil {
			c.alert(err.Error())
			c.alert("Please reedit file and run /rostereditdone again")
			continue
		}

		toDelete, toEdit, toAdd := diffRoster(parsedRoster, edit.Roster)

		//DELETE
		for _, jid := range toDelete {
			c.info("Deleting roster entry for " + jid)
			_, _, err := c.session.Conn().SendIQ("" /* to the server */, "set", data.RosterRequest{
				Item: data.RosterRequestItem{
					Jid:          jid,
					Subscription: "remove",
				},
			})

			if err != nil {
				c.alert("Failed to remove roster entry: " + err.Error())
			}

			c.session.GetConfig().RemovePeer(jid)
			c.session.CommandManager().ExecuteCmd(client.SaveApplicationConfigCmd{})
		}

		//EDIT
		for _, entry := range toEdit {
			c.info("Updating roster entry for " + entry.Jid)
			_, _, err := c.session.Conn().SendIQ("" /* to the server */, "set", data.RosterRequest{
				Item: data.RosterRequestItem{
					Jid:   entry.Jid,
					Name:  entry.Name,
					Group: entry.Group,
				},
			})

			if err != nil {
				c.alert("Failed to update roster entry: " + err.Error())
			}
		}

		//ADD
		for _, entry := range toAdd {
			c.info("Adding roster entry for " + entry.Jid)
			_, _, err := c.session.Conn().SendIQ("" /* to the server */, "set", data.RosterRequest{
				Item: data.RosterRequestItem{
					Jid:   entry.Jid,
					Name:  entry.Name,
					Group: entry.Group,
				},
			})

			if err != nil {
				c.alert("Failed to add roster entry: " + err.Error())
			}
		}

		c.PendingRosterEdit = nil
	}
}

func (c *cliUI) withCurrentPeer(peer string, f func(*roster.Peer)) {
	p, ok := c.session.R().Get(peer)

	if ok {
		f(p)
	}
}

func (c *cliUI) currentResourceFor(peer string) string {
	resource := ""
	c.withCurrentPeer(peer, func(p *roster.Peer) {
		resource = p.ResourceToUse()
	})
	return resource
}

func (c *cliUI) watchInputCommands() {
	defer c.quit()

	commandChan := make(chan interface{})
	go c.input.processCommands(commandChan)
	c.term.SetPrompt("> ")

	s := c.session
	conf := s.GetConfig()

	var err error

CommandLoop:
	for {
		select {
		case cmd, ok := <-commandChan:
			if !ok {
				c.warn("Exiting because command channel closed")
				break CommandLoop
			}
			s.SetLastActionTime(time.Now())
			switch cmd := cmd.(type) {
			case quitCommand:
				break CommandLoop
			case versionCommand:
				replyChan, cookie, err := s.Conn().SendIQ(cmd.User, "get", data.VersionQuery{})
				if err != nil {

					c.alert("Error sending version request: " + err.Error())
					continue
				}

				s.Timeout(cookie, time.Now().Add(5*time.Second))
				go s.AwaitVersionReply(replyChan, cmd.User)
			case rosterCommand:
				c.info("Current roster:")
				maxLen := 0
				for _, item := range s.R().ToSlice() {
					if maxLen < len(item.Jid) {
						maxLen = len(item.Jid)
					}
				}

				for _, item := range s.R().ToSlice() {
					state, _, ok := s.R().StateOf(item.Jid)

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
					c.info(line)
				}
			case rosterEditCommand:
				if c.PendingRosterEdit != nil {
					c.warn("Aborting previous roster edit")
					c.PendingRosterEdit = nil
				}

				currR := s.R().ToSlice()

				c.RosterEditor.Roster = make([]data.RosterEntry, len(currR))
				rosterCopy := make([]data.RosterEntry, len(currR))
				for ix, e := range currR {
					c.RosterEditor.Roster[ix] = e.ToEntry()
					rosterCopy[ix] = e.ToEntry()
				}

				go func(rosterCopy []data.RosterEntry) {
					err := c.EditRoster(rosterCopy)
					if err != nil {
						c.alert(err.Error())
					}
				}(rosterCopy)

			case rosterEditDoneCommand:
				if c.PendingRosterEdit == nil {
					c.warn("No roster edit in progress. Use /rosteredit to start one")
					continue
				}

				go func(edit RosterEdit) {
					err := c.LoadEditedRoster(edit)
					if err != nil {
						c.alert(err.Error())
					}
				}(*c.PendingRosterEdit)

			case toggleStatusUpdatesCommand:
				s.GetConfig().HideStatusUpdates = !s.GetConfig().HideStatusUpdates
				s.CommandManager().ExecuteCmd(client.SaveApplicationConfigCmd{})

				// Tell the user the current state of the statuses
				if s.GetConfig().HideStatusUpdates {
					c.info("Status updates disabled")
				} else {
					c.info("Status updates enabled")
				}
			case confirmCommand:
				s.HandleConfirmOrDeny(cmd.User, true /* confirm */)
			case denyCommand:
				s.HandleConfirmOrDeny(cmd.User, false /* deny */)
			case addCommand:
				s.RequestPresenceSubscription(cmd.User, "") // second argument is the potential message
			case msgCommand:
				message := []byte(cmd.msg)
				conversation, _ := s.ConversationManager().EnsureConversationWith(cmd.to, c.currentResourceFor(cmd.to))
				isEncrypted := conversation.IsEncrypted()
				if cmd.setPromptIsEncrypted != nil {
					cmd.setPromptIsEncrypted <- isEncrypted
				}

				if !isEncrypted && conf.ShouldEncryptTo(cmd.to) {
					c.warn(fmt.Sprintf("Did not send: no encryption established with %s", cmd.to))
					continue
				}

				_, err := conversation.Send(s, c.currentResourceFor(cmd.to), message)
				if err != nil {
					c.alert(err.Error())
					break
				}

			case otrCommand:
				resource := c.currentResourceFor(string(cmd.User))
				conversation, _ := s.ConversationManager().EnsureConversationWith(string(cmd.User), resource)
				conversation.StartEncryptedChat(s, resource)
			case otrInfoCommand:
				for _, pk := range s.PrivateKeys() {
					c.info(fmt.Sprintf("Your OTR fingerprint is %x", pk.PublicKey().Fingerprint()))
				}
				for to, conversation := range s.ConversationManager().Conversations() {
					if conversation.IsEncrypted() {
						c.info(fmt.Sprintf("Secure session with %s underway:", to))
						c.printConversationInfo(to, conversation)
					}
				}
			case endOTRCommand:
				to := string(cmd.User)
				resource := c.currentResourceFor(to)
				conversation, exists := s.ConversationManager().GetConversationWith(to, resource)

				if !exists {
					c.alert("No secure session established")
					break
				}

				err := conversation.EndEncryptedChat(s, resource)
				if err != nil {
					c.alert("Can't end the conversation - it seems there is no randomness in your system. This could be a significant problem.")
					break
				}

				c.input.SetPromptForTarget(cmd.User, false)
				c.warn("OTR conversation ended with " + cmd.User)
			case authQACommand:
				to := string(cmd.User)
				resource := c.currentResourceFor(to)
				conversation, exists := s.ConversationManager().GetConversationWith(to, resource)
				if !exists {
					c.alert("Can't authenticate without a secure conversation established")
					break
				}

				if s.OtrEventHandler()[to].WaitingForSecret {
					s.OtrEventHandler()[to].WaitingForSecret = false
					err = conversation.ProvideAuthenticationSecret(s, resource, []byte(cmd.Secret))
				} else {
					err = conversation.StartAuthenticate(s, resource, cmd.Question, []byte(cmd.Secret))
				}

				if err != nil {
					c.alert("Error while starting authentication with " + to + ": " + err.Error())
				}

			case authOobCommand:
				fpr, err := hex.DecodeString(cmd.Fingerprint)
				if err != nil {
					c.alert(fmt.Sprintf("Invalid fingerprint %s - not authenticated", cmd.Fingerprint))
					break
				}

				idForFpr := s.GetConfig().UserIDForVerifiedFingerprint(fpr)
				if len(idForFpr) != 0 {
					c.alert(fmt.Sprintf("Fingerprint %s already belongs to %s", cmd.Fingerprint, idForFpr))
					break
				}

				s.CommandManager().ExecuteCmd(client.AuthorizeFingerprintCmd{
					Account:     s.GetConfig(),
					Session:     s,
					Peer:        cmd.User,
					Fingerprint: fpr,
				})

				c.info(fmt.Sprintf("Saved manually verified fingerprint %s for %s", cmd.Fingerprint, cmd.User))
			case awayCommand:
				s.Conn().SignalPresence("away")
			case chatCommand:
				s.Conn().SignalPresence("chat")
			case dndCommand:
				s.Conn().SignalPresence("dnd")
			case xaCommand:
				s.Conn().SignalPresence("xa")
			case onlineCommand:
				s.Conn().SignalPresence("")
			}
		}
	}
}

func (c *cliUI) enroll(conf *config.ApplicationConfig, currentConf *config.Account) bool {
	var err error
	c.warn("Enrolling new config file")

	var domain string
	for {
		c.term.SetPrompt("Account (i.e. user@example.com, enter to quit): ")
		if currentConf.Account, err = c.term.ReadLine(); err != nil || len(currentConf.Account) == 0 {
			return false
		}

		parts := strings.SplitN(currentConf.Account, "@", 2)
		if len(parts) != 2 {
			c.alert("invalid username (want user@domain): " + currentConf.Account)
			continue
		}
		domain = parts[1]
		break
	}

	c.term.SetPrompt("Enable debug logging to /tmp/xmpp-client-debug.log? ")
	if debugLog, err := c.term.ReadLine(); err != nil || !config.ParseYes(debugLog) {
		c.info("Not enabling debug logging...")
	} else {
		c.info("Debug logging enabled...")
		conf.RawLogFile = "/tmp/xmpp-client-debug.log"
	}

	c.term.SetPrompt("Use Tor?: ")
	if useTorQuery, err := c.term.ReadLine(); err != nil || len(useTorQuery) == 0 || !config.ParseYes(useTorQuery) {
		c.info("Not using Tor...")
		currentConf.Proxies = []string{}
	} else {
		c.info("Using Tor...")
	}

	c.term.SetPrompt("File to import libotr private key from (enter to generate): ")

	var pkeys []otr3.PrivateKey
	for {
		importFile, err := c.term.ReadLine()
		if err != nil {
			return false
		}
		if len(importFile) > 0 {
			privKeyBytes, err := ioutil.ReadFile(importFile)
			if err != nil {
				c.alert("Failed to open private key file: " + err.Error())
				continue
			}
			var priv otr3.DSAPrivateKey
			if !priv.Import(privKeyBytes) {
				c.alert("Failed to parse libotr private key file (the parser is pretty simple I'm afraid)")
				continue
			}
			pkeys = append(pkeys, &priv)
			break
		} else {
			c.info("Generating private key...")
			pkeys, err = otr3.GenerateMissingKeys([][]byte{})
			if err != nil {
				c.alert("Failed to generate private key - this implies something is really bad with your system, so we bail out now")
				return false
			}
			break
		}
	}

	currentConf.PrivateKeys = config.SerializedKeys(pkeys)
	currentConf.OTRAutoAppendTag = true
	currentConf.OTRAutoStartSession = true
	currentConf.OTRAutoTearDown = false

	// Force Tor for servers with well known Tor hidden services.
	if _, ok := servers.Get(domain); ok && currentConf.HasTorAuto() {
		const torProxyURL = "socks5://127.0.0.1:9050"
		c.info("It appears that you are using a well known server and we will use its Tor hidden service to connect.")
		currentConf.Proxies = []string{torProxyURL}
		c.term.SetPrompt("> ")
		return true
	}

	var proxyStr string
	proxyDefaultPrompt := ", enter for none"
	if currentConf.HasTorAuto() {
		proxyDefaultPrompt = ", which is the default"
	}
	c.term.SetPrompt("Proxy (i.e socks5://127.0.0.1:9050" + proxyDefaultPrompt + "): ")

	for {
		if proxyStr, err = c.term.ReadLine(); err != nil {
			return false
		}
		if len(proxyStr) == 0 {
			if !currentConf.HasTorAuto() {
				break
			} else {
				proxyStr = "socks5://127.0.0.1:9050"
			}
		}
		u, err := url.Parse(proxyStr)
		if err != nil {
			c.alert("Failed to parse " + proxyStr + " as a URL: " + err.Error())
			continue
		}
		if _, err = proxy.FromURL(u, proxy.Direct); err != nil {
			c.alert("Failed to parse " + proxyStr + " as a proxy: " + err.Error())
			continue
		}
		break
	}

	c.term.SetPrompt("> ")

	return true
}

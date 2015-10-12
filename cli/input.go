package cli

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh/terminal"
)

type uiCommand struct {
	name      string
	prototype interface{}
	desc      string
}

var uiCommands = []uiCommand{
	{"add", addCommand{}, "Request a subscription to another user's presence"},
	{"away", awayCommand{}, "Set your status to Away"},
	{"chat", chatCommand{}, "Set your status to Available for Chat"},
	{"close", closeCommand{}, "Forget current chat target"},
	{"confirm", confirmCommand{}, "Confirm an inbound subscription request"},
	{"deny", denyCommand{}, "Deny an inbound subscription request"},
	{"dnd", dndCommand{}, "Set your status to Busy / Do Not Disturb"},
	{"help", helpCommand{}, "List known commands"},
	{"nopaste", noPasteCommand{}, "Stop interpreting text verbatim"},
	{"online", onlineCommand{}, "Set your status to Available / Online"},
	{"otr-auth", authCommand{}, "Authenticate a secure peer with a mutual, shared secret"},
	{"otr-authoob", authOobCommand{}, "Authenticate a secure peer with out-of-band fingerprint verification"},
	{"otr-authqa", authQACommand{}, "Authenticate a secure peer with a question and answer"},
	{"otr-end", endOTRCommand{}, "End an OTR session"},
	{"otr-info", otrInfoCommand{}, "Print OTR information such as OTR fingerprint"},
	{"otr-start", otrCommand{}, "Start an OTR session with the given user"},
	{"paste", pasteCommand{}, "Start interpreting text verbatim"},
	{"quit", quitCommand{}, "Quit the program"},
	{"rostereditdone", rosterEditDoneCommand{}, "Load the edited roster from disk"},
	{"rosteredit", rosterEditCommand{}, "Write the roster to disk"},
	{"roster", rosterCommand{}, "Display the current roster"},
	{"statusupdates", toggleStatusUpdatesCommand{}, "Toggle if status updates are displayed"},
	{"version", versionCommand{}, "Ask a Jabber client for its version"},
	{"xa", xaCommand{}, "Set your status to Extended Away"},
}

type addCommand struct {
	User string "uid"
}

type authCommand struct {
	User   string "uid"
	Secret string
}

type authOobCommand struct {
	User        string "uid"
	Fingerprint string
}

type authQACommand struct {
	User     string "uid"
	Question string
	Secret   string
}

type awayCommand struct{}
type chatCommand struct{}
type closeCommand struct{}

type confirmCommand struct {
	User string "uid"
}

type denyCommand struct {
	User string "uid"
}

type dndCommand struct{}

type endOTRCommand struct {
	User string "uid"
}

type helpCommand struct{}

type msgCommand struct {
	to  string
	msg string
	// setPromptIsEncrypted is used to synchonously indicate whether the
	// prompt should show the contact as encrypted, before the prompt is
	// redrawn. It may be nil to indicate that the prompt cannot be
	// updated but otherwise must be sent to.
	setPromptIsEncrypted chan<- bool
}

type noPasteCommand struct{}
type onlineCommand struct{}

type otrCommand struct {
	User string "uid"
}

type otrInfoCommand struct{}

type pasteCommand struct{}
type quitCommand struct{}

type rosterCommand struct {
	OnlineOnly bool "flag:online"
}

type rosterEditCommand struct{}
type rosterEditDoneCommand struct{}
type toggleStatusUpdatesCommand struct{}

type versionCommand struct {
	User string "uid"
}

type xaCommand struct{}

func numPositionalFields(t reflect.Type) int {
	for i := 0; i < t.NumField(); i++ {
		if strings.HasPrefix(string(t.Field(i).Tag), "flag:") {
			return i
		}
	}
	return t.NumField()
}

func parseCommandForCompletion(commands []uiCommand, line string) (before, prefix string, isCommand, ok bool) {
	if len(line) == 0 || line[0] != '/' {
		return
	}

	spacePos := strings.IndexRune(line, ' ')
	if spacePos == -1 {
		// We're completing a command name.
		before = line[:1]
		prefix = line[1:]
		isCommand = true
		ok = true
		return
	}

	command := line[1:spacePos]
	var prototype interface{}

	for _, cmd := range commands {
		if cmd.name == command {
			prototype = cmd.prototype
			break
		}
	}
	if prototype == nil {
		return
	}

	t := reflect.TypeOf(prototype)
	fieldNum := -1
	fieldStart := 0
	inQuotes := false
	lastWasEscape := false
	numFields := numPositionalFields(t)

	skippingWhitespace := true
	for pos, r := range line[spacePos:] {
		if skippingWhitespace {
			if r == ' ' {
				continue
			}
			skippingWhitespace = false
			fieldNum++
			fieldStart = pos + spacePos
		}

		if lastWasEscape {
			lastWasEscape = false
			continue
		}

		if r == '\\' {
			lastWasEscape = true
			continue
		}

		if r == '"' {
			inQuotes = !inQuotes
		}

		if r == ' ' && !inQuotes {
			skippingWhitespace = true
		}
	}

	if skippingWhitespace {
		return
	}
	if fieldNum >= numFields {
		return
	}
	f := t.Field(fieldNum)
	if f.Tag != "uid" {
		return
	}
	ok = true
	isCommand = false
	before = line[:fieldStart]
	prefix = line[fieldStart:]
	return
}

// setOption updates the uiCommand, v, of type t given an option string with
// the "--" prefix already removed. It returns true on success.
func setOption(v reflect.Value, t reflect.Type, option string) bool {
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		tag := string(fieldType.Tag)
		if strings.HasPrefix(tag, "flag:") && tag[5:] == option {
			field := v.Field(i)
			if field.Bool() {
				return false // already set
			} else {
				field.SetBool(true)
				return true
			}
		}
	}

	return false
}

func parseCommand(commands []uiCommand, line []byte) (interface{}, string) {
	if len(line) == 0 || line[0] != '/' {
		panic("not a command")
	}

	spacePos := bytes.IndexByte(line, ' ')
	if spacePos == -1 {
		spacePos = len(line)
	}
	command := string(line[1:spacePos])
	var prototype interface{}

	for _, cmd := range commands {
		if cmd.name == command {
			prototype = cmd.prototype
			break
		}
	}
	if prototype == nil {
		return nil, "Unknown command: " + command
	}

	t := reflect.TypeOf(prototype)
	v := reflect.New(t)
	v = reflect.Indirect(v)
	pos := spacePos
	fieldNum := -1
	inQuotes := false
	lastWasEscape := false
	numFields := numPositionalFields(t)
	var field []byte

	skippingWhitespace := true
	for ; pos <= len(line); pos++ {
		if !skippingWhitespace && (pos == len(line) || (line[pos] == ' ' && !inQuotes && !lastWasEscape)) {
			skippingWhitespace = true
			strField := string(field)

			switch {
			case fieldNum < numFields:
				f := v.Field(fieldNum)
				f.Set(reflect.ValueOf(strField))
			case strings.HasPrefix(strField, "--"):
				if !setOption(v, t, strField[2:]) {
					return nil, "No such option " + strField + " for command"
				}
			default:
				return nil, "Too many arguments for command " + command + ". Expected " + strconv.Itoa(v.NumField())
			}
			field = field[:0]
			continue
		}

		if pos == len(line) {
			break
		}

		if lastWasEscape {
			field = append(field, line[pos])
			lastWasEscape = false
			continue
		}

		if skippingWhitespace {
			if line[pos] == ' ' {
				continue
			}
			skippingWhitespace = false
			fieldNum++
		}

		if line[pos] == '\\' {
			lastWasEscape = true
			continue
		}

		if line[pos] == '"' {
			inQuotes = !inQuotes
			continue
		}

		field = append(field, line[pos])
	}

	if fieldNum < numFields-1 {
		return nil, "Too few arguments for command " + command + ". Expected " + strconv.Itoa(v.NumField()) + ", but found " + strconv.Itoa(fieldNum+1)
	}

	return v.Interface(), ""
}

type Input struct {
	term                 *terminal.Terminal
	commands             *priorityList
	lastKeyWasCompletion bool

	// lock protects uids, uidComplete and lastTarget.
	lock        sync.Mutex
	uids        []string
	uidComplete *priorityList
	lastTarget  string
}

func (i *Input) AddUser(uid string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	for _, existingUid := range i.uids {
		if existingUid == uid {
			return
		}
	}

	i.uidComplete.Insert(uid)
	i.uids = append(i.uids, uid)
}

func (i *Input) ProcessCommands(commandsChan chan<- interface{}) {
	i.commands = new(priorityList)
	for _, command := range uiCommands {
		i.commands.Insert(command.name)
	}

	autoCompleteCallback := func(line string, pos int, key rune) (string, int, bool) {
		return i.AutoComplete(line, pos, key)
	}

	paste := false
	setPromptIsEncrypted := make(chan bool)

	for {
		if paste {
			i.term.AutoCompleteCallback = nil
		} else {
			i.term.AutoCompleteCallback = autoCompleteCallback
		}

		line, err := i.term.ReadLine()
		if err == terminal.ErrPasteIndicator {
			if len(i.lastTarget) == 0 {
				alert(i.term, "Pasted line ignored. Send a message to someone to select the destination.")
			} else {
				commandsChan <- msgCommand{i.lastTarget, string(line), nil}
			}
			continue
		}
		if err != nil {
			close(commandsChan)
			return
		}
		if paste {
			l := string(line)
			if l == "/nopaste" {
				paste = false
			} else {
				commandsChan <- msgCommand{i.lastTarget, l, nil}
			}
			continue
		}
		if len(line) == 0 {
			continue
		}
		if line[0] == '/' {
			cmd, err := parseCommand(uiCommands, []byte(line))
			if len(err) != 0 {
				alert(i.term, err)
				continue
			}
			// authCommand is turned into authQACommand with an
			// empty question.
			if authCmd, ok := cmd.(authCommand); ok {
				cmd = authQACommand{
					User:   authCmd.User,
					Secret: authCmd.Secret,
				}
			}
			if _, ok := cmd.(helpCommand); ok {
				i.showHelp()
				continue
			}
			if _, ok := cmd.(pasteCommand); ok {
				if len(i.lastTarget) == 0 {
					alert(i.term, "Can't enter paste mode without a destination. Send a message to someone to select the destination.")
					continue
				}
				paste = true
				continue
			}
			if _, ok := cmd.(noPasteCommand); ok {
				paste = false
				continue
			}
			if _, ok := cmd.(closeCommand); ok {
				i.lastTarget = ""
				i.term.SetPrompt("> ")
				continue
			}
			if cmd != nil {
				commandsChan <- cmd
			}
			continue
		}

		i.lock.Lock()
		if pos := strings.Index(line, string(nameTerminator)); pos > 0 {
			possibleName := line[:pos]
			for _, uid := range i.uids {
				if possibleName == uid {
					i.lastTarget = possibleName
					line = line[pos+2:]
					break
				}
			}
		}
		i.lock.Unlock()

		if len(i.lastTarget) == 0 {
			warn(i.term, "Start typing a Jabber address and hit tab to send a message to someone")
			continue
		}
		commandsChan <- msgCommand{i.lastTarget, string(line), setPromptIsEncrypted}
		isEncrypted := <-setPromptIsEncrypted
		i.SetPromptForTarget(i.lastTarget, isEncrypted)
	}
}

func (input *Input) SetPromptForTarget(target string, isEncrypted bool) {
	input.lock.Lock()
	isCurrent := input.lastTarget == target
	input.lock.Unlock()

	if !isCurrent {
		return
	}

	prompt := make([]byte, 0, len(target)+16)
	if isEncrypted {
		prompt = append(prompt, input.term.Escape.Green...)
	} else {
		prompt = append(prompt, input.term.Escape.Red...)
	}

	prompt = append(prompt, target...)
	prompt = append(prompt, input.term.Escape.Reset...)
	prompt = append(prompt, '>', ' ')
	input.term.SetPrompt(string(prompt))
}

func (input *Input) showHelp() {
	examples := make([]string, len(uiCommands))
	maxLen := 0

	for i, cmd := range uiCommands {
		line := "/" + cmd.name
		prototype := reflect.TypeOf(cmd.prototype)
		for j := 0; j < prototype.NumField(); j++ {
			if strings.HasPrefix(string(prototype.Field(j).Tag), "flag:") {
				line += " [--" + strings.ToLower(string(prototype.Field(j).Tag[5:])) + "]"
			} else {
				line += " <" + strings.ToLower(prototype.Field(j).Name) + ">"
			}
		}
		if l := len(line); l > maxLen {
			maxLen = l
		}
		examples[i] = line
	}

	for i, cmd := range uiCommands {
		line := examples[i]
		numSpaces := 1 + (maxLen - len(line))
		for j := 0; j < numSpaces; j++ {
			line += " "
		}
		line += cmd.desc
		info(input.term, line)
	}
}

const nameTerminator = ": "

func (i *Input) AutoComplete(line string, pos int, key rune) (string, int, bool) {
	const keyTab = 9

	if key != keyTab {
		i.lastKeyWasCompletion = false
		return "", -1, false
	}

	i.lock.Lock()
	defer i.lock.Unlock()

	prefix := line[:pos]
	if i.lastKeyWasCompletion {
		// The user hit tab right after a completion, so we got
		// it wrong.
		if len(prefix) > 0 && prefix[0] == '/' {
			if strings.IndexRune(prefix, ' ') == len(prefix)-1 {
				// We just completed a command.
				newCommand := i.commands.Next()
				newLine := "/" + string(newCommand) + " " + line[pos:]
				return newLine, len(newCommand) + 2, true
			} else if prefix[len(prefix)-1] == ' ' {
				// We just completed a uid in a command.
				newUser := i.uidComplete.Next()
				spacePos := strings.LastIndex(prefix[:len(prefix)-1], " ")

				newLine := prefix[:spacePos] + " " + string(newUser) + " " + line[pos:]
				return newLine, spacePos + 1 + len(newUser) + 1, true
			}
		} else if len(prefix) > 0 && prefix[0] != '/' && strings.HasSuffix(prefix, nameTerminator) {
			// We just completed a uid at the start of a
			// conversation line.
			newUser := i.uidComplete.Next()
			newLine := string(newUser) + nameTerminator + line[pos:]
			return newLine, len(newUser) + 2, true
		}
	} else {
		if len(prefix) > 0 && prefix[0] == '/' {
			a, b, isCommand, ok := parseCommandForCompletion(uiCommands, prefix)
			if !ok {
				return "", -1, false
			}
			var newValue string
			if isCommand {
				newValue, ok = i.commands.Find(b)
			} else {
				newValue, ok = i.uidComplete.Find(b)
			}
			if !ok {
				return "", -1, false
			}

			newLine := string(a) + newValue + " " + line[pos:]
			i.lastKeyWasCompletion = true
			return newLine, len(a) + len(newValue) + 1, true
		} else if len(prefix) > 0 && strings.IndexAny(prefix, ": \t") == -1 {
			// We're completing a uid at the start of a
			// conversation line.
			newUser, ok := i.uidComplete.Find(prefix)
			if !ok {
				return "", -1, false
			}

			newLine := newUser + nameTerminator + line[pos:]
			i.lastKeyWasCompletion = true
			return newLine, len(newUser) + len(nameTerminator), true
		}
	}

	i.lastKeyWasCompletion = false
	return "", -1, false
}

type priorityListEntry struct {
	value string
	next  *priorityListEntry
}

type priorityList struct {
	head       *priorityListEntry
	lastPrefix string
	lastResult string
	n          int
}

func (pl *priorityList) Insert(value string) {
	ent := new(priorityListEntry)
	ent.next = pl.head
	ent.value = value
	pl.head = ent
}

func (pl *priorityList) findNth(prefix string, nth int) (string, bool) {
	var cur, last *priorityListEntry
	cur = pl.head
	for n := 0; cur != nil; cur = cur.next {
		if strings.HasPrefix(cur.value, prefix) {
			if n == nth {
				// move this entry to the top
				if last != nil {
					last.next = cur.next
				} else {
					pl.head = cur.next
				}
				cur.next = pl.head
				pl.head = cur
				pl.lastResult = cur.value
				return cur.value, true
			}
			n++
		}
		last = cur
	}

	return "", false
}

func (pl *priorityList) Find(prefix string) (string, bool) {
	pl.lastPrefix = prefix
	pl.n = 0

	return pl.findNth(prefix, 0)
}

func (pl *priorityList) Next() string {
	pl.n++
	result, ok := pl.findNth(pl.lastPrefix, pl.n)
	if !ok {
		pl.n = 1
		result, ok = pl.findNth(pl.lastPrefix, pl.n)
	}
	// In this case, there's only one matching entry in the list.
	if !ok {
		pl.n = 0
		result, _ = pl.findNth(pl.lastPrefix, pl.n)
	}
	return result
}

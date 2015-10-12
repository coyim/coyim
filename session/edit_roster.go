package session

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/twstrike/coyim/config"
	"github.com/twstrike/coyim/ui"
	"github.com/twstrike/coyim/xmpp"
)

//TODO Move to CLI

// RosterEdit contains information about a pending roster edit. Roster edits
// occur by writing the roster to a file and inviting the user to edit the
// file.
//TODO Is this only relevant for CLI? How should this be handled by the GUI?
type RosterEdit struct {
	// FileName is the name of the file containing the roster information.
	FileName string
	// Roster contains the state of the roster at the time of writing the
	// file. It's what we diff against when reading the file.
	Roster []xmpp.RosterEntry
	// isComplete is true if this is the result of reading an edited
	// roster, rather than a report that the file has been written.
	IsComplete bool
	// contents contains the edited roster, if isComplete is true.
	Contents []byte
}

// editRoster runs in a goroutine and writes the roster to a file that the user
// can edit.
func (s *Session) EditRoster(roster []xmpp.RosterEntry) {
	// In case the editor rewrites the file, we work inside a temp
	// directory.
	dir, err := ioutil.TempDir("" /* system default temp dir */, "xmpp-client")
	if err != nil {
		s.alert("Failed to create temp dir to edit roster: " + err.Error())
		return
	}

	mode, err := os.Stat(dir)
	if err != nil || mode.Mode()&os.ModePerm != 0700 {
		panic("broken system libraries gave us an insecure temp dir")
	}

	fileName := filepath.Join(dir, "roster")
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		s.alert("Failed to create temp file: " + err.Error())
		return
	}

	io.WriteString(f, `# Use this file to edit your roster.
# The file is tab deliminated and you need to preserve that. Otherwise you
# can delete lines to remove roster entries, add lines to subscribe (only
# the account is needed when adding a line) and change lines to change the
# corresponding entry.

# Once you are done, use the /rostereditdone command to process the result.

# Since there are multiple levels of unspecified character encoding, we give up
# and hex escape anything outside of printable ASCII in "\x01" form.

`)

	// Calculate the number of tabs which covers the longest escaped JID.
	maxLen := 0
	escapedJids := make([]string, len(roster))
	for i, item := range roster {
		escapedJids[i] = ui.EscapeNonASCII(item.Jid)
		if l := len(escapedJids[i]); l > maxLen {
			maxLen = l
		}
	}
	tabs := (maxLen + 7) / 8

	for i, item := range s.Roster {
		line := escapedJids[i]
		tabsUsed := len(escapedJids[i]) / 8

		if len(item.Name) > 0 || len(item.Group) > 0 {
			// We're going to put something else on the line to tab
			// across to the next column.
			for i := 0; i < tabs-tabsUsed; i++ {
				line += "\t"
			}
		}

		if len(item.Name) > 0 {
			line += "name:" + ui.EscapeNonASCII(item.Name)
			if len(item.Group) > 0 {
				line += "\t"
			}
		}

		for j, group := range item.Group {
			if j > 0 {
				line += "\t"
			}
			line += "group:" + ui.EscapeNonASCII(group)
		}
		line += "\n"
		io.WriteString(f, line)
	}
	f.Close()

	s.PendingRosterChan <- &RosterEdit{
		FileName: fileName,
		Roster:   roster,
	}
}

func (s *Session) LoadEditedRoster(edit RosterEdit) {
	contents, err := ioutil.ReadFile(edit.FileName)
	if err != nil {
		s.alert("Failed to load edited roster: " + err.Error())
		return
	}
	os.Remove(edit.FileName)
	os.Remove(filepath.Dir(edit.FileName))

	edit.IsComplete = true
	edit.Contents = contents
	s.PendingRosterChan <- &edit
}

func (s *Session) processEditedRoster(edit *RosterEdit) bool {
	parsedRoster := make(map[string]xmpp.RosterEntry)
	lines := bytes.Split(edit.Contents, ui.NewLine)
	tab := []byte{'\t'}

	// Parse roster entries from the file.
	for i, line := range lines {
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts := bytes.Split(line, tab)

		var entry xmpp.RosterEntry
		var err error

		if entry.Jid, err = ui.UnescapeNonASCII(string(string(parts[0]))); err != nil {
			s.alert(fmt.Sprintf("Failed to parse JID on line %d: %s", i+1, err))
			return false
		}
		for _, part := range parts[1:] {
			if len(part) == 0 {
				continue
			}

			pos := bytes.IndexByte(part, ':')
			if pos == -1 {
				s.alert(fmt.Sprintf("Failed to find colon in item on line %d", i+1))
				return false
			}

			typ := string(part[:pos])
			value, err := ui.UnescapeNonASCII(string(part[pos+1:]))
			if err != nil {
				s.alert(fmt.Sprintf("Failed to unescape item on line %d: %s", i+1, err))
				return false
			}

			switch typ {
			case "name":
				if len(entry.Name) > 0 {
					s.alert(fmt.Sprintf("Multiple names given for contact on line %d", i+1))
					return false
				}
				entry.Name = value
			case "group":
				if len(value) > 0 {
					entry.Group = append(entry.Group, value)
				}
			default:
				s.alert(fmt.Sprintf("Unknown item tag '%s' on line %d", typ, i+1))
				return false
			}
		}

		parsedRoster[entry.Jid] = entry
	}

	// Now diff them from the original roster
	var toDelete []string
	var toEdit []xmpp.RosterEntry
	var toAdd []xmpp.RosterEntry

	for _, entry := range edit.Roster {
		newEntry, ok := parsedRoster[entry.Jid]
		if !ok {
			toDelete = append(toDelete, entry.Jid)
			continue
		}
		if newEntry.Name != entry.Name || !setEqual(newEntry.Group, entry.Group) {
			toEdit = append(toEdit, newEntry)
		}
	}

NextAdd:
	for jid, newEntry := range parsedRoster {
		for _, entry := range edit.Roster {
			if entry.Jid == jid {
				continue NextAdd
			}
		}
		toAdd = append(toAdd, newEntry)
	}

	for _, jid := range toDelete {
		s.info("Deleting roster entry for " + jid)
		_, _, err := s.Conn.SendIQ("" /* to the server */, "set", xmpp.RosterRequest{
			Item: xmpp.RosterRequestItem{
				Jid:          jid,
				Subscription: "remove",
			},
		})
		if err != nil {
			s.alert("Failed to remove roster entry: " + err.Error())
		}

		// Filter out any known fingerprints.
		newKnownFingerprints := make([]config.KnownFingerprint, 0, len(s.Config.KnownFingerprints))
		for _, fpr := range s.Config.KnownFingerprints {
			if fpr.UserId == jid {
				continue
			}
			newKnownFingerprints = append(newKnownFingerprints, fpr)
		}
		s.Config.KnownFingerprints = newKnownFingerprints
		s.Config.Save()
	}

	for _, entry := range toEdit {
		s.info("Updating roster entry for " + entry.Jid)
		_, _, err := s.Conn.SendIQ("" /* to the server */, "set", xmpp.RosterRequest{
			Item: xmpp.RosterRequestItem{
				Jid:   entry.Jid,
				Name:  entry.Name,
				Group: entry.Group,
			},
		})
		if err != nil {
			s.alert("Failed to update roster entry: " + err.Error())
		}
	}

	for _, entry := range toAdd {
		s.info("Adding roster entry for " + entry.Jid)
		_, _, err := s.Conn.SendIQ("" /* to the server */, "set", xmpp.RosterRequest{
			Item: xmpp.RosterRequestItem{
				Jid:   entry.Jid,
				Name:  entry.Name,
				Group: entry.Group,
			},
		})
		if err != nil {
			s.alert("Failed to add roster entry: " + err.Error())
		}
	}

	return true
}

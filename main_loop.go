package main


//func rosterLoop(term *terminal.Terminal, s *Session, rosterReply <-chan xmpp.Stanza, done chan<- bool) {
//RosterLoop:
//	for {
//		var err error
//
//		select {
//		case rosterStanza, ok := <-rosterReply:
//			if !ok {
//				//TODO was this error supposed to print the latest error regardless of where it came from?
//				alert(term, "Failed to read roster: "+err.Error())
//				break RosterLoop
//			}
//			if s.roster, err = xmpp.ParseRoster(rosterStanza); err != nil {
//				alert(term, "Failed to parse roster: "+err.Error())
//				break RosterLoop
//			}
//
//			s.rosterReceived()
//			info(term, "Roster received")
//
//		case edit := <-s.pendingRosterChan:
//			if !edit.IsComplete {
//				info(term, "Please edit "+edit.FileName+" and run /rostereditdone when complete")
//				s.pendingRosterEdit = edit
//				continue
//			}
//			if s.processEditedRoster(edit) {
//				s.pendingRosterEdit = nil
//			} else {
//				alert(term, "Please reedit file and run /rostereditdone again")
//			}
//
//		}
//	}
//
//	done <- true
//}

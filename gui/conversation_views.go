package gui

import (
	"sync"
)

// : fix the place where we can open to a specific resource, roster.go:296
// -- Make sure that the resources actually show up in the list!!!!!!!!!! I must have broke that functionality.
// : find all the ways and places and differences in how conversationViews can be created and gotten
//     for resource opening, this is used:
//					r.openConversationView(account, jid, true, resource)
//  r.openConversationView and r.ui.openConversationView seem to be the same
// - gtkUI.openConversationView (moved)
//      will try to find an existing cv, and if not, will create a new one, and then display it
//      used: all over the place
// - gtkUI.getConversationView (moved)
//      tries to get an existing cv from account, and then verifies that the window type is correct - potentially fixing it if not
//      used: in implementations of other stuff.
//            also in places that want to do something if the CV exists, but otherwise not
// - account.getConversationView (moved)
//      simply tries to get the cv from the inner list, and returns if it can find it. refactored to use cvs.getConversationView
//      used: same as above
// - cvs.getConversationView (moved)
//      simply gets the cv from the internal lists
//      used: same as above
// - account.getConversationWith (moved)
//      gets the cv from the account.getConversationView. it also fixing the type, just as in gtkUI.getConversationView
//      used: only in one place - for presence updates
// - gtkUI.createConversationViewBasedOnWindowLayout (moved)
//      creates a new wrapping type for an existing pane
//      used: only for implementation of other stuff
// - gtkUI.createConversationView (moved)
//      same as createConversationViewBasedOnWindowLayout with a nil argument
//      used: only in implementing other stuff
// - unified.createConversation (moved)
//      used: only in implementing other stuff
// - newConversationWindow (moved)
//      used: only in implementing other stuff

// FOUR: remove the JID and resource argument from almost all places
// FIVE: fix all the jid todos

// This file contains logic for dealing with conversation views and how they are locked
// and unlocked based on JID changes.
// For non-OTR conversations:
//    a conversation view can have a specific target resource or not
//    if there is a specific target resource, we will always use that.
//    if not, we will let roster.Peer define the behavior, based on the RFC around resources
// For OTR conversations:
//    the view will have a flag that is set when an OTR session is active.
//    when it's active, we will treat the window as locked to that resource.

// There is some weirdness here about what happens when you have an OTR conversation locked to someone
// and then want to open another window... No. You can't do that.

type conversationViews struct {
	// c contains all conversations. the ones indexed with a "resourced" JID will be locked to that view
	// everything else will be indexed with a bare jid
	c map[string]conversationView

	sync.RWMutex
}

func newConversationViews() *conversationViews {
	return &conversationViews{
		c: make(map[string]conversationView),
	}
}

func (cvs *conversationViews) enableExistingConversationWindows(enable bool) {
	cvs.RLock()
	defer cvs.RUnlock()

	for _, cv := range cvs.c {
		cv.setEnabled(enable)
	}
}

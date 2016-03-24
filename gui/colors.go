package gui

type colorSet struct {
	rosterPeerBackground               string
	rosterPeerOfflineForeground        string
	rosterPeerOnlineForeground         string
	rosterGroupBackground              string
	rosterAccountOnlineBackground      string
	rosterAccountOfflineBackground     string
	conversationOutgoingUserForeground string
	conversationIncomingUserForeground string
	conversationOutgoingTextForeground string
	conversationIncomingTextForeground string
	conversationStatusTextForeground   string
}

func (u *gtkUI) currentColorSet() colorSet {
	return u.defaultLightColorSet()
}

func (u *gtkUI) defaultLightColorSet() colorSet {
	return colorSet{
		rosterPeerBackground:               "#ffffff",
		rosterPeerOfflineForeground:        "#aaaaaa",
		rosterPeerOnlineForeground:         "#000000",
		rosterGroupBackground:              "#e9e7f3",
		rosterAccountOnlineBackground:      "#918caa",
		rosterAccountOfflineBackground:     "#d5d3de",
		conversationOutgoingUserForeground: "#3465a4",
		conversationIncomingUserForeground: "#a40000",
		conversationOutgoingTextForeground: "#555753",
		conversationIncomingTextForeground: "#000000",
		conversationStatusTextForeground:   "#4e9a06",
	}
}

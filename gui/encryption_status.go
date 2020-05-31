package gui

type encryptionStatus struct {
	encrypted   bool
	newKey      bool
	verifiedKey bool
}

func (conv *conversationPane) savePeerFingerprint(u *gtkUI) {
	conversation, exists := conv.getConversation()
	if !exists {
		conv.account.log.Warn("Conversation does not exist - this shouldn't happen")
		return
	}

	conf := conv.account.session.GetConfig()
	strP := conv.currentPeerForSending().NoResource().String()
	p, hasPeer := conf.GetPeer(strP)

	if !hasPeer {
		p = conf.EnsurePeer(strP)
	}

	p.EnsureHasFingerprint(conversation.TheirFingerprint())

	err := u.saveConfigInternal()
	if err != nil {
		conv.account.log.WithError(err).Warn("Failed to save config")
	}
}

func (conv *conversationPane) calculateNewKeyStatus() {
	conversation, exists := conv.getConversation()
	if !exists {
		conv.account.log.Warn("Conversation does not exist - this shouldn't happen")
		return
	}

	fingerprint := conversation.TheirFingerprint()

	strP := conv.currentPeerForSending().NoResource().String()

	conv.encryptionStatus.newKey = true

	p, hasPeer := conv.account.session.GetConfig().GetPeer(strP)
	if hasPeer {
		_, has := p.GetFingerprint(fingerprint)
		conv.encryptionStatus.newKey = !has
	}
}

func (conv *conversationPane) updateSecurityStatus() {
	conversation, exists := conv.getConversation()
	e := false
	if exists {
		e = conversation.IsEncrypted()
	}

	conv.encryptionStatus.encrypted = e
	if e {
		strP := conv.currentPeerForSending().NoResource().String()

		p, hasPeer := conv.account.session.GetConfig().GetPeer(strP)

		if hasPeer {
			conv.encryptionStatus.verifiedKey = p.HasTrustedFingerprint(conversation.TheirFingerprint())
		}
	} else {
		conv.encryptionStatus.newKey = false
		conv.encryptionStatus.verifiedKey = false
	}

	conv.updateIdentityVerificationWarning()
	conv.updateSecurityWarning()
}

func (conv *conversationPane) isEncrypted() bool {
	return conv.encryptionStatus.encrypted
}

func (conv *conversationPane) hasNewKey() bool {
	return conv.encryptionStatus.newKey
}

func (conv *conversationPane) hasVerifiedKey() bool {
	return conv.encryptionStatus.verifiedKey
}

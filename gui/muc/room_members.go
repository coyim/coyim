package muc

import (
	"fmt"
)

type member struct {
	*rosterItem
	role roomRole
}

type membersList map[string]*member

const (
	indexMemberStatusIcon     = 0
	indexMemberDisplayName    = 1
	indexMemberDisplayRole    = 2
	indexMemberDisplayTooltip = 3
)

func (r *roomUI) showRoomMembers() {
	for _, member := range r.room.members {
		r.addRoomMember(member)
	}

	r.roomMembersView.ExpandAll()
}

func (r *roomUI) addRoomMember(m *member) {
	parentIter := r.roomMembersModel.Append()

	_ = r.roomMembersModel.SetValue(parentIter, indexMemberDisplayName, m.displayName())
	_ = r.roomMembersModel.SetValue(parentIter, indexMemberDisplayRole, m.displayRole())
	_ = r.roomMembersModel.SetValue(parentIter, indexMemberDisplayTooltip, m.displayInfo())
	_ = r.roomMembersModel.SetValue(parentIter, indexMemberStatusIcon, statusIcons[m.getStatus()].GetPixbuf())
}

func (m *member) displayInfo() string {
	if m.rosterItem.name != "" {
		return fmt.Sprintf("%s [%s]", m.rosterItem.name, m.rosterItem.id)
	}

	return m.rosterItem.id
}

func (m *member) displayRole() string {
	if m.role == roleAdministrator {
		return "Administrator"
	}

	if m.role == roleModerator {
		return "Moderator"
	}

	if m.role == roleParticipant {
		return "Participant"
	}

	if m.role == roleVisitor {
		return "Visitor"
	}

	return "None"
}

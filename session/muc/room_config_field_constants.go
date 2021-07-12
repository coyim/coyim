package muc

import (
	"strings"
)

// RoomConfigFieldType represents a known field
type RoomConfigFieldType int

const (
	// RoomConfigFieldUnexpected represents an unknown field
	RoomConfigFieldUnexpected RoomConfigFieldType = iota
	// RoomConfigFieldName represents the "name" field
	RoomConfigFieldName
	// RoomConfigFieldDescription represents the "description" field
	RoomConfigFieldDescription
	// RoomConfigFieldEnableLogging represents the room "enable logging" field
	RoomConfigFieldEnableLogging
	// RoomConfigFieldLanguage represents the "language" field
	RoomConfigFieldLanguage
	// RoomConfigFieldPubsub represents the "pubsub" field
	RoomConfigFieldPubsub
	// RoomConfigFieldCanChangeSubject represents the "can change subject" field
	RoomConfigFieldCanChangeSubject
	// RoomConfigFieldAllowInvites represents the "allow invites" field
	RoomConfigFieldAllowInvites
	// RoomConfigFieldAllowPrivateMessages represents the "allow private messages" field
	RoomConfigFieldAllowPrivateMessages
	// RoomConfigFieldMaxOccupantsNumber represents the "max occupants number" field
	RoomConfigFieldMaxOccupantsNumber
	// RoomConfigFieldIsPublic represents the "public" field
	RoomConfigFieldIsPublic
	// RoomConfigFieldIsPersistent represents the "persistent" field
	RoomConfigFieldIsPersistent
	// RoomConfigFieldPresenceBroadcast represents the "presence broadcast" field
	RoomConfigFieldPresenceBroadcast
	// RoomConfigFieldIsModerated represents the "moderated" field
	RoomConfigFieldIsModerated
	// RoomConfigFieldIsMembersOnly represents the "members only" field
	RoomConfigFieldIsMembersOnly
	// RoomConfigFieldMembers represents the "members" field
	RoomConfigFieldMembers
	// RoomConfigFieldIsPasswordProtected represents the "password protected" field
	RoomConfigFieldIsPasswordProtected
	// RoomConfigFieldPassword represents the "password" field
	RoomConfigFieldPassword
	// RoomConfigFieldOwners represents the "owners list" field
	RoomConfigFieldOwners
	// RoomConfigFieldWhoIs represents the "who is" field
	RoomConfigFieldWhoIs
	// RoomConfigFieldMaxHistoryFetch represents the "max history fecth" field
	RoomConfigFieldMaxHistoryFetch
	// RoomConfigFieldAdmins represents the "admins list" field
	RoomConfigFieldAdmins
	// RoomConfigFieldAllowQueryUsers represents the "allow query users" field
	RoomConfigFieldAllowQueryUsers
	// RoomConfigFieldAllowVisitorStatus represents the "allow visitor status" field
	RoomConfigFieldAllowVisitorStatus
	// RoomConfigFieldAllowVisitorNickchange represents the "allow visitor nickchange" field
	RoomConfigFieldAllowVisitorNickchange
	// RoomConfigFieldAllowVoiceRequest represents the "allow voice request" field
	RoomConfigFieldAllowVoiceRequest
	// RoomConfigFieldAllowSubscription represents the "allow subscription" field
	RoomConfigFieldAllowSubscription
	// RoomConfigFieldMembersByDefault represents the "members by default" field
	RoomConfigFieldMembersByDefault
	// RoomConfigFieldVoiceRequestMinInteval represents the "voice request minimum interval" field
	RoomConfigFieldVoiceRequestMinInteval
	// RoomConfigAllowPrivateMessagesFromVisitors represents the "allow private messages from visitors" field
	RoomConfigAllowPrivateMessagesFromVisitors
	// RoomConfigPublicList represents the "public list" field
	RoomConfigPublicList
)

var roomConfigKnownFields = map[RoomConfigFieldType][]string{
	RoomConfigFieldName:                        {configFieldRoomName},
	RoomConfigFieldDescription:                 {configFieldRoomDescription},
	RoomConfigFieldEnableLogging:               {configFieldEnableLogging, configFieldEnableArchiving, configFieldMessageArchiveManagement},
	RoomConfigFieldLanguage:                    {configFieldLanguage},
	RoomConfigFieldPubsub:                      {configFieldPubsub},
	RoomConfigFieldCanChangeSubject:            {configFieldCanChangeSubject},
	RoomConfigFieldAllowInvites:                {configFieldAllowInvites, configFieldAllowMemberInvites},
	RoomConfigFieldAllowPrivateMessages:        {configFieldAllowPM, configFieldAllowPrivateMessages},
	RoomConfigFieldMaxOccupantsNumber:          {configFieldMaxOccupantsNumber},
	RoomConfigFieldIsPublic:                    {configFieldIsPublic},
	RoomConfigFieldIsPersistent:                {configFieldIsPersistent},
	RoomConfigFieldPresenceBroadcast:           {configFieldPresenceBroadcast},
	RoomConfigFieldIsModerated:                 {configFieldModerated},
	RoomConfigFieldIsMembersOnly:               {configFieldMembersOnly},
	RoomConfigFieldMembers:                     {configFieldMemberList},
	RoomConfigFieldIsPasswordProtected:         {configFieldPasswordProtected},
	RoomConfigFieldPassword:                    {configFieldPassword},
	RoomConfigFieldOwners:                      {configFieldOwners},
	RoomConfigFieldWhoIs:                       {configFieldWhoIs},
	RoomConfigFieldMaxHistoryFetch:             {configFieldMaxHistoryFetch, configFieldMaxHistoryLength},
	RoomConfigFieldAdmins:                      {configFieldRoomAdmins},
	RoomConfigFieldAllowQueryUsers:             {allowQueryUsers},
	RoomConfigFieldAllowVisitorStatus:          {allowVisitorStatus},
	RoomConfigFieldAllowVisitorNickchange:      {allowVisitorNickchange},
	RoomConfigFieldAllowVoiceRequest:           {allowVoiceRequests},
	RoomConfigFieldAllowSubscription:           {allowSubscription},
	RoomConfigFieldMembersByDefault:            {membersByDefault},
	RoomConfigFieldVoiceRequestMinInteval:      {voiceRequestMinInterval},
	RoomConfigAllowPrivateMessagesFromVisitors: {allowPrivateMessagesFromVisitors},
	RoomConfigPublicList:                       {publicList},
}

type roomConfigFieldsNames []string

func (l roomConfigFieldsNames) includes(fieldName string) bool {
	for _, fn := range l {
		if fn == fieldName {
			return true
		}
	}
	return false
}

func getKnownRoomConfigFieldKey(fieldName string) (RoomConfigFieldType, bool) {
	for key, fieldNames := range roomConfigKnownFields {
		names := roomConfigFieldsNames(fieldNames)
		if names.includes(fieldName) {
			return key, true
		}
	}
	return RoomConfigFieldUnexpected, false
}

const (
	// RoomConfigFieldText represents a "text-single" config field type
	RoomConfigFieldText = "text-single"
	// RoomConfigFieldTextPrivate represents a "text-private" config field type
	RoomConfigFieldTextPrivate = "text-private"
	// RoomConfigFieldTextMulti represents a "text-multi" config field type
	RoomConfigFieldTextMulti = "text-multi"
	// RoomConfigFieldBoolean represents a "boolean" config field type
	RoomConfigFieldBoolean = "boolean"
	// RoomConfigFieldList represents a "list-single" config field type
	RoomConfigFieldList = "list-single"
	// RoomConfigFieldListMulti represents a "list-multi" config field type
	RoomConfigFieldListMulti = "list-multi"
	// RoomConfigFieldJidMulti represents a "jid-multi" config field type
	RoomConfigFieldJidMulti = "jid-multi"
	// RoomConfigFieldFixed represents a "fixed" config field type
	RoomConfigFieldFixed = "fixed"
	// RoomConfigFieldHidden represents a "hidden" config field type
	RoomConfigFieldHidden = "hidden"
)

const (
	// RoomConfigOptionModerators represents the field option for "moderators"
	RoomConfigOptionModerators = "moderators"
	// RoomConfigOptionParticipants represents the field option for "participants"
	RoomConfigOptionParticipants = "participants"
	// RoomConfigOptionNobody represents the field opion for "nobody"
	RoomConfigOptionNobody = "nobody"
	// RoomConfigOptionAnyone represents the field opion for "anyone"
	RoomConfigOptionAnyone = "anyone"
	// RoomConfigOptionModerator represents the field option for "moderator"
	RoomConfigOptionModerator = "moderator"
	// RoomConfigOptionParticipant represents the field option for "participant"
	RoomConfigOptionParticipant = "participant"
	// RoomConfigOptionVisitor represents the field option for "visitor"
	RoomConfigOptionVisitor = "visitor"
	// RoomConfigOptionNone represents the field option for "none"
	RoomConfigOptionNone = "0"
	// RoomConfigOption10 represents the field option for "10"
	RoomConfigOption10 = "10"
	// RoomConfigOption20 represents the field option for "20"
	RoomConfigOption20 = "20"
	// RoomConfigOption30 represents the field option for "30"
	RoomConfigOption30 = "30"
	// RoomConfigOption50 represents the field option for "50"
	RoomConfigOption50 = "50"
	// RoomConfigOption100 represents the field option for "100"
	RoomConfigOption100 = "100"
)

var maxHistoryFetchDefaultOptions = []*RoomConfigFieldOption{
	newRoomConfigFieldOption(RoomConfigOption10, RoomConfigOption10),
	newRoomConfigFieldOption(RoomConfigOption20, RoomConfigOption20),
	newRoomConfigFieldOption(RoomConfigOption30, RoomConfigOption30),
	newRoomConfigFieldOption(RoomConfigOption50, RoomConfigOption50),
	newRoomConfigFieldOption(RoomConfigOption100, RoomConfigOption100),
	newRoomConfigFieldOption(RoomConfigOptionNone, RoomConfigOptionNone),
}

// SubmitFormError contains information related to a bad request error caused by a room configuration form submitted
type SubmitFormError struct {
	err        error
	fieldError RoomConfigFieldType
}

// NewSubmitFormError returns a new SubmitFormError with the specified error
func NewSubmitFormError(err error) *SubmitFormError {
	return &SubmitFormError{
		err:        err,
		fieldError: RoomConfigFieldUnexpected,
	}
}

// Error returns the error from SubmitFormError
func (sfe *SubmitFormError) Error() error {
	return sfe.err
}

// Field returns specific information about field in bad request error
func (sfe *SubmitFormError) Field() RoomConfigFieldType {
	return sfe.fieldError
}

// SetFieldErrorBadResponseText tries to identify a known field within the error message of a bad request.
// It assigns the found known field in fieldError
func (sfe *SubmitFormError) SetFieldErrorBadResponseText(errorText string) {
	for key, fieldNames := range roomConfigKnownFields {
		for _, fieldName := range fieldNames {
			if strings.Contains(errorText, fieldName) {
				sfe.fieldError = key
				return
			}
		}
	}
}

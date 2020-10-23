package data


// RoomConfig represents the room's configuration values
type RoomConfig struct {
	SupportsVoiceRequests     bool
	AllowsRegistration        bool
	Anonymity                 string
	Persistent                bool
	Moderated                 bool
	Open                      bool
	PasswordProtected         bool
	Public                    bool
	Language                  string
	OccupantsCanChangeSubject bool
	Title                     string
	Description               string
	Occupants                 int
	MembersCanInvite          bool
	OccupantsCanInvite        bool
	AllowPrivateMessages      string // This can be 'anyone', 'participants', 'moderators', 'none'
	ContactJid                string
	Logged                    bool // Notice that this will not always be correct for all servers
}

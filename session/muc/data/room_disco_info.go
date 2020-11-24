package data

// RoomDiscoInfo represents the room discovery information values
type RoomDiscoInfo struct {
	SupportsVoiceRequests bool
	AllowsRegistration    bool
	AnonymityLevel        string
	Persistent            bool
	Moderated             bool
	Open                  bool
	PasswordProtected     bool
	Public                bool

	Language                  string
	OccupantsCanChangeSubject bool
	Logged                    bool // Notice that this will not always be correct for all servers
	Title                     string
	Description               string
	Occupants                 int
	MembersCanInvite          bool
	OccupantsCanInvite        bool
	AllowPrivateMessages      string // This can be 'anyone', 'participants', 'moderators', 'none'
	ContactJid                string
	MaxHistoryFetch           int
}

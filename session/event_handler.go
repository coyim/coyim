package session

import "time"

// EventHandler represents the main notifications that the session can emit
// It's really more an observer than an even handler
type EventHandler interface {
	Debug(string)
	Info(string)
	Warn(string)
	Alert(string)
	MessageReceived(s *Session, from string, timestamp time.Time, encrypted bool, message []byte)
	ProcessPresence(from, to, show, status string, gone bool)
	Subscribed(account, peer string)
	Unsubscribe(account, peer string)
	RegisterCallback(title, instructions string, fields []interface{}) error
}

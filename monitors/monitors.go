package monitors

import "github.com/syntaqx/elkish"

// Monitor ...
type Monitor interface {
	String() string
	Add(entry elkish.LogEntry)
}

package alerts

import (
	"time"

	"github.com/syntaqx/elkish"
)

// Alert ...
type Alert interface {
	Len() int
	Add(entry elkish.LogEntry)
	Check(when time.Time)
}

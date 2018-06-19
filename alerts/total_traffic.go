package alerts

import (
	"fmt"
	"io"
	"time"

	"github.com/syntaqx/elkish"
)

// TotalTrafficAlert defines an alert that will aggregate total traffic over a
// period of time.
type TotalTrafficAlert struct {
	writer      io.Writer
	threshold   int
	duration    time.Duration
	isTriggered bool
	entries     []elkish.LogEntry
}

// NewTotalTrafficAlert will create a new TotalTrafficAlert with given values.
func NewTotalTrafficAlert(writer io.Writer, threshold int, duration time.Duration) *TotalTrafficAlert {
	return &TotalTrafficAlert{
		writer:    writer,
		threshold: threshold,
		duration:  duration,
	}
}

// Len provides support for Length.
func (t *TotalTrafficAlert) Len() int {
	return len(t.entries)
}

// Add will add a new entry to the stack, and perform an alerts.
func (t *TotalTrafficAlert) Add(entry elkish.LogEntry) {
	t.entries = append(t.entries, entry)
	t.Check(entry.Date)
}

// Check will check the alert stack.
func (t *TotalTrafficAlert) Check(when time.Time) {
	t.clean()

	if t.thresholdExceeded() && !t.isTriggered {
		fmt.Fprintf(t.writer, "High traffic generated an alert - hits = %d, triggered at %s\n", len(t.entries), time.Now())
		t.isTriggered = true
	} else if !t.thresholdExceeded() && t.isTriggered {
		fmt.Fprintf(t.writer, "High traffic has recovered, triggered at %s", time.Now())
		t.isTriggered = false
	}
}

// thresholdExceeded is a simple DRY for checking if the traffic threshold has
// been exceeded
func (t *TotalTrafficAlert) thresholdExceeded() bool {
	return len(t.entries) > 0 && t.Len() >= t.threshold
}

// clean will remove outdated entries from the stack.
func (t *TotalTrafficAlert) clean() {
	var idx int
	for i := len(t.entries) - 1; i >= 0; i-- {
		if time.Now().Sub(t.entries[i].Date) > t.duration {
			idx = i + 1
			break
		}
	}

	t.entries = t.entries[idx:]
}

package monitors

import (
	"fmt"

	"github.com/syntaqx/elkish"
)

var totalsMonitorFormat = `Successes: %d
Redirects: %d
Client Failures: %d
Server Failures: %d
Responses: %d`

type TotalsMonitor struct {
	Successes      int
	Redirects      int
	ClientFailures int
	ServerFailures int
	Total          int64
}

func NewTotalsMonitor() *TotalsMonitor {
	return &TotalsMonitor{}
}

func (m *TotalsMonitor) String() string {
	return fmt.Sprintf(totalsMonitorFormat, m.Successes, m.Redirects, m.ClientFailures, m.ServerFailures, m.Total)
}

func (m *TotalsMonitor) Add(entry elkish.LogEntry) {
	switch {
	case entry.StatusCode >= 500:
		m.ServerFailures++
	case entry.StatusCode >= 400:
		m.ClientFailures++
	case entry.StatusCode >= 300:
		m.Redirects++
	case entry.StatusCode >= 200:
		m.Successes++
	}

	m.Total++
}

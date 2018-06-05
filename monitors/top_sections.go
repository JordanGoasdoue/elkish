package monitors

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/codesuki/go-time-series"
	"github.com/syntaqx/elkish"
)

type TopSectionsMonitor struct {
	interval time.Duration
	sections map[string]*timeseries.TimeSeries
}

func NewTopSectionsMonitor(interval time.Duration) *TopSectionsMonitor {
	return &TopSectionsMonitor{interval: interval, sections: make(map[string]*timeseries.TimeSeries)}
}

func (m *TopSectionsMonitor) String() string {
	var resp string
	s := make(PairList, len(m.sections))

	i := 0
	for section, series := range m.sections {
		dat, _ := series.Recent(m.interval)
		s[i] = Pair{Key: section, Value: int(dat)}
		i++
	}

	sort.Sort(sort.Reverse(s))

	max := len(s)
	if max > 5 {
		max = 5
	}

	if max > 0 {
		for i := 0; i < max; i++ {
			resp += fmt.Sprintf("%s %d\n", s[i].Key, s[i].Value)
		}
	} else {
		resp += "No sections have been recorded"
	}

	return resp
}

func (m *TopSectionsMonitor) Add(entry elkish.LogEntry) {
	section := strings.Split(strings.TrimLeft(entry.Request.Resource, "/"), "/")[0]

	if _, ok := m.sections[section]; !ok {
		ts, _ := timeseries.NewTimeSeries()
		m.sections[section] = ts
	}

	m.sections[section].Increase(1)
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

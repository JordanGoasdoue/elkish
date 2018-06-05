package elkish

import (
	"net/http"
	"testing"
	"time"
)

var entryTests = []struct {
	in  string
	out LogEntry
}{
	{
		"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report HTTP/1.0\" 200 1234",
		LogEntry{
			IPAddress:  "127.0.0.1",
			Identifier: "-",
			UserID:     "james",
			Date:       time.Date(2018, time.May, 9, 16, 0, 39, 0, time.UTC),
			Request: LogEntryRequest{
				Method:   http.MethodGet,
				Resource: "/report",
				Protocol: "HTTP/1.0",
			},
			StatusCode: http.StatusOK,
			Size:       1234,
		},
	},
	{
		"127.0.0.1 - james [09/May/2018:16:00:39 +0000] \"GET /report\" 200 1234",
		LogEntry{
			IPAddress:  "127.0.0.1",
			Identifier: "-",
			UserID:     "james",
			Date:       time.Date(2018, time.May, 9, 16, 0, 39, 0, time.UTC),
			Request: LogEntryRequest{
				Method:   http.MethodGet,
				Resource: "/report",
				Protocol: "",
			},
			StatusCode: http.StatusOK,
			Size:       1234,
		},
	},
}

func TestFlagParser(t *testing.T) {
	for _, tt := range entryTests {
		entry, err := ParseEntry(tt.in)
		if err != nil {
			t.Errorf("got err %q, expected nil", err)
		}

		if entry.String() != tt.out.String() {
			t.Errorf("got %v, want %v", entry, tt.out)
		}
	}
}

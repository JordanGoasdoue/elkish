package elkish

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// LogEntry represents an individual entry from the Common Log Format, also
// known as the NCSA Common log format (after NCSA_HTTPd).
// https://en.wikipedia.org/wiki/Common_Log_Format
type LogEntry struct {
	// IPAddress is the IP address of the client (remote host) which made the
	// request to the server.
	IPAddress string

	// Identifier is thhe RFC 1413 identity of the client.
	Identifier string

	// UserID is the userid of the person requesting the document.
	UserID string

	// Time is the date, time, and time zone that the request was received, by
	// default in strftime format %d/%b/%Y:%H:%M:%S %z
	Date time.Time

	// Request is the request line from the client. The method GET, the resource
	// requested, and the HTTP protocol.
	Request LogEntryRequest

	// StatusCode is the HTTP status code returned to the client. 2xx is a
	// successful response, 3xx a redirection, 4xx a client error, and 5xx a
	// server error.
	StatusCode int

	// Size is the size of the object returned to the client, measured in bytes.
	Size int64
}

// LogEntryRequest stores the method, the resource, and protocol for a LogEntry
// request item.
type LogEntryRequest struct {
	Method   string
	Resource string
	Protocol string
}

// String implements the stringer interface for LogEntryRequest.
func (e LogEntryRequest) String() string {
	if e.Protocol != "" {
		return fmt.Sprintf("%s %s %s", e.Method, e.Resource, e.Protocol)
	}

	return fmt.Sprintf("%s %s", e.Method, e.Resource)
}

// String implements the stringer interface for LogEntry.
func (e LogEntry) String() string {
	return fmt.Sprintf("%s %s %s [%s] \"%s\" %d %v",
		e.IPAddress,
		e.Identifier,
		e.UserID,
		e.Date.Format("02/Jan/2006:15:04:05 -0700"),
		e.Request,
		e.StatusCode,
		e.Size,
	)
}

// ParseEntry parses a given common log line into a LogEntry.
func ParseEntry(line string) (LogEntry, error) {
	var entry LogEntry

	// Remove space/newline characters from the entry.
	line = strings.TrimSpace(line)

	parts := strings.SplitN(line, " ", 4)
	entry.IPAddress = parts[0]
	entry.Identifier = parts[1]
	entry.UserID = parts[2]

	// Strip used values from the line
	line = parts[3]

	// Extract date from the first [ and last `"] "`
	parts = strings.SplitAfterN(line, "] ", 2)

	date, err := time.Parse("02/Jan/2006:15:04:05 -0700", parts[0][1:len(parts[0])-2])
	if err != nil {
		return entry, fmt.Errorf("could not parse time: %v", err)
	}
	entry.Date = date

	// Strip the date from the line
	line = parts[1]

	// extract request values, stripping the first " and last `" `, then
	// splitting on space.
	parts = strings.SplitAfterN(line, "\" ", 2)

	requestValues := strings.Fields(parts[0][1 : len(parts[0])-2])

	// Ensure we have at least a method and resource, protocol is optional.
	if len(requestValues) < 2 {
		return entry, fmt.Errorf("invalid request value length %d, expected at least 2", len(requestValues))
	}

	request := LogEntryRequest{
		Method:   requestValues[0],
		Resource: requestValues[1],
	}

	if len(requestValues) > 2 {
		request.Protocol = requestValues[2]
	}

	entry.Request = request

	// Ensure we have enough parts to finish.
	if len(parts) != 2 {
		return entry, errors.New("Truncated log")
	}

	// Strip the request from the line
	line = parts[1]

	parts = strings.SplitN(line, " ", 3)

	statusCode, err := strconv.Atoi(parts[0])
	if err != nil {
		return entry, fmt.Errorf("Unable to parse status code: %v", err)
	}
	entry.StatusCode = statusCode

	// When given a valid size, parse it.
	if parts[1] != "-" {
		bytes, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return entry, fmt.Errorf("Unable to parse size bytes: %v", err)
		}
		entry.Size = bytes
	}

	return entry, nil
}

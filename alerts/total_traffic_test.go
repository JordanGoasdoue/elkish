package alerts

import (
	"bytes"
	"testing"
	"time"

	"github.com/syntaqx/elkish"
)

func TestTotalTrafficAlert(t *testing.T) {
	var b bytes.Buffer

	entry := elkish.LogEntry{Date: time.Now()}
	alert := NewTotalTrafficAlert(&b, 1, time.Second)

	alert.Check(time.Now())

	if len(b.String()) != 0 {
		t.Errorf("Check should not alert anything, as it's empty")
	}

	alert.Add(entry)

	l := len(b.String())
	if l <= 0 {
		t.Error("Expected contents to be written to buffer")
	}

	time.Sleep(2 * time.Second)
	alert.Check(time.Now())

	if len(b.String()) <= l {
		t.Errorf("Expected more content to be written to buffer when notification ended.")
	}
}

func TestTotalTrafficAlertMultiple(t *testing.T) {
	var b bytes.Buffer

	entry := elkish.LogEntry{Date: time.Now()}
	alert := NewTotalTrafficAlert(&b, 5, time.Minute)

	alert.Check(time.Now())

	if len(b.String()) != 0 {
		t.Errorf("Check should not alert anything, as it's empty")
	}

	alert.Add(entry)
	alert.Add(entry)
	alert.Add(entry)
	alert.Add(entry)

	if len(b.String()) != 0 {
		t.Errorf("Not alert should have occured yet")
	}

	alert.Add(entry)

	if len(b.String()) <= 0 {
		t.Error("Expected contents to be written to buffer")
	}
}

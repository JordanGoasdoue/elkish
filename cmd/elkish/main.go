package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hpcloud/tail"
	"github.com/syntaqx/elkish"
	"github.com/syntaqx/elkish/alerts"
	"github.com/syntaqx/elkish/monitors"
)

var filepath string
var alertThreshold int
var alertDuration time.Duration

func init() {
	flag.StringVar(&filepath, "filepath", "/var/log/access.log", "Log file to monitor")
	flag.IntVar(&alertThreshold, "alert-threshold", 10, "Total traffic threshhold amount for a given alertDuration")
	flag.DurationVar(&alertDuration, "alert-duration", 2*time.Minute, "Total traffic durration for a given alertThreshold")
}

func main() {
	flag.Parse()

	// Ensure they gave us a filepath.
	if len(filepath) < 1 {
		log.Fatalf("You must specify a valid log file to monitor\n")
	}

	// Begin tailing the file for new line entries, reopening when truncated.
	t, err := tail.TailFile(filepath, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		log.Fatalf("Could not start monitoring the given file `%s`: %v\n", filepath, err)
	}

	// Initialize our alerters.
	totalTrafficAlert := alerts.NewTotalTrafficAlert(os.Stdout, alertThreshold, alertDuration)

	// Initialize our monitors.
	topSectionsMonitor := monitors.NewTopSectionsMonitor(11 * time.Second) // offset for sleep time.
	totalsMonitor := monitors.NewTotalsMonitor()

	// Background tasks for alerting/monitoring.
	go func() {
		for {
			totalTrafficAlert.Check(time.Now())
			time.Sleep(time.Millisecond * 500) // Sleep so we don't hog the CPU
		}
	}()

	// As new lines are written to the log file, inform alerting and monitoring.
	// How they handle this information is agnostic of this loop.
	go func() {
		for line := range t.Lines {
			entry, err := elkish.ParseEntry(line.Text)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Unable to parse log entry `%s`: %v\n", line, err)
				continue
			}

			totalTrafficAlert.Add(entry)
			topSectionsMonitor.Add(entry)
			totalsMonitor.Add(entry)
		}
	}()

	// Every 10 seconds, request information from our monitors.
	for {
		time.Sleep(10 * time.Second)

		fmt.Fprintf(os.Stdout, "## Top Sections:\n\n")
		fmt.Fprintf(os.Stdout, "%s\n", topSectionsMonitor)
		fmt.Fprintf(os.Stdout, "\n")

		fmt.Fprintf(os.Stdout, "## Totals:\n\n")
		fmt.Fprintf(os.Stdout, "%s\n", totalsMonitor)
		fmt.Fprintf(os.Stdout, "\n")
	}
}

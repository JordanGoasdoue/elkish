# Elk'ish

[CLF]: https://en.wikipedia.org/wiki/Common_Log_Format

> An interview puzzle to create a simple HTTP log monitoring program.

A simple program that mimics alerting and monitoring behavior of the ELK stack,
done for an interview. This program will monitor a given [HTTP access log][CLF]
file and provide alerting and monitoring.

## Challenge Requirements

* Consume an actively written-to w3c-formatted HTTP access log. It should
  default to reading /var/log/access.log and be overridable.
* Display stats every 10s about the traffic during those 10s: the sections of
  the web site with the most hits, as well as interesting summary statistics on
  the traffic as a whole. A section is defined as being what's before the
  second `/` in the path. For example, the section for
  `http://my.site.com/pages/create` is `http://my.site.com/pages`.
* Make sure a user can keep the app running and monitor the log file
  continuously
* Whenever total traffic for the past 2 minutes exceeds a certain number on
  average, add a message saying that:
  `High traffic generated an alert - hits = {value}, triggered at {time}`.
  The default threshold should be 10 requests per second and should be overridable.
* Whenever the total traffic drops again below that value on average for the
  past 2 minutes, print or displays another message detailing when the alert
  recovered.
* Write a test for the alerting logic.
* Explain how you’d improve on this application design.
* If you have access to a linux docker environment, we'd love to be able to
  docker build and run your project! If you don't though, don't sweat it.

## General Improvement Plan

This code is meant to reimplement some standard functionality that I would
generally not implement myself:

* Total traffic alerting implements its own "cleanup" concept, where as the top
  section monitor implements a simple in-memory time series database. I would
  use ElasticSearch for both of these, as there's really no reason to try and
  rewrite that platform.
* Monitors and Alerts aren't interfaced. Although I kept with a common pattern,
  I didn't account for some `Add` functions needing errorability, and given my
  attempt to timebox this, I didn't want to iterate on it forever. There's some
  room for improvement there though.
* Creating new monitors and alerts shouldn't require code. Perhaps a DSL of some
  sort, or just relying on Kibana/Grafana rather than implementing this. ELK is
  pretty powerful without needing to reinvent the wheel.
* Test coverage of Monitors is non-existant, but I believe they work.

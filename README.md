# Elk'ish

[CLF]: https://en.wikipedia.org/wiki/Common_Log_Format

> An interview puzzle to create a simple HTTP log monitoring program.

A simple program that mimics alerting and monitoring behavior of the ELK stack,
done for an interview. This program will monitor a given [HTTP access log][CLF]
file and provide alerting and monitoring.

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

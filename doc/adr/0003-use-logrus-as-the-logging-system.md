# 3. Use logrus as the logging system

Date: 2023-06-24

## Status

Accepted

## Context

The project needs a consistent logging platform.  I have used logrus before.

### External References that Influenced my Decision

- [Logging in Go: Choosing a System and Using it](https://www.honeybadger.io/blog/golang-logging/), [Ayooluwa Isaiah](https://www.honeybadger.io/blog/golang-logging/#authorDetails), 2020-04-01

## Decision

This project will use the `github.com/sirupsen/logrus` package to provide a
standard logging framework for capturing important application debug and error
information.

The application will log in `JSON` format to provide a consistant, machine
readable format that is easy to ingest by log aggregation systems.

## Consequences

- `github.com/sirupsen/logrus` is now a dependency that must be monitored for
  vulnerabilities and license compliance.
# clog - Go package for Canonical Logging

[![GoDoc](https://godoc.org/github.com/jwilder/clog?status.svg)](https://godoc.org/github.com/jwilder/clog)

clog provides a simple way to support canonical logging in Go applications for easily
constructing wide events.  Wide events, also called canonical logs, are
a way to log structured data that can be easily queried and analyzed.  This package
makes it easy to construct wide events by providing a simple API for setting and
retrieving structured logging data.  Using canonical logs complements (and can sometimes replace)
typical places where you might use metrics, logs or traces.  This approach is central
to Observability 2.0 ideas.

Learn more:

* [A Practitioner's Guide to Wide Events](https://jeremymorrell.dev/blog/a-practitioners-guide-to-wide-events/)
* [All you need is Wide Events, not “Metrics, Logs and Traces”](https://isburmistrov.substack.com/p/all-you-need-is-wide-events-not-metrics)
* [Observability wide events 101](https://boristane.com/blog/observability-wide-events-101/)
* [Instrumenting distributed systems for operational visibility](https://aws.amazon.com/builders-library/instrumenting-distributed-systems-for-operational-visibility/)
* [Observability 2.0](https://charity.wtf/tag/observability-2-0/)

# Installation

```bash
go get github.com/jwilder/clog
```

# Usage

```go
// Initialize the canonical logging context. This should be done at the 
// beginning of your unit of work (http request, background process, etc).
ctx := clog.Init(context.Background())

// Set some values in the logging context as needed.  As long as the context
// is passed around, these values will be aggregated into the same 
// logging context.  Dotted key names are used to group related values.
clog.SetString(ctx, "http.request.method", "GET")
clog.SetString(ctx, "http.request.path", "/example")
clog.SetInt(ctx, "http.response.status_code", 200)
clog.SetFloat64(ctx, "http.response.duration_ms", 123.45)

// Add values to existing keys
clog.AddInt(ctx, "http.response.status_code", 1)
clog.AddFloat64(ctx, "http.response.duration_ms", 10.5)

// Marshal the logging context to JSON
log := clog.MarshalJSON(ctx)
fmt.Println(log)
```

# License
 MIT

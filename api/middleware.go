package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/trace"
)

var traceIDRegex = regexp.MustCompile("(.*)/([^;]*)")

const (
	traceIDRegexGroups = 2
	spanIDBase         = 10
	spanIDBitSize      = 64
)

type invalidTraceHeaderError struct {
	Message     string
	TraceHeader string
}

func (e *invalidTraceHeaderError) Error() string {
	return fmt.Sprintf("%s: '%s'", e.Message, e.TraceHeader)
}

type serveHTTP func(http.ResponseWriter, *http.Request)

func (a *API) traceMiddleware(next serveHTTP) serveHTTP {
	return func(writer http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		sc, err := parseTraceHeader(req.Header.Get("X-Cloud-Trace-Context"))
		if err != nil {
			a.log.WithError(err).Error("could not parse trace header")
		} else {
			ctx = trace.ContextWithSpanContext(ctx, *sc)
		}

		next(writer, req.WithContext(ctx))
	}
}

func parseTraceHeader(header string) (*trace.SpanContext, error) {
	if header == "" {
		return nil, &invalidTraceHeaderError{"empty trace header", header}
	}

	traceIDSpanID := traceIDRegex.FindStringSubmatch(header)

	if len(traceIDSpanID) != traceIDRegexGroups+1 {
		return nil, &invalidTraceHeaderError{"could not match pattern", header}
	}

	spanIDInt, err := strconv.ParseUint(traceIDSpanID[2], spanIDBase, spanIDBitSize)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse span ID integer")
	}

	spanID, err := trace.SpanIDFromHex(fmt.Sprintf("%x", spanIDInt))
	if err != nil {
		return nil, errors.Wrap(err, "could not parse span ID from hex string")
	}

	traceID, err := trace.TraceIDFromHex(traceIDSpanID[1])
	if err != nil {
		return nil, errors.Wrap(err, "could not parse trace ID from hex string")
	}

	sc := trace.NewSpanContext(trace.SpanContextConfig{TraceID: traceID, SpanID: spanID})

	return &sc, nil
}

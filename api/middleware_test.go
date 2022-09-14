package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFullTraceHeader(t *testing.T) {
	require := require.New(t)

	// When run with a trace header containing traceID, spanID, and collection flag
	sc, err := parseTraceHeader("105445aa7843bc8bf206b12000100000/15641949547778198677;o=1")

	// Then no errors are returned
	require.NoError(err, "unexpected error")

	// And the traceID string is returned as is from the trace header
	require.Equal("105445aa7843bc8bf206b12000100000", sc.TraceID().String(), "wrong TraceID")

	// And the spanIDis the returned as a the hex value of the trace header value
	require.Equal("d9135e3d35e6cc95", sc.SpanID().String(), "wrong SpanID")
}

func TestParseShortTraceHeader(t *testing.T) {
	require := require.New(t)

	// When run with a trace header containing traceID and spanID
	sc, err := parseTraceHeader("ef4dddcfbe73ca491f98048cdb0561fc/12602853554865825763")

	// Then no errors are returned
	require.NoError(err, "unexpected error")

	// And the traceID string is returned as is from the trace header
	require.Equal("ef4dddcfbe73ca491f98048cdb0561fc", sc.TraceID().String(), "wrong TraceID")

	// And the spanIDis the returned as a the hex value of the trace header value
	require.Equal("aee654890e0be7e3", sc.SpanID().String(), "wrong SpanID")
}

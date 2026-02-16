package stream

import "errors"

var (
	// the client disconnected before the stream completed.
	ErrStreamClosed = errors.New("stream closed by client")

	// a block fetch exceeded the timeout.
	ErrBlockTimeout = errors.New("block fetch timed out")

	// all retry attempts failed.
	ErrMaxRetriesExceeded = errors.New("max retries exceeded")

	// the pipe was closed and all data was consumed.
	ErrPipeDrained = errors.New("pipe drained")
)

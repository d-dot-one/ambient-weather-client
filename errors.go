package awn

import (
	"errors"
	"fmt"
)

type errorType int

const (
	_ errorType = iota // so we don't start at 0
	errContextTimeoutExceeded
	errMalformedDate
	errRegexFailed
	errAPIKeyMissing
	errAppKeyMissing
	errInvalidDateFormat
	errMacAddressMissing
)

var (
	ErrContextTimeoutExceeded = ClientError{kind: errContextTimeoutExceeded} //nolint:exhaustruct
	ErrMalformedDate          = ClientError{kind: errMalformedDate}          //nolint:exhaustruct
	ErrRegexFailed            = ClientError{kind: errRegexFailed}            //nolint:exhaustruct
	ErrAPIKeyMissing          = ClientError{kind: errAPIKeyMissing}          //nolint:exhaustruct
	ErrAppKeyMissing          = ClientError{kind: errAppKeyMissing}          //nolint:exhaustruct
	ErrInvalidDateFormat      = ClientError{kind: errInvalidDateFormat}      //nolint:exhaustruct
	ErrMacAddressMissing      = ClientError{kind: errMacAddressMissing}      //nolint:exhaustruct
)

// ClientError is a public custom error type that is used to return errors from the client.
type ClientError struct {
	kind  errorType // errKind in example
	value int
	err   error
}

// todo: should all of the be passing a pointer?

// Error is a public function that returns the error message.
func (c ClientError) Error() string {
	switch c.kind {
	case errContextTimeoutExceeded:
		return fmt.Sprintf("context timeout exceeded: %v", c.value)
	case errMalformedDate:
		return fmt.Sprintf("date format is malformed. should be YYYY-MM-DD: %v", c.value)
	case errRegexFailed:
		return fmt.Sprintf("regex failed: %v", c.value)
	case errAPIKeyMissing:
		return fmt.Sprintf("api key is missing: %v", c.value)
	case errAppKeyMissing:
		return fmt.Sprintf("application key is missing: %v", c.value)
	case errInvalidDateFormat:
		return fmt.Sprintf("date is invalid. It should be in epoch time in milliseconds: %v", c.value)
	case errMacAddressMissing:
		return fmt.Sprintf("mac address is missing: %v", c.value)
	default:
		return fmt.Sprintf("unknown error: %v", c.value)
	}
}

// from is a private function that returns an error with a particular location and the
// underlying error.
func (c ClientError) from(pos int, err error) ClientError {
	ce := c
	ce.value = pos
	ce.err = err
	return ce
}

// with is a private function that returns an error with a particular value.
func (c ClientError) with(val int) ClientError {
	ce := c
	ce.value = val
	return ce
}

// Is is a public function that reports whether any error in the error's chain matches target.
func (c ClientError) Is(err error) bool {
	var clientError ClientError
	ok := errors.As(err, &clientError) // reflection
	if !ok {
		return false
	}

	return clientError.kind == c.kind
}

// Unwrap is a public function that returns the underlying error by unwrapping it.
func (c ClientError) Unwrap() error {
	return c.err
}

// Wrap is a public function that allows for errors to be propagated up correctly.
func (c ClientError) Wrap() error {
	return fmt.Errorf("error: %w", c)
}

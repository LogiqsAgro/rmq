package web

import (
	"errors"
	"fmt"
)

var (
	errRequest  *requestError  = newRequestErrorf(nil, "request processor failed")
	errResponse *responseError = newResponseErrorf(nil, "response processor failed")
)

// IsRequestProcessorError returns true if the error is caused by
// one of the RequestProcessor functions returning an error
func IsRequestProcessorError(e error) bool {
	return errors.Is(e, errRequest)
}

// IsResponseProcessorError returns true if the error is caused by
// one of the ResponseProcessor functions returning an error
func IsResponseProcessorError(e error) bool {
	return errors.Is(e, errResponse)
}

type requestError struct {
	msg string
	err error
}

func newRequestErrorf(wrapped error, format string, args ...interface{}) *requestError {
	return &requestError{
		msg: getErrMsg(wrapped, format, args...),
		err: wrapped,
	}
}

func (e *requestError) Error() string { return e.msg }
func (e *requestError) Unwrap() error { return e.err }
func (e *requestError) Is(target error) bool {
	_, ok := target.(*requestError)
	return ok
}

type responseError struct {
	msg string
	err error
}

func newResponseErrorf(wrapped error, format string, args ...interface{}) *responseError {
	return &responseError{
		msg: getErrMsg(wrapped, format, args...),
		err: wrapped,
	}
}

func (e *responseError) Error() string { return e.msg }
func (e *responseError) Unwrap() error { return e.err }
func (e *responseError) Is(target error) bool {
	_, ok := target.(*responseError)
	return ok
}

func getErrMsg(wrapped error, format string, args ...interface{}) string {
	msg := ""
	if len(format) == 0 && len(args) == 0 {
		if wrapped != nil {
			msg = wrapped.Error()
		} else {
			msg = "response processor failed"
		}
	} else {
		msg = fmt.Sprintf(format, args...)
	}
	return msg
}

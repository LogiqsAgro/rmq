package web

//go:generate go run ..\web-gen\main.go package=rs

import (
	"context"
	"io"
	"net/http"
)

type (
	ResponseProcessor func(*http.Response) error

	Response interface {
		// MaxSize sets the maximum number of bytes read from the body, this is a client-side limit.
		// By default this is set to 1MB (see )
		MaxSize(size int64) Response

		// Ensure sets the ResponseProcessor functions the are run after the response has been received,
		// and before the body is processed. If any of functions return a non-nil error,
		// response processing is aborted
		Ensure(h ...ResponseProcessor) Response

		// Body set the handler used to process the response body
		Body(h func(context.Context, io.Reader) error) Response

		// BodyFile saves the the response body to the given filepath.
		// The file is truncated and overwritten when it already exists.
		BodyFile(path string) Response

		// BodyCopyTo copies the response body to the given io.Writer.
		BodyCopyTo(w io.Writer) Response

		// BodyEncoding sets the body encoding to use for decoding the response ans sets the Accept header
		// if more than one encoding is specified, the Accept header will be a quality value list
		// ( https://developer.mozilla.org/en-US/docs/Glossary/Quality_values )
		BodyEncoding(encoding ...ResponseEncoding) Response

		// BodyDecode reads the encoded value from the response
		// body and stores it in the value pointed to by v.
		// If multiple encodings are specified the encoding will be selected
		// based on the response Content-Type header
		BodyDecode(v interface{}, encoding ...ResponseEncoding) Response

		// apply response configuration
		apply(configure ...func(Response))
	}
)

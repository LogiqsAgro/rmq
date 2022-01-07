package rs

import (
	"context"
	"io"

	"github.com/LogiqsAgro/rmq/web"
)

// MaxSize sets the maximum number of bytes read from the body, this is a client-side limit.
// By default this is set to 1MB (see )
func MaxSize(size int64) func(web.Response) {
	return func(r web.Response) { r.MaxSize(size) }
}

// Ensure sets the ResponseProcessor functions the are run after the response has been received,
// and before the body is processed. If any of functions return a non-nil error,
// response processing is aborted
func Ensure(h ...web.ResponseProcessor) func(web.Response) {
	return func(r web.Response) { r.Ensure(h...) }
}

// Body set the handler used to process the response body
func Body(h func(context.Context, io.Reader) error) func(web.Response) {
	return func(r web.Response) { r.Body(h) }
}

// BodyFile saves the the response body to the given filepath.
// The file is truncated and overwritten when it already exists.
func BodyFile(path string) func(web.Response) {
	return func(r web.Response) { r.BodyFile(path) }
}

// BodyCopyTo copies the response body to the given io.Writer.
func BodyCopyTo(w io.Writer) func(web.Response) {
	return func(r web.Response) { r.BodyCopyTo(w) }
}

// BodyEncoding sets the body encoding to use for decoding the response ans sets the Accept header
// if more than one encoding is specified, the Accept header will be a quality value list
// ( https://developer.mozilla.org/en-US/docs/Glossary/Quality_values )
func BodyEncoding(encoding ...web.ResponseEncoding) func(web.Response) {
	return func(r web.Response) { r.BodyEncoding(encoding...) }
}

// BodyDecode reads the encoded value from the response
// body and stores it in the value pointed to by v.
// If multiple encodings are specified the encoding will be selected
// based on the response Content-Type header
func BodyDecode(v interface{}, encoding ...web.ResponseEncoding) func(web.Response) {
	return func(r web.Response) { r.BodyDecode(v, encoding...) }
}

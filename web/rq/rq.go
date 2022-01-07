package rq

import (
	"io"
	urlpkg "net/url"

	"github.com/LogiqsAgro/rmq/web"
)

// Method sets the request httmp method
func Method(method string) func(web.Request) {
	return func(r web.Request) { r.Method(method) }
}

// BaseUrl sets the request base url
func BaseURL(baseURL string) func(web.Request) {
	return func(r web.Request) { r.BaseURL(baseURL) }
}

// Scheme sets the request url scheme
func Scheme(scheme string) func(web.Request) {
	return func(r web.Request) { r.Scheme(scheme) }
}

// Host sets the request url host portion
func Host(host string) func(web.Request) {
	return func(r web.Request) { r.Host(host) }
}

// Hostf sets the request url host portion to the formatted host string
func Hostf(host string, args ...interface{}) func(web.Request) {
	return func(r web.Request) { r.Hostf(host, args...) }
}

// HostAndPort sets the request url host portion to the provided host and port value
func HostAndPort(host string, port int) func(web.Request) {
	return func(r web.Request) { r.HostAndPort(host, port) }
}

// Path appends a path segment to the request base url, if it starts with a /, the path is replaced
// Path segments are appended using the following code:
// for _, p := range pathSegments {
//   baseUrl.Path = baseUrl.ResolveReference(&url.URL{Path: p}).Path
// }
func Path(path string) func(web.Request) {
	return func(r web.Request) { r.Path(path) }
}

// Pathf formats the path with the given args, and calls Path(string)
func Pathf(path string, args ...interface{}) func(web.Request) {
	return func(r web.Request) { r.Pathf(path, args...) }
}

// Param sets the values for the given query parameter, overwriting existing values in the url.
func Param(name string, values ...string) func(web.Request) {
	return func(r web.Request) { r.Param(name, values...) }
}

// Fragment sets the url fragment value
func Fragment(fragment string) func(web.Request) {
	return func(r web.Request) { r.Fragment(fragment) }
}

// ContentType sets the Content-Type header
func ContentType(contentType string) func(web.Request) {
	return func(r web.Request) { r.ContentType(contentType) }
}

// Accept sets the Accept header
func Accept(contentType string) func(web.Request) {
	return func(r web.Request) { r.Accept(contentType) }
}

// AcceptRange sets the Accept header with the given content types, sorted by their quality value
func AcceptRange(contentTypes map[string]float64) func(web.Request) {
	return func(r web.Request) { r.AcceptRange(contentTypes) }
}

// BasicAuth sets the Authorization header to Basic, with the base64 encoded user and password
func BasicAuth(user, password string) func(web.Request) {
	return func(r web.Request) { r.BasicAuth(user, password) }
}

// BearerAuth sets the Authorization header to Bearer, with the given token.
func BearerAuth(token string) func(web.Request) {
	return func(r web.Request) { r.BearerAuth(token) }
}

// Header sets the given header, overwriting existing values
func Header(key string, values ...string) func(web.Request) {
	return func(r web.Request) { r.Header(key, values...) }
}

// Body uses the getBody function to read the request body contents.
// getBody may be called multiple times in case of request redirection
func Body(getBody web.GetBody) func(web.Request) {
	return func(r web.Request) { r.Body(getBody) }
}

// BodyCached uses getBody once, caching the resulting bytes should they be needed again.
func BodyCached(getBody web.GetBody) func(web.Request) {
	return func(r web.Request) { r.BodyCached(getBody) }
}

// BodyReader uses body once to read the request body, caching the resulting bytes should they be needed again.
func BodyReader(body io.Reader) func(web.Request) {
	return func(r web.Request) { r.BodyReader(body) }
}

// BodyBytes uses body as the request body, modifying the contents of
// the slice isn't recommended, as BodyBytes doesn't make a copy
func BodyBytes(body []byte) func(web.Request) {
	return func(r web.Request) { r.BodyBytes(body) }
}

// BodyString uses body as the request body
func BodyString(body string) func(web.Request) {
	return func(r web.Request) { r.BodyString(body) }
}

// BodyForm uses the form values to send a `application/x-www-form-urlencoded` encoded request
func BodyForm(form urlpkg.Values) func(web.Request) {
	return func(r web.Request) { r.BodyForm(form) }
}

// BodyFile read the file from the given path, and uses it as the request body
func BodyFile(path string) func(web.Request) {
	return func(r web.Request) { r.BodyFile(path) }
}

// BodyEncoding sets the body encoding to use
func BodyEncoding(encoding web.RequestEncoding) func(web.Request) {
	return func(r web.Request) { r.BodyEncoding(encoding) }
}

// BodyEncode encodes v, and uses that as the request body.
// If encoding is not specified here, you must set it using BodyEncoding(...)
// If more than one encoding is specified only the first encoding is used.
func BodyEncode(v interface{}, encoding ...web.RequestEncoding) func(web.Request) {
	return func(r web.Request) { r.BodyEncode(v, encoding...) }
}

// Ensure calls the RequestProcessors right after a
// *http.Request instance is created but before it is used
func Ensure(h ...web.RequestProcessor) func(web.Request) {
	return func(r web.Request) { r.Ensure(h...) }
}

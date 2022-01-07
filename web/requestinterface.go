package web

//go:generate go run ..\web-gen\main.go package=rq

import (
	"io"
	"net/http"
	urlpkg "net/url"
)

var _ Request = newRequest()

type (
	// GetBody is the function that retrieves a reader for the request body.
	// GetBody can be called multiple times if for example the request is redirected.
	// See also http.Request.GetBody
	GetBody func() (io.ReadCloser, error)

	// Called once on every http.Request after it has been constructed by the builder.
	RequestProcessor func(*http.Request) error

	Request interface {
		// Method sets the request httmp method
		Method(method string) Request

		// URL

		// BaseUrl sets the request base url
		BaseURL(baseURL string) Request

		// Scheme sets the request url scheme
		Scheme(scheme string) Request

		// Host sets the request url host portion
		Host(host string) Request

		// Hostf sets the request url host portion to the formatted host string
		Hostf(host string, args ...interface{}) Request

		// HostAndPort sets the request url host portion to the provided host and port value
		HostAndPort(host string, port int) Request

		// Path appends a path segment to the request base url, if it starts with a /, the path is replaced
		// Path segments are appended using the following code:
		// for _, p := range pathSegments {
		//   baseUrl.Path = baseUrl.ResolveReference(&url.URL{Path: p}).Path
		// }
		Path(path string) Request

		// Pathf formats the path with the given args, and calls Path(string)
		Pathf(path string, args ...interface{}) Request

		// Param sets the values for the given query parameter, overwriting existing values in the url.
		Param(name string, values ...string) Request

		// Fragment sets the url fragment value
		Fragment(fragment string) Request

		// Headers

		// ContentType sets the Content-Type header
		ContentType(contentType string) Request

		// Accept sets the Accept header
		Accept(contentType string) Request

		// AcceptRange sets the Accept header with the given content types, sorted by their quality value
		AcceptRange(contentTypes map[string]float64) Request

		// BasicAuth sets the Authorization header to Basic, with the base64 encoded user and password
		BasicAuth(user, password string) Request

		// BearerAuth sets the Authorization header to Bearer, with the given token.
		BearerAuth(token string) Request

		// Header sets the given header, overwriting existing values
		Header(key string, values ...string) Request

		// Body uses the getBody function to read the request body contents.
		// getBody may be called multiple times in case of request redirection
		Body(getBody GetBody) Request

		// BodyCached uses getBody once, caching the resulting bytes should they be needed again.
		BodyCached(getBody GetBody) Request

		// BodyReader uses body once to read the request body, caching the resulting bytes should they be needed again.
		BodyReader(body io.Reader) Request

		// BodyBytes uses body as the request body, modifying the contents of
		// the slice isn't recommended, as BodyBytes doesn't make a copy
		BodyBytes(body []byte) Request

		// BodyString uses body as the request body
		BodyString(body string) Request

		// BodyForm uses the form values to send a `application/x-www-form-urlencoded` encoded request
		BodyForm(form urlpkg.Values) Request

		// BodyFile read the file from the given path, and uses it as the request body
		BodyFile(path string) Request

		// BodyEncoding sets the body encoding to use
		BodyEncoding(encoding RequestEncoding) Request

		// BodyEncode encodes v, and uses that as the request body.
		// If encoding is not specified here, you must set it using BodyEncoding(...)
		// If more than one encoding is specified only the first encoding is used.
		BodyEncode(v interface{}, encoding ...RequestEncoding) Request

		// Ensure calls the RequestProcessors right after a
		// *http.Request instance is created but before it is used
		Ensure(h ...RequestProcessor) Request

		// apply request configuration
		apply(configure ...func(Request))
	}
)

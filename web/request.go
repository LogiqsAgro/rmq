package web

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"net/http"
	urlpkg "net/url"
	"os"
	"sort"
	"strings"
	"sync"
)

type (
	request struct {
		baseURL  string
		scheme   string
		host     string
		paths    []string
		query    urlpkg.Values
		fragment string

		method  string
		headers http.Header
		getBody GetBody

		encoding   RequestEncoding
		processors []RequestProcessor
	}
)

// NewRequest creates a new request
func newRequest() *request {
	return &request{}
}

func (r *request) Scheme(scheme string) Request {
	r.scheme = scheme
	return r
}

func (r *request) Host(host string) Request {
	r.host = host
	return r
}

func (r *request) Hostf(host string, args ...interface{}) Request {
	return r.Host(fmt.Sprintf(host, args...))
}

func (r *request) HostAndPort(host string, port int) Request {
	return r.Hostf("%s:%d", host, port)
}

// ContentType is a convenience func for Header("Content-Type", contentType)
func (r *request) ContentType(contentType string) Request {
	return r.Header("Content-Type", contentType)
}

// Accept is a convenience func for Header("Accept", contentType)
func (r *request) Accept(contentType string) Request {
	return r.Header("Accept", contentType)
}

// Accept is a convenience func for Header("Accept", contentType)
func (r *request) AcceptRange(contentTypes map[string]float64) Request {
	// sort higher quality items first
	ss := makeSortedRange(contentTypes)

	return r.Header("Accept", ss...)
}

func makeSortedRange(qRange map[string]float64) []string {
	ss := make([]string, 0, len(qRange))
	for k := range qRange {
		ss = append(ss, k)
	}
	sort.Slice(ss, func(i, j int) bool {
		delta := qRange[ss[i]] - qRange[ss[j]]
		if delta == 0 {
			return ss[i] < ss[j]
		}
		return delta > 0
	})

	for i := 0; i < len(ss); i++ {
		q := qRange[ss[i]]
		q = math.Max(0, math.Min(q, 1))
		ss[i] = fmt.Sprintf("%s;q=%1.3g", ss[i], q)
	}
	return ss
}

// BasicAuth sets the Authorization header with the given user and password
func (r *request) BasicAuth(user, password string) Request {
	auth := user + ":" + password
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return r.Header("Authorization", "Basic "+encoded)
}

// BearerAuth sets the Authorization header to a bearer token.
func (r *request) BearerAuth(token string) Request {
	return r.Header("Authorization", "Bearer "+token)
}

// Header sets a request header, overwrites previous values set for the header
func (r *request) Header(key string, values ...string) Request {
	if r.headers == nil {
		r.headers = make(http.Header)
	}
	key = http.CanonicalHeaderKey(key)
	r.headers[key] = values
	return r
}

// Method sets the http request method
func (r *request) Method(method string) Request {
	r.method = method
	return r
}

// BaseURL sets the http request base url, if baseURL does not end with a '/' it is appended for you
func (r *request) BaseURL(baseURL string) Request {
	if len(baseURL) > 0 && !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}
	r.baseURL = baseURL
	return r
}

// Path appends a new path segment to the base url
func (r *request) Path(path string) Request {
	r.paths = append(r.paths, path)
	return r
}

// Pathf formats path template with the given args, and appends the new path segment to the base url
func (r *request) Pathf(path string, args ...interface{}) Request {
	r.Path(fmt.Sprintf(path, args...))
	return r
}

// Fragment sets the url fragment
func (r *request) Fragment(fragment string) Request {
	r.fragment = fragment
	return r
}

// Param sets the query parameter, overwriting existing values
func (r *request) Param(name string, values ...string) Request {
	if r.query == nil {
		r.query = make(urlpkg.Values)
	}
	r.query[name] = values
	return r
}

// Validate adds request handlers that are executed after a request is created and before it is sent.
// If any handler returns an error, processing is stopped and the request is not sent.
func (r *request) Ensure(h ...RequestProcessor) Request {
	r.processors = append(r.processors, h...)
	return r
}

func (r *request) URL() (*urlpkg.URL, error) {
	url, err := urlpkg.Parse(r.baseURL)
	if err != nil {
		return nil, fmt.Errorf("could not initialize with base URL %q: %w", r.baseURL, err)
	}

	if url.Scheme == "" {
		url.Scheme = "https"
	}

	if r.scheme != "" {
		url.Scheme = r.scheme
	}

	if r.host != "" {
		url.Host = r.host
	}

	for _, p := range r.paths {
		url.Path = url.ResolveReference(&urlpkg.URL{Path: p}).Path
	}

	if len(r.query) > 0 {
		q := url.Query()
		for k, vv := range r.query {
			q[k] = vv
		}
		url.RawQuery = q.Encode()
	}

	if r.fragment != "" {
		url.Fragment = r.fragment
	}

	return url, nil
}

// Clone creates a deep copy of the request.
func (r *request) Clone() *request {
	clone := *r
	clone.headers = r.cloneHeaders()
	clone.paths = r.clonePaths()
	clone.query = r.cloneQuery()
	clone.processors = r.cloneProcessors()
	return &clone
}

// Body sets the source of the request body
// use nil to send an empty body.
func (r *request) Body(getBody GetBody) Request {
	r.getBody = getBody
	return r
}

// BodyCached wraps the GetBody in a caching GetBody.
// Be careful not to cache multi gigabyte files.
// Use when you expect to do the request multiple times
// and the request body is small, but expensive to generate.
func (r *request) BodyCached(getBody GetBody) Request {
	return r.Body(cached(getBody))
}

// cached wraps getBody in a caching function
func cached(getBody GetBody) GetBody {
	cached := &struct {
		lock  sync.Mutex
		error error
		data  *[]byte
	}{}

	return func() (io.ReadCloser, error) {
		cached.lock.Lock()
		defer cached.lock.Unlock()

		if cached.error != nil {
			return nil, cached.error
		} else if cached.data != nil {
			newBuf := bytes.NewBuffer(*cached.data)
			return io.NopCloser(newBuf), nil
		} else {
			body, err := getBody()
			if err != nil {
				cached.error = err
				return nil, err
			}
			defer body.Close()
			data, err := io.ReadAll(body)
			if err != nil && err != io.EOF {
				cached.error = err
				return nil, err
			}
			cached.data = &data
			cached.error = nil
			newBuf := bytes.NewBuffer(data)
			return io.NopCloser(newBuf), nil
		}
	}
}

// BodyReader reads the request body from body
func (r *request) BodyReader(body io.Reader) Request {
	return r.Body(func() (io.ReadCloser, error) {
		if rc, ok := body.(io.ReadCloser); ok {
			return rc, nil
		}
		return io.NopCloser(body), nil
	})
}

// BodyBytes sets body as the source of the request body
func (r *request) BodyBytes(body []byte) Request {
	r.getBody = func() (io.ReadCloser, error) {
		buf := bytes.NewBuffer(body)
		return io.NopCloser(buf), nil
	}
	return r
}

// BodyString sets body as the source of the request body
func (r *request) BodyString(body string) Request {
	r.getBody = func() (io.ReadCloser, error) {
		buf := bytes.NewBufferString(body)
		return io.NopCloser(buf), nil
	}
	return r
}

// BodyForm sets the Content-Type header to "application/x-www-form-urlencoded".
// and sets the body to the encoded form
func (r *request) BodyForm(form urlpkg.Values) Request {
	return r.
		ContentType("application/x-www-form-urlencoded").
		Body(func() (io.ReadCloser, error) {
			body := form.Encode()
			buf := bytes.NewBufferString(body)
			rc := io.NopCloser(buf)
			return rc, nil
		})
}

func (r *request) BodyFile(path string) Request {
	return r.Body(func() (io.ReadCloser, error) {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("could not open file for reading '%s': %w", path, err)
		}
		return f, nil
	})
}

// BodyEncoding sets the body encoding, and the Content-Type header
func (r *request) BodyEncoding(encoding RequestEncoding) Request {
	r.encoding = encoding
	if encoding != nil {
		r.ContentType(encoding.ContentType())
	}
	return r
}

// BodyEncode encodes v and sets it as the request body
func (r *request) BodyEncode(v interface{}, encoding ...RequestEncoding) Request {
	if len(encoding) > 0 {
		r.BodyEncoding(encoding[0])
	}

	return r.Body(func() (io.ReadCloser, error) {
		if r.encoding == nil {
			return nil, newRequestErrorf(nil, "no request encoder configured")
		}

		buf := &bytes.Buffer{}
		encoder := r.encoding.NewEncoder(buf)
		if err := encoder.Encode(v); err != nil {
			return nil, newRequestErrorf(err, "request body encoding failed: %s", err.Error())
		}
		return io.NopCloser(buf), nil
	})
}

// Build creates a new *http.Request from the configured values
func (r *request) Build(ctx context.Context) (*http.Request, error) {
	url, err := r.URL()
	if err != nil {
		return nil, err
	}

	body, err := r.getBodyReader()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, r.method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.GetBody = r.getBodyReader

	r.applyHeaders(req)

	err = r.process(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// getBodyReader returns (http.NoBody, nil) if request.getBody is nil, else it returns the
// result of request.getBody()
func (r *request) getBodyReader() (io.ReadCloser, error) {
	if r.getBody == nil {
		return http.NoBody, nil
	}

	return r.getBody()
}

func (r *request) applyHeaders(req *http.Request) {
	for k, vv := range r.headers {
		delete(req.Header, k)
		for i := range vv {
			req.Header.Add(k, vv[i])
		}
	}
}

func (r *request) process(req *http.Request) error {
	hh := r.processors
	for i := 0; i < len(hh); i++ {
		if h := hh[i]; h != nil {
			if err := h(req); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *request) cloneHeaders() http.Header {
	return r.headers.Clone()
}

func (r *request) clonePaths() []string {
	pp := make([]string, len(r.paths))
	copy(pp, r.paths)
	return pp
}
func (r *request) cloneQuery() urlpkg.Values {
	q := r.query
	if q == nil {
		return nil
	}

	numValues := 0
	for _, vv := range q {
		numValues += len(vv)
	}

	// shared backing array for values
	newValues := make([]string, numValues)

	q2 := make(urlpkg.Values, len(q))
	for key, oldValues := range q {
		n := copy(newValues, oldValues)
		q2[key] = newValues[:n:n]
		newValues = newValues[n:]
	}
	return q2
}

func (r *request) cloneProcessors() []RequestProcessor {
	hh := make([]RequestProcessor, len(r.processors))
	copy(hh, r.processors)
	return hh
}

func (r *request) apply(configure ...func(Request)) {
	for i := 0; i < len(configure); i++ {
		if configure[i] == nil {
			continue
		}
		configure[i](r)
	}
}

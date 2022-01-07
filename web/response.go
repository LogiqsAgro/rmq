package web

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	// DefaultMaxResponseSize is set to 1MB ( 1 << 20 bytes )
	// This limit is set low intentionally, to make unwise decisions explicit.
	// ( like reading a response without limits from a random server on the internet and writing it to your disk. )
	DefaultMaxResponseSize int64 = 1 << 20
)

type (
	response struct {
		builder *builder
		maxSize int64

		encodings   responseEncodings
		processors  []ResponseProcessor
		bodyHandler ResponseProcessor
	}
)

var _ Response = newResponse()

func newResponse() *response {
	return &response{
		maxSize: DefaultMaxResponseSize,
	}
}

func (rsp *response) Clone() *response {
	clone := *rsp
	clone.processors = rsp.cloneHandlers()
	return &clone
}

// Ensure adds response handlers that are executed after a response has been received, but before the response body is read.
// If any handler returns an error, processing is stopped and the response body is discarded.
func (r *response) Ensure(rv ...ResponseProcessor) Response {
	r.processors = append(r.processors, rv...)
	return r
}

// MaxSize sets the maximum number of bytes that is processed by the response body handler.
// More bytes might be read to detect if the maximum size is exceeded, but are never passed on to the handler.
// Set the size to -1 to disable the response size limit.
// Setting the maximum size to zero is valid, if you expect the response body to be empty.
func (r *response) MaxSize(size int64) Response {
	if size < -1 {
		r.maxSize = -1
	} else {
		r.maxSize = size
	}

	return r
}

// ResponseBody sets the handler of the response body,
// If it is nil, the response body will be discarded.
func (r *response) Body(rh func(context.Context, io.Reader) error) Response {
	r.bodyHandler = func(res *http.Response) error {
		ctx := context.Background()
		if res.Request != nil {
			ctx = res.Request.Context()
		}
		// defer res.Body.Close() is done in builder.Do(...)
		body := r.getResponseBodyReader(res)
		return rh(ctx, body)
	}
	return r
}

func (r *response) BodyFile(path string) Response {
	return r.Body(func(ctx context.Context, rdr io.Reader) error {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("could not open file for writing '%s': %w", path, err)
		}

		_, err = io.Copy(file, rdr)
		return err
	})
}

func (r *response) BodyCopyTo(w io.Writer) Response {
	r.bodyHandler = func(res *http.Response) error {
		body := r.getResponseBodyReader(res)
		_, err := io.Copy(w, body)
		return err
	}
	return r
}

// BodyEncoding sets the body encoding, and the Content-Type header
func (r *response) BodyEncoding(encodings ...ResponseEncoding) Response {
	r.encodings = encodings

	switch len(encodings) {
	case 0:
		break
	case 1:
		r.builder.request.Accept(encodings[0].ContentType())
	default:
		qv := r.encodings.QualityValues()
		r.builder.request.AcceptRange(qv)
	}
	return r
}

func (r *response) BodyDecode(v interface{}, encoding ...ResponseEncoding) Response {
	r.BodyEncoding(encoding...)
	r.bodyHandler = func(res *http.Response) error {
		contentType := res.Header.Get("Content-Type")
		if i := strings.Index(contentType, ";"); i >= 0 {
			contentType = contentType[:i]
		}

		for i := range r.encodings {
			encoding := r.encodings[i]
			if len(contentType) == 0 || strings.EqualFold(encoding.ContentType(), contentType) {
				body := r.getResponseBodyReader(res)
				decoder := encoding.NewDecoder(body)
				err := decoder.Decode(v)
				if err != nil {
					return newResponseErrorf(err, "decoding the response failed: %s", err.Error())
				}
				return nil
			}

		}

		return newResponseErrorf(nil, "no response encoding configured for content type %s", contentType)
	}

	return r
}

func (r *response) apply(configure ...func(Response)) {
	for i := 0; i < len(configure); i++ {
		if configure[i] == nil {
			continue
		}
		configure[i](r)
	}
}

func (r *response) process(rsp *http.Response) error {
	hh := r.processors
	for i := 0; i < len(hh); i++ {
		if h := hh[i]; h != nil {
			if err := h(rsp); err != nil {
				return newResponseErrorf(err, "")
			}
		}
	}
	return nil
}

func (r *response) getResponseBodyReader(hr *http.Response) io.Reader {
	if r.maxSize >= 0 {
		return io.LimitReader(hr.Body, r.maxSize)
	}
	return hr.Body
}
func (r *response) invokeBodyHandler(hr *http.Response) error {
	h := r.bodyHandler
	body := r.getResponseBodyReader(hr)
	if h == nil {
		_, err := io.Copy(io.Discard, body)
		return err
	} else {
		return h(hr)
	}
}

func (r *response) cloneHandlers() []ResponseProcessor {
	hh := make([]ResponseProcessor, len(r.processors))
	copy(hh, r.processors)
	return hh
}

type responseEncodings []ResponseEncoding

func (r responseEncodings) QualityValues() map[string]float64 {
	qv := make(map[string]float64)
	if len(r) == 0 {
		// do nothing
	} else if len(r) == 1 {
		qv[r[0].ContentType()] = 1.0
	} else {
		seen := make(map[string]bool)
		dec := 1.0 / float64(len(r)+1)
		q := 1.0
		for i := 0; i < len(r); i++ {
			ct := r[i].ContentType()

			// ignore duplicate content types
			lct := strings.ToLower(ct)
			if seen[lct] {
				continue
			} else {
				seen[lct] = true
			}

			qv[ct] = q
			q -= dec
		}
	}

	return qv
}

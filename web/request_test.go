package web

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_makeSortedRange(t *testing.T) {
	ss := makeSortedRange(map[string]float64{
		"quarter": 0.2554,
		"chalf":   0.5,
		"ahalf":   0.5,
		"bhalf":   0.500,
		"whole":   1,
	})

	t.Logf("sorted range: %s", strings.Join(ss, ", "))

	if !strings.HasPrefix(ss[0], "whole") {
		t.Errorf("Expected 'whole' as the first value")
	}

	if !strings.HasPrefix(ss[2], "bhalf") {
		t.Errorf("Expected 'bhalf' as the middle value")
	}

	if !strings.HasSuffix(ss[2], "0.5") {
		t.Errorf("Expected '0.5' as the middle value's quality")
	}

	if !strings.HasPrefix(ss[4], "quarter") {
		t.Errorf("Expected 'quarter' as the last value")
	}

	if !strings.HasSuffix(ss[4], "0.255") {
		t.Errorf("Expected '0.255' as the last value's quality")
	}
}

func TestRequest_Clone(t *testing.T) {
	a := New().
		Request(func(r Request) {
			r.
				Path("a").
				ContentType("application/json").
				Param("p", "p").
				Ensure(func(*http.Request) error {
					return nil
				})
		}).
		Response(func(r Response) {
			r.Ensure(func(*http.Response) error {
				return nil
			})
		}).(*builder)

	b := a.clone()

	a.request.paths[0] = ""
	if b.request.paths[0] != "a" {
		t.Errorf("paths were not cloned properly")
	}

	a.request.headers["Content-Type"] = []string{}
	if b.request.headers["Content-Type"][0] != "application/json" {
		t.Errorf("headers were not cloned properly")
	}

	a.request.query["p"] = []string{}
	if b.request.query["p"][0] != "p" {
		t.Errorf("query parameters were not cloned properly")
	}

	a.request.processors[0] = nil
	if b.request.processors[0] == nil {
		t.Errorf("requestHandlers were not cloned properly")
	}

	a.response.processors[0] = nil
	if b.response.processors[0] == nil {
		t.Errorf("responseHandlers were not cloned properly")
	}

}

func Test_request_URL(t *testing.T) {
	url, err := Get("http://localhost/").URL()
	if err != nil {
		t.Errorf("error generating url %v", err)
	}

	t.Logf("Url: %s", url.String())
}

func Test_what_happens_if_you_handle_request_body_in_OnRequest_URL(t *testing.T) {
	server := newTestServer(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(200)
		rw.Header().Add("Content-Type", "application/json")
		rw.Write([]byte("{\"text\": \"Hello, world!\"}"))
	})
	defer server.Close()

	err := Get(server.URL).
		Response(func(r Response) {
			r.Ensure(func(r *http.Response) error {
				buf := &bytes.Buffer{}
				_, err := io.Copy(buf, r.Body)
				t.Logf("Response %d: %s", r.StatusCode, buf.String())
				return err
			})
		}).
		Invoke(context.Background())

	if err != nil {
		t.Errorf("error reading response url %v", err)
	}
}

func newTestServer(h http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(h)
}

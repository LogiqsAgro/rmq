package examples

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/LogiqsAgro/rmq/web"
	"github.com/LogiqsAgro/rmq/web/rq"
	"github.com/LogiqsAgro/rmq/web/rs"
)

// The web builder exposes a fluent interface to configure a http.Request,
// send it and verify responses in one method chain.
//
// Below is a full example of using just the web package. To really make it
// nice to use, the rq and rs packages define helper methods that make most the
// boilerplate func declarations dissapear.
func Example() {

	// Set up the request context, with some context value and a timeout, because we can.
	ctx := context.WithValue(context.Background(), traceParentKey, traceParent)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Set up a example request handler and server
	example := &exampleHandler{
		statusCode:  200,
		contentType: "text/plain",
		response:    "Hello to you too!",
	}
	server := httptest.NewServer(example)

	// Get configures a http GET request, there are methods for all common HTTP verbs
	err := web.
		Get(server.URL,
			// there are utility methods in the rq package
			// to help you easily set up headers...
			rq.ContentType("text/plain"),
			rq.Accept("text/plain"),
			// and the body...
			rq.BodyString("Hello!"),
			// and methods to ensure the http.Request meets expectations
			rq.Ensure(traceIdIsSet),
		).
		// Response is used to set up response handling
		Response(
			// there are utility methods in the rs package
			// to help you easily set up response expectations...
			rs.EnsureStatus(200),
			// and to process the body...
			rs.MaxSize(1024),
		).
		// Invoke creates and sends the request,
		// verifies the response, and processes the response body,
		Invoke(ctx, printBody)

	// if the request the returned error
	if err != nil {
		panic(err)
	}

	// Output:
	// request  body: Hello!
	// response body: Hello to you too!
}

func printBody(ctx context.Context, body io.Reader) error {
	os.Stdout.WriteString("response body: ")
	s, err := readBodyAsString(body)
	if err != nil {
		return err
	}
	_, err = os.Stdout.WriteString(s + "\n")
	if err != nil {
		return err
	}
	return err
}

// ensures the http.Request.Context() has a traceId value set.
func traceIdIsSet(r *http.Request) error {
	value, ok := r.Context().Value(traceParentKey).(string)
	if ok && len(value) > 0 {
		r.Header.Set(string(traceParentKey), value)
		return nil
	} else {
		return fmt.Errorf("Expected a '%s' value to be present in the request context", traceParentKey)
	}
}

// readBodyAsString drains the reader, trimming all trailing newlines
func readBodyAsString(body io.Reader) (string, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return "nil", err
	}
	s := string(data)
	s = strings.TrimRight(s, "\n")
	return s, nil
}

type exampleHandler struct {
	statusCode  int
	contentType string
	response    string
}

func (h *exampleHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := readBodyAsString(r.Body)
	if err == nil {
		os.Stdout.WriteString("request  body: ")
		os.Stdout.WriteString(body)
		os.Stdout.WriteString("\n")

		rw.Header().Add("Content-Type", h.contentType)

		rw.WriteHeader(h.statusCode)
		rw.Write([]byte(h.response))
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
	}
}

type contextKey string

const (
	traceParentKey contextKey = "traceparent"

	// see https://www.w3.org/TR/trace-context/#trace-context-http-headers-format
	traceParent string = "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
)

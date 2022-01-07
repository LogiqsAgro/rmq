package examples

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"os"
	"time"

	"github.com/LogiqsAgro/rmq/web"
	"github.com/LogiqsAgro/rmq/web/rs"
)

// The web builder exposes a fluent interface to configure a http.Request,
// send it and verify responses in one method chain.
//
// Below is a full example of using just the web package. To really make it
// nice to use, the rq and rs packages define helper methods that make most the
// boilerplate func declarations dissapear.
func ExampleBuilder_UseJSON() {
	//
	// Set up the request context, with some context value and a timeout, because we can.
	ctx := context.WithValue(context.Background(), traceParentKey, traceParent)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Set up a example request handler and server
	example := &exampleHandler{
		statusCode:  200,
		contentType: "application/json",
		response:    `{"name": "jane", "age": 36 }`,
	}
	server := httptest.NewServer(example)

	query := &struct{ ByName string }{"jane"}

	value := &struct {
		Name string
		Age  int
	}{}

	// Get and decode the json response
	err := web.
		Get(server.URL).
		Response(rs.EnsureStatus(200)).
		UseJSON(query, value).
		Invoke(ctx)

	// if the request the returned error, panic
	if err != nil {
		panic(err)
	}

	os.Stdout.WriteString("encoded value: ")
	err = json.NewEncoder(os.Stdout).Encode(value)
	// if the encoding failed, panic
	if err != nil {
		panic(err)
	}

	// Output:
	// request  body: {"ByName":"jane"}
	// encoded value: {"Name":"jane","Age":36}
}

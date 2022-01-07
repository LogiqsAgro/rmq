package web

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type (
	// Builder is the main entrypoint to the web package
	Builder interface {

		// Client sets the *http.Client to use for the request
		Client(client *http.Client) Builder

		// Transport is shorthand for Client(&http.Client{Transport: transport})
		Transport(transport *http.Transport) Builder

		// Codec configures this request to use codec.ContentType() as the Content-Type and Accept headers.
		Codec(codec Codec) Builder

		// UseCodec configures this request to use codec.ContentType() as the Content-Type and Accept headers.
		// req is encoded and used as the request body, the response is decoded into rsp.
		UseCodec(req, rsp interface{}, codec Codec) Builder

		// UseJSON configures this request to use "application/json" as the Content-Type and Accept headers.
		// the encoding/json package is used to encode req and decode rsp
		UseJSON(req, rsp interface{}) Builder

		// Request configures the request building
		Request(cfg ...func(Request)) Builder

		// Response configures the response handling
		Response(cfg ...func(Response)) Builder

		// URL returns the url that was configured, or an error if it wasn't confgured properly
		URL() (*url.URL, error)

		// HttpRequest returns the configured and processed *http.Request instance.
		HttpRequest(ctx context.Context) (*http.Request, error)

		// Do invokes the request using the configured client
		// and invokes the configured response processors and body handlers
		Do(req *http.Request) error

		// Invoke call HttpRequest(ctx), and calls Do(...) on the result.
		// If a handler is specified, that is called with the request context and the response body as parameters.
		// It is an error to specify more than one handler
		Invoke(ctx context.Context, handler ...func(context.Context, io.Reader) error) error

		// Clone creates a copy of this builder.
		// This is useful if you want to use this library to do a lot of requests
		// that only vary by path or query parameters.
		Clone() Builder
	}

	// Encoder encodes values
	Encoder interface {
		// Encode writes the encoded value of v to its output writer
		Encode(v interface{}) error
	}

	// Decoder decodes values
	Decoder interface {
		// Decode reads the next encoded value from its input
		// reader and stores it in the value pointed to by v.
		Decode(v interface{}) error
	}

	// RequestEncoding creates request body encoders,
	// and specifies the content type of the encoded value
	RequestEncoding interface {
		ContentType() string
		NewEncoder(w io.Writer) Encoder
	}

	// ResponseEncoding creates response body decoders,
	// and specifies the content type of the encoded value
	ResponseEncoding interface {
		ContentType() string
		NewDecoder(r io.Reader) Decoder
	}

	// Codec creates request body encoders and response body decoders
	// and specifies the content type of the encoded values
	Codec interface {
		RequestEncoding
		ResponseEncoding
	}

	builder struct {
		client   *http.Client
		request  *request
		response *response
	}
)

func New() Builder {
	return newBuilder()
}

func newBuilder() *builder {
	b := &builder{
		request:  newRequest(),
		response: newResponse(),
	}
	b.response.builder = b
	return b
}

func (b *builder) Client(client *http.Client) Builder {
	b.client = client
	return b
}

func (b *builder) Transport(transport *http.Transport) Builder {
	return b.Client(&http.Client{Transport: transport})
}

func (b *builder) Codec(codec Codec) Builder {
	b.request.BodyEncoding(codec)
	b.response.BodyEncoding(codec)
	return b
}

func (b *builder) Clone() Builder {
	return b.clone()
}

func (b *builder) clone() *builder {
	clone := &builder{
		request:  b.request.Clone(),
		response: b.response.Clone(),
	}
	clone.response.builder = clone
	return clone
}
func (b *builder) URL() (*url.URL, error) {
	return b.request.URL()
}

func (b *builder) HttpRequest(ctx context.Context) (*http.Request, error) {
	return b.request.Build(ctx)
}

func (b *builder) Request(cfg ...func(Request)) Builder {
	b.request.apply(cfg...)
	return b
}

func (b *builder) Response(cfg ...func(Response)) Builder {
	b.response.apply(cfg...)
	return b
}

func (b *builder) Invoke(ctx context.Context, handler ...func(context.Context, io.Reader) error) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := b.request.Build(ctx)
	if err != nil {
		return err
	}

	switch len(handler) {
	case 0:
		break
	case 1:
		b = b.clone()
		b.response.Body(handler[0])
	default:
		return fmt.Errorf("only one response handler allowed")
	}

	return b.Do(req)
}

func (b *builder) Do(req *http.Request) error {
	client := b.client
	if client == nil {
		client = http.DefaultClient
	}

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer rsp.Body.Close()

	if err := b.response.process(rsp); err != nil {
		return err
	}

	return b.response.invokeBodyHandler(rsp)
}

func (b *builder) UseJSON(req, rsp interface{}) Builder {
	return b.UseCodec(req, rsp, &jsonCodec{})
}

func (b *builder) UseCodec(req interface{}, rsp interface{}, codec Codec) Builder {
	b.request.BodyEncode(req, codec)
	b.response.BodyDecode(rsp, codec)
	return b
}

type jsonCodec struct{}

func (c *jsonCodec) ContentType() string            { return "application/json" }
func (c *jsonCodec) NewEncoder(w io.Writer) Encoder { return json.NewEncoder(w) }
func (c *jsonCodec) NewDecoder(w io.Reader) Decoder { return json.NewDecoder(w) }

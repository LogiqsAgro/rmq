package web

import "net/http"

// Connect is a convenience func for creating a CONNECT request
func Connect(baseURL string, configure ...func(Request)) Builder {
	b := newBuilder()
	b.request.
		Method(http.MethodConnect).
		BaseURL(baseURL).
		apply(configure...)

	return b
}

// Delete is a convenience func for creating a DELETE request
func Delete(baseURL string, configure ...func(Request)) Builder {
	b := newBuilder()
	b.request.
		Method(http.MethodDelete).
		BaseURL(baseURL).
		apply(configure...)

	return b
}

// Get is a convenience func for creating a GET request
func Get(baseURL string, configure ...func(Request)) Builder {
	b := newBuilder()
	b.request.
		Method(http.MethodGet).
		BaseURL(baseURL).
		apply(configure...)

	return b
}

// Head is a convenience func for creating a HEAD request
func Head(baseURL string, configure ...func(Request)) Builder {
	b := newBuilder()
	b.request.
		Method(http.MethodHead).
		BaseURL(baseURL).
		apply(configure...)

	return b
}

// Options is a convenience func for creating an OPTIONS request
func Options(baseURL string, configure ...func(Request)) Builder {
	b := newBuilder()
	b.request.
		Method(http.MethodOptions).
		BaseURL(baseURL).
		apply(configure...)

	return b
}

// Patch is a convenience func for creating a PATCH request
func Patch(baseURL string, configure ...func(Request)) Builder {
	b := newBuilder()
	b.request.
		Method(http.MethodPatch).
		BaseURL(baseURL).
		apply(configure...)

	return b
}

// Post is a convenience func for creating a POST request
func Post(baseURL string, configure ...func(Request)) Builder {
	b := newBuilder()
	b.request.
		Method(http.MethodPost).
		BaseURL(baseURL).
		apply(configure...)

	return b
}

// Put is a convenience func for creating a PUT request
func Put(baseURL string, configure ...func(Request)) Builder {
	b := newBuilder()
	b.request.
		Method(http.MethodPut).
		BaseURL(baseURL).
		apply(configure...)

	return b
}

// Trace is a convenience func for creating a TRACE request
func Trace(baseURL string, configure ...func(Request)) Builder {
	b := newBuilder()
	b.request.
		Method(http.MethodTrace).
		BaseURL(baseURL).
		apply(configure...)

	return b
}

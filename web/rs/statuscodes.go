// rs stands for response
package rs

import (
	"fmt"
	"net/http"

	"github.com/LogiqsAgro/rmq/web"
)

//  EnsureStatus ensures that the response has one of the expected status codes
func EnsureStatus(expectedStatusCodes ...int) func(web.Response) {
	return func(r web.Response) {
		r.Ensure(func(r *http.Response) error {
			statusCode := r.StatusCode
			for i := 0; i < len(expectedStatusCodes); i++ {
				if statusCode == expectedStatusCodes[i] {
					return nil
				}
			}
			return fmt.Errorf("expected one of %v status codes but got %v", expectedStatusCodes, r.StatusCode)
		})
	}
}

//  EnsureStatusInformational ensures that the response has a status code in the range 100-199
func EnsureStatusInformational() func(web.Response) {
	return func(r web.Response) {
		r.Ensure(func(r *http.Response) error {
			if http.StatusContinue <= r.StatusCode && r.StatusCode < http.StatusOK {
				return fmt.Errorf("expected 100-199 status, got %s", r.Status)
			}
			return nil
		})
	}
}

// EnsureStatusOK ensures that the response has a status code 200 OK
func EnsureStatusOK() func(web.Response) {
	return func(r web.Response) {
		r.Ensure(func(r *http.Response) error {
			if r.StatusCode != http.StatusOK {
				return fmt.Errorf("expected 200 OK status, got %s", r.Status)
			}
			return nil
		})
	}
}

// EnsureStatusCreated ensures that the response has a status code 201 Created
func EnsureStatusCreated() func(web.Response) {
	return func(r web.Response) {
		r.Ensure(func(r *http.Response) error {
			if r.StatusCode != http.StatusCreated {
				return fmt.Errorf("expected 201 Created status, got %s", r.Status)
			}
			return nil
		})
	}
}

//  EnsureStatusAccepted ensures that the response has a status code 202 Accepted
func EnsureStatusAccepted() func(web.Response) {
	return func(r web.Response) {
		r.Ensure(func(rsp *http.Response) error {
			if rsp.StatusCode != http.StatusAccepted {
				return fmt.Errorf("expected 202 Accepted status, got %s", rsp.Status)
			}
			return nil
		})
	}
}

//  EnsureStatusSuccess ensures that the response has a status code in the range 200-299
func EnsureStatusSuccess() func(web.Response) {
	return func(r web.Response) {
		r.Ensure(func(rsp *http.Response) error {
			if http.StatusOK <= rsp.StatusCode && rsp.StatusCode < http.StatusMultipleChoices {
				return fmt.Errorf("expected 200-299 status, got %s", rsp.Status)
			}
			return nil
		})
	}
}

// EnsureStatusRedirect ensures that the response has a status code in the range 300-399
func EnsureStatusRedirect() func(web.Response) {
	return func(r web.Response) {
		r.Ensure(func(rsp *http.Response) error {
			if http.StatusMultipleChoices <= rsp.StatusCode && rsp.StatusCode < http.StatusBadRequest {
				return fmt.Errorf("expected 300-399 status, got %s", rsp.Status)
			}
			return nil
		})
	}
}

// EnsureStatusClientError ensures that the response has a status code in the range 400-499
func EnsureStatusClientError() func(web.Response) {
	return func(r web.Response) {
		r.Ensure(func(rsp *http.Response) error {
			if http.StatusBadRequest <= rsp.StatusCode && rsp.StatusCode < http.StatusInternalServerError {
				return fmt.Errorf("expected 300-399 status, got %s", rsp.Status)
			}
			return nil
		})
	}
}

// EnsureStatusServerError ensures that the response has a status code in the range 500-599
func EnsureStatusServerError() func(web.Response) {
	return func(r web.Response) {
		r.Ensure(func(rsp *http.Response) error {
			if http.StatusInternalServerError <= rsp.StatusCode && rsp.StatusCode < 600 {
				return fmt.Errorf("expected 300-399 status, got %s", rsp.Status)
			}
			return nil
		})
	}
}

/*
Copyright © 2021 Remco Schoeman <remco.schoeman@logiqs.nl>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Print writes the response to stdout, and returns any errors
func Print(resp *http.Response, err error) error {
	if resp == nil {
		return err
	}

	traceRequest(resp.Request)
	traceResponse(resp, err)

	if Config.IndentJson {
		body := &bytes.Buffer{}
		if _, err := io.Copy(body, resp.Body); err != nil {
			return err
		}

		formatted := &bytes.Buffer{}
		if err := json.Indent(formatted, body.Bytes(), "", "\t"); err != nil {
			return err
		}

		if _, err = io.Copy(os.Stdout, formatted); err != nil {
			return err
		}
	} else {
		if io.Copy(os.Stdout, resp.Body); err != nil {
			return err
		}
	}
	os.Stdout.WriteString("\n")
	if err == nil && !isSuccess(resp.StatusCode) {
		return fmt.Errorf("request failed: %s ( url: %s )", resp.Status, resp.Request.URL.Redacted())
	}

	return err
}

func isSuccess(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}

func traceRequest(r *http.Request) {
	if !Config.Debug {
		return
	}
	fmt.Println(">==>==>==>")
	fmt.Println(r.Method, r.URL)
	for k, v := range r.Header {
		fmt.Println(k, ": ", v)
	}
	fmt.Println(">==>==>==>")
}

func traceResponse(r *http.Response, err error) {
	if !Config.Debug {
		return
	}
	fmt.Println("<==<==<==<")
	fmt.Println(r.Status)
	for k, v := range r.Header {
		fmt.Println(k, ": ", v)
	}
	fmt.Println("<==<==<==<")
}

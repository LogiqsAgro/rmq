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
	"net/http"
	"os"
)

// Print writes the json to stdout, or the error to stderr
func Print(json []byte, err error) error {
	if err == nil {
		if Config.IndentJson {
			if json, err = indentJson(json); err != nil {
				return err
			}
		}
		_, err := os.Stdout.Write(json)
		if err != nil {
			return err
		}
	} else {
		_, err := os.Stderr.WriteString("ERROR: " + err.Error())
		if err != nil {
			return err
		}
	}
	return nil
}

func indentJson(data []byte) ([]byte, error) {
	dst := &bytes.Buffer{}
	err := json.Indent(dst, data, "", "  ")
	if err != nil {
		return nil, err
	}
	return dst.Bytes(), nil
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

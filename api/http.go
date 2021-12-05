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
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	client *http.Client = &http.Client{}
)

// Get invokes the Management api with the relative path, using the GET http verb
func Get(pathAndQuery string) (*http.Response, error) {
	return call(http.MethodGet, pathAndQuery, nil)
}

// Post invokes the Management api with the relative path, using the POST http verb and the given body
func Post(pathAndQuery string, body interface{}) (*http.Response, error) {
	return call(http.MethodPost, pathAndQuery, body)
}

// Put invokes the Management api with the relative path, using the PUT http verb and the given body
func Put(pathAndQuery string, body interface{}) (*http.Response, error) {
	return call(http.MethodPut, pathAndQuery, body)
}

// ReadBody reads and closes the response body in 1 call
func ReadBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// Call invokes the Management api with the relative path, using the given http method: e.g. http://localhost:15672/api/[relPath]
func call(method, pathAndQuery string, body interface{}) (*http.Response, error) {
	pathAndQuery = appendGlobalQueryParameters(pathAndQuery)
	url := Config.url(pathAndQuery)
	return Config.call(method, url, body)
}

func appendGlobalQueryParameters(pathAndQuery string) string {
	q := NewQuery()

	if len(Config.Columns) > 0 {
		q.Add("columns", strings.Join(Config.Columns, ","))
	}

	if len(Config.Sort) > 0 {
		q.Add("sort", Config.Sort)
	}

	q.AddIf(Config.SortReverse, "sort_reverse", "true")

	if q.Empty() {
		return pathAndQuery
	}

	if strings.Contains(pathAndQuery, "?") {
		return pathAndQuery + "&" + q.String()
	}
	return pathAndQuery + q.QueryString()
}

// url returns the absolute url for the given path and query.
func (cfg *cfg) url(pathAndQuery string) string {
	pathAndQuery = strings.TrimLeft(pathAndQuery, "/")
	return fmt.Sprintf("%s://%s:%d/api/%s", cfg.Scheme, cfg.Host, cfg.ApiPort, pathAndQuery)
}

// call invokes the http method on the given url and returns the response.
// it also adds Authorization, Conten-Type and Accept headers.
func (cfg *cfg) call(method, url string, body interface{}) (*http.Response, error) {

	req, err := cfg.createRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	traceRequest(req)
	resp, err := client.Do(req)
	if resp.StatusCode >= http.StatusBadRequest {
		err = fmt.Errorf("RabbitMQ returned error: %s", resp.Status)
	}
	traceResponse(resp, err)
	return resp, err
}

func (cfg *cfg) createRequest(method, url string, body interface{}) (*http.Request, error) {
	req, err := newRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("got error %s", err.Error())
	}

	req.SetBasicAuth(cfg.User, cfg.Password)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	return req, nil
}

func newRequest(method, url string, body interface{}) (*http.Request, error) {
	if body == nil {
		return http.NewRequest(method, url, nil)
	}

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}
	return http.NewRequest(method, url, buf)
}

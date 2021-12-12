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
//go:generate go run ..\api-gen\main.go

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	UnitDays   = "days"
	UnitWeeks  = "weeks"
	UnitMonths = "months"
	UnitYears  = "years"
)

// CertificateExpirationTimeUnits returns the list of time units usable in GetHealthChecksCertificateExpirationJson(within int, unit string)
func CertificateExpirationTimeUnits() []string {
	return []string{
		UnitDays,
		UnitWeeks,
		UnitMonths,
		UnitYears,
	}
}

type (
	builder struct {
		baseUrl  string
		path     string
		method   string
		user     string
		password string

		page      int
		pageSize  int
		usePaging bool

		sortColumn  string
		sortReverse bool
		useSort     bool

		columns    []string
		useColumns bool

		filter         string
		filterUseRegex bool
		useFilter      bool

		memory bool
		binary bool

		body interface{}

		addQueryParams func(Query)
	}

	Builder interface {
		Path(path string) Builder
		Method(method string) Builder
		BaseUrl(baseUrl string) Builder
		BasicAuth(user, password string) Builder
		Body(body interface{}) Builder
		Page(page, size int) Builder
		Columns(columns ...string) Builder
		Sort(column string, reverse bool) Builder
		Filter(filter string, useRegex bool) Builder
		Memory(memory bool) Builder
		Binary(binary bool) Builder
		QueryParameters(addQueryParams func(Query)) Builder

		Url() string
		QueryString() string
		Build() (*http.Request, error)
	}
)

func Request() *builder {
	return &builder{
		baseUrl:  "http://localhost:15672",
		path:     "/",
		method:   http.MethodGet,
		user:     "guest",
		password: "guest",
	}
}

// Do builds a request and uses http.DefaultClient.Do(req *http.Request) to execute it
func Do(b Builder) (*http.Response, error) {
	req, err := b.Build()
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

// Clone returns a copy of the builder with fields baseUrl, path, method, user and password copied
func (b *builder) Clone() *builder {
	return &builder{
		baseUrl:  b.baseUrl,
		path:     b.path,
		method:   b.method,
		user:     b.user,
		password: b.password,
	}
}

// Url returns the full url for the request, including the query string
func (b *builder) Url() string {
	url := strings.TrimRight(b.baseUrl, "/")
	if !strings.HasPrefix(b.path, "/") {
		url += "/"
	}
	url += b.path
	url += b.QueryString()
	return url
}

// QueryString returns the query string including the leading '?'
func (b *builder) QueryString() string {
	// clone the query, so QueryString can be called
	// multiple times on the same builder instance
	q := NewQuery()
	if b.addQueryParams != nil {
		b.addQueryParams(q)
	}

	q.AddIf(b.usePaging, "page", fmt.Sprintf("%d", b.page))
	q.AddIf(b.usePaging, "page_size", fmt.Sprintf("%d", b.pageSize))

	q.AddIf(b.useSort, "sort", b.sortColumn)
	q.AddIf(b.useSort && b.sortReverse, "sort_reverse", "true")

	q.AddIf(b.useFilter, "name", b.filter)
	q.AddIf(b.useFilter && b.filterUseRegex, "use_regex", "true")

	q.AddIf(b.memory, "memory", "true")
	q.AddIf(b.binary, "binary", "true")

	q.AddIf(b.useColumns, "columns", strings.Join(b.columns, ","))
	return q.QueryString()
}

func (b *builder) BasicAuth(user, password string) Builder {
	b.user = user
	b.password = password
	return b
}

func (b *builder) BaseUrl(baseUrl string) Builder {
	b.baseUrl = baseUrl
	return b
}

func (b *builder) Path(path string) Builder {
	b.path = path
	return b
}

func (b *builder) Method(method string) Builder {
	b.method = method
	return b
}

func (b *builder) Body(body interface{}) Builder {
	b.body = body
	return b
}

// QueryParameters can be used to set additional query parameters
// the supplied func will be invoked when Build(), Url() or QueryString() is called
func (b *builder) QueryParameters(addQueryParams func(Query)) Builder {
	b.addQueryParams = addQueryParams
	return b
}

// Build creates a new http.Request
func (b *builder) Build() (*http.Request, error) {
	url := b.Url()
	body := io.Reader(nil)
	if b.body != nil {
		buf := &bytes.Buffer{}
		json.NewEncoder(buf).Encode(b.body)
		body = buf
	}
	req, err := http.NewRequest(b.method, url, body)
	req.SetBasicAuth(b.user, b.password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	if err != nil {
		return nil, err
	}

	return req, nil
}

// Page can be applied to the endpoints that list queues, exchanges, connections and channels
func (b *builder) Page(page, size int) Builder {
	b.page = page
	b.pageSize = size
	b.usePaging = page > 0 && size > 0
	return b
}

// Columns can be applied to most endpoints, and limits the results to the listed fields,
// subfields can be selected by using . to separate fields, e.g. message_stats.publish_details.rate
func (b *builder) Columns(columns ...string) Builder {
	b.columns = columns
	b.useColumns = len(columns) > 0
	return b
}

// Sort applies sorting to the returned list, only 1 sort column is supported.
func (b *builder) Sort(column string, reverse bool) Builder {
	b.sortColumn = column
	b.sortReverse = reverse && len(column) > 0
	b.useSort = len(column) > 0
	return b
}

// Filter can be applied to the endpoints that list queues, exchanges, connections and channels. If useRegex is true, filter is interpreted as a regex.
func (b *builder) Filter(filter string, useRegex bool) Builder {
	b.filter = filter
	b.filterUseRegex = useRegex
	b.useFilter = len(filter) > 0
	return b
}

// Memory adds the memory=true query parameter.
// Only useful when fetching node statistics
func (b *builder) Memory(memory bool) Builder {
	b.memory = memory
	return b
}

// Binary adds the binary=true query parameter.
// Only useful when fetching node statistics
func (b *builder) Binary(binary bool) Builder {
	b.binary = binary
	return b
}

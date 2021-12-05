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
	"fmt"
	"net/url"
	"strings"
)

type (
	pageFilter struct {
		Page     int
		PageSize int
		Name     string
		UseRegex bool
	}
)

func NewPage(page, pageSize int) *pageFilter {
	return NewPageFilter(page, pageSize, "", false)
}

func NewPageFilter(page, pageSize int, name string, useRegex bool) *pageFilter {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 100
	}
	return &pageFilter{
		Page:     page,
		PageSize: pageSize,
		Name:     name,
		UseRegex: useRegex,
	}
}

func newPage() *pageFilter {
	return &pageFilter{
		Page:     1,
		PageSize: 100,
		Name:     "",
		UseRegex: false,
	}
}

// ToQuery returns the Query representing this page
func (p *pageFilter) ToQuery() (Query, bool) {
	q := NewQuery()

	if p.Page > 1 {
		q.Add("page", fmt.Sprintf("%d", p.Page))
	}

	if p.PageSize > 0 && p.PageSize != 100 {
		q.Add("page_size", fmt.Sprintf("%d", p.PageSize))
	}

	if p.Name != "" {
		q.Add("name", p.Name)
		if p.UseRegex {
			q.Add("use_regex", "true")
		}
	}

	return q, !q.Empty()
}

// ToUrlSuffix returns the page parameters as a query string including the leading '?'
// or it returns the empty string if this is the default page settings
func (p *pageFilter) ToUrlSuffix() string {
	if p != nil {
		if q, ok := p.ToQuery(); ok {
			return q.QueryString()
		}
	}
	return ""
}

type (
	// Query helps with building a url query string
	Query interface {

		// Empty returns true if there are no parameters added to the query
		Empty() bool

		// Add adds a query parameter, both name and value are query escaped for you.
		// You can add multiple parameters with the same name, if you think that makes sense.
		Add(name, value string) Query

		// AddEscaped adds a query parameter, name and value are assumed to be already query escaped.
		// You can add multiple parameters with the same name, if you think that makes sense.
		AddEscaped(name, value string) Query

		// AddIf adds a query parameter if condition is true.
		AddIf(condition bool, name, value string) Query

		// String returns the encoded query string excluding the leading question mark
		String() string

		// QueryString returns the encoded query string including the leading question mark
		QueryString() string
	}

	query struct {
		Params []param
	}

	param struct {
		Name      string
		Value     string
		IsEscaped bool
	}
)

// NewQuery create a new Query builder
func NewQuery() Query {
	return &query{}
}

func (q *query) Add(name, value string) Query {
	q.Params = append(q.Params, param{Name: name, Value: value, IsEscaped: false})
	return q
}

func (q *query) AddEscaped(name, value string) Query {
	q.Params = append(q.Params, param{Name: name, Value: value, IsEscaped: true})
	return q
}

func (q *query) AddIf(condition bool, name, value string) Query {
	if condition {
		return q.Add(name, value)
	}
	return q
}

func (q *query) Empty() bool {
	return len(q.Params) == 0
}

func (q *query) String() string {
	return q.buildQueryString(false)
}

func (q *query) QueryString() string {
	return q.buildQueryString(true)
}

func (q *query) buildQueryString(includeQuestionMark bool) string {
	if q.Empty() {
		return ""
	}

	str := &strings.Builder{}

	if includeQuestionMark {
		str.WriteString("?")
	}

	for i := 0; i < len(q.Params); i++ {
		p := q.Params[i]
		if i > 0 {
			str.WriteString("&")
		}

		name := p.Name
		value := p.Value

		if !p.IsEscaped {
			name = url.QueryEscape(name)
			value = url.QueryEscape(value)
		}

		str.WriteString(name)
		str.WriteString("=")
		str.WriteString(value)
	}
	return str.String()
}

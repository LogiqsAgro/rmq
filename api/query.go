package api

import (
	"fmt"
	"net/url"
	"strings"
)

type (
	page struct {
		Page     int
		Size     int
		Name     string
		UseRegex bool
	}
)

func NewPage(num, size int) *page {
	return NewPageFilter(num, size, "", false)
}

func NewPageFilter(num, size int, name string, useRegex bool) *page {
	if num < 1 {
		num = 1
	}
	if size < 1 {
		size = 100
	}
	return &page{
		Page:     num,
		Size:     size,
		Name:     name,
		UseRegex: useRegex,
	}
}

func newPage() *page {
	return &page{
		Page:     1,
		Size:     100,
		Name:     "",
		UseRegex: false,
	}
}

// ToQuery returns the Query representing this page
func (p *page) ToQuery() (Query, bool) {
	q := NewQuery()

	if p.Page > 1 {
		q.Add("page", fmt.Sprintf("%d", p.Page))
	}

	if p.Size > 0 && p.Size != 100 {
		q.Add("page_size", fmt.Sprintf("%d", p.Size))
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
func (p *page) ToUrlSuffix() string {
	if p != nil {
		if q, ok := p.ToQuery(); ok {
			return "?" + q.String()
		}
	}
	return ""
}

type (
	// Query helps with building a url query string
	Query interface {

		// Empty returns true if there are no parameters added to the query
		Empty() bool

		// Add adds a query parameter, both name and value are query encoded for you.
		// You can add multiple parameters with the same name, if you think that makes sense.
		Add(name, value string) Query

		// AddIf adds a query parameter if condition is true.
		AddIf(condition bool, name, value string) Query

		// String returns the encoded query string not including the leading question mark
		String() string
	}

	query struct {
		Params []param
	}

	param struct {
		Name  string
		Value string
	}
)

// NewQuery create a new Query builder
func NewQuery() Query {
	return &query{}
}

func (q *query) Add(name, value string) Query {
	q.Params = append(q.Params, param{Name: name, Value: value})
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
	if q.Empty() {
		return ""
	}

	str := &strings.Builder{}
	for i := 0; i < len(q.Params); i++ {
		p := q.Params[i]
		if i > 0 {
			str.WriteString("&")
		}
		str.WriteString(url.QueryEscape(p.Name))
		str.WriteString("=")
		str.WriteString(url.QueryEscape(p.Value))
	}
	return str.String()
}

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

package definitions

import (
	"bytes"
	"encoding/json"
	"errors"
	"sort"
)

type Schema struct {
	Version   string      `json:"version"`
	Endpoints []*Endpoint `json:"endpoints,omitempty"`
}

func NewSchema() *Schema {
	return &Schema{}
}

func (s *Schema) initialize() {
	for _, ep := range s.Endpoints {
		ep.initialize()
	}
}

type FeatureSet struct {
	Version  string              `json:"version"`
	Features map[string]*Feature `json:"features,omitempty"`
}

func NewFeatureSet() *FeatureSet {
	fs := &FeatureSet{}
	fs.initialize()
	return fs
}

func (fs *FeatureSet) initialize() {
	if fs.Features == nil {
		fs.Features = map[string]*Feature{}
	} else {
		for _, f := range fs.Features {
			f.initialize()
		}
	}
}

type Feature struct {
	Description     string                `json:"description,omitempty"`
	QueryParameters map[string]*Parameter `json:"query-parameters,omitempty"`
}

func NewFeature() *Feature {
	f := &Feature{}
	f.initialize()
	return f
}

func (f *Feature) initialize() {
	if f.QueryParameters == nil {
		f.QueryParameters = make(map[string]*Parameter)
	} else {
		for _, p := range f.QueryParameters {
			p.initialize()
		}
	}
}

type Parameter struct {
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Mandatory   bool   `json:"mandatory,omitempty"`
	Default     string `json:"default,omitempty"`
}

func NewParameter() *Parameter {
	p := &Parameter{}
	p.initialize()
	return p
}

func (p *Parameter) initialize() {

}

type Endpoint struct {
	Description string `json:"description"`
	Path        string `json:"path"`

	// PathParameters is not encoded to JSON, after decoding use InitPathParameters to populate this field
	PathParameters []string                  `json:"-"`
	Methods        map[string]*MethodDetails `json:"methods,omitempty"`
}

func NewEndpoint() *Endpoint {
	ep := &Endpoint{}
	ep.initialize()
	return ep
}

func (ep *Endpoint) initialize() {
	if ep.Methods == nil {
		ep.Methods = make(map[string]*MethodDetails)
	} else {
		for _, details := range ep.Methods {
			details.initialize()
		}
	}
}

func (ep *Endpoint) HasParameter(name string) bool {
	for i := 0; i < len(ep.PathParameters); i++ {
		if name == ep.PathParameters[i] {
			return true
		}
	}
	return false
}

func (ep *Endpoint) AllMethods() []string {
	ep.initialize()
	kk := make([]string, 0, len(ep.Methods))
	for k := range ep.Methods {
		kk = append(kk, k)
	}
	sort.Strings(kk)
	return kk
}

// InitPathParameters clears the PathParameters fields, and repopulates it from '{param}' parameter declarations in the Path field.
func (ep *Endpoint) InitPathParameters() error {
	param := &bytes.Buffer{}
	ep.PathParameters = []string{}
	inParam := false
	for i := 0; i < len(ep.Path); i++ {
		switch ep.Path[i] {
		case '{':
			if inParam {
				return errors.New("invalid path parameter sytax in path: nested '{' before '}'")
			}
			inParam = true
		case '}':
			inParam = false
			ep.PathParameters = append(ep.PathParameters, param.String())
			param.Reset()
		default:
			if inParam {
				param.WriteByte(ep.Path[i])
			}
		}
	}
	return nil
}

type MethodDetails struct {
	Features map[string]bool `json:"features,omitempty"`
	CliVerbs []string        `json:"cli-verbs,omitempty"`
}

func NewMethodDetails() *MethodDetails {
	details := &MethodDetails{}
	details.initialize()
	return details
}
func (m *MethodDetails) initialize() {
	if m.Features == nil {
		m.Features = make(map[string]bool)
	}
}

func (m *MethodDetails) AddFeature(f ...string) {
	m.initialize()
	for j := 0; j < len(f); j++ {
		m.Features[f[j]] = true
	}
}

func (m *MethodDetails) DelFeature(f ...string) {
	m.initialize()
	for j := 0; j < len(f); j++ {
		delete(m.Features, f[j])
	}
}

func (m *MethodDetails) AddCliVerb(v ...string) {
	m.CliVerbs = append(m.CliVerbs, v...)
}

type List struct {
	items []*Endpoint
}

func NewList() *List {
	return &List{
		items: make([]*Endpoint, 0, 64),
	}
}

func (list *List) Len() int {
	return len(list.items)
}

func (list *List) AppendNew() *Endpoint {
	ep := NewEndpoint()
	list.Append(ep)
	return ep
}

func (list *List) Append(ep *Endpoint) {
	if ep == nil {
		return
	}
	list.items = append(list.items, ep)
}

func (list *List) ToSlice() []*Endpoint {
	var items = make([]*Endpoint, list.Len())
	copy(items, list.items)
	return items
}

func (list *List) ToJson() ([]byte, error) {
	return json.Marshal(list.items)
}

func (list *List) ToIndentedJson() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(list.items); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

type builder struct {
	ep *Endpoint
}

func Builder() *builder {
	return &builder{ep: NewEndpoint()}
}

func (b *builder) AppendPath(s string) *builder {
	b.ep.Path += s
	return b
}

func (b *builder) Path(path string) *builder {
	b.ep.Path = path
	return b
}

func (b *builder) AddMethod(method string) *builder {
	if len(method) == 0 {
		return b
	}

	if _, ok := b.ep.Methods[method]; ok {
		return b
	}

	b.ep.Methods[method] = NewMethodDetails()
	return b
}

func (b *builder) Methods(methods ...string) *builder {
	for i := 0; i < len(methods); i++ {
		b.AddMethod(methods[i])
	}
	return b
}

func (b *builder) PathParameters(parameters ...string) *builder {
	b.ep.PathParameters = parameters
	return b
}

func (b *builder) Features(verb string, features ...string) *builder {
	b.AddMethod(verb)
	b.ep.Methods[verb].AddFeature(features...)
	return b
}

func (b *builder) Description(description string) *builder {
	b.ep.Description = description
	return b
}

func (b *builder) Clear() *builder {
	b.ep = NewEndpoint()
	return b
}

func (b *builder) Build() (*Endpoint, error) {
	ep := b.ep
	err := ep.InitPathParameters()
	if err != nil {
		return nil, err
	}
	b.Clear()
	return ep, nil
}

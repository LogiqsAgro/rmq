package endpoint

import (
	"bytes"
	"encoding/json"
	"errors"
)

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
type (
	Endpoint struct {
		Path           string   `json:"path"`
		Verbs          []string `json:"verbs"`
		PathParameters []string `json:"pathParameters"`
		Description    string   `json:"description"`
	}
)

func New() *Endpoint {
	return &Endpoint{}
}

func (ep *Endpoint) HasParameter(name string) bool {
	for i := 0; i < len(ep.PathParameters); i++ {
		if name == ep.PathParameters[i] {
			return true
		}
	}
	return false
}

// InitPathParameters clears the PathParameters fields, and repopulates it from '{param}' parameter declarations in the Path field.
func (ep *Endpoint) InitPathParameters() error {
	param := new(bytes.Buffer)
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

type (
	List struct {
		items []*Endpoint
	}
)

func NewList() *List {
	return &List{
		items: make([]*Endpoint, 0, 64),
	}
}

func (list *List) Len() int {
	return len(list.items)
}

func (list *List) AppendNew() *Endpoint {
	ep := New()
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

type (
	builder struct {
		ep Endpoint
	}
)

func Builder() *builder {
	return &builder{}
}

func (b *builder) AppendPath(s string) *builder {
	b.ep.Path += s
	return b
}

func (b *builder) Path(path string) *builder {
	b.ep.Path = path
	return b
}

func (b *builder) AddVerb(verb string) *builder {
	b.ep.Verbs = append(b.ep.Verbs, verb)
	return b
}

func (b *builder) Verbs(verbs ...string) *builder {
	b.ep.Verbs = verbs
	return b
}

func (b *builder) PathParameters(parameters ...string) *builder {
	b.ep.PathParameters = parameters
	return b
}

func (b *builder) Description(description string) *builder {
	b.ep.Description = description
	return b
}

func (b *builder) Clear() *builder {
	b.ep = Endpoint{}
	return b
}

func (b *builder) Build() (*Endpoint, error) {
	ep := b.ep
	err := ep.InitPathParameters()
	if err != nil {
		return nil, err
	}
	b.Clear()
	return &ep, nil
}

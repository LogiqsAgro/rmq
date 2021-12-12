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
package generator

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/LogiqsAgro/rmq/api-gen/endpoint"
)

type (
	generator struct {
		Package           string
		RabbitMQVersion   string
		namer             *uniqueNamer
		eps               []*endpoint.Endpoint
		paramTypes        map[string]string
		pathParamReplacer *regexp.Regexp
	}
)

func New(list *endpoint.List) *generator {
	g := &generator{
		namer:             newUniqueNamer(),
		eps:               list.ToSlice(),
		paramTypes:        buildParamTypeMap(list),
		pathParamReplacer: regexp.MustCompile(`\{\w+\}`),
		Package:           "",
		RabbitMQVersion:   "???",
	}
	return g
}

func (g *generator) Generate(w io.Writer) {
	eps := g.eps

	fmt.Fprintln(w, "// Code generated by go generate; DO NOT EDIT.")
	fmt.Fprintf(w, "//\n// last generated at %s\n//\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(w, "package %s\n", g.Package)
	fmt.Fprintln(w, "")

	// Imports
	fmt.Fprintln(w, "import (")
	imports := []string{"net/http", "net/url", "fmt"}
	sort.Strings(imports)
	for i := 0; i < len(imports); i++ {

		fmt.Fprintf(w, "\t%q\n", imports[i])
	}
	fmt.Fprintln(w, ")")

	// generate the RabbitMQ version
	if len(g.RabbitMQVersion) > 0 {
		fmt.Fprintf(w, "\n//RabbitMQVersion show the version against which this api was generated\nfunc RabbitMQVersion() string { return %q }\n\n", g.RabbitMQVersion)
	}

	for i := 0; i < len(eps); i++ {
		ep := eps[i]
		for j := 0; j < len(ep.Verbs); j++ {
			verb := ep.Verbs[j]

			funcName := g.getEndpointMethodName(verb, ep)

			fmt.Fprint(w, "// ")
			fmt.Fprint(w, strings.ReplaceAll(ep.Description, "\n", "\n// "))
			fmt.Fprintln(w)

			var args = ""

			for i := 0; i < len(ep.PathParameters); i++ {
				if i > 0 {
					args += ", "
				}
				param := ep.PathParameters[i]
				args += param + " " + g.paramTypes[param]
			}

			fmt.Fprintf(w, "// %s %s\n", verb, ep.Path)
			fmt.Fprintf(w, "func %s(%s) Builder {\n", funcName, args)
			path := ""
			if len(ep.PathParameters) == 0 {
				path = fmt.Sprintf("%q", ep.Path)
			} else {
				urlParams := ""
				for i := 0; i < len(ep.PathParameters); i++ {
					p := ep.PathParameters[i]
					pt := g.paramTypes[p]
					switch pt {
					case "string":
						urlParams += fmt.Sprintf(", url.PathEscape(%s)", p)
					case "int":
						urlParams += fmt.Sprintf(", %s", p)
					default:
						err := fmt.Errorf("unsupported path parameter type: %q", pt)
						panic(err)
					}
				}
				path = fmt.Sprintf("fmt.Sprintf(%q%s)", g.pathParamReplacer.ReplaceAllString(ep.Path, "%v"), urlParams)

			}
			fmt.Fprintf(w, "\tpath := %s\n", path)
			fmt.Fprintf(w, "\treturn Request().Method(http.Method%s).Path(path)\n", capitalize(verb))
			fmt.Fprintf(w, "}\n\n")
		}
	}

	//fmt.Fprintf(w, "// All function names:\n// %s", strings.Join(g.namer.AllNames(), "\n// "))
}

func (g *generator) getEndpointMethodName(verb string, ep *endpoint.Endpoint) string {
	funcName := capitalize(verb)
	pathParts := strings.Split(ep.Path, "/")

	for i := 0; i < len(pathParts); i++ {
		if pathParts[i] == "api" {
			continue
		}

		if strings.Contains(pathParts[i], "{") {
			paramName := strings.Trim(pathParts[i], "{}")
			if paramName == "name" || paramName == "channel" {
				if strings.HasSuffix(funcName, "ies") {
					funcName = funcName[:len(funcName)-3] + "y"
				} else {
					funcName = strings.TrimRight(funcName, "s")
				}
			}
		} else {
			funcName += capitalize(pathParts[i])
		}
	}

	if len(ep.PathParameters) == 1 {
		switch ep.PathParameters[0] {
		case "vhost":
			funcName = strings.ReplaceAll(funcName, "Vhosts", "") + "ForVhost"
		}
	} else {
		first := true
		for i := 0; i < len(ep.PathParameters); i++ {
			param := ep.PathParameters[i]
			if param == "name" {
				continue
			}
			if first {
				funcName += "For"
				first = false
			} else {
				funcName += "And"
			}
			funcName += capitalize(param)
		}
	}

	funcName = g.namer.UniqueName(funcName)
	return funcName
}

func buildParamTypeMap(list *endpoint.List) map[string]string {
	eps := list.ToSlice()
	paramTypes := make(map[string]string)
	for i := 0; i < len(eps); i++ {
		ep := eps[i]
		for i := 0; i < len(ep.PathParameters); i++ {
			p := ep.PathParameters[i]
			paramTypes[p] = "string"
		}
	}

	paramTypes["port"] = "int"
	paramTypes["within"] = "int"
	return paramTypes
}

func capitalize(s string) string {
	s = strings.ToLower(s)
	s = strings.Title(s)
	s = strings.ReplaceAll(s, "-", "")
	return s
}

type (
	uniqueNamer struct {
		exists map[string]bool
	}
)

func newUniqueNamer() *uniqueNamer {
	return &uniqueNamer{
		exists: make(map[string]bool),
	}
}

func (un *uniqueNamer) AllNames() []string {
	names := make([]string, 0, len(un.exists))
	for k := range un.exists {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
func (un *uniqueNamer) Exists(name string) bool {
	return un.exists[name]
}

func (un *uniqueNamer) TryAdd(name string) bool {
	if un.exists[name] {
		return false
	}

	un.exists[name] = true
	return true
}

func (un *uniqueNamer) UniqueName(name string) string {
	n := 1
	format := name + "%d"
	for {
		if un.TryAdd(name) {
			return name
		} else {
			n += 1
			name = fmt.Sprintf(format, n)
			continue
		}
	}
}
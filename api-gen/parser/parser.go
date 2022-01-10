package parser

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
import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"strings"

	endpoint "github.com/LogiqsAgro/rmq/api/definitions"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

func ParseApiDocs(r io.Reader) (*endpoint.List, error) {
	tokenizer := html.NewTokenizer(r)
	endpoints := endpoint.NewList()

	for {
		tokenType := tokenizer.Next()

		switch {
		case tokenType == html.ErrorToken:
			return nil, errors.New("Error parsing html")
		case tokenType == html.StartTagToken:
			token := tokenizer.Token()
			switch {
			case token.Data == "table":
				err := parseTable(tokenizer, endpoints)
				if err != nil {
					return nil, err
				} else {

					return endpoints, nil
				}
			}
		default:
			// do nothing
		}

	}
}

func parseTable(tokenizer *html.Tokenizer, endpoints *endpoint.List) error {

	consumeHeaderRow(tokenizer)

	for {
		tt := tokenizer.Next()
		tagBytes, _ := tokenizer.TagName()
		tag := string(tagBytes)

		switch {
		case tag == "table" && tt == html.EndTagToken:
			return nil
		case tag == "tr" && tt == html.StartTagToken:
			if err := parseEndpointRow(tokenizer, endpoints); err != nil {
				return err
			}
		}
	}
}

var (
	newLineAndTrailingWhitespace *regexp.Regexp
)

func init() {
	re, err := regexp.Compile("\n\\s*")
	if err != nil {
		panic(err)
	}
	newLineAndTrailingWhitespace = re
}

func parseEndpointRow(tokenizer *html.Tokenizer, endpoints *endpoint.List) error {
	b := endpoint.Builder()
	if get, _, err := getNextCellContents(tokenizer); err == nil {
		if strings.Contains(get, "X") {
			b.AddMethod(http.MethodGet)
		}
	} else {
		return err
	}

	if put, _, err := getNextCellContents(tokenizer); err == nil {
		if strings.Contains(put, "X") {
			b.AddMethod(http.MethodPut)
		}
	} else {
		return err
	}

	if del, _, err := getNextCellContents(tokenizer); err == nil {
		if strings.Contains(del, "X") {
			b.AddMethod(http.MethodDelete)
		}
	} else {
		return err
	}

	if post, _, err := getNextCellContents(tokenizer); err == nil {
		if strings.Contains(post, "X") {
			b.AddMethod(http.MethodPost)
		}
	} else {
		return err
	}

	if path, _, err := getNextCellContents(tokenizer); err == nil {
		path = newLineAndTrailingWhitespace.ReplaceAllString(path, "\n")
		path = strings.Trim(path, "\n ")
		b.Path(path)
	} else {
		return err
	}

	if desc, features, err := getNextCellContents(tokenizer); err == nil {
		desc = newLineAndTrailingWhitespace.ReplaceAllString(desc, "\n")
		desc = strings.Trim(desc, "\n ")
		b.Features(http.MethodGet, features...)
		b.Description(desc)
	} else {
		return err
	}

	ep, err := b.Build()
	if err != nil {
		return err
	}

	//
	// Fixup multiple paths in the path cell, and deprecation notices.
	//
	if strings.Contains(ep.Path, "\n") {
		paths := strings.Split(ep.Path, "\n")
		for i := 0; i < len(paths); i++ {
			path := strings.Trim(paths[i], "\n ")
			newep := *ep
			deprecated := "(deprecated)"
			if strings.Contains(path, deprecated) {
				path = strings.ReplaceAll(path, deprecated, "")
				path = strings.Trim(path, "\n ")
				newep.Description = deprecated + " " + ep.Description
			}

			newep.Path = path
			newep.InitPathParameters()
			endpoints.Append(&newep)
		}

	} else {
		endpoints.Append(ep)
	}

	if err := moveToRowEnd(tokenizer); err != nil {
		return err
	}

	return nil
}

func consumeHeaderRow(tokenizer *html.Tokenizer) error {
	if err := moveToRowStart(tokenizer); err != nil {
		return err
	}
	if err := moveToRowEnd(tokenizer); err != nil {
		return err
	}
	return nil

}

func moveToRowStart(tokenizer *html.Tokenizer) error {
	return moveTo(tokenizer, html.StartTagToken, "tr")
}

func moveToRowEnd(tokenizer *html.Tokenizer) error {
	return moveTo(tokenizer, html.EndTagToken, "tr")
}
func moveToCellStart(tokenizer *html.Tokenizer) error {
	return moveTo(tokenizer, html.StartTagToken, "td")
}

func moveTo(tokenizer *html.Tokenizer, tokenType html.TokenType, tagName string) error {

	token := tokenizer.Token()
	tt := token.Type

	for {
		if tt == html.ErrorToken {
			return tokenizer.Err()
		}

		if tt == tokenType {
			name, _ := tokenizer.TagName()
			if bytes.Equal(name, []byte(tagName)) {
				return nil
			}
		}

		tt = tokenizer.Next()
	}
}

func getNextCellContents(tokenizer *html.Tokenizer) (text string, features []string, err error) {
	if err = moveToCellStart(tokenizer); err != nil {
		return "", nil, err
	}
	href := []byte("href")
	hash := []byte("#")
	contents := &bytes.Buffer{}
	features = []string{}
	for {
		tt := tokenizer.Next()
		tagBytes, _ := tokenizer.TagName()
		tag := string(tagBytes)

		switch {
		case tt == html.TextToken:
			contents.Write(tokenizer.Text())
		case tag == "i" && tt == html.StartTagToken:
			contents.WriteString("{")
		case tag == "i" && tt == html.EndTagToken:
			contents.WriteString("}")
		case tag == "td" && tt == html.EndTagToken:
			return contents.String(), features, nil
		case tag == "a" && tt == html.StartTagToken:
			for {
				// if the href of the link only contains a fragment, like #pagination
				// we assume this to be a link to an endpoint feature
				k, v, next := tokenizer.TagAttr()
				if bytes.Equal(k, href) && bytes.HasPrefix(v, hash) {
					if err == nil {
						feature := string(v[len(hash):])
						features = append(features, feature)
					}
				}
				if !next {
					break
				}
			}
		case tt == html.ErrorToken:
			return "", nil, tokenizer.Err()
		default:
			contents.Write(tokenizer.Text())
		}
	}
}

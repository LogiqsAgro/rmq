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

	"github.com/LogiqsAgro/rmq/api-gen/endpoint"
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
	if get, err := getNextCellContents(tokenizer); err == nil {
		if strings.Contains(get, "X") {
			b.AddVerb(http.MethodGet)
		}
	} else {
		return err
	}

	if put, err := getNextCellContents(tokenizer); err == nil {
		if strings.Contains(put, "X") {
			b.AddVerb(http.MethodPut)
		}
	} else {
		return err
	}

	if del, err := getNextCellContents(tokenizer); err == nil {
		if strings.Contains(del, "X") {
			b.AddVerb(http.MethodDelete)
		}
	} else {
		return err
	}

	if post, err := getNextCellContents(tokenizer); err == nil {
		if strings.Contains(post, "X") {
			b.AddVerb(http.MethodPost)
		}
	} else {
		return err
	}

	if path, err := getNextCellContents(tokenizer); err == nil {
		path = newLineAndTrailingWhitespace.ReplaceAllString(path, "\n")
		path = strings.Trim(path, "\n ")
		b.Path(path)
	} else {
		return err
	}

	if desc, err := getNextCellContents(tokenizer); err == nil {
		desc = newLineAndTrailingWhitespace.ReplaceAllString(desc, "\n")
		desc = strings.Trim(desc, "\n ")
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

func getNextCellContents(tokenizer *html.Tokenizer) (string, error) {
	if err := moveToCellStart(tokenizer); err != nil {
		return "", err
	}
	contents := make([]byte, 0, 64)
	for {
		tt := tokenizer.Next()
		tagBytes, _ := tokenizer.TagName()
		tag := string(tagBytes)

		switch {
		case tt == html.TextToken:
			contents = append(contents, tokenizer.Text()...)
		case tag == "i" && tt == html.StartTagToken:
			contents = append(contents, []byte("{")...)
		case tag == "i" && tt == html.EndTagToken:
			contents = append(contents, []byte("}")...)
		case tag == "td" && tt == html.EndTagToken:
			return string(contents), nil
		case tt == html.ErrorToken:
			return "", tokenizer.Err()
		default:
			contents = append(contents, tokenizer.Text()...)
		}
	}
}

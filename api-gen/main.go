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
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/LogiqsAgro/rmq/api-gen/generator"
	"github.com/LogiqsAgro/rmq/api-gen/parser"
	"github.com/pkg/errors"
)

func main() {

	run()
	// data, err := endpoints.ToIndentedJson()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(data))
}

func run() {
	client := &http.Client{}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fileName := os.Getenv("GOFILE")
	if strings.HasSuffix(fileName, ".go") {
		fileName = fileName[:len(fileName)-3]
		fileName = filepath.Join(cwd, fileName+".g.go")
	} else {
		panic(errors.New("Expected a go file as input but got: '" + fileName + "'"))
	}

	resp, err := client.Get("http://localhost:15672/api/index.html")
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := errors.Errorf("Fetching RabbitMQ api documentation failed: " + resp.Status)
		panic(err)
	}

	endpoints, err := parser.ParseApiDocs(resp.Body)
	if err != nil {
		panic(err)
	}

	overviewUrl := "http://localhost:15672/api/overview?columns=rabbitmq_version,management_version"
	req, err := http.NewRequest(http.MethodGet, overviewUrl, nil)
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth("guest", "guest")
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("could not get rabbitmq version from %s: %s", overviewUrl, resp.Status))
	}

	versions := struct {
		RabbitMQVersion   string `json:"rabbitmq_version"`
		ManagementVersion string `json:"management_version"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&versions)
	if err != nil {
		panic(err)
	}

	var g = generator.New(endpoints)
	g.Package = os.Getenv("GOPACKAGE")
	g.RabbitMQVersion = versions.RabbitMQVersion
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	g.Generate(f)

}

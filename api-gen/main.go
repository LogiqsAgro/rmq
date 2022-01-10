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
	endpoint "github.com/LogiqsAgro/rmq/api/definitions"
	"github.com/pkg/errors"
)

func main() {
	run()
}

func run() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}()

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	inFile := os.Getenv("GOFILE")
	inFile = filepath.Join(cwd, inFile)

	client := &http.Client{}
	endpoints, err := parseApiDocs(client)
	if err != nil {
		panic(err)
	}

	rmqVersion, _, err := fetchVersions(client, "guest", "guest")
	if err != nil {
		panic(err)
	}

	cmd := command()
	switch cmd {
	case "dump-endpoints":
		outFile := endpoint.SchemaFileName(rmqVersion)
		outFile = filepath.Join(cwd, outFile)

		obj := struct {
			Version   string               `json:"version"`
			Endpoints []*endpoint.Endpoint `json:"endpoints"`
		}{
			Version:   rmqVersion,
			Endpoints: endpoints.ToSlice(),
		}

		f, err := os.Create(outFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		err = enc.Encode(obj)
		if err != nil {
			panic(err)
		}

	case "gen-api":
		outFile := ""

		if strings.HasSuffix(inFile, ".go") {
			outFile = inFile[:len(inFile)-3]
			outFile = outFile + ".g.go"
		} else {
			panic(errors.New("Expected a go file as input but got: '" + inFile + "'"))
		}
		var g = generator.New(endpoints)
		g.Package = os.Getenv("GOPACKAGE")
		g.RabbitMQVersion = rmqVersion
		f, err := os.Create(outFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		g.Generate(f)
	default:
		panic(fmt.Errorf("Unrecognized command: '%s'", cmd))

	}

}

func command() string {
	args := os.Args
	if len(args) > 1 {
		return args[1]
	}
	return "gen-api"
}

func parseApiDocs(client *http.Client) (*endpoint.List, error) {
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
		return nil, err
	}
	return endpoints, nil
}

func fetchVersions(client *http.Client, username, password string) (rmqVersion, mgmtVersion string, err error) {
	overviewUrl := "http://localhost:15672/api/overview?columns=rabbitmq_version,management_version"
	req, err := http.NewRequest(http.MethodGet, overviewUrl, nil)
	if err != nil {
		return "", "", err
	}
	req.SetBasicAuth("guest", "guest")
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("could not get rabbitmq version from %s: %s", overviewUrl, resp.Status)
	}

	versions := struct {
		RabbitMQVersion   string `json:"rabbitmq_version"`
		ManagementVersion string `json:"management_version"`
	}{
		RabbitMQVersion:   "",
		ManagementVersion: "",
	}

	err = json.NewDecoder(resp.Body).Decode(&versions)
	if err != nil {
		return "", "", err
	}

	return versions.RabbitMQVersion, versions.ManagementVersion, nil
}

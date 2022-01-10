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
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:generate go run ..\..\api-gen\main.go dump-endpoints

//go:embed *.json
var fs embed.FS

// LoadSchema loads the embedded schema json for the given RabbitMQ version
func LoadFeatures(version string) (*FeatureSet, error) {
	if !IsValidVersionSpec(version) {
		return nil, fmt.Errorf("'%s' is not a valid version number", version)
	}

	f, err := fs.Open(FeatureFileName(version))
	if err != nil {
		return nil, fmt.Errorf("error loading feature set version '%s': %w", version, err)
	}

	featureSet := NewFeatureSet()
	err = json.NewDecoder(f).Decode(featureSet)
	if err != nil {
		return nil, fmt.Errorf("error loading feature set version '%s': %w", version, err)
	}

	featureSet.initialize()
	return featureSet, err
}

// LoadSchema loads the embedded schema json for the given RabbitMQ version
func LoadSchema(version string) (*Schema, error) {
	if !IsValidVersionSpec(version) {
		return nil, fmt.Errorf("'%s' is not a valid version number", version)
	}

	f, err := fs.Open(SchemaFileName(version))
	if err != nil {
		return nil, fmt.Errorf("error loading schema version '%s': %w", version, err)
	}

	schema := NewSchema()
	err = json.NewDecoder(f).Decode(schema)
	if err != nil {
		return nil, fmt.Errorf("error loading schema version '%s': %w", version, err)
	}

	schema.initialize()
	return schema, err
}

// AvailableSchemaVersions returns all the versiona you can use with LoadSchema(...)
func AvailableSchemaVersions() []string {
	files, err := fs.ReadDir(".")
	if err != nil {
		return nil
	}

	versions := make([]string, 0, 8)
	for i := 0; i < len(files); i++ {
		file := files[i]
		name := file.Name()
		if strings.HasPrefix(name, "v") && strings.HasSuffix(name, schemaFileSuffix) {
			version := name[1 : len(name)-len(schemaFileSuffix)]
			if IsValidVersionSpec(version) {
				versions = append(versions, version)
			}
		}
	}

	return versions
}

const (
	schemaFileSuffix  string = "-endpoints.json"
	featureFileSuffix string = "-features.json"
)

func SchemaFileName(version string) string {
	return fmt.Sprintf("v%s%s", version, schemaFileSuffix)
}

func FeatureFileName(version string) string {
	return fmt.Sprintf("v%s%s", version, featureFileSuffix)
}

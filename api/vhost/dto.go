//
// DTO's for the /api/definitions/[vhost] managment api endpoint
//

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
package vhost

type (
	Definition struct {
		RabbitVersion string        `json:"rabbit_version"`
		Parameters    []interface{} `json:"parameters"`
		Policies      []interface{} `json:"policies"`
		Queues        []*Queue      `json:"queues"`
		Exchanges     []*Exchange   `json:"exchanges"`
		Bindings      []*Binding    `json:"bindings"`
	}

	Queue struct {
		Name       string                 `json:"name"`
		Durable    bool                   `json:"durable"`
		AutoDelete bool                   `json:"auto_delete"`
		Arguments  map[string]interface{} `json:"arguments"`
	}

	Exchange struct {
		Name       string                 `json:"name"`
		Type       string                 `json:"type"`
		Durable    bool                   `json:"durable"`
		AutoDelete bool                   `json:"auto_delete"`
		Internal   bool                   `json:"internal"`
		Arguments  map[string]interface{} `json:"arguments"`
	}

	Binding struct {
		VHost           string                 `json:"vhost"`
		Source          string                 `json:"source"`
		Destination     string                 `json:"destination"`
		DestinationType string                 `json:"destination_type"`
		PropertiesKey   string                 `json:"properties_key"`
		RoutingKey      string                 `json:"routing_key"`
		Arguments       map[string]interface{} `json:"arguments"`
	}
)

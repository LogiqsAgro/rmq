//
// DTO's for the /api/definitions/[vhost] managment api endpoint
//

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

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
package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/LogiqsAgro/rmq/api/vhost"
)

const (
	UnitDays   = "days"
	UnitWeeks  = "weeks"
	UnitMonths = "months"
	UnitYears  = "years"
)

// CertificateExpirationTimeUnits returns the list of time units usable in GetHealthChecksCertificateExpirationJson(within int, unit string)
func CertificateExpirationTimeUnits() []string {
	return []string{
		UnitDays,
		UnitWeeks,
		UnitMonths,
		UnitYears,
	}
}

// GetOverviewJson returns various random bits of information that describe the whole system. ( GET /api/overview )
func GetOverviewJson() ([]byte, error) {
	resp, err := Get("overview")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetFederationLinksJson returns status for all federation links. Requires the rabbitmq_federation_management plugin to be enabled. ( GET /api/federation-links )
func GetFederationLinksJson() ([]byte, error) {
	resp, err := Get("federation-links")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVhostFederationLinksJson returns status for federation links of a vhost. Requires the rabbitmq_federation_management plugin to be enabled. ( GET /api/federation-links/<vhost> )
func GetVhostFederationLinksJson(name string) ([]byte, error) {
	pnq := fmt.Sprintf("federation-links/%s", url.PathEscape(name))
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetClusterName returns the name identifying this RabbitMQ cluster.
// ( GET /api/cluster-name )
func GetClusterName() ([]byte, error) {
	resp, err := Get("cluster-name")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetNodesJson returns a list of nodes in the RabbitMQ cluster.
// ( GET /api/nodes )
func GetNodesJson() ([]byte, error) {
	resp, err := Get("nodes")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetNodeJson returns An individual node in the RabbitMQ cluster.
// Add "?memory=true" to get memory statistics,
// and "?binary=true" to get a breakdown of binary memory use
// (may be expensive if there are many small binaries in the system).
// ( GET /api/nodes/<name> )
func GetNodeJson(name string, memory, binary bool) ([]byte, error) {
	pnq := "nodes/" +
		url.PathEscape(name) +
		NewQuery().
			AddIf(memory, "memory", "true").
			AddIf(binary, "binary", "true").
			QueryString()

	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetExtensionsJson returns a list of extensions to the management plugin.
// ( GET /api/extensions)
func GetExtensionsJson() ([]byte, error) {
	resp, err := Get("extensions")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetDefinitionsJson returns the server definitions - exchanges, queues, bindings, users,
// virtual hosts, permissions, topic permissions, and parameters.
// Everything apart from messages.
// ( GET /api/definitions )
func GetDefinitionsJson() ([]byte, error) {
	resp, err := Get("definitions")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostDefinitionsJson returns the server definitions for a given virtual host -
// exchanges, queues, bindings and policies.
// ( GET /api/definitions/<name> )
func GetVHostDefinitionsJson(name string) ([]byte, error) {
	pnq := "definitions/" + url.PathEscape(name)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostDefinitions returns the server definitions for a given virtual host -
// exchanges, queues, bindings and policies.
// ( GET /api/definitions/<name> )
func GetVHostDefinitions(name string) (*vhost.Definition, error) {
	if data, err := GetVHostDefinitionsJson(name); err != nil {
		return nil, err
	} else {
		definition := &vhost.Definition{}
		err := json.Unmarshal(data, definition)
		if err != nil {
			return nil, err
		}

		return definition, nil
	}
}

// GetConnectionsJson returns a list of all open connections.
// Use pagination parameter page to filter connections.
// Use nil for default page (page 1, size 100, no name filter)
// see api.NewPage(...)  or api.NewPageFilter(...)
// ( GET /api/connections )
func GetConnectionsJson(page *pageFilter) ([]byte, error) {
	pnq := "connections" + page.ToUrlSuffix()
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostConnectionsJson returns a list of all open connections for the vhost.
// Use pagination parameter page to filter connections.
// Use nil for default page (page 1, size 100, no name filter)
// see api.NewPage(...)  or api.NewPageFilter(...)
// ( GET /api/vhosts/<vhost>/connections )
func GetVHostConnectionsJson(vhost string, page *pageFilter) ([]byte, error) {
	pnq := "vhosts/" + url.PathEscape(vhost) + "/connections" + page.ToUrlSuffix()
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetConnectionJson returns the named connection
// ( GET /api/connections/<name> )
func GetConnectionJson(name string) ([]byte, error) {
	pnq := "connections/" + url.PathEscape(name)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetConnectionChannelsJson list of all channels for the given connection.
// ( GET /api/connections/<name>/channels )
func GetConnectionChannelsJson(name string) ([]byte, error) {
	pnq := "connections/" + url.PathEscape(name) + "/channels"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetChannelsJson lists  all open channels.
// ( GET /api/channels )
func GetChannelsJson(page *pageFilter) ([]byte, error) {
	pnq := "channels" + page.ToUrlSuffix()
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetChannelJson lists details about an individual channel.
// ( GET /api/channels/<channel> )
func GetChannelJson(channel string) ([]byte, error) {
	pnq := "channels/" + url.PathEscape(channel)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetChannelsJson list all open channels in a specific virtual host
// ( GET /api/vhosts/<vhost>/channels )
func GetVhostChannelsJson(vhost string, page *pageFilter) ([]byte, error) {
	pnq := "vhosts/" + url.PathEscape(vhost) + "/channels" + page.ToUrlSuffix()
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetConsumersJson lists all consumers.
// ( GET /api/consumers )
func GetConsumersJson(page *pageFilter) ([]byte, error) {
	pnq := "consumers" + page.ToUrlSuffix()
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostConsumersJson lists all consumers in a given virtual host.
// ( GET /api/consumers/vhost )
func GetVHostConsumersJson(vhost string, page *pageFilter) ([]byte, error) {
	pnq := "consumers/" + url.PathEscape(vhost) + page.ToUrlSuffix()
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetExchangesJson returns A list of all exchanges
// ( GET /api/exchanges )
func GetExchangesJson(page *pageFilter) ([]byte, error) {
	pnq := "exchanges" + page.ToUrlSuffix()
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostExchangesJson returns a list of all exchanges in a given virtual host.
// ( GET /api/exchanges/<vhost>)
func GetVHostExchangesJson(vhost string) ([]byte, error) {
	pnq := "exchanges/" + url.PathEscape(vhost)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostExchanges returns a list of all exchanges in a given virtual host.
// ( GET /api/exchanges/<vhost>)
func GetVHostExchanges(name string) ([]*vhost.Exchange, error) {
	if data, err := GetVHostExchangesJson(name); err != nil {
		return nil, err
	} else {
		exchanges := []*vhost.Exchange{}
		err := json.Unmarshal(data, &exchanges)
		if err != nil {
			return nil, err
		}
		return exchanges, nil
	}
}

// GetVHostExchangeJson returns details of an individual exchange
// ( GET /api/exchanges/<vhost>)
func GetVHostExchangeJson(vhost, exchange string) ([]byte, error) {
	pnq := "exchanges/" + url.PathEscape(vhost) + "/" + url.PathEscape(exchange)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostExchangeBindingsJson returns a list of all bindings in which a given exchange is the source.
// ( GET /api/exchanges/<vhost>/bindings/source )
func GetVHostExchangeBindingsSourceJson(vhost, exchange string) ([]byte, error) {
	pnq := "exchanges/" + url.PathEscape(vhost) + "/" + url.PathEscape(exchange) + "/bindings/source"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostExchangeBindingsDestinationJson returns a list of all bindings in which a given exchange is the destination.
// ( GET /api/exchanges/<vhost>/bindings/destination )
func GetVHostExchangeBindingsDestinationJson(vhost, exchange string) ([]byte, error) {
	pnq := "exchanges/" + url.PathEscape(vhost) + "/" + url.PathEscape(exchange) + "/bindings/destination"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetQueuesJson() ([]byte, error) {
	resp, err := Get("queues")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetVHostQueuesJson(vhost string) ([]byte, error) {
	pnq := "queues/" + url.PathEscape(vhost)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetVHostQueueJson(vhost, queue string) ([]byte, error) {
	pnq := "queues/" + url.PathEscape(vhost) + "/" + url.PathEscape(queue)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetVHostQueueBindingsJson(vhost, queue string) ([]byte, error) {
	pnq := "queues/" + url.PathEscape(vhost) + "/" + url.PathEscape(queue) + "/bindings"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetVHostQueues(name string) ([]*vhost.Queue, error) {
	if data, err := GetVHostQueuesJson(name); err != nil {
		return nil, err
	} else {
		queues := []*vhost.Queue{}
		err := json.Unmarshal(data, &queues)
		if err != nil {
			return nil, err
		}
		return queues, nil
	}
}

// GetBindingsJson returns
// a list of all bindings
// ( GET /api/bindings  )
func GetBindingsJson() ([]byte, error) {
	resp, err := Get("bindings")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostBindingsJson returns
// a list of all bindings defined in  the vhost
// ( GET /api/bindings/<vhost> )
func GetVHostBindingsJson(vhost string) ([]byte, error) {
	pnq := "bindings/" + url.PathEscape(vhost)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostBindings returns
// a list of all bindings defined in  the vhost
// ( GET /api/bindings/<vhost> )
func GetVHostBindings(name string) ([]*vhost.Binding, error) {
	if data, err := GetVHostBindingsJson(name); err != nil {
		return nil, err
	} else {
		bindings := []*vhost.Binding{}
		err := json.Unmarshal(data, &bindings)
		if err != nil {
			return nil, err
		}
		return bindings, nil
	}
}

// GetVHostExchangeQueueBindingsJson returns
// a list of all bindings between an exchange and a queue
// ( GET /api/bindings/<vhost>/e/<exchange>/q/<queue> )
func GetVHostExchangeQueueBindingsJson(vhost, exchange, queue string) ([]byte, error) {
	pnq := "bindings/" +
		url.PathEscape(vhost) +
		"/e/" + url.PathEscape(exchange) +
		"/q/" + url.PathEscape(queue)

	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostExchangeQueueBindingJson returns
// the binding between an exchange and a queue. The props part of the URI is a "name"
// for the binding composed of its routing key and a hash of its arguments.
// props is the field named "properties_key" from a bindings listing response.
// ( GET /api/bindings/<vhost>/e/<exchange>/q/<queue>/<props> )
func GetVHostExchangeQueueBindingJson(vhost, exchange, queue, propsKey string) ([]byte, error) {
	pnq := "bindings/" +
		url.PathEscape(vhost) +
		"/e/" + url.PathEscape(exchange) +
		"/q/" + url.PathEscape(queue) +
		"/" + url.PathEscape(propsKey)

	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostExchangeExchangeBindingsJson returns
// a list of all bindings between an exchange and another exchange
// ( GET /api/bindings/<vhost>/e/<source>/e/<destination> )
func GetVHostExchangeExchangeBindingsJson(vhost, source, destination string) ([]byte, error) {
	pnq := "bindings/" +
		url.PathEscape(vhost) +
		"/e/" + url.PathEscape(source) +
		"/e/" + url.PathEscape(destination)

	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostExchangeExchangeBindingJson returns
// the binding between an exchange and another exchange with the given properties key
// ( GET /api/bindings/<vhost>/e/<source>/e/<destination>/<props> )
func GetVHostExchangeExchangeBindingJson(vhost, source, destination, propsKey string) ([]byte, error) {
	pnq := "bindings/" +
		url.PathEscape(vhost) +
		"/e/" + url.PathEscape(source) +
		"/e/" + url.PathEscape(destination) +
		"/" + url.PathEscape(propsKey)

	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetVHostsJson() ([]byte, error) {
	resp, err := Get("vhosts")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetVHostsLimitsJson() ([]byte, error) {
	resp, err := Get("vhost-limits")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetVHostLimitsJson(name string) ([]byte, error) {
	pnq := fmt.Sprintf("vhost-limits/%s", url.PathEscape(name))
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetVHostJson(name string) ([]byte, error) {
	pnq := "vhosts/" + url.PathEscape(name)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetVHostPermissionsJson(name string) ([]byte, error) {
	pnq := "vhosts/" + url.PathEscape(name) + "/permissions"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetVHostTopicPermissionsJson(name string) ([]byte, error) {
	pnq := "vhosts/" + url.PathEscape(name) + "/topic-permissions"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetAuthJson() ([]byte, error) {
	pnq := "auth"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetAuthAttemptsJson(node string) ([]byte, error) {
	pnq := fmt.Sprintf("auth/attempts/%s", url.PathEscape(node))
	if len(node) == 0 {
		pnq = strings.TrimRight(pnq, "/")
	}
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

func GetAuthAttemptsBySourceJson(node string) ([]byte, error) {
	pnq := fmt.Sprintf("auth/attempts/%s/source", url.PathEscape(node))
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetUsersJson returns the list of users ( GET /api/users )
func GetUsersJson() ([]byte, error) {
	resp, err := Get("users")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetUsersWithoutPermissionsJson returns a list of users that do not have access to any virtual host.
// ( GET /api/users/without-permissions )
func GetUsersWithoutPermissionsJson() ([]byte, error) {
	resp, err := Get("users/without-permissions")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetUserJson returns the user's details
// ( GET /api/users/<name> )
func GetUserJson(name string) ([]byte, error) {
	pnq := "users/" + url.PathEscape(name)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetUserPermissionsJson returns the user's permissions
// () GET /api/users/<name>/permissions )
func GetUserPermissionsJson(name string) ([]byte, error) {
	pnq := "users/" + url.PathEscape(name) + "/permissions"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetUserTopicPermissionsJson returns the user's topic permissions
// () GET /api/users/<name>/topic-permissions )
func GetUserTopicPermissionsJson(name string) ([]byte, error) {
	pnq := "users/" + url.PathEscape(name) + "/topic-permissions"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetUsersLimitsJson returns per-user limits for all users.
// () GET /api/user-limits )
func GetUsersLimitsJson() ([]byte, error) {
	pnq := "user-limits"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetUserLimitsJson returns per-user limits for a specific user.
// () GET /api/user-limits/<name> )
func GetUserLimitsJson(name string) ([]byte, error) {
	pnq := "user-limits/" + url.PathEscape(name)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetWhoAmIJson returns details of the currently authenticated user.
//  GET /api/whoami
func GetWhoAmIJson() ([]byte, error) {
	resp, err := Get("whoami")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetPermissionsJson returns a list of all permissions for all users.
//  GET /api/permissions
func GetPermissionsJson() ([]byte, error) {
	resp, err := Get("permissions")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostUserPermissionJson returns an individual permission of a user and virtual host
// ( GET /api/permissions/<vhost>/<user> )
func GetVHostUserPermissionJson(vhost, user string) ([]byte, error) {
	pnq := "permissions/" + url.PathEscape(vhost) + "/" + url.PathEscape(user)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetTopicPermissionsJson returns a list of all topic permissions for all users.
//  GET /api/topic-permissions
func GetTopicPermissionsJson() ([]byte, error) {
	resp, err := Get("topic-permissions")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostUserTopicPermissionsJson returns an individual permission of a user and virtual host
// ( GET /api/topic-permissions/<vhost>/<user> )
func GetVHostUserTopicPermissionsJson(vhost, user string) ([]byte, error) {
	pnq := "topic-permissions/" + url.PathEscape(vhost) + "/" + url.PathEscape(user)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetParametersJson returns a list of all vhost-scoped parameters.
// ( GET /api/parameters )
func GetParametersJson() ([]byte, error) {
	resp, err := Get("parameters")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetComponentParametersJson returns a list of all
// vhost-scoped parameters for a given component.
// ( GET /api/parameters/<component> )
func GetComponentParametersJson(component string) ([]byte, error) {
	pnq := "parameters/" + url.PathEscape(component)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetComponentVHostParametersJson returns a list of all vhost-scoped
// parameters for a given component and virtual host.
// ( GET /api/parameters/<component>/<vhost> )
func GetComponentVHostParametersJson(component, vhost string) ([]byte, error) {
	pnq := "parameters" +
		"/" + url.PathEscape(component) +
		"/" + url.PathEscape(vhost)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetComponentVHostParameterJson returns an individual vhost-scoped parameter.
// parameters for a given component and virtual host.
// ( GET /api/parameters/<component>/<vhost> )
func GetComponentVHostParameterJson(component, vhost, name string) ([]byte, error) {
	pnq := "parameters" +
		"/" + url.PathEscape(component) +
		"/" + url.PathEscape(vhost) +
		"/" + url.PathEscape(name)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetGlobalParametersJson returns a list of all global parameters.
// ( GET /api/parameters )
func GetGlobalParametersJson() ([]byte, error) {
	resp, err := Get("global-parameters")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetComponentParameters returns a list of all
// vhost-scoped parameters for a given component.
// ( GET /api/parameters/<component> )
func GetGlobalParameterJson(name string) ([]byte, error) {
	pnq := "global-parameters/" + url.PathEscape(name)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetPoliciesJson returns a list of all policies.
// ( GET /api/policies )
func GetPoliciesJson() ([]byte, error) {
	resp, err := Get("policies")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostPoliciesJson returns a list of all policies in a given vhost.
// ( GET /api/policies/<vhost> )
func GetVHostPoliciesJson(vhost string) ([]byte, error) {
	pnq := "policies/" + url.PathEscape(vhost)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostPolicyJson returns a policy in a given vhost.
// ( GET /api/policies/<vhost>/<name> )
func GetVHostPolicyJson(vhost, name string) ([]byte, error) {
	pnq := "policies" +
		"/" + url.PathEscape(vhost) +
		"/" + url.PathEscape(name)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetOperatorPoliciesJson returns a list of all operator-policies.
// ( GET /api/operator-policies )
func GetOperatorPoliciesJson() ([]byte, error) {
	resp, err := Get("operator-policies")
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostOperatorPoliciesJson returns a list of all operator-policies in a given vhost.
// ( GET /api/operator-policies/<vhost> )
func GetVHostOperatorPoliciesJson(vhost string) ([]byte, error) {
	pnq := "operator-policies/" + url.PathEscape(vhost)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetVHostOperatorPolicyJson returns a operator-policy in a given vhost.
// ( GET /api/operator-policies/<vhost>/<name> )
func GetVHostOperatorPolicyJson(vhost, name string) ([]byte, error) {
	pnq := "operator-policies" +
		"/" + url.PathEscape(vhost) +
		"/" + url.PathEscape(name)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetAlivenessTestJson declares a test queue on the target node, then publishes and consumes a message.
// Intended to be used as a very basic health check.
// Responds a 200 OK if the check succeeded, otherwise responds with a 503 Service Unavailable.
// ( GET /api/aliveness-test/<vhost> )
func GetAlivenessTestJson(vhost string) ([]byte, error) {
	pnq := "aliveness-test/" + url.PathEscape(vhost)
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetHealthChecksVHosts responds a 200 OK if all virtual hosts and running on the target node, otherwise responds with a 503 Service Unavailable.
// ( GET /api/health/checks/virtual-hosts )
func GetHealthChecksVHosts() ([]byte, error) {
	pnq := "health/checks/virtual-hosts"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetHealthChecksNodeIsMirrorSyncCritical Checks if there are classic mirrored queues without synchronised mirrors online
// (queues that would potentially lose data if the target node is shut down).
// Responds a 200 OK if there are no such classic mirrored queues, otherwise responds with a 503 Service Unavailable.
// ( GET /api/health/checks/node-is-mirror-sync-critical )
func GetHealthChecksNodeIsMirrorSyncCritical() ([]byte, error) {
	pnq := "health/checks/node-is-mirror-sync-critical"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetHealthChecksNodeIsQuorumCritical Checks if there are quorum queues with minimum online quorum
// (queues that would lose their quorum and availability if the target node is shut down).
// Responds a 200 OK if there are no such quorum queues, otherwise responds with a 503 Service Unavailable.
// ( GET /api/health/checks/node-is-quorum-critical )
func GetHealthChecksNodeIsQuorumCritical() ([]byte, error) {
	pnq := "health/checks/node-is-quorum-critical"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetHealthChecksAlarmsJson responds a 200 OK if there are no alarms in effect in the cluster, otherwise responds with a 503 Service Unavailable.
// ( GET /api/health/checks/alarms )
func GetHealthChecksAlarmsJson() ([]byte, error) {
	pnq := "health/checks/alarms"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetHealthChecksLocalAlarmsJson responds a 200 OK if there are no local alarms in effect on the target node, otherwise responds with a 503 Service Unavailable.
// ( GET /api/health/checks/local-alarms )
func GetHealthChecksLocalAlarmsJson() ([]byte, error) {
	pnq := "health/checks/local-alarms"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetHealthChecksCertificateExpirationJson checks the expiration date on the certificates for every listener configured to use TLS.
// Responds a 200 OK if all certificates are valid (have not expired), otherwise responds with a 503 Service Unavailable.
// Valid units: days, weeks, months, years. The value of the within argument is the number of units. So, when within is 2 and unit is "months", the expiration period used by the check will be the next two months.
// ( GET health/checks/certificate-expiration/<within>/<days|weeks|months|years> )
func GetHealthChecksCertificateExpirationJson(within int, unit string) ([]byte, error) {
	pnq := fmt.Sprintf("health/checks/certificate-expiration/%d/%s", within, url.PathEscape(unit))

	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetHealthChecksPortListenerJson Responds a 200 OK if there is an active listener on the give port, otherwise responds with a 503 Service Unavailable.
// ( GET /api/health/checks/port-listener/<port> )
func GetHealthChecksPortListenerJson(port uint16) ([]byte, error) {
	pnq := fmt.Sprintf("health/checks/port-listener/%d", port)

	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetHealthChecksProtocolListenerJson responds a 200 OK if there is an active listener for the given protocol,
// otherwise responds with a 503 Service Unavailable. Valid protocol names are: amqp091, amqp10, mqtt, stomp, web-mqtt, web-stomp
// ( GET /api/health/checks/protocol-listener/<protocol> )
func GetHealthChecksProtocolListenerJson(protocol string) ([]byte, error) {
	pnq := fmt.Sprintf("health/checks/protocol-listener/%s", protocol)

	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

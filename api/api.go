package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/LogiqsAgro/rmq/api/vhost"
)

const (
	UnitDays   = "days"
	UnitWeeks  = "weeks"
	UnitMonths = "months"
	UnitYears  = "years"
)

func AllUnits() []string {
	return []string{
		UnitDays,
		UnitWeeks,
		UnitMonths,
		UnitYears,
	}
}

// GetOverviewJson returns tarious random bits of information that describe the whole system. ( GET /api/overview )
func GetOverviewJson() ([]byte, error) {
	resp, err := Get("overview")
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
			UrlSuffix()

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

// GetUsersJson returns the result of GET /api/users
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
// () GET /api/users/<name> )
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

// GetHealthChecksAlarms responds a 200 OK if there are no alarms in effect in the cluster, otherwise responds with a 503 Service Unavailable.
// ( GET /api/health/checks/alarms )
func GetHealthChecksLocalAlarmsJson() ([]byte, error) {
	pnq := "health/checks/local-alarms"
	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

// GetHealthChecksAlarms responds a 200 OK if there are no alarms in effect in the cluster, otherwise responds with a 503 Service Unavailable.
// ( GET /api/health/checks/alarms )
func GetHealthChecksCertificateExpirationJson(within int, unit string) ([]byte, error) {
	pnq := "health/checks/certificate-expiration" +
		"/" + fmt.Sprintf("%d", within) +
		"/" + url.PathEscape(unit)

	resp, err := Get(pnq)
	if err != nil {
		return nil, err
	}

	return ReadBody(resp)
}

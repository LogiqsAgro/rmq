/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"strings"

	"github.com/LogiqsAgro/rmq/api"
	"github.com/LogiqsAgro/rmq/api/vhost"
	"github.com/awalterschulze/gographviz"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates an image from the rabbit mq topology",
	Long:  ` `,
	Run:   generate,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	api.AddConfigFlags(generateCmd)

}

func generate(cmd *cobra.Command, args []string) {
	definition, err := api.GetVHostDefinitions(api.Config.VHost)
	api.PanicIf(err)

	g, err := buildGraph(definition)
	panicIf(err)
	fmt.Println("// RabbitMQ version: ", definition.RabbitVersion)
	fmt.Println("// node host: ", api.Config.Host)
	fmt.Println("// vhost: ", api.Config.VHost)
	fmt.Println(g.String())
}

func quoted(s string) string {
	s = strings.Replace(s, ":", "\\n:", -1)
	return fmt.Sprintf("\"%s\"", s)
}

func buildGraph(definition *vhost.Definition) (*gographviz.Graph, error) {
	g := gographviz.NewGraph()
	graphName := "RabbitMQ"
	g.SetName(graphName)
	g.SetDir(true)

	qId := 0
	nextQID := func() string {
		qId += 1
		return fmt.Sprintf("Q%03d", qId)
	}

	queues := make(map[string]string)
	for i := 0; i < len(definition.Queues); i++ {
		q := definition.Queues[i]
		qid := nextQID()
		queues[q.Name] = qid
		g.AddNode(graphName, qid, map[string]string{
			"label": quoted(q.Name),
			"shape": "box",
			"color": "blue",
		})
		// fmt.Println("Q", q.Name)
	}

	eId := 0
	nextEID := func() string {
		eId += 1
		return fmt.Sprintf("E%03d", eId)
	}

	exchanges := make(map[string]string)
	for i := 0; i < len(definition.Exchanges); i++ {
		e := definition.Exchanges[i]
		eid := nextEID()
		exchanges[e.Name] = eid
		g.AddNode(graphName, eid, map[string]string{
			"label": quoted(e.Name + "\\ntype = " + e.Type),
			"shape": "octagon",
			"color": "gray",
		})
		// fmt.Println("X", e.Name)
	}

	for i := 0; i < len(definition.Bindings); i++ {
		b := definition.Bindings[i]
		//fmt.Println("B", b.Source, "->", b.Destination, "(", b.DestinationType, ")")

		source, ok := exchanges[b.Source]
		if !ok {
			return nil, fmt.Errorf("could not find source exchange %s for binding", b.Source)
		}

		if b.DestinationType == "queue" {
			target, ok := queues[b.Destination]
			if !ok {
				return nil, fmt.Errorf("could not find destination %s %s for binding", b.DestinationType, b.Destination)
			}
			g.AddEdge(source, target, true, map[string]string{
				// "label": quoted("x to q"),
				"color": "blue",
			})
		} else {
			target, ok := exchanges[b.Destination]
			if !ok {
				return nil, fmt.Errorf("could not find destination %s %s for binding", b.DestinationType, b.Destination)
			}
			g.AddEdge(source, target, true, map[string]string{
				// "label": quoted("x to x"),
				"color": "black",
			})
		}
	}

	return g, nil
}

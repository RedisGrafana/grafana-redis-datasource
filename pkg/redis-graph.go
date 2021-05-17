package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
)

/**
 * GRAPH.QUERY <Graph name> {query}
 *
 * Executes the given query against a specified graph.
 * @see https://oss.redislabs.com/redisgraph/commands/#graphquery
 */
func queryGraphQuery(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	var result []interface{}

	// Run command
	err := client.RunFlatCmd(&result, qm.Command, qm.Key, qm.Cypher)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame for nodes
	frameWithNodes := data.NewFrame("nodes")
	frameWithNodes.Meta = &data.FrameMeta{
		PreferredVisualization: "nodeGraph",
	}
	frameWithNodes.Fields = append(frameWithNodes.Fields, data.NewField("id", nil, []string{}))
	frameWithNodes.Fields = append(frameWithNodes.Fields, data.NewField("title", nil, []string{}))
	frameWithNodes.Fields = append(frameWithNodes.Fields, data.NewField("subTitle", nil, []string{}))
	frameWithNodes.Fields = append(frameWithNodes.Fields, data.NewField("mainStat", nil, []string{}))
	frameWithNodes.Fields = append(frameWithNodes.Fields, data.NewField("arc__", nil, []int64{}))

	// New Frame for edges
	frameWithEdges := data.NewFrame("edges")
	frameWithEdges.Meta = &data.FrameMeta{
		PreferredVisualization: "nodeGraph",
	}
	frameWithEdges.Fields = append(frameWithEdges.Fields, data.NewField("id", nil, []string{}))
	frameWithEdges.Fields = append(frameWithEdges.Fields, data.NewField("source", nil, []string{}))
	frameWithEdges.Fields = append(frameWithEdges.Fields, data.NewField("target", nil, []string{}))
	frameWithEdges.Fields = append(frameWithEdges.Fields, data.NewField("mainStat", nil, []string{}))

	// New Frame for data
	frameWithData := data.NewFrame("data")

	// Parse entries
	existingNodes := map[string]bool{}
	for _, entries := range result[1].([]interface{}) {
		// Parse results
		nodes, edges, dataEntries := findAllNodesAndEdges(entries)

		// Add Nodes
		for _, node := range nodes {
			// Add each nodeEntry only once
			if _, ok := existingNodes[node.Id]; !ok {
				frameWithNodes.AppendRow(node.Id, node.Title, node.SubTitle, node.MainStat, node.Arc)
				existingNodes[node.Id] = true
			}
		}

		// Add Edges
		for _, edge := range edges {
			frameWithEdges.AppendRow(edge.Id, edge.Source, edge.Target, edge.MainStat)
		}

		// Skip empty data lines
		entries := dataEntries.([]interface{})
		if len(entries) == 0 {
			continue
		}

		// Add fields based on the first row to know types
		if len(frameWithData.Fields) == 0 {
			// Parse data records header
			for i, record := range result[0].([]interface{}) {
				field := string(record.([]byte))

				switch entries[i].(type) {
				case int64:
					frameWithData.Fields = append(frameWithData.Fields, data.NewField(field, nil, []int64{}))
				case float64:
					frameWithData.Fields = append(frameWithData.Fields, data.NewField(field, nil, []float64{}))
				default:
					frameWithData.Fields = append(frameWithData.Fields, data.NewField(field, nil, []string{}))
				}
			}
		}

		frameWithData.AppendRow(dataEntries.([]interface{})...)
	}

	// Add Frames with Nodes if found
	if frameWithNodes.Rows() > 0 {
		response.Frames = append(response.Frames, frameWithNodes)
	}

	// Add Frames with Edges if found
	if frameWithEdges.Rows() > 0 {
		response.Frames = append(response.Frames, frameWithEdges)
	}

	// Add Frames with Data if found
	if frameWithData.Rows() > 0 {
		response.Frames = append(response.Frames, frameWithData)
	}

	// New Frame for metadata
	frameWithMeta := data.NewFrame("metadata")
	frameWithMeta.Fields = append(frameWithMeta.Fields, data.NewField("data", nil, []string{}))
	frameWithMeta.Fields = append(frameWithMeta.Fields, data.NewField("value", nil, []string{}))

	// Parse meta
	for _, meta := range result[2].([]interface{}) {
		switch m := meta.(type) {
		case []byte:
			fields := strings.Split(string(m), ":")
			if len(fields) != 2 {
				continue
			}

			frameWithMeta.AppendRow(fields[0], strings.TrimSpace(fields[1]))
		}
	}

	// Add Frames with Metadata if found
	if frameWithMeta.Rows() > 0 {
		response.Frames = append(response.Frames, frameWithMeta)
	}

	return response
}

/**
 * Parse array of entries and find
 * either Nodes https://oss.redislabs.com/redisgraph/result_structure/#nodes
 * or Relations https://oss.redislabs.com/redisgraph/result_structure/#relations
 * and create corresponding nodeEntry or edgeEntry
 **/
func findAllNodesAndEdges(input interface{}) ([]models.NodeEntry, []models.EdgeEntry, interface{}) {
	nodes := []models.NodeEntry{}
	edges := []models.EdgeEntry{}
	dataEntries := []interface{}{}

	// Parse entries
	for _, entry := range input.([]interface{}) {
		switch e := entry.(type) {
		case []interface{}:
			props := []string{}

			// Node https://oss.redislabs.com/redisgraph/result_structure/#nodes
			if len(e) == 3 {
				node := models.NodeEntry{Arc: 1}

				// Id
				idArray := e[0].([]interface{})
				node.Id = strconv.FormatInt(idArray[1].(int64), 10)

				// Labels
				labelsArray := e[1].([]interface{})
				labels := labelsArray[1].([]interface{})

				// Assume first label will be a title if exists
				if len(labels) > 0 {
					node.Title = string(labels[0].([]byte))
				}

				// Properties
				propertiesArray := e[2].([]interface{})
				properties := propertiesArray[1].([]interface{})

				// Assume first property will be a mainStat if exists
				for _, prop := range properties {
					propertyArray := prop.([]interface{})
					value := ""

					// Get value
					switch propValue := propertyArray[1].(type) {
					case []byte:
						value = string(propValue)
					case int64:
						value = strconv.FormatInt(propValue, 10)
					}

					// Set MainStat
					if node.MainStat == "" {
						node.MainStat = value
					}

					// Add property
					propString := fmt.Sprintf("\"%s\"=\"%s\"", propertyArray[0], value)
					props = append(props, propString)
				}

				// Add data entry and node
				dataEntries = append(dataEntries, strings.Join(props, ", "))
				nodes = append(nodes, node)
			}

			// Relation https://oss.redislabs.com/redisgraph/result_structure/#relations
			if len(e) == 5 {
				edge := models.EdgeEntry{}

				// Id
				idArray := e[0].([]interface{})
				edge.Id = strconv.FormatInt(idArray[1].(int64), 10)

				// Main Stat
				typeArray := e[1].([]interface{})
				edge.MainStat = string(typeArray[1].([]byte))

				// Source Id
				srcArray := e[2].([]interface{})
				edge.Source = strconv.FormatInt(srcArray[1].(int64), 10)

				// Target Id
				destArray := e[3].([]interface{})
				edge.Target = strconv.FormatInt(destArray[1].(int64), 10)

				// Properties
				propertiesArray := e[4].([]interface{})
				properties := propertiesArray[1].([]interface{})
				for _, prop := range properties {
					propertyArray := prop.([]interface{})
					value := ""

					// Get value
					switch propValue := propertyArray[1].(type) {
					case []byte:
						value = string(propValue)
					case int64:
						value = strconv.FormatInt(propValue, 10)
					}

					// Add property
					propString := fmt.Sprintf("\"%s\"=\"%s\"", propertyArray[0], value)
					props = append(props, propString)
				}

				// Add data entry and edge
				dataEntries = append(dataEntries, strings.Join(props, ", "))
				edges = append(edges, edge)
			}
		case []byte:
			dataEntries = append(dataEntries, string(e))
		default:
			dataEntries = append(dataEntries, e)
		}
	}

	// Return
	return nodes, edges, dataEntries
}

/**
 * GRAPH.SLOWLOG <Graph name>
 *
 * Returns a list containing up to 10 of the slowest queries issued against the given graph ID.
 * @see https://oss.redislabs.com/redisgraph/commands/#graphslowlog
 */
func queryGraphSlowlog(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	var result [][]string

	// Run command
	err := client.RunFlatCmd(&result, qm.Command, qm.Key)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Command)
	frame.Fields = append(frame.Fields, data.NewField("timestamp", nil, []time.Time{}))
	frame.Fields = append(frame.Fields, data.NewField("command", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("query", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("duration", nil, []float64{}))
	response.Frames = append(response.Frames, frame)

	// Set Field Config
	frame.Fields[3].Config = &data.FieldConfig{Unit: "µs"}

	// Entries
	for _, entry := range result {
		timestamp, _ := strconv.ParseInt(entry[0], 10, 64)
		duration, _ := strconv.ParseFloat(entry[3], 64)
		frame.AppendRow(time.Unix(timestamp, 0), entry[1], entry[2], duration)
	}

	// Return
	return response
}

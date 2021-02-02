package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * Graph data for GRAPH.QUERY Radix marshaling
 */
// type graphData struct {
//	ID         string                 `redis:"id"`
//	Labels     map[string]interface{} `redis:"labels"`
//	Properties map[string]interface{} `redis:"properties"`
//}

/**
 * GRAPH.QUERY <Graph name> {query}
 *
 * Executes the given query against a specified graph.
 * @see https://oss.redislabs.com/redisgraph/commands/#graphquery
 */
func queryGraphQuery(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	var result interface{}

	// Run command
	err := client.RunFlatCmd(&result, "GRAPH.QUERY", qm.Key, qm.Cypher)

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

	// Adding frames to response
	response.Frames = append(response.Frames, frameWithNodes)
	response.Frames = append(response.Frames, frameWithEdges)

	// Data
	frameWithNodes.AppendRow("5e5aeb4a820126475cb09ccb", "Writer", "", "George R. R. Martin", int64(1))
	frameWithNodes.AppendRow("5e5aeb4a820126475cb09ccd", "Writer", "", "Linda Antonsson", int64(1))
	frameWithNodes.AppendRow("5e5aeb4a820126475cb09ccc", "Writer", "", "Elio Garcia", int64(1))
	frameWithNodes.AppendRow("5e5aeadb820126475cb09cbf", "Book", "", "A Game of Thrones", int64(1))
	frameWithNodes.AppendRow("5e5aeadb820126475cb09cc0", "Book", "", "A Clash of Kings", int64(1))
	frameWithNodes.AppendRow("5e5aeadb820126475cb09cc9", "Book", "", "The World of Ice and Fire", int64(1))
	frameWithNodes.AppendRow("5e5aeadb820126475cb09cc2", "Book", "", "The Hedge Knight", int64(1))
	frameWithNodes.AppendRow("5e5aeadb820126475cb09cc3", "Book", "", "A Feast for Crows", int64(1))

	frameWithEdges.AppendRow("1", "5e5aeb4a820126475cb09ccb", "5e5aeadb820126475cb09cbf", "wrote")
	frameWithEdges.AppendRow("2", "5e5aeb4a820126475cb09ccb", "5e5aeadb820126475cb09cc0", "wrote")
	frameWithEdges.AppendRow("3", "5e5aeb4a820126475cb09ccb", "5e5aeadb820126475cb09cc9", "wrote")
	frameWithEdges.AppendRow("3", "5e5aeb4a820126475cb09ccb", "5e5aeadb820126475cb09cc2", "wrote")
	frameWithEdges.AppendRow("4", "5e5aeb4a820126475cb09ccd", "5e5aeadb820126475cb09cc9", "wrote")
	frameWithEdges.AppendRow("5", "5e5aeb4a820126475cb09ccc", "5e5aeadb820126475cb09cc9", "wrote")
	frameWithEdges.AppendRow("5", "5e5aeb4a820126475cb09ccb", "5e5aeadb820126475cb09cc3", "wrote")

	return response
}

package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * XINFO [CONSUMERS key groupname] [GROUPS key] [STREAM key] [HELP]
 *
 * @see https://redis.io/commands/xinfo
 */
func queryXInfoStream(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result map[string]string
	err := client.RunFlatCmd(&result, "XINFO", "STREAM", qm.Key)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	fields := []string{}
	values := []string{}

	// Add fields and values
	for k := range result {
		fields = append(fields, k)
		values = append(values, result[k])
	}

	// New Frame
	frame := data.NewFrame(qm.Key,
		data.NewField("Field", nil, fields),
		data.NewField("Value", nil, values))

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

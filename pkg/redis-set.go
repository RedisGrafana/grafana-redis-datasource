package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * SMEMBERS key
 *
 * @see https://redis.io/commands/smembers
 */
func querySMembers(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var values []string
	err := client.RunFlatCmd(&values, qm.Command, qm.Key)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Key,
		data.NewField("Value", nil, values))

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

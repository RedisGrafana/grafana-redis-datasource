package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
)

/**
 * SMEMBERS key
 *
 * @see https://redis.io/commands/smembers
 */
func (ds *redisDatasource) querySMembers(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var values []string
	err := client.Do(radix.FlatCmd(&values, qm.Command, qm.Key))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Key,
		data.NewField("Value", nil, values))

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
)

/**
 * HGETALL key
 *
 * @see https://redis.io/commands/hgetall
 */
func (ds *redisDatasource) queryHGetAll(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result []string
	err := client.Do(radix.FlatCmd(&result, qm.Command, qm.Key))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	fields := []string{}
	values := []string{}

	// Add fields and values
	for i := 0; i < len(result); i += 2 {
		fields = append(fields, result[i])
		values = append(values, result[i+1])
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

/**
 * HGET key field
 *
 * @see https://redis.io/commands/hget
 */
func (ds *redisDatasource) queryHGet(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var value string
	err := client.Do(radix.FlatCmd(&value, qm.Command, qm.Key, qm.Field))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, ds.createFrameValue(qm.Key, value))

	// Return
	return response
}

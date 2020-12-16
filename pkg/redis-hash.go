package main

import (
	"fmt"
	"strconv"

	"bitbucket.org/creachadair/shell"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
)

/**
 * HGETALL key
 *
 * @see https://redis.io/commands/hgetall
 */
func (ds *redisDatasource) queryHGetAll(qm queryModel, client ClientInterface) backend.DataResponse {
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
func (ds *redisDatasource) queryHGet(qm queryModel, client ClientInterface) backend.DataResponse {
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

/**
 * HMGET key field [field ...]
 *
 * @see https://redis.io/commands/hmget
 */
func (ds *redisDatasource) queryHMGet(qm queryModel, client ClientInterface) backend.DataResponse {
	response := backend.DataResponse{}

	// Split Field to array
	fields, ok := shell.Split(qm.Field)

	// Check if filter is valid
	if !ok {
		response.Error = fmt.Errorf("Field is not valid")
		return response
	}

	// Execute command
	var result []string
	err := client.Do(radix.FlatCmd(&result, qm.Command, qm.Key, fields))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Command)

	// Parse results and add fields
	for i, value := range result {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			frame.Fields = append(frame.Fields, data.NewField(fields[i], nil, []float64{floatValue}))
		} else {
			frame.Fields = append(frame.Fields, data.NewField(fields[i], nil, []string{value}))
		}
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

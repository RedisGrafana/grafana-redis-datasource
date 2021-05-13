package main

import (
	"fmt"
	"strconv"

	"bitbucket.org/creachadair/shell"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * HGETALL key
 *
 * @see https://redis.io/commands/hgetall
 */
func queryHGetAll(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result []string
	err := client.RunFlatCmd(&result, qm.Command, qm.Key)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Command)

	// Add fields and values
	for i := 0; i < len(result); i += 2 {
		if floatValue, err := strconv.ParseFloat(result[i+1], 64); err == nil {
			frame.Fields = append(frame.Fields, data.NewField(result[i], nil, []float64{floatValue}))
		} else {
			frame.Fields = append(frame.Fields, data.NewField(result[i], nil, []string{result[i+1]}))
		}
	}

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
func queryHGet(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var value string
	err := client.RunFlatCmd(&value, qm.Command, qm.Key, qm.Field)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, createFrameValue(qm.Field, value, qm.Field))

	// Return
	return response
}

/**
 * HMGET key field [field ...]
 *
 * @see https://redis.io/commands/hmget
 */
func queryHMGet(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Split Field to array
	fields, ok := shell.Split(qm.Field)

	// Check if filter is valid
	if !ok {
		response.Error = fmt.Errorf("field is not valid")
		return response
	}

	// Execute command
	var result []string
	err := client.RunFlatCmd(&result, qm.Command, qm.Key, fields)

	// Check error
	if err != nil {
		return errorHandler(response, err)
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

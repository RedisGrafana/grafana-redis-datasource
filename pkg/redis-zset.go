package main

import (
	"strconv"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * ZRANGE key min max [BYSCORE|BYLEX] [REV] [LIMIT offset count] [WITHSCORES]
 *
 * @see https://redis.io/commands/zrange
 */
func queryZRange(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result []string
	var err error

	if qm.ZRangeQuery == "" {
		err = client.RunFlatCmd(&result, qm.Command, qm.Key, qm.Min, qm.Max, "WITHSCORES")
	} else {
		err = client.RunFlatCmd(&result, qm.Command, qm.Key, qm.Min, qm.Max, qm.ZRangeQuery, "WITHSCORES")
	}

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Command)

	// Add fields and scores
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

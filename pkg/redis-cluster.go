package main

import (
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
)

/**
 * CLUSTER INFO
 *
 * @see https://redis.io/commands/cluster-info
 */
func (ds *redisDatasource) queryClusterInfo(qm queryModel, client ClientInterface) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result string
	err := client.Do(radix.Cmd(&result, "CLUSTER", "INFO"))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// Split lines
	lines := strings.Split(strings.Replace(result, "\r\n", "\n", -1), "\n")

	// New Frame
	frame := data.NewFrame(qm.Command)

	// Parse lines
	for _, line := range lines {
		fields := strings.Split(line, ":")

		if len(fields) < 2 {
			continue
		}

		// Add Field
		if floatValue, err := strconv.ParseFloat(fields[1], 64); err == nil {
			frame.Fields = append(frame.Fields, data.NewField(fields[0], nil, []float64{floatValue}))
		} else {
			frame.Fields = append(frame.Fields, data.NewField(fields[0], nil, []string{fields[1]}))
		}
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

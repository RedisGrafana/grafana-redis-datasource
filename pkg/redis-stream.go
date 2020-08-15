package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
)

/**
 * XINFO [CONSUMERS key groupname] [GROUPS key] [STREAM key] [HELP]
 *
 * @see https://redis.io/commands/xinfo
 */
func (ds *redisDatasource) queryXInfoStream(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result []string
	err := client.Do(radix.FlatCmd(&result, "XINFO", "STREAM", qm.Key))

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

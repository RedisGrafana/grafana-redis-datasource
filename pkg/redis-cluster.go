package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * CLUSTER INFO
 *
 * @see https://redis.io/commands/cluster-info
 */
func queryClusterInfo(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result string
	err := client.RunCmd(&result, "CLUSTER", "INFO")

	// Check error
	if err != nil {
		return errorHandler(response, err)
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

/**
 * CLUSTER NODES
 *
 * @see https://redis.io/commands/cluster-nodes
 */
func queryClusterNodes(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result string
	err := client.RunCmd(&result, "CLUSTER", "NODES")

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Split lines
	lines := strings.Split(strings.Replace(result, "\r\n", "\n", -1), "\n")

	// New Frame
	frame := data.NewFrame(qm.Command,
		data.NewField("Id", nil, []string{}),
		data.NewField("Address", nil, []string{}),
		data.NewField("Flags", nil, []string{}),
		data.NewField("Master", nil, []string{}),
		data.NewField("Ping", nil, []int64{}),
		data.NewField("Pong", nil, []int64{}),
		data.NewField("Epoch", nil, []int64{}),
		data.NewField("State", nil, []string{}),
		data.NewField("Slot", nil, []string{}))

	// Set Field Config
	frame.Fields[4].Config = &data.FieldConfig{Unit: "ms"}
	frame.Fields[5].Config = &data.FieldConfig{Unit: "ms"}

	// Parse lines
	for _, line := range lines {
		fields := strings.Split(line, " ")

		// Check number of fields
		if len(fields) < 8 {
			continue
		}

		var ping int64
		var pong int64
		var epoch int64
		slot := ""

		// Parse values
		ping, _ = strconv.ParseInt(fields[4], 10, 64)
		pong, _ = strconv.ParseInt(fields[5], 10, 64)
		epoch, _ = strconv.ParseInt(fields[6], 10, 64)

		// Check Ping and convert 0
		if ping == 0 {
			ping = time.Now().UnixNano() / 1e6
		}

		// Check Pong and convert 0
		if pong == 0 {
			pong = time.Now().UnixNano() / 1e6
		}

		// Add slots which is missing for slaves
		if len(fields) > 8 {
			slot = fields[8]
		}

		// Add Query
		frame.AppendRow(fields[0], fields[1], fields[2], fields[3], ping, pong, epoch, fields[7], slot)
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

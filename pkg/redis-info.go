package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * INFO [section]
 *
 * @see https://redis.io/commands/info
 */
func queryInfo(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result string
	err := client.RunCmd(&result, qm.Command, qm.Section)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Split lines
	lines := strings.Split(strings.Replace(result, "\r\n", "\n", -1), "\n")

	// New Frame
	frame := data.NewFrame(qm.Command)

	// Command stats
	if qm.Section == "commandstats" {
		frame.Fields = append(frame.Fields, data.NewField("Command", nil, []string{}),
			data.NewField("Calls", nil, []float64{}),
			data.NewField("Usec", nil, []float64{}).SetConfig(&data.FieldConfig{Unit: "µs"}),
			data.NewField("Usec_per_call", nil, []float64{}).SetConfig(&data.FieldConfig{Unit: "µs"}),
			data.NewField("RejectedCalls", nil, []float64{}),
			data.NewField("FailedCalls", nil, []float64{}),
			data.NewField("CallsMaster", nil, []float64{}),
		)

		// Parse lines
		for _, line := range lines {
			fields := strings.Split(line, ":")

			if len(fields) < 2 {
				continue
			}

			// Stats
			stats := strings.Split(fields[1], ",")
			values := map[string]float64{}

			for _, stat := range stats {
				value := strings.Split(stat, "=")
				values[value[0]], _ = strconv.ParseFloat(value[1], 64)
			}

			// Command name
			cmd := strings.Replace(fields[0], "cmdstat_", "", 1)

			// Add Command
			frame.AppendRow(cmd, values["calls"], values["usec"], values["usec_per_call"], values["rejected_calls"], values["failed_calls"], values["calls_master"])
		}

		// Add the frames to the response
		response.Frames = append(response.Frames, frame)

		// Return
		return response
	}

	// Error stats ( added in Redis >= v6.2 )
	if qm.Section == "errorstats" {
		// Not Streaming
		if !qm.Streaming {
			frame.Fields = append(frame.Fields,
				data.NewField("Error", nil, []string{}),
				data.NewField("Count", nil, []int64{}))
		}

		// Parse lines
		for _, line := range lines {
			fields := strings.Split(line, ":")

			if len(fields) < 2 {
				continue
			}

			// Parse Error Stats
			count := strings.Split(fields[1], "=")
			var errorValue int64

			// Parse Error
			if len(count) == 2 {
				errorValue, _ = strconv.ParseInt(count[1], 10, 64)
			}

			// Error prefix
			error := strings.Replace(fields[0], "errorstat_", "", 1)

			// Streaming
			if qm.Streaming {
				frame.Fields = append(frame.Fields, data.NewField(error, nil, []int64{errorValue}))
			} else {
				frame.AppendRow(error, errorValue)
			}
		}

		// Add the frames to the response
		response.Frames = append(response.Frames, frame)

		// Return
		return response
	}

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
 * CLIENT LIST [TYPE normal|master|replica|pubsub]
 *
 * @see https://redis.io/commands/client-list
 */
func queryClientList(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result string
	err := client.RunCmd(&result, "CLIENT", "LIST")

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Split lines
	lines := strings.Split(strings.Replace(result, "\r\n", "\n", -1), "\n")

	// New Frame
	frame := data.NewFrame(qm.Command)

	// Parse lines
	for i, line := range lines {
		var values []interface{}

		// Split line to array
		fields := strings.Fields(line)

		// Parse lines
		for _, field := range fields {
			// Split properties
			value := strings.Split(field, "=")

			// Skip if less than 2 elements
			if len(value) < 2 {
				continue
			}

			// Add Header for first row
			if i == 0 {
				if _, err := strconv.ParseInt(value[1], 10, 64); err == nil {
					frame.Fields = append(frame.Fields, data.NewField(value[0], nil, []int64{}))
				} else {
					frame.Fields = append(frame.Fields, data.NewField(value[0], nil, []string{}))
				}
			}

			// Add Int64 or String value
			if intValue, err := strconv.ParseInt(value[1], 10, 64); err == nil {
				values = append(values, intValue)
			} else {
				values = append(values, value[1])
			}
		}

		// Add Row
		frame.AppendRow(values...)
	}

	// Add the frame to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

/**
 * SLOWLOG subcommand [argument]
 *
 * @see https://redis.io/commands/slowlog
 */
func querySlowlogGet(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result interface{}
	var err error

	if qm.Size > 0 {
		err = client.RunFlatCmd(&result, "SLOWLOG", "GET", qm.Size)
	} else {
		err = client.RunCmd(&result, "SLOWLOG", "GET")
	}

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Command,
		data.NewField("Id", nil, []int64{}),
		data.NewField("Timestamp", nil, []time.Time{}),
		data.NewField("Duration", nil, []int64{}),
		data.NewField("Command", nil, []string{}))

	// Set Field Config
	frame.Fields[2].Config = &data.FieldConfig{Unit: "µs"}

	// Parse Time-Series data
	for _, innerArray := range result.([]interface{}) {
		query := innerArray.([]interface{})
		command := ""

		/**
		 * Redis OSS has arguments as forth element of array
		 * Redis Enterprise has arguments as fifth
		 * Redis prior to 4.0 has only 4 fields.
		 */
		argumentsID := 3
		if len(query) > 4 {
			switch query[4].(type) {
			case []interface{}:
				argumentsID = 4
			default:
			}
		}

		/**
		 * Merge all arguments
		 */
		for _, arg := range query[argumentsID].([]interface{}) {

			// Add space between command and arguments
			if command != "" {
				command += " "
			}

			// Combine args into single command
			switch arg := arg.(type) {
			case int64:
				command += strconv.FormatInt(arg, 10)
			case []byte:
				command += string(arg)
			case string:
				command += arg
			default:
				log.DefaultLogger.Debug("Slowlog", "default", arg)
			}
		}

		// Add Query
		frame.AppendRow(query[0].(int64), time.Unix(query[1].(int64), 0), query[2].(int64), command)
	}

	// Add the frame to the response
	response.Frames = append(response.Frames, frame)

	// Return Response
	return response
}

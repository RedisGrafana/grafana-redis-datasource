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
	lines := splitInfoLines(result)
	section := strings.ToLower(strings.TrimSpace(qm.Section))

	/**
	 * INFO all / INFO everything return every section at once. Flattening them
	 * into a single frame loses the record sections (commandstats, latencystats,
	 * keysizes, ...), so each block is parsed on its own and returned as its own
	 * frame named after the section.
	 */
	if section == sectionAll || section == sectionEverything {
		for _, block := range splitInfoSections(lines) {
			if len(block.Lines) == 0 {
				continue
			}

			response.Frames = append(response.Frames, infoSectionFrame(block.Name, block.Name, block.Lines, qm.Streaming))
		}

		return response
	}

	// Single section keeps the historical frame name
	response.Frames = append(response.Frames, infoSectionFrame(qm.Command, section, lines, qm.Streaming))

	// Return
	return response
}

/**
 * INFO commandstats
 *
 * cmdstat_get:calls=326,usec=1406,usec_per_call=4.31,rejected_calls=0,failed_calls=0
 */
func infoCommandstatsFrame(name string, lines []string) *data.Frame {
	frame := data.NewFrame(name,
		data.NewField("Command", nil, []string{}),
		data.NewField("Calls", nil, []float64{}),
		data.NewField("Usec", nil, []float64{}).SetConfig(&data.FieldConfig{Unit: "µs"}),
		data.NewField("Usec_per_call", nil, []float64{}).SetConfig(&data.FieldConfig{Unit: "µs"}),
		data.NewField("RejectedCalls", nil, []float64{}),
		data.NewField("FailedCalls", nil, []float64{}),
		data.NewField("CallsMaster", nil, []float64{}),
	)

	// Parse lines
	for _, line := range lines {
		key, value, ok := splitInfoLine(line)
		if !ok {
			continue
		}

		// Stats
		values := map[string]float64{}
		for _, pair := range parseInfoPairs(value) {
			values[pair.Key], _ = strconv.ParseFloat(pair.Value, 64)
		}

		// Command name
		cmd := strings.Replace(key, "cmdstat_", "", 1)

		// Add Command
		frame.AppendRow(cmd, values["calls"], values["usec"], values["usec_per_call"], values["rejected_calls"], values["failed_calls"], values["calls_master"])
	}

	return frame
}

/**
 * INFO errorstats ( Redis >= 6.2 )
 *
 * errorstat_ERR:count=2850
 */
func infoErrorstatsFrame(name string, lines []string, streaming bool) *data.Frame {
	frame := data.NewFrame(name)

	// Not Streaming
	if !streaming {
		frame.Fields = append(frame.Fields,
			data.NewField("Error", nil, []string{}),
			data.NewField("Count", nil, []int64{}))
	}

	// Parse lines
	for _, line := range lines {
		key, value, ok := splitInfoLine(line)
		if !ok {
			continue
		}

		// Parse Error Stats
		var errorValue int64
		for _, pair := range parseInfoPairs(value) {
			if pair.Key == "count" {
				errorValue, _ = strconv.ParseInt(pair.Value, 10, 64)
			}
		}

		// Error prefix
		errorName := strings.Replace(key, "errorstat_", "", 1)

		// Streaming
		if streaming {
			frame.Fields = append(frame.Fields, data.NewField(errorName, nil, []int64{errorValue}))
		} else {
			frame.AppendRow(errorName, errorValue)
		}
	}

	return frame
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
			value := strings.SplitN(field, "=", 2)

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

// slowlogTruncationMarker is appended by Redis when an entry has more than
// SLOWLOG_ENTRY_MAX_ARGC arguments.
const slowlogTruncationMarker = "more arguments)"

/**
 * SLOWLOG GET [count]
 *
 * Entry layout, by server version:
 *   Redis < 4.0        [id, timestamp, duration, args]
 *   Redis >= 4.0       [id, timestamp, duration, args, client-address, client-name]
 *   Redis >= 8.10      [id, timestamp, duration, args, client-address, client-name, argc]
 *   Redis Enterprise   the arguments array is shifted one position to the right
 *
 * The arguments array itself is capped at 32 entries, the last one being
 * "... (N more arguments)", so argc is the only way to know the real size of the
 * command that was actually slow.
 *
 * @see https://redis.io/commands/slowlog-get
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

	// An empty slowlog replies with an empty array; anything else means no entries
	entries, _ := result.([]interface{})

	// New Frame
	frame := data.NewFrame(qm.Command,
		data.NewField("Id", nil, []int64{}),
		data.NewField("Timestamp", nil, []time.Time{}),
		data.NewField("Duration", nil, []int64{}),
		data.NewField("Command", nil, []string{}),
		data.NewField("Client Address", nil, []string{}),
		data.NewField("Client Name", nil, []string{}),
		data.NewField("Arg Count", nil, []*int64{}),
		data.NewField("Truncated", nil, []bool{}))

	// Set Field Config
	frame.Fields[2].Config = &data.FieldConfig{Unit: "µs"}

	// Parse entries
	for _, entry := range entries {
		query, ok := entry.([]interface{})
		if !ok {
			continue
		}

		// Id, timestamp and duration always lead the entry
		if len(query) < 4 {
			continue
		}

		/**
		 * Locate the arguments array. Redis Enterprise inserts an extra scalar
		 * between the duration and the arguments, so the position is discovered
		 * rather than assumed.
		 */
		argumentsID := -1
		for i := 3; i < len(query); i++ {
			if _, isArray := query[i].([]interface{}); isArray {
				argumentsID = i
				break
			}
		}

		if argumentsID < 0 {
			continue
		}

		arguments, _ := query[argumentsID].([]interface{})

		// Merge all arguments into a single command string
		parts := make([]string, 0, len(arguments))
		truncated := false

		for _, arg := range arguments {
			value, ok := redisValueToStringOK(arg)
			if !ok {
				log.DefaultLogger.Debug("Slowlog", "unsupported argument type", arg)
				continue
			}

			parts = append(parts, value)

			if strings.HasSuffix(value, slowlogTruncationMarker) {
				truncated = true
			}
		}

		command := strings.Join(parts, " ")

		// Fields that only newer servers report
		clientAddress := ""
		clientName := ""
		var argCount *int64

		if index := argumentsID + 1; index < len(query) {
			clientAddress = redisValueToString(query[index])
		}

		if index := argumentsID + 2; index < len(query) {
			clientName = redisValueToString(query[index])
		}

		// Argument count, added in Redis 8.10
		if index := argumentsID + 3; index < len(query) {
			if value, ok := redisValueToInt64(query[index]); ok {
				argCount = &value
			}
		}

		// The count disagreeing with the array is truncation too
		if argCount != nil && *argCount > int64(len(arguments)) {
			truncated = true
		}

		id, _ := redisValueToInt64(query[0])
		timestamp, _ := redisValueToInt64(query[1])
		duration, _ := redisValueToInt64(query[2])

		// Add Query
		frame.AppendRow(id, time.Unix(timestamp, 0), duration, command, clientAddress, clientName, argCount, truncated)
	}

	// Add the frame to the response
	response.Frames = append(response.Frames, frame)

	// Return Response
	return response
}

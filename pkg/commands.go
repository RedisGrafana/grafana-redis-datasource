package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
	"github.com/mediocregopher/radix/v3/resp/resp2"
)

/**
 * Query for all commands
 */
func (ds *redisDatasource) query(ctx context.Context, query backend.DataQuery, client *radix.Pool) backend.DataResponse {
	var qm queryModel

	// Unmarshal the json into our queryModel
	err := json.Unmarshal(query.JSON, &qm)
	log.DefaultLogger.Debug("QueryData", "JSON", query.JSON)

	// Error
	if err != nil {
		response := backend.DataResponse{}
		response.Error = err
		return response
	}

	// From and To
	from := query.TimeRange.From.UnixNano() / 1000000
	to := query.TimeRange.To.UnixNano() / 1000000

	// Handle Panic from any command
	defer func() {
		if err := recover(); err != nil {
			log.DefaultLogger.Error("PANIC", "command", err)
		}
	}()

	/**
	 * Custom Command using Query
	 */
	if qm.Query != "" {
		return ds.queryCustomCommand(qm, client)
	}

	/**
	 * Commands switch
	 */
	switch qm.Command {
	case "ts.range":
		return ds.queryTsRange(from, to, qm, client)
	case "ts.mrange":
		return ds.queryTsMRange(from, to, qm, client)
	case "hgetall":
		return ds.queryHGetAll(qm, client)
	case "smembers", "hkeys":
		return ds.querySMembers(qm, client)
	case "hget":
		return ds.queryHGet(qm, client)
	case "info":
		return ds.queryInfo(qm, client)
	case "clientList":
		return ds.queryClientList(qm, client)
	case "slowlogGet":
		return ds.querySlowlogGet(qm, client)
	case "type", "get", "ttl", "hlen", "xlen", "llen", "scard":
		return ds.queryKeyCommand(qm, client)
	case "xinfoStream":
		return ds.queryXInfoStream(qm, client)
	default:
		response := backend.DataResponse{}
		response.Error = fmt.Errorf("Unknown command")
		return response
	}
}

/**
 * Error Handler
 */
func (ds *redisDatasource) errorHandler(response backend.DataResponse, err error) backend.DataResponse {
	var redisErr resp2.Error

	// Check for RESP2 Error
	if errors.As(err, &redisErr) {
		response.Error = redisErr.E
	} else {
		response.Error = err
	}

	// Return Response
	return response
}

/**
 * Custom Command, used for CLI and Variables
 */
func (ds *redisDatasource) queryCustomCommand(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Query is empty
	if qm.Query == "" {
		response.Error = fmt.Errorf("Command is empty")
		return response
	}

	// Split query and parse command
	query := strings.Fields(qm.Query)
	command, query := query[0], query[1:]

	var result interface{}
	var err error

	// Run command
	if len(query) == 0 {
		err = client.Do(radix.Cmd(&result, command))
	} else {
		key, query := query[0], query[1:]
		err = client.Do(radix.FlatCmd(&result, command, key, query))
	}

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	/**
	 * Check results and add frames
	 */
	switch result.(type) {
	case int64:
		// Format number
		value := result.(int64)

		// Add Frame
		response.Frames = append(response.Frames,
			data.NewFrame(qm.Key,
				data.NewField("Value", nil, []int64{value})))
	case []byte:
		value := string(result.([]byte))

		// Split lines
		values := strings.Split(strings.Replace(value, "\r\n", "\n", -1), "\n")

		// Add Frame
		response.Frames = append(response.Frames,
			data.NewFrame(qm.Key,
				data.NewField("Value", nil, values)))
	case string:
		// Add Frame
		response.Frames = append(response.Frames,
			data.NewFrame(qm.Key,
				data.NewField("Value", nil, []string{result.(string)})))
	case []interface{}:
		var values []string

		// Parse array values
		for _, value := range result.([]interface{}) {
			switch value.(type) {
			case []byte:
				values = append(values, string(value.([]byte)))
			default:
				response.Error = fmt.Errorf("Unsupported array return type")
				return response
			}
		}

		// Add Frame
		response.Frames = append(response.Frames,
			data.NewFrame(qm.Key,
				data.NewField("Value", nil, values)))
	default:
		response.Error = fmt.Errorf("Unsupported return type")
		return response
	}

	// Return Response
	return response
}

/**
 * Commands with one key parameter and return value
 *
 * @see https://redis.io/commands/type
 * @see https://redis.io/commands/ttl
 * @see https://redis.io/commands/hlen
 */
func (ds *redisDatasource) queryKeyCommand(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var value string
	err := client.Do(radix.Cmd(&value, qm.Command, qm.Key))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Key)

	// Parse Float
	if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
		frame.Fields = append(frame.Fields, data.NewField("Value", nil, []int64{intValue}))
	} else {
		frame.Fields = append(frame.Fields, data.NewField("Value", nil, []string{value}))
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return Response
	return response
}

/**
 * TS.RANGE key fromTimestamp toTimestamp [COUNT count] [AGGREGATION aggregationType timeBucket]
 *
 * @see https://oss.redislabs.com/redistimeseries/commands/#tsrangetsrevrange
 */
func (ds *redisDatasource) queryTsRange(from int64, to int64, qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	var result [][]string
	var err error

	// Execute command
	if qm.Aggregation != "" {
		err = client.Do(radix.FlatCmd(&result, qm.Command, qm.Key, from, to, "AGGREGATION", qm.Aggregation, qm.Bucket))
	} else {
		err = client.Do(radix.FlatCmd(&result, qm.Command, qm.Key, from, to))
	}

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// Legend
	legend := qm.Key
	if qm.Legend != "" {
		legend = qm.Legend
	}

	// Create data frame response
	frame := data.NewFrame(qm.Key,
		data.NewField("time", nil, []time.Time{}),
		data.NewField(legend, nil, []float64{}))

	// Add rows
	for _, row := range result {
		t, _ := strconv.ParseInt(row[0], 10, 64)
		ts := time.Unix(t/1000, 0)
		v, _ := strconv.ParseFloat(row[1], 64)
		frame.AppendRow(ts, v)
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return Response
	return response
}

/**
 * TS.MRANGE fromTimestamp toTimestamp [COUNT count] [AGGREGATION aggregationType timeBucket] [WITHLABELS] FILTER filter..
 *
 * @see https://oss.redislabs.com/redistimeseries/commands/#tsmrangetsmrevrange
 */
func (ds *redisDatasource) queryTsMRange(from int64, to int64, qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	var result interface{}
	var err error

	// Split Filter to array
	filter := strings.Fields(qm.Filter)

	// Execute command
	if qm.Aggregation != "" {
		err = client.Do(radix.FlatCmd(&result, qm.Command, strconv.FormatInt(from, 10), to, "AGGREGATION", qm.Aggregation, qm.Bucket, "WITHLABELS", "FILTER", filter))
	} else {
		err = client.Do(radix.FlatCmd(&result, qm.Command, strconv.FormatInt(from, 10), to, "WITHLABELS", "FILTER", filter))
	}

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// Check results
	switch result.(type) {
	case string:
		response.Error = fmt.Errorf(result.(string))
		return response
	default:
	}

	// Parse Time-Series data
	for _, innerArray := range result.([]interface{}) {
		tsArrReply := innerArray.([]interface{})

		// Labels
		labelsRaw := tsArrReply[1].([]interface{})
		labels := make(map[string]string, len(labelsRaw))

		// Parse Labels
		for _, labelRaw := range labelsRaw {
			kvPair := labelRaw.([]interface{})
			k := string(kvPair[0].([]byte))
			v := string(kvPair[1].([]byte))
			labels[k] = v
		}

		// Use Time-series's name as Legend if Legend label is not specified
		legend := string(tsArrReply[0].([]byte))
		if qm.Legend != "" {
			legend = labels[qm.Legend]
		}

		// Use value's label if specified
		value := "value"
		if qm.Value != "" {
			value = labels[qm.Value]
		}

		// Create Frame
		frame := data.NewFrame(legend,
			data.NewField("time", nil, []time.Time{}))

		// Return labels if legend is not specified
		if qm.Legend != "" {
			frame.Fields = append(frame.Fields,
				data.NewField(value, nil, []float64{}),
			)
		} else {
			frame.Fields = append(frame.Fields,
				data.NewField(value, labels, []float64{}),
			)
		}

		// Values
		for _, valueRaw := range tsArrReply[2].([]interface{}) {
			kvPair := valueRaw.([]interface{})
			var k int64
			var v float64

			// Key
			switch kvPair[0].(type) {
			case []byte:
				k, _ = strconv.ParseInt(string(kvPair[0].([]byte)), 10, 64)
			default:
				k = kvPair[0].(int64)
			}

			// Value
			switch kvPair[1].(type) {
			case []byte:
				v, _ = strconv.ParseFloat(string(kvPair[1].([]byte)), 64)
			default:
				v, _ = strconv.ParseFloat(kvPair[1].(string), 64)
			}

			// Append Row to Frame
			frame.AppendRow(time.Unix(k/1000, 0), v)
		}

		// add the frames to the response
		response.Frames = append(response.Frames, frame)
	}

	// Return Response
	return response
}

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

	// New Frame
	frame := data.NewFrame(qm.Key,
		data.NewField("Value", nil, []string{value}))

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

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

/**
 * INFO [section]
 *
 * @see https://redis.io/commands/info
 */
func (ds *redisDatasource) queryInfo(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result string
	err := client.Do(radix.Cmd(&result, qm.Command, qm.Section))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// Split lines
	lines := strings.Split(strings.Replace(result, "\r\n", "\n", -1), "\n")

	// Command stats
	if qm.Section == "commandstats" {
		// New Frame
		frame := data.NewFrame(qm.Command,
			data.NewField("Command", nil, []string{}),
			data.NewField("Calls", nil, []int64{}),
			data.NewField("Usec", nil, []float64{}),
			data.NewField("Usec_per_call", nil, []float64{}))

		// Parse lines
		for _, line := range lines {
			fields := strings.Split(line, ":")

			if len(fields) < 2 {
				continue
			}

			// Stats
			stats := strings.Split(fields[1], ",")

			if len(stats) < 3 {
				continue
			}

			// Parse Stats
			calls := strings.Split(stats[0], "=")
			usec := strings.Split(stats[1], "=")
			usecPerCall := strings.Split(stats[2], "=")

			var callsValue int64
			var usecValue float64
			var usecPerCallValue float64

			// Parse Calls
			if len(calls) == 2 {
				callsValue, _ = strconv.ParseInt(calls[1], 10, 64)
			}

			// Parse Usec
			if len(usec) == 2 {
				usecValue, _ = strconv.ParseFloat(usec[1], 64)
			}

			// Parse Usec per Call
			if len(usecPerCall) == 2 {
				usecPerCallValue, _ = strconv.ParseFloat(usecPerCall[1], 64)
			}

			// Add Query
			frame.AppendRow(fields[0], callsValue, usecValue, usecPerCallValue)
		}

		// Add the frames to the response
		response.Frames = append(response.Frames, frame)

		// Return
		return response
	}

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
 * CLIENT LIST [TYPE normal|master|replica|pubsub]
 *
 * @see https://redis.io/commands/client-list
 */
func (ds *redisDatasource) queryClientList(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result string
	err := client.Do(radix.Cmd(&result, "CLIENT", "LIST"))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
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

			// Add Header for first row
			if i == 0 {
				if _, err := strconv.ParseInt(value[1], 10, 64); err == nil {
					frame.Fields = append(frame.Fields, data.NewField(value[0], nil, []int64{}))
				} else {
					frame.Fields = append(frame.Fields, data.NewField(value[0], nil, []string{}))
				}
			}

			// Skip if less than 2 elements
			if len(value) < 2 {
				continue
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
func (ds *redisDatasource) querySlowlogGet(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result interface{}
	err := client.Do(radix.Cmd(&result, "SLOWLOG", "GET"))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Command,
		data.NewField("Id", nil, []int64{}),
		data.NewField("Timestamp", nil, []int64{}),
		data.NewField("Duration", nil, []int64{}),
		data.NewField("Command", nil, []string{}))

	// Parse Time-Series data
	for _, innerArray := range result.([]interface{}) {
		query := innerArray.([]interface{})
		command := ""

		// Merge all args
		for _, arg := range query[3].([]interface{}) {

			// Add space between command and arguments
			if command != "" {
				command += " "
			}

			// Combine args into single command
			switch arg.(type) {
			case int64:
				command += strconv.FormatInt(arg.(int64), 10)
			case []byte:
				command += string(arg.([]byte))
			case string:
				command += arg.(string)
			default:
				log.DefaultLogger.Debug("Slowlog", "default", arg)
			}
		}

		// Add Query
		frame.AppendRow(query[0].(int64), query[1].(int64), query[2].(int64), command)
	}

	// Add the frame to the response
	response.Frames = append(response.Frames, frame)

	// Return Response
	return response
}

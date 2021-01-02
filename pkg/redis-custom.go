package main

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/creachadair/shell"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
)

// EmptyArray for (empty array)
const EmptyArray = "(empty array)"

/**
 * Execute Query
 * Can PANIC if command is wrong
 */
func (ds *redisDatasource) executeCustomQuery(qm queryModel, client ClientInterface) (interface{}, error) {
	var result interface{}
	var err error

	// Split query
	query, ok := shell.Split(qm.Query)

	// Check if query is valid
	if !ok {
		err = fmt.Errorf("Query is not valid")
		return result, err
	}

	// Separate command from params
	command, params := query[0], query[1:]

	// Handle Panic from custom command to catch "should never get here"
	defer func() {
		if err := recover(); err != nil {
			log.DefaultLogger.Error("PANIC", "command", err, "query", qm.Query)
		}
	}()

	// Run command without params
	if len(params) == 0 {
		err = client.Do(radix.Cmd(&result, command))
		return result, err
	}

	// Extract key or 1st parameter as required for FlatCmd
	key, params := params[0], params[1:]
	err = client.Do(radix.FlatCmd(&result, command, key, params))

	return result, err
}

/**
 * Parse Value
 */
func (ds *redisDatasource) parseInterfaceValue(value []interface{}, response backend.DataResponse) ([]string, backend.DataResponse) {
	var values []string

	for _, element := range value {
		switch element := element.(type) {
		case []byte:
			values = append(values, string(element))
		case int64:
			values = append(values, strconv.FormatInt(element, 10))
		case string:
			values = append(values, element)
		case []interface{}:
			var parsedValues []string
			parsedValues, response = ds.parseInterfaceValue(element, response)

			// If no values
			if len(parsedValues) == 0 {
				parsedValues = append(parsedValues, EmptyArray)
			}

			values = append(values, parsedValues...)
		default:
			response.Error = fmt.Errorf("Unsupported array return type")
			return values, response
		}
	}

	return values, response
}

/**
 * Custom Command, used for CLI and Variables
 */
func (ds *redisDatasource) queryCustomCommand(qm queryModel, client ClientInterface) backend.DataResponse {
	response := backend.DataResponse{}

	// Query is empty
	if qm.Query == "" {
		response.Error = fmt.Errorf("Command is empty")
		return response
	}

	var result interface{}
	var err error

	// Parse and execute query
	result, err = ds.executeCustomQuery(qm, client)

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	/**
	 * Check results and add frames
	 */
	switch result := result.(type) {
	case int64:
		// Add Frame
		response.Frames = append(response.Frames,
			data.NewFrame(qm.Key,
				data.NewField("Value", nil, []int64{result})))
	case []byte:
		value := string(result)

		// Split lines
		values := strings.Split(strings.Replace(value, "\r\n", "\n", -1), "\n")

		// Parse float if only one value
		if len(values) == 1 {
			response.Frames = append(response.Frames, ds.createFrameValue(qm.Key, values[0]))
			break
		}

		// Add Frame
		response.Frames = append(response.Frames,
			data.NewFrame(qm.Key,
				data.NewField("Value", nil, values)))
	case string:
		// Add Frame
		response.Frames = append(response.Frames, ds.createFrameValue(qm.Key, result))
	case []interface{}:
		var values []string

		// Parse values
		if len(values) == 0 {
			values = append(values, EmptyArray)
		} else {
			values, response = ds.parseInterfaceValue(result, response)
		}

		// Error when parsing intarface
		if response.Error != nil {
			return response
		}

		// Add Frame
		response.Frames = append(response.Frames,
			data.NewFrame(qm.Key,
				data.NewField("Value", nil, values)))
	case nil:
		response.Error = fmt.Errorf("Wrong command")
	default:
		response.Error = fmt.Errorf("Unsupported return type")
	}

	// Return Response
	return response
}

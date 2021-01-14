package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"

	"bitbucket.org/creachadair/shell"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// EmptyArray for (empty array)
const EmptyArray = "(empty array)"

/**
 * Execute Query
 * Can PANIC if command is wrong
 */
func executeCustomQuery(qm queryModel, client redisClient) (interface{}, error) {
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
		err = client.RunCmd(&result, command)
		return result, err
	}

	// Extract key or 1st parameter as required for RunFlatCmd
	key, params := params[0], params[1:]
	err = client.RunFlatCmd(&result, command, key, params)

	return result, err
}

/**
 * Parse Value
 */
func parseInterfaceValue(value []interface{}, response backend.DataResponse) ([]string, backend.DataResponse) {
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
			parsedValues, response = parseInterfaceValue(element, response)

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
func queryCustomCommand(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Query is empty
	if qm.Query == "" {
		response.Error = fmt.Errorf("Command is empty")
		return response
	}

	var result interface{}
	var err error

	// Parse and execute query
	result, err = executeCustomQuery(qm, client)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Command-line mode enabled
	if qm.CLI {
		var builder strings.Builder

		// Use a tab writer for having CLI-like tabulation aligned right
		tabWriter := tabwriter.NewWriter(&builder, 0, 1, 0, ' ', tabwriter.AlignRight)

		// Concatenate everything to string with proper tabs and newlines and pass it to tabWriter for tab formatting
		_, err := fmt.Fprint(tabWriter, convertToCLI(result, ""))

		// Check formatting error
		if err != nil {
			log.DefaultLogger.Error("Error when writing to TabWriter", "error", err.Error(), "query", qm.Query)
		}

		// Check tab writer error
		err = tabWriter.Flush()
		if err != nil {
			log.DefaultLogger.Error("Error when flushing TabWriter", "error", err.Error(), "query", qm.Query)
		}

		// Get the properly formatted string from the string builder
		processed := builder.String()

		// Write result string as a single frame with a single field with name "Value"
		response.Frames = append(response.Frames, data.NewFrame(qm.Key, data.NewField("Value", nil, []string{processed})))
	} else {
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
				response.Frames = append(response.Frames, createFrameValue(qm.Key, values[0]))
				break
			}

			// Add Frame
			response.Frames = append(response.Frames,
				data.NewFrame(qm.Key,
					data.NewField("Value", nil, values)))
		case string:
			// Add Frame
			response.Frames = append(response.Frames, createFrameValue(qm.Key, result))
		case []interface{}:
			var values []string

			// Parse values
			if len(result) == 0 {
				values = append(values, EmptyArray)
			} else {
				values, response = parseInterfaceValue(result, response)
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
	}
	// Return Response
	return response
}

/**
 * Convert results to CLI format
 */

func convertToCLI(input interface{}, tabs string) string {
	switch value := input.(type) {
	case int64:
		return fmt.Sprintf("(integer) %d\n", value)
	case []byte:
		return fmt.Sprintf("\"%v\"\n", string(value))
	case string:
		return fmt.Sprintf("\"%v\"\n", value)
	case []interface{}:
		if len(value) < 1 {
			return EmptyArray + "\n"
		}

		var builder strings.Builder
		for i, member := range value {
			additionalTabs := ""
			if i != 0 {
				additionalTabs = tabs
			}

			builder.WriteString(fmt.Sprintf("%v%d)\t %v", additionalTabs, i+1, convertToCLI(member, tabs+"\t")))
		}

		return builder.String()
	case nil:
		return "(nil)\n"
	default:
		log.DefaultLogger.Error("Unsupported type for CLI mode", "value", value, "type", reflect.TypeOf(value).String())
		return fmt.Sprintf("\"%v\"\n", value)
	}
}

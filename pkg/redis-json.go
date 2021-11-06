package main

import (
	"encoding/json"
	"reflect"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
)

/**
 * JSON.OBJKEYS <key> [path]
 *
 * @see https://oss.redis.com/redisjson/commands/#jsonobjkeys
 */
func queryJsonObjKeys(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var values []string
	err := client.RunFlatCmd(&values, qm.Command, qm.Key, qm.Path)

	// Check error
	if err != nil {
		return errorHandler(response, err)
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
 * JSON.OBJLEN <key> [path]
 *
 * @see https://oss.redis.com/redisjson/commands/#jsonobjlen
 */
func queryJsonObjLen(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var value string
	err := client.RunCmd(&value, qm.Command, qm.Key, qm.Path)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, createFrameValue(qm.Key, value, "Value"))

	// Return Response
	return response
}

/**
 * JSON.GET <key> [path]
 *
 * @see https://oss.redis.com/redisjson/commands/#jsonget
 */
func queryJsonGet(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var value string
	err := client.RunCmd(&value, qm.Command, qm.Key, qm.Path)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	var result interface{}
	err = json.Unmarshal([]byte(value), &result)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Command)

	// Parse result
	switch value := result.(type) {
	case string:
		frame.Fields = append(frame.Fields, data.NewField(qm.Key, nil, []string{value}))
	case bool:
		frame.Fields = append(frame.Fields, data.NewField(qm.Key, nil, []bool{value}))
	case map[string]interface{}:
		for i, value := range value {
			// Value
			switch v := value.(type) {
			case string:
				frame.Fields = append(frame.Fields, data.NewField(i, nil, []string{v}))
			case bool:
				frame.Fields = append(frame.Fields, data.NewField(i, nil, []bool{v}))
			case float64:
				frame.Fields = append(frame.Fields, data.NewField(i, nil, []float64{v}))
			default:
				log.DefaultLogger.Error(models.JsonGet, "Conversion Error", "Unsupported Value type")
			}
		}
	case []interface{}:
		// Map for storing all the fields found in entries
		fields := map[string]*data.Field{}
		rowscount := 0

		for _, entry := range value {
			keysFoundInCurrentEntry := map[string]bool{}
			rowscount++

			for i, value := range entry.(map[string]interface{}) {
				// Value
				switch v := value.(type) {
				case bool:
					if _, ok := fields[i]; !ok {
						fields[i] = data.NewField(i, nil, []bool{})
						frame.Fields = append(frame.Fields, fields[i])

						// Generate empty values for all previous rows
						for j := 0; j < rowscount-1; j++ {
							fields[i].Append(false)
						}
					}

					// Insert value for current row
					fields[i].Append(v)
					keysFoundInCurrentEntry[i] = true
				case string:
					if _, ok := fields[i]; !ok {
						fields[i] = data.NewField(i, nil, []string{})
						frame.Fields = append(frame.Fields, fields[i])

						// Generate empty values for all previous rows
						for j := 0; j < rowscount-1; j++ {
							fields[i].Append("")
						}
					}

					// Insert value for current row
					fields[i].Append(v)
					keysFoundInCurrentEntry[i] = true
				case float64:
					if _, ok := fields[i]; !ok {
						fields[i] = data.NewField(i, nil, []float64{})
						frame.Fields = append(frame.Fields, fields[i])

						// Generate empty values for all previous rows
						for j := 0; j < rowscount-1; j++ {
							fields[i].Append(0.0)
						}
					}

					// Insert value for current row
					fields[i].Append(v)
					keysFoundInCurrentEntry[i] = true
				default:
					log.DefaultLogger.Error(models.JsonGet, "Conversion Error", "Unsupported Value type")
				}
			}

			// Iterate over all keys found so far for stream
			for key, field := range fields {
				// Check if key exist in entry
				if _, ok := keysFoundInCurrentEntry[key]; !ok {
					if field.Type() == data.FieldTypeFloat64 {
						field.Append(0.0)
						continue
					}

					if field.Type() == data.FieldTypeBool {
						field.Append(false)
						continue
					}

					field.Append("")
				}
			}
		}
	default:
		log.DefaultLogger.Error("Unexpected type received", "value", value, "type", reflect.TypeOf(value).String())
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

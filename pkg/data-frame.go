package main

import (
	"strconv"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * Create frame with single value
 *
 * @param {string} key Key
 * @param {string} value Value
 */
func createFrameValue(key string, value string, field string) *data.Frame {
	frame := data.NewFrame(key)

	// Parse Float
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		frame.Fields = append(frame.Fields, data.NewField(field, nil, []float64{floatValue}))
	} else {
		frame.Fields = append(frame.Fields, data.NewField(field, nil, []string{value}))
	}

	// Return
	return frame
}

/**
 * Add Frame Fields from Array
 */
func addFrameFieldsFromArray(values []interface{}, frame *data.Frame) *data.Frame {
	for _, value := range values {
		pair := value.([]interface{})
		var key string

		// Key
		switch k := pair[0].(type) {
		case []byte:
			key = string(k)
		default:
			log.DefaultLogger.Error("addFrameFieldsFromArray", "Conversion Error", "Unsupported Key type")
			continue
		}

		// Value
		switch v := pair[1].(type) {
		case []byte:
			value := string(v)

			// Is it Integer?
			if valueInt, err := strconv.ParseInt(value, 10, 64); err == nil {
				frame.Fields = append(frame.Fields, data.NewField(key, nil, []int64{valueInt}))
				break
			}

			// Add as string
			frame.Fields = append(frame.Fields, data.NewField(key, nil, []string{value}))
		case int64:
			frame.Fields = append(frame.Fields, data.NewField(key, nil, []int64{v}))
		case float64:
			frame.Fields = append(frame.Fields, data.NewField(key, nil, []float64{v}))
		default:
			log.DefaultLogger.Error("addFrameFieldsFromArray", "Conversion Error", "Unsupported Value type")
		}
	}

	return frame
}

package main

import (
	"strconv"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
)

/**
 * FT.INFO {index}
 *
 * @see https://oss.redislabs.com/redisearch/Commands/#ftinfo
 */
func queryFtInfo(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result map[string]interface{}
	err := client.RunCmd(&result, qm.Command, qm.Key)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Create data frame response
	frame := data.NewFrame(qm.Key)

	// Add fields and values
	for key := range result {
		// Value
		switch value := result[key].(type) {
		case int64:
			// Add field
			field := data.NewField(key, nil, []int64{value})
			frame.Fields = append(frame.Fields, field)
		case []byte:
			// Parse Float
			if floatValue, err := strconv.ParseFloat(string(value), 64); err == nil {
				field := data.NewField(key, nil, []float64{floatValue})

				// Set unit
				if models.SearchInfoConfig[key] != "" {
					field.Config = &data.FieldConfig{Unit: models.SearchInfoConfig[key]}
				}

				frame.Fields = append(frame.Fields, field)
			} else {
				frame.Fields = append(frame.Fields, data.NewField(key, nil, []string{string(value)}))
			}
		case string:
			frame.Fields = append(frame.Fields, data.NewField(key, nil, []string{string(value)}))
		case []interface{}:
		default:
			log.DefaultLogger.Error(models.SearchInfo, "Conversion Error", "Unsupported Value type")
		}
	}

	// Add the frame to the response
	response.Frames = append(response.Frames, frame)

	// Return Response
	return response
}

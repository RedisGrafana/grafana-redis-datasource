package main

import (
	"strconv"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
)

/**
 * FT.INFO {index}
 *
 * @see https://oss.redislabs.com/redisearch/Commands/#ftinfo
 */
func (ds *redisDatasource) queryFtInfo(qm queryModel, client ClientInterface) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result map[string]interface{}
	err := client.Do(radix.Cmd(&result, qm.Command, qm.Key))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
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

				// Field Units
				config := map[string]string{"inverted_sz_mb": "decmbytes", "offset_vectors_sz_mb": "decmbytes",
					"doc_table_size_mb": "decmbytes", "sortable_values_size_mb": "decmbytes",
					"key_table_size_mb": "decmbytes", "percent_indexed": "percentunit"}

				// Set unit
				if config[key] != "" {
					field.Config = &data.FieldConfig{Unit: config[key]}
				}

				frame.Fields = append(frame.Fields, field)
			} else {
				frame.Fields = append(frame.Fields, data.NewField(key, nil, []string{string(value)}))
			}
		case string:
			frame.Fields = append(frame.Fields, data.NewField(key, nil, []string{string(value)}))
		case []interface{}:
		default:
			log.DefaultLogger.Error("queryTsInfo", "Conversion Error", "Unsupported Value type")
		}
	}

	// Add the frame to the response
	response.Frames = append(response.Frames, frame)

	// Return Response
	return response
}

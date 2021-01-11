package main

import (
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/creachadair/shell"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * TS.RANGE key fromTimestamp toTimestamp [COUNT count] [AGGREGATION aggregationType timeBucket]
 *
 * @see https://oss.redislabs.com/redistimeseries/commands/#tsrangetsrevrange
 */
func queryTsRange(from int64, to int64, qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	var result [][]string
	var err error

	// Execute command
	if qm.Aggregation != "" {
		err = client.RunFlatCmd(&result, qm.Command, qm.Key, from, to, "AGGREGATION", qm.Aggregation, qm.Bucket)
	} else {
		err = client.RunFlatCmd(&result, qm.Command, qm.Key, from, to)
	}

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Legend
	legend := qm.Key
	if qm.Legend != "" {
		legend = qm.Legend
	}

	// Create data frame response
	frame := data.NewFrame(legend,
		data.NewField("time", nil, []time.Time{}),
		data.NewField(qm.Value, nil, []float64{}))

	// Previous time to fill missing intervals
	var prevTime time.Time

	// Add rows
	for _, row := range result {
		t, _ := strconv.ParseInt(row[0], 10, 64)
		ts := time.Unix(0, t*int64(time.Millisecond))

		// Fill missing intervals
		if qm.Fill && qm.Bucket != 0 {
			if !prevTime.IsZero() {
				for ts.Sub(prevTime) > time.Duration(qm.Bucket)*time.Millisecond {
					prevTime = prevTime.Add(time.Duration(qm.Bucket) * time.Millisecond)
					frame.AppendRow(prevTime, float64(0))
				}
			}

			prevTime = ts
		}

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
func queryTsMRange(from int64, to int64, qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	var result interface{}
	var err error

	// Split Filter to array
	filter, ok := shell.Split(qm.Filter)

	// Check if filter is valid
	if !ok {
		response.Error = fmt.Errorf("Filter is not valid")
		return response
	}

	// Execute command
	if qm.Aggregation != "" {
		err = client.RunFlatCmd(&result, qm.Command, strconv.FormatInt(from, 10), to, "AGGREGATION", qm.Aggregation, qm.Bucket, "WITHLABELS", "FILTER", filter)
	} else {
		err = client.RunFlatCmd(&result, qm.Command, strconv.FormatInt(from, 10), to, "WITHLABELS", "FILTER", filter)
	}

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Check results
	switch result := result.(type) {
	case string:
		response.Error = fmt.Errorf(result)
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
		value := ""
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

		// Previous time to fill missing intervals
		var prevTime time.Time

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

			ts := time.Unix(0, k*int64(time.Millisecond))

			// Fill missing intervals
			if qm.Fill && qm.Bucket != 0 {
				if !prevTime.IsZero() {
					for ts.Sub(prevTime) > time.Duration(qm.Bucket)*time.Millisecond {
						prevTime = prevTime.Add(time.Duration(qm.Bucket) * time.Millisecond)
						frame.AppendRow(prevTime, float64(0))
					}
				}

				prevTime = ts
			}

			// Value
			switch kvPair[1].(type) {
			case []byte:
				v, _ = strconv.ParseFloat(string(kvPair[1].([]byte)), 64)
			default:
				v, _ = strconv.ParseFloat(kvPair[1].(string), 64)
			}

			// Append Row to Frame
			frame.AppendRow(ts, v)
		}

		// add the frames to the response
		response.Frames = append(response.Frames, frame)
	}

	// Return Response
	return response
}

/**
 * TS.GET key
 *
 * @see https://oss.redislabs.com/redistimeseries/1.4/commands/#tsget
 */
func queryTsGet(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result []string
	err := client.RunCmd(&result, qm.Command, qm.Key)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Create data frame response
	frame := data.NewFrame(qm.Key,
		data.NewField("time", nil, []time.Time{}),
		data.NewField("value", nil, []float64{}))

	// Add row
	t, _ := strconv.ParseInt(result[0], 10, 64)
	v, _ := strconv.ParseFloat(result[1], 64)
	frame.AppendRow(time.Unix(0, t*int64(time.Millisecond)), v)

	// Add the frame to the response
	response.Frames = append(response.Frames, frame)

	// Return Response
	return response
}

/**
 * TS.INFO key
 *
 * @see https://oss.redislabs.com/redistimeseries/1.4/commands/#tsinfo
 */
func queryTsInfo(qm queryModel, client redisClient) backend.DataResponse {
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
			// Return timestamp as time
			if key == "firstTimestamp" || key == "lastTimestamp" {
				frame.Fields = append(frame.Fields, data.NewField(key, nil, []time.Time{time.Unix(0, value*int64(time.Millisecond))}))
				break
			}

			// Add field
			field := data.NewField(key, nil, []int64{value})

			// Set unit
			if key == "memoryUsage" {
				field.Config = &data.FieldConfig{Unit: "decbytes"}
			} else if key == "retentionTime" {
				field.Config = &data.FieldConfig{Unit: "ms"}
			}

			frame.Fields = append(frame.Fields, field)
		case []byte:
			frame.Fields = append(frame.Fields, data.NewField(key, nil, []string{string(value)}))
		case []interface{}:
			frame = addFrameFieldsFromArray(value, frame)
		default:
			log.DefaultLogger.Error("queryTsInfo", "Conversion Error", "Unsupported Value type")
		}
	}
	// Add the frame to the response
	response.Frames = append(response.Frames, frame)

	// Return Response
	return response
}

/**
 * TS.QUERYINDEX filter...
 *
 * @see https://oss.redislabs.com/redistimeseries/commands/#tsqueryindex
 */
func queryTsQueryIndex(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Split Filter to array
	filter, ok := shell.Split(qm.Filter)

	// Check if filter is valid
	if !ok {
		response.Error = fmt.Errorf("Filter is not valid")
		return response
	}

	// Execute command
	var values []string
	err := client.RunCmd(&values, qm.Command, filter...)

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

package main

import (
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/creachadair/shell"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
)

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

	// Previous time and bucket to fill missing intervals
	var prevTime time.Time
	var bucket, _ = strconv.ParseInt(qm.Bucket, 10, 64)

	// Add rows
	for _, row := range result {
		t, _ := strconv.ParseInt(row[0], 10, 64)
		ts := time.Unix(t/1000, 0)

		// Fill missing intervals
		if qm.Fill && bucket != 0 {
			if !prevTime.IsZero() {
				for ts.Sub(prevTime) > time.Duration(bucket)*time.Millisecond {
					prevTime = prevTime.Add(time.Duration(bucket) * time.Millisecond)
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
func (ds *redisDatasource) queryTsMRange(from int64, to int64, qm queryModel, client *radix.Pool) backend.DataResponse {
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
		err = client.Do(radix.FlatCmd(&result, qm.Command, strconv.FormatInt(from, 10), to, "AGGREGATION", qm.Aggregation, qm.Bucket, "WITHLABELS", "FILTER", filter))
	} else {
		err = client.Do(radix.FlatCmd(&result, qm.Command, strconv.FormatInt(from, 10), to, "WITHLABELS", "FILTER", filter))
	}

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
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

		// Previous time and bucket to fill missing intervals
		var prevTime time.Time
		var bucket, _ = strconv.ParseInt(qm.Bucket, 10, 64)

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

			ts := time.Unix(k/1000, 0)

			// Fill missing intervals
			if qm.Fill && bucket != 0 {
				if !prevTime.IsZero() {
					for ts.Sub(prevTime) > time.Duration(bucket)*time.Millisecond {
						prevTime = prevTime.Add(time.Duration(bucket) * time.Millisecond)
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
func (ds *redisDatasource) queryTsGet(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result []string
	err := client.Do(radix.Cmd(&result, qm.Command, qm.Key))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// Create data frame response
	frame := data.NewFrame(qm.Key,
		data.NewField("time", nil, []time.Time{}),
		data.NewField("value", nil, []float64{}))

	// Add row
	t, _ := strconv.ParseInt(result[0], 10, 64)
	v, _ := strconv.ParseFloat(result[1], 64)
	frame.AppendRow(time.Unix(t/1000, 0), v)

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
func (ds *redisDatasource) queryTsInfo(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result interface{}
	err := client.Do(radix.Cmd(&result, qm.Command, qm.Key))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// Create data frame response
	frame := data.NewFrame(qm.Key)

	// Add fields and values
	for i := 0; i < len(result.([]interface{})); i += 2 {

		// Parameter
		var param string
		switch value := result.([]interface{})[i].(type) {
		case string:
			param = value
		default:
			log.DefaultLogger.Error("queryTsInfo", "Conversion Error", "Unsupported Key type")
		}

		// Value
		switch value := result.([]interface{})[i+1].(type) {
		case int64:
			// Return timestamp as time
			if param == "firstTimestamp" || param == "lastTimestamp" {
				frame.Fields = append(frame.Fields, data.NewField(param, nil, []time.Time{time.Unix(value/1000, 0)}))
				break
			}

			// Add field
			field := data.NewField(param, nil, []int64{value})

			// Set unit
			if param == "memoryUsage" {
				field.Config = &data.FieldConfig{Unit: "decbytes"}
			} else if param == "retentionTime" {
				field.Config = &data.FieldConfig{Unit: "ms"}
			}

			frame.Fields = append(frame.Fields, field)
		case []byte:
			frame.Fields = append(frame.Fields, data.NewField(param, nil, []string{string(value)}))
		case []interface{}:
			frame = ds.addFrameFieldsFromArray(value, frame)
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
func (ds *redisDatasource) queryTsQueryIndex(qm queryModel, client *radix.Pool) backend.DataResponse {
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
	err := client.Do(radix.Cmd(&values, qm.Command, filter...))

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

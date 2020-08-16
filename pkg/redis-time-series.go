package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
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
 * TS.GET key
 *
 * @see https://oss.redislabs.com/redistimeseries/1.4/commands/#tsget
 */
func (ds *redisDatasource) queryTsGet(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	var result []string

	// Execute command
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

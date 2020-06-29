package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
	"github.com/mediocregopher/radix/v3/resp/resp2"
)

// Error Handler
func (ds *redisDatasource) errorHandler(response backend.DataResponse, err error) backend.DataResponse {
	var redisErr resp2.Error

	// Check for RESP2 Error
	if errors.As(err, &redisErr) {
		response.Error = redisErr.E
	} else {
		response.Error = err
	}

	return response
}

// TS.RANGE
func (ds *redisDatasource) queryTsRange(from int64, to int64, qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	var res [][]string
	var err error

	// Execute command
	if qm.Aggregation != "" {
		err = client.Do(radix.FlatCmd(&res, "TS.RANGE", qm.Key, from, to, "AGGREGATION", qm.Aggregation, qm.Bucket))
	} else {
		err = client.Do(radix.FlatCmd(&res, "TS.RANGE", qm.Key, from, to))
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
	for _, row := range res {
		t, _ := strconv.ParseInt(row[0], 10, 64)
		ts := time.Unix(t/1000, 0)
		v, _ := strconv.ParseFloat(row[1], 64)
		frame.AppendRow(ts, v)
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

// TS.MRANGE
func (ds *redisDatasource) queryTsMRange(from int64, to int64, qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	var res interface{}
	var err error

	// Split Filter to array
	filter := strings.Fields(qm.Filter)

	// Execute command
	if qm.Aggregation != "" {
		err = client.Do(radix.FlatCmd(&res, "TS.MRANGE", strconv.FormatInt(from, 10), to, "AGGREGATION", qm.Aggregation, qm.Bucket, "WITHLABELS", "FILTER", filter))
	} else {
		err = client.Do(radix.FlatCmd(&res, "TS.MRANGE", strconv.FormatInt(from, 10), to, "WITHLABELS", "FILTER", filter))
	}

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// Check results
	switch res.(type) {
	case string:
		response.Error = fmt.Errorf(res.(string))
		return response
	default:
	}

	// Time-Series
	for _, innerArray := range res.([]interface{}) {
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

	// Return
	return response
}

// HGETALL
func (ds *redisDatasource) queryHGetAll(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var res []string
	err := client.Do(radix.FlatCmd(&res, "HGETALL", qm.Key))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	fields := []string{}
	values := []string{}

	// Add fields and values
	for i := 0; i < len(res); i += 2 {
		fields = append(fields, res[i])
		values = append(values, res[i+1])
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

// SMEMBERS
func (ds *redisDatasource) querySMembers(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var values []string
	err := client.Do(radix.FlatCmd(&values, "SMEMBERS", qm.Key))

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

// HGET
func (ds *redisDatasource) queryHGet(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var value string
	err := client.Do(radix.FlatCmd(&value, "HGET", qm.Key, qm.Field))

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

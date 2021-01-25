package main

import (
	"bytes"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * XINFO Radix marshaling
 */
type xinfo struct {
	Length          int64         `redis:"length"`
	RadixTreeKeys   int64         `redis:"radix-tree-keys"`
	RadixTreeNodes  int64         `redis:"radix-tree-nodes"`
	Groups          int64         `redis:"groups"`
	LastGeneratedID string        `redis:"last-generated-id"`
	FirstEntry      []interface{} `redis:"first-entry"`
	LastEntry       []interface{} `redis:"last-entry"`
}

/**
 * XINFO [CONSUMERS key groupname] [GROUPS key] [STREAM key] [HELP]
 *
 * @see https://redis.io/commands/xinfo
 */
func queryXInfoStream(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var model xinfo
	err := client.RunFlatCmd(&model, "XINFO", "STREAM", qm.Key)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Key)

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Add plain fields to frame
	frame.Fields = append(frame.Fields, data.NewField("length", nil, []int64{model.Length}))
	frame.Fields = append(frame.Fields, data.NewField("radix-tree-keys", nil, []int64{model.RadixTreeKeys}))
	frame.Fields = append(frame.Fields, data.NewField("radix-tree-nodes", nil, []int64{model.RadixTreeNodes}))
	frame.Fields = append(frame.Fields, data.NewField("groups", nil, []int64{model.Groups}))
	frame.Fields = append(frame.Fields, data.NewField("last-generated-id", nil, []string{model.LastGeneratedID}))

	// First entry
	if model.FirstEntry != nil {
		frame.Fields = append(frame.Fields, data.NewField("first-entry-id", nil, []string{string(model.FirstEntry[0].([]byte))}))
		entryFields := model.FirstEntry[1].([]interface{})
		fields := new(bytes.Buffer)

		// Merging args to string like "key"="value"\n
		for i := 0; i < len(entryFields); i += 2 {
			field := string(entryFields[i].([]byte))
			value := string(entryFields[i+1].([]byte))
			fmt.Fprintf(fields, "\"%s\"=\"%s\"\n", field, value)
		}

		frame.Fields = append(frame.Fields, data.NewField("first-entry-fields", nil, []string{fields.String()}))
	}

	// Last entry
	if model.LastEntry != nil {
		frame.Fields = append(frame.Fields, data.NewField("last-entry-id", nil, []string{string(model.LastEntry[0].([]byte))}))
		entryFields := model.LastEntry[1].([]interface{})
		fields := new(bytes.Buffer)

		// Merging args to string like "key"="value"\n
		for i := 0; i < len(entryFields); i += 2 {
			field := string(entryFields[i].([]byte))
			value := string(entryFields[i+1].([]byte))
			fmt.Fprintf(fields, "\"%s\"=\"%s\"\n", field, value)
		}

		frame.Fields = append(frame.Fields, data.NewField("last-entry-fields", nil, []string{fields.String()}))
	}

	// Return
	return response
}

/**
 * XRANGE key start end [COUNT count]
 *
 * @see https://redis.io/commands/xrange
 */
func queryXRange(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Start
	start := "-"
	if qm.Start != "" {
		start = qm.Start
	}

	// End
	end := "+"
	if qm.End != "" {
		end = qm.End
	}

	// Arguments
	args := []interface{}{start, end}
	if qm.Count > 0 {
		args = append(args, "COUNT", qm.Count)
	}

	var result []interface{}

	// Execute command
	err := client.RunFlatCmd(&result, "XRANGE", qm.Key, args...)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Create frame
	frame := createFrameFromRangeResponse(qm.Command, result)

	// Add the frame to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

/**
 * XREVRANGE key end start [COUNT count]
 *
 * @see https://redis.io/commands/xrevrange
 */
func queryXRevRange(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Start
	start := "-"
	if qm.Start != "" {
		start = qm.Start
	}

	// End
	end := "+"
	if qm.End != "" {
		end = qm.End
	}

	// Arguments
	args := []interface{}{end, start}
	if qm.Count > 0 {
		args = append(args, "COUNT", qm.Count)
	}

	var result []interface{}

	// Execute command
	err := client.RunFlatCmd(&result, "XREVRANGE", qm.Key, args...)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Create frame
	frame := createFrameFromRangeResponse(qm.Command, result)

	// Add the frame to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

/**
 * Iterate over xrange/xrevrange result and build new Frame with required fields
 */
func createFrameFromRangeResponse(command string, result []interface{}) *data.Frame {
	// Create new frame
	frame := data.NewFrame(command)

	// Create field to store entry id
	idField := data.NewField("$streamId", nil, []string{})

	// Add id field to the response
	frame.Fields = append(frame.Fields, idField)

	// Map for storing all the fields found in entries
	fields := map[string]*data.Field{}

	for _, entry := range result {
		id := string(entry.([]interface{})[0].([]byte))

		idField.Append(id)

		keysFoundInCurrentEntry := map[string]bool{}

		keyValuePairs := entry.([]interface{})[1].([]interface{})
		for i := 0; i < len(keyValuePairs); i += 2 {
			key := string(keyValuePairs[i].([]byte))
			value := string(keyValuePairs[i+1].([]byte))

			// Check if field has been already created before
			if _, ok := fields[key]; !ok {
				// Create new field
				newField := data.NewField(key, nil, []string{})
				fields[key] = newField

				// Append field to frame
				frame.Fields = append(frame.Fields, newField)

				// Get the number of rows we processed previously
				rowsCount := idField.Len() - 1

				// Generate empty values for all previous rows
				for j := 0; j < rowsCount; j++ {
					newField.Append("")
				}
			}

			// Insert value for current row
			fields[key].Append(value)
			keysFoundInCurrentEntry[key] = true
		}

		// Iterate over all keys found so far for stream
		for key, field := range fields {
			// Check if key exist in entry
			if _, ok := keysFoundInCurrentEntry[key]; !ok {
				// If key is missed in entry insert empty value
				field.Append("")
			}
		}
	}

	return frame
}

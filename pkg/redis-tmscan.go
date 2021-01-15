package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * TMSCAN cursor
 *
 * Iterates over the collection of keys and query type and memory usage
 * Cursor iteration similar to SCAN command
 * @see https://redis.io/commands/scan
 * @see https://redis.io/commands/type
 * @see https://redis.io/commands/memory-usage
 */
func queryTMScan(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute SCAN
	var result []interface{}
	cursor := "0"
	if qm.Cursor != "" {
		cursor = qm.Cursor
	}
	var args []interface{}
	if qm.Match != "" {
		args = append(args, "match", qm.Match)
	}
	if qm.Count != 0 {
		args = append(args, "count", qm.Count)
	}
	// Running CURSOR command
	err := client.RunFlatCmd(&result, "SCAN", cursor, args...)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}
	// New Frame
	frame := data.NewFrame(qm.Command)

	// Next cursor value is first value ([]byte) in result array see https://redis.io/commands/scan
	nextCursor := string(result[0].([]byte))
	// Add cursor field to frame
	frame.Fields = append(frame.Fields, data.NewField("cursor", nil, []string{nextCursor}))
	// Array with keys is second value in result array see https://redis.io/commands/scan
	keys := result[1].([]interface{})

	var typeCommands []flatCommandArgs
	var memoryCommands []flatCommandArgs
	// Slices with batch receiver pointers
	var typePointers []*string
	var memoryPointers []*int64
	// Slices with output values
	var names []string
	var types []string
	var memory []int64

	for _, key := range keys {
		name := string(key.([]byte))
		names = append(names, name)
		var keyType string
		var keyMemory int64
		typePointers = append(typePointers, &keyType)
		memoryPointers = append(memoryPointers, &keyMemory)
		typeCommands = append(typeCommands, flatCommandArgs{cmd: "TYPE", key: name, rcv: &keyType})
		memoryCommands = append(memoryCommands, flatCommandArgs{cmd: "MEMORY", key: "USAGE", args: []interface{}{name}, rcv: &keyMemory})
	}
	// Send batch with TYPE commands
	err = client.RunBatchFlatCmd(typeCommands)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}
	// Get the values stored by pointers and apply it to result slice
	for _, typePointer := range typePointers {
		types = append(types, *typePointer)
	}

	// Send batch with MEMORY USAGE commands
	err = client.RunBatchFlatCmd(memoryCommands)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Get the values stored by pointers and apply it to result slice
	for _, memoryPointer := range memoryPointers {
		memory = append(memory, *memoryPointer)
	}

	// Add key names field to frame
	frame.Fields = append(frame.Fields, data.NewField("key", nil, names))
	// Add key types field to frame
	frame.Fields = append(frame.Fields, data.NewField("type", nil, types))
	// Add key memory to frame with a proper config
	memoryField := data.NewField("memory", nil, memory)
	memoryField.Config = &data.FieldConfig{Unit: "decbytes"}
	frame.Fields = append(frame.Fields, memoryField)

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

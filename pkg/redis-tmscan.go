package main

import (
	"sort"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * TMSCAN cursor match count
 *
 * Iterates over the collection of keys and query type and memory usage
 * Cursor iteration similar to SCAN command
 * @see https://redis.io/commands/scan
 * @see https://redis.io/commands/type
 * @see https://redis.io/commands/memory-usage
 */
func queryTMScan(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	var result []interface{}

	// Cursor
	cursor := "0"
	if qm.Cursor != "" {
		cursor = qm.Cursor
	}

	// Match
	var args []interface{}
	if qm.Match != "" {
		args = append(args, "match", qm.Match)
	}

	// Count
	if qm.Count != 0 {
		args = append(args, "count", qm.Count)
	}

	// Running CURSOR command
	err := client.RunFlatCmd(&result, "SCAN", cursor, args...)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frames
	frame := data.NewFrame(qm.Command)
	frameCursor := data.NewFrame("Cursor")

	/**
	 * Next cursor value is first value ([]byte) in result array
	 * @see https://redis.io/commands/scan
	 */
	nextCursor := string(result[0].([]byte))

	// Add cursor field to frame
	frameCursor.Fields = append(frame.Fields, data.NewField("cursor", nil, []string{nextCursor}))

	/**
	 * Array with keys is second value in result array
	 * @see https://redis.io/commands/scan
	 */
	keys := result[1].([]interface{})

	var typeCommands []flatCommandArgs
	var memoryCommands []flatCommandArgs

	// Slices with output values
	var rows []*tmscanRow

	// Check memory usage for all keys
	for i, key := range keys {
		rows = append(rows, &tmscanRow{keyName: string(key.([]byte))})

		// Commands
		memoryCommandArgs := []interface{}{rows[i].keyName}
		if qm.Samples > 0 {
			memoryCommandArgs = append(memoryCommandArgs, "SAMPLES", qm.Samples)
		}
		memoryCommands = append(memoryCommands, flatCommandArgs{cmd: "MEMORY", key: "USAGE", args: memoryCommandArgs, rcv: &(rows[i].keyMemory)})
	}

	// Send batch with MEMORY USAGE commands
	err = client.RunBatchFlatCmd(memoryCommands)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Check if size is less than the number of rows and we need to select biggest keys
	if qm.Size > 0 && qm.Size < len(rows) {
		// Sort by memory usage
		sort.Slice(rows, func(i, j int) bool {
			// Use reversed condition for Descending sort
			return rows[i].keyMemory > rows[j].keyMemory
		})
		// Get first qm.Size keys
		rows = rows[:qm.Size]
	}

	// Check type for all keys
	for _, row := range rows {
		// Commands
		typeCommands = append(typeCommands, flatCommandArgs{cmd: "TYPE", key: row.keyName, rcv: &(row.keyType)})
	}

	// Send batch with TYPE commands
	err = client.RunBatchFlatCmd(typeCommands)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Add key names field to frame
	frame.Fields = append(frame.Fields, data.NewField("key", nil, []string{}))

	// Add key types field to frame
	frame.Fields = append(frame.Fields, data.NewField("type", nil, []string{}))

	// Add key ьemory to frame with a proper config
	memoryField := data.NewField("memory", nil, []int64{})
	memoryField.Config = &data.FieldConfig{Unit: "decbytes"}
	frame.Fields = append(frame.Fields, memoryField)

	// Append result rows to frame
	for _, row := range rows {
		frame.AppendRow(row.keyName, row.keyType, row.keyMemory)
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, frame, frameCursor)

	// Return
	return response
}

/**
 * TMSCAN result row entity
 */
type tmscanRow struct {
	keyName   string
	keyMemory int64
	keyType   string
}

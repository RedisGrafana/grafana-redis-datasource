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
	LastGeneratedId string        `redis:"last-generated-id"`
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
	frame.Fields = append(frame.Fields, data.NewField("last-generated-id", nil, []string{model.LastGeneratedId}))
	if model.FirstEntry != nil {
		frame.Fields = append(frame.Fields, data.NewField("first-entry-id", nil, []string{string(model.FirstEntry[0].([]byte))}))
		// Merging args to string like "key"="value"\n
		entryFields := model.FirstEntry[1].([]interface{})
		fields := new(bytes.Buffer)
		for i := 0; i < len(entryFields); i += 2 {
			field := string(entryFields[i].([]byte))
			value := string(entryFields[i+1].([]byte))
			fmt.Fprintf(fields, "\"%s\"=\"%s\"\n", field, value)
		}
		frame.Fields = append(frame.Fields, data.NewField("first-entry-fields", nil, []string{fields.String()}))
	}
	if model.LastEntry != nil {
		frame.Fields = append(frame.Fields, data.NewField("last-entry-id", nil, []string{string(model.LastEntry[0].([]byte))}))
		// Merging args to string like "key"="value"\n
		entryFields := model.LastEntry[1].([]interface{})
		fields := new(bytes.Buffer)
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

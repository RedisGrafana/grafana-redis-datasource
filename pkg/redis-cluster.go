package main

import (
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * CLUSTER INFO
 *
 * @see https://redis.io/commands/cluster-info
 */
func queryClusterInfo(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result string
	err := client.RunCmd(&result, "CLUSTER", "INFO")

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Split lines
	lines := strings.Split(strings.Replace(result, "\r\n", "\n", -1), "\n")

	// New Frame
	frame := data.NewFrame(qm.Command)

	// Parse lines
	for _, line := range lines {
		fields := strings.Split(line, ":")

		if len(fields) < 2 {
			continue
		}

		// Add Field
		if floatValue, err := strconv.ParseFloat(fields[1], 64); err == nil {
			frame.Fields = append(frame.Fields, data.NewField(fields[0], nil, []float64{floatValue}))
		} else {
			frame.Fields = append(frame.Fields, data.NewField(fields[0], nil, []string{fields[1]}))
		}
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

/**
 * CLUSTER NODES
 *
 * Line layout:
 *   <id> <ip:port@cport[,hostname[,aux=val]*]> <flags> <master> <ping-sent>
 *   <pong-recv> <config-epoch> <link-state> <slot> <slot> ... <slot>
 *
 * A node commonly owns several slot ranges and may be importing or migrating
 * slots, so every trailing field is kept instead of only the first one.
 *
 * @see https://redis.io/commands/cluster-nodes
 */
func queryClusterNodes(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result string
	err := client.RunCmd(&result, "CLUSTER", "NODES")

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Split lines
	lines := strings.Split(strings.Replace(result, "\r\n", "\n", -1), "\n")

	// New Frame
	frame := data.NewFrame(qm.Command,
		data.NewField("Id", nil, []string{}),
		data.NewField("Address", nil, []string{}),
		data.NewField("Hostname", nil, []string{}),
		data.NewField("Flags", nil, []string{}),
		data.NewField("Master", nil, []string{}),
		data.NewField("Ping", nil, []*int64{}),
		data.NewField("Pong", nil, []*int64{}),
		data.NewField("Epoch", nil, []int64{}),
		data.NewField("State", nil, []string{}),
		data.NewField("Slot", nil, []string{}),
		data.NewField("Slots", nil, []int64{}))

	/**
	 * ping-sent and pong-received are unix timestamps in milliseconds, not
	 * durations. 0 means "no ping currently pending", which is null rather than
	 * the epoch.
	 */
	frame.Fields[5].Config = &data.FieldConfig{Unit: "dateTimeAsIso"}
	frame.Fields[6].Config = &data.FieldConfig{Unit: "dateTimeAsIso"}

	// Parse lines
	for _, line := range lines {
		fields := strings.Fields(line)

		// Check number of fields
		if len(fields) < 8 {
			continue
		}

		// Parse values
		var ping *int64
		var pong *int64

		if value, err := strconv.ParseInt(fields[4], 10, 64); err == nil && value > 0 {
			ping = &value
		}

		if value, err := strconv.ParseInt(fields[5], 10, 64); err == nil && value > 0 {
			pong = &value
		}

		epoch, _ := strconv.ParseInt(fields[6], 10, 64)

		// Address carries an optional hostname and auxiliary fields
		address := fields[1]
		hostname := ""

		if parts := strings.Split(fields[1], ","); len(parts) > 1 {
			address = parts[0]
			hostname = parts[1]
		}

		// Every remaining field is a slot, a slot range, or a migration marker
		slots := fields[8:]
		served := countClusterSlots(slots)

		// Add Query
		frame.AppendRow(fields[0], address, hostname, fields[2], fields[3], ping, pong, epoch, fields[7], strings.Join(slots, " "), served)
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, frame)

	// Return
	return response
}

/**
 * countClusterSlots counts the slots served by a node.
 *
 * Handles "0-5460" ranges and bare "1234" slots; importing and migrating
 * markers such as "[1234-<-<node-id>]" are not served yet and are skipped.
 */
func countClusterSlots(slots []string) int64 {
	var total int64

	for _, slot := range slots {
		if strings.HasPrefix(slot, "[") {
			continue
		}

		bounds := strings.SplitN(slot, "-", 2)

		from, err := strconv.ParseInt(bounds[0], 10, 64)
		if err != nil {
			continue
		}

		if len(bounds) == 1 {
			total++
			continue
		}

		to, err := strconv.ParseInt(bounds[1], 10, 64)
		if err != nil {
			continue
		}

		if to >= from {
			total += to - from + 1
		}
	}

	return total
}

package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * HOTKEYS GET ( Redis >= 8.6 )
 *
 * Returns the top-K keys collected since HOTKEYS START, as a flat map:
 *
 *   tracking-active, sample-ratio, selected-slots, collection-start-time-unix-ms,
 *   collection-duration-ms, total-cpu-time-user-ms, total-cpu-time-sys-ms,
 *   total-net-bytes, by-cpu-time-us, by-net-bytes
 *
 * by-cpu-time-us and by-net-bytes are flat key/value arrays. They are merged
 * into one long frame so a single query drives a top-keys table, while the
 * scalars are returned as a separate summary frame.
 *
 * Requires the @admin ACL category; a read-only Grafana user will get NOPERM.
 *
 * @see https://redis.io/commands/hotkeys-get
 */
func queryHotkeysGet(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var result interface{}
	err := client.RunCmd(&result, "HOTKEYS", "GET")

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// Tracking was never started
	if result == nil {
		response.Frames = append(response.Frames, data.NewFrame("summary",
			data.NewField("tracking-active", nil, []int64{0})))
		return response
	}

	elements, ok := result.([]interface{})
	if !ok {
		response.Error = errUnexpectedHotkeysReply
		return response
	}

	summary := data.NewFrame("summary")

	// Key -> metric -> value, keeping the order keys were first seen in
	order := []string{}
	metrics := map[string]map[string]int64{}

	collect := func(metric string, values []interface{}) {
		for i := 0; i+1 < len(values); i += 2 {
			key := redisValueToString(values[i])

			amount, ok := redisValueToInt64(values[i+1])
			if !ok {
				continue
			}

			if _, seen := metrics[key]; !seen {
				order = append(order, key)
				metrics[key] = map[string]int64{}
			}

			metrics[key][metric] = amount
		}
	}

	// Walk the flat map
	for i := 0; i+1 < len(elements); i += 2 {
		name := redisValueToString(elements[i])
		value := elements[i+1]

		switch name {
		case "by-cpu-time-us", "by-net-bytes":
			if values, ok := value.([]interface{}); ok {
				collect(name, values)
			}
		case "selected-slots":
			// Array of [from, to] pairs, flattened to "0-16383, ..."
			summary.Fields = append(summary.Fields,
				data.NewField(name, nil, []string{formatSlotRanges(value)}))
		default:
			if amount, ok := redisValueToInt64(value); ok {
				summary.Fields = append(summary.Fields, data.NewField(name, nil, []int64{amount}))
			} else {
				summary.Fields = append(summary.Fields, data.NewField(name, nil, []string{redisValueToString(value)}))
			}
		}
	}

	keys := data.NewFrame("hotkeys",
		data.NewField("Key", nil, []string{}),
		data.NewField("CPU Time", nil, []*int64{}).SetConfig(&data.FieldConfig{Unit: "µs"}),
		data.NewField("Net Bytes", nil, []*int64{}).SetConfig(&data.FieldConfig{Unit: "bytes"}))

	for _, key := range order {
		var cpu *int64
		var net *int64

		if value, ok := metrics[key]["by-cpu-time-us"]; ok {
			cpu = &value
		}

		if value, ok := metrics[key]["by-net-bytes"]; ok {
			net = &value
		}

		keys.AppendRow(key, cpu, net)
	}

	response.Frames = append(response.Frames, summary, keys)

	// Return
	return response
}

/**
 * formatSlotRanges renders the selected-slots reply as "from-to, from-to".
 */
func formatSlotRanges(value interface{}) string {
	ranges, ok := value.([]interface{})
	if !ok {
		return redisValueToString(value)
	}

	formatted := ""

	for _, element := range ranges {
		bounds, ok := element.([]interface{})
		if !ok || len(bounds) < 2 {
			continue
		}

		if formatted != "" {
			formatted += ", "
		}

		formatted += redisValueToString(bounds[0]) + "-" + redisValueToString(bounds[1])
	}

	return formatted
}

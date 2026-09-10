package main

import (
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * Shared helpers to parse the INFO reply.
 *
 * The reply is a bulk string made of "# Section" headers and "key:value" lines.
 * Some sections encode a whole record on the value side as "a=1,b=2,c=3".
 *
 * @see https://redis.io/commands/info
 */

// infoSection is a single "# Name" block of an INFO reply.
type infoSection struct {
	Name  string
	Lines []string
}

/**
 * splitInfoLines normalizes line endings and returns the raw lines.
 */
func splitInfoLines(result string) []string {
	return strings.Split(strings.Replace(result, "\r\n", "\n", -1), "\n")
}

/**
 * splitInfoSections groups INFO lines by their "# Name" header.
 *
 * Lines emitted before any header (never the case for real servers, but cheap to
 * support) are collected under an empty section name.
 */
func splitInfoSections(lines []string) []infoSection {
	sections := []infoSection{}
	current := infoSection{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "#") {
			// Flush the previous block
			if current.Name != "" || len(current.Lines) > 0 {
				sections = append(sections, current)
			}

			current = infoSection{Name: strings.ToLower(strings.TrimSpace(strings.TrimPrefix(trimmed, "#")))}
			continue
		}

		if trimmed == "" {
			continue
		}

		current.Lines = append(current.Lines, trimmed)
	}

	// Flush the last block
	if current.Name != "" || len(current.Lines) > 0 {
		sections = append(sections, current)
	}

	return sections
}

/**
 * splitInfoLine splits "key:value" keeping any ":" that belongs to the value.
 *
 * Values legitimately contain colons (IPv6 endpoints in replication and listener
 * lines, absolute paths, "8.10.1 - oss" style versions), so strings.Split must
 * never be used here.
 */
func splitInfoLine(line string) (string, string, bool) {
	trimmed := strings.TrimSpace(line)

	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return "", "", false
	}

	fields := strings.SplitN(trimmed, ":", 2)
	if len(fields) < 2 {
		return "", "", false
	}

	return fields[0], fields[1], true
}

/**
 * splitTopLevel splits on sep, ignoring separators nested in "[...]".
 *
 * INFO modules joins its own lists with "|" inside brackets, but a bracket aware
 * split keeps the parser correct if that ever changes.
 */
func splitTopLevel(value string, sep byte) []string {
	parts := []string{}
	depth := 0
	start := 0

	for i := 0; i < len(value); i++ {
		switch value[i] {
		case '[':
			depth++
		case ']':
			if depth > 0 {
				depth--
			}
		case sep:
			if depth == 0 {
				parts = append(parts, value[start:i])
				start = i + 1
			}
		}
	}

	return append(parts, value[start:])
}

// infoPair is one "key=value" element of an INFO record, kept ordered.
type infoPair struct {
	Key   string
	Value string
}

/**
 * parseInfoPairs parses "a=1,b=2" into ordered key/value pairs.
 *
 * Elements without "=" are skipped, so a malformed field cannot shift the rest
 * of the record.
 */
func parseInfoPairs(value string) []infoPair {
	pairs := []infoPair{}

	for _, element := range splitTopLevel(value, ',') {
		element = strings.TrimSpace(element)
		if element == "" {
			continue
		}

		kv := strings.SplitN(element, "=", 2)
		if len(kv) < 2 {
			continue
		}

		pairs = append(pairs, infoPair{Key: kv[0], Value: kv[1]})
	}

	return pairs
}

// Column types used when a record section is turned into a long frame.
const (
	columnKindInt = iota
	columnKindFloat
	columnKindString
)

/**
 * columnKindOf reports the narrowest type able to hold value.
 */
func columnKindOf(value string) int {
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return columnKindInt
	}

	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return columnKindFloat
	}

	return columnKindString
}

// infoRecord is one parsed "name:a=1,b=2" line.
type infoRecord struct {
	Name  string
	Pairs []infoPair
}

/**
 * parseInfoRecords parses every "name:a=1,b=2" line of a section.
 *
 * trimPrefix is removed from the name, so "cmdstat_get" becomes "get" and
 * "latency_percentiles_usec_client|setinfo" becomes "client|setinfo".
 */
func parseInfoRecords(lines []string, trimPrefix string) []infoRecord {
	records := []infoRecord{}

	for _, line := range lines {
		name, value, ok := splitInfoLine(line)
		if !ok {
			continue
		}

		pairs := parseInfoPairs(value)
		if len(pairs) == 0 {
			continue
		}

		if trimPrefix != "" {
			name = strings.TrimPrefix(name, trimPrefix)
		}

		records = append(records, infoRecord{Name: name, Pairs: pairs})
	}

	return records
}

/**
 * recordsToFrame turns parsed records into a long frame.
 *
 * The first column holds the record name; one column is added per key seen in
 * any record, in first-seen order. A column is numeric only when every value
 * observed for it is numeric, so a single "n/a" cannot silently zero a series.
 * Missing values are null rather than 0.
 */
func recordsToFrame(frameName string, nameColumn string, records []infoRecord, unit string) *data.Frame {
	// Ordered union of the keys, plus the widest type required by each
	order := []string{}
	kinds := map[string]int{}

	for _, record := range records {
		for _, pair := range record.Pairs {
			kind := columnKindOf(pair.Value)

			previous, seen := kinds[pair.Key]
			if !seen {
				order = append(order, pair.Key)
				kinds[pair.Key] = kind
				continue
			}

			// Widen: int -> float -> string
			if kind > previous {
				kinds[pair.Key] = kind
			}
		}
	}

	frame := data.NewFrame(frameName, data.NewField(nameColumn, nil, []string{}))

	for _, key := range order {
		switch kinds[key] {
		case columnKindInt:
			field := data.NewField(key, nil, []*int64{})
			if unit != "" {
				field.Config = &data.FieldConfig{Unit: unit}
			}
			frame.Fields = append(frame.Fields, field)
		case columnKindFloat:
			field := data.NewField(key, nil, []*float64{})
			if unit != "" {
				field.Config = &data.FieldConfig{Unit: unit}
			}
			frame.Fields = append(frame.Fields, field)
		default:
			frame.Fields = append(frame.Fields, data.NewField(key, nil, []*string{}))
		}
	}

	// Rows
	for _, record := range records {
		values := map[string]string{}
		for _, pair := range record.Pairs {
			values[pair.Key] = pair.Value
		}

		row := make([]interface{}, 0, len(order)+1)
		row = append(row, record.Name)

		for _, key := range order {
			raw, present := values[key]

			switch kinds[key] {
			case columnKindInt:
				var parsed *int64
				if present {
					if value, err := strconv.ParseInt(raw, 10, 64); err == nil {
						parsed = &value
					}
				}
				row = append(row, parsed)
			case columnKindFloat:
				var parsed *float64
				if present {
					if value, err := strconv.ParseFloat(raw, 64); err == nil {
						parsed = &value
					}
				}
				row = append(row, parsed)
			default:
				var parsed *string
				if present {
					value := raw
					parsed = &value
				}
				row = append(row, parsed)
			}
		}

		frame.AppendRow(row...)
	}

	return frame
}

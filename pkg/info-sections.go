package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * INFO section names understood by the backend.
 *
 * The frontend only offers a subset in a dropdown, but the backend accepts any
 * string, so a dashboard may always ask for a section by name.
 */
const (
	sectionCommandstats = "commandstats"
	sectionErrorstats   = "errorstats"
	sectionLatencystats = "latencystats"
	sectionKeysizes     = "keysizes"
	sectionKeyspace     = "keyspace"
	sectionModules      = "modules"
	sectionThreads      = "threads"
	sectionEverything   = "everything"
	sectionAll          = "all"

	// Key shared by every line of the modules section
	moduleKey = "module"
)

/**
 * INFO latencystats ( Redis >= 7.0 )
 *
 * latency_percentiles_usec_<command>:p50=1.003,p99=2.007,p99.9=2.007
 *
 * The percentile list is configurable through latency-tracking-info-percentiles,
 * so the columns are discovered from the reply instead of being hardcoded.
 */
func infoLatencystatsFrame(name string, lines []string) *data.Frame {
	records := parseInfoRecords(lines, "latency_percentiles_usec_")
	return recordsToFrame(name, "Command", records, "µs")
}

/**
 * INFO threads ( Redis >= 8.0 )
 *
 * io_thread_0:clients=7,reads=48677,writes=44861
 */
func infoThreadsFrame(name string, lines []string) *data.Frame {
	records := parseInfoRecords(lines, "io_thread_")
	return recordsToFrame(name, "Thread", records, "")
}

/**
 * INFO keyspace
 *
 * db0:keys=426,expires=410,avg_ttl=2326254245,subexpiry=254
 *
 * subexpiry is Redis >= 7.4 (hash field expiration); columns are discovered so
 * newer counters show up without another plugin release.
 */
func infoKeyspaceFrame(name string, lines []string) *data.Frame {
	records := parseInfoRecords(lines, "")
	return recordsToFrame(name, "Database", records, "")
}

/**
 * INFO modules
 *
 * module:name=search,ver=81000,api=1,filters=0,usedby=[],using=[ReJSON],options=[handle-io-errors]
 *
 * Every line shares the "module" key, so the module name is promoted to the row
 * label; parsing it as a plain key/value section would produce one frame field
 * per module, all called "module".
 */
func infoModulesFrame(name string, lines []string) *data.Frame {
	records := []infoRecord{}

	for _, line := range lines {
		key, value, ok := splitInfoLine(line)
		if !ok || key != moduleKey {
			continue
		}

		moduleName := ""
		pairs := []infoPair{}

		for _, pair := range parseInfoPairs(value) {
			if pair.Key == "name" {
				moduleName = pair.Value
				continue
			}

			// Flatten the bracketed lists Redis joins with "|"
			if strings.HasPrefix(pair.Value, "[") && strings.HasSuffix(pair.Value, "]") {
				inner := strings.TrimSuffix(strings.TrimPrefix(pair.Value, "["), "]")
				pair.Value = strings.Join(strings.Split(inner, "|"), ", ")
			}

			pairs = append(pairs, pair)
		}

		if moduleName == "" {
			continue
		}

		records = append(records, infoRecord{Name: moduleName, Pairs: pairs})
	}

	return recordsToFrame(name, "Name", records, "")
}

/**
 * INFO keysizes ( Redis >= 8.0 )
 *
 * db0_distrib_strings_sizes:1=2,2=3,4=1
 * db0_distrib_lists_items:1=27
 *
 * Buckets are power-of-two lower bounds: bucket 4 counts keys of size 4 to 7.
 * The section is emitted as a long frame so a single query can drive a
 * histogram per database and data type.
 */
func infoKeysizesFrame(name string, lines []string) *data.Frame {
	frame := data.NewFrame(name,
		data.NewField("Database", nil, []string{}),
		data.NewField("Type", nil, []string{}),
		data.NewField("Metric", nil, []string{}),
		data.NewField("Bucket", nil, []int64{}),
		data.NewField("Range", nil, []string{}),
		data.NewField("Count", nil, []int64{}))

	for _, line := range lines {
		key, value, ok := splitInfoLine(line)
		if !ok {
			continue
		}

		// db0_distrib_strings_sizes -> db0 / strings / sizes
		parts := strings.SplitN(key, "_distrib_", 2)
		if len(parts) < 2 {
			continue
		}

		database := parts[0]
		dataType := parts[1]
		metric := ""

		if index := strings.LastIndex(parts[1], "_"); index > 0 {
			dataType = parts[1][:index]
			metric = parts[1][index+1:]
		}

		for _, pair := range parseInfoPairs(value) {
			bucket, err := strconv.ParseInt(pair.Key, 10, 64)
			if err != nil {
				continue
			}

			count, err := strconv.ParseInt(pair.Value, 10, 64)
			if err != nil {
				continue
			}

			frame.AppendRow(database, dataType, metric, bucket, keysizesRange(bucket), count)
		}
	}

	return frame
}

/**
 * keysizesRange renders the closed interval covered by a power-of-two bucket.
 */
func keysizesRange(bucket int64) string {
	if bucket < 1 {
		return strconv.FormatInt(bucket, 10)
	}

	upper := bucket*2 - 1
	if upper == bucket {
		return strconv.FormatInt(bucket, 10)
	}

	return fmt.Sprintf("%d-%d", bucket, upper)
}

/**
 * infoGenericFrame builds the historical wide frame: one row, one field per key.
 *
 * Two fixes over a naive implementation: the value keeps every ":" it contains
 * (IPv6 endpoints, paths, "8.10.1 - oss"), and repeated keys are suffixed so the
 * frame never carries two fields with the same name.
 */
func infoGenericFrame(name string, lines []string) *data.Frame {
	frame := data.NewFrame(name)
	seen := map[string]int{}

	// Disambiguate repeated keys instead of emitting duplicate fields
	unique := func(key string) string {
		if count, exists := seen[key]; exists {
			seen[key] = count + 1
			return fmt.Sprintf("%s_%d", key, count)
		}

		seen[key] = 1
		return key
	}

	appendValue := func(key string, value string) {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			frame.Fields = append(frame.Fields, data.NewField(unique(key), nil, []float64{floatValue}))
			return
		}

		frame.Fields = append(frame.Fields, data.NewField(unique(key), nil, []string{value}))
	}

	for _, line := range lines {
		key, value, ok := splitInfoLine(line)
		if !ok {
			continue
		}

		appendValue(key, value)

		/**
		 * Some sections hide records inside a scalar line, most importantly
		 *
		 *   slave0:ip=10.0.0.2,port=6379,state=online,offset=2212222,lag=0
		 *   search_fields_text:Text=2,IndexErrors=0
		 *
		 * Upstream returned those as one opaque string, so replication offset,
		 * lag and per field index errors could not be charted at all. The
		 * combined string is kept for dashboards that already reference it and
		 * every member is added as its own <key>_<member> field.
		 *
		 * module lines are left alone. A server can load a dozen modules, each
		 * with seven members, which would bury the rest of INFO under expanded
		 * columns, and the modules section already parses them into a frame of
		 * one row per module.
		 */
		if key != moduleKey {
			for _, pair := range recordValuePairs(value) {
				appendValue(key+"_"+pair.Key, pair.Value)
			}
		}
	}

	return frame
}

/**
 * Returns the members of a value shaped exactly like k=v,k=v, and nil for
 * anything else, so ordinary values such as mem_allocator:jemalloc-5.3.0 or
 * master_failover_state:no-failover are never expanded.
 */
func recordValuePairs(value string) []infoPair {
	if !strings.Contains(value, "=") {
		return nil
	}

	elements := splitTopLevel(value, ',')
	pairs := make([]infoPair, 0, len(elements))

	for _, element := range elements {
		kv := strings.SplitN(strings.TrimSpace(element), "=", 2)

		// Every member has to be a key=value pair, otherwise this is a plain string
		if len(kv) < 2 || kv[0] == "" || strings.Contains(kv[1], "=") {
			return nil
		}

		pairs = append(pairs, infoPair{Key: kv[0], Value: kv[1]})
	}

	return pairs
}

/**
 * infoSectionFrame routes one INFO section to the parser that fits its shape.
 *
 * name is used as the frame name so a multi-section reply (INFO all/everything)
 * produces frames a panel can tell apart.
 */
func infoSectionFrame(name string, section string, lines []string, streaming bool) *data.Frame {
	switch section {
	case sectionCommandstats:
		return infoCommandstatsFrame(name, lines)
	case sectionErrorstats:
		return infoErrorstatsFrame(name, lines, streaming)
	case sectionLatencystats:
		return infoLatencystatsFrame(name, lines)
	case sectionKeysizes:
		return infoKeysizesFrame(name, lines)
	case sectionKeyspace:
		return infoKeyspaceFrame(name, lines)
	case sectionModules:
		return infoModulesFrame(name, lines)
	case sectionThreads:
		return infoThreadsFrame(name, lines)
	default:
		return infoGenericFrame(name, lines)
	}
}

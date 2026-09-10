package main

import (
	"testing"

	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

/**
 * Payloads captured from Redis 8.10.1.
 */
const (
	latencystatsPayload = "# Latencystats\r\n" +
		"latency_percentiles_usec_hincrby:p50=13.055,p99=21.119,p99.9=21.119\r\n" +
		"latency_percentiles_usec_client|setinfo:p50=1.003,p99=2.007,p99.9=2.007\r\n"

	keysizesPayload = "# Keysizes\r\n" +
		"db0_distrib_strings_sizes:1=2,2=3,4=1\r\n" +
		"db0_distrib_lists_items:1=27\r\n" +
		"db0_distrib_sets_items:2=9,4=4,8=5,16=2,128=1\r\n"

	modulesPayload = "# Modules\r\n" +
		"module:name=vectorset,ver=1,api=1,filters=0,usedby=[],using=[],options=[handle-io-errors|handle-repl-async-load]\r\n" +
		"module:name=ReJSON,ver=81000,api=1,filters=0,usedby=[search],using=[],options=[handle-io-errors]\r\n"

	threadsPayload = "# Threads\r\nio_thread_0:clients=7,reads=48677,writes=44861\r\n"

	keyspacePayload = "# Keyspace\r\ndb0:keys=426,expires=410,avg_ttl=2326254245,subexpiry=254\r\n"
)

/**
 * INFO latencystats
 */
func TestInfoLatencystatsFrame(t *testing.T) {
	t.Parallel()

	frame := infoLatencystatsFrame(models.Info, splitInfoLines(latencystatsPayload))

	require.Len(t, frame.Fields, 4, "Command plus one field per percentile")
	require.Equal(t, "Command", frame.Fields[0].Name)
	require.Equal(t, "p50", frame.Fields[1].Name)
	require.Equal(t, "p99", frame.Fields[2].Name)
	require.Equal(t, "p99.9", frame.Fields[3].Name)
	require.Equal(t, "µs", frame.Fields[1].Config.Unit, "Percentiles are microseconds")

	require.Equal(t, 2, frame.Fields[0].Len())
	require.Equal(t, "hincrby", frame.Fields[0].At(0))
	require.Equal(t, 13.055, *frame.Fields[1].At(0).(*float64))

	// Container commands keep their "|" separator
	require.Equal(t, "client|setinfo", frame.Fields[0].At(1))
	require.Equal(t, 2.007, *frame.Fields[3].At(1).(*float64))
}

/**
 * INFO keysizes
 */
func TestInfoKeysizesFrame(t *testing.T) {
	t.Parallel()

	frame := infoKeysizesFrame(models.Info, splitInfoLines(keysizesPayload))

	require.Len(t, frame.Fields, 6)
	require.Equal(t, 9, frame.Fields[0].Len(), "One row per bucket across every distribution")

	require.Equal(t, "db0", frame.Fields[0].At(0))
	require.Equal(t, "strings", frame.Fields[1].At(0))
	require.Equal(t, "sizes", frame.Fields[2].At(0))
	require.Equal(t, int64(1), frame.Fields[3].At(0))
	require.Equal(t, "1", frame.Fields[4].At(0))
	require.Equal(t, int64(2), frame.Fields[5].At(0))

	// Power of two buckets are half open ranges
	require.Equal(t, int64(4), frame.Fields[3].At(2))
	require.Equal(t, "4-7", frame.Fields[4].At(2))

	require.Equal(t, "lists", frame.Fields[1].At(3))
	require.Equal(t, "items", frame.Fields[2].At(3))

	require.Equal(t, "sets", frame.Fields[1].At(8))
	require.Equal(t, "128-255", frame.Fields[4].At(8))
}

/**
 * INFO modules
 */
func TestInfoModulesFrame(t *testing.T) {
	t.Parallel()

	frame := infoModulesFrame(models.Info, splitInfoLines(modulesPayload))

	// Name, ver, api, filters, usedby, using, options
	require.Len(t, frame.Fields, 7)
	require.Equal(t, "Name", frame.Fields[0].Name)
	require.Equal(t, 2, frame.Fields[0].Len(), "Every line must produce a row, not a field")

	require.Equal(t, "vectorset", frame.Fields[0].At(0))
	require.Equal(t, int64(81000), *frame.Fields[1].At(1).(*int64))
	require.Equal(t, "search", *frame.Fields[4].At(1).(*string), "usedby loses its brackets")
	require.Equal(t, "handle-io-errors, handle-repl-async-load", *frame.Fields[6].At(0).(*string))
}

/**
 * INFO threads
 */
func TestInfoThreadsFrame(t *testing.T) {
	t.Parallel()

	frame := infoThreadsFrame(models.Info, splitInfoLines(threadsPayload))

	require.Len(t, frame.Fields, 4)
	require.Equal(t, "Thread", frame.Fields[0].Name)
	require.Equal(t, "0", frame.Fields[0].At(0))
	require.Equal(t, int64(48677), *frame.Fields[2].At(0).(*int64))
}

/**
 * INFO keyspace
 */
func TestInfoKeyspaceFrame(t *testing.T) {
	t.Parallel()

	frame := infoKeyspaceFrame(models.Info, splitInfoLines(keyspacePayload))

	require.Len(t, frame.Fields, 5, "Database plus keys, expires, avg_ttl and subexpiry")
	require.Equal(t, "db0", frame.Fields[0].At(0))
	require.Equal(t, int64(426), *frame.Fields[1].At(0).(*int64))
	require.Equal(t, int64(254), *frame.Fields[4].At(0).(*int64), "subexpiry is Redis 7.4+")
}

/**
 * A record section with columns missing from some rows
 */
func TestRecordsToFrameSparseColumns(t *testing.T) {
	t.Parallel()

	payload := "a:x=1,y=2\r\nb:x=3\r\nc:y=4,z=n/a\r\n"
	frame := recordsToFrame(models.Info, "Name", parseInfoRecords(splitInfoLines(payload), ""), "")

	require.Len(t, frame.Fields, 4)
	require.Equal(t, int64(1), *frame.Fields[1].At(0).(*int64))
	require.Equal(t, int64(3), *frame.Fields[1].At(1).(*int64))
	require.Nil(t, frame.Fields[1].At(2), "A missing value is null, not zero")
	require.Nil(t, frame.Fields[2].At(1))

	// A single non numeric value makes the whole column text
	require.Equal(t, "n/a", *frame.Fields[3].At(2).(*string))
	require.Nil(t, frame.Fields[3].At(0))
}

/**
 * The generic parser must not lose colons or drop repeated keys
 */
func TestInfoGenericFrame(t *testing.T) {
	t.Parallel()

	payload := "# Server\r\n" +
		"redis_version:8.10.1\r\n" +
		"listener0:name=tcp,bind=127.0.0.1,port=6379\r\n" +
		"slave0:ip=::1,port=6380,state=online\r\n" +
		"search_redis_version:8.10.1 - oss\r\n" +
		"module:name=bf\r\n" +
		"module:name=search\r\n"

	frame := infoGenericFrame(models.Info, splitInfoLines(payload))

	names := []string{}
	for _, field := range frame.Fields {
		names = append(names, field.Name)
	}

	require.Equal(t, []string{
		"redis_version",
		"listener0", "listener0_name", "listener0_bind", "listener0_port",
		"slave0", "slave0_ip", "slave0_port", "slave0_state",
		"search_redis_version",
		"module", "module_1",
	}, names, "records are expanded in place, modules are left to the modules section")

	require.Equal(t, "8.10.1", frame.Fields[0].At(0), "A three part version is not a float")

	require.Equal(t, "ip=::1,port=6380,state=online", frame.Fields[5].At(0), "IPv6 keeps its colons")
	require.Equal(t, "::1", frame.Fields[6].At(0), "an IPv6 address is not split on its colons")
	require.Equal(t, "8.10.1 - oss", frame.Fields[9].At(0))

	// Repeated keys are disambiguated instead of colliding
	require.Equal(t, "name=bf", frame.Fields[10].At(0))
	require.Equal(t, "name=search", frame.Fields[11].At(0))
}

/**
 * INFO all is returned as one frame per section
 */
func TestQueryInfoEverything(t *testing.T) {
	t.Parallel()

	payload := "# Server\r\nredis_version:8.10.1\r\n\r\n" +
		"# Keyspace\r\ndb0:keys=426,expires=410\r\n\r\n" +
		"# Latencystats\r\nlatency_percentiles_usec_get:p50=4.015\r\n"

	client := testClient{rcv: payload}
	response := queryInfo(queryModel{Command: models.Info, Section: "everything"}, &client)

	require.NoError(t, response.Error)
	require.Len(t, response.Frames, 3, "One frame per section")

	require.Equal(t, "server", response.Frames[0].Name)
	require.Equal(t, "8.10.1", response.Frames[0].Fields[0].At(0))

	// Sections that carry records stay parsed instead of collapsing to text
	require.Equal(t, "keyspace", response.Frames[1].Name)
	require.Equal(t, "Database", response.Frames[1].Fields[0].Name)

	require.Equal(t, "latencystats", response.Frames[2].Name)
	require.Equal(t, "Command", response.Frames[2].Fields[0].Name)
	require.Equal(t, "get", response.Frames[2].Fields[0].At(0))
}

/**
 * Replication records
 */
func TestInfoGenericFrameExpandsRecordValues(t *testing.T) {
	t.Parallel()

	lines := splitInfoLines("# Replication\r\n" +
		"role:master\r\n" +
		"connected_slaves:1\r\n" +
		"slave0:ip=10.0.0.2,port=6379,state=online,offset=2212222,lag=0\r\n" +
		"master_failover_state:no-failover\r\n" +
		"master_replid:97463b446cd183cb8141e2349f4a07e8bcafa480\r\n")

	frame := infoGenericFrame("info", lines)

	fields := map[string]interface{}{}
	for _, field := range frame.Fields {
		fields[field.Name] = field.At(0)
	}

	require.Equal(t, "master", fields["role"], "role stays a string")
	require.Equal(t, float64(1), fields["connected_slaves"], "connected_slaves stays numeric")
	require.Equal(t, "ip=10.0.0.2,port=6379,state=online,offset=2212222,lag=0", fields["slave0"],
		"the combined value is kept for dashboards that already use it")
	require.Equal(t, "10.0.0.2", fields["slave0_ip"], "the replica address is exposed on its own")
	require.Equal(t, float64(6379), fields["slave0_port"], "the replica port is numeric")
	require.Equal(t, "online", fields["slave0_state"], "the replica state is exposed on its own")
	require.Equal(t, float64(2212222), fields["slave0_offset"], "the replication offset can be charted")
	require.Equal(t, float64(0), fields["slave0_lag"], "the replication lag can be charted")

	// Values that merely look unusual must not be split apart
	require.Equal(t, "no-failover", fields["master_failover_state"], "a dashed value is left alone")
	require.NotContains(t, fields, "master_replid_", "a plain identifier is left alone")
}

/**
 * INFO search ( Redis >= 8.0, query engine )
 */
func TestInfoSearchFrame(t *testing.T) {
	t.Parallel()

	// Captured from redis 8.10.1 with one index present
	lines := splitInfoLines("# search_version\r\n" +
		"search_version:8.10.0\r\n" +
		"search_redis_version:8.10.1 - oss\r\n" +
		"\r\n" +
		"# search_index\r\n" +
		"search_number_of_indexes:1\r\n" +
		"\r\n" +
		"# search_fields_statistics\r\n" +
		"search_fields_text:Text=2,IndexErrors=0\r\n" +
		"search_fields_numeric:Numeric=1,IndexErrors=0\r\n" +
		"\r\n" +
		"# search_runtime_configurations\r\n" +
		"search_extension_load:\r\n" +
		"search_default_scorer:BM25STD\r\n")

	frame := infoSectionFrame("info", "search", lines, false)

	fields := map[string]interface{}{}
	for _, field := range frame.Fields {
		fields[field.Name] = field.At(0)
	}

	require.Equal(t, "8.10.0", fields["search_version"], "the module version is a string")
	require.Equal(t, "8.10.1 - oss", fields["search_redis_version"], "the build suffix is preserved")
	require.Equal(t, float64(1), fields["search_number_of_indexes"], "index counters are numeric")
	require.Equal(t, float64(2), fields["search_fields_text_Text"], "field statistics records are expanded")
	require.Equal(t, float64(0), fields["search_fields_text_IndexErrors"], "index errors can be alerted on")
	require.Equal(t, float64(1), fields["search_fields_numeric_Numeric"], "every field type is expanded")
	require.Equal(t, "", fields["search_extension_load"], "an empty configuration value stays empty")
	require.Equal(t, "BM25STD", fields["search_default_scorer"], "sub section headers are skipped")
}

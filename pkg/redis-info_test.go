package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQueryInfo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		qm                      queryModel
		rcv                     interface{}
		fieldsCount             int
		rowsPerField            int
		valuesToCheckInResponse []valueToCheckInResponse
		err                     error
	}{
		{
			"should parse default bulk string",
			queryModel{Command: "info"},
			"# Server\r\nredis_version:6.0.1\r\nredis_git_sha1:00000000\r\nredis_git_dirty:0\r\nredis_build_id:e02d1d807e41d65\r\nredis_mode:standalone\r\nos:Linux 5.4.0-58-generic x86_64\r\narch_bits:64\r\nmultiplexing_api:epoll\r\natomicvar_api:atomic-builtin\r\ngcc_version:8.3.0\r\nprocess_id:1\r\nrun_id:63886ab60ce4a06f9b5c2101dd08b0a1c3099a64\r\ntcp_port:6379\r\nuptime_in_seconds:224\r\nuptime_in_days:0\r\nhz:10\r\nconfigured_hz:10\r\nlru_clock:15846554\r\nexecutable:/data/redis-server\r\nconfig_file:\r\n\r\n# Clients\r\nconnected_clients:2\r\nclient_recent_max_input_buffer:0\r\nclient_recent_max_output_buffer:0\r\nblocked_clients:0\r\ntracking_clients:0\r\nclients_in_timeout_table:0\r\n\r\n# Memory\r\nused_memory:5377000\r\nused_memory_human:5.13M\r\nused_memory_rss:22200320\r\nused_memory_rss_human:21.17M\r\nused_memory_peak:5377000\r\nused_memory_peak_human:5.13M\r\nused_memory_peak_perc:101.60%\r\nused_memory_overhead:5292568\r\nused_memory_startup:5292568\r\nused_memory_dataset:84432\r\nused_memory_dataset_perc:100.00%\r\nallocator_allocated:5588768\r\nallocator_active:5984256\r\nallocator_resident:8859648\r\ntotal_system_memory:33509154816\r\ntotal_system_memory_human:31.21G\r\nused_memory_lua:37888\r\nused_memory_lua_human:37.00K\r\nused_memory_scripts:0\r\nused_memory_scripts_human:0B\r\nnumber_of_cached_scripts:0\r\nmaxmemory:0\r\nmaxmemory_human:0B\r\nmaxmemory_policy:noeviction\r\nallocator_frag_ratio:1.07\r\nallocator_frag_bytes:395488\r\nallocator_rss_ratio:1.48\r\nallocator_rss_bytes:2875392\r\nrss_overhead_ratio:2.51\r\nrss_overhead_bytes:13340672\r\nmem_fragmentation_ratio:4.19\r\nmem_fragmentation_bytes:16907752\r\nmem_not_counted_for_evict:0\r\nmem_replication_backlog:0\r\nmem_clients_slaves:0\r\nmem_clients_normal:0\r\nmem_aof_buffer:0\r\nmem_allocator:jemalloc-5.1.0\r\nactive_defrag_running:0\r\nlazyfree_pending_objects:0\r\n\r\n# Persistence\r\nloading:0\r\nrdb_changes_since_last_save:0\r\nrdb_bgsave_in_progress:0\r\nrdb_last_save_time:1609681850\r\nrdb_last_bgsave_status:ok\r\nrdb_last_bgsave_time_sec:-1\r\nrdb_current_bgsave_time_sec:-1\r\nrdb_last_cow_size:0\r\naof_enabled:0\r\naof_rewrite_in_progress:0\r\naof_rewrite_scheduled:0\r\naof_last_rewrite_time_sec:-1\r\naof_current_rewrite_time_sec:-1\r\naof_last_bgrewrite_status:ok\r\naof_last_write_status:ok\r\naof_last_cow_size:0\r\nmodule_fork_in_progress:0\r\nmodule_fork_last_cow_size:0\r\n\r\n# Stats\r\ntotal_connections_received:2\r\ntotal_commands_processed:5\r\ninstantaneous_ops_per_sec:0\r\ntotal_net_input_bytes:14\r\ntotal_net_output_bytes:0\r\ninstantaneous_input_kbps:0.00\r\ninstantaneous_output_kbps:0.00\r\nrejected_connections:0\r\nsync_full:0\r\nsync_partial_ok:0\r\nsync_partial_err:0\r\nexpired_keys:0\r\nexpired_stale_perc:0.00\r\nexpired_time_cap_reached_count:0\r\nexpire_cycle_cpu_milliseconds:6\r\nevicted_keys:0\r\nkeyspace_hits:0\r\nkeyspace_misses:0\r\npubsub_channels:0\r\npubsub_patterns:0\r\nlatest_fork_usec:0\r\nmigrate_cached_sockets:0\r\nslave_expires_tracked_keys:0\r\nactive_defrag_hits:0\r\nactive_defrag_misses:0\r\nactive_defrag_key_hits:0\r\nactive_defrag_key_misses:0\r\ntracking_total_keys:0\r\ntracking_total_items:0\r\nunexpected_error_replies:0\r\n\r\n# Replication\r\nrole:master\r\nconnected_slaves:0\r\nmaster_replid:9a0cb33c79e465dad8b2468c958a523d82f6b572\r\nmaster_replid2:0000000000000000000000000000000000000000\r\nmaster_repl_offset:0\r\nmaster_repl_meaningful_offset:0\r\nsecond_repl_offset:-1\r\nrepl_backlog_active:0\r\nrepl_backlog_size:1048576\r\nrepl_backlog_first_byte_offset:0\r\nrepl_backlog_histlen:0\r\n\r\n# CPU\r\nused_cpu_sys:0.393216\r\nused_cpu_user:0.403734\r\nused_cpu_sys_children:0.000000\r\nused_cpu_user_children:0.000000\r\n\r\n# Modules\r\nmodule:name=ReJSON,ver=10007,api=1,filters=0,usedby=[],using=[],options=[]\r\nmodule:name=search,ver=20005,api=1,filters=0,usedby=[],using=[],options=[]\r\nmodule:name=graph,ver=20212,api=1,filters=0,usedby=[],using=[],options=[]\r\nmodule:name=rg,ver=10003,api=1,filters=0,usedby=[],using=[ai],options=[]\r\nmodule:name=timeseries,ver=10407,api=1,filters=0,usedby=[],using=[],options=[]\r\nmodule:name=bf,ver=20205,api=1,filters=0,usedby=[],using=[],options=[]\r\nmodule:name=ai,ver=10002,api=1,filters=0,usedby=[rg],using=[],options=[]\r\n\r\n# Cluster\r\ncluster_enabled:0\r\n\r\n# Keyspace\r\n",
			137,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "6.0.1"},
				{frameIndex: 0, fieldIndex: 6, rowIndex: 0, value: float64(64)},
				{frameIndex: 0, fieldIndex: 2, rowIndex: 0, value: float64(0)},
			},
			nil,
		},
		{
			"should parse bulk string with 'commandstats' section",
			queryModel{Command: "info", Section: "commandstats"},
			"# Commandstats\r\ncmdstat_info:calls=5,usec=203,usec_per_call=40.60\r\ncmdstat_config:calls=1,usec=29,usec_per_call=29.00\r\n",
			4,
			2,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "info"},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: int64(5)},
				{frameIndex: 0, fieldIndex: 2, rowIndex: 0, value: float64(203)},
				{frameIndex: 0, fieldIndex: 3, rowIndex: 0, value: float64(40.60)},
			},
			nil,
		},
		{
			"should parse bulk string with 'commandstats' and ignore stats if only 2 stats per command",
			queryModel{Command: "info", Section: "commandstats"},
			"# Commandstats\r\ncmdstat_info:calls=5,usec_per_call=40.60\r\ncmdstat_config:calls=1,usec_per_call=29.00\r\n",
			4,
			0,
			nil,
			nil,
		},
		{
			"should parse bulk string with 'commandstats' section ans streaming true",
			queryModel{Command: "info", Section: "commandstats", Streaming: true},
			"# Commandstats\r\ncmdstat_info:calls=5,usec=203,usec_per_call=40.60\r\ncmdstat_config:calls=1,usec=29,usec_per_call=29.00\r\n",
			2,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: int64(5)},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: int64(1)},
			},
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "info"},
			nil,
			0,
			0,
			nil,
			errors.New("error occurred"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ds := redisDatasource{}
			client := TestClient{tt.rcv, tt.err}
			response := ds.queryInfo(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Command, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created from bulk string")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")
				require.NoError(t, response.Error, "Should not return error")
				if tt.valuesToCheckInResponse != nil {
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				}
			}
		})
	}
}

func TestQueryClientList(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		qm                      queryModel
		rcv                     interface{}
		fieldsCount             int
		rowsPerField            int
		valuesToCheckInResponse []valueToCheckInResponse
		err                     error
	}{
		{
			"should parse default bulk string",
			queryModel{Command: "clientList"},
			"id=81 addr=172.18.0.1:33504 fd=13 name= age=0 idle=0 flags=N db=0 sub=0 psub=0 multi=-1 qbuf=26 qbuf-free=32742 obl=0 oll=0 omem=0 events=r cmd=client user=default\nid=82 addr=172.18.0.1:33508 fd=14 name= age=0 idle=0 flags=N db=0 sub=0 psub=0 multi=-1 qbuf=0 qbuf-free=0 obl=0 oll=0 omem=0 events=r cmd=NULL user=default\n",
			19,
			2,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: int64(81)},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: "172.18.0.1:33504"},
			},
			nil,
		},
		{
			"should parse default bulk string and ignore elements without the =",
			queryModel{Command: "clientList"},
			"id=81 dummy addr=172.18.0.1:33504 fd=13 name= age=0 idle=0 flags=N db=0 sub=0 psub=0 multi=-1 qbuf=26 qbuf-free=32742 obl=0 oll=0 omem=0 events=r cmd=client user=default\nid=82 dummy addr=172.18.0.1:33508 fd=14 name= age=0 idle=0 flags=N db=0 sub=0 psub=0 multi=-1 qbuf=0 qbuf-free=0 obl=0 oll=0 omem=0 events=r cmd=NULL user=default\n",
			19,
			2,
			nil,
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "clientList"},
			nil,
			0,
			0,
			nil,
			errors.New("error occurred"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ds := redisDatasource{}
			client := TestClient{tt.rcv, tt.err}
			response := ds.queryClientList(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Command, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created from bulk string")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")
				require.NoError(t, response.Error, "Should not return error")
				if tt.valuesToCheckInResponse != nil {
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				}
			}
		})
	}
}

func TestQuerySlowlogGet(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		qm                      queryModel
		rcv                     interface{}
		fieldsCount             int
		rowsPerField            int
		valuesToCheckInResponse []valueToCheckInResponse
		err                     error
	}{
		{
			"should parse payload for redis prior to version 4.0",
			queryModel{Command: "slowlogGet"},
			[]interface{}{
				[]interface{}{int64(14), int64(1309448221), int64(15), []interface{}{"ping"}},
				[]interface{}{int64(13), int64(1309448128), int64(30), []interface{}{"slowlog", "get", "100"}},
			},
			4,
			2,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: int64(14)},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: time.Unix(1309448221, 0)},
				{frameIndex: 0, fieldIndex: 2, rowIndex: 0, value: int64(15)},
				{frameIndex: 0, fieldIndex: 3, rowIndex: 1, value: "slowlog get 100"},
			},
			nil,
		},
		{
			"should parse payload for redis starting version 4.0",
			queryModel{Command: "slowlogGet"},
			[]interface{}{
				[]interface{}{int64(14), int64(1309448221), int64(15), []interface{}{"ping"}, "127.0.0.1:58217", "worker-123"},
				[]interface{}{int64(13), int64(1309448128), int64(30), []interface{}{"slowlog", "get", "100"}, "127.0.0.1:58217", "worker-123"},
			},
			4,
			2,
			nil,
			nil,
		},
		{
			"should parse payload with array of command arguments having specific types",
			queryModel{Command: "slowlogGet"},
			[]interface{}{
				[]interface{}{int64(14), int64(1309448221), int64(15), []interface{}{"ping", int32(3), []byte("pong"), []interface{}{}}, "127.0.0.1:58217", "worker-123"},
			},
			4,
			1,
			nil,
			nil,
		},
		{
			"should parse payload with size provided",
			queryModel{Command: "slowlogGet", Size: 2},
			[]interface{}{
				[]interface{}{int64(14), int64(1309448221), int64(15), []interface{}{"ping"}, "127.0.0.1:58217", "worker-123"},
			},
			4,
			1,
			nil,
			nil,
		},
		{
			"should parse payload for redis Enterprise (command is in 5th field)",
			queryModel{Command: "slowlogGet"},
			[]interface{}{
				[]interface{}{int64(14), int64(1309448221), int64(15), "N:886,M:885", []interface{}{"ping"}},
				[]interface{}{int64(13), int64(1309448128), int64(30), "N:886,M:885", []interface{}{"slowlog", "get", "100"}},
			},
			4,
			2,
			nil,
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "slowlogGet"},
			nil,
			0,
			0,
			nil,
			errors.New("error occurred"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ds := redisDatasource{}
			client := TestClient{tt.rcv, tt.err}
			response := ds.querySlowlogGet(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Command, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created from bulk string")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")
				require.NoError(t, response.Error, "Should not return error")
				if tt.valuesToCheckInResponse != nil {
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				}
			}
		})
	}
}

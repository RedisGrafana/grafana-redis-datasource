package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryClusterInfo(t *testing.T) {
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
			"should parse clusterInfo bulk string",
			queryModel{Command: "clusterInfo"},
			"cluster_state:ok\r\ncluster_slots_assigned:16384\r\ncluster_slots_ok:16384\r\ncluster_slots_pfail:0\r\ncluster_slots_fail:0\r\ncluster_known_nodes:6\r\ncluster_size:3\r\ncluster_current_epoch:6\r\ncluster_my_epoch:2\r\ncluster_stats_messages_sent:1483972\r\ncluster_stats_messages_received:1483968",
			11,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "ok"},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: float64(16384)},
				{frameIndex: 0, fieldIndex: 6, rowIndex: 0, value: float64(3)},
			},
			nil,
		},
		{
			"should parse string and ignore non-pairing param",
			queryModel{Command: "clusterInfo"},
			"cluster_state:ok\r\ncluster_slots_assigned\r\ncluster_slots_ok:16384\r\ncluster_slots_pfail:0\r\ncluster_slots_fail:0\r\ncluster_known_nodes:6\r\ncluster_size:3\r\ncluster_current_epoch:6\r\ncluster_my_epoch:2\r\ncluster_stats_messages_sent:1483972\r\ncluster_stats_messages_received:1483968",
			10,
			1,
			nil,
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
			client := testClient{tt.rcv, tt.err}
			response := queryClusterInfo(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Command, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")

				if tt.valuesToCheckInResponse != nil {
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				}
			}
		})
	}
}

func TestQueryClusterNodes(t *testing.T) {
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
			"should parse clusterNodes bulk string",
			queryModel{Command: "clusterNodes"},
			"07c37dfeb235213a872192d90877d0cd55635b91 127.0.0.1:30004@31004 slave e7d1eecce10fd6bb5eb35b9f99a514335d9ba9ca 1609783649927 1426238317239 4 connected\r\n67ed2db8d677e59ec4a4cefb06858cf2a1a89fa1 127.0.0.1:30002@31002 master - 0 1426238316232 2 connected 5461-10922\r\n292f8b365bb7edb5e285caf0b7e6ddc7265d2f4f 127.0.0.1:30003@31003 master - 0 1426238318243 3 connected 10923-16383\r\n6ec23923021cf3ffec47632106199cb7f496ce01 127.0.0.1:30005@31005 slave 67ed2db8d677e59ec4a4cefb06858cf2a1a89fa1 0 1426238316232 5 connected\r\n824fe116063bc5fcf9f4ffd895bc17aee7731ac3 127.0.0.1:30006@31006 slave 292f8b365bb7edb5e285caf0b7e6ddc7265d2f4f 0 1426238317741 6 connected\r\ne7d1eecce10fd6bb5eb35b9f99a514335d9ba9ca 127.0.0.1:30001@31001 myself,master - 0 0 1 connected 0-5460",
			9,
			6,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 2, value: "292f8b365bb7edb5e285caf0b7e6ddc7265d2f4f"},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 3, value: "127.0.0.1:30005@31005"},
				{frameIndex: 0, fieldIndex: 2, rowIndex: 0, value: "slave"},
				{frameIndex: 0, fieldIndex: 4, rowIndex: 0, value: int64(1609783649927)},
				{frameIndex: 0, fieldIndex: 6, rowIndex: 3, value: int64(5)},
				{frameIndex: 0, fieldIndex: 8, rowIndex: 2, value: "10923-16383"},
			},
			nil,
		},
		{
			"should handle string with invalid number of values",
			queryModel{Command: "clusterNodes"},
			"07c37dfeb235213a872192d90877d0cd55635b91 127.0.0.1:30004@31004 e7d1eecce10fd6bb5eb35b9f99a514335d9ba9ca 1609783649927 1426238317239 4 connected",
			9,
			0,
			nil,
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "clusterNodes"},
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
			client := testClient{tt.rcv, tt.err}
			response := queryClusterNodes(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Command, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")

				if tt.valuesToCheckInResponse != nil {
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				}
			}
		})
	}
}

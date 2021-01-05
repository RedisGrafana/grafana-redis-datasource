package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/mediocregopher/radix/v3/resp/resp2"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	t.Parallel()
	tests := []struct {
		qm queryModel
	}{
		{queryModel{Command: "ts.get"}},
		{queryModel{Command: "ts.info"}},
		{queryModel{Command: "ts.queryindex"}},
		{queryModel{Command: "ts.range"}},
		{queryModel{Command: "ts.mrange"}},
		{queryModel{Command: "hgetall"}},
		{queryModel{Command: "smembers"}},
		{queryModel{Command: "hkeys"}},
		{queryModel{Command: "hget"}},
		{queryModel{Command: "hmget"}},
		{queryModel{Command: "info"}},
		{queryModel{Command: "clientList"}},
		{queryModel{Command: "slowlogGet"}},
		{queryModel{Command: "type"}},
		{queryModel{Command: "xinfoStream"}},
		{queryModel{Command: "clusterInfo"}},
		{queryModel{Command: "clusterNodes"}},
		{queryModel{Command: "ft.info"}},
		{queryModel{Command: "xinfoStream"}},
		{queryModel{Query: "DO something"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.qm.Command, func(t *testing.T) {
			t.Parallel()
			ds := redisDatasource{}
			client := TestClient{nil, nil}
			var marshaled, _ = json.Marshal(tt.qm)
			response := ds.query(context.TODO(), backend.DataQuery{
				RefID:         "",
				QueryType:     "",
				MaxDataPoints: 100,
				Interval:      10,
				TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
				JSON:          marshaled,
			}, client)
			require.NoError(t, response.Error, "Should not return error")
		})
	}
}

func TestQueryWithErrors(t *testing.T) {
	t.Parallel()

	t.Run("Marshalling failure", func(t *testing.T) {
		t.Parallel()
		ds := redisDatasource{}
		client := TestClient{nil, nil}
		response := ds.query(context.TODO(), backend.DataQuery{
			RefID:         "",
			QueryType:     "",
			MaxDataPoints: 100,
			Interval:      10,
			TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
			JSON:          []byte{31, 17, 45},
		}, client)

		require.EqualError(t, response.Error, "invalid character '\\x1f' looking for beginning of value", "Should return marshalling error")
	})

	t.Run("Unknown command failure", func(t *testing.T) {
		t.Parallel()
		ds := redisDatasource{}
		client := TestClient{nil, nil}
		var marshaled, _ = json.Marshal(queryModel{Command: "unknown"})
		response := ds.query(context.TODO(), backend.DataQuery{
			RefID:         "",
			QueryType:     "",
			MaxDataPoints: 100,
			Interval:      10,
			TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
			JSON:          marshaled,
		}, client)

		require.EqualError(t, response.Error, "Unknown command", "Should return unknown command error")
	})

}

func TestErrorHandler(t *testing.T) {
	t.Parallel()

	t.Run("Common error", func(t *testing.T) {
		t.Parallel()
		ds := redisDatasource{}
		resp := ds.errorHandler(backend.DataResponse{}, errors.New("common error"))
		require.EqualError(t, resp.Error, "common error", "Should return marshalling error")
	})

	t.Run("Redis error", func(t *testing.T) {
		t.Parallel()
		ds := redisDatasource{}
		resp := ds.errorHandler(backend.DataResponse{}, resp2.Error{E: errors.New("redis error")})
		require.EqualError(t, resp.Error, "redis error", "Should return marshalling error")
	})

}

func TestQueryKeyCommand(t *testing.T) {
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
			"should handle string value",
			queryModel{Command: "get", Key: "test1"},
			"someStr",
			1,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "someStr"},
			},
			nil,
		},
		{
			"should handle float64 value",
			queryModel{Command: "get", Key: "test1"},
			"3.14",
			1,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: 3.14},
			},
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "get"},
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
			response := ds.queryKeyCommand(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Key, response.Frames[0].Name, "Invalid frame name")
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

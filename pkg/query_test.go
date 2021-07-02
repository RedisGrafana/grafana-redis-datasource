package main

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/mediocregopher/radix/v3/resp/resp2"
	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

/**
 * Query
 */
func TestQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		qm queryModel
	}{
		{queryModel{Command: models.TimeSeriesGet}},
		{queryModel{Command: models.TimeSeriesMGet}},
		{queryModel{Command: models.TimeSeriesInfo}},
		{queryModel{Command: models.TimeSeriesQueryIndex}},
		{queryModel{Command: models.TimeSeriesRange}},
		{queryModel{Command: models.TimeSeriesMRange}},
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
		{queryModel{Command: "tmscan"}},
		{queryModel{Command: models.GearsPyStats}},
		{queryModel{Command: models.GearsDumpRegistrations}},
		{queryModel{Command: models.GearsPyExecute}},
		{queryModel{Command: models.GearsPyDumpReqs}},
		{queryModel{Command: "xrange"}},
		{queryModel{Command: "xrevrange"}},
		{queryModel{Command: models.GraphConfig}},
		{queryModel{Command: models.GraphExplain}},
		{queryModel{Command: models.GraphProfile}},
		{queryModel{Command: models.GraphQuery}},
		{queryModel{Command: models.GraphSlowlog}},
	}

	// Run Tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.qm.Command, func(t *testing.T) {
			t.Parallel()

			// Client
			client := testClient{rcv: nil, err: nil}

			// Response
			response := query(context.TODO(), backend.DataQuery{
				RefID:         "",
				QueryType:     "",
				MaxDataPoints: 100,
				Interval:      10,
				TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
			}, &client, tt.qm)
			require.NoError(t, response.Error, "Should not return error")
		})
	}

	// Custom Query
	t.Run("custom query", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{rcv: []interface{}{}, err: nil}
		qm := queryModel{Query: "Test"}

		// Response
		response := query(context.TODO(), backend.DataQuery{
			RefID:         "",
			QueryType:     "",
			MaxDataPoints: 100,
			Interval:      10,
			TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
		}, &client, qm)
		require.NoError(t, response.Error, "Should not return error")
	})
}

/**
 * Query with Error
 */
func TestQueryWithErrors(t *testing.T) {
	t.Parallel()

	// Unknown command
	t.Run("Unknown command failure", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{rcv: nil, err: nil}
		qm := queryModel{Command: "unknown"}

		// Response
		response := query(context.TODO(), backend.DataQuery{
			RefID:         "",
			QueryType:     "",
			MaxDataPoints: 100,
			Interval:      10,
			TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
		}, &client, qm)

		require.NoError(t, response.Error, "Should not return error")
	})

}

/**
 * Error Handler
 */
func TestErrorHandler(t *testing.T) {
	t.Parallel()

	t.Run("Common error", func(t *testing.T) {
		t.Parallel()
		resp := errorHandler(backend.DataResponse{}, errors.New("common error"))
		require.EqualError(t, resp.Error, "common error", "Should return marshalling error")
	})

	t.Run("Redis error", func(t *testing.T) {
		t.Parallel()
		resp := errorHandler(backend.DataResponse{}, resp2.Error{E: errors.New("redis error")})
		require.EqualError(t, resp.Error, "redis error", "Should return marshalling error")
	})

}

/**
 * Query Command with Key
 */
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

	// Run Tests
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Client
			client := testClient{rcv: tt.rcv, err: tt.err}

			// Response
			response := queryKeyCommand(tt.qm, &client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Key, response.Frames[0].Name, "Invalid frame name")
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

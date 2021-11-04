package main

import (
	"errors"
	"testing"

	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

/**
 * ZRANGE
 */
func TestQueryZRange(t *testing.T) {
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
			"should handle default array of strings",
			queryModel{Command: models.ZRange, Key: "test:zset", Min: "0", Max: "-1"},
			[]string{"member1", "10", "member2", "2", "member3", "15"},
			3,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: float64(10)},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: float64(2)},
				{frameIndex: 0, fieldIndex: 2, rowIndex: 0, value: float64(15)},
			},
			nil,
		},
		{
			"should handle default array of strings",
			queryModel{Command: models.ZRange, Key: "test:zset", ZRangeQuery: "BYSCORE", Min: "-inf", Max: "+inf"},
			[]string{"member1", "test", "member2", "2", "member3", "15"},
			3,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "test"},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: float64(2)},
				{frameIndex: 0, fieldIndex: 2, rowIndex: 0, value: float64(15)},
			},
			nil,
		},
		{
			"should handle error",
			queryModel{Command: models.ZRange},
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

			client := testClient{rcv: tt.rcv, err: tt.err}
			response := queryZRange(tt.qm, &client)
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

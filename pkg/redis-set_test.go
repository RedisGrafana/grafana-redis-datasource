package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuerySMembers(t *testing.T) {
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
			queryModel{Command: "smembers", Key: "test1"},
			[]string{"value1", "2", "3.14"},
			1,
			3,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "value1"},
				{frameIndex: 0, fieldIndex: 0, rowIndex: 1, value: "2"},
				{frameIndex: 0, fieldIndex: 0, rowIndex: 2, value: "3.14"},
			},
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "smembers"},
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
			client := testClient{tt.rcv, tt.err}
			response := ds.querySMembers(tt.qm, client)
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

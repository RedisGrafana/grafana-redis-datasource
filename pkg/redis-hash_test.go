package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryHGetAll(t *testing.T) {
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
			queryModel{Command: "hgetall", Key: "test1"},
			[]string{"key1", "value1", "key2", "2", "key3", "3.14"},
			3,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "value1"},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: float64(2)},
				{frameIndex: 0, fieldIndex: 2, rowIndex: 0, value: 3.14},
			},
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "hgetall"},
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
			response := ds.queryHGetAll(tt.qm, client)
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

func TestQueryHGet(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		qm           queryModel
		rcv          interface{}
		fieldsCount  int
		rowsPerField int
		value        interface{}
		err          error
	}{
		{
			"should handle simple string",
			queryModel{Command: "hget", Key: "test1", Field: "field1"},
			"value1",
			1,
			1,
			"value1",
			nil,
		},
		{
			"should handle string with underlying float64 value",
			queryModel{Command: "hget", Key: "test1", Field: "key1"},
			"3.14",
			1,
			1,
			3.14,
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "hget"},
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
			response := ds.queryHGet(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Field, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")
				require.Equal(t, tt.value, response.Frames[0].Fields[0].At(0), "Invalid value contained in frame")

			}
		})
	}
}

func TestQueryHMGet(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		qm                      queryModel
		rcv                     interface{}
		fieldsCount             int
		rowsPerField            int
		shouldCreateFrames      bool
		valuesToCheckInResponse []valueToCheckInResponse
		err                     error
	}{
		{
			"should handle 3 fields with different underlying types",
			queryModel{Command: "hmget", Key: "test1", Field: "field1 field2 field3"},
			[]string{"value1", "2", "3.14"},
			3,
			1,
			true,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "value1"},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: float64(2)},
				{frameIndex: 0, fieldIndex: 2, rowIndex: 0, value: 3.14},
			},
			nil,
		},
		{
			"should handle Field string parsing error and create no fields",
			queryModel{Command: "hmget", Key: "test1", Field: "field1 field2\"field3"},
			nil,
			0,
			0,
			false,
			nil,
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "hmget"},
			nil,
			0,
			0,
			false,
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
			response := ds.queryHMGet(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				if tt.shouldCreateFrames {
					require.Equal(t, tt.qm.Command, response.Frames[0].Name, "Invalid frame name")
					require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
					require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				} else {
					require.Nil(t, response.Frames, "Should not create frames in response")
				}
			}
		})
	}
}

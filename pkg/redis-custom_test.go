package main

import (
	"errors"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExecuteCustomQuery(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		qm   queryModel
		rcv  interface{}
		err  error
	}{
		{
			"should parse correct real-world command with params",
			queryModel{Query: "config get *max-*-entries*"},
			[]interface{}{
				[]byte("hash-max-ziplist-entries"),
				[]byte("512"),
				[]byte("set-max-intset-entries"),
				[]byte("512"),
				[]byte("zset-max-ziplist-entries"),
				[]byte("128"),
			},
			nil,
		},
		{
			"should parse correct real-world command without params",
			queryModel{Query: "lastsave"},
			int64(1609840612),
			nil,
		},
		{
			"should handle error if invalid command string",
			queryModel{Query: "lastsave \""},
			nil,
			errors.New("Query is not valid"),
		},
		{
			"should handle error",
			queryModel{Query: "lastsave"},
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
			result, err := ds.executeCustomQuery(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, result, "No result should be created if failed")
			} else {
				require.Equal(t, tt.rcv, result, "Should return receiver value")
			}
		})
	}
}

func TestExecuteCustomQueryWithPanic(t *testing.T) {
	t.Parallel()
	ds := redisDatasource{}
	client := panickingClient{}
	result, err := ds.executeCustomQuery(queryModel{Query: "panic"}, client)
	require.NoError(t, err, "Should return error")
	require.Nil(t, result, "No result if panicked")
}

func TestParseInterfaceValue(t *testing.T) {
	t.Parallel()
	t.Run("should parse complex input", func(t *testing.T) {
		t.Parallel()
		ds := redisDatasource{}
		inputResponse := backend.DataResponse{}
		input := []interface{}{
			"str",
			[]byte("str2"),
			int64(42),
			[]interface{}{},
			[]interface{}{
				"str",
				[]byte("str3"),
				int64(66),
			},
		}

		expected := []string{"str", "str2", "42", "(empty array)", "str", "str3", "66"}
		result, response := ds.parseInterfaceValue(input, inputResponse)
		require.NoError(t, response.Error, "Should return error")
		require.Equal(t, expected, result, "Invalid function return value")

	})
	t.Run("should fail on unsupported type", func(t *testing.T) {
		t.Parallel()
		ds := redisDatasource{}
		inputResponse := backend.DataResponse{}
		input := []interface{}{
			"str",
			[]byte("str2"),
			int64(42),
			3.14,
			[]interface{}{},
			[]interface{}{
				"str",
				[]byte("str2"),
				int64(42),
				3.14,
			},
		}

		expected := []string{"str", "str2", "42"}
		result, response := ds.parseInterfaceValue(input, inputResponse)
		require.EqualError(t, response.Error, "Unsupported array return type", "Should return error")
		require.Equal(t, expected, result, "Should contain results before unsupported parameter")
	})

}

func TestQueryCustomCommand(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		qm                      queryModel
		rcv                     interface{}
		fieldsCount             int
		rowsPerField            int
		valuesToCheckInResponse []valueToCheckInResponse
		err                     error
		errToCheck              string
	}{
		{
			"should handle empty interface array without values",
			queryModel{Query: "test"},
			[]interface{}{},
			1,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "(empty array)"},
			},
			nil,
			"",
		},
		{
			"should handle empty interface array with nesting",
			queryModel{Query: "test"},
			[]interface{}{
				"str",
				[]byte("str2"),
				int64(42),
				[]interface{}{},
				[]interface{}{
					"str",
					[]byte("str3"),
					int64(66),
				},
			},
			1,
			7,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 6, value: "66"},
			},
			nil,
			"",
		},
		{
			"should handle string",
			queryModel{Query: "test"},
			"str",
			1,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "str"},
			},
			nil,
			"",
		},
		{
			"should handle []byte with single string inside",
			queryModel{Query: "test"},
			[]byte("str"),
			1,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "str"},
			},
			nil,
			"",
		},
		{
			"should handle []byte with bulk string inside",
			queryModel{Query: "test"},
			[]byte("str\r\nstr2"),
			1,
			2,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "str"},
				{frameIndex: 0, fieldIndex: 0, rowIndex: 1, value: "str2"},
			},
			nil,
			"",
		},
		{
			"should handle int64 and return field as int64",
			queryModel{Query: "test"},
			int64(42),
			1,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: int64(42)},
			},
			nil,
			"",
		},
		{
			"should fail with emtpy command",
			queryModel{Query: ""},
			nil,
			0,
			0,
			nil,
			nil,
			"Command is empty",
		},
		{
			"should handle error",
			queryModel{Query: "test"},
			nil,
			0,
			0,
			nil,
			errors.New("error occurred"),
			"error occurred",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ds := redisDatasource{}
			client := testClient{tt.rcv, tt.err}
			response := ds.queryCustomCommand(tt.qm, client)
			if tt.errToCheck != "" {
				require.EqualError(t, response.Error, tt.errToCheck, "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Command, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")
				for _, value := range tt.valuesToCheckInResponse {
					require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
				}
			}
		})
	}
}

package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryXInfoStream(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		qm           queryModel
		rcv          interface{}
		fieldsCount  int
		rowsPerField int
		err          error
	}{
		{
			"should handle default payload, but collect only top-level key-value pairs",
			queryModel{Command: "xinfoStream", Key: "test1"},
			map[string]string{
				"length":            "2",
				"radix-tree-keys":   "1",
				"radix-tree-nodes":  "2",
				"groups":            "2",
				"last-generated-id": "1538385846314-0",
			},
			2,
			5,
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "xinfoStream"},
			nil,
			0,
			0,
			errors.New("error occurred"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			client := testClient{tt.rcv, tt.err}
			response := queryXInfoStream(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Key, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")

			}
		})
	}
}

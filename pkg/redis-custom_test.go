package main

import (
	"errors"
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

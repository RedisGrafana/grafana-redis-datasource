// +build integration

package main

import (
	"fmt"
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/require"
)

func TestTMScanIntegration(t *testing.T) {
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	types := map[string]string{
		"test:string": "string",
		"test:stream": "stream",
		"test:set":    "set",
		"test:list":   "list",
		"test:float":  "string",
		"test:hash":   "hash",
	}

	memory := map[string]int64{
		"test:string": int64(59),
		"test:stream": int64(612),
		"test:set":    int64(265),
		"test:list":   int64(140),
		"test:float":  int64(59),
		"test:hash":   int64(108),
	}

	client := radixV3Impl{radixClient: radixClient}
	resp := queryTMScan(queryModel{Cursor: "0"}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 1)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.Equal(t, 6, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 6, resp.Frames[0].Fields[1].Len())
	require.Equal(t, 6, resp.Frames[0].Fields[2].Len())
	require.Equal(t, "0", resp.Frames[1].Fields[0].At(0))

	keys := []string{
		resp.Frames[0].Fields[0].At(0).(string),
		resp.Frames[0].Fields[0].At(1).(string),
		resp.Frames[0].Fields[0].At(2).(string),
		resp.Frames[0].Fields[0].At(3).(string),
		resp.Frames[0].Fields[0].At(4).(string),
		resp.Frames[0].Fields[0].At(5).(string),
	}

	for i, key := range keys {
		require.Equal(t, types[key], resp.Frames[0].Fields[1].At(i), "Invalid type returned")
		require.Equal(t, memory[key], resp.Frames[0].Fields[2].At(i), "Invalid memory size returned")
	}
}

func TestTMScanIntegrationWithNoMatched(t *testing.T) {
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)

	client := radixV3Impl{radixClient: radixClient}
	resp := queryTMScan(queryModel{Cursor: "0", Match: "nomatch"}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 1)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.Equal(t, 0, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 0, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 0, resp.Frames[0].Fields[0].Len())
	require.Equal(t, "0", resp.Frames[1].Fields[0].At(0))
}

func TestTMScanIntegrationWithMatched(t *testing.T) {
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)

	client := radixV3Impl{radixClient: radixClient}
	resp := queryTMScan(queryModel{Cursor: "0", Match: "test:*"}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 1)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.Equal(t, 6, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 6, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 6, resp.Frames[0].Fields[0].Len())
	require.Equal(t, "0", resp.Frames[1].Fields[0].At(0))
}

func TestTMScanIntegrationWithCount(t *testing.T) {
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)

	client := radixV3Impl{radixClient: radixClient}
	resp := queryTMScan(queryModel{Cursor: "0", Count: 1}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 1)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
	require.NotEqual(t, "0", resp.Frames[1].Fields[0].At(0))
}

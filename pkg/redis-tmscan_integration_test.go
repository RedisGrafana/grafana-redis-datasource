// +build integration

package main

import (
	"fmt"
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/require"
)

// Types
var types = map[string]string{
	"test:string":                "string",
	"test:stream":                "stream",
	"test:set":                   "set",
	"test:list":                  "list",
	"test:float":                 "string",
	"test:hash":                  "hash",
	"ts:enqueue:queue:customers": "TSDB-TYPE",
	"queue:complete":             "stream",
	"queue:orders":               "stream",
	"ts:len:queue:complete":      "TSDB-TYPE",
	"queue:customers":            "stream",
	"product":                    "string",
	"ts:enqueue:queue:complete":  "TSDB-TYPE",
	"ts:enqueue:queue:orders":    "TSDB-TYPE",
	"ts:len:queue:customers":     "TSDB-TYPE",
	"ts:len:queue:orders":        "TSDB-TYPE",
}

// Memory
var memory = map[string]int64{
	"test:string":                int64(59),
	"test:stream":                int64(612),
	"test:set":                   int64(265),
	"test:list":                  int64(140),
	"test:float":                 int64(59),
	"test:hash":                  int64(108),
	"ts:enqueue:queue:customers": int64(4236),
	"queue:complete":             int64(12876),
	"queue:orders":               int64(1219),
	"ts:len:queue:complete":      int64(4231),
	"queue:customers":            int64(2035),
	"product":                    int64(49),
	"ts:enqueue:queue:complete":  int64(4235),
	"ts:enqueue:queue:orders":    int64(4233),
	"ts:len:queue:customers":     int64(4232),
	"ts:len:queue:orders":        int64(4229),
}

/**
 * TMSCAN
 */

func TestTMScanIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryTMScan(queryModel{Cursor: "0", Count: 20}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 1)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())

	// Keys
	keys := map[string]int{}
	for i := 0; i < resp.Frames[0].Fields[0].Len(); i++ {
		if _, ok := types[resp.Frames[0].Fields[0].At(i).(string)]; ok {
			keys[resp.Frames[0].Fields[0].At(i).(string)] = i
		}
	}
	for key, value := range keys {
		require.Equal(t, types[key], resp.Frames[0].Fields[1].At(value), "Invalid type returned")
		require.Equal(t, memory[key], resp.Frames[0].Fields[2].At(value), "Invalid memory size returned")
	}
}

/**
 * TMSCAN with Match nomatch
 */
func TestTMScanIntegrationWithNoMatched(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryTMScan(queryModel{Cursor: "0", Match: "nomatch"}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 1)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.Equal(t, 0, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 0, resp.Frames[0].Fields[1].Len())
	require.Equal(t, 0, resp.Frames[0].Fields[2].Len())
	require.NotEqual(t, "0", resp.Frames[1].Fields[0].At(0))
}

/**
 * TMSCAN with Match test:*
 */
func TestTMScanIntegrationWithMatched(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryTMScan(queryModel{Cursor: "0", Match: "test:*", Count: 20}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 1)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.Equal(t, 6, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 6, resp.Frames[0].Fields[1].Len())
	require.Equal(t, 6, resp.Frames[0].Fields[2].Len())
	require.Equal(t, "0", resp.Frames[1].Fields[0].At(0))
}

func TestTMScanIntegrationWithSamples(t *testing.T) {
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)

	client := radixV3Impl{radixClient: radixClient}
	resp := queryTMScan(queryModel{Cursor: "0", Samples: 10}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 1)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.GreaterOrEqual(t, resp.Frames[0].Fields[0].Len(), 10)
	require.GreaterOrEqual(t, resp.Frames[0].Fields[1].Len(), 10)
	require.GreaterOrEqual(t, resp.Frames[0].Fields[2].Len(), 10)
	require.NotEqual(t, "0", resp.Frames[1].Fields[0].At(0))
}

/**
 * TMSCAN with Size 10
 */
func TestTMScanIntegrationWithSize(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryTMScan(queryModel{Cursor: "0", Count: 20, Size: 10}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 1)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.Equal(t, 10, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 10, resp.Frames[0].Fields[1].Len())
	require.Equal(t, 10, resp.Frames[0].Fields[2].Len())
	require.Equal(t, memory["queue:complete"], resp.Frames[0].Fields[2].At(0))
	require.Equal(t, "queue:complete", resp.Frames[0].Fields[0].At(0))
	require.Equal(t, memory["test:stream"], resp.Frames[0].Fields[2].At(9))
	require.Equal(t, "test:stream", resp.Frames[0].Fields[0].At(9))
	require.Equal(t, "0", resp.Frames[1].Fields[0].At(0))
}

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
	"test:string": "string",
	"test:stream": "stream",
	"test:set":    "set",
	"test:list":   "list",
	"test:float":  "string",
	"test:hash":   "hash",
}

// Memory
var memory = map[string]int64{
	"test:string": int64(59),
	"test:stream": int64(612),
	"test:set":    int64(265),
	"test:list":   int64(140),
	"test:float":  int64(59),
	"test:hash":   int64(108),
}

/**
 * TMSCAN
 */
func TestTMScanIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("%s:%d", integrationTestIP, integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryTMScan(queryModel{Cursor: "0", Count: 5}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 2)
	require.Equal(t, "cursor", resp.Frames[1].Fields[0].Name)
	require.Equal(t, "count", resp.Frames[1].Fields[1].Name)
	require.GreaterOrEqual(t, resp.Frames[0].Fields[0].Len(), 5)
	require.GreaterOrEqual(t, resp.Frames[0].Fields[1].Len(), 5)
	require.GreaterOrEqual(t, resp.Frames[0].Fields[2].Len(), 5)
	require.IsType(t, "", resp.Frames[0].Fields[0].At(0))
	require.IsType(t, "", resp.Frames[0].Fields[1].At(0))
	require.IsType(t, int64(0), resp.Frames[0].Fields[2].At(0))
	require.NotEqual(t, "0", resp.Frames[1].Fields[0].At(0))
	require.Equal(t, int64(resp.Frames[0].Fields[0].Len()), resp.Frames[1].Fields[1].At(0))
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())

}

/**
 * TMSCAN with Match nomatch
 */
func TestTMScanIntegrationWithNoMatched(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("%s:%d", integrationTestIP, integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryTMScan(queryModel{Cursor: "0", Match: "nomatch"}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 2)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.Equal(t, 0, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 0, resp.Frames[0].Fields[1].Len())
	require.Equal(t, 0, resp.Frames[0].Fields[2].Len())
	require.NotEqual(t, "0", resp.Frames[1].Fields[0].At(0))
	require.Equal(t, int64(resp.Frames[0].Fields[0].Len()), resp.Frames[1].Fields[1].At(0))
}

/**
 * TMSCAN with Match test:*
 */
func TestTMScanIntegrationWithMatched(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("%s:%d", integrationTestIP, integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryTMScan(queryModel{Cursor: "0", Match: "test:*", Count: 20}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 2)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.GreaterOrEqual(t, 9, resp.Frames[0].Fields[0].Len())
	require.GreaterOrEqual(t, 9, resp.Frames[0].Fields[1].Len())
	require.GreaterOrEqual(t, 9, resp.Frames[0].Fields[2].Len())
	require.Equal(t, int64(resp.Frames[0].Fields[0].Len()), resp.Frames[1].Fields[1].At(0))

	// Keys
	keys := map[string]int{}
	for i := 0; i < resp.Frames[0].Fields[0].Len(); i++ {
		if _, ok := types[resp.Frames[0].Fields[0].At(i).(string)]; ok {
			keys[resp.Frames[0].Fields[0].At(i).(string)] = i
		}
	}

	for key, value := range keys {
		require.Equal(t, types[key], resp.Frames[0].Fields[1].At(value), "Invalid type returned")
		require.LessOrEqual(t, memory[key], resp.Frames[0].Fields[2].At(value), "Invalid memory size returned")
	}
}

/**
 * TMSCAN with Samples count
 */
func TestTMScanIntegrationWithSamples(t *testing.T) {
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("%s:%d", integrationTestIP, integrationTestPort), 10)

	// Client
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryTMScan(queryModel{Cursor: "0", Samples: 10}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 2)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.GreaterOrEqual(t, resp.Frames[0].Fields[0].Len(), 10)
	require.GreaterOrEqual(t, resp.Frames[0].Fields[1].Len(), 10)
	require.GreaterOrEqual(t, resp.Frames[0].Fields[2].Len(), 10)
	require.NotEqual(t, "0", resp.Frames[1].Fields[0].At(0))
	require.Equal(t, int64(resp.Frames[0].Fields[0].Len()), resp.Frames[1].Fields[1].At(0))
}

/**
 * TMSCAN with Size 10
 */
func TestTMScanIntegrationWithSize(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("%s:%d", integrationTestIP, integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryTMScan(queryModel{Cursor: "0", Count: 10, Size: 8}, &client)
	require.Len(t, resp.Frames, 2)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.Len(t, resp.Frames[1].Fields, 2)
	require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
	require.Equal(t, 8, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 8, resp.Frames[0].Fields[1].Len())
	require.Equal(t, 8, resp.Frames[0].Fields[2].Len())

	// Check proper sorting by memory
	for i := 0; i < 7; i++ {
		require.LessOrEqual(t, resp.Frames[0].Fields[2].At(i+1), resp.Frames[0].Fields[2].At(i))
	}

	require.NotEqual(t, "0", resp.Frames[1].Fields[0].At(0))
	require.GreaterOrEqual(t, int64(11), resp.Frames[1].Fields[1].At(0))
}

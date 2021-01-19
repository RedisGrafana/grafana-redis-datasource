// +build integration

package main

import (
	"fmt"
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/require"
)

/**
 * RG.PYSTATS
 */
func TestRgPystatsIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryRgPystats(queryModel{Command: "rg.pystats"}, &client)
	require.Len(t, resp.Frames, 1)
	require.Len(t, resp.Frames[0].Fields, 3)
	require.IsType(t, int64(0), resp.Frames[0].Fields[0].At(0))
	require.IsType(t, int64(0), resp.Frames[0].Fields[1].At(0))
	require.IsType(t, int64(0), resp.Frames[0].Fields[2].At(0))
	require.Equal(t, "TotalAllocated", resp.Frames[0].Fields[0].Name)
	require.Equal(t, "PeakAllocated", resp.Frames[0].Fields[1].Name)
	require.Equal(t, "CurrAllocated", resp.Frames[0].Fields[2].Name)
}

/**
 * RG.DUMPREGISTRATIONS
 */
func TestRgDumpregistrationsIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryRgDumpregistrations(queryModel{Command: "rg.dumpregistrations"}, &client)
	require.Len(t, resp.Frames[0].Fields, 12)
	require.Equal(t, "id", resp.Frames[0].Fields[0].Name)
	require.Equal(t, "reader", resp.Frames[0].Fields[1].Name)
	require.Equal(t, "desc", resp.Frames[0].Fields[2].Name)
	require.Equal(t, "PD", resp.Frames[0].Fields[3].Name)
	require.Equal(t, "mode", resp.Frames[0].Fields[4].Name)
	require.Equal(t, "numTriggered", resp.Frames[0].Fields[5].Name)
	require.Equal(t, "numSuccess", resp.Frames[0].Fields[6].Name)
	require.Equal(t, "numFailures", resp.Frames[0].Fields[7].Name)
	require.Equal(t, "numAborted", resp.Frames[0].Fields[8].Name)
	require.Equal(t, "lastError", resp.Frames[0].Fields[9].Name)
	require.Equal(t, "args", resp.Frames[0].Fields[10].Name)
	require.Equal(t, 3, resp.Frames[0].Fields[0].Len())
	require.Equal(t, 3, resp.Frames[0].Fields[1].Len())
	require.Equal(t, 3, resp.Frames[0].Fields[2].Len())
	require.Equal(t, 3, resp.Frames[0].Fields[3].Len())
	require.Equal(t, 3, resp.Frames[0].Fields[4].Len())
	require.Equal(t, 3, resp.Frames[0].Fields[5].Len())
	require.Equal(t, 3, resp.Frames[0].Fields[6].Len())
	require.Equal(t, 3, resp.Frames[0].Fields[7].Len())
	require.Equal(t, 3, resp.Frames[0].Fields[8].Len())
	require.Equal(t, 3, resp.Frames[0].Fields[9].Len())
	require.Equal(t, 3, resp.Frames[0].Fields[10].Len())
}

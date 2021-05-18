// +build integration

package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/mediocregopher/radix/v3"
	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
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
	resp := queryRgPystats(queryModel{Command: models.GearsPyStats}, &client)
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
	resp := queryRgDumpregistrations(queryModel{Command: models.GearsDumpRegistrations}, &client)
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
	for i := 0; i < len(resp.Frames[0].Fields); i++ {
		require.Equal(t, 3, resp.Frames[0].Fields[0].Len())
	}
}

/**
 * RG.PYEXECUTE
 */
func TestRgPyexecuteIntegration(t *testing.T) {
	// Increase timeout to 30 seconds for requirements
	customConnFunc := func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialTimeout(30*time.Second),
		)
	}

	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10, radix.PoolConnFunc(customConnFunc))
	client := radixV3Impl{radixClient: radixClient}

	// Results
	t.Run("Test command with full response", func(t *testing.T) {
		resp := queryRgPyexecute(queryModel{Command: models.GearsPyExecute, Key: "GB().run()"}, &client)
		require.Len(t, resp.Frames, 2)
		require.Len(t, resp.Frames[0].Fields, 1)
		require.Equal(t, "results", resp.Frames[0].Name)
		require.Equal(t, "results", resp.Frames[0].Fields[0].Name)
		require.Greater(t, resp.Frames[0].Fields[0].Len(), 0)
		require.IsType(t, "", resp.Frames[0].Fields[0].At(0))
		require.Len(t, resp.Frames[1].Fields, 1)
		require.Equal(t, "errors", resp.Frames[1].Name)
		require.Equal(t, "errors", resp.Frames[1].Fields[0].Name)
		require.NoError(t, resp.Error)
	})

	// UNBLOCKING and REQUIREMENTS
	t.Run("Test command with UNBLOCKING and REQUIREMENTS", func(t *testing.T) {
		resp := queryRgPyexecute(queryModel{Command: models.GearsPyExecute, Key: "GearsBuilder(reader=\"KeysReader\").run()", Unblocking: true, Requirements: "numpy"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 1)
		require.Equal(t, "operationId", resp.Frames[0].Name)
		require.Equal(t, "operationId", resp.Frames[0].Fields[0].Name)
		require.Greater(t, resp.Frames[0].Fields[0].Len(), 0)
		require.IsType(t, "", resp.Frames[0].Fields[0].At(0))
	})

	// OK
	t.Run("Test command with full OK string", func(t *testing.T) {
		resp := queryRgPyexecute(queryModel{Command: models.GearsPyExecute, Key: "GB('CommandReader')"}, &client)
		require.Len(t, resp.Frames, 2)
		require.Len(t, resp.Frames[0].Fields, 1)
		require.Equal(t, "results", resp.Frames[0].Name)
		require.Equal(t, "results", resp.Frames[0].Fields[0].Name)
		require.Equal(t, 0, resp.Frames[0].Fields[0].Len())
		require.Len(t, resp.Frames[1].Fields, 1)
		require.Equal(t, "errors", resp.Frames[1].Name)
		require.Equal(t, "errors", resp.Frames[1].Fields[0].Name)
		require.Equal(t, 0, resp.Frames[1].Fields[0].Len())
		require.NoError(t, resp.Error)
	})

	// Error
	t.Run("Test command with error", func(t *testing.T) {
		resp := queryRgPyexecute(queryModel{Command: models.GearsPyExecute, Key: "some key"}, &client)
		require.Len(t, resp.Frames, 0)
		require.Error(t, resp.Error)
	})
}

/**
 * RG.DUMPREQS
 */
func TestRgDumpReqsIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryRgPydumpReqs(queryModel{Command: models.GearsPyDumpReqs}, &client)

	require.Len(t, resp.Frames[0].Fields, 6)
	require.Equal(t, "GearReqVersion", resp.Frames[0].Fields[0].Name)
	require.Equal(t, "Name", resp.Frames[0].Fields[1].Name)
	require.Equal(t, "IsDownloaded", resp.Frames[0].Fields[2].Name)
	require.Equal(t, "IsInstalled", resp.Frames[0].Fields[3].Name)
	require.Equal(t, "CompiledOs", resp.Frames[0].Fields[4].Name)
	require.Equal(t, "Wheels", resp.Frames[0].Fields[5].Name)
	for i := 0; i < len(resp.Frames[0].Fields); i++ {
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
	}
}

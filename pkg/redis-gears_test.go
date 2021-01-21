package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

/**
 * RG.PYSTATS
 */
func TestRgPystats(t *testing.T) {
	t.Parallel()

	/**
	 * Success
	 */
	t.Run("should process command", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: pystats{
				TotalAllocated: int64(11),
				PeakAllocated:  int64(12),
				CurrAllocated:  int64(13),
			},
			err: nil,
		}

		// Response
		resp := queryRgPystats(queryModel{Command: "rg.pystats"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 3)
		require.Equal(t, int64(11), resp.Frames[0].Fields[0].At(0))
		require.Equal(t, int64(12), resp.Frames[0].Fields[1].At(0))
		require.Equal(t, int64(13), resp.Frames[0].Fields[2].At(0))
		require.Equal(t, "TotalAllocated", resp.Frames[0].Fields[0].Name)
		require.Equal(t, "PeakAllocated", resp.Frames[0].Fields[1].Name)
		require.Equal(t, "CurrAllocated", resp.Frames[0].Fields[2].Name)

	})

	/**
	 * Error
	 */
	t.Run("should handle error", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv:      nil,
			batchRcv: nil,
			err:      errors.New("error occurred")}

		// Response
		resp := queryRgPystats(queryModel{Command: "rg.pystats"}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})
}

/**
 * RG.DUMPREGISTRATIONS
 */
func TestRgDumpregistrations(t *testing.T) {
	t.Parallel()

	/**
	 * Success
	 */
	t.Run("should process command", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: []dumpregistrations{{
				ID:     "123",
				Reader: "reader",
				Desc:   "desc",
				RegistrationData: registrationData{
					Mode:         "async",
					NumTriggered: 1,
					NumSuccess:   2,
					NumFailures:  3,
					NumAborted:   4,
					LastError:    "some err",
					Args:         map[string]string{"mytrigger": "trigger"},
				},
				PD: "some_pd",
			}},
			err: nil,
		}

		// Response
		resp := queryRgDumpregistrations(queryModel{Command: "rg.dumpregistrations"}, &client)
		require.Len(t, resp.Frames, 1)
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
		require.Equal(t, "123", resp.Frames[0].Fields[0].At(0))
		require.Equal(t, "reader", resp.Frames[0].Fields[1].At(0))
		require.Equal(t, "desc", resp.Frames[0].Fields[2].At(0))
		require.Equal(t, "some_pd", resp.Frames[0].Fields[3].At(0))
		require.Equal(t, "async", resp.Frames[0].Fields[4].At(0))
		require.Equal(t, int64(1), resp.Frames[0].Fields[5].At(0))
		require.Equal(t, int64(2), resp.Frames[0].Fields[6].At(0))
		require.Equal(t, int64(3), resp.Frames[0].Fields[7].At(0))
		require.Equal(t, int64(4), resp.Frames[0].Fields[8].At(0))
		require.Equal(t, "some err", resp.Frames[0].Fields[9].At(0))
		require.Equal(t, "\"mytrigger\"=\"trigger\"\n", resp.Frames[0].Fields[10].At(0))
	})

	/**
	 * Error
	 */
	t.Run("should handle error", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv:      nil,
			batchRcv: nil,
			err:      errors.New("error occurred")}

		// Response
		resp := queryRgDumpregistrations(queryModel{Command: "rg.dumpregistrations"}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})
}

/**
 * RG.PYEXECUTE
 */
func TestRgPyexecute(t *testing.T) {
	t.Parallel()

	/**
	 * Success with OK
	 */
	t.Run("should process command with OK result", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: "OK",
			err: nil,
		}

		// Response
		resp := queryRgPyexecute(queryModel{Command: "rg.pyexecute"}, &client)
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

	/**
	 * Success with 2 arrays in result
	 */
	t.Run("should process command with 2 arrays in result", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: []interface{}{
				[]interface{}{
					[]byte("success info"),
				},
				[]interface{}{
					[]byte("error info"),
				},
			},
			err: nil,
		}

		// Response
		resp := queryRgPyexecute(queryModel{Command: "rg.pyexecute"}, &client)
		require.Len(t, resp.Frames, 2)
		require.Len(t, resp.Frames[0].Fields, 1)
		require.Equal(t, "results", resp.Frames[0].Name)
		require.Equal(t, "results", resp.Frames[0].Fields[0].Name)
		require.IsType(t, "", resp.Frames[0].Fields[0].At(0))
		require.Len(t, resp.Frames[1].Fields, 1)
		require.Equal(t, "errors", resp.Frames[1].Name)
		require.Equal(t, "errors", resp.Frames[1].Fields[0].Name)
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
		require.Equal(t, "success info", resp.Frames[0].Fields[0].At(0))
		require.Equal(t, 1, resp.Frames[1].Fields[0].Len())
		require.Equal(t, "error info", resp.Frames[1].Fields[0].At(0))
		require.NoError(t, resp.Error)
	})

	/**
	 * Error
	 */
	t.Run("should handle error", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv:      nil,
			batchRcv: nil,
			err:      errors.New("error occurred")}

		// Response
		resp := queryRgPyexecute(queryModel{Command: "rg.pyexecute"}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})
}

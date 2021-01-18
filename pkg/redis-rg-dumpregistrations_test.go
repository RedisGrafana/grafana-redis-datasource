package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRgDumpregistrations(t *testing.T) {
	t.Parallel()

	t.Run("should process command", func(t *testing.T) {
		t.Parallel()
		client := testClient{
			rcv: []dumpregistrations{{
				Id:     "123",
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
		resp := queryRgDumpregistrations(queryModel{Command: "rg.dumpregistrations"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 11)
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

	t.Run("should handle error", func(t *testing.T) {
		t.Parallel()
		client := testClient{
			rcv:      nil,
			batchRcv: nil,
			err:      errors.New("error occurred")}
		resp := queryRgDumpregistrations(queryModel{Command: "rg.dumpregistrations"}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})
}

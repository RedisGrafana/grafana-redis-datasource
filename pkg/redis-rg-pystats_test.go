package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRgPystats(t *testing.T) {
	t.Parallel()

	t.Run("should process command", func(t *testing.T) {
		t.Parallel()
		client := testClient{
			rcv: pystats{
				TotalAllocated: int64(11),
				PeakAllocated:  int64(12),
				CurrAllocated:  int64(13),
			},
			err: nil,
		}
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

	t.Run("should handle error", func(t *testing.T) {
		t.Parallel()
		client := testClient{
			rcv:      nil,
			batchRcv: nil,
			err:      errors.New("error occurred")}
		resp := queryRgPystats(queryModel{Command: "rg.pystats"}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})

}

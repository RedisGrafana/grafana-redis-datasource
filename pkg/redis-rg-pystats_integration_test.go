// +build integration

package main

import (
	"fmt"
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/require"
)

func TestRgPystatsIntegration(t *testing.T) {
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)

	client := radixV3Impl{radixClient: radixClient}
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

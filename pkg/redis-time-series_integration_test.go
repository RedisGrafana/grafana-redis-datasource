// +build integration

package main

import (
	"fmt"
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/require"
)

/**
 * TS.INFO
 */
func TestTSInfoIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryTsInfo(queryModel{Command: "ts.info", Key: "test:timeseries2"}, &client)
	require.Len(t, resp.Frames, 1)
	require.Len(t, resp.Frames[0].Fields, 12)
}

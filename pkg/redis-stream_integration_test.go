// +build integration

package main

import (
	"fmt"
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/require"
)

/**
 * TMSCAN
 */

func TestXInfoStreamIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	t.Run("query stream queue:customers", func(t *testing.T) {
		resp := queryXInfoStream(queryModel{Key: "queue:customers"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 9)
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
	})

	t.Run("query stream queue:orders", func(t *testing.T) {
		resp := queryXInfoStream(queryModel{Key: "queue:orders"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 9)
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
	})
}

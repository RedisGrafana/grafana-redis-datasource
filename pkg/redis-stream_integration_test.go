// +build integration

package main

import (
	"fmt"
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/require"
)

/**
 * XINFO
 */
func TestXInfoStreamIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Customers
	t.Run("query stream queue:customers", func(t *testing.T) {
		resp := queryXInfoStream(queryModel{Key: "queue:customers"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 9)
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
	})

	// Orders
	t.Run("query stream queue:orders", func(t *testing.T) {
		resp := queryXInfoStream(queryModel{Key: "queue:orders"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 9)
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
	})
}

/**
 * XRANGE
 */

func TestXRangeStreamIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	t.Run("query stream queue:customers", func(t *testing.T) {
		resp := queryXRange(queryModel{Key: "queue:customers"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 2)
		require.Equal(t, "$streamId", resp.Frames[0].Fields[0].Name)
		require.Equal(t, "id", resp.Frames[0].Fields[1].Name)
		require.Equal(t, resp.Frames[0].Fields[1].Len(), resp.Frames[0].Fields[0].Len())
	})

	t.Run("query stream queue:customers with COUNT", func(t *testing.T) {
		resp := queryXRange(queryModel{Key: "queue:customers", Count: 3}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 2)
		require.Equal(t, "$streamId", resp.Frames[0].Fields[0].Name)
		require.Equal(t, "id", resp.Frames[0].Fields[1].Name)
		require.Equal(t, 3, resp.Frames[0].Fields[0].Len())
		require.Equal(t, 3, resp.Frames[0].Fields[1].Len())
	})

	t.Run("query stream queue:customers with start and end", func(t *testing.T) {
		resp := queryXRange(queryModel{Key: "queue:customers", Start: "1611019111439-0", End: "1611019111985-0"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 2)
		require.Equal(t, "$streamId", resp.Frames[0].Fields[0].Name)
		require.Equal(t, "id", resp.Frames[0].Fields[1].Name)
		require.Equal(t, 7, resp.Frames[0].Fields[0].Len())
		require.Equal(t, 7, resp.Frames[0].Fields[1].Len())
		require.Equal(t, "1611019111439-0", resp.Frames[0].Fields[0].At(0))
		require.Equal(t, "1611019111985-0", resp.Frames[0].Fields[0].At(6))
	})
}

/**
 * XREVRANGE
 */

func TestXRevRangeStreamIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	t.Run("query stream queue:customers", func(t *testing.T) {
		resp := queryXRange(queryModel{Key: "queue:customers"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 2)
		require.Equal(t, "$streamId", resp.Frames[0].Fields[0].Name)
		require.Equal(t, "id", resp.Frames[0].Fields[1].Name)
		require.Equal(t, resp.Frames[0].Fields[1].Len(), resp.Frames[0].Fields[0].Len())
	})

	t.Run("query stream queue:customers with COUNT", func(t *testing.T) {
		resp := queryXRevRange(queryModel{Key: "queue:customers", Count: 3}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 2)
		require.Equal(t, "$streamId", resp.Frames[0].Fields[0].Name)
		require.Equal(t, "id", resp.Frames[0].Fields[1].Name)
		require.Equal(t, 3, resp.Frames[0].Fields[0].Len())
		require.Equal(t, 3, resp.Frames[0].Fields[1].Len())
	})

	t.Run("query stream queue:customers with start and end", func(t *testing.T) {
		resp := queryXRevRange(queryModel{Key: "queue:customers", End: "1611019111985-0", Start: "1611019111439-0"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 2)
		require.Equal(t, "$streamId", resp.Frames[0].Fields[0].Name)
		require.Equal(t, "id", resp.Frames[0].Fields[1].Name)
		require.Equal(t, 7, resp.Frames[0].Fields[0].Len())
		require.Equal(t, 7, resp.Frames[0].Fields[1].Len())
		require.Equal(t, "1611019111439-0", resp.Frames[0].Fields[0].At(6))
		require.Equal(t, "1611019111985-0", resp.Frames[0].Fields[0].At(0))
	})
}

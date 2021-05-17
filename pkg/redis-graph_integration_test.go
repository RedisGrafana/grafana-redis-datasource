// +build integration

package main

import (
	"fmt"
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/require"
)

/**
 * GRAPH.QUERY
 */
func TestGraphQueryIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryGraphQuery(queryModel{Command: "graph.query", Key: "GOT_DEMO", Cypher: "MATCH (w:writer)-[r:wrote]->(b:book) return w,r,b"}, &client)
	require.Len(t, resp.Frames, 4)
	require.Len(t, resp.Frames[0].Fields, 5)
	require.Equal(t, "id", resp.Frames[0].Fields[0].Name)
	require.Equal(t, "title", resp.Frames[0].Fields[1].Name)
	require.Equal(t, "subTitle", resp.Frames[0].Fields[2].Name)
	require.Equal(t, "mainStat", resp.Frames[0].Fields[3].Name)
	require.Equal(t, "arc__", resp.Frames[0].Fields[4].Name)
	require.Equal(t, 15, resp.Frames[0].Fields[0].Len())
	require.Len(t, resp.Frames[1].Fields, 4)
	require.Equal(t, "id", resp.Frames[1].Fields[0].Name)
	require.Equal(t, "source", resp.Frames[1].Fields[1].Name)
	require.Equal(t, "target", resp.Frames[1].Fields[2].Name)
	require.Equal(t, "mainStat", resp.Frames[1].Fields[3].Name)
	require.Equal(t, 14, resp.Frames[1].Fields[0].Len())
}

func TestGraphQueryIntegrationWithoutRelations(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryGraphQuery(queryModel{Command: "graph.query", Key: "GOT_DEMO", Cypher: "MATCH (w:writer)-[wrote]->(b:book) return w,b"}, &client)
	require.Len(t, resp.Frames, 3)
	require.Len(t, resp.Frames[0].Fields, 5)
	require.Equal(t, 15, resp.Frames[0].Fields[0].Len())
	require.Len(t, resp.Frames[1].Fields, 4)
	require.Equal(t, 0, resp.Frames[1].Fields[0].Len())
}

func TestGraphQueryIntegrationWithoutNodes(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryGraphQuery(queryModel{Command: "graph.query", Key: "GOT_DEMO", Cypher: "MATCH (w:writer)-[r:wrote]->(b:book) return r"}, &client)
	require.Len(t, resp.Frames, 3)
	require.Len(t, resp.Frames[0].Fields, 5)
	require.Equal(t, 0, resp.Frames[0].Fields[0].Len())
	require.Len(t, resp.Frames[1].Fields, 4)
	require.Equal(t, 14, resp.Frames[1].Fields[0].Len())
}

/**
 * GRAPH.SLOWLOG
 */
func TestGraphSlowlogIntegration(t *testing.T) {
	// Client
	radixClient, _ := radix.NewPool("tcp", fmt.Sprintf("127.0.0.1:%d", integrationTestPort), 10)
	client := radixV3Impl{radixClient: radixClient}

	// Response
	resp := queryGraphSlowlog(queryModel{Command: "graph.slowlog", Key: "GOT_DEMO"}, &client)
	require.Len(t, resp.Frames, 1)
	require.Len(t, resp.Frames[0].Fields, 4)
}

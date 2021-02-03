package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

/**
 * GRAPH.QUERY
 */
func TestGraphQuery(t *testing.T) {
	t.Parallel()

	/**
	 * Success
	 */
	t.Run("should process command", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: []interface{}{
				nil,
				[]interface{}{
					// nodes + relations
					// Entry 1
					[]interface{}{
						// node
						[]interface{}{
							// id Array
							[]interface{}{
								[]byte("id"),
								int64(333),
							},
							// Labels Array
							[]interface{}{
								[]byte("labels"),
								[]interface{}{
									[]byte("writer"),
								},
							},
							// Properties Array
							[]interface{}{
								[]byte("properties"),
								[]interface{}{
									[]interface{}{[]byte("name"), []byte("Mark Twain")},
								},
							},
						},
						// node
						[]interface{}{
							// id Array
							[]interface{}{
								[]byte("id"),
								int64(444),
							},
							// Labels Array
							[]interface{}{
								[]byte("labels"),
								[]interface{}{
									[]byte("book"),
								},
							},
							// Properties Array
							[]interface{}{
								[]byte("properties"),
								[]interface{}{
									[]interface{}{[]byte("name"), []byte("The Adventures of Tom Sawyer")},
								},
							},
						},
						// relations
						[]interface{}{
							// id Array
							[]interface{}{
								[]byte("id"),
								int64(1),
							},
							// type Array
							[]interface{}{
								[]byte("type"),
								[]byte("wrote"),
							},
							// src Array
							[]interface{}{
								[]byte("src_node"),
								int64(333),
							},
							// dest Array
							[]interface{}{
								[]byte("dest_node"),
								int64(444),
							},
							// Properties Array
							[]interface{}{
								[]byte("properties"),
								[]interface{}{
									[]interface{}{[]byte("relation"), []byte("author")},
								},
							},
						},
					},
					// Entry 2
					[]interface{}{
						// node
						[]interface{}{
							// id Array
							[]interface{}{
								[]byte("id"),
								int64(333),
							},
							// Labels Array
							[]interface{}{
								[]byte("labels"),
								[]interface{}{
									[]byte("writer"),
								},
							},
							// Properties Array
							[]interface{}{
								[]byte("properties"),
								[]interface{}{
									[]interface{}{[]byte("name"), []byte("Mark Twain")},
								},
							},
						},
						// node
						[]interface{}{
							// id Array
							[]interface{}{
								[]byte("id"),
								int64(555),
							},
							// Labels Array
							[]interface{}{
								[]byte("labels"),
								[]interface{}{
									[]byte("book"),
								},
							},
							// Properties Array
							[]interface{}{
								[]byte("properties"),
								[]interface{}{
									[]interface{}{[]byte("count"), int64(10)},
								},
							},
						},
						// relations
						[]interface{}{
							// id Array
							[]interface{}{
								[]byte("id"),
								int64(2),
							},
							// type Array
							[]interface{}{
								[]byte("type"),
								[]byte("wrote"),
							},
							// src Array
							[]interface{}{
								[]byte("src_node"),
								int64(333),
							},
							// dest Array
							[]interface{}{
								[]byte("dest_node"),
								int64(555),
							},
							// Properties Array
							[]interface{}{
								[]byte("properties"),
								[]interface{}{
									[]interface{}{[]byte("relation"), []byte("author")},
								},
							},
						},
					},
				},
				nil,
			},
			err: nil,
		}

		// Response
		resp := queryGraphQuery(queryModel{Command: "graph.query", Key: "GOT_DEMO", Cypher: "MATCH (w:writer)-[r:wrote]->(b:book) return w,r,b"}, &client)
		require.Len(t, resp.Frames, 2)
		require.Len(t, resp.Frames[0].Fields, 5)
		require.Equal(t, "id", resp.Frames[0].Fields[0].Name)
		require.Equal(t, "333", resp.Frames[0].Fields[0].At(0))
		require.Equal(t, "writer", resp.Frames[0].Fields[1].At(0))
		require.Equal(t, "", resp.Frames[0].Fields[2].At(0))
		require.Equal(t, "Mark Twain", resp.Frames[0].Fields[3].At(0))
		require.Equal(t, int64(1), resp.Frames[0].Fields[4].At(0))
		require.Equal(t, "444", resp.Frames[0].Fields[0].At(1))
		require.Equal(t, "555", resp.Frames[0].Fields[0].At(2))
		require.Equal(t, "title", resp.Frames[0].Fields[1].Name)
		require.Equal(t, "subTitle", resp.Frames[0].Fields[2].Name)
		require.Equal(t, "mainStat", resp.Frames[0].Fields[3].Name)
		require.Equal(t, "arc__", resp.Frames[0].Fields[4].Name)
		require.Equal(t, 3, resp.Frames[0].Fields[0].Len())
		require.Len(t, resp.Frames[1].Fields, 4)
		require.Equal(t, "id", resp.Frames[1].Fields[0].Name)
		require.Equal(t, "1", resp.Frames[1].Fields[0].At(0))
		require.Equal(t, "333", resp.Frames[1].Fields[1].At(0))
		require.Equal(t, "444", resp.Frames[1].Fields[2].At(0))
		require.Equal(t, "wrote", resp.Frames[1].Fields[3].At(0))
		require.Equal(t, "2", resp.Frames[1].Fields[0].At(1))
		require.Equal(t, "source", resp.Frames[1].Fields[1].Name)
		require.Equal(t, "target", resp.Frames[1].Fields[2].Name)
		require.Equal(t, "mainStat", resp.Frames[1].Fields[3].Name)
		require.Equal(t, 2, resp.Frames[1].Fields[0].Len())

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
		resp := queryGraphQuery(queryModel{Command: "graph.query", Key: "GOT_DEMO", Cypher: "MATCH (w:writer)-[r:wrote]->(b:book) return w,r,b"}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})
}

/**
 * GRAPH.SLOWLOG
 */
func TestGraphSlowlog(t *testing.T) {
	t.Parallel()

	/**
	 * Success
	 */
	t.Run("should process command", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: [][]string{
				{"1612352919", "GRAPH.QUERY", "MATCH (w:writer)-[wrote]->(b:book) return w,r,b", "0.929"},
			},
			err: nil,
		}

		// Response
		resp := queryGraphSlowlog(queryModel{Command: "graph.slowlog", Key: "GOT_DEMO"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 4)
		require.Equal(t, "timestamp", resp.Frames[0].Fields[0].Name)
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
		require.Equal(t, "command", resp.Frames[0].Fields[1].Name)
		require.Equal(t, "query", resp.Frames[0].Fields[2].Name)
		require.Equal(t, "duration", resp.Frames[0].Fields[3].Name)

		require.Equal(t, int64(1612352919), resp.Frames[0].Fields[0].At(0))
		require.Equal(t, "GRAPH.QUERY", resp.Frames[0].Fields[1].At(0))
		require.Equal(t, "MATCH (w:writer)-[wrote]->(b:book) return w,r,b", resp.Frames[0].Fields[2].At(0))
		require.Equal(t, float64(0.929), resp.Frames[0].Fields[3].At(0))

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
		resp := queryGraphSlowlog(queryModel{Command: "graph.slowlog", Key: "GOT_DEMO"}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})
}

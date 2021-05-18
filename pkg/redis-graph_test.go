package main

import (
	"errors"
	"testing"
	"time"

	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

/**
 * GRAPH.QUERY
 */
func TestGraphQuery(t *testing.T) {
	t.Parallel()

	/**
	 * Results with Nodes and Edges
	 */
	t.Run("should process command", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: []interface{}{
				[]interface{}{
					[]byte("w"),
					[]byte("r"),
					[]byte("b"),
				},
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
									[]interface{}{[]byte("count"), int64(10)},
								},
							},
						},
					},
				},
				[]interface{}{
					[]byte("Cached execution: 1"),
					[]byte("Query internal execution time: 0.402967 milliseconds"),
				},
			},
			err: nil,
		}

		// Response
		resp := queryGraphQuery(queryModel{Command: models.GraphQuery, Key: "GOT_DEMO", Cypher: "MATCH (w:writer)-[r:wrote]->(b:book) return w,r,b"}, &client)
		require.Len(t, resp.Frames, 4)
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
	 * Results with data only
	 */
	t.Run("should process command with data", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: []interface{}{
				[]interface{}{
					[]byte("r.name"),
					[]byte("t.name"),
					[]byte("t.gp"),
					[]byte("t.float"),
				},
				[]interface{}{
					// Entry 1
					[]interface{}{
						[]byte("Green Den"),
						[]byte("27 semiprecious stones"),
						int64(225),
						float64(3.14),
					},
					// Entry 2
					[]interface{}{
						[]byte("Green Warren"),
						[]byte("sack of coins"),
						int64(500),
						float64(15),
					},
				},
				[]interface{}{
					[]byte("Cached execution: 0"),
					[]byte("Query internal execution time: 0.524894 milliseconds"),
					[]byte("Should skip this line"),
				},
			},
			err: nil,
		}

		// Response
		resp := queryGraphQuery(queryModel{Command: models.GraphQuery, Key: "dungeon", Cypher: "MATCH (r:Room)-[:CONTAINS]->(t:Treasure) RETURN r.name, t.name, t.gp, t.float"}, &client)
		require.Len(t, resp.Frames, 2)
		require.Len(t, resp.Frames[0].Fields, 4)

		// Data
		require.Equal(t, "r.name", resp.Frames[0].Fields[0].Name)
		require.Equal(t, "Green Den", resp.Frames[0].Fields[0].At(0))
		require.Equal(t, "t.gp", resp.Frames[0].Fields[2].Name)
		require.Equal(t, int64(225), resp.Frames[0].Fields[2].At(0))
		require.Equal(t, "t.float", resp.Frames[0].Fields[3].Name)
		require.Equal(t, float64(3.14), resp.Frames[0].Fields[3].At(0))
		require.Equal(t, 2, resp.Frames[0].Fields[0].Len())

		// Meta
		require.Equal(t, "data", resp.Frames[1].Fields[0].Name)
		require.Equal(t, "Cached execution", resp.Frames[1].Fields[0].At(0))
		require.Equal(t, "value", resp.Frames[1].Fields[1].Name)
		require.Equal(t, "0", resp.Frames[1].Fields[1].At(0))
		require.Equal(t, 2, resp.Frames[1].Fields[0].Len())
	})

	/**
	 * Results with no data
	 */
	t.Run("should process command with data", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: []interface{}{
				[]interface{}{
					[]byte("r.name"),
					[]byte("t.name"),
					[]byte("t.gp"),
					[]byte("t.float"),
				},
				[]interface{}{[]interface{}{}},
				[]interface{}{
					[]byte("Cached execution: 0"),
					[]byte("Query internal execution time: 0.524894 milliseconds"),
				},
			},
			err: nil,
		}

		// Response
		resp := queryGraphQuery(queryModel{Command: models.GraphQuery, Key: "dungeon", Cypher: "MATCH (r:Room)-[:CONTAINS]->(t:Treasure) RETURN r.name, t.name, t.gp, t.float"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 2)
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
		resp := queryGraphQuery(queryModel{Command: models.GraphQuery, Key: "GOT_DEMO", Cypher: "MATCH (w:writer)-[r:wrote]->(b:book) return w,r,b"}, &client)
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
		resp := queryGraphSlowlog(queryModel{Command: models.GraphSlowlog, Key: "GOT_DEMO"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 4)
		require.Equal(t, "timestamp", resp.Frames[0].Fields[0].Name)
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
		require.Equal(t, "command", resp.Frames[0].Fields[1].Name)
		require.Equal(t, "query", resp.Frames[0].Fields[2].Name)
		require.Equal(t, "duration", resp.Frames[0].Fields[3].Name)

		require.Equal(t, time.Unix(1612352919, 0), resp.Frames[0].Fields[0].At(0))
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
		resp := queryGraphSlowlog(queryModel{Command: models.GraphSlowlog, Key: "GOT_DEMO"}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})
}

/**
 * GRAPH.EXPLAIN
 */
func TestGraphExplain(t *testing.T) {
	t.Parallel()

	/**
	 * Success
	 */
	t.Run("should process command", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: []string{
				"Results     Project         Conditional Traverse | (w:writer)-[r:wrote]->(b:book)             Node By Label Scan | (w:writer)",
			},
			err: nil,
		}

		// Response
		resp := queryGraphExplain(queryModel{Command: models.GraphExplain, Key: "GOT_DEMO", Cypher: "MATCH (r:Room)-[:CONTAINS]->(t:Treasure) RETURN r.name, t.name, t.gp"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 1)
		require.Equal(t, "execution plan", resp.Frames[0].Fields[0].Name)
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
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
		resp := queryGraphExplain(queryModel{Command: models.GraphExplain, Key: "GOT_DEMO", Cypher: "MATCH (r:Room)-[:CONTAINS]->(t:Treasure) RETURN r.name, t.name, t.gp"}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})
}

/**
 * GRAPH.PROFILE
 */
func TestGraphProfile(t *testing.T) {
	t.Parallel()

	/**
	 * Success
	 */
	t.Run("should process command", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: []string{
				"Results | Records produced: 17, Execution time: 0.003256 ms",
				"Project | Records produced: 17, Execution time: 0.019947 ms",
				"Conditional Traverse | (r:Room)->(t:Treasure) | Records produced: 17, Execution time: 0.161637 ms",
				"Conditional Traverse | (r:Room)->(t:Treasure) | Records produced: 17, Execution time: 1 s",
			},
			err: nil,
		}

		// Response
		resp := queryGraphProfile(queryModel{Command: models.GraphProfile, Key: "GOT_DEMO", Cypher: "MATCH (r:Room)-[:CONTAINS]->(t:Treasure) RETURN r.name, t.name, t.gp"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 3)
		require.Equal(t, "operation", resp.Frames[0].Fields[0].Name)
		require.Equal(t, "Results", resp.Frames[0].Fields[0].At(0))
		require.Equal(t, "records", resp.Frames[0].Fields[1].Name)
		require.Equal(t, int64(17), resp.Frames[0].Fields[1].At(0))
		require.Equal(t, "execution", resp.Frames[0].Fields[2].Name)
		require.Equal(t, float64(0.003256), resp.Frames[0].Fields[2].At(0))
		require.Equal(t, 3, resp.Frames[0].Fields[0].Len())
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
		resp := queryGraphProfile(queryModel{Command: models.GraphProfile, Key: "GOT_DEMO", Cypher: "MATCH (r:Room)-[:CONTAINS]->(t:Treasure) RETURN r.name, t.name, t.gp"}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})
}

/**
 * GRAPH.CONFIG
 */
func TestGraphConfig(t *testing.T) {
	t.Parallel()

	/**
	 * Success
	 */
	t.Run("should process command", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: []interface{}{
				[]interface{}{[]byte("CACHE_SIZE"), int64(25)},
				[]interface{}{[]byte("ASYNC_DELETE"), int64(1)},
				[]interface{}{[]byte("THREAD_COUNT"), float64(3.14)},
			},
			err: nil,
		}

		// Response
		resp := queryGraphConfig(queryModel{Command: models.GraphConfig}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 3)
		require.Equal(t, "CACHE_SIZE", resp.Frames[0].Fields[0].Name)
		require.Equal(t, int64(25), resp.Frames[0].Fields[0].At(0))
		require.Equal(t, "THREAD_COUNT", resp.Frames[0].Fields[2].Name)
		require.Equal(t, float64(3.14), resp.Frames[0].Fields[2].At(0))
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
		resp := queryGraphConfig(queryModel{Command: models.GraphConfig}, &client)
		require.EqualError(t, resp.Error, "error occurred")
	})
}

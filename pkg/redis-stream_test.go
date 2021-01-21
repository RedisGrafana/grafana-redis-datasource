package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

/**
 * XINFO STREAM
 */
func TestQueryXInfoStream(t *testing.T) {
	t.Parallel()

	t.Run("should handle response with FirstEntry and LasEntry", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: xinfo{
				Length:          5,
				RadixTreeKeys:   4,
				RadixTreeNodes:  3,
				Groups:          2,
				LastGeneratedId: "id",
				FirstEntry: []interface{}{
					[]byte("id2"),
					[]interface{}{
						[]byte("key1"),
						[]byte("value1"),
						[]byte("key2"),
						[]byte("value2"),
					},
				},
				LastEntry: []interface{}{
					[]byte("id3"),
					[]interface{}{
						[]byte("key3"),
						[]byte("value3"),
						[]byte("key4"),
						[]byte("value4"),
					},
				},
			},
		}

		// Response
		resp := queryXInfoStream(queryModel{Command: "xinfoStream", Key: "test1"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 9)
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
		require.Equal(t, "length", resp.Frames[0].Fields[0].Name)
		require.Equal(t, int64(5), resp.Frames[0].Fields[0].At(0))
		require.Equal(t, "radix-tree-keys", resp.Frames[0].Fields[1].Name)
		require.Equal(t, int64(4), resp.Frames[0].Fields[1].At(0))
		require.Equal(t, "radix-tree-nodes", resp.Frames[0].Fields[2].Name)
		require.Equal(t, int64(3), resp.Frames[0].Fields[2].At(0))
		require.Equal(t, "groups", resp.Frames[0].Fields[3].Name)
		require.Equal(t, int64(2), resp.Frames[0].Fields[3].At(0))
		require.Equal(t, "last-generated-id", resp.Frames[0].Fields[4].Name)
		require.Equal(t, "id", resp.Frames[0].Fields[4].At(0))
		require.Equal(t, "first-entry-id", resp.Frames[0].Fields[5].Name)
		require.Equal(t, "id2", resp.Frames[0].Fields[5].At(0))
		require.Equal(t, "first-entry-fields", resp.Frames[0].Fields[6].Name)
		require.Equal(t, "\"key1\"=\"value1\"\n\"key2\"=\"value2\"\n", resp.Frames[0].Fields[6].At(0))
		require.Equal(t, "last-entry-id", resp.Frames[0].Fields[7].Name)
		require.Equal(t, "id3", resp.Frames[0].Fields[7].At(0))
		require.Equal(t, "last-entry-fields", resp.Frames[0].Fields[8].Name)
		require.Equal(t, "\"key3\"=\"value3\"\n\"key4\"=\"value4\"\n", resp.Frames[0].Fields[8].At(0))
	})

	t.Run("should handle response without FirstEntry and LasEntry", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{rcv: xinfo{
			Length:          5,
			RadixTreeKeys:   4,
			RadixTreeNodes:  3,
			Groups:          2,
			LastGeneratedId: "id",
			FirstEntry:      nil,
			LastEntry:       nil,
		}}

		// Response
		resp := queryXInfoStream(queryModel{Command: "xinfoStream", Key: "test1"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 5)
		require.Equal(t, 1, resp.Frames[0].Fields[0].Len())
	})

	t.Run("should handle rerror", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{err: errors.New("some error")}

		// Response
		resp := queryXInfoStream(queryModel{Command: "xinfoStream", Key: "test1"}, &client)
		require.Len(t, resp.Frames, 0)
		require.EqualError(t, resp.Error, "some error")
	})
}

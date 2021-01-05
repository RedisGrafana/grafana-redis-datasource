package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryFtInfo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		qm                      queryModel
		rcv                     interface{}
		fieldsCount             int
		rowsPerField            int
		valuesToCheckInResponse []valueToCheckInResponse
		err                     error
	}{
		{
			"should parse default bulk string",
			queryModel{Command: "ft.info", Key: "wik{0}"},
			[]interface{}{
				"index_name",
				"wikipedia",
				"index_options",
				[]interface{}{},
				"index_definition",
				[]interface{}{
					"key_type",
					"HASH",
					"prefixes",
					[]interface{}{"thing"},
					"filter",
					"startswith(@__key, \"thing:\")",
					"language_field",
					"__language",
					"default_score",
					"1",
					"score_field",
					"__score",
					"payload_field",
					"__payload",
				},
				"fields",
				[]interface{}{
					[]interface{}{
						"title",
						"type",
						"TEXT",
						"WEIGHT",
						"1",
						"SORTABLE",
					},
					[]interface{}{
						"body",
						"type",
						"TEXT",
						"WEIGHT",
					},
					[]interface{}{
						"id",
						"type",
						"NUMERIC",
					},
					[]interface{}{
						"subject location",
						"type",
						"GEO",
					},
				},
				"num_docs",
				"0",
				"max_doc_id",
				"345678",
				"num_terms",
				"691356",
				"num_records",
				"0",
				"inverted_sz_mb",
				"0",
				"total_inverted_index_blocks",
				"933290",
				"offset_vectors_sz_mb",
				"0.65932846069335938",
				"doc_table_size_mb",
				"29.893482208251953",
				"sortable_values_size_mb",
				"11.432285308837891",
				"key_table_size_mb",
				"1.239776611328125e-05",
				"records_per_doc_avg",
				"-nan",
				"bytes_per_record_avg",
				"-nan",
				"offsets_per_term_avg",
				"inf",
				"offset_bits_per_record_avg",
				"8",
				"hash_indexing_failures",
				"0",
				"indexing",
				"0",
				"percent_indexed",
				"1",
				"gc_stats",
				[]interface{}{
					"bytes_collected",
					"4148136",
					"total_ms_run",
					"14796",
					"total_cycles",
					"1",
					"average_cycle_time_ms",
					"14796",
					"last_run_time_ms",
					"14796",
					"gc_numeric_trees_missed",
					"0",
					"gc_blocks_denied",
					"0",
				},
				"cursor_stats",
				[]interface{}{
					"global_idle",
					int64(0),
					"global_total",
					int64(0),
					"index_capacity",
					int64(128),
					"index_total",
					int64(0),
				},
				"stopwords_list",
				[]interface{}{
					"tlv",
					"summer",
					"2020",
				},
			},
			18,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "wikipedia"},
				{frameIndex: 0, fieldIndex: 3, rowIndex: 0, value: float64(691356)},
				{frameIndex: 0, fieldIndex: 7, rowIndex: 0, value: float64(0.6593284606933594)},
				{frameIndex: 0, fieldIndex: 7, rowIndex: 0, value: float64(0.6593284606933594)},
			},
			nil,
		},
		{
			"should handle error",
			queryModel{Command: "info"},
			nil,
			0,
			0,
			nil,
			errors.New("error occurred"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ds := redisDatasource{}
			client := TestClient{tt.rcv, tt.err}
			response := ds.queryFtInfo(tt.qm, client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Key, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created from bulk string")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")
				require.NoError(t, response.Error, "Should not return error")
				if tt.valuesToCheckInResponse != nil {
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				}
			}
		})
	}
}

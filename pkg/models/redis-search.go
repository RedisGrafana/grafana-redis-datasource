package models

/**
 * RediSearch Commands
 */
const (
	SearchInfo = "ft.info"
	Search     = "ft.search"
)

/**
 * FT.INFO field configuration
 */
var SearchInfoConfig = map[string]string{
	"inverted_sz_mb":          "decmbytes",
	"offset_vectors_sz_mb":    "decmbytes",
	"doc_table_size_mb":       "decmbytes",
	"sortable_values_size_mb": "decmbytes",
	"key_table_size_mb":       "decmbytes",
	"percent_indexed":         "percentunit",
}

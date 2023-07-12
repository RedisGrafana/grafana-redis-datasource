package main

import (
	"errors"
	"testing"

	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

/**
 * Type and Length commands with Key and Path
 */
func TestQueryJsonObjLen(t *testing.T) {
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
			"should handle string value",
			queryModel{Command: models.JsonObjLen, Key: "test:json", Path: "."},
			"someStr",
			1,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "someStr"},
			},
			nil,
		},
		{
			"should handle float64 value",
			queryModel{Command: models.JsonObjLen, Key: "test:json", Path: "."},
			"3.14",
			1,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: 3.14},
			},
			nil,
		},
		{
			"should handle error",
			queryModel{Command: models.JsonObjLen},
			nil,
			0,
			0,
			nil,
			errors.New("error occurred"),
		},
	}

	// Run Tests
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Client
			client := testClient{rcv: tt.rcv, err: tt.err}

			// Response
			response := queryJsonObjLen(tt.qm, &client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Key, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")

				if tt.valuesToCheckInResponse != nil {
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				}
			}
		})
	}
}

/**
 * JSON.OBJKEYS
 */
func TestQueryJsonObjKeys(t *testing.T) {
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
			"should handle default array of strings",
			queryModel{Command: models.JsonObjKeys, Key: "test:json", Path: "."},
			[]string{"value1", "2", "3.14"},
			1,
			3,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "value1"},
				{frameIndex: 0, fieldIndex: 0, rowIndex: 1, value: "2"},
				{frameIndex: 0, fieldIndex: 0, rowIndex: 2, value: "3.14"},
			},
			nil,
		},
		{
			"should handle error",
			queryModel{Command: models.JsonObjKeys},
			nil,
			0,
			0,
			nil,
			errors.New("error occurred"),
		},
	}

	// Run Tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Client
			client := testClient{rcv: tt.rcv, err: tt.err}

			// Response
			response := queryJsonObjKeys(tt.qm, &client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Key, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")

				if tt.valuesToCheckInResponse != nil {
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				}
			}
		})
	}
}

/**
 * JSON.GET
 */
func TestQueryJsonGet(t *testing.T) {
	t.Parallel()

	t.Run("Should return four strings in frame", func(t *testing.T) {
		t.Parallel()

		client := testClient{rcv: "[[],\"gin\",\"rum\",\"whiskey\"]"}

		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "$.num"}, &client)

		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 1)
		require.Equal(t, resp.Frames[0].Fields[0].Len(), 4)
		require.Equal(t, resp.Frames[0].Fields[0].At(0), "")
		require.Equal(t, resp.Frames[0].Fields[0].At(1), "gin")
		require.Equal(t, resp.Frames[0].Fields[0].At(2), "rum")
		require.Equal(t, resp.Frames[0].Fields[0].At(3), "whiskey")
	})

	t.Run("Should return four booleans in frame", func(t *testing.T) {
		t.Parallel()

		client := testClient{rcv: "[[],true,false,true]"}

		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "$.num"}, &client)

		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 1)
		require.Equal(t, resp.Frames[0].Fields[0].Len(), 4)
		require.Equal(t, resp.Frames[0].Fields[0].At(0), false)
		require.Equal(t, resp.Frames[0].Fields[0].At(1), true)
		require.Equal(t, resp.Frames[0].Fields[0].At(2), false)
		require.Equal(t, resp.Frames[0].Fields[0].At(3), true)
	})

	t.Run("Should return four float64 in frame", func(t *testing.T) {
		t.Parallel()

		client := testClient{rcv: "[[],42,43,44]"}

		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "$.num"}, &client)

		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 1)
		require.Equal(t, resp.Frames[0].Fields[0].Len(), 4)
		require.Equal(t, resp.Frames[0].Fields[0].At(0), float64(0))
		require.Equal(t, resp.Frames[0].Fields[0].At(1), float64(42))
		require.Equal(t, resp.Frames[0].Fields[0].At(2), float64(43))
		require.Equal(t, resp.Frames[0].Fields[0].At(3), float64(44))

	})

	t.Run("Should return a single float64 in frame", func(t *testing.T) {
		t.Parallel()

		client := testClient{rcv: "[42]"}

		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "$.num"}, &client)

		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 1)
		require.Equal(t, resp.Frames[0].Fields[0].At(0), float64(42))

	})

	t.Run("Should return a single boolean in frame", func(t *testing.T) {
		t.Parallel()

		client := testClient{rcv: "[true]"}
		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "$.bool"}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 1)
		require.Equal(t, resp.Frames[0].Fields[0].At(0), true)
	})

	t.Run("should handle encoded JSON with string", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: "{\"name\":\"Leonard Cohen\",\"lastSeen\":1478476800,\"loggedOut\":true, \"key3\":3.14}",
		}

		// Response
		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "."}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 4)
	})

	t.Run("should handle encoded JSON with string", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: "\"test\"",
		}

		// Response
		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "."}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 1)
	})

	t.Run("should handle encoded JSON with bool", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: "true",
		}

		// Response
		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "."}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 1)
	})

	t.Run("should handle encoded JSON with array of objects", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: "[{\"timestamp\":1559174400,\"redislabs\\/redisearch\":66216,\"redislabs\\/redisgraph\":33133,\"redislabs\\/rebloom\":5922,\"redislabs\\/rejson\":43140},{\"timestamp\":1560211200,\"redislabs\\/redisearch\":71615,\"redislabs\\/redisgraph\":34217,\"redislabs\\/rebloom\":6144,\"redislabs\\/redistimeseries\":1123,\"redislabs\\/rejson\":45050,\"redisai\\/redisai\":926},{\"timestamp\":1560729600,\"redislabs\\/redisearch\":74351,\"redislabs\\/redisgraph\":34779,\"redislabs\\/rebloom\":6312,\"redislabs\\/redistimeseries\":1271,\"redislabs\\/rejson\":47177,\"redisai\\/redisai\":1002},{\"timestamp\":1562025600,\"redislabs\\/redisearch\":82291,\"redislabs\\/redisgraph\":36147,\"redislabs\\/rebloom\":6819,\"redislabs\\/redistimeseries\":1562,\"redislabs\\/rejson\":52525,\"redisai\\/redisai\":1227},{\"timestamp\":1562544000,\"redislabs\\/redisearch\":83740,\"redislabs\\/redisgraph\":36478,\"redislabs\\/rebloom\":7014,\"redislabs\\/redistimeseries\":1653,\"redislabs\\/rejson\":54278,\"redisai\\/redisai\":1263},{\"timestamp\":1565740800,\"redislabs\\/redisearch\":91943,\"redislabs\\/redisgraph\":40232,\"redislabs\\/rebloom\":8022,\"redislabs\\/redistimeseries\":2260,\"redislabs\\/rejson\":67772,\"redisai\\/redisai\":1634},{\"timestamp\":1572134400,\"redislabs\\/redisearch\":121171,\"redislabs\\/redisgraph\":70668,\"redislabs\\/rebloom\":10411,\"redislabs\\/redistimeseries\":3466,\"redislabs\\/rejson\":83498,\"redisai\\/redisai\":13368},{\"timestamp\":1573948800,\"redislabs\\/redisearch\":136122,\"redislabs\\/redisgraph\":79692,\"redislabs\\/rebloom\":11325,\"redislabs\\/redistimeseries\":3823,\"redislabs\\/rejson\":86438,\"redisai\\/redisai\":13517},{\"timestamp\":1577491200,\"redislabs\\/redisearch\":157442,\"redislabs\\/redisgraph\":102747,\"redislabs\\/rebloom\":12463,\"redislabs\\/redistimeseries\":4802,\"redislabs\\/rejson\":94166,\"redisai\\/redisai\":13785},{\"timestamp\":1582243200,\"redislabs\\/redisearch\":195170,\"redislabs\\/redisgraph\":135159,\"redislabs\\/rebloom\":14869,\"redislabs\\/redistimeseries\":6940,\"redislabs\\/rejson\":126303,\"redislabs\\/redisgears\":1223,\"redisai\\/redisai\":14609},{\"timestamp\":1585065499.89136,\"redislabs\\/redisearch\":231666,\"redislabs\\/redisgraph\":144947,\"redislabs\\/rebloom\":16432,\"redislabs\\/redistimeseries\":9391,\"redislabs\\/rejson\":207763,\"redislabs\\/redisjson2\":2238,\"redislabs\\/redisgears\":1859,\"redisai\\/redisai\":15162},{\"timestamp\":1585066940.058363,\"redislabs\\/redisearch\":231693,\"redislabs\\/redisgraph\":145004,\"redislabs\\/rebloom\":16433,\"redislabs\\/redistimeseries\":9394,\"redislabs\\/rejson\":207808,\"redislabs\\/redisjson2\":2238,\"redislabs\\/redisgears\":1860,\"redisai\\/redisai\":15164},{\"timestamp\":1585066997.3418889,\"redislabs\\/redisearch\":231693,\"redislabs\\/redisgraph\":145004,\"redislabs\\/rebloom\":16433,\"redislabs\\/redistimeseries\":9395,\"redislabs\\/rejson\":207812,\"redislabs\\/redisjson2\":2238,\"redislabs\\/redisgears\":1860,\"redisai\\/redisai\":15164},{\"timestamp\":1585067702.008811,\"redislabs\\/redisearch\":231707,\"redislabs\\/redisgraph\":145031,\"redislabs\\/rebloom\":16433,\"redislabs\\/redistimeseries\":9396,\"redislabs\\/rejson\":207834,\"redislabs\\/redisjson2\":2238,\"redislabs\\/redisgears\":1860,\"redisai\\/redisai\":15164},{\"timestamp\":1585070847.8951271,\"redislabs\\/redisearch\":231759,\"redislabs\\/redisgraph\":145088,\"redislabs\\/rebloom\":16433,\"redislabs\\/redistimeseries\":9406,\"redislabs\\/rejson\":207933,\"redislabs\\/redisjson2\":2238,\"redislabs\\/redisgears\":1861,\"redisai\\/redisai\":15164},{\"timestamp\":1585071149.7625909,\"redislabs\\/redisearch\":231763,\"redislabs\\/redisgraph\":145088,\"redislabs\\/rebloom\":16433,\"redislabs\\/redistimeseries\":9407,\"redislabs\\/rejson\":207944,\"redislabs\\/redisjson2\":2238,\"redislabs\\/redisgears\":1861,\"redisai\\/redisai\":15164},{\"timestamp\":1585071829.6869831,\"redislabs\\/redisearch\":231775,\"redislabs\\/redisgraph\":145089,\"redislabs\\/rebloom\":16434,\"redislabs\\/redistimeseries\":9410,\"redislabs\\/rejson\":207972,\"redislabs\\/redisjson2\":2238,\"redislabs\\/redisgears\":1862,\"redisai\\/redisai\":15165}]",
		}

		// Response
		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "."}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 9)
	})

	t.Run("should handle encoded JSON with array of objects", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: "[{\"timestamp\":4,\"key1\":\"test\",\"key2\":3.3,\"key4\":true},{\"key1\":\"test\",\"key2\":3.3},{\"string\":\"test\"},{\"key2\":3.3},{\"true\":false}]",
		}

		// Response
		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "."}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 6)
	})

	t.Run("should handle encoded JSON with array of strings", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: "[\"string\",\"test\"]",
		}

		// Response
		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "."}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 1)
	})

	t.Run("should handle encoded JSON with array of map of maps", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: "[{\"test\":{\"inside\":{\"string\":\"test\"}}},{\"timestamp\":4,\"key1\":\"test\"},{\"string\":\"test\"},{\"key2\":3.3},{\"true\":false}]",
		}

		// Response
		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "."}, &client)
		require.Len(t, resp.Frames, 1)
		require.Len(t, resp.Frames[0].Fields, 6)
	})

	t.Run("should handle unmarshall error", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{
			rcv: "JSON",
		}

		// Response
		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "."}, &client)
		require.Len(t, resp.Frames, 0)
		require.EqualError(t, resp.Error, "invalid character 'J' looking for beginning of value")
	})

	t.Run("should handle error", func(t *testing.T) {
		t.Parallel()

		// Client
		client := testClient{err: errors.New("some error")}

		// Response
		resp := queryJsonGet(queryModel{Command: models.JsonGet, Key: "test:json", Path: "."}, &client)
		require.Len(t, resp.Frames, 0)
		require.EqualError(t, resp.Error, "some error")
	})
}

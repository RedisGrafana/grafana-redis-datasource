package main

import (
	"errors"
	"testing"
	"time"

	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestQueryTsRange(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		qm                      queryModel
		rcv                     interface{}
		from                    int64
		to                      int64
		fieldsCount             int
		rowsPerField            int
		valuesToCheckInResponse []valueToCheckInResponse
		err                     error
	}{
		{
			"should process receiver without aggregation and legend provided",
			queryModel{Command: models.TimeSeriesRange, Key: "test1"},
			[][]string{
				{"1548149180000", "26.199999999999999"},
				{"1548149195000", "27.399999999999999"},
				{"1548149220000", "24.800000000000001"},
				{"1548149215000", "23.199999999999999"},
				{"1548149230000", "25.199999999999999"},
				{"1548149285000", "28"},
				{"1548149150000", "20"},
			},
			0,
			0,
			2,
			7,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			nil,
		},
		{
			"should process receiver with aggregation and legend",
			queryModel{Command: models.TimeSeriesRange, Aggregation: "avg", Bucket: 5000, Key: "test1", Legend: "Legend"},
			[][]string{
				{"1548149180000", "26.199999999999999"},
				{"1548149185000", "27.399999999999999"},
				{"1548149190000", "24.800000000000001"},
				{"1548149195000", "23.199999999999999"},
				{"1548149200000", "25.199999999999999"},
				{"1548149205000", "28"},
				{"1548149210000", "20"},
			},
			0,
			0,
			2,
			7,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			nil,
		},
		{
			"should process receiver with fill",
			queryModel{Command: models.TimeSeriesRange, Bucket: 5000, Key: "test1", Fill: true},
			[][]string{
				{"1548149180000", "26.199999999999999"},
				{"1548149195000", "27.399999999999999"},
				{"1548149220000", "24.800000000000001"},
				{"1548149215000", "23.199999999999999"},
				{"1548149230000", "25.199999999999999"},
				{"1548149285000", "28"},
				{"1548149150000", "20"},
			},
			0,
			0,
			2,
			25,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 1, value: time.Unix(0, (1548149180000+5000)*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 1, value: float64(0)},
			},
			nil,
		},
		{
			"should process receiver error",
			queryModel{Command: models.TimeSeriesRange, Key: "test1"},
			nil,
			0,
			0,
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

			client := testClient{rcv: tt.rcv, err: tt.err}
			response := queryTsRange(tt.from, tt.to, tt.qm, &client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				if tt.qm.Legend != "" {
					require.Equal(t, tt.qm.Legend, response.Frames[0].Name, "Invalid frame name")
				} else {
					require.Equal(t, tt.qm.Key, response.Frames[0].Name, "Invalid frame name")
				}
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

func TestQueryTsMRange(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		qm                      queryModel
		rcv                     interface{}
		from                    int64
		to                      int64
		fieldsCount             int
		rowsPerField            int
		valuesToCheckInResponse []valueToCheckInResponse
		expectedFrameName       string
		expectedValueFieldName  string
		expectedError           string
		err                     error
	}{
		{
			"should process receiver without aggregation and legend provided but with labels",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
					[]interface{}{
						[]interface{}{int64(1548149180000), []byte("26.199999999999999")},
						[]interface{}{int64(1548149195000), []byte("27.399999999999999")},
						[]interface{}{int64(1548149220000), []byte("24.800000000000001")},
						[]interface{}{int64(1548149215000), []byte("23.199999999999999")},
						[]interface{}{int64(1548149230000), []byte("25.199999999999999")},
						[]interface{}{int64(1548149285000), []byte("28")},
						[]interface{}{int64(1548149150000), []byte("20")},
					},
				},
				[]interface{}{
					[]byte("temperature:3:32"),
					[]interface{}{},
					[]interface{}{
						[]interface{}{int64(1548149180000), []byte("26.7")},
						[]interface{}{int64(1548149195000), []byte("27.8")},
						[]interface{}{int64(1548149220000), []byte("24.4")},
						[]interface{}{int64(1548149215000), []byte("26.199999999999999")},
						[]interface{}{int64(1548149230000), []byte("25.199999999999999")},
						[]interface{}{int64(1548149285000), []byte("27")},
						[]interface{}{int64(1548149150000), []byte("22")},
					},
				},
			},
			0,
			0,
			2,
			7,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			"temperature:2:32",
			"",
			"",
			nil,
		},
		{
			"should process receiver with aggregation and legend",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1", Aggregation: "avg", Legend: "Legend", Bucket: 5000, Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{},
					[]interface{}{
						[]interface{}{int64(1548149180000), []byte("26.199999999999999")},
						[]interface{}{int64(1548149185000), []byte("27.399999999999999")},
						[]interface{}{int64(1548149190000), []byte("24.800000000000001")},
						[]interface{}{int64(1548149195000), []byte("23.199999999999999")},
						[]interface{}{int64(1548149200000), []byte("25.199999999999999")},
						[]interface{}{int64(1548149205000), []byte("28")},
						[]interface{}{int64(1548149210000), []byte("20")},
					},
				},
			},
			0,
			0,
			2,
			7,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			"",
			"",
			"",
			nil,
		},
		{
			"should process receiver with labels specified in legend",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1", Legend: "area_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
					[]interface{}{
						[]interface{}{int64(1548149210000), []byte("20")},
					},
				},
			},
			0,
			0,
			2,
			1,
			nil,
			"32",
			"",
			"",
			nil,
		},
		{
			"should process receiver with value field existed in labels",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1", Value: "sensor_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
					[]interface{}{
						[]interface{}{int64(1548149210000), []byte("20")},
					},
				},
			},
			0,
			0,
			2,
			1,
			nil,
			"temperature:2:32",
			"2",
			"",
			nil,
		},

		{
			"should process receiver with []byte field instead of int",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1", Legend: "area_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
					[]interface{}{
						[]interface{}{[]byte("1548149180000"), []byte("26.199999999999999")},
					},
				},
			},
			0,
			0,
			2,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			"32",
			"",
			"",
			nil,
		},
		{
			"should process receiver with string field instead of int",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1", Legend: "area_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
					[]interface{}{
						[]interface{}{[]byte("1548149180000"), "26.199999999999999"},
					},
				},
			},
			0,
			0,
			2,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			"32",
			"",
			"",
			nil,
		},
		{
			"should process receiver with Fill",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1", Fill: true, Bucket: 5000, Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{},
					[]interface{}{
						[]interface{}{int64(1548149180000), []byte("26.199999999999999")},
						[]interface{}{int64(1548149195000), []byte("27.399999999999999")},
						[]interface{}{int64(1548149220000), []byte("24.800000000000001")},
						[]interface{}{int64(1548149215000), []byte("23.199999999999999")},
						[]interface{}{int64(1548149230000), []byte("25.199999999999999")},
						[]interface{}{int64(1548149285000), []byte("28")},
						[]interface{}{int64(1548149150000), []byte("20")},
					},
				},
			},
			0,
			0,
			2,
			25,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			"temperature:2:32",
			"",
			"",
			nil,
		},
		{
			"test groupby/reduction",
			queryModel{Command: models.TimeSeriesMRange, TsReducer: "SUM", TsGroupByLabel: "reduceLabel"},
			[]interface{}{
				[]interface{}{
					[]byte("foo=bar"),
					[]interface{}{
						[]interface{}{
							[]byte("foo"),
							[]byte("bar"),
						},
						[]interface{}{
							[]byte("__reducer__"),
							[]byte("sum"),
						},
						[]interface{}{
							[]byte("__source__"),
							[]byte("ts:1,ts:2,ts:3"),
						},
					},
					[]interface{}{
						[]interface{}{int64(1686835010300), []byte("2102")},
						[]interface{}{int64(1686835011312), []byte("1882")},
						[]interface{}{int64(1686835013348), []byte("2378")},
						[]interface{}{int64(1686835014362), []byte("3007")},
					},
				},
			},
			1686835010300,
			1686835014362,
			2,
			4,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1686835010300*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: float64(2102)},
			},
			"foo=bar",
			"",
			"",
			nil,
		},
		{"should return error because we missed an actual reducer",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1", Filter: "filter", TsGroupByLabel: "foo"},
			interface{}("someString"),
			0,
			0,
			0,
			0,
			nil,
			"",
			"",
			"reducer not provided for groups, please provide a reducer (e.g. avg, sum) and try again",
			nil,
		},
		{
			"should return error if result is string",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1", Filter: "filter"},
			interface{}("someString"),
			0,
			0,
			0,
			0,
			nil,
			"",
			"",
			"someString",
			nil,
		},
		{
			"should return error on bad filter",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1", Filter: "\""},
			nil,
			0,
			0,
			0,
			0,
			nil,
			"",
			"",
			"filter is not valid",
			nil,
		},
		{
			"should process receiver error",
			queryModel{Command: models.TimeSeriesMRange, Key: "test1"},
			nil,
			0,
			0,
			0,
			0,
			nil,
			"",
			"",
			"error occurred",
			errors.New("error occurred"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := testClient{rcv: tt.rcv, err: tt.err}
			response := queryTsMRange(tt.from, tt.to, tt.qm, &client)
			if tt.expectedError != "" {
				require.EqualError(t, response.Error, tt.expectedError, "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.expectedFrameName, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")
				require.Equal(t, tt.expectedValueFieldName, response.Frames[0].Fields[1].Name, "Invalid field name")

				if tt.valuesToCheckInResponse != nil {
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				}
			}
		})
	}
}

func TestQueryTsGet(t *testing.T) {
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
			"should process receiver without aggregation and legend provided",
			queryModel{Command: models.TimeSeriesGet, Key: "test1"},
			[]string{"1548149180000", "26.199999999999999"},
			2,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			nil,
		},
		{
			"should process receiver error",
			queryModel{Command: models.TimeSeriesRange, Key: "test1"},
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

			client := testClient{rcv: tt.rcv, err: tt.err}
			response := queryTsGet(tt.qm, &client)
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

func TestQueryTsInfo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                          string
		qm                            queryModel
		rcv                           interface{}
		fieldsCount                   int
		rowsPerField                  int
		valueToCheckByLabelInResponse []valueToCheckByLabelInResponse
		err                           error
	}{
		{
			"should process receiver",
			queryModel{Command: models.TimeSeriesQueryIndex, Filter: "test1"},
			map[string]interface{}{
				"totalSamples":    int64(100),
				"memoryUsage":     int64(4184),
				"firstTimestamp":  int64(1548149180),
				"lastTimestamp":   int64(1548149279),
				"retentionTime":   int64(0),
				"chunkCount":      int64(1),
				"byteField":       []byte("byteField"),
				"chunkSize":       int64(256),
				"chunkType":       "compressed",
				"duplicatePolicy": interface{}(nil),
				"labels":          []interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
				"sourceKey":       interface{}(nil),
				"unsupportedType": nil,
				"rules":           []interface{}{},
			},
			11,
			1,
			[]valueToCheckByLabelInResponse{
				{frameIndex: 0, fieldName: "totalSamples", rowIndex: 0, value: int64(100)},
				{frameIndex: 0, fieldName: "chunkType", rowIndex: 0, value: "compressed"},
			},
			nil,
		},
		{
			"should process receiver error",
			queryModel{Command: models.TimeSeriesInfo, Filter: "test1"},
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
			client := testClient{rcv: tt.rcv, err: tt.err}
			response := queryTsInfo(tt.qm, &client)
			if tt.err != nil {
				require.EqualError(t, response.Error, tt.err.Error(), "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.qm.Key, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")

				if tt.valueToCheckByLabelInResponse != nil {
					for _, value := range tt.valueToCheckByLabelInResponse {
						for _, field := range response.Frames[value.frameIndex].Fields {
							if field.Name == value.fieldName {
								require.Equalf(t, value.value, field.At(value.rowIndex), "Invalid value at Frame[%v]:Field[Name:%v]:Row[%v]", value.frameIndex, value.fieldName, value.rowIndex)
							}
						}

					}
				}
			}
		})
	}
}

func TestQueryTsQueryIndex(t *testing.T) {
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
			"should process receiver without aggregation and legend provided",
			queryModel{Command: models.TimeSeriesQueryIndex, Filter: "sensor_id=2"},
			[]string{"temperature:2:32", "temperature:2:33"},
			1,
			2,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: "temperature:2:32"},
				{frameIndex: 0, fieldIndex: 0, rowIndex: 1, value: "temperature:2:33"},
			},
			nil,
		},
		{
			"should process error on bad filter",
			queryModel{Command: models.TimeSeriesQueryIndex, Filter: "\""},
			nil,

			0,
			0,
			nil,
			errors.New("filter is not valid"),
		},
		{
			"should process receiver error",
			queryModel{Command: models.TimeSeriesRange, Filter: "test1"},
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

			client := testClient{rcv: tt.rcv, err: tt.err}
			response := queryTsQueryIndex(tt.qm, &client)
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

func TestQueryTsMGet(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name                    string
		qm                      queryModel
		rcv                     interface{}
		fieldsCount             int
		rowsPerField            int
		valuesToCheckInResponse []valueToCheckInResponse
		expectedFrameName       string
		expectedValueFieldName  string
		expectedError           string
		err                     error
	}{
		{
			"should process receiver and legend provided but with labels",
			queryModel{Command: models.TimeSeriesMGet, Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
					[]interface{}{int64(1548149180000), []byte("26.199999999999999")},
				},
				[]interface{}{
					[]byte("temperature:3:32"),
					[]interface{}{},
					[]interface{}{int64(1548149180000), []byte("26.7")},
				},
			},
			2,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			"temperature:2:32",
			"",
			"",
			nil,
		},
		{
			"should process receiver with labels specified in legend",
			queryModel{Command: models.TimeSeriesMGet, Legend: "area_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
					[]interface{}{int64(1548149210000), []byte("20")},
				},
			},
			2,
			1,
			nil,
			"32",
			"",
			"",
			nil,
		},
		{
			"should process receiver with value field existed in labels",
			queryModel{Command: models.TimeSeriesMGet, Key: "test1", Value: "sensor_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
					[]interface{}{int64(1548149210000), []byte("20")},
				},
			},
			2,
			1,
			nil,
			"temperature:2:32",
			"2",
			"",
			nil,
		},
		{
			"should process receiver with []byte field instead of int",
			queryModel{Command: models.TimeSeriesMGet, Key: "test1", Legend: "area_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
					[]interface{}{[]byte("1548149180000"), []byte("26.199999999999999")},
				},
			},
			2,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			"32",
			"",
			"",
			nil,
		},
		{
			"should process receiver with string field instead of int",
			queryModel{Command: models.TimeSeriesMGet, Key: "test1", Legend: "area_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					[]byte("temperature:2:32"),
					[]interface{}{[]interface{}{[]byte("sensor_id"), []byte("2")}, []interface{}{[]byte("area_id"), []byte("32")}},
					[]interface{}{[]byte("1548149180000"), "26.199999999999999"},
				},
			},
			2,
			1,
			[]valueToCheckInResponse{
				{frameIndex: 0, fieldIndex: 0, rowIndex: 0, value: time.Unix(0, 1548149180000*int64(time.Millisecond))},
				{frameIndex: 0, fieldIndex: 1, rowIndex: 0, value: 26.2},
			},
			"32",
			"",
			"",
			nil,
		},
		{
			"should return error if result is string",
			queryModel{Command: models.TimeSeriesMGet, Filter: "filter"},
			interface{}("someString"),
			0,
			0,
			nil,
			"",
			"",
			"someString",
			nil,
		},
		{
			"should return error on bad filter",
			queryModel{Command: models.TimeSeriesMGet, Key: "test1", Filter: "\""},
			nil,
			0,
			0,
			nil,
			"",
			"",
			"filter is not valid",
			nil,
		},
		{
			"should process receiver error",
			queryModel{Command: models.TimeSeriesMGet, Key: "test1"},
			nil,
			0,
			0,
			nil,
			"",
			"",
			"error occurred",
			errors.New("error occurred"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := testClient{rcv: tt.rcv, err: tt.err}
			response := queryTsMGet(tt.qm, &client)
			if tt.expectedError != "" {
				require.EqualError(t, response.Error, tt.expectedError, "Should set error to response if failed")
				require.Nil(t, response.Frames, "No frames should be created if failed")
			} else {
				require.Equal(t, tt.expectedFrameName, response.Frames[0].Name, "Invalid frame name")
				require.Len(t, response.Frames[0].Fields, tt.fieldsCount, "Invalid number of fields created ")
				require.Equal(t, tt.rowsPerField, response.Frames[0].Fields[0].Len(), "Invalid number of values in field vectors")
				require.Equal(t, tt.expectedValueFieldName, response.Frames[0].Fields[1].Name, "Invalid field name")

				if tt.valuesToCheckInResponse != nil {
					for _, value := range tt.valuesToCheckInResponse {
						require.Equalf(t, value.value, response.Frames[value.frameIndex].Fields[value.fieldIndex].At(value.rowIndex), "Invalid value at Frame[%v]:Field[%v]:Row[%v]", value.frameIndex, value.fieldIndex, value.rowIndex)
					}
				}
			}
		})
	}
}

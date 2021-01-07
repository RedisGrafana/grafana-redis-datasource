package main

import (
	"errors"
	"testing"
	"time"

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
			queryModel{Command: "ts.range", Key: "test1"},
			[]interface{}{
				[]interface{}{1548149180000, "26.199999999999999"},
				[]interface{}{1548149195000, "27.399999999999999"},
				[]interface{}{1548149220000, "24.800000000000001"},
				[]interface{}{1548149215000, "23.199999999999999"},
				[]interface{}{1548149230000, "25.199999999999999"},
				[]interface{}{1548149285000, "28"},
				[]interface{}{1548149150000, "20"},
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
			queryModel{Command: "ts.range", Aggregation: "avg", Bucket: 5000, Key: "test1", Legend: "Legend"},
			[]interface{}{
				[]interface{}{1548149180000, "26.199999999999999"},
				[]interface{}{1548149185000, "27.399999999999999"},
				[]interface{}{1548149190000, "24.800000000000001"},
				[]interface{}{1548149195000, "23.199999999999999"},
				[]interface{}{1548149200000, "25.199999999999999"},
				[]interface{}{1548149205000, "28"},
				[]interface{}{1548149210000, "20"},
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
			queryModel{Command: "ts.range", Bucket: 5000, Key: "test1", Fill: true},
			[]interface{}{
				[]interface{}{1548149180000, "26.199999999999999"},
				[]interface{}{1548149195000, "27.399999999999999"},
				[]interface{}{1548149220000, "24.800000000000001"},
				[]interface{}{1548149215000, "23.199999999999999"},
				[]interface{}{1548149230000, "25.199999999999999"},
				[]interface{}{1548149285000, "28"},
				[]interface{}{1548149150000, "20"},
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
			queryModel{Command: "ts.range", Key: "test1"},
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
			ds := redisDatasource{}
			client := testClient{tt.rcv, tt.err}
			response := ds.queryTsRange(tt.from, tt.to, tt.qm, client)
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
			queryModel{Command: "ts.mrange", Key: "test1", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					"temperature:2:32",
					[][]string{{"sensor_id", "2"}, {"area_id", "32"}},
					[]interface{}{
						[]interface{}{1548149180000, "26.199999999999999"},
						[]interface{}{1548149195000, "27.399999999999999"},
						[]interface{}{1548149220000, "24.800000000000001"},
						[]interface{}{1548149215000, "23.199999999999999"},
						[]interface{}{1548149230000, "25.199999999999999"},
						[]interface{}{1548149285000, "28"},
						[]interface{}{1548149150000, "20"},
					},
				},
				[]interface{}{
					"temperature:3:32",
					[][]string{},
					[]interface{}{
						[]interface{}{1548149180000, "26.7"},
						[]interface{}{1548149195000, "27.8"},
						[]interface{}{1548149220000, "24.4"},
						[]interface{}{1548149215000, "26.199999999999999"},
						[]interface{}{1548149230000, "25.199999999999999"},
						[]interface{}{1548149285000, "27"},
						[]interface{}{1548149150000, "22"},
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
			queryModel{Command: "ts.mrange", Key: "test1", Aggregation: "avg", Legend: "Legend", Bucket: 5000, Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					"temperature:2:32",
					[][]string{},
					[]interface{}{
						[]interface{}{1548149180000, "26.199999999999999"},
						[]interface{}{1548149185000, "27.399999999999999"},
						[]interface{}{1548149190000, "24.800000000000001"},
						[]interface{}{1548149195000, "23.199999999999999"},
						[]interface{}{1548149200000, "25.199999999999999"},
						[]interface{}{1548149205000, "28"},
						[]interface{}{1548149210000, "20"},
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
			queryModel{Command: "ts.mrange", Key: "test1", Legend: "area_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					"temperature:2:32",
					[][]string{{"sensor_id", "2"}, {"area_id", "32"}},
					[]interface{}{
						[]interface{}{1548149210000, "20"},
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
			queryModel{Command: "ts.mrange", Key: "test1", Value: "sensor_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					"temperature:2:32",
					[][]string{{"sensor_id", "2"}, {"area_id", "32"}},
					[]interface{}{
						[]interface{}{1548149210000, "20"},
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
			queryModel{Command: "ts.mrange", Key: "test1", Legend: "area_id", Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					"temperature:2:32",
					[][]string{{"sensor_id", "2"}, {"area_id", "32"}},
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
			queryModel{Command: "ts.mrange", Key: "test1", Fill: true, Bucket: 5000, Filter: "area_id=32 sensor_id!=1"},
			[]interface{}{
				[]interface{}{
					"temperature:2:32",
					[][]string{},
					[]interface{}{
						[]interface{}{1548149180000, "26.199999999999999"},
						[]interface{}{1548149195000, "27.399999999999999"},
						[]interface{}{1548149220000, "24.800000000000001"},
						[]interface{}{1548149215000, "23.199999999999999"},
						[]interface{}{1548149230000, "25.199999999999999"},
						[]interface{}{1548149285000, "28"},
						[]interface{}{1548149150000, "20"},
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
			"should return error on bad filter",
			queryModel{Command: "ts.mrange", Key: "test1", Filter: "\""},
			nil,
			0,
			0,
			0,
			0,
			nil,
			"",
			"",
			"Filter is not valid",
			nil,
		},
		{
			"should process receiver error",
			queryModel{Command: "ts.mrange", Key: "test1"},
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
			ds := redisDatasource{}
			client := testClient{tt.rcv, tt.err}
			response := ds.queryTsMRange(tt.from, tt.to, tt.qm, client)
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
			queryModel{Command: "ts.get", Key: "test1"},
			[]interface{}{1548149180000, "26.199999999999999"},
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
			queryModel{Command: "ts.range", Key: "test1"},
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
			client := testClient{tt.rcv, tt.err}
			response := ds.queryTsGet(tt.qm, client)
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
			queryModel{Command: "ts.queryindex", Filter: "test1"},
			[]interface{}{
				"totalSamples", int64(100),
				"memoryUsage", int64(4184),
				"firstTimestamp", int64(1548149180),
				"lastTimestamp", int64(1548149279),
				"retentionTime", int64(0),
				"chunkCount", int64(1),
				"chunkSize", int64(256),
				"chunkType", "compressed",
				"duplicatePolicy", nil,
				"labels", [][]string{{"sensor_id", "2"}, {"area_id", "32"}},
				"sourceKey", nil,
				"rules", []interface{}{},
			},
			12,
			1,
			[]valueToCheckByLabelInResponse{
				{frameIndex: 0, fieldName: "totalSamples", rowIndex: 0, value: int64(100)},
				{frameIndex: 0, fieldName: "chunkType", rowIndex: 0, value: "compressed"},
			},
			nil,
		},
		{
			"should process receiver error",
			queryModel{Command: "ts.info", Filter: "test1"},
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
			client := testClient{tt.rcv, tt.err}
			response := ds.queryTsInfo(tt.qm, client)
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
			queryModel{Command: "ts.queryindex", Filter: "sensor_id=2"},
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
			queryModel{Command: "ts.queryindex", Filter: "\""},
			nil,

			0,
			0,
			nil,
			errors.New("Filter is not valid"),
		},
		{
			"should process receiver error",
			queryModel{Command: "ts.range", Filter: "test1"},
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
			client := testClient{tt.rcv, tt.err}
			response := ds.queryTsQueryIndex(tt.qm, client)
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

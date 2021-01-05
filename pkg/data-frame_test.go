package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateFrameValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		value    string
		expected interface{}
	}{
		{"3.14", 3.14},
		{"3", float64(3)},
		{"somestring", "somestring"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.value, func(t *testing.T) {
			t.Parallel()
			ds := redisDatasource{}
			frame := ds.createFrameValue("keyName", tt.value)
			field := frame.Fields[0].At(0)
			require.Equal(t, tt.expected, field, "Unexpected conversation")
		})
	}
}

func TestAddFrameFieldsFromArray(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		values      []interface{}
		fieldsCount int
	}{
		{
			"should not parse key of not []byte type, and should not create field",
			[]interface{}{
				[]interface{}{"sensor_id", []byte("2")},
			},
			0,
		},
		{
			"should parse value of type bytes[] with underlying int",
			[]interface{}{
				[]interface{}{[]byte("sensor_id"), []byte("2")},
				[]interface{}{[]byte("area_id"), []byte("32")},
			},
			2,
		},
		{
			"should parse value of type bytes[] with underlying non-int value",
			[]interface{}{
				[]interface{}{[]byte("sensor_id"), []byte("some_string")},
			},
			1,
		},
		{
			"should parse value of type int64",
			[]interface{}{
				[]interface{}{[]byte("sensor_id"), int64(145)},
			},
			1,
		},
		{
			"should not parse value of not bytes[] or int64",
			[]interface{}{
				[]interface{}{[]byte("sensor_id"), float32(3.14)},
			},
			0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ds := redisDatasource{}
			frame := data.NewFrame("name")
			frame = ds.addFrameFieldsFromArray(tt.values, frame)
			require.Len(t, frame.Fields, tt.fieldsCount, "Invalid number of fields created in Frame")
		})
	}
}

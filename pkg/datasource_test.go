package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

/**
 * Client Config
 */
func TestCreateRedisClientConfig(t *testing.T) {
	tests := []struct {
		name      string
		dataModel dataModel
		settings  backend.DataSourceInstanceSettings
		expected  redisClientConfiguration
		err       string
	}{
		{
			"should parse settings without secureData",
			dataModel{
				PoolSize:       1,
				Timeout:        1000,
				PingInterval:   100,
				PipelineWindow: 100,
				TLSAuth:        true,
				TLSSkipVerify:  true,
				Client:         "socket",
				SentinelName:   "",
				ACL:            true,
				User:           "",
				SentinelACL:    false,
				SentinelUser:   "",
			},
			backend.DataSourceInstanceSettings{},
			redisClientConfiguration{
				Timeout:        1000,
				PoolSize:       1,
				PingInterval:   100,
				PipelineWindow: 100,
				User:           "",
				ACL:            true,
				SentinelACL:    false,
				SentinelUser:   "",
				TLSAuth:        true,
				TLSSkipVerify:  true,
				Client:         "socket",
			},
			"",
		},
		{
			"should parse settings with secureData and default settings",
			dataModel{
				User: "username",
			},
			backend.DataSourceInstanceSettings{
				URL: "localhost:6379",
				DecryptedSecureJSONData: map[string]string{
					"password":         "1234",
					"sentinelPassword": "123456",
					"tlsCACert":        "BEGIN CERTIFICATE",
					"tlsClientCert":    "BEGIN CERTIFICATE",
					"tlsClientKey":     "BEGIN RSA PRIVATE KEY",
				},
			},
			redisClientConfiguration{
				URL:              "localhost:6379",
				Timeout:          10,
				PoolSize:         5,
				PingInterval:     0,
				PipelineWindow:   0,
				User:             "username",
				Password:         "1234",
				SentinelPassword: "123456",
				TLSCACert:        "BEGIN CERTIFICATE",
				TLSClientCert:    "BEGIN CERTIFICATE",
				TLSClientKey:     "BEGIN RSA PRIVATE KEY",
			},
			"",
		},
		{
			"should parse sentinel settings with secureData and default settings",
			dataModel{
				Client:       "sentinel",
				User:         "username",
				SentinelUser: "admin",
			},
			backend.DataSourceInstanceSettings{
				URL: "localhost:6379",
				DecryptedSecureJSONData: map[string]string{
					"password":         "1234",
					"sentinelPassword": "123456",
					"tlsCACert":        "BEGIN CERTIFICATE",
					"tlsClientCert":    "BEGIN CERTIFICATE",
					"tlsClientKey":     "BEGIN RSA PRIVATE KEY",
				},
			},
			redisClientConfiguration{
				URL:              "localhost:6379",
				Timeout:          10,
				PoolSize:         5,
				PingInterval:     0,
				PipelineWindow:   0,
				Client:           "sentinel",
				User:             "username",
				SentinelUser:     "admin",
				Password:         "1234",
				SentinelPassword: "123456",
				TLSCACert:        "BEGIN CERTIFICATE",
				TLSClientCert:    "BEGIN CERTIFICATE",
				TLSClientKey:     "BEGIN RSA PRIVATE KEY",
			},
			"",
		},
	}

	// Run Tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			raw, _ := json.Marshal(&tt.dataModel)
			tt.settings.JSONData = raw

			// Config
			config, err := createRedisClientConfig(tt.settings)
			if tt.err != "" {
				require.EqualError(t, err, tt.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, config)
			}
			fmt.Println(config)
		})
	}

	// Error
	t.Run("handle marshalling error", func(t *testing.T) {
		t.Parallel()
		_, err := createRedisClientConfig(backend.DataSourceInstanceSettings{})
		require.EqualError(t, err, "unexpected end of JSON input")
	})
}

/**
 * Dispose
 */
func TestDispose(t *testing.T) {
	// Client
	client := &testClient{}
	client.On("Close").Return(nil)

	// Instance
	is := instanceSettings{client}
	is.Dispose()
	client.AssertNumberOfCalls(t, "Close", 1)
}

/**
 * Get Instance
 */
func TestGetInstance(t *testing.T) {
	// Data Source
	client := &testClient{}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}
	ctx := context.Background()

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)
	actualClient, err := ds.getInstance(ctx, backend.PluginContext{})
	require.Equal(t, client, actualClient)
	require.NoError(t, err)
}

/**
 * Get Instance Error
 */
func TestGetInstanceError(t *testing.T) {
	// Data Source
	client := &testClient{}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}
	ctx := context.Background()

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, errors.New("some_err"))
	_, err := ds.getInstance(ctx, backend.PluginContext{})
	require.EqualError(t, err, "some_err")
}

/**
 * Query Data
 */
func TestQueryData(t *testing.T) {
	// Data Source
	client := &testClient{rcv: "3.14", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)

	// HGET
	dm := queryModel{Command: models.HGet, Key: "test1", Field: "key1"}
	marshaled, _ := json.Marshal(dm)

	// Response
	response, err := ds.QueryData(context.TODO(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Headers:       nil,
		Queries: []backend.DataQuery{
			{
				RefID:         "",
				QueryType:     "",
				MaxDataPoints: 100,
				Interval:      10,
				TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
				JSON:          marshaled,
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, response.Responses, 1)
}

/**
 * Query Data with Error
 */
func TestQueryDataWithError(t *testing.T) {
	// Client
	client := &testClient{rcv: "3.14", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, errors.New("some_err"))

	// HGET
	dm := queryModel{Command: models.HGet, Key: "test1", Field: "key1"}
	marshaled, _ := json.Marshal(dm)

	// Query
	_, err := ds.QueryData(context.TODO(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Headers:       nil,
		Queries: []backend.DataQuery{
			{
				RefID:         "",
				QueryType:     "",
				MaxDataPoints: 100,
				Interval:      10,
				TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
				JSON:          marshaled,
			},
		},
	})
	require.EqualError(t, err, "some_err")
}

/**
 * Query Data with Marshalling Error
 */
func TestQueryDataWithMarshallingError(t *testing.T) {
	// Client
	client := &testClient{rcv: "3.14", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)

	// Query
	resp, _ := ds.QueryData(context.TODO(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Headers:       nil,
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     "",
				MaxDataPoints: 100,
				Interval:      10,
				TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
				JSON:          []byte{31, 17, 45},
			},
		},
	})
	require.EqualError(t, resp.Responses["A"].Error, "invalid character '\\x1f' looking for beginning of value", "Should return marshalling error")
}

/**
 * Check Health
 */
func TestCheckHealth(t *testing.T) {
	// Client
	client := &testClient{rcv: "PONG", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)

	// Result
	result, err := ds.CheckHealth(context.TODO(), &backend.CheckHealthRequest{})
	require.NoError(t, err)
	require.Equal(t, result.Status, backend.HealthStatusOk)
}

/**
 * Check Health with Error
 */
func TestCheckHealthWithErrorFromIm(t *testing.T) {
	// Client
	client := &testClient{rcv: "PONG", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, errors.New("some_err"))

	// Result
	result, err := ds.CheckHealth(context.TODO(), &backend.CheckHealthRequest{})
	require.NoError(t, err)
	require.Equal(t, result.Status, backend.HealthStatusError)
}

/**
 * Check Health with Client Error
 */
func TestCheckHealthWithErrorFromClient(t *testing.T) {
	// Client
	client := &testClient{rcv: "PONG", err: errors.New("some_err")}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)

	// Result
	result, err := ds.CheckHealth(context.TODO(), &backend.CheckHealthRequest{})
	require.NoError(t, err)
	require.Equal(t, result.Status, backend.HealthStatusError)
}

/**
 * Test Streaming TimeSeries
 */
func TestStreamingTimeSeries(t *testing.T) {
	// Data Source
	client := &testClient{rcv: "# Server\nredis_version:6.0.1\nredis_git_sha1:00000000\nredis_git_dirty:0\nredis_build_id:e02d1d807e41d65\nredis_mode:standalone\nos:Linux 5.10.25-linuxkit x86_64", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)

	// INFO
	dm := queryModel{Command: models.Info, Section: "server", Streaming: true}
	marshaled, _ := json.Marshal(dm)

	// Response
	response, err := ds.QueryData(context.TODO(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Headers:       nil,
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     "",
				MaxDataPoints: 100,
				Interval:      10,
				TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
				JSON:          marshaled,
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, response.Responses, 1)
	require.Len(t, response.Responses["A"].Frames, 1)

	frame := response.Responses["A"].Frames[0]
	require.Len(t, frame.Fields, 7)
	require.Equal(t, 1, frame.Fields[0].Len())
	require.Equal(t, "#time", frame.Fields[0].Name)
	require.LessOrEqual(t, time.Now().Unix(), frame.Fields[0].At(0).(time.Time).Unix())
	require.Equal(t, "redis_version", frame.Fields[1].Name)
	require.LessOrEqual(t, "6.0.1", frame.Fields[1].At(0))
}

/**
 * Test Streaming TimeSeries with Field
 */
func TestStreamingTimeSeriesWithField(t *testing.T) {
	// Data Source
	client := &testClient{rcv: "3.14", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)

	// HGET
	dm := queryModel{Command: models.HGet, Key: "test1", Field: "key1", Streaming: true}
	marshaled, _ := json.Marshal(dm)

	// Response
	response, err := ds.QueryData(context.TODO(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Headers:       nil,
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     "",
				MaxDataPoints: 100,
				Interval:      10,
				TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
				JSON:          marshaled,
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, response.Responses, 1)
	require.Len(t, response.Responses["A"].Frames, 1)

	frame := response.Responses["A"].Frames[0]
	require.Len(t, frame.Fields, 2)
	require.Equal(t, 1, frame.Fields[0].Len())
	require.Equal(t, "#time", frame.Fields[0].Name)
	require.LessOrEqual(t, time.Now().Unix(), frame.Fields[0].At(0).(time.Time).Unix())
	require.Equal(t, "key1", frame.Fields[1].Name)
	require.Equal(t, 3.14, frame.Fields[1].At(0))
}

/**
 * Test Streaming TimeSeries with Error Field
 */
func TestStreamingTimeSeriesWithErrorField(t *testing.T) {
	// Data Source
	client := &testClient{rcv: "3.14", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)

	// INFO
	dm := queryModel{Command: models.Info, Field: "\"", Streaming: true}
	marshaled, _ := json.Marshal(dm)

	// Response
	response, err := ds.QueryData(context.TODO(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Headers:       nil,
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     "",
				MaxDataPoints: 100,
				Interval:      10,
				TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
				JSON:          marshaled,
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, response.Responses, 1)
	require.Len(t, response.Responses["A"].Frames, 1)

	error := response.Responses["A"].Error
	require.EqualError(t, error, "field is not valid")
}

/**
 * Test Streaming TimeSeries with wrong Field
 */
func TestStreamingTimeSeriesWithWrongField(t *testing.T) {
	// Data Source
	client := &testClient{rcv: "# Server\nredis_version:6.0.1\nredis_git_sha1:00000000\nredis_git_dirty:0\nredis_build_id:e02d1d807e41d65\nredis_mode:standalone\nos:Linux 5.10.25-linuxkit x86_64", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}

	// Instance
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)

	// INFO
	dm := queryModel{Command: models.Info, Section: "server", Field: "key1", Streaming: true}
	marshaled, _ := json.Marshal(dm)

	// Response
	response, err := ds.QueryData(context.TODO(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Headers:       nil,
		Queries: []backend.DataQuery{
			{
				RefID:         "A",
				QueryType:     "",
				MaxDataPoints: 100,
				Interval:      10,
				TimeRange:     backend.TimeRange{From: time.Now(), To: time.Now()},
				JSON:          marshaled,
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, response.Responses, 1)
	require.Len(t, response.Responses["A"].Frames, 1)

	frame := response.Responses["A"].Frames[0]
	require.Len(t, frame.Fields, 1)
	require.Equal(t, 1, frame.Fields[0].Len())
	require.Equal(t, "#time", frame.Fields[0].Name)
	require.LessOrEqual(t, time.Now().Unix(), frame.Fields[0].At(0).(time.Time).Unix())
}

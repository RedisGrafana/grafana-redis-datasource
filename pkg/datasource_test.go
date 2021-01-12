package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

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
			},
			backend.DataSourceInstanceSettings{},
			redisClientConfiguration{
				Timeout:        1000,
				PoolSize:       1,
				PingInterval:   100,
				PipelineWindow: 100,
				ACL:            true,
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
					"password":      "1234",
					"tlsCACert":     "BEGIN CERTIFICATE",
					"tlsClientCert": "BEGIN CERTIFICATE",
				},
			},
			redisClientConfiguration{
				Url:            "localhost:6379",
				Timeout:        10,
				PoolSize:       5,
				PingInterval:   0,
				PipelineWindow: 0,
				User:           "username",
				Password:       "1234",
				TlsCACert:      "BEGIN CERTIFICATE",
				TlsClientCert:  "BEGIN CERTIFICATE",
			},
			"",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			raw, _ := json.Marshal(&tt.dataModel)
			tt.settings.JSONData = raw
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

	t.Run("handle marshalling error", func(t *testing.T) {
		t.Parallel()
		_, err := createRedisClientConfig(backend.DataSourceInstanceSettings{})
		require.EqualError(t, err, "unexpected end of JSON input")
	})
}

func TestDispose(t *testing.T) {
	client := &testClient{}
	client.On("Close").Return(nil)
	is := instanceSettings{client}
	is.Dispose()
	client.AssertNumberOfCalls(t, "Close", 1)
}

func TestGetInstance(t *testing.T) {
	client := &testClient{}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)
	actualClient, err := ds.getInstance(backend.PluginContext{})
	require.Equal(t, client, actualClient)
	require.NoError(t, err)
}
func TestGetInstanceError(t *testing.T) {
	client := &testClient{}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, errors.New("some_err"))
	_, err := ds.getInstance(backend.PluginContext{})
	require.EqualError(t, err, "some_err")
}

func TestQueryData(t *testing.T) {
	client := &testClient{rcv: "3.14", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)
	dm := queryModel{Command: "hget", Key: "test1", Field: "key1"}
	marshaled, _ := json.Marshal(dm)
	response, err := ds.QueryData(context.TODO(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Headers:       nil,
		Queries: []backend.DataQuery{
			backend.DataQuery{
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

func TestQueryDataWithError(t *testing.T) {
	client := &testClient{rcv: "3.14", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, errors.New("some_err"))
	dm := queryModel{Command: "hget", Key: "test1", Field: "key1"}
	marshaled, _ := json.Marshal(dm)
	_, err := ds.QueryData(context.TODO(), &backend.QueryDataRequest{
		PluginContext: backend.PluginContext{},
		Headers:       nil,
		Queries: []backend.DataQuery{
			backend.DataQuery{
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

func TestCheckHealth(t *testing.T) {
	client := &testClient{rcv: "PONG", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)
	result, err := ds.CheckHealth(context.TODO(), &backend.CheckHealthRequest{})
	require.NoError(t, err)
	require.Equal(t, result.Status, backend.HealthStatusOk)
}

func TestCheckHealthWithErrorFromIm(t *testing.T) {
	client := &testClient{rcv: "PONG", err: nil}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, errors.New("some_err"))
	result, err := ds.CheckHealth(context.TODO(), &backend.CheckHealthRequest{})
	require.NoError(t, err)
	require.Equal(t, result.Status, backend.HealthStatusError)
}

func TestCheckHealthWithErrorFromClient(t *testing.T) {
	client := &testClient{rcv: "PONG", err: errors.New("some_err")}
	im := fakeInstanceManager{}
	ds := redisDatasource{&im}
	is := instanceSettings{client}
	im.On("Get", mock.Anything).Return(&is, nil)
	result, err := ds.CheckHealth(context.TODO(), &backend.CheckHealthRequest{})
	require.NoError(t, err)
	require.Equal(t, result.Status, backend.HealthStatusError)
}

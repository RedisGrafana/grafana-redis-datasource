package main

import (
	"encoding/json"
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"
	"testing"
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

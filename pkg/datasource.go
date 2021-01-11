package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/mediocregopher/radix/v3"
)

/**
 * The function is called when the instance is created for the first time or when a datasource configuration changed.
 */
func newDatasource() datasource.ServeOpts {
	im := datasource.NewInstanceManager(newDataSourceInstance)

	ds := &redisDatasource{
		im: im,
	}

	// Returns datasource.ServeOpts
	return datasource.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	}
}

/**
 * QueryData handles multiple queries and returns multiple responses.
 * req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
 * The QueryDataResponse contains a map of RefID to the response for each query, and each response contains Frames ([]*Frame).
 */
func (ds *redisDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Debug("QueryData", "request", req)

	// Get Instance
	client, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return nil, err
	}

	// Create response struct
	response := backend.NewQueryDataResponse()

	// Loop over queries and execute them individually
	for _, q := range req.Queries {
		res := query(ctx, q, client)

		// save the response in a hashmap based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

/**
 * CheckHealth handles health checks sent from Grafana to the plugin
 *
 * @see https://redis.io/commands/ping
 */
func (ds *redisDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status backend.HealthStatus
	message := "Data Source health is yet to become known."

	// Get Instance
	client, err := ds.getInstance(req.PluginContext)

	if err != nil {
		status = backend.HealthStatusError
		message = fmt.Sprintf("getInstance error: %s", err.Error())
	} else {
		err = client.Do(radix.Cmd(&message, "PING"))

		// Check errors
		if err != nil {
			status = backend.HealthStatusError
			message = fmt.Sprintf("PING command failed: %s", err.Error())
		} else {
			status = backend.HealthStatusOk
			message = "Data Source is working as expected."
		}
	}

	// Return Health result
	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

/**
 * Return Instance
 */
func (ds *redisDatasource) getInstance(ctx backend.PluginContext) (ClientInterface, error) {
	s, err := ds.im.Get(ctx)

	if err != nil {
		return nil, err
	}

	// Return client
	return s.(*instanceSettings).client, nil
}

/**
 * New Datasource Instance
 *
 * @see https://github.com/mediocregopher/radix
 */
func newDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var jsonData dataModel

	// Unmarshal Configuration
	var dataError = json.Unmarshal(setting.JSONData, &jsonData)
	if dataError != nil {
		log.DefaultLogger.Error("JSONData", "Error", dataError)
		return nil, dataError
	}

	// Debug
	log.DefaultLogger.Debug("JSONData", "Values", jsonData)

	// Pool size
	poolSize := 5
	if jsonData.PoolSize > 0 {
		poolSize = jsonData.PoolSize
	}

	// Connect, Read and Write Timeout
	timeout := 10
	if jsonData.Timeout > 0 {
		timeout = jsonData.Timeout
	}

	// Ping Interval, disabled by default
	pingInterval := 0
	if jsonData.PingInterval > 0 {
		pingInterval = jsonData.PingInterval
	}

	// Pipeline Window, disabled by default
	pipelineWindow := 0
	if jsonData.PipelineWindow > 0 {
		pipelineWindow = jsonData.PipelineWindow
	}

	// Secured Data
	var secureData = setting.DecryptedSecureJSONData

	// Set up connection
	connFunc := func(network, addr string) (radix.Conn, error) {
		opts := []radix.DialOpt{radix.DialTimeout(time.Duration(timeout) * time.Second)}

		// Authentication
		if secureData != nil && secureData["password"] != "" {
			// If ACL enabled
			if jsonData.ACL {
				opts = append(opts, radix.DialAuthUser(jsonData.User, secureData["password"]))
			} else {
				opts = append(opts, radix.DialAuthPass(secureData["password"]))
			}
		}

		// TLS Authentication
		if jsonData.TLSAuth {
			// TLS Config
			tlsConfig := &tls.Config{
				InsecureSkipVerify: jsonData.TLSSkipVerify,
			}

			// Certification Authority
			if secureData["tlsCACert"] != "" {
				caPool := x509.NewCertPool()
				ok := caPool.AppendCertsFromPEM([]byte(secureData["tlsCACert"]))
				if ok {
					tlsConfig.RootCAs = caPool
				}
			}

			// Certificate and Key
			if secureData["tlsClientCert"] != "" {
				cert, err := tls.X509KeyPair([]byte(secureData["tlsClientCert"]), []byte(secureData["tlsClientKey"]))
				if err == nil {
					tlsConfig.Certificates = []tls.Certificate{cert}
				} else {
					log.DefaultLogger.Error("X509KeyPair", "Error", err)
					return nil, err
				}
			}

			// Add TLS Config
			opts = append(opts, radix.DialUseTLS(tlsConfig))
		}

		return radix.Dial(network, addr, opts...)
	}

	// Pool with specified Ping Interval, Pipeline Window and Timeout
	poolFunc := func(network, addr string) (radix.Client, error) {
		return radix.NewPool(network, addr, poolSize, radix.PoolConnFunc(connFunc),
			radix.PoolPingInterval(time.Duration(pingInterval)*time.Second/time.Duration(poolSize+1)),
			radix.PoolPipelineWindow(time.Duration(pipelineWindow)*time.Microsecond, 0))
	}

	var client ClientInterface
	var err error

	// Client Type
	switch jsonData.Client {
	case "cluster":
		client, err = radix.NewCluster(strings.Split(setting.URL, ","), radix.ClusterPoolFunc(poolFunc))
	case "sentinel":
		client, err = radix.NewSentinel(jsonData.SentinelName, strings.Split(setting.URL, ","), radix.SentinelConnFunc(connFunc),
			radix.SentinelPoolFunc(poolFunc))
	case "socket":
		client, err = poolFunc("unix", setting.URL)
	default:
		client, err = poolFunc("tcp", setting.URL)
	}

	if err != nil {
		return nil, err
	}

	return &instanceSettings{
		client,
	}, nil
}

/**
 * Called before creating a new instance to close Redis connection pool
 */
func (s *instanceSettings) Dispose() {
	s.client.Close()
}

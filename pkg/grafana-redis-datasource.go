package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
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
		res := ds.query(ctx, q, client)

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
	var status = backend.HealthStatusUnknown
	var message = "Data source health is yet to become known"

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
		} else {
			status = backend.HealthStatusOk
			message = "Data source working as expected"
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
func (ds *redisDatasource) getInstance(ctx backend.PluginContext) (*radix.Pool, error) {
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

	// Default Pool size
	poolSize := 5

	// Default Connect, Read and Write Timeout
	timeout := 10

	// Default Ping Interval disabled
	pingInterval := 0

	// Default Pipeline Window disabled
	pipelineWindow := 0

	if dataError != nil {
		log.DefaultLogger.Error("JSONData", "Error", dataError)
	} else {
		log.DefaultLogger.Debug("JSONData", "Values", jsonData)

		// Set values
		poolSize = jsonData.PoolSize
		timeout = jsonData.Timeout
		pingInterval = jsonData.PingInterval
		pipelineWindow = jsonData.PipelineWindow
	}

	// Secured Data
	var secureData = setting.DecryptedSecureJSONData

	// Set up connection
	connFunc := func(network, addr string) (radix.Conn, error) {
		opts := []radix.DialOpt{radix.DialTimeout(time.Duration(timeout) * time.Second)}

		// Add Password
		if secureData != nil && secureData["password"] != "" {
			opts = append(opts, radix.DialAuthPass(secureData["password"]))
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
			cert, err := tls.X509KeyPair([]byte(secureData["tlsClientCert"]), []byte(secureData["tlsClientKey"]))
			if err == nil {
				tlsConfig.Certificates = []tls.Certificate{cert}
			}

			// Add TLS Config
			opts = append(opts, radix.DialUseTLS(tlsConfig))
		}

		return radix.Dial(network, addr, opts...)
	}

	// Pool with specified Ping Interval, Pipeline Window and Timeout
	pool, err := radix.NewPool("tcp", setting.URL, poolSize, radix.PoolConnFunc(connFunc),
		radix.PoolPingInterval(time.Duration(pingInterval)*time.Second/time.Duration(poolSize+1)),
		radix.PoolPipelineWindow(time.Duration(pipelineWindow)*time.Microsecond, 0))

	if err != nil {
		return nil, err
	}

	return &instanceSettings{
		client: pool,
	}, nil
}

/**
 * Called before creating a new instance to close Redis connection pool
 */
func (s *instanceSettings) Dispose() {
	s.client.Close()
}

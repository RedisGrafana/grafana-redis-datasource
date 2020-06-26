package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/mediocregopher/radix/v3"
)

// newDatasource returns datasource.ServeOpts.
func newDatasource() datasource.ServeOpts {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.
	im := datasource.NewInstanceManager(newDataSourceInstance)
	ds := &redisDatasource{
		im: im,
	}

	return datasource.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	}
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (ds *redisDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Debug("QueryData", "request", req)

	client, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return nil, err
	}

	// Create response struct
	response := backend.NewQueryDataResponse()

	// Loop over queries and execute them individually
	for _, q := range req.Queries {
		res := ds.query(ctx, q, client)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (ds *redisDatasource) query(ctx context.Context, query backend.DataQuery, client *radix.Pool) backend.DataResponse {
	var qm queryModel

	// Unmarshal the json into our queryModel
	err := json.Unmarshal(query.JSON, &qm)
	log.DefaultLogger.Debug("QueryData", "JSON", query.JSON)

	// Error
	if err != nil {
		response := backend.DataResponse{}
		response.Error = err
		return response
	}

	// From and To
	from := query.TimeRange.From.UnixNano() / 1000000
	to := query.TimeRange.To.UnixNano() / 1000000

	// Handle Panic from any command
	defer func() {
		if err := recover(); err != nil {
			log.DefaultLogger.Error("PANIC", "occurred", err)
		}
	}()

	// Commands
	switch qm.Command {
	case "tsrange":
		return ds.queryTsRange(from, to, qm, client)
	case "tsmrange":
		return ds.queryTsMRange(from, to, qm, client)
	case "hgetall":
		return ds.queryHGetAll(qm, client)
	case "smembers":
		return ds.querySMembers(qm, client)
	case "hget":
		return ds.queryHGet(qm, client)
	default:
		response := backend.DataResponse{}
		response.Error = fmt.Errorf("Unknown command")
		return response
	}
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (ds *redisDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusUnknown
	var message = "Data source health is yet to become known"

	r, err := ds.getInstance(req.PluginContext)
	if err != nil {
		status = backend.HealthStatusError
		message = fmt.Sprintf("getInstance error: %s", err.Error())
	} else {
		err = r.Do(radix.Cmd(&message, "PING"))
		if err != nil {
			status = backend.HealthStatusError
		} else {
			status = backend.HealthStatusOk
			message = "Data source working as expected"
		}
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (ds *redisDatasource) getInstance(ctx backend.PluginContext) (*radix.Pool, error) {
	s, err := ds.im.Get(ctx)
	if err != nil {
		return nil, err
	}

	return s.(*instanceSettings).client, nil
}

// NewClient creates a new Client with or without authentication.
func newClient(setting backend.DataSourceInstanceSettings) (*radix.Pool, error) {
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

	// Check if password specified in Secured Data
	password := ""
	if secureData != nil && secureData["password"] != "" {
		password = secureData["password"]
	}

	// Set up a connection which is authenticated and has a 10 seconds timeout on all operations
	connFunc := func(network, addr string) (radix.Conn, error) {
		return radix.Dial(network, addr,
			radix.DialTimeout(time.Duration(timeout)*time.Second),
			radix.DialAuthPass(password),
		)
	}

	// Return Pool with specified Ping Interval, Pipeline Window and Timeout
	return radix.NewPool("tcp", setting.URL, poolSize, radix.PoolConnFunc(connFunc),
		radix.PoolPingInterval(time.Duration(pingInterval)*time.Second/time.Duration(poolSize+1)),
		radix.PoolPipelineWindow(time.Duration(pipelineWindow)*time.Microsecond, 0))
}

func newDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	pool, err := newClient(setting)
	if err != nil {
		return nil, err
	}
	return &instanceSettings{
		client: pool,
	}, nil
}

func (s *instanceSettings) Dispose() {
	// Called before creating a a new instance to allow plugin authors
	// to cleanup.
	// s.Client.Close()
}

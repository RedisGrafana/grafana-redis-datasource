package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/mediocregopher/radix/v3"
	"github.com/mediocregopher/radix/v3/resp/resp2"
)

type instanceSettings struct {
	client *radix.Pool
}

// newDatasource returns datasource.ServeOpts.
func newDatasource() datasource.ServeOpts {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.
	im := datasource.NewInstanceManager(newDataSourceInstance)
	ds := &RedisTimeSeriesDatasource{
		im: im,
	}

	return datasource.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	}
}

// RedisTimeSeriesDatasource is an example datasource used to scaffold
// new datasource plugins with an backend.
type RedisTimeSeriesDatasource struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirements
	// but a best practice that we recommend that you follow.
	im instancemgmt.InstanceManager
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (ds *RedisTimeSeriesDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData", "request", req)

	r, err := ds.getInstance(req.PluginContext)
	if err != nil {
		return nil, err
	}

	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := ds.query(ctx, q, r)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}

	return response, nil
}

type queryModel struct {
	QueryText string `json:"queryText"`
}

func (ds *RedisTimeSeriesDatasource) query(ctx context.Context, query backend.DataQuery, r *radix.Pool) backend.DataResponse {
	// Unmarshal the json into our queryModel
	var qm queryModel

	response := backend.DataResponse{}

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	var res [][]string
	err := r.Do(radix.FlatCmd(&res, "TS.RANGE", qm.QueryText, query.TimeRange.From.UnixNano()/1000000, query.TimeRange.To.UnixNano()/1000000))
	if err != nil {
		var redisErr resp2.Error
		if errors.As(err, &redisErr) {
			response.Error = redisErr.E
		} else {
			response.Error = err
		}
		return response
	}

	// create data frame response
	frame := data.NewFrame(qm.QueryText)

	// add the time dimension
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{}),
	)

	// add values
	frame.Fields = append(frame.Fields,
		data.NewField("values", nil, []float64{}),
	)

	// add rows
	for _, row := range res {
		t, _ := strconv.ParseInt(row[0], 10, 64)
		ts := time.Unix(t/1000, 0)
		v, _ := strconv.ParseFloat(row[1], 64)
		frame.AppendRow(ts, v)
	}

	// add the frames to the response
	response.Frames = append(response.Frames, frame)

	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (ds *RedisTimeSeriesDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
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
			// message = fmt.Sprintf("getInstance PING error: %s", err.Error())
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

func (ds *RedisTimeSeriesDatasource) getInstance(ctx backend.PluginContext) (*radix.Pool, error) {
	s, err := ds.im.Get(ctx)
	if err != nil {
		return nil, err
	}

	return s.(*instanceSettings).client, nil
}

func newDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	pool, err := radix.NewPool("tcp", setting.URL, 1)
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

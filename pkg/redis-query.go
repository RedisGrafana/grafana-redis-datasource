package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/mediocregopher/radix/v3"
	"github.com/mediocregopher/radix/v3/resp/resp2"
)

/**
 * Query commands
 */
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
			log.DefaultLogger.Error("PANIC", "command", err)
		}
	}()

	/**
	 * Custom Command using Query
	 */
	if qm.Query != "" {
		return ds.queryCustomCommand(qm, client)
	}

	/**
	 * Commands switch
	 */
	switch qm.Command {
	case "ts.get":
		return ds.queryTsGet(qm, client)
	case "ts.info":
		return ds.queryTsInfo(qm, client)
	case "ts.queryindex":
		return ds.queryTsQueryIndex(qm, client)
	case "ts.range":
		return ds.queryTsRange(from, to, qm, client)
	case "ts.mrange":
		return ds.queryTsMRange(from, to, qm, client)
	case "hgetall":
		return ds.queryHGetAll(qm, client)
	case "smembers", "hkeys":
		return ds.querySMembers(qm, client)
	case "hget":
		return ds.queryHGet(qm, client)
	case "info":
		return ds.queryInfo(qm, client)
	case "clientList":
		return ds.queryClientList(qm, client)
	case "slowlogGet":
		return ds.querySlowlogGet(qm, client)
	case "type", "get", "ttl", "hlen", "xlen", "llen", "scard":
		return ds.queryKeyCommand(qm, client)
	case "xinfoStream":
		return ds.queryXInfoStream(qm, client)
	default:
		response := backend.DataResponse{}
		response.Error = fmt.Errorf("Unknown command")
		return response
	}
}

/**
 * Error Handler
 */
func (ds *redisDatasource) errorHandler(response backend.DataResponse, err error) backend.DataResponse {
	var redisErr resp2.Error

	// Check for RESP2 Error
	if errors.As(err, &redisErr) {
		response.Error = redisErr.E
	} else {
		response.Error = err
	}

	// Return Response
	return response
}

/**
 * Commands with one key parameter and return value
 *
 * @see https://redis.io/commands/type
 * @see https://redis.io/commands/ttl
 * @see https://redis.io/commands/hlen
 */
func (ds *redisDatasource) queryKeyCommand(qm queryModel, client *radix.Pool) backend.DataResponse {
	response := backend.DataResponse{}

	// Execute command
	var value string
	err := client.Do(radix.Cmd(&value, qm.Command, qm.Key))

	// Check error
	if err != nil {
		return ds.errorHandler(response, err)
	}

	// Add the frames to the response
	response.Frames = append(response.Frames, ds.createFrameValue(qm.Key, value))

	// Return Response
	return response
}

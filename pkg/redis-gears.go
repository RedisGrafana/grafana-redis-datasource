package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
)

/**
 * RG.PYSTATS
 *
 * Returns memory usage statistics from the Python interpreter
 * @see https://oss.redislabs.com/redisgears/commands.html#rgpystats
 */
func queryRgPystats(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Using radix marshaling of key-value arrays to structs
	var stats models.PyStats

	// Run command
	err := client.RunCmd(&stats, models.GearsPyStats)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Command)
	response.Frames = append(response.Frames, frame)

	// New Fields
	frame.Fields = append(frame.Fields, data.NewField("TotalAllocated", nil, []int64{stats.TotalAllocated}))
	frame.Fields = append(frame.Fields, data.NewField("PeakAllocated", nil, []int64{stats.PeakAllocated}))
	frame.Fields = append(frame.Fields, data.NewField("CurrAllocated", nil, []int64{stats.CurrAllocated}))

	return response
}

/**
 * RG.DUMPREGISTRATIONS
 *
 * Returns the list of function registrations
 * @see https://oss.redislabs.com/redisgears/commands.html#rgdumpregistrations
 */
func queryRgDumpregistrations(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Using radix marshaling of key-value arrays to structs
	var registrations []models.DumpRegistrations

	// Run command
	err := client.RunCmd(&registrations, models.GearsDumpRegistrations)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame for all data except of RegistrationData.args
	frame := data.NewFrame(qm.Command)
	response.Frames = append(response.Frames, frame)

	// New Fields
	frame.Fields = append(frame.Fields, data.NewField("id", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("reader", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("desc", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("PD", nil, []string{}))

	frame.Fields = append(frame.Fields, data.NewField("mode", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("numTriggered", nil, []int64{}))
	frame.Fields = append(frame.Fields, data.NewField("numSuccess", nil, []int64{}))
	frame.Fields = append(frame.Fields, data.NewField("numFailures", nil, []int64{}))
	frame.Fields = append(frame.Fields, data.NewField("numAborted", nil, []int64{}))
	frame.Fields = append(frame.Fields, data.NewField("lastError", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("args", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("status", nil, []string{}))

	// Registrations
	for _, registration := range registrations {
		// Merging args to string like "key"="value"\n
		args := new(bytes.Buffer)
		for key, value := range registration.RegistrationData.Args {
			fmt.Fprintf(args, "\"%s\"=\"%s\"\n", key, value)
		}

		frame.AppendRow(registration.ID, registration.Reader, registration.Desc, registration.PD, registration.RegistrationData.Mode,
			registration.RegistrationData.NumTriggered, registration.RegistrationData.NumSuccess, registration.RegistrationData.NumFailures,
			registration.RegistrationData.NumAborted, registration.RegistrationData.LastError, args.String(), registration.RegistrationData.Status)
	}

	return response
}

/**
 * RG.PYEXECUTE "<function>" [UNBLOCKING] [REQUIREMENTS "<dep> ..."]
 *
 * Executes a Python function
 * @see https://oss.redislabs.com/redisgears/commands.html#rgpyexecute
 */
func queryRgPyexecute(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	var result interface{}

	// Check and create list of optional parameters
	var args []interface{}
	if qm.Unblocking {
		args = append(args, "UNBLOCKING")
	}

	if qm.Requirements != "" {
		args = append(args, "REQUIREMENTS", qm.Requirements)
	}

	// Run command
	err := client.RunFlatCmd(&result, models.GearsPyExecute, qm.Key, args...)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// UNBLOCKING
	if qm.Unblocking {
		// when running with UNBLOCKING only operationId is returned
		frame := data.NewFrame("operationId")
		frame.Fields = append(frame.Fields, data.NewField("operationId", nil, []string{string(result.([]byte))}))

		// Adding frame to response
		response.Frames = append(response.Frames, frame)
		return response
	}

	// New Frame for results
	frameWithResults := data.NewFrame("results")
	frameWithResults.Fields = append(frameWithResults.Fields, data.NewField("results", nil, []string{}))

	// New Frame for errors
	frameWithErrors := data.NewFrame("errors")
	frameWithErrors.Fields = append(frameWithErrors.Fields, data.NewField("errors", nil, []string{}))

	// Adding frames to response
	response.Frames = append(response.Frames, frameWithResults)
	response.Frames = append(response.Frames, frameWithErrors)

	// Parse result
	switch value := result.(type) {
	case string:
		return response
	case []interface{}:
		// Inserting results
		for _, entry := range value[0].([]interface{}) {
			frameWithResults.AppendRow(string(entry.([]byte)))
		}

		// Inserting errors
		for _, entry := range value[1].([]interface{}) {
			frameWithErrors.AppendRow(string(entry.([]byte)))
		}
		return response
	default:
		log.DefaultLogger.Error("Unexpected type received", "value", value, "type", reflect.TypeOf(value).String())
		return response
	}
}

/**
 * RG.PYDUMPREQS
 *
 * Returns a list of all the python requirements available (with information about each requirement).
 * @see https://oss.redislabs.com/redisgears/commands.html#rgpydumpreqs
 */
func queryRgPydumpReqs(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Using radix marshaling of key-value arrays to structs
	var reqs []models.PyDumpReq

	// Run command
	err := client.RunCmd(&reqs, models.GearsPyDumpReqs)

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frame
	frame := data.NewFrame(qm.Command)
	response.Frames = append(response.Frames, frame)

	// New Fields
	frame.Fields = append(frame.Fields, data.NewField("GearReqVersion", nil, []int64{}))
	frame.Fields = append(frame.Fields, data.NewField("Name", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("IsDownloaded", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("IsInstalled", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("CompiledOs", nil, []string{}))
	frame.Fields = append(frame.Fields, data.NewField("Wheels", nil, []string{}))

	// Requirements
	for _, req := range reqs {
		var wheels string

		// Parse wheels
		switch value := req.Wheels.(type) {
		case []byte:
			wheels = string(value)
		case []string:
			wheels = strings.Join(value, ", ")
		case []interface{}:
			var values []string
			for _, entry := range value {
				values = append(values, string(entry.([]byte)))
			}

			wheels = strings.Join(values, ", ")
		default:
			log.DefaultLogger.Error("Unexpected type received", "value", value, "type", reflect.TypeOf(value).String())
			wheels = "Can't parse output"
		}

		frame.AppendRow(req.GearReqVersion, req.Name, req.IsDownloaded, req.IsInstalled, req.CompiledOs, wheels)
	}

	return response
}

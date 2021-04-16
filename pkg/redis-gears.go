package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

/**
 * RG.PYSTATS Radix marshaling
 */
type pystats struct {
	TotalAllocated int64 `redis:"TotalAllocated"`
	PeakAllocated  int64 `redis:"PeakAllocated"`
	CurrAllocated  int64 `redis:"CurrAllocated"`
}

/**
 * RG.DUMPREGISTRATIONS Radix marshaling
 */
type dumpregistrations struct {
	ID               string           `redis:"id"`
	Reader           string           `redis:"reader"`
	Desc             string           `redis:"desc"`
	RegistrationData registrationData `redis:"RegistrationData"`
	PD               string           `redis:"PD"`
}

/**
 * Registration data for RG.DUMPREGISTRATIONS Radix marshaling
 */
type registrationData struct {
	Mode         string                 `redis:"mode"`
	NumTriggered int64                  `redis:"numTriggered"`
	NumSuccess   int64                  `redis:"numSuccess"`
	NumFailures  int64                  `redis:"numFailures"`
	NumAborted   int64                  `redis:"numAborted"`
	LastError    string                 `redis:"lastError"`
	Args         map[string]interface{} `redis:"args"`
	Status       string                 `redis:"status"`
}

/**
 * RG.PYDUMPREQS Radix marshaling
 */
type pydumpreq struct {
	GearReqVersion int64    `redis:"GearReqVersion"`
	Name           string   `redis:"Name"`
	IsDownloaded   string   `redis:"IsDownloaded"`
	IsInstalled    string   `redis:"IsInstalled"`
	CompiledOs     string   `redis:"CompiledOs"`
	Wheels         []string `redis:"Wheels"`
}

/**
 * RG.PYSTATS
 *
 * Returns memory usage statistics from the Python interpreter
 * @see https://oss.redislabs.com/redisgears/commands.html#rgpystats
 */
func queryRgPystats(qm queryModel, client redisClient) backend.DataResponse {
	response := backend.DataResponse{}

	// Using radix marshaling of key-value arrays to structs
	var stats pystats

	// Run command
	err := client.RunCmd(&stats, "RG.PYSTATS")

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
	var models []dumpregistrations

	// Run command
	err := client.RunCmd(&models, "RG.DUMPREGISTRATIONS")

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
	for _, model := range models {
		// Merging args to string like "key"="value"\n
		args := new(bytes.Buffer)
		for key, value := range model.RegistrationData.Args {
			fmt.Fprintf(args, "\"%s\"=\"%s\"\n", key, value)
		}

		frame.AppendRow(model.ID, model.Reader, model.Desc, model.PD, model.RegistrationData.Mode,
			model.RegistrationData.NumTriggered, model.RegistrationData.NumSuccess, model.RegistrationData.NumFailures,
			model.RegistrationData.NumAborted, model.RegistrationData.LastError, args.String(), model.RegistrationData.Status)
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
	err := client.RunFlatCmd(&result, "RG.PYEXECUTE", qm.Key, args...)

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
	var reqs []pydumpreq

	// Run command
	err := client.RunCmd(&reqs, "RG.PYDUMPREQS")

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
		frame.AppendRow(req.GearReqVersion, req.Name, req.IsDownloaded, req.IsInstalled, req.CompiledOs, strings.Join(req.Wheels, ", "))
	}

	return response
}

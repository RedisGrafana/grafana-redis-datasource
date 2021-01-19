package main

import (
	"bytes"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
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
	Mode         string            `redis:"mode"`
	NumTriggered int64             `redis:"numTriggered"`
	NumSuccess   int64             `redis:"numSuccess"`
	NumFailures  int64             `redis:"numFailures"`
	NumAborted   int64             `redis:"numAborted"`
	LastError    string            `redis:"lastError"`
	Args         map[string]string `redis:"args"`
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

	for _, model := range models {
		// Merging args to string like "key"="value"\n
		args := new(bytes.Buffer)
		for key, value := range model.RegistrationData.Args {
			fmt.Fprintf(args, "\"%s\"=\"%s\"\n", key, value)
		}

		frame.AppendRow(model.ID, model.Reader, model.Desc, model.PD, model.RegistrationData.Mode,
			model.RegistrationData.NumTriggered, model.RegistrationData.NumSuccess, model.RegistrationData.NumFailures,
			model.RegistrationData.NumAborted, model.RegistrationData.LastError, args.String())
	}

	return response
}

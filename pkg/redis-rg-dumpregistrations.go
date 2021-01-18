package main

import (
	"bytes"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

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

	//New Fields
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
		frame.AppendRow(model.Id, model.Reader, model.Desc, model.PD, model.RegistrationData.Mode,
			model.RegistrationData.NumTriggered, model.RegistrationData.NumSuccess, model.RegistrationData.NumFailures,
			model.RegistrationData.NumAborted, model.RegistrationData.LastError, args.String())
	}

	return response
}

type dumpregistrations struct {
	Id               string           `redis:"id"`
	Reader           string           `redis:"reader"`
	Desc             string           `redis:"desc"`
	RegistrationData registrationData `redis:"RegistrationData"`
	PD               string           `redis:"PD"`
}

type registrationData struct {
	Mode         string            `redis:"mode"`
	NumTriggered int64             `redis:"numTriggered"`
	NumSuccess   int64             `redis:"numSuccess"`
	NumFailures  int64             `redis:"numFailures"`
	NumAborted   int64             `redis:"numAborted"`
	LastError    string            `redis:"lastError"`
	Args         map[string]string `redis:"args"`
}

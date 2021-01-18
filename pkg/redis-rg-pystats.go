package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
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
	var stats pystats

	// Run command
	err := client.RunCmd(&stats, "RG.PYSTATS")

	// Check error
	if err != nil {
		return errorHandler(response, err)
	}

	// New Frames
	frame := data.NewFrame(qm.Command)
	response.Frames = append(response.Frames, frame)

	//New Fields
	frame.Fields = append(frame.Fields, data.NewField("TotalAllocated", nil, []int64{stats.TotalAllocated}))
	frame.Fields = append(frame.Fields, data.NewField("PeakAllocated", nil, []int64{stats.PeakAllocated}))
	frame.Fields = append(frame.Fields, data.NewField("CurrAllocated", nil, []int64{stats.CurrAllocated}))
	return response
}

type pystats struct {
	TotalAllocated int64 `redis:"TotalAllocated"`
	PeakAllocated  int64 `redis:"PeakAllocated"`
	CurrAllocated  int64 `redis:"CurrAllocated"`
}

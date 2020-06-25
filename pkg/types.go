package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/mediocregopher/radix/v3"
)

type instanceSettings struct {
	client *radix.Pool
}

type redisDatasource struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirements
	// but a best practice that we recommend that you follow.
	im instancemgmt.InstanceManager
}

// Data Model
type dataModel struct {
	Size int `json:"size"`
}

// Query Model
type queryModel struct {
	Key         string `json:"key"`
	Field       string `json:"field"`
	Filter      string `json:"filter"`
	Command     string `json:"command"`
	Aggregation string `json:"aggregation"`
	Bucket      string `json:"bucket"`
	Legend      string `json:"legend"`
	Value       string `json:"value"`
}

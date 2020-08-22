package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/mediocregopher/radix/v3"
)

/**
 * Instance Settings
 */
type instanceSettings struct {
	client *radix.Pool
}

/**
 * 	The instance manager can help with lifecycle management of datasource instances in plugins.
 */
type redisDatasource struct {
	im instancemgmt.InstanceManager
}

/**
 * Configuration Data Model
 */
type dataModel struct {
	PoolSize       int  `json:"poolSize"`
	Timeout        int  `json:"timeout"`
	PingInterval   int  `json:"pingInterval"`
	PipelineWindow int  `json:"pipelineWindow"`
	TLSAuth        bool `json:"tlsAuth"`
	TLSSkipVerify  bool `json:"tlsSkipVerify"`
}

/*
 * Query Model
 */
type queryModel struct {
	Type        string `json:"type"`
	Query       string `json:"query"`
	Key         string `json:"key"`
	Field       string `json:"field"`
	Filter      string `json:"filter"`
	Command     string `json:"command"`
	Aggregation string `json:"aggregation"`
	Bucket      string `json:"bucket"`
	Legend      string `json:"legend"`
	Value       string `json:"value"`
	Section     string `json:"section"`
	Fill        bool   `json:"fill"`
}

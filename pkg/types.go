package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/mediocregopher/radix/v3"
)

// ClientInterface is an interface that represents the skeleton of a connection to Redis ( cluster, standalone, or sentinel )
type ClientInterface interface {
	Do(a radix.Action) error
	Close() error
}

/**
 * Instance Settings
 */
type instanceSettings struct {
	client ClientInterface
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
	PoolSize       int    `json:"poolSize"`
	Timeout        int    `json:"timeout"`
	PingInterval   int    `json:"pingInterval"`
	PipelineWindow int    `json:"pipelineWindow"`
	TLSAuth        bool   `json:"tlsAuth"`
	TLSSkipVerify  bool   `json:"tlsSkipVerify"`
	Client         string `json:"client"`
	SentinelName   string `json:"sentinelName"`
	ACL            bool   `json:"acl"`
	User           string `json:"user"`
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
	Size        int    `json:"size"`
	Fill        bool   `json:"fill"`
	Streaming   bool   `json:"streaming"`
}

package main

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
)

/**
 * Instance Settings
 */
type instanceSettings struct {
	client redisClient
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
	Type         string `json:"type"`
	Query        string `json:"query"`
	Key          string `json:"keyName"`
	Field        string `json:"field"`
	Filter       string `json:"filter"`
	Command      string `json:"command"`
	Aggregation  string `json:"aggregation"`
	Bucket       int    `json:"bucket"`
	Legend       string `json:"legend"`
	Value        string `json:"value"`
	Section      string `json:"section"`
	Size         int    `json:"size"`
	Fill         bool   `json:"fill"`
	Streaming    bool   `json:"streaming"`
	CLI          bool   `json:"cli"`
	Cursor       string `json:"cursor"`
	Match        string `json:"match"`
	Count        int    `json:"count"`
	Samples      int    `json:"samples"`
	Unblocking   bool   `json:"unblocking"`
	Requirements string `json:"requirements"`
	Start        string `json:"start"`
	End          string `json:"end"`
	Cypher       string `json:"cypher"`
}

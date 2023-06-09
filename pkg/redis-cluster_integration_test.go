//go:build clusterIntegration
// +build clusterIntegration

package main

import (
	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCluster(t *testing.T) {
	radixClient, err := radix.NewCluster([]string{"redis://redis-cluster1:6379", "redis://redis-cluster2:6379", "redis://redis-cluster3:6379"})

	require.Nil(t, err)
	var client = &radixV3Impl{radixClient}
	var result interface{}

	client.RunCmd(&result, "PING")

	require.Equal(t, "PONG", result.(string))
}

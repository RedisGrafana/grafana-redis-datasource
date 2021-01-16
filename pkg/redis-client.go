package main

import (
	"crypto/tls"
	"crypto/x509"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/mediocregopher/radix/v3"
)

/**
 * Configuration Data Model for redisClient
 */
type redisClientConfiguration struct {
	URL            string
	Client         string
	Timeout        int
	PoolSize       int
	PingInterval   int
	PipelineWindow int
	ACL            bool
	TLSAuth        bool
	TLSSkipVerify  bool
	User           string
	Password       string
	TLSCACert      string
	TLSClientCert  string
	SentinelName   string
}

/**
 * Interface for running redis commands without explicit dependencies to 3-rd party libraries
 */
type redisClient interface {
	RunFlatCmd(rcv interface{}, cmd, key string, args ...interface{}) error
	RunCmd(rcv interface{}, cmd string, args ...string) error
	RunBatchFlatCmd(commands []flatCommandArgs) error
	Close() error
}

type flatCommandArgs struct {
	rcv  interface{}
	cmd  string
	key  string
	args []interface{}
}

// radixClient is an interface that represents the skeleton of a connection to Redis ( cluster, standalone, or sentinel)
type radixClient interface {
	Do(a radix.Action) error
	Close() error
}

// radixV3Impl is an implementation of redisClient using the radix/v3 library
type radixV3Impl struct {
	radixClient radixClient
}

// Execute Radix FlatCmd
func (client *radixV3Impl) RunFlatCmd(rcv interface{}, cmd, key string, args ...interface{}) error {
	return client.radixClient.Do(radix.FlatCmd(rcv, cmd, key, args...))
}

// Execute Batch FlatCmd
func (client *radixV3Impl) RunBatchFlatCmd(commands []flatCommandArgs) error {
	var actions []radix.CmdAction
	for _, command := range commands {
		actions = append(actions, radix.FlatCmd(command.rcv, command.cmd, command.key, command.args...))
	}

	// Pipeline commands
	pipeline := radix.Pipeline(actions...)
	return client.radixClient.Do(pipeline)
}

// Execute Radix Cmd
func (client *radixV3Impl) RunCmd(rcv interface{}, cmd string, args ...string) error {
	return client.radixClient.Do(radix.Cmd(rcv, cmd, args...))
}

// Close connection
func (client *radixV3Impl) Close() error {
	return client.radixClient.Close()
}

// creates new radixV3Impl implementation of redisClient interface
func newRadixV3Client(configuration redisClientConfiguration) (redisClient, error) {
	var radixClient radixClient

	var err error

	// Set up connection
	connFunc := func(network, addr string) (radix.Conn, error) {
		opts := []radix.DialOpt{radix.DialTimeout(time.Duration(configuration.Timeout) * time.Second)}

		// Authentication
		if configuration.Password != "" {
			// If ACL enabled
			if configuration.ACL {
				opts = append(opts, radix.DialAuthUser(configuration.User, configuration.Password))
			} else {
				opts = append(opts, radix.DialAuthPass(configuration.Password))
			}
		}

		// TLS Authentication
		if configuration.TLSAuth {
			// TLS Config
			tlsConfig := &tls.Config{
				InsecureSkipVerify: configuration.TLSSkipVerify,
			}

			// Certification Authority
			if configuration.TLSCACert != "" {
				caPool := x509.NewCertPool()
				ok := caPool.AppendCertsFromPEM([]byte(configuration.TLSCACert))
				if ok {
					tlsConfig.RootCAs = caPool
				}
			}

			// Certificate and Key
			if configuration.TLSClientCert != "" {
				cert, err := tls.X509KeyPair([]byte(configuration.TLSClientCert), []byte(configuration.TLSClientCert))
				if err == nil {
					tlsConfig.Certificates = []tls.Certificate{cert}
				} else {
					log.DefaultLogger.Error("X509KeyPair", "Error", err)
					return nil, err
				}
			}

			// Add TLS Config
			opts = append(opts, radix.DialUseTLS(tlsConfig))
		}

		return radix.Dial(network, addr, opts...)
	}

	// Pool with specified Ping Interval, Pipeline Window and Timeout
	poolFunc := func(network, addr string) (radix.Client, error) {
		return radix.NewPool(network, addr, configuration.PoolSize, radix.PoolConnFunc(connFunc),
			radix.PoolPingInterval(time.Duration(configuration.PingInterval)*time.Second/time.Duration(configuration.PoolSize+1)),
			radix.PoolPipelineWindow(time.Duration(configuration.PipelineWindow)*time.Microsecond, 0))
	}

	// Client Type
	switch configuration.Client {
	case "cluster":
		radixClient, err = radix.NewCluster(strings.Split(configuration.URL, ","), radix.ClusterPoolFunc(poolFunc))
	case "sentinel":
		radixClient, err = radix.NewSentinel(configuration.SentinelName, strings.Split(configuration.URL, ","), radix.SentinelConnFunc(connFunc),
			radix.SentinelPoolFunc(poolFunc))
	case "socket":
		radixClient, err = poolFunc("unix", configuration.URL)
	default:
		radixClient, err = poolFunc("tcp", configuration.URL)
	}

	if err != nil {
		return nil, err
	}

	// Return Radix client
	var client = &radixV3Impl{radixClient}
	return client, nil
}

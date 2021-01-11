package main

import (
	"github.com/mediocregopher/radix/v3"
)

type testClient struct {
	rcv interface{}
	err error
}

func (client testClient) RunFlatCmd(rcv interface{}, cmd, key string, args ...interface{}) error {
	stub := radix.Stub("tcp", "127.0.0.1:6379", func(args []string) interface{} {
		return client.rcv
	})
	var _ = stub.Do(radix.FlatCmd(rcv, cmd, key, args...))
	return client.err
}
func (client testClient) RunCmd(rcv interface{}, cmd string, args ...string) error {
	stub := radix.Stub("tcp", "127.0.0.1:6379", func(args []string) interface{} {
		return client.rcv
	})
	var _ = stub.Do(radix.Cmd(rcv, cmd, args...))
	return client.err
}
func (client testClient) Close() error {
	return client.err
}

type panickingClient struct {
}

func (client panickingClient) RunFlatCmd(rcv interface{}, cmd, key string, args ...interface{}) error {
	panic("Panic")
}
func (client panickingClient) RunCmd(rcv interface{}, cmd string, args ...string) error {
	panic("Panic")
}
func (client panickingClient) Close() error {
	return nil
}

type valueToCheckInResponse struct {
	frameIndex int
	fieldIndex int
	rowIndex   int
	value      interface{}
}

type valueToCheckByLabelInResponse struct {
	frameIndex int
	fieldName  string
	rowIndex   int
	value      interface{}
}

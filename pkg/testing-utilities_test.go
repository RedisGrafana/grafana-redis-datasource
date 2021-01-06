package main

import (
	"github.com/mediocregopher/radix/v3"
)

type testClient struct {
	rcv interface{}
	err error
}

func (client testClient) Do(action radix.Action) error {
	stub := radix.Stub("tcp", "127.0.0.1:6379", func(args []string) interface{} {
		return client.rcv
	})
	var _ = stub.Do(action)
	return client.err
}
func (client testClient) Close() error {
	return client.err
}

type panickingClient struct {
}

func (client panickingClient) Do(action radix.Action) error {
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

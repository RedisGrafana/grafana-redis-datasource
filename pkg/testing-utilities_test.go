package main

import (
  "github.com/mediocregopher/radix/v3"
)

type TestClient struct {
  rcv interface{}
  err error
}

func (client TestClient) Do(action radix.Action) error {
  stub := radix.Stub("tcp", "127.0.0.1:6379", func(args []string) interface{} {
    return client.rcv
  })
  var _ = stub.Do(action)
  return client.err
}
func (client TestClient) Close() error {
  return client.err
}

type valueToCheckInResponse struct {
  frameIndex int
  fieldIndex int
  rowIndex int
  value interface{}
}

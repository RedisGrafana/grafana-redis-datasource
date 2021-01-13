package main

import (
	"reflect"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/stretchr/testify/mock"
)

type testClient struct {
	rcv interface{}
	err error
	mock.Mock
}

func (client *testClient) RunFlatCmd(rcv interface{}, cmd, key string, args ...interface{}) error {
	if client.err != nil {
		return client.err
	} else {
		client.assignReceiver(rcv)
		return nil
	}

}
func (client *testClient) RunCmd(rcv interface{}, cmd string, args ...string) error {
	if client.err != nil {
		return client.err
	} else {
		client.assignReceiver(rcv)
		return nil
	}
}

func (client *testClient) assignReceiver(rcv interface{}) {
	switch rcv.(type) {
	case int:
		*(rcv.(*int)) = client.rcv.(int)
	case int64:
		*(rcv.(*int64)) = client.rcv.(int64)
	case float64:
		*(rcv.(*float64)) = client.rcv.(float64)
	case []string:
		*(rcv.(*[]string)) = client.rcv.([]string)
	case []interface{}:
		*(rcv.(*[]interface{})) = client.rcv.([]interface{})
	case [][]string:
		*(rcv.(*[][]string)) = client.rcv.([][]string)
	case map[string]string:
		*(rcv.(*map[string]string)) = client.rcv.(map[string]string)
	case map[string]interface{}:
		*(rcv.(*map[string]interface{})) = client.rcv.(map[string]interface{})
	case *string:
		*(rcv.(*string)) = client.rcv.(string)
	case interface{}:
		switch client.rcv.(type) {
		case int:
			*(rcv.(*interface{})) = client.rcv.(int)
		case int64:
			*(rcv.(*interface{})) = client.rcv.(int64)
		case float64:
			*(rcv.(*interface{})) = client.rcv.(float64)
		case []string:
			*(rcv.(*[]string)) = client.rcv.([]string)
		case []interface{}:
			*(rcv.(*interface{})) = client.rcv.([]interface{})
		case [][]string:
			*(rcv.(*[][]string)) = client.rcv.([][]string)
		case map[string]string:
			*(rcv.(*map[string]string)) = client.rcv.(map[string]string)
		case *string:
			*(rcv.(*string)) = client.rcv.(string)
		case string:
			*(rcv.(*interface{})) = client.rcv.(string)
		case []uint8:
			*(rcv.(*interface{})) = client.rcv.([]uint8)
		case map[string]interface{}:
			*(rcv.(*map[string]interface{})) = client.rcv.(map[string]interface{})
		default:
			panic("Unsupported type of client.rcv: " + reflect.TypeOf(client.rcv).String())
		}
	default:
		panic("Unsupported type of rcv: " + reflect.TypeOf(rcv).String())
	}
}
func (client *testClient) Close() error {
	args := client.Called()
	return args.Error(0)
}

type panickingClient struct {
}

func (client *panickingClient) RunFlatCmd(rcv interface{}, cmd, key string, args ...interface{}) error {
	panic("Panic")
}
func (client *panickingClient) RunCmd(rcv interface{}, cmd string, args ...string) error {
	panic("Panic")
}
func (client *panickingClient) Close() error {
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

type fakeInstanceManager struct {
	mock.Mock
}

func (im *fakeInstanceManager) Get(pluginContext backend.PluginContext) (instancemgmt.Instance, error) {
	args := im.Called(pluginContext)
	return args.Get(0), args.Error(1)
}

func (im *fakeInstanceManager) Do(pluginContext backend.PluginContext, fn instancemgmt.InstanceCallbackFunc) error {
	args := im.Called(pluginContext, fn)
	return args.Error(0)
}

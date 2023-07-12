package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/redisgrafana/grafana-redis-datasource/pkg/models"
	"github.com/stretchr/testify/mock"
)

/**
 * Test client
 */
type testClient struct {
	rcv          interface{}
	batchRcv     [][]interface{}
	batchErr     []error
	expectedArgs []string
	expectedCmd  string
	err          error
	batchCalls   int
	mock.Mock
}

/**
 * PANIC
 */
type panickingClient struct {
}

/**
 * Response
 */
type valueToCheckInResponse struct {
	frameIndex int
	fieldIndex int
	rowIndex   int
	value      interface{}
}

/**
 * Response
 */
type valueToCheckByLabelInResponse struct {
	frameIndex int
	fieldName  string
	rowIndex   int
	value      interface{}
}

/**
 * Fake Instance manager
 */
type fakeInstanceManager struct {
	mock.Mock
}

/**
 * FlatCmd()
 */
func (client *testClient) RunFlatCmd(rcv interface{}, cmd, key string, args ...interface{}) error {
	if client.err != nil {
		return client.err
	}

	assignReceiver(rcv, client.rcv)
	return nil
}

/**
 * Cmd()
 */
func (client *testClient) RunCmd(rcv interface{}, cmd string, args ...string) error {
	if client.err != nil {
		return client.err
	}

	if client.expectedArgs != nil {
		if !reflect.DeepEqual(args, client.expectedArgs) {
			return fmt.Errorf("expected args did not match actuall args\nExpected:%s\nActual:%s\n", client.expectedArgs, args)
		}
	}

	if client.expectedCmd != "" && client.expectedCmd != cmd {
		return fmt.Errorf("incorrect command, Expected:%s - Actual: %s", client.expectedCmd, cmd)
	}

	assignReceiver(rcv, client.rcv)
	return nil
}

/**
 * Pipeline execution using Batch
 */
func (client *testClient) RunBatchFlatCmd(commands []flatCommandArgs) error {
	for i, args := range commands {
		assignReceiver(args.rcv, client.batchRcv[client.batchCalls][i])
	}

	var err error
	if client.batchErr != nil && client.batchErr[client.batchCalls] != nil {
		err = client.batchErr[client.batchCalls]
	}

	client.batchCalls++
	return err
}

/**
 * Receiver
 */
func assignReceiver(to interface{}, from interface{}) {
	switch to.(type) {
	case int:
		*(to.(*int)) = from.(int)
	case int64:
		*(to.(*int64)) = from.(int64)
	case float64:
		*(to.(*float64)) = from.(float64)
	case []string:
		*(to.(*[]string)) = from.([]string)
	case []interface{}:
		*(to.(*[]interface{})) = from.([]interface{})
	case [][]string:
		*(to.(*[][]string)) = from.([][]string)
	case map[string]string:
		*(to.(*map[string]string)) = from.(map[string]string)
	case map[string]interface{}:
		*(to.(*map[string]interface{})) = from.(map[string]interface{})
	case *string:
		*(to.(*string)) = from.(string)
	case *models.PyStats:
		*(to.(*models.PyStats)) = from.(models.PyStats)
	case *xinfo:
		*(to.(*xinfo)) = from.(xinfo)
	case *[]models.DumpRegistrations:
		*(to.(*[]models.DumpRegistrations)) = from.([]models.DumpRegistrations)
	case *[]models.PyDumpReq:
		*(to.(*[]models.PyDumpReq)) = from.([]models.PyDumpReq)
	case interface{}:
		switch from.(type) {
		case int:
			*(to.(*interface{})) = from.(int)
		case int64:
			_, ok := to.(*int64)
			if ok {
				*(to.(*int64)) = from.(int64)
			} else {
				_, ok = to.(*interface{})
				if ok {
					*(to.(*interface{})) = from.(int64)
				}
			}
		case float64:
			*(to.(*interface{})) = from.(float64)
		case []string:
			*(to.(*[]string)) = from.([]string)
		case []interface{}:
			_, ok := to.(*[]interface{})
			if ok {
				*(to.(*[]interface{})) = from.([]interface{})
			} else {
				_, ok = to.(*interface{})
				if ok {
					*(to.(*interface{})) = from.([]interface{})
				}
			}
		case [][]string:
			*(to.(*[][]string)) = from.([][]string)
		case map[string]string:
			*(to.(*map[string]string)) = from.(map[string]string)
		case *string:
			*(to.(*string)) = from.(string)
		case string:
			*(to.(*interface{})) = from.(string)
		case []uint8:
			*(to.(*interface{})) = from.([]uint8)
		case map[string]interface{}:
			*(to.(*map[string]interface{})) = from.(map[string]interface{})
		default:
			panic("Unsupported type of from rcv: " + reflect.TypeOf(from).String())
		}
	default:
		panic("Unsupported type of to rcv: " + reflect.TypeOf(to).String())
	}
}

/**
 * Close session
 */
func (client *testClient) Close() error {
	args := client.Called()
	return args.Error(0)
}

/**
 * FlatCmd() Error
 */
func (client *panickingClient) RunFlatCmd(rcv interface{}, cmd, key string, args ...interface{}) error {
	panic("Panic")
}

/**
 * Cmd() Error
 */
func (client *panickingClient) RunCmd(rcv interface{}, cmd string, args ...string) error {
	panic("Panic")
}

/**
 * Close() Error
 */
func (client *panickingClient) Close() error {
	return nil
}

/**
 * Batch command
 */
func (client *panickingClient) RunBatchFlatCmd(commands []flatCommandArgs) error {
	panic("Panic")
}

/**
 * Get
 */
func (im *fakeInstanceManager) Get(ctx context.Context, pluginContext backend.PluginContext) (instancemgmt.Instance, error) {
	args := im.Called(pluginContext)
	return args.Get(0), args.Error(1)
}

/**
 * Do
 */
func (im *fakeInstanceManager) Do(ctx context.Context, pluginContext backend.PluginContext, fn instancemgmt.InstanceCallbackFunc) error {
	args := im.Called(pluginContext, fn)
	return args.Error(0)
}

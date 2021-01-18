package main

import (
	"testing"

	"github.com/mediocregopher/radix/v3"
	"github.com/stretchr/testify/require"
)

func TestRadixV3Impl(t *testing.T) {
	t.Parallel()
	t.Run("should run Cmd", func(t *testing.T) {
		t.Parallel()
		client := radixV3Impl{radix.Stub("tcp", "127.0.0.1:6379", func(args []string) interface{} {
			return args
		})}
		var result []string
		err := client.RunCmd(&result, "Command1", "Arg1", "Arg2")
		require.NoError(t, err)
		require.Equal(t, []string{"Command1", "Arg1", "Arg2"}, result)

	})
	t.Run("should run flatCmd", func(t *testing.T) {
		t.Parallel()
		client := radixV3Impl{radix.Stub("tcp", "127.0.0.1:6379", func(args []string) interface{} {
			return args
		})}
		var result []string
		err := client.RunFlatCmd(&result, "Command2", "SomeKey", "Arg1", "Arg2")
		require.NoError(t, err)
		require.Equal(t, []string{"Command2", "SomeKey", "Arg1", "Arg2"}, result)
	})
	t.Run("should have RunBatchFlatCmd", func(t *testing.T) {
		t.Parallel()
		client := radixV3Impl{radix.Stub("tcp", "127.0.0.1:6379", func(args []string) interface{} {
			return args
		})}
		var result []string
		err := client.RunBatchFlatCmd([]flatCommandArgs{{
			rcv:  &result,
			cmd:  "Command2",
			key:  "SomeKey",
			args: []interface{}{"Arg1", "Arg2"},
		}})
		require.NoError(t, err)
		require.Equal(t, []string{"Command2", "SomeKey", "Arg1", "Arg2"}, result)
	})
	t.Run("should have close method", func(t *testing.T) {
		t.Parallel()
		client := radixV3Impl{radix.Stub("tcp", "127.0.0.1:6379", func(args []string) interface{} {
			return args
		})}
		err := client.Close()
		require.NoError(t, err)
	})
}

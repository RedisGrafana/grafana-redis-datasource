package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/require"
)

// recordingDialer records requested addresses and connects them all to the fake server.
type recordingDialer struct {
	mu     sync.Mutex
	addrs  []string
	target string
}

func (d *recordingDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	d.mu.Lock()
	d.addrs = append(d.addrs, addr)
	d.mu.Unlock()

	return (&net.Dialer{}).DialContext(ctx, "tcp", d.target)
}

func (d *recordingDialer) dialed() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]string{}, d.addrs...)
}

// fakeRedis answers just enough RESP for the radix client and records AUTH commands.
type fakeRedis struct {
	listener net.Listener
	mu       sync.Mutex
	auths    []string
}

func newFakeRedis(t *testing.T) *fakeRedis {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	server := &fakeRedis{listener: listener}
	t.Cleanup(func() { listener.Close() })

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.serve(conn)
		}
	}()

	return server
}

func (s *fakeRedis) addr() string {
	return s.listener.Addr().String()
}

func (s *fakeRedis) authCommands() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string{}, s.auths...)
}

func readCommand(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, err
	}

	args := make([]string, 0, count)
	for i := 0; i < count; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		size, err := strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			return nil, err
		}

		buf := make([]byte, size+2)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return nil, err
		}
		args = append(args, string(buf[:size]))
	}

	return args, nil
}

func (s *fakeRedis) serve(conn net.Conn) {
	defer conn.Close()

	_, port, _ := net.SplitHostPort(s.addr())
	reader := bufio.NewReader(conn)

	for {
		args, err := readCommand(reader)
		if err != nil {
			return
		}

		switch strings.ToUpper(args[0]) {
		case "AUTH":
			s.mu.Lock()
			s.auths = append(s.auths, strings.Join(args[1:], " "))
			s.mu.Unlock()
			fmt.Fprint(conn, "+OK\r\n")
		case "PING":
			fmt.Fprint(conn, "+PONG\r\n")
		case "CLUSTER":
			// A single node owning every slot
			fmt.Fprintf(conn, "*1\r\n*3\r\n:0\r\n:16383\r\n*2\r\n$9\r\n127.0.0.1\r\n:%s\r\n", port)
		default:
			fmt.Fprint(conn, "+OK\r\n")
		}
	}
}

func TestProxiedStandaloneClient(t *testing.T) {
	tests := []struct {
		name          string
		configuration redisClientConfiguration
		expectedAuth  []string
	}{
		{
			"should connect through the proxy without auth",
			redisClientConfiguration{URL: "redis.internal:6379"},
			[]string{},
		},
		{
			"should authenticate with password only",
			redisClientConfiguration{URL: "redis.internal:6379", Password: "secret"},
			[]string{"secret"},
		},
		{
			"should authenticate with ACL user",
			redisClientConfiguration{URL: "redis.internal:6379", ACL: true, User: "alice", Password: "secret"},
			[]string{"alice secret"},
		},
		{
			"should ignore the user when ACL is disabled",
			redisClientConfiguration{URL: "redis.internal:6379", User: "alice", Password: "secret"},
			[]string{"secret"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := newFakeRedis(t)
			dialer := &recordingDialer{target: server.addr()}

			tt.configuration.PoolSize = 1
			tt.configuration.Timeout = 5
			tt.configuration.ProxyDialer = dialer

			client, err := newRadixV3Client(tt.configuration)
			require.NoError(t, err)
			defer client.Close()

			var reply string
			require.NoError(t, client.RunCmd(&reply, "PING"))
			require.Equal(t, "PONG", reply)

			require.Contains(t, dialer.dialed(), "redis.internal:6379")
			require.Equal(t, tt.expectedAuth, server.authCommands())
		})
	}
}

func TestProxiedClusterClient(t *testing.T) {
	server := newFakeRedis(t)
	dialer := &recordingDialer{target: server.addr()}

	client, err := newRadixV3Client(redisClientConfiguration{
		URL:         "seed.internal:7000",
		Client:      "cluster",
		PoolSize:    1,
		Timeout:     5,
		ProxyDialer: dialer,
	})
	require.NoError(t, err)
	defer client.Close()

	// The seed and the node learned from the topology are both dialed through the proxy
	require.Contains(t, dialer.dialed(), "seed.internal:7000")
	require.Contains(t, dialer.dialed(), server.addr())
}

func TestProxiedClientRejectsSocket(t *testing.T) {
	_, err := newRadixV3Client(redisClientConfiguration{
		URL:         "/var/run/redis.sock",
		Client:      "socket",
		PoolSize:    1,
		ProxyDialer: &recordingDialer{},
	})
	require.Error(t, err)
}

func TestNewProxyDialer(t *testing.T) {
	t.Run("should not use the proxy when disabled on the data source", func(t *testing.T) {
		dialer, err := newProxyDialer(context.Background(), backend.DataSourceInstanceSettings{
			JSONData: []byte(`{"enableSecureSocksProxy": false}`),
		})
		require.NoError(t, err)
		require.Nil(t, dialer)
	})

	t.Run("should fail instead of connecting directly when Grafana has no proxy configured", func(t *testing.T) {
		dialer, err := newProxyDialer(context.Background(), backend.DataSourceInstanceSettings{
			JSONData: []byte(`{"enableSecureSocksProxy": true}`),
		})
		require.Error(t, err)
		require.Nil(t, dialer)
	})
}

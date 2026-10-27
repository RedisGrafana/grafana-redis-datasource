package main

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/proxy"
	"github.com/mediocregopher/radix/v3"
)

// Matches radix's default ACL user, for which it sends AUTH with the password only.
const defaultAuthUser = "default"

// contextDialer dials connections through the Secure SOCKS Proxy.
type contextDialer interface {
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}

// newProxyDialer returns nil when the Secure SOCKS Proxy is not enabled on the data source.
func newProxyDialer(ctx context.Context, setting backend.DataSourceInstanceSettings) (contextDialer, error) {
	opts, err := setting.ProxyOptionsFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if opts == nil {
		return nil, nil
	}

	client := proxy.New(opts)
	if !client.SecureSocksProxyEnabled() {
		return nil, errors.New("secure socks proxy is enabled on the data source but not configured in Grafana")
	}

	dialer, err := client.NewSecureSocksProxyContextDialer()
	if err != nil {
		return nil, err
	}

	ctxDialer, ok := dialer.(contextDialer)
	if !ok {
		return nil, errors.New("secure socks proxy dialer does not support contexts")
	}
	return ctxDialer, nil
}

// deadlineConn applies the timeout to every read and write, like radix's own dial options.
type deadlineConn struct {
	net.Conn
	timeout time.Duration
}

func (c *deadlineConn) Read(b []byte) (int, error) {
	if c.timeout > 0 {
		if err := c.Conn.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
			return 0, err
		}
	}
	return c.Conn.Read(b)
}

func (c *deadlineConn) Write(b []byte) (int, error) {
	if c.timeout > 0 {
		if err := c.Conn.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
			return 0, err
		}
	}
	return c.Conn.Write(b)
}

// dialProxied replaces radix.Dial, which cannot use a custom dialer. An empty user sends AUTH with the password only.
func dialProxied(dialer contextDialer, network, addr string, configuration redisClientConfiguration, user, password string) (radix.Conn, error) {
	timeout := time.Duration(configuration.Timeout) * time.Second

	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}

	netConn, err := dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}

	tlsConfig, err := getTLSConfig(configuration)
	if err != nil {
		netConn.Close()
		return nil, err
	}

	if tlsConfig != nil {
		// tls.Client does not infer the server name from the address
		if tlsConfig.ServerName == "" {
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				host = addr
			}
			tlsConfig.ServerName = host
		}

		tlsConn := tls.Client(netConn, tlsConfig)
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			netConn.Close()
			return nil, err
		}
		netConn = tlsConn
	}

	conn := radix.NewConn(&deadlineConn{Conn: netConn, timeout: timeout})

	if user != "" && user != defaultAuthUser {
		err = conn.Do(radix.Cmd(nil, "AUTH", user, password))
	} else if password != "" {
		err = conn.Do(radix.Cmd(nil, "AUTH", password))
	}
	if err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

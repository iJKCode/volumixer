package pulse

import (
	"fmt"
	"github.com/jfreymuth/pulse/proto"
	"os"
	"path"
	"sync"
	"time"
)

type Client struct {
	index uint32
	api   *proto.Client
	close func() error
}

type clientConfig struct {
	props   proto.PropList
	server  string
	events  chan any
	timeout time.Duration
}

type Option func(*clientConfig)

func NewClient(opts ...Option) (*Client, error) {

	// process client options
	config := clientConfig{
		server: "",
		props: proto.PropList{
			"application.name":           proto.PropListString(path.Base(os.Args[0])),
			"application.icon_name":      proto.PropListString("audio-x-generic"),
			"application.process.id":     proto.PropListString(fmt.Sprintf("%d", os.Getpid())),
			"application.process.binary": proto.PropListString(os.Args[0]),
		},
	}
	for _, opt := range opts {
		opt(&config)
	}

	// configure connection
	api, conn, err := proto.Connect(config.server)
	if err != nil {
		return nil, err
	}
	if config.timeout != 0 {
		api.SetTimeout(config.timeout)
	}

	// set client properties
	var info proto.SetClientNameReply
	err = api.Request(&proto.SetClientName{Props: config.props}, &info)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}

	// return client instance
	return &Client{
		api: api,
		close: sync.OnceValue(func() error {
			return conn.Close()
		}),
		index: info.ClientIndex,
	}, nil
}

func (c *Client) Close() error {
	return c.close()
}

func WithName(name string) Option {
	return WithProperty("application.name", proto.PropListString(name))
}

func WithProperty(key string, value proto.PropListEntry) Option {
	return func(c *clientConfig) { c.props[key] = value }
}

func WithServerString(server string) Option {
	return func(c *clientConfig) { c.server = server }
}

func WithTimeout(d time.Duration) Option {
	return func(c *clientConfig) { c.timeout = d }
}

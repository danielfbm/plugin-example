package extension

import (
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-plugin"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Manager interface {
	Load(name string, opts PluginLoadOptions) (err error)
	Get(name string) (client plugin.ClientProtocol, err error)
	Stop()
}

var pluginMap = map[string]plugin.Plugin{
	"foo":      &FooPlugin{},
	"bar":      &BarPlugin{},
	"bar_grpc": &BarPlugin{},
}

var pluginScheme = schema.GroupResource{
	Group:    "plugin",
	Resource: "plugin",
}

var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

type PluginLoadOptions struct {
	Config *plugin.ClientConfig

	NetworkPluginAddress string
	networkAddress       net.Addr
	LocalPluginPath      string
}

func (opts PluginLoadOptions) Copy() PluginLoadOptions {
	return opts
}

func (opts PluginLoadOptions) Validate() (res PluginLoadOptions, err error) {
	if opts.Config == nil {
		err = errors.NewBadRequest("Config is empty")
		return
	}
	if opts.NetworkPluginAddress == "" && opts.LocalPluginPath == "" && opts.Config.Cmd == nil && opts.Config.Reattach == nil {
		err = errors.NewBadRequest("No local or network plugin configuration provided")
	}
	if opts.LocalPluginPath != "" {
		if _, err = filepath.EvalSymlinks(opts.LocalPluginPath); err != nil {
			err = errors.NewInternalError(err)
		}

	} else if opts.NetworkPluginAddress != "" {
		opts.networkAddress, err = net.ResolveTCPAddr("tcp", opts.NetworkPluginAddress)
	}

	res = opts.init()
	return
}

func (opts PluginLoadOptions) init() PluginLoadOptions {
	opts.Config.Plugins = pluginMap
	if opts.Config.HandshakeConfig.MagicCookieKey == "" {
		opts.Config.HandshakeConfig = HandshakeConfig
	}
	if opts.LocalPluginPath != "" {
		opts.Config.Cmd = exec.Command(opts.LocalPluginPath)
		opts.Config.Reattach = nil
	} else if opts.NetworkPluginAddress != "" {
		opts.Config.Cmd = nil
		opts.Config.Reattach = &plugin.ReattachConfig{
			Protocol: plugin.ProtocolNetRPC,
			Addr:     opts.networkAddress,
		}
	}
	fmt.Println("opts", opts, "reatach", opts.Config.Reattach, "cmd", opts.Config.Cmd, "netAddr", opts.networkAddress)
	return opts
}

type clientOptions struct {
	Client  *plugin.Client
	Options PluginLoadOptions
}

func NewManager() Manager {
	mgr := &manager{}
	mgr.init()
	return mgr
}

type manager struct {
	clients  map[string]*clientOptions
	initOnce sync.Once
}

func (m *manager) init() {
	m.initOnce.Do(func() {
		m.clients = make(map[string]*clientOptions)
	})
}

func (m *manager) Load(name string, opts PluginLoadOptions) (err error) {
	m.init()
	// if m.clients[name]
	if opts, err = opts.Validate(); err != nil {
		return
	}
	clt := plugin.NewClient(opts.Config)
	if _, err = clt.Client(); err != nil {
		return
	}

	m.clients[name] = &clientOptions{Client: clt, Options: opts}
	return
}

func (m *manager) Get(name string) (client plugin.ClientProtocol, err error) {
	m.init()
	c, ok := m.clients[name]
	if !ok {
		err = errors.NewNotFound(pluginScheme, name)
		// err = fmt.Errorf("No %s plugin found", name)
		return
	}
	if client, err = c.Client.Client(); err != nil {
		err = errors.NewInternalError(err)
	}
	return

}

func (m *manager) Stop() {
	m.init()
	for _, c := range m.clients {
		if c.Options.Config.Cmd != nil {
			c.Client.Kill()
		}
	}
}

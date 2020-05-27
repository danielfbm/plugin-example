package extension

import (
	"fmt"
	"net"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/hashicorp/go-plugin"
)

type Manager interface {
	Load(name string, opts PluginLoadOptions) (err error)
	Get(name string) (client plugin.ClientProtocol, err error)
	Stop()
}

var pluginMap = map[string]plugin.Plugin{
	"foo": &FooPlugin{},
	"bar": &BarPlugin{},
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

func (opts PluginLoadOptions) Validate() (res PluginLoadOptions, err error) {
	if opts.Config == nil {
		err = fmt.Errorf("Config is empty")
		return
	}
	if opts.NetworkPluginAddress == "" && opts.LocalPluginPath == "" && opts.Config.Cmd == nil && opts.Config.Reattach == nil {
		err = fmt.Errorf("No local or network plugin configuration provided")
	}
	if opts.LocalPluginPath != "" {
		_, err = filepath.EvalSymlinks(opts.LocalPluginPath)

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
		err = fmt.Errorf("No %s plugin found", name)
		return
	}
	client, err = c.Client.Client()
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

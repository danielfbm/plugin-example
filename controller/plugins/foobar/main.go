package main

import (
	"os"

	ext "github.com/danielfbm/plugin-example/controller/extension"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type FooBarPlugin struct {
	logger hclog.Logger
}

func (b *FooBarPlugin) Bars() []string {
	b.logger.Debug("FooBarPlugin.Bars: called")
	return []string{"foo", "bar"}
}

func (b *FooBarPlugin) Foos() (string, error) {
	b.logger.Debug("FooBarPlugin.Foos: called")
	return "FooBars!", nil
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	plug := &FooBarPlugin{logger: logger}

	var pluginMap = map[string]plugin.Plugin{
		"foo": &ext.FooPlugin{Impl: plug},
		"bar": &ext.BarPlugin{Impl: plug},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: ext.HandshakeConfig,
		Plugins:         pluginMap,
	})
}

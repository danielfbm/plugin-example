package main

import (
	"os"

	ext "github.com/danielfbm/plugin-example/controller/extension"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type FooPlugin struct {
	logger hclog.Logger
}

func (b *FooPlugin) Foos() (string, error) {
	b.logger.Debug("FooPlugin.Foos: called")
	return "Foos!", nil
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	plug := &FooPlugin{logger: logger}

	var pluginMap = map[string]plugin.Plugin{
		"foo": &ext.FooPlugin{Impl: plug},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: ext.HandshakeConfig,
		Plugins:         pluginMap,
	})
}

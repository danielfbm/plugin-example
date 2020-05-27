package main

import (
	"os"

	ext "github.com/danielfbm/plugin-example/controller/extension"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

type BarPlugin struct {
	logger hclog.Logger
}

func (b *BarPlugin) Bars() []string {
	b.logger.Debug("BaarPlugin.Bars: called")
	return []string{}
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	plug := &BarPlugin{logger: logger}

	var pluginMap = map[string]plugin.Plugin{
		"bar": &ext.BarPlugin{Impl: plug},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: ext.HandshakeConfig,
		Plugins:         pluginMap,
	})
}

package main

import (
	"fmt"
	"os"

	example "github.com/danielfbm/plugin-example/demo/commons"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of Greeter
type GreeterHello struct {
	logger hclog.Logger
}

func (g *GreeterHello) Greet(name string) string {
	g.logger.Debug("message from zh - GreeterHello.Greet")
	return fmt.Sprintf("您好 %s", name)
}

func (g *GreeterHello) Hi() (string, error) {
	return "大家好！", nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	greeter := &GreeterHello{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"greeter": &example.GreeterPlugin{Impl: greeter},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}

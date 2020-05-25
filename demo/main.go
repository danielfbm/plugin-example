package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	example "github.com/danielfbm/plugin-example/demo/commons"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
)

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	plugins, err := loadPlugins(logger)
	if err != nil {
		log.Fatal(err)
	}

	for _, client := range plugins {
		defer client.Kill()
		// Connect via RPC
		rpcClient, err := client.Client()
		if err != nil {
			log.Fatal(err)
		}

		// Request the plugin
		raw, err := rpcClient.Dispense("greeter")
		if err != nil {
			log.Fatal(err)
		}

		// We should have a Greeter now! This feels like a normal interface
		// implementation but is in fact over an RPC connection.
		greeter := raw.(example.Greeter)
		fmt.Println(greeter.Greet("some name here"))

		fmt.Println(greeter.Hi(2))

		fmt.Println("Will try pingponger...")
		raw2, err2 := rpcClient.Dispense("pingponger")
		if err2 != nil {
			fmt.Println("err", err2)
			continue
		}

		pinger := raw2.(example.PingPonger)
		pong, err2 := pinger.Ping()
		fmt.Println("ping?", pong, "err", err2)
	}
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

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"greeter":    &example.GreeterPlugin{},
	"pingponger": &example.PingPongerPlugin{},
}

func loadPlugins(logger hclog.Logger) (plugins []*plugin.Client, err error) {
	var found []string
	found, err = plugin.Discover("*.po", "./bin")

	fmt.Println("found", found, "err", err)
	plugins = make([]*plugin.Client, 0, len(found))

	for _, f := range found {
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         pluginMap,
			Cmd:             exec.Command(f),
			Logger:          logger,
		})
		plugins = append(plugins, client)
	}

	return
}

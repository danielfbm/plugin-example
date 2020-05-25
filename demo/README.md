# Plugin demo


## Intro

1. Create a hello world plugin with a name as parameter, implement two versions of the plugin, load them dynamically
2. Make a change to the interface, add a new method, compile the main program and run
3. Make a breaking change to the interface and try the step 2


### Basic plugin struct

Inside a new `demo` folder:

1. copy basic example from plugin [github.com/hashicorp/go-plugin/examples/basic]
2. add a parameter name and printing method displaying name
3. implement another plugin with a different response
4. change the main app to load plugins from file


### Run

Compile the plugin itself via:

    go build -o ./plugin/greeter ./plugin/greeter_impl.go

Compile this driver via:

    go build -o basic .

You can then launch the plugin sample via:

    ./basic


### Problems to solve:

#### 1. Can a change in interface on the host adding a new method utilizes an old plugin?

1. Add a Hi() string, error to the interface `commons/greeter_interface.go`
2. Compile basic again and run


*Result: *Yes, it works when adding new methods, but if the new method is invoked will return an error:

```
2020-05-25T10:48:14.797+0800 [DEBUG] plugin.greeter: message from GreeterHello.Greet: timestamp=2020-05-25T10:48:14.796+0800
Hello someone
 rpc: can't find method Plugin.Hi
```


3. If the `ProtocolVersion` changes,increase a number

*Result: *Returns a new error and any calls to the plugin will fail

```
2020-05-25T10:49:37.917+0800 [DEBUG] plugin.greeter: message from plugin: foo=bar timestamp=2020-05-25T10:49:37.917+0800
2020/05/25 10:49:37 Incompatible API version with plugin. Plugin version: 1, Client versions: [2]
```

#### 2. How to dynamically call different implementations?

1. Copy the `plugin` folder as `pluginzh` folder
2. Implement the new method Hi
3. Change the implementation to anything else (change the returned message)
4. Create a method to load multiple `*plugin.Client`s
5. Use `plugin.Discover` method to discover plugins
6. Move all plugins to a folder using a "file format" .po: `bin/greeter-en.po` `bin/greeter-zh.po`
7. Change the main method to load and manager multiple clients and make calls

```

func main() {
	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// We're a host! Start by launching the plugin process.
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
		fmt.Println(greeter.Greet("someone"))

		fmt.Println(greeter.Hi())
	}
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
	"greeter": &example.GreeterPlugin{},
}
```

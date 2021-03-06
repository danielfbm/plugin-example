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

**Result:** Yes, it works when adding new methods, but if the new method is invoked will return an error:

```
2020-05-25T10:48:14.797+0800 [DEBUG] plugin.greeter: message from GreeterHello.Greet: timestamp=2020-05-25T10:48:14.796+0800
Hello someone
 rpc: can't find method Plugin.Hi
```


3. If the `ProtocolVersion` changes,increase a number

**Result: **Returns a new error and any calls to the plugin will fail

```
2020-05-25T10:49:37.917+0800 [DEBUG] plugin.greeter: message from plugin: foo=bar timestamp=2020-05-25T10:49:37.917+0800
2020/05/25 10:49:37 Incompatible API version with plugin. Plugin version: 1, Client versions: [2]
```


#### 2. Can a breaking change in the interface cause any trouble? (without changing the version)

1. Add a parameter to  `Hi` method `Hi(int) (string, error)`
2. Remove the parameter from `Greet` method `Greet() string`
3. Build and try again

**Result:** 

- Step 1 can work successfully, will use a default value (zero value)
- Step 2 will break the plugin and will not work any calls.

If during the call the `host` uses a default value then it will work:

```
func (g *GreeterRPC) Greet() string {
	var resp string
	err := g.client.Call("Plugin.Greet", map[string]interface{}{
		"name": "default?",
	}, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}
```

#### 3. How to dynamically call different implementations?

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

**Result:** It is possible but needs to add a new layer to manage clients


#### 4. Can one plugin implement multiple interfaces?

1. Add another interface `commons/pingpong_interface.go`

```
// PingPonger is the interface that we're exposing as a plugin.
type PingPonger interface {
	Ping() (string, error)
}
```

2. Implement RPC server and client on interface
3. Add interface methods to `pluginzh/greeter_impl.go`

```
// Ping adds implementation for PingPonger plugin
func (g *GreeterHello) Ping() (string, error) {
	return "pong!", nil
}

//[...]

func main() {
    //[...]
    // pluginMap is the map of plugins we can dispense.
    var pluginMap = map[string]plugin.Plugin{
        "greeter":    &example.GreeterPlugin{Impl: greeter},
        "pingponger": &example.PingPongerPlugin{Impl: greeter},
    }
    //[...]
}
```
4. Add new interface to `main.go`


```
//[...]
func main() {
    // [...]

    for _, client := range plugins {
        //[...]
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

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"greeter":    &example.GreeterPlugin{},
	"pingponger": &example.PingPongerPlugin{},
}
```

**Result:** It is possible, the host needs to add this kind of logic to try multiple plugins over the same client

#### 5. Can a host use plugins in the network (localhost/kubernetes)?

[Hashicorp's go-plugin](https://github.com/hashicorp/go-plugin) does not support connection over the network on linux systems, and only uses this kind of method on windows but binding to 127.0.0.1 making it impossible to access it in the network.

[I forked the repo](https://github.com/danielfbm/go-plugin) and changed few lines of code to support starting a plugin using tcp bindings:

`server.go`:

```
func serverListener() (net.Listener, error) {
	if runtime.GOOS == "windows" || os.Getenv("PLUGIN_HOST") != "" {
		return serverListener_tcp()
	}

	return serverListener_unix()
}

// and
func serverListener_tcp() (net.Listener, error) {
	//...
	for port := minPort; port <= maxPort; port++ {
		address := fmt.Sprintf("0.0.0.0:%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			return listener, nil
		}
	}

	return nil, errors.New("Couldn't bind plugin TCP listener")
}

```



##### 5.1 localhost

1. Change the `go.mod` to replace the module
2. Change the `main.go` to load a plugin over the network

Inside the `loadPlugins` function:

```
/...
	if os.Getenv("NETWORK_PLUGIN") != "" {
		var netTCP net.Addr
		if netTCP, err = net.ResolveTCPAddr("tcp", os.Getenv("NETWORK_PLUGIN")); err != nil {
			return
		}
		client := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: handshakeConfig,
			Plugins:         pluginMap,
			Reattach: &plugin.ReattachConfig{
				Protocol: plugin.ProtocolNetRPC,
				Addr:     netTCP,
			},
			Logger: logger,
		})
		plugins = append(plugins, client)
	}
/...
```

3. Recompile `pluginzh` using new version, start using the keywork:

```
BASIC_PLUGIN=hello PLUGIN_HOST=true ./bin/greeter-zh.po
```

4. Recompile `basic` and start using the `NETWORK_PLUGIN` envvar. *PS: the port needs to be fetched from the plugin log*

```
NETWORK_PLUGIN=127.0.0.1:57674 ./basic
```

**Result:** Yes, but with a custom fork, but the client should be managed separatedly.

##### 5.2 kubernetes

5. Add 2 [`Dockerfiles`](Dockerfile), one for basic and another for the [`pluginzh`](pluginzh/Dockerfile)
6. Build and push to any registry
7. Create and expose the plugin

```
kubectl run pluginzh --image=danielfbm/pluginzh --env=BASIC_PLUGIN=hello --port=7000
kubectl expose pod pluginzh --port=7000 --target-port=7000
```

8. Run the basic

```
kubectl run basic --image=danielfbm/basic --env=NETWORK_PLUGIN=pluginzh:7000
```

**PS: When using this approach the host should not kill plugins**

*Result: *Yes
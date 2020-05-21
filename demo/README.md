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

package example

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Greeter is the interface that we're exposing as a plugin.
type Greeter interface {
	Greet(name string) string
	Hi() (string, error)
}

// Here is an implementation that talks over RPC
type GreeterRPC struct{ client *rpc.Client }

func (g *GreeterRPC) Greet(name string) string {
	var resp string
	err := g.client.Call("Plugin.Greet", map[string]interface{}{
		"name": name,
	}, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

func (g *GreeterRPC) Hi() (string, error) {
	var resp string
	err := g.client.Call("Plugin.Hi", map[string]interface{}{}, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		return "", err
	}

	return resp, nil
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type GreeterRPCServer struct {
	// This is the real implementation
	Impl Greeter
}

func (s *GreeterRPCServer) Greet(args map[string]interface{}, resp *string) error {
	*resp = s.Impl.Greet(args["name"].(string))
	return nil
}

func (s *GreeterRPCServer) Hi(args map[string]interface{}, resp *string) error {
	var err error
	*resp, err = s.Impl.Hi()
	return err
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type GreeterPlugin struct {
	// Impl Injection
	Impl Greeter
}

func (p *GreeterPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &GreeterRPCServer{Impl: p.Impl}, nil
}

func (GreeterPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &GreeterRPC{client: c}, nil
}
